package land

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

const vnmlSize = int(65)

// VertexField represents a 3-component int8 vertex
type VertexField struct {
	X int8
	Y int8
	Z int8
}

// Vertex Normals. A 65×65 array of: int8 - X, int8 - Y, int8 - Z.
// Note that the Y-direction of the data is from the bottom up.
const VNML = esm.SubrecordTag("VNML")

// Heights for world map. Derived from VHGT data.
type VNMLField struct {
	Vertices [][]VertexField
}

func (s *VNMLField) Tag() esm.SubrecordTag { return VNML }

func (s *VNMLField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	vertexSlice, err := util.SliceFromBytes[VertexField](vnmlSize*vnmlSize, sub.Data)
	if err != nil {
		return fmt.Errorf("slice from bytes: %w", err)
	}
	s.Vertices, err = util.SliceAsGrid(vnmlSize, vertexSlice)
	if err != nil {
		return fmt.Errorf("slice as grid: %w", err)
	}
	return nil
}

func (s *VNMLField) Marshal() (*esm.Subrecord, error) {
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
