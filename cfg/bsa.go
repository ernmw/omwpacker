package cfg

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func (e *Environment) ReadFile(path string) ([]byte, error) {
	// load data folders first
	for _, dataFolder := range slices.Backward(e.Data) {
		raw, err := os.ReadFile(filepath.Join(dataFolder, path))
		if err != nil {
			continue
		}
		return raw, nil
	}
	// then BSAs
	for _, bsaFile := range slices.Backward(e.BSA) {
		raw, err := e.extractFile(bsaFile, path)
		if err != nil {
			continue
		}
		return raw, nil
	}
	return nil, fmt.Errorf("%q not found", path)
}

func (e *Environment) cachedEntries(bsaFile string) ([]*entry, error) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if entries, ok := e.bsaIndices[bsaFile]; ok {
		return entries, nil
	}

	f, err := os.Open(bsaFile)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", bsaFile, err)
	}
	defer f.Close()

	indices, err := parseTES3Index(f)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", bsaFile, err)
	}
	e.bsaIndices[bsaFile] = indices
	return indices, nil
}

func (e *Environment) extractFile(bsaFile, path string) ([]byte, error) {
	indices, err := e.cachedEntries(bsaFile)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", bsaFile, err)
	}
	for _, idx := range indices {
		if idx.Name == path {
			f, err := os.Open(bsaFile)
			if err != nil {
				return nil, fmt.Errorf("open %q: %w", bsaFile, err)
			}
			defer f.Close()
			out := make([]byte, idx.Size)
			readCount, err := f.ReadAt(out, int64(idx.Offset))
			if err != nil {
				return nil, fmt.Errorf("read %d to %d: %w", idx.Offset, idx.Offset+idx.Size, err)
			}
			if readCount != int(idx.Size) {
				return nil, fmt.Errorf("expected %d bytes, got %d", idx.Size, readCount)
			}
			return out, nil
		}
	}
	return nil, fmt.Errorf("%q not found in %q", path, bsaFile)
}

// entry represents one file inside a TES3 (Morrowind) BSA.
// Offset is an absolute offset from the start of the archive file.
type entry struct {
	Name   string
	Size   uint32
	Offset uint32 // absolute offset in file (safe to cast to int64 for Seek)
}

// parseTES3Index reads the index from a Morrowind (TES3) BSA (r must be seekable).
// It returns the list of entries with absolute offsets (ready to be read).
func parseTES3Index(r io.ReadSeeker) ([]*entry, error) {
	// Save starting position and determine archive length
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek start: %w", err)
	}
	// get file length
	end, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("seek end: %w", err)
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek start 2: %w", err)
	}

	var magic uint32
	if err := binary.Read(r, binary.LittleEndian, &magic); err != nil {
		return nil, fmt.Errorf("read magic: %w", err)
	}
	// magic must be 0x00000100 per TES3 spec
	const tes3Magic = 0x00000100
	if magic != tes3Magic {
		return nil, fmt.Errorf("not a TES3 BSA (magic=0x%08x)", magic)
	}

	var hashTableOffsetMinusHeader uint32
	var fileCount uint32
	if err := binary.Read(r, binary.LittleEndian, &hashTableOffsetMinusHeader); err != nil {
		return nil, fmt.Errorf("read hashTableOffsetMinusHeader: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &fileCount); err != nil {
		return nil, fmt.Errorf("read fileCount: %w", err)
	}

	// sanity checks
	if fileCount == 0 || fileCount > 200000 {
		return nil, fmt.Errorf("unreasonable fileCount %d", fileCount)
	}

	// Read file size/offset pairs (uint32 size, uint32 offset) count = fileCount
	type rawPair struct {
		Size   uint32
		Offset uint32 // this is "offset of the file in the data section" (relative; we'll convert later)
	}
	pairs := make([]rawPair, fileCount)
	for i := uint32(0); i < fileCount; i++ {
		if err := binary.Read(r, binary.LittleEndian, &pairs[i].Size); err != nil {
			return nil, fmt.Errorf("entry %d: read size: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &pairs[i].Offset); err != nil {
			return nil, fmt.Errorf("entry %d: read offset: %w", i, err)
		}
	}

	// Read name offsets (uint32[fileCount]) - offsets are relative to the start of the Names section,
	// which begins immediately after this array.
	nameOffsets := make([]uint32, fileCount)
	for i := uint32(0); i < fileCount; i++ {
		if err := binary.Read(r, binary.LittleEndian, &nameOffsets[i]); err != nil {
			return nil, fmt.Errorf("name offset %d: %w", i, err)
		}
	}

	// The current position is the start of the Names section.
	namesSectionStart, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("tell namesSectionStart: %w", err)
	}

	// Read each name by seeking to (namesSectionStart + nameOffsets[i]) and reading a null-terminated zstring.
	names := make([]string, fileCount)
	for i := uint32(0); i < fileCount; i++ {
		off := int64(namesSectionStart) + int64(nameOffsets[i])
		if off < 0 || off >= end {
			return nil, fmt.Errorf("name %d: invalid name offset %d (out of bounds)", i, nameOffsets[i])
		}
		if _, err := r.Seek(off, io.SeekStart); err != nil {
			return nil, fmt.Errorf("name %d: seek: %w", i, err)
		}

		// read bytes until 0
		var bbuf [1]byte
		var sb strings.Builder
		for {
			if _, err := r.Read(bbuf[:]); err != nil {
				return nil, fmt.Errorf("name %d: read byte: %w", i, err)
			}
			if bbuf[0] == 0 {
				break
			}
			// names are ASCII lowercase per spec
			sb.WriteByte(bbuf[0])
			// guard: if name grows absurdly long, abort
			if sb.Len() > 4096 {
				return nil, fmt.Errorf("name %d too long", i)
			}
		}
		names[i] = strings.ToLower(strings.ReplaceAll(sb.String(), "\\", "/"))
	}

	// Compute absolute locations:
	// The header's second field is "offset of the hash table in the file, minus the header size (12)".
	// So absoluteHashTableOffset = hashTableOffsetMinusHeader + 12.
	absHashTableOffset := int64(hashTableOffsetMinusHeader) + 12
	if absHashTableOffset < 0 || absHashTableOffset > end {
		return nil, fmt.Errorf("invalid hash table offset: %d", absHashTableOffset)
	}

	// Filename hashes occupy 8 * fileCount bytes starting at absHashTableOffset.
	hashesEnd := absHashTableOffset + int64(8*fileCount)
	if hashesEnd > end {
		return nil, fmt.Errorf("hash table exceeds file size")
	}

	// Raw data (the data section) starts immediately after the filename-hash table.
	dataSectionStart := hashesEnd

	// Build entries: size from pairs[i].Size, absolute offset = dataSectionStart + pairs[i].Offset
	entries := make([]*entry, 0, fileCount)
	for i := uint32(0); i < fileCount; i++ {
		size := pairs[i].Size
		relOff := pairs[i].Offset
		absOff := int64(dataSectionStart) + int64(relOff)

		// bounds checks
		if absOff < 0 || absOff > end {
			return nil, fmt.Errorf("entry %d: computed absolute offset out of bounds (%d)", i, absOff)
		}
		if int64(size) < 0 || absOff+int64(size) > end {
			return nil, fmt.Errorf("entry %d: file data out of bounds (off %d size %d)", i, absOff, size)
		}
		entries = append(entries, &entry{
			Name:   names[i],
			Size:   size,
			Offset: uint32(absOff),
		})
	}

	return entries, nil
}
