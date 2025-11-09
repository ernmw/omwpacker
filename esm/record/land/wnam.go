package land

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
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
	if err := fillGridFromBytes(s.Heights, wnamSize, wnamSize, sub.Data); err != nil {
		return fmt.Errorf("parsing 2d array")
	}
	return nil
}

func (s *WNAMField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}

	return &esm.Subrecord{Tag: s.Tag(), Data: flattenGrid(s.Heights, wnamSize, wnamSize)}, nil
}
