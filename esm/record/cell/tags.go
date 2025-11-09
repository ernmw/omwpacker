// Package cell handles https://en.uesp.net/wiki/Morrowind_Mod:Mod_File_Format/CELL
package cell

import "github.com/ernmw/omwpacker/esm"

const (
	CELL esm.RecordTag = "CELL"
)

const (
	// DATA is a 12 byte struct containing flags and position.
	DATA esm.SubrecordTag = "DATA"
	// AMBI is a 16 byte struct for Ambient light (ambient, sunlight, fog colors, and fog density).
	AMBI esm.SubrecordTag = "AMBI"
)
