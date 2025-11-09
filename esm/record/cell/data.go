package cell

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ernmw/omwpacker/esm"
)

// DATA is a 12 byte struct containing flags and position.
const DATA esm.SubrecordTag = "DATA"

type DATAField struct {
	Flags uint32
	GridX int32
	GridY int32
}

func (s *DATAField) Tag() esm.SubrecordTag { return DATA }

func (s *DATAField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if len(sub.Data) < 12 {
		return fmt.Errorf("CELL.DATA too short: %d < 12", len(sub.Data))
	}
	s.Flags = binary.LittleEndian.Uint32(sub.Data[0:4])
	s.GridX = int32(binary.LittleEndian.Uint32(sub.Data[4:8]))
	s.GridY = int32(binary.LittleEndian.Uint32(sub.Data[8:12]))
	return nil
}

func (s *DATAField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, s.Flags)
	binary.Write(buf, binary.LittleEndian, s.GridX)
	binary.Write(buf, binary.LittleEndian, s.GridY)
	return &esm.Subrecord{Tag: s.Tag(), Data: buf.Bytes()}, nil
}
