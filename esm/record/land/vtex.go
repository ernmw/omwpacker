package land

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

const vtexSize = int(16)

// Vertex Normals. A 65×65 array of: int8 - X, int8 - Y, int8 - Z.
// Note that the Y-direction of the data is from the bottom up.
const VTEX = esm.SubrecordTag("VTEX")

// Heights for world map. Derived from VHGT data.
type VTEXField struct {
	Vertices [][]uint16
}

func (s *VTEXField) Tag() esm.SubrecordTag { return VTEX }

func (s *VTEXField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	vertexSlice, err := util.SliceFromBytes[uint16](vtexSize*vtexSize, sub.Data)
	if err != nil {
		return fmt.Errorf("slice from bytes: %w", err)
	}
	s.Vertices, err = util.SliceAsGrid(vtexSize, vertexSlice)
	if err != nil {
		return fmt.Errorf("slice as grid: %w", err)
	}
	return nil
}

func (s *VTEXField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}
	vertexSlice, err := util.GridAsSlice(s.Vertices)
	if err != nil {
		return nil, fmt.Errorf("grid as slice: %w", err)
	}
	outData, err := util.BytesFromSlice(vertexSlice)
	if err != nil {
		return nil, fmt.Errorf("bytes from slice: %w", err)
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: outData}, nil
}
