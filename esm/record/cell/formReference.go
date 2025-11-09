package cell

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
)

// References to objects in cells are listed as part of the cell data, each beginning with FRMR and NAME fields, followed by a list of fields specific to the object type.
type FormReference struct {
	// Reference ID.
	// Type: uint32
	// Required.
	FRMR *FRMRField
	// Object ID or "PlayerSaveGame".
	// zstring
	// Required.
	NAME *NAMEField
	// Reference blocked (value is always 0; present if Blocked is set in the reference's record header, otherwise absent).
	// uint8
	// Optional.
	UNAM *UNAMField
	// Reference's scale, if applicable and not 1.0.
	// float32
	// Optional.
	XSCL *XSCLField
	// NPC ID, if applicable (NPC-only).
	// zstring
	// Optional, exclusive with BNAM.
	ANAM *ANAMField
	// Global variable name
	// zstring
	// Optional, exclusive with ANAM.
	BNAM *BNAMField
	// Faction ID (not light, NPC, or static)
	// zstring
	// Optional, if present then INDX must also exist.
	CNAM *CNAMField
	// Faction rank.
	// uint32
	// Optional, if present then CNAM must also exist.
	INDX *INDXField
	// ID of soul in gem (soul gems only)
	// zstring
	// Optional.
	XSOL *XSOLField
	// Enchantment charge (charged items with non-zero charges).
	// float32
	// Optional.
	XCHG *XCHGField
	// Depends on the object type.
	//   uint32 - health remaining (weapons and armor)
	//   uint32 - uses remaining (locks, probes, repair items)
	//   float32 - time remaining (lights)
	// Optional.
	INTV *INTVField
	// Value (in gold)
	// uint32
	// Optional.
	NAM9 *NAM9Field
	// Cell Travel Destination (Rotations are in radians)
	//   float32 - Position X
	//   float32 - Position Y
	//   float32 - Position Z
	//   //   float32 - Rotation X
	//   float32 - Rotation Y
	//   float32 - Rotation Z
	// Optional.
	DODT *DODTField // Cell name for previous DODT, if interior.
	// zstring
	// Optional, must accompany DODT if present.
	DNAM *DNAMField
	// Lock difficulty
	// uint32
	// Optional.
	FLTV *FLTVField
	// Key name
	// zstring
	// Optional.
	KNAM *KNAMField
	// Trap name
	// zstring
	// Optional.
	TNAM *TNAMField
	// Reference is disabled (always 0). Like UNAM, this will be emitted if the relevant flag is set in the reference's record header. This may only be possible via scripting. Also, even if present in the file, the field appears to be ignored on loading.
	// uint8
	// Optional.
	ZNAM *ZNAMField
	// Reference position (Rotations are in radians)
	//   float32 - Position X
	//   float32 - Position Y
	//   float32 - Position Z
	//   float32 - Rotation X
	//   float32 - Rotation Y
	//   float32 - Rotation Z
	// Optional.
	DATA *DATAFormReferenceField
}

func (f *FormReference) OrderedRecords() ([]*esm.Subrecord, error) {
	if f == nil {
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

	if err := add(f.FRMR); err != nil {
		return nil, err
	}
	if err := add(f.NAME); err != nil {
		return nil, err
	}
	if err := add(f.UNAM); err != nil {
		return nil, err
	}
	if err := add(f.XSCL); err != nil {
		return nil, err
	}
	if err := add(f.ANAM); err != nil {
		return nil, err
	}
	if err := add(f.BNAM); err != nil {
		return nil, err
	}
	if err := add(f.INDX); err != nil {
		return nil, err
	}
	if err := add(f.XSOL); err != nil {
		return nil, err
	}
	if err := add(f.XCHG); err != nil {
		return nil, err
	}
	if err := add(f.INTV); err != nil {
		return nil, err
	}
	if err := add(f.NAM9); err != nil {
		return nil, err
	}
	if err := add(f.DODT); err != nil {
		return nil, err
	}
	if err := add(f.DNAM); err != nil {
		return nil, err
	}
	if err := add(f.FLTV); err != nil {
		return nil, err
	}
	if err := add(f.KNAM); err != nil {
		return nil, err
	}
	if err := add(f.TNAM); err != nil {
		return nil, err
	}
	if err := add(f.ZNAM); err != nil {
		return nil, err
	}
	if err := add(f.DATA); err != nil {
		return nil, err
	}
	return orderedSubrecords, nil
}

// returns formref + how many records it ate
func ParseFormRef(subs []*esm.Subrecord) (*FormReference, int, error) {
	if subs == nil {
		return nil, 0, esm.ErrArgumentNil
	}
	fr := &FormReference{}
	processed := 0
subber:
	for i := 0; i < len(subs); i++ {
		sub := subs[i]
		switch sub.Tag {
		case FRMR:
			if fr.FRMR != nil {
				break subber
			}
			fr.FRMR = &FRMRField{}
			if err := fr.FRMR.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case NAME:
			fr.NAME = &NAMEField{}
			if err := fr.NAME.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case UNAM:
			fr.UNAM = &UNAMField{}
			if err := fr.UNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case XSCL:
			fr.XSCL = &XSCLField{}
			if err := fr.XSCL.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case ANAM:
			fr.ANAM = &ANAMField{}
			if err := fr.ANAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case BNAM:
			fr.BNAM = &BNAMField{}
			if err := fr.BNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case CNAM:
			fr.CNAM = &CNAMField{}
			if err := fr.CNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case INDX:
			fr.INDX = &INDXField{}
			if err := fr.INDX.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case XSOL:
			fr.XSOL = &XSOLField{}
			if err := fr.XSOL.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case XCHG:
			fr.XCHG = &XCHGField{}
			if err := fr.XCHG.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case INTV:
			fr.INTV = &INTVField{}
			if err := fr.INTV.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case NAM9:
			fr.NAM9 = &NAM9Field{}
			if err := fr.NAM9.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case DODT:
			fr.DODT = &DODTField{}
			if err := fr.DODT.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case DNAM:
			fr.DNAM = &DNAMField{}
			if err := fr.DNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case FLTV:
			fr.FLTV = &FLTVField{}
			if err := fr.FLTV.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case KNAM:
			fr.KNAM = &KNAMField{}
			if err := fr.KNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case TNAM:
			fr.TNAM = &TNAMField{}
			if err := fr.TNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case ZNAM:
			fr.ZNAM = &ZNAMField{}
			if err := fr.ZNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case DATAFormReference:
			fr.DATA = &DATAFormReferenceField{}
			if err := fr.DATA.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		default:
			break subber
		}
		processed++
	}
	return fr, processed, nil
}
