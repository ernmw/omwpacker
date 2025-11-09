//go:generate go run ../generator/gen.go subrecords.json
package cell

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
)

// CellRecord represents a full CellRecord record composed of subrecords.
type CellRecord struct {
	NAME               *NAMEField
	DATA               *DATAField
	RGNN               *RGNNField
	NAM5               *NAM5Field
	WHGT               *WHGTField
	AMBI               *AMBIField
	MovedReferences    []*MoveReference
	PersistentChildren []*FormReference
	// Count of temporaray children
	NAM0              *NAM0Field
	TemporaryChildren []*FormReference
}

func (c *CellRecord) OrderedRecords() ([]*esm.Subrecord, error) {
	if c == nil {
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

	if err := add(c.NAME); err != nil {
		return nil, err
	}
	if err := add(c.DATA); err != nil {
		return nil, err
	}
	if err := add(c.RGNN); err != nil {
		return nil, err
	}
	if err := add(c.NAM5); err != nil {
		return nil, err
	}
	if err := add(c.WHGT); err != nil {
		return nil, err
	}
	if err := add(c.AMBI); err != nil {
		return nil, err
	}

	for _, mr := range c.MovedReferences {
		if recs, err := mr.OrderedRecords(); err != nil {
			return nil, err
		} else {
			orderedSubrecords = append(orderedSubrecords, recs...)
		}
	}
	for _, fr := range c.PersistentChildren {
		if recs, err := fr.OrderedRecords(); err != nil {
			return nil, err
		} else {
			orderedSubrecords = append(orderedSubrecords, recs...)
		}
	}

	// deal with temp children in cell
	tempChildrenCount := uint32(len(c.TemporaryChildren))
	if tempChildrenCount > 0 {
		if c.NAM0 == nil {
			c.NAM0 = &NAM0Field{}
		}
		c.NAM0.Value = uint32(len(c.TemporaryChildren))
		if err := add(c.NAM0); err != nil {
			return nil, err
		}
		for _, fr := range c.TemporaryChildren {
			if recs, err := fr.OrderedRecords(); err != nil {
				return nil, err
			} else {
				orderedSubrecords = append(orderedSubrecords, recs...)
			}
		}
	} else {
		c.NAM0 = nil
	}

	return orderedSubrecords, nil
}

// ParseCELL builds a CELL record from a list of subrecords.
func ParseCELL(rec *esm.Record) (*CellRecord, error) {
	if rec == nil {
		return nil, esm.ErrArgumentNil
	}
	if rec.Tag != CELL {
		return nil, esm.ErrTagMismatch
	}
	c := &CellRecord{
		MovedReferences:    []*MoveReference{},
		PersistentChildren: []*FormReference{},
		TemporaryChildren:  []*FormReference{},
	}
	for i := 0; i < len(rec.Subrecords); i++ {
		sub := rec.Subrecords[i]
		switch sub.Tag {
		case NAME:
			c.NAME = &NAMEField{}
			if err := c.NAME.Unmarshal(sub); err != nil {
				return nil, err
			}
		case DATA:
			c.DATA = &DATAField{}
			if err := c.DATA.Unmarshal(sub); err != nil {
				return nil, err
			}
		case RGNN:
			c.RGNN = &RGNNField{}
			if err := c.RGNN.Unmarshal(sub); err != nil {
				return nil, err
			}
		case NAM5:
			c.NAM5 = &NAM5Field{}
			if err := c.NAM5.Unmarshal(sub); err != nil {
				return nil, err
			}
		case WHGT:
			c.WHGT = &WHGTField{}
			if err := c.WHGT.Unmarshal(sub); err != nil {
				return nil, err
			}
		case AMBI:
			c.AMBI = &AMBIField{}
			if err := c.AMBI.Unmarshal(sub); err != nil {
				return nil, err
			}
		case MVRF:
			newMoveRef, consumed, err := ParseMoveRef(rec.Subrecords[i:])
			if err != nil {
				return nil, fmt.Errorf("parse form reference: %w", err)
			}
			c.MovedReferences = append(c.MovedReferences, newMoveRef)
			i = i + consumed
		case FRMR:
			newFormRef, consumed, err := ParseFormRef(rec.Subrecords[i:])
			if err != nil {
				return nil, fmt.Errorf("parse form reference: %w", err)
			}
			if c.NAM0 != nil {
				c.TemporaryChildren = append(c.TemporaryChildren, newFormRef)
			} else {
				c.PersistentChildren = append(c.PersistentChildren, newFormRef)
			}
			i = i + consumed
		default:
			return nil, fmt.Errorf("unknown CELL subrecord %q", sub.Tag)
		}
	}
	return c, nil
}
