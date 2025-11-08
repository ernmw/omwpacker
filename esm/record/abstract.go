package record

import (
	"bytes"
	"encoding/binary"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

var _ esm.ParsedSubrecord = (*AnonymousZstringSubrecord)(nil)
var _ esm.ParsedSubrecord = (*AnonymousFloat32Subrecord)(nil)
var _ esm.ParsedSubrecord = (*AnonymousUint32Subrecord)(nil)

type AnonymousZstringSubrecord struct {
	Value       string
	EmbeddedTag esm.SubrecordTag
}

func (s *AnonymousZstringSubrecord) Tag() esm.SubrecordTag {
	return s.EmbeddedTag
}

func (s *AnonymousZstringSubrecord) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Value = string(sub.Data)
	return nil
}

func (s *AnonymousZstringSubrecord) Marshal() (*esm.Subrecord, error) {
	return &esm.Subrecord{Tag: s.Tag(), Data: []byte(s.Value)}, nil
}

type AnonymousFloat32Subrecord struct {
	Value       float32
	EmbeddedTag esm.SubrecordTag
}

func (s *AnonymousFloat32Subrecord) Tag() esm.SubrecordTag {
	return s.EmbeddedTag
}

func (s *AnonymousFloat32Subrecord) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Value = util.BytesToFloat32(sub.Data[0:4])
	return nil
}

func (s *AnonymousFloat32Subrecord) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, s.Value); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: buff.Bytes()}, nil
}

type AnonymousUint32Subrecord struct {
	Value       uint32
	EmbeddedTag esm.SubrecordTag
}

func (s *AnonymousUint32Subrecord) Tag() esm.SubrecordTag {
	return s.EmbeddedTag
}

func (s *AnonymousUint32Subrecord) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Value = binary.LittleEndian.Uint32(sub.Data[0:4])
	return nil
}

func (s *AnonymousUint32Subrecord) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, s.Value); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: buff.Bytes()}, nil
}

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
