// Package esm deals with parsing ESM/omwaddon files.
package esm

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"
	"strings"
)

// RecordTag identifies the type of record.
// See https://github.com/OpenMW/openmw/blob/39d117e362808dc13cd411debcb48e363e11639c/components/esm/defs.hpp#L78
type RecordTag string
type SubrecordTag string

var ErrArgumentNil error
var ErrTagMismatch error

func newErrTagMismatch(expected SubrecordTag, got SubrecordTag) error {
	if expected != got {
		return fmt.Errorf("expected %q, got %q: %w", expected, got, ErrTagMismatch)
	}
	return nil
}

// ParsedSubrecord is an unmarshalled Subrecord.
type ParsedSubrecord interface {
	// Unmarshal sub into this instance.
	Unmarshal(sub *Subrecord) error
	// Marshal the parsed subrecord into a raw binary representation.
	// If this instance is nil, this should return (nil, nil).
	Marshal() (*Subrecord, error)
	Tag() SubrecordTag
}

// Subrecord is a marshalled component of a Record.
type Subrecord struct {
	Tag  SubrecordTag
	Data []byte
}

// Write the data in the subrecord to the writer w.
func (s *Subrecord) Write(w io.Writer) error {
	if _, err := w.Write([]byte(s.Tag)[0:4]); err != nil {
		return fmt.Errorf("write subrecord tag %q: %v", s.Tag, err)
	}
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(s.Data)))
	if _, err := w.Write(size); err != nil {
		return fmt.Errorf("write subrecord size %d: %v", len(s.Data), err)
	}
	if _, err := w.Write(s.Data); err != nil {
		return fmt.Errorf("write subrecord data: %v", err)
	}
	return nil
}

// UnmarshalTo the parsed subrecord p.
func (s *Subrecord) UnmarshalTo(p ParsedSubrecord) error {
	if s == nil {
		return ErrArgumentNil
	}
	if s.Tag != p.Tag() {
		return newErrTagMismatch(p.Tag(), s.Tag)
	}
	if err := p.Unmarshal(s); err != nil {
		return fmt.Errorf("unmarshal %q: %w", s.Tag, err)
	}
	return nil
}

// Record is an unmarshalled component of an ESM file.
type Record struct {
	// tag, size, padding, flags
	Tag        RecordTag
	Flags      uint32
	Subrecords []*Subrecord
	// PluginName is just metadata; it is not written to the file.
	PluginName string
	// PluginOffset is just metadata; it is not written to the file.
	PluginOffset int64
}

var padding = []byte{0, 0, 0, 0}

// Write the record to the writer w.
func (r *Record) Write(w io.Writer) error {
	// tag
	if _, err := w.Write([]byte(r.Tag)[0:4]); err != nil {
		return fmt.Errorf("write record tag %q: %v", r.Tag, err)
	}

	// brief intermission so we can determine size
	var buff bytes.Buffer
	for i, sub := range r.Subrecords {
		if err := sub.Write(&buff); err != nil {
			return fmt.Errorf("write %q subrecord %d: %v", r.Tag, i, err)
		}
	}

	// size
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(buff.Len()))
	if _, err := w.Write(size); err != nil {
		return fmt.Errorf("write %q record size %d: %v", r.Tag, buff.Len(), err)
	}

	// padding
	if _, err := w.Write(padding); err != nil {
		return fmt.Errorf("write record padding %q: %v", r.Tag, err)
	}

	// flags
	flags := make([]byte, 4)
	binary.LittleEndian.PutUint32(flags, r.Flags)
	if _, err := w.Write(flags); err != nil {
		return fmt.Errorf("write record %q flags %d: %v", r.Tag, r.Flags, err)
	}

	// now dump data out
	if _, err := io.Copy(w, &buff); err != nil {
		return fmt.Errorf("write subrecord data: %v", err)
	}
	return nil
}

func readUint32LE(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func readNextRecord(headerBuffer []byte, pluginName string, br io.Reader) (*Record, error) {
	n, err := io.ReadFull(br, headerBuffer)
	if err == io.EOF || (err == io.ErrUnexpectedEOF && n == 0) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tag := RecordTag(string(headerBuffer[0:4]))
	size := readUint32LE(headerBuffer[4:8])
	flags := readUint32LE(headerBuffer[12:16])

	body := make([]byte, size)
	if _, err := io.ReadFull(br, body); err != nil {
		return nil, fmt.Errorf("record %q: %w", tag, err)
	}

	rec := &Record{
		Tag:        tag,
		Flags:      flags,
		PluginName: pluginName,
		Subrecords: []*Subrecord{},
		// PluginOffset can be tracked externally if needed
	}

	// Parse the body buffer without more I/O
	pos := 0
	for pos < int(size) {
		if pos+8 > int(size) {
			return nil, fmt.Errorf("corrupt subrecord header in %q", tag)
		}

		subtag := SubrecordTag(string(body[pos : pos+4]))
		subsize := readUint32LE(body[pos+4 : pos+8])
		pos += 8

		if pos+int(subsize) > int(size) {
			return nil, fmt.Errorf("corrupt subrecord %q in %q", subtag, tag)
		}

		sr := &Subrecord{
			Tag:  subtag,
			Data: body[pos : pos+int(subsize)],
		}
		rec.Subrecords = append(rec.Subrecords, sr)
		pos += int(subsize)
	}

	return rec, nil
}

// ParsePluginFile extracts records from some esm or omwaddon file.
// See https://en.uesp.net/wiki/Morrowind_Mod:Mod_File_Format
func ParsePluginFile(path string) ([]*Record, error) {
	pluginName := strings.ToLower(filepath.Base(path))
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bufferedFile := bufio.NewReader(f)

	return ParsePluginData(pluginName, bufferedFile)
}

// ParsePluginData extracts records from an io.Reader.
func ParsePluginData(pluginName string, f io.Reader) ([]*Record, error) {
	records := []*Record{}
	hdr := make([]byte, 16)
	for {
		rec, err := readNextRecord(hdr, pluginName, f)
		if err != nil {
			return nil, err
		}
		if rec == nil {
			break
		}
		records = append(records, rec)
	}
	return records, nil
}

func WriteRecords(w io.Writer, recs iter.Seq[*Record]) error {
	for rec := range recs {
		if err := rec.Write(w); err != nil {
			return fmt.Errorf("write %q record: %w", rec.Tag, err)
		}
	}
	return nil
}
