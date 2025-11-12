package land

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

const wnamSize = 9

// Heights for world map. Derived from VHGT data.
const WNAM = esm.SubrecordTag("WNAM")

// Heights for world map. Derived from VHGT data.
type WNAMField struct {
	Heights [][]uint8
}

func (s *WNAMField) Tag() esm.SubrecordTag { return WNAM }

func (s *WNAMField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	var err error
	s.Heights, err = util.SliceAsGrid(wnamSize, sub.Data)
	if err != nil {
		return fmt.Errorf("slice as grid: %w", err)
	}
	return nil
}

func (s *WNAMField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}
	outData, err := util.GridAsSlice(s.Heights)
	if err != nil {
		return nil, fmt.Errorf("grid as slice: %w", err)
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: outData}, nil
}
