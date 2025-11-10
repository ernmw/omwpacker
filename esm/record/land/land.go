//go:generate go run ../generator/gen.go subrecords.json
package land

import "github.com/ernmw/omwpacker/esm"

// LAND handles https://en.uesp.net/wiki/Morrowind_Mod:Mod_File_Format/LAND
const LAND esm.RecordTag = "LAND"
