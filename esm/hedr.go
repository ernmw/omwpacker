package esm

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ernmw/omwpacker/esm/tags"
)

type HEDRdata struct {
	Version     float32
	Flags       uint32
	Name        string
	Description string
	NumRecords  uint32
}

func (h *HEDRdata) Tag() tags.SubrecordTag {
	return tags.HEDR
}

func (h *HEDRdata) Unmarshal(sub *Subrecord) error {
	if h == nil || sub == nil {
		return ErrArgumentNil
	}
	// require full HEDR payload size (300 bytes)
	if len(sub.Data) < 300 {
		return fmt.Errorf("%q subrecord too short: %d < 300", h.Tag(), len(sub.Data))
	}
	h.Version = bytesToFloat32(sub.Data[0:4])
	h.Flags = binary.LittleEndian.Uint32(sub.Data[4:8])
	h.Name = readPaddedString(sub.Data[8 : 8+32])
	h.Description = readPaddedString(sub.Data[8+32 : 8+32+256])
	h.NumRecords = binary.LittleEndian.Uint32(sub.Data[8+32+256 : 8+32+256+4])
	return nil
}

func (h *HEDRdata) Marshal() (*Subrecord, error) {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, h.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.LittleEndian, h.Flags); err != nil {
		return nil, err
	}
	if err := writePaddedString(buff, []byte(h.Name), 32); err != nil {
		return nil, err
	}
	if err := writePaddedString(buff, []byte(h.Description), 256); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.LittleEndian, h.NumRecords); err != nil {
		return nil, err
	}
	return &Subrecord{Tag: h.Tag(), Data: buff.Bytes()}, nil
}
