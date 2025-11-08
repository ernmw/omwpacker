package tes3

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

type HEDRdata struct {
	Version     float32
	Flags       uint32
	Name        string
	Description string
	NumRecords  uint32
}

func (h *HEDRdata) Tag() esm.SubrecordTag {
	return HEDR
}

func (h *HEDRdata) Unmarshal(sub *esm.Subrecord) error {
	if h == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	// require full HEDR payload size (300 bytes)
	if len(sub.Data) < 300 {
		return fmt.Errorf("%q subrecord too short: %d < 300", h.Tag(), len(sub.Data))
	}
	h.Version = util.BytesToFloat32(sub.Data[0:4])
	h.Flags = binary.LittleEndian.Uint32(sub.Data[4:8])
	h.Name = util.ReadPaddedString(sub.Data[8 : 8+32])
	h.Description = util.ReadPaddedString(sub.Data[8+32 : 8+32+256])
	h.NumRecords = binary.LittleEndian.Uint32(sub.Data[8+32+256 : 8+32+256+4])
	return nil
}

func (h *HEDRdata) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, h.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.LittleEndian, h.Flags); err != nil {
		return nil, err
	}
	if err := util.WritePaddedString(buff, []byte(h.Name), 32); err != nil {
		return nil, err
	}
	if err := util.WritePaddedString(buff, []byte(h.Description), 256); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.LittleEndian, h.NumRecords); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: h.Tag(), Data: buff.Bytes()}, nil
}
