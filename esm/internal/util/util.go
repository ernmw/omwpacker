package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

func WritePaddedString(out *bytes.Buffer, s []byte, size int) error {
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

func ReadPaddedString(raw []byte) string {
	if i := bytes.IndexByte(raw, 0); i >= 0 {
		return string(raw[:i])
	}
	return string(raw)
}

func BytesToFloat32(bytes []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(bytes))
}

func Float32ToBytes(float float32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, math.Float32bits(float))
	return bytes
}
