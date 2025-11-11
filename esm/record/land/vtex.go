package land

import (
	"fmt"

	"unsafe"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

const vtexSize = int(16)

type UInt16Field uint16

func (u *UInt16Field) Data() []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(u)), 2)
}

func (u *UInt16Field) ByteSize() int {
	return 2
}

// Vertex Normals. A 65×65 array of: int8 - X, int8 - Y, int8 - Z.
// Note that the Y-direction of the data is from the bottom up.
const VTEX = esm.SubrecordTag("VTEX")

// Heights for world map. Derived from VHGT data.
type VTEXField struct {
	Vertices [][]*UInt16Field
}

func (s *VTEXField) Tag() esm.SubrecordTag { return VTEX }

func (s *VTEXField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if err := util.FillGridFromBytes(s.Vertices, vtexSize, vtexSize, sub.Data); err != nil {
		return fmt.Errorf("parsing 2d array: %w", err)
	}
	return nil
}

func (s *VTEXField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}
	outBuff := make([]byte, 2*vtexSize*vtexSize)
	if err := util.FlattenGrid(s.Vertices, vtexSize, vtexSize, outBuff); err != nil {
		return nil, fmt.Errorf("flatten grid: %w", err)
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: outBuff}, nil
}
