package land

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

const vclrSize = int(65)

// ColorField represents a 3-component int8 vertex
type ColorField struct {
	R uint8
	G uint8
	B uint8
}

// Vertex Normals. A 65×65 array of: int8 - X, int8 - Y, int8 - Z.
// Note that the Y-direction of the data is from the bottom up.
const VCLR = esm.SubrecordTag("VCLR")

// Heights for world map. Derived from VHGT data.
type VCLRField struct {
	Colors [][]ColorField
}

func (s *VCLRField) Tag() esm.SubrecordTag { return VCLR }

func (s *VCLRField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	colorSlice, err := util.SliceFromBytes[ColorField](vclrSize*vclrSize, sub.Data)
	if err != nil {
		return fmt.Errorf("slice from bytes: %w", err)
	}
	s.Colors, err = util.SliceAsGrid(vclrSize, colorSlice)
	if err != nil {
		return fmt.Errorf("slice as grid: %w", err)
	}
	return nil
}

func (s *VCLRField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}
	colorSlice, err := util.GridAsSlice(s.Colors)
	if err != nil {
		return nil, fmt.Errorf("grid as slice: %w", err)
	}
	outData, err := util.BytesFromSlice(colorSlice)
	if err != nil {
		return nil, fmt.Errorf("bytes from slice: %w", err)
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: outData}, nil
}
