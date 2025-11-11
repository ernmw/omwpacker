// LAND records contain information about the landscape of exterior cells.
// More specifically, it defines 65x65 arrays of vertex heights, normals,
// colors, and a smaller 16x16 array of textures. It also defines a 9x9 array
// of heights that the game can load to quickly build the world map.
//
// https://github.com/OpenMW/openmw/blob/429305401ee7486f160cb1bbd2196fc80d33dc3a/components/esm3/loadland.cpp
//
//go:generate go run ../generator/gen.go subrecords.json
package land

import "github.com/ernmw/omwpacker/esm"

// LAND handles https://en.uesp.net/wiki/Morrowind_Mod:Mod_File_Format/LAND
const LAND esm.RecordTag = "LAND"
