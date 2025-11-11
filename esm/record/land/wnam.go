package land

import (
	"fmt"

	"unsafe"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

const wnamSize = 9

type ByteField uint8

func (b *ByteField) Data() []byte {
	// Return a 1-byte slice pointing to the underlying ByteField.
	return unsafe.Slice((*byte)(b), 1)
}

func (b *ByteField) ByteSize() int {
	return 1
}

// Heights for world map. Derived from VHGT data.
const WNAM = esm.SubrecordTag("WNAM")

// Heights for world map. Derived from VHGT data.
type WNAMField struct {
	Heights [][]*ByteField
}

func (s *WNAMField) Tag() esm.SubrecordTag { return WNAM }

func (s *WNAMField) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if err := util.FillGridFromBytes(s.Heights, wnamSize, wnamSize, sub.Data); err != nil {
		return fmt.Errorf("parsing 2d array")
	}
	return nil
}

func (s *WNAMField) Marshal() (*esm.Subrecord, error) {
	if s == nil {
		return nil, nil
	}
	outBuff := make([]byte, 1*wnamSize*wnamSize)
	if err := util.FlattenGrid(s.Heights, wnamSize, wnamSize, outBuff); err != nil {
		return nil, fmt.Errorf("flatten grid: %w", err)
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: outBuff}, nil
}
