package land

import (
	"fmt"

	"unsafe"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

const vclrSize = int(65)
const vclrDepth = int(3)

// ColorField represents a 3-component int8 vertex
type ColorField [vclrDepth]uint8

// Data returns a []byte slice pointing to the underlying bytes of the ColorField.
// Zero-allocation: no new slices are created.
func (v *ColorField) Data() []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(v)), len(v))
}

// ByteSize returns the number of bytes in the ColorField (always 3)
func (v *ColorField) ByteSize() int {
	return len(v)
}

func (v *ColorField) GetR() uint8 {
	return v[0]
}
func (v *ColorField) GetG() uint8 {
	return v[1]
}
func (v *ColorField) GetB() uint8 {
	return v[2]
}

// Vertex Normals. A 65×65 array of: int8 - X, int8 - Y, int8 - Z.
// Note that the Y-direction of the data is from the bottom up.
const VCLR = esm.SubrecordTag("VCLR")

// Heights for world map. Derived from VHGT data.
type VCLRField struct {
	Colors [][]*ColorField
}

func (s *VCLRField) Tag() esm.SubrecordTag { return VCLR }

func (s *VCLRField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if err := util.FillGridFromBytes(s.Colors, vclrSize, vclrSize, sub.Data); err != nil {
		return fmt.Errorf("parsing 2d array: %w", err)
	}
	return nil
}

func (s *VCLRField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}
	outBuff := make([]byte, vclrDepth*vclrSize*vclrSize)
	if err := util.FlattenGrid(s.Colors, vclrSize, vclrSize, outBuff); err != nil {
		return nil, fmt.Errorf("flatten grid: %w", err)
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: outBuff}, nil
}
