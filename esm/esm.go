// Package esm deals with parsing ESM/omwaddon files.
package esm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"iter"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/ernmw/omwpacker/esm/tags"
)

var ErrArgumentNil error
var ErrTagMismatch error

func newErrTagMismatch(expected tags.SubrecordTag, got tags.SubrecordTag) error {
	if expected != got {
		return fmt.Errorf("expected %q, got %q: %w", expected, got, ErrTagMismatch)
	}
	return nil
}

type ParsedSubrecord interface {
	Unmarshal(sub *Subrecord) error
	Marshal() (*Subrecord, error)
	Tag() tags.SubrecordTag
}

type Subrecord struct {
	Tag  tags.SubrecordTag
	Data []byte
}

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

func (s *Subrecord) Unmarshal(p ParsedSubrecord) error {
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

type Record struct {
	// tag, size, padding, flags
	Tag        tags.RecordTag
	Flags      uint32
	Subrecords []*Subrecord
	// PluginName is just metadata; it is not written to the file.
	PluginName string
	// PluginOffset is just metadata; it is not written to the file.
	PluginOffset int64
}

var padding = []byte{0, 0, 0, 0}

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

func readNextRecord(pluginName string, f io.ReadSeeker) (*Record, error) {
	start, _ := f.Seek(0, io.SeekCurrent)
	hdr := make([]byte, 16)
	n, err := io.ReadFull(f, hdr)
	if err == io.EOF || (err == io.ErrUnexpectedEOF && n == 0) {
		return nil, nil // end of file
	}
	if err != nil {
		return nil, err
	}
	rec := &Record{
		Subrecords:   []*Subrecord{},
		PluginOffset: start,
		PluginName:   pluginName,
	}
	rec.Tag = tags.RecordTag(string(hdr[0:4]))
	size := readUint32LE(hdr[4:8])
	// hdr[8:12] are padding
	rec.Flags = readUint32LE(hdr[12:16])

	limit := start + 16 + int64(size)
	for {
		pos, _ := f.Seek(0, io.SeekCurrent)
		if pos >= limit {
			break
		}
		// read subrecord header
		subhdr := make([]byte, 8)
		if _, err := io.ReadFull(f, subhdr); err != nil {
			return nil, fmt.Errorf("read subrecord header for %q: %w", rec.Tag, err)
		}
		tag := tags.SubrecordTag(string(subhdr[0:4]))
		size := readUint32LE(subhdr[4:8])
		data := make([]byte, size)
		if size > 0 {
			if _, err := io.ReadFull(f, data); err != nil {
				return nil, fmt.Errorf("read subrecord %q data (size %d) for %q: %w", tag, size, rec.Tag, err)
			}
		}
		rec.Subrecords = append(rec.Subrecords, &Subrecord{Tag: tag, Data: data})
	}
	return rec, nil
}

// GetSubrecord returns the first subrecord with the matching tag.
func (r *Record) GetSubrecord(tag tags.SubrecordTag) *Subrecord {
	for _, s := range r.Subrecords {
		if s.Tag == tag {
			return s
		}
	}
	return nil
}

// UpsertSubrecord replaces or inserts the subrecord.
func (r *Record) UpsertSubrecord(s *Subrecord) {
	// replace first occurrence or append
	for i, ex := range r.Subrecords {
		if ex.Tag == s.Tag {
			r.Subrecords[i] = s
			return
		}
	}
	r.Subrecords = append(r.Subrecords, s)
}

// DeleteSubrecord deletes the first subrecord with the matching tag.
func (r *Record) DeleteSubrecord(tag tags.SubrecordTag) {
	for i, ex := range r.Subrecords {
		if ex.Tag == tag {
			r.Subrecords = append(r.Subrecords[:i], r.Subrecords[i+1:]...)
			return
		}
	}
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

	return ParsePluginData(pluginName, f)
}

// ParsePluginData extracts records from an io.ReadSeeker.
func ParsePluginData(pluginName string, f io.ReadSeeker) ([]*Record, error) {
	records := []*Record{}
	for {
		rec, err := readNextRecord(pluginName, f)
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

func writePaddedString(out *bytes.Buffer, s []byte, size int) error {
	if len(s) > size {
		return errors.New("string too big")
	}
	if _, err := out.Write(s); err != nil {
		return err
	}
	if _, err := out.Write(make([]byte, size-len(s))); err != nil {
		return err
	}
	return nil
}

func readPaddedString(raw []byte) string {
	if i := bytes.IndexByte(raw, 0); i >= 0 {
		return string(raw[:i])
	}
	return string(raw)
}

func bytesToFloat32(bytes []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(bytes))
}

func float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func WriteRecords(w io.Writer, recs iter.Seq[*Record]) error {
	for rec := range recs {
		if err := rec.Write(w); err != nil {
			return fmt.Errorf("write %q record: %w", rec.Tag, err)
		}
	}
	return nil
}

func NewTES3Record(name string, description string) (*Record, error) {
	// make new empty records
	hedr := &HEDRdata{
		Version:     1.3,
		Flags:       0,
		Name:        name,
		Description: description,
		NumRecords:  0,
	}
	hedrSubRec, err := hedr.Marshal()
	if err != nil {
		return nil, fmt.Errorf("Write HEDR subrecord: %w", err)
	}
	return &Record{
		Tag: tags.TES3,
		Subrecords: []*Subrecord{
			hedrSubRec,
		},
	}, nil
}
