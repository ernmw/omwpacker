package cell

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
)

// These only appear in mod files when creatures or NPCs are moved from one cell to another; they commonly appear in saved game files as things move around.
type MoveReference struct {
	// Reference ID (always the same as the attached FRMR value).
	// Required.
	MVRF *MVRFField
	// Name of the cell the reference was moved to (interior cells only)
	// zstring
	// Optional.
	CNAM *CNAMField
	// Coordinates of the cell the reference was moved to (exterior cells only)
	//   int32 - Grid X
	//   int32 - Grid Y
	// Optional.
	CNDT *CNDTField
	// Reference to the form that was moved.
	// Optional.
	Moved *FormReference
}

func (m *MoveReference) OrderedRecords() ([]*esm.Subrecord, error) {
	if m == nil {
		return nil, nil
	}
	orderedSubrecords := []*esm.Subrecord{}
	add := func(p esm.ParsedSubrecord) error {
		if p != nil {
			subRec, err := p.Marshal()
			if err != nil {
				return fmt.Errorf("marshal %q to subrec", p.Tag())
			}
			if subRec != nil {
				orderedSubrecords = append(orderedSubrecords, subRec)
			}
		}
		return nil
	}

	if err := add(m.MVRF); err != nil {
		return nil, err
	}
	if err := add(m.CNAM); err != nil {
		return nil, err
	}
	if err := add(m.CNDT); err != nil {
		return nil, err
	}

	if m.Moved != nil {
		if recs, err := m.Moved.OrderedRecords(); err != nil {
			return nil, err
		} else {
			orderedSubrecords = append(orderedSubrecords, recs...)
		}
	}

	return orderedSubrecords, nil
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
			mr.MVRF = &MVRFField{}
			if err := mr.MVRF.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case CNAM:
			mr.CNAM = &CNAMField{}
			if err := mr.CNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case CNDT:
			mr.CNDT = &CNDTField{}
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
