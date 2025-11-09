package land

import (
	"fmt"

	"unsafe"

	"github.com/ernmw/omwpacker/esm"
)

const vnmlSize = int(65)
const vnmlDepth = int(3)

// VertexField represents a 3-component int8 vertex
type VertexField [vnmlDepth]int8

// Data returns a []byte slice pointing to the underlying bytes of the VertexField.
// Zero-allocation: no new slices are created.
func (v *VertexField) Data() []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(v)), len(v))
}

// ByteSize returns the number of bytes in the VertexField (always 3)
func (v *VertexField) ByteSize() int {
	return len(v)
}

func (v *VertexField) GetX() int8 {
	return v[0]
}
func (v *VertexField) GetY() int8 {
	return v[1]
}
func (v *VertexField) GetZ() int8 {
	return v[2]
}

// Vertex Normals. A 65×65 array of: int8 - X, int8 - Y, int8 - Z.
// Note that the Y-direction of the data is from the bottom up.
const VNML = esm.SubrecordTag("VNML")

// Heights for world map. Derived from VHGT data.
type VNMLField struct {
	Vertices [][]*VertexField
}

func (s *VNMLField) Tag() esm.SubrecordTag { return VNML }

func (s *VNMLField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if err := FillGridFromBytes(s.Vertices, vnmlSize, vnmlSize, sub.Data); err != nil {
		return fmt.Errorf("parsing 2d array")
	}
	return nil
}

func (s *VNMLField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}
	outBuff := make([]byte, vnmlDepth*vnmlSize*vnmlSize)
	if err := FlattenGrid(s.Vertices, vnmlSize, vnmlSize, outBuff); err != nil {
		return nil, fmt.Errorf("flatten grid: %w", err)
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: outBuff}, nil
}
