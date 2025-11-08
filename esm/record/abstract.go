package record

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
)

type GenericSubrecord[T any] struct {
	Value T
	Tag   esm.SubrecordTag
}

func (s *GenericSubrecord[T]) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	switch t := (s.Value).(type) {
	case string:
		s.Value = string(sub.Data)
	}

	return s
}

func (s *GenericSubrecord[T]) Marshal() (*esm.Subrecord, error) {
	switch t := (s.Value).(type) {
	case string:
		return &esm.Subrecord{Tag: s.Tag, Data: s.Value}, nil
	}
	return nil, fmt.Errorf("unhandled type %T", s)
}
