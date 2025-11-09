package cell

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ernmw/omwpacker/esm"
)

type CNDTdata struct {
	X int32
	Y int32
}

func (s *CNDTdata) Tag() esm.SubrecordTag {
	return "CNDT"
}

func (s *CNDTdata) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.X = int32(binary.LittleEndian.Uint32(sub.Data[0:4]))
	s.Y = int32(binary.LittleEndian.Uint32(sub.Data[4:8]))
	return nil
}

func (s *CNDTdata) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, s.X); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.LittleEndian, s.Y); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: buff.Bytes()}, nil
}

// These only appear in mod files when creatures or NPCs are moved from one cell to another; they commonly appear in saved game files as things move around.
type MoveReference struct {
	// Reference ID (always the same as the attached FRMR value).
	// Required.
	MVRF *MVRFdata
	// Name of the cell the reference was moved to (interior cells only)
	// zstring
	// Optional.
	CNAM *CNAMdata
	// Coordinates of the cell the reference was moved to (exterior cells only)
	//   int32 - Grid X
	//   int32 - Grid Y
	// Optional.
	CNDT *CNDTdata
	// Reference to the form that was moved.
	// Optional.
	Moved *FormReference
}

// returns formref + how many records it ate
func ParseMoveRef(subs []*esm.Subrecord) (*MoveReference, int, error) {
	if subs == nil {
		return nil, 0, esm.ErrArgumentNil
	}
	mr := &MoveReference{}
	processed := 0
subber:
	for i := 0; i < len(subs); i++ {
		sub := subs[i]
		switch sub.Tag {
		case MVRF:
			if mr.MVRF != nil {
				break subber
			}
			mr.MVRF = &MVRFdata{}
			if err := mr.MVRF.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case CNAM:
			mr.CNAM = &CNAMdata{}
			if err := mr.CNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case CNDT:
			mr.CNDT = &CNDTdata{}
			if err := mr.CNDT.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case FRMR:
			newFormRef, consumed, err := ParseFormRef(subs[i:])
			if err != nil {
				return nil, 0, fmt.Errorf("parse form reference: %w", err)
			}
			mr.Moved = newFormRef
			i = i + consumed
		default:
			break subber
		}
		processed++
	}
	return mr, processed, nil
}
