package cell

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

// AMBI is a 16 byte struct for Ambient light (ambient, sunlight, fog colors, and fog density).
const AMBI esm.SubrecordTag = "AMBI"

type AMBIField struct {
	AmbientColor [3]uint8
	Sunlight     [3]uint8
	FogColor     [3]uint8
	FogDensity   float32
}

func (s *AMBIField) Tag() esm.SubrecordTag { return AMBI }

func (s *AMBIField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if len(sub.Data) < 16 {
		return fmt.Errorf("CELL.AMBI too short: %d < 16", len(sub.Data))
	}
	copy(s.AmbientColor[:], sub.Data[0:3]) // 4 is padding
	copy(s.Sunlight[:], sub.Data[4:7])     // 8 is padding
	copy(s.FogColor[:], sub.Data[8:11])    // 12 is padding
	s.FogDensity = util.BytesToFloat32(sub.Data[12:])
	return nil
}

func (s *AMBIField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}
	buf := new(bytes.Buffer)
	if _, err := buf.Write(s.AmbientColor[:]); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(0); err != nil {
		return nil, err
	}
	if _, err := buf.Write(s.Sunlight[:]); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(0); err != nil {
		return nil, err
	}
	if _, err := buf.Write(s.FogColor[:]); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(0); err != nil {
		return nil, err
	}
	binary.Write(buf, binary.LittleEndian, s.FogDensity)
	return &esm.Subrecord{Tag: s.Tag(), Data: buf.Bytes()}, nil
}
