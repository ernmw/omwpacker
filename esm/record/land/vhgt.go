package land

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

const vhgtSize = int(65)
const vhgtDepth = int(1)

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
	Heights [][]*ByteField
}

func (s *VHGTField) Tag() esm.SubrecordTag { return VHGT }

func (s *VHGTField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Offset = util.BytesToFloat32(sub.Data[0:4])
	if err := FillGridFromBytes(s.Heights, vhgtSize, vhgtSize, sub.Data[4:]); err != nil {
		return fmt.Errorf("parsing 2d array")
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

	if err := FlattenGrid(s.Heights, vhgtSize, vhgtSize, outBuff[4:gridSize+4]); err != nil {
		return nil, fmt.Errorf("flatten grid: %w", err)
	}

	// last 3 bytes are junk

	return &esm.Subrecord{Tag: s.Tag(), Data: outBuff}, nil
}
