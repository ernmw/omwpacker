package record

import (
	"bytes"
	"encoding/binary"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

// These exist to get compile-time type assurances and also get some
// code re-use between primitive subrecord types.
// It's kind of a pain to set them up, but using them is nice.

type Tagged interface {
	Tag() esm.SubrecordTag
}

type ZstringSubrecord[T Tagged] struct {
	Value string
}

func (s *ZstringSubrecord[T]) Tag() esm.SubrecordTag {
	var zero T
	return zero.Tag()
}

func (s *ZstringSubrecord[T]) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Value = string(sub.Data)
	return nil
}

func (s *ZstringSubrecord[T]) Marshal() (*esm.Subrecord, error) {
	return &esm.Subrecord{Tag: s.Tag(), Data: []byte(s.Value)}, nil
}

type Float32Subrecord[T Tagged] struct {
	Value float32
}

func (s *Float32Subrecord[T]) Tag() esm.SubrecordTag {
	var zero T
	return zero.Tag()
}

func (s *Float32Subrecord[T]) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Value = util.BytesToFloat32(sub.Data[0:4])
	return nil
}

func (s *Float32Subrecord[T]) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, s.Value); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: buff.Bytes()}, nil
}

type Uint32Subrecord[T Tagged] struct {
	Value uint32
}

func (s *Uint32Subrecord[T]) Tag() esm.SubrecordTag {
	var zero T
	return zero.Tag()
}

func (s *Uint32Subrecord[T]) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Value = binary.LittleEndian.Uint32(sub.Data[0:4])
	return nil
}

func (s *Uint32Subrecord[T]) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, s.Value); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: buff.Bytes()}, nil
}

type Uint8Subrecord[T Tagged] struct {
	Value uint8
}

func (s *Uint8Subrecord[T]) Tag() esm.SubrecordTag {
	var zero T
	return zero.Tag()
}

func (s *Uint8Subrecord[T]) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Value = sub.Data[0]
	return nil
}

func (s *Uint8Subrecord[T]) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, s.Value); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: []byte{s.Value}}, nil
}

type BytesSubrecord[T Tagged] struct {
	Value []byte
}

func (s *BytesSubrecord[T]) Tag() esm.SubrecordTag {
	var zero T
	return zero.Tag()
}

func (s *BytesSubrecord[T]) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Value = sub.Data
	return nil
}

func (s *BytesSubrecord[T]) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, s.Value); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: s.Value}, nil
}
