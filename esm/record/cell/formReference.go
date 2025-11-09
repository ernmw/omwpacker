package cell

import (
	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/record"
)

type frmrTagger struct{}

func (t *frmrTagger) Tag() esm.SubrecordTag { return "FRMR" }

type FRMRdata = record.Uint32Subrecord[*frmrTagger]

type anamTagger struct{}

func (t *anamTagger) Tag() esm.SubrecordTag { return "ANAM" }

type ANAMdata = record.ZstringSubrecord[*anamTagger]

type bnamTagger struct{}

func (t *bnamTagger) Tag() esm.SubrecordTag { return "BNAM" }

type BNAMdata = record.ZstringSubrecord[*bnamTagger]

type unamTagger struct{}

func (t *unamTagger) Tag() esm.SubrecordTag { return "UNAM" }

type UNAMdata = record.Uint8Subrecord[*unamTagger]

type xsclTagger struct{}

func (t *xsclTagger) Tag() esm.SubrecordTag { return "XSCL" }

type XSCLdata = record.Float32Subrecord[*xsclTagger]

type cnamTagger struct{}

func (t *cnamTagger) Tag() esm.SubrecordTag { return "CNAM" }

type CNAMdata = record.ZstringSubrecord[*cnamTagger]

type indxTagger struct{}

func (t *indxTagger) Tag() esm.SubrecordTag { return "INDX" }

type INDXdata = record.Uint32Subrecord[*indxTagger]

type xsolTagger struct{}

func (t *xsolTagger) Tag() esm.SubrecordTag { return "XSOL" }

type XSOLdata = record.ZstringSubrecord[*xsolTagger]

type xchgTagger struct{}

func (t *xchgTagger) Tag() esm.SubrecordTag { return "XCHG" }

type XCHGdata = record.Float32Subrecord[*xchgTagger]

type nam9Tagger struct{}

func (t *nam9Tagger) Tag() esm.SubrecordTag { return "NAM9" }

type NAM9data = record.Uint32Subrecord[*nam9Tagger]

type dnamTagger struct{}

func (t *dnamTagger) Tag() esm.SubrecordTag { return "DNAM" }

type DNAMdata = record.ZstringSubrecord[*dnamTagger]

type fltvTagger struct{}

func (t *fltvTagger) Tag() esm.SubrecordTag { return "FLTV" }

type FLTVdata = record.Uint32Subrecord[*fltvTagger]

type knamTagger struct{}

func (t *knamTagger) Tag() esm.SubrecordTag { return "KNAM" }

type KNAMdata = record.ZstringSubrecord[*knamTagger]

type tnamTagger struct{}

func (t *tnamTagger) Tag() esm.SubrecordTag { return "TNAM" }

type TNAMdata = record.ZstringSubrecord[*tnamTagger]

type znamTagger struct{}

func (t *znamTagger) Tag() esm.SubrecordTag { return "ZNAM" }

type ZNAMdata = record.Uint8Subrecord[*znamTagger]

type intvTagger struct{}

func (t *intvTagger) Tag() esm.SubrecordTag { return "INTV" }

type INTVdata = record.BytesSubrecord[*intvTagger]

type dodtTagger struct{}

func (t *dodtTagger) Tag() esm.SubrecordTag { return "DODT" }

type DODTdata = record.BytesSubrecord[*dodtTagger]

type formReferenceDATAtagger struct{}

func (t *formReferenceDATAtagger) Tag() esm.SubrecordTag { return "DATA" }

type FormReferenceDATAdata = record.BytesSubrecord[*formReferenceDATAtagger]

// References to objects in cells are listed as part of the cell data, each beginning with FRMR and NAME fields, followed by a list of fields specific to the object type.
type FormReference struct {
	// Reference ID.
	// Type: uint32
	// Required.
	FRMR *FRMRdata
	// Object ID or "PlayerSaveGame".
	// zstring
	// Required.
	NAME *NAMEdata
	// Reference blocked (value is always 0; present if Blocked is set in the reference's record header, otherwise absent).
	// uint8
	// Optional.
	UNAM *UNAMdata
	// Reference's scale, if applicable and not 1.0.
	// float32
	// Optional.
	XSCL *XSCLdata
	// NPC ID, if applicable (NPC-only).
	// zstring
	// Optional, exclusive with BNAM.
	ANAM *ANAMdata
	// Global variable name
	// zstring
	// Optional, exclusive with ANAM.
	BNAM *BNAMdata
	// Faction ID (not light, NPC, or static)
	// zstring
	// Optional, if present then INDX must also exist.
	CNAM *CNAMdata
	// Faction rank.
	// uint32
	// Optional, if present then CNAM must also exist.
	INDX *INDXdata
	// ID of soul in gem (soul gems only)
	// zstring
	// Optional.
	XSOL *XSOLdata
	// Enchantment charge (charged items with non-zero charges).
	// float32
	// Optional.
	XCHG *XCHGdata
	// Depends on the object type.
	//   uint32 - health remaining (weapons and armor)
	//   uint32 - uses remaining (locks, probes, repair items)
	//   float32 - time remaining (lights)
	// Optional.
	INTV *INTVdata
	// Value (in gold)
	// uint32
	// Optional.
	NAM9 *NAM9data
	// Cell Travel Destination (Rotations are in radians)
	//   float32 - Position X
	//   float32 - Position Y
	//   float32 - Position Z
	//   //   float32 - Rotation X
	//   float32 - Rotation Y
	//   float32 - Rotation Z
	// Optional.
	DODT *DODTdata
	// Cell name for previous DODT, if interior.
	// zstring
	// Optional, must accompany DODT if present.
	DNAM *DNAMdata
	// Lock difficulty
	// uint32
	// Optional.
	FLTV *FLTVdata
	// Key name
	// zstring
	// Optional.
	KNAM *KNAMdata
	// Trap name
	// zstring
	// Optional.
	TNAM *TNAMdata
	// Reference is disabled (always 0). Like UNAM, this will be emitted if the relevant flag is set in the reference's record header. This may only be possible via scripting. Also, even if present in the file, the field appears to be ignored on loading.
	// uint8
	// Optional.
	ZNAM *ZNAMdata
	// Reference position (Rotations are in radians)
	//   float32 - Position X
	//   float32 - Position Y
	//   float32 - Position Z
	//   float32 - Rotation X
	//   float32 - Rotation Y
	//   float32 - Rotation Z
	// Optional.
	DATA *FormReferenceDATAdata
}

func (f *FormReference) OrderedRecords() ([]*esm.Subrecord, error) {
	if f == nil {
		return nil, nil
	}
	orderedSubrecords := []*esm.Subrecord{}
	add := func(p esm.ParsedSubrecord) error {
		if p != nil {
			subRec := esm.Subrecord{}
			if err := subRec.Unmarshal(p); err != nil {
				return err
			}
			orderedSubrecords = append(orderedSubrecords, &subRec)
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
			fr.FRMR = &FRMRdata{}
			if err := fr.FRMR.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case NAME:
			fr.NAME = &NAMEdata{}
			if err := fr.NAME.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case UNAM:
			fr.UNAM = &UNAMdata{}
			if err := fr.UNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case XSCL:
			fr.XSCL = &XSCLdata{}
			if err := fr.XSCL.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case ANAM:
			fr.ANAM = &ANAMdata{}
			if err := fr.ANAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
		case BNAM:
			fr.BNAM = &BNAMdata{}
			if err := fr.BNAM.Unmarshal(sub); err != nil {
				return nil, 0, err
			}
			// todo: finish
		default:
			break subber
		}
		processed++
	}
	return fr, processed, nil
}
