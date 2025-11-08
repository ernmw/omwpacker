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

type unamTagger struct{}

func (t *unamTagger) Tag() esm.SubrecordTag { return "UNAM" }

type UNAMdata = record.Uint8Subrecord[*anamTagger]

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
