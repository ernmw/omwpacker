package esm

import (
	"bytes"
	"encoding/binary"
)

type Header struct {
	Version     float32
	Flags       uint32
	Name        string
	Description string
	NumRecords  uint32
}

func (h *Header) Unmarshal(raw []byte) error {
	h.Version = bytesToFloat32(raw[0:4])
	h.Flags = binary.LittleEndian.Uint32(raw[4:8])
	h.Name = readPaddedString(raw[8 : 8+32])
	h.Description = readPaddedString(raw[8+32 : 8+32+256])
	h.NumRecords = binary.LittleEndian.Uint32(raw[8+32+256 : 8+32+256+4])
	return nil
}

func (h *Header) Marshal() ([]byte, error) {
	var buff bytes.Buffer
	if _, err := buff.Write(float32ToBytes(h.Version)); err != nil {
		return nil, err
	}

	flags := make([]byte, 4)
	binary.LittleEndian.PutUint32(flags, h.Flags)
	if _, err := buff.Write(flags); err != nil {
		return nil, err
	}

	if err := writePaddedString(&buff, []byte(h.Name), 32); err != nil {
		return nil, err
	}

	if err := writePaddedString(&buff, []byte(h.Description), 256); err != nil {
		return nil, err
	}

	num := make([]byte, 4)
	binary.LittleEndian.PutUint32(num, h.NumRecords)
	if _, err := buff.Write(num); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
