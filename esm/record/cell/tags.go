// Package cell handles https://en.uesp.net/wiki/Morrowind_Mod:Mod_File_Format/CELL
package cell

import "github.com/ernmw/omwpacker/esm"

const (
	CELL esm.RecordTag = "CELL"
)

const (
	// NAME is the Cell Name. Unlike other NAME fields, this is the localized, human-readable name of the cell, not a language-agnostic ID string. Exterior regions are mostly empty strings; for these, the region name is used in the Construction Set.
	NAME esm.SubrecordTag = "NAME"
	// DATA is a 12 byte struct containing flags and position.
	DATA esm.SubrecordTag = "DATA"
	// RGNN is the Region name (exterior and like-exterior only).
	RGNN esm.SubrecordTag = "RGNN"
	// NAM5 is the Map color (exterior and like-exterior only).
	NAM5 esm.SubrecordTag = "NAM5"
	// WHGT is the Water height (interior only), a float32.
	WHGT esm.SubrecordTag = "WHGT"
	// AMBI is a 16 byte struct for Ambient light (ambient, sunlight, fog colors, and fog density).
	AMBI esm.SubrecordTag = "AMBI"
	// MVRF is the Reference ID for a Moved Reference, always the same as the attached FRMR value.
	MVRF esm.SubrecordTag = "MVRF"
	// CNAM is the Name of the cell the reference was moved to (interior cells only) or Faction ID (not light, NPC, or static).
	CNAM esm.SubrecordTag = "CNAM"
	// CNDT is an 8 byte struct containing the Coordinates of the cell the reference was moved to (exterior cells only).
	CNDT esm.SubrecordTag = "CNDT"
	// FRMR is the Reference ID for a Form Reference.
	FRMR esm.SubrecordTag = "FRMR"
	// UNAM is the Reference blocked flag (always 0, present if Blocked is set in the header).
	UNAM esm.SubrecordTag = "UNAM"
	// XSCL is the Reference's scale, if applicable and not 1.0.
	XSCL esm.SubrecordTag = "XSCL"
	// ANAM is the NPC ID, if applicable (NPC-only).
	ANAM esm.SubrecordTag = "ANAM"
	// BNAM is the Global variable name.
	BNAM esm.SubrecordTag = "BNAM"
	// INDX is the Faction rank (uint32).
	INDX esm.SubrecordTag = "INDX"
	// XSOL is the ID of soul in gem (soul gems only).
	XSOL esm.SubrecordTag = "XSOL"
	// XCHG is the Enchantment charge (charged items with non-zero charges), a float32.
	XCHG esm.SubrecordTag = "XCHG"
	// INTV is the Remaining usage (health, uses, or time remaining), depending on object type.
	INTV esm.SubrecordTag = "INTV"
	// NAM9 is a Value (uint32).
	NAM9 esm.SubrecordTag = "NAM9"
	// DODT is a 24 byte struct for Cell Travel Destination (position and rotation in radians).
	DODT esm.SubrecordTag = "DODT"
	// DNAM is the Cell name for previous DODT, if interior.
	DNAM esm.SubrecordTag = "DNAM"
	// FLTV is the Lock difficulty (uint32).
	FLTV esm.SubrecordTag = "FLTV"
	// KNAM is the Key name.
	KNAM esm.SubrecordTag = "KNAM"
	// TNAM is the Trap name.
	TNAM esm.SubrecordTag = "TNAM"
	// ZNAM is the Reference is disabled flag (always 0, present if the relevant flag is set in the header).
	ZNAM esm.SubrecordTag = "ZNAM"
	// DATA is a 24 byte struct for Reference position (position and rotation in radians). Note: This is an inner DATA subrecord for a reference, distinct from the outer CELL DATA.
	DATAReferencePosition esm.SubrecordTag = "DATA"
)
