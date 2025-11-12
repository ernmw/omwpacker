package land

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

const vhgtSize = int(65)

// Height data.
const VHGT = esm.SubrecordTag("VHGT")

// Height Data.
type VHGTField struct {
	// A height offset for the entire cell.
	// Decreasing this value will shift the entire cell land down (by 8 units).
	Offset float32
	// Contains the height data for the cell in the form of a 65×65 pixel array.
	// The height data is not absolute values but uses differences between adjacent pixels.
	// Thus a pixel value of 0 means it has the same height as the last pixel.
	// Note that the Y-direction of the data is from the bottom up.
	Heights [][]uint8
}

func (s *VHGTField) Tag() esm.SubrecordTag { return VHGT }

func (s *VHGTField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Offset = util.BytesToFloat32(sub.Data[0:4])
	var err error
	s.Heights, err = util.SliceAsGrid(vhgtSize, sub.Data[4:len(sub.Data)-3])
	if err != nil {
		return fmt.Errorf("slice as grid: %w", err)
	}
	return nil
}

func (s *VHGTField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}
	gridSize := 1 * vhgtSize * vhgtSize
	outBuff := make([]byte, gridSize+4+3)

	copy(outBuff, util.Float32ToBytes(s.Offset))

	outData, err := util.GridAsSlice(s.Heights)
	if err != nil {
		return nil, fmt.Errorf("grid as slice: %w", err)
	}
	copy(outBuff[1:], outData[:])

	// last 3 bytes are junk

	return &esm.Subrecord{Tag: s.Tag(), Data: outBuff}, nil
}

// LandHeightScale is the factor OpenMW uses to scale VHGT heights.
// (From OpenMW’s Land::sHeightScale = 8.0f)
const LandHeightScale = 8.0

// ComputeAbsoluteHeights reconstructs the absolute height map from the
// differential VHGT data and returns a 65×65 [][]float32 slice.
// Each value already includes the Offset and height deltas.
//
// The resulting matrix follows the same bottom-up Y order as the
// stored VHGT data, i.e. [0][0] corresponds to the bottom-left corner.
func (s *VHGTField) ComputeAbsoluteHeights() [][]float32 {
	if s == nil || len(s.Heights) == 0 {
		return nil
	}

	heights := make([][]float32, vhgtSize)
	for i := range heights {
		heights[i] = make([]float32, vhgtSize)
	}

	// Start from the offset
	rowOffset := s.Offset

	for y := range vhgtSize {
		// First column of each row
		rowOffset += float32(s.Heights[y][0])
		colOffset := rowOffset

		h := colOffset * LandHeightScale
		heights[y][0] = h

		// Remaining columns in row
		for x := 1; x < vhgtSize; x++ {
			colOffset += float32(s.Heights[y][x])
			h := colOffset * LandHeightScale
			heights[y][x] = h
		}
	}

	return heights
}
