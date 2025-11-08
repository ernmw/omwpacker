package tags

// RecordTag identifies the type of record.
// See https://github.com/OpenMW/openmw/blob/39d117e362808dc13cd411debcb48e363e11639c/components/esm/defs.hpp#L78
type RecordTag string

const (
	// Lua record.
	LUAL RecordTag = "LUAL"

	// Activator.
	ACTI RecordTag = "ACTI"

	// Potion.
	ALCH RecordTag = "ALCH"

	// Alchemy Apparatus.
	APPA RecordTag = "APPA"

	// Armor.
	ARMO RecordTag = "ARMO"

	// Body Parts.
	BODY RecordTag = "BODY"

	// Book.
	BOOK RecordTag = "BOOK"

	// Birthsign.
	BSGN RecordTag = "BSGN"

	// Cell.
	CELL RecordTag = "CELL"

	// Class.
	CLAS RecordTag = "CLAS"

	// Clothing.
	CLOT RecordTag = "CLOT"

	// Container.
	CONT RecordTag = "CONT"

	// Creature.
	CREA RecordTag = "CREA"

	// Dialog Topic.
	DIAL RecordTag = "DIAL"

	// Door.
	DOOR RecordTag = "DOOR"

	// Enchantment.
	ENCH RecordTag = "ENCH"

	// Data Object.
	TYPE RecordTag = "Type"

	// Faction.
	FACT RecordTag = "FACT"

	// Global.
	GLOB RecordTag = "GLOB"

	// Game Setting.
	GMST RecordTag = "GMST"

	// Dialog response.
	INFO RecordTag = "INFO"

	// Ingredient.
	INGR RecordTag = "INGR"

	// Land.
	LAND RecordTag = "LAND"

	// Leveled Creature.
	LEVC RecordTag = "LEVC"

	// Leveled Item.
	LEVI RecordTag = "LEVI"

	// Light.
	LIGH RecordTag = "LIGH"

	// Lockpicking Items.
	LOCK RecordTag = "LOCK"

	// Land Texture.
	LTEX RecordTag = "LTEX"

	// Magic Effect.
	MGEF RecordTag = "MGEF"

	// Misc. Item.
	MISC RecordTag = "MISC"

	// Non-Player Character.
	NPC_ RecordTag = "NPC_"

	// Path grid.
	PGRD RecordTag = "PGRD"

	// Probe Items.
	PROB RecordTag = "PROB"

	// Race.
	RACE RecordTag = "RACE"

	// Region.
	REGN RecordTag = "REGN"

	// Repair Items.
	REPA RecordTag = "REPA"

	// Script.
	SCPT RecordTag = "SCPT"

	// Skill.
	SKIL RecordTag = "SKIL"

	// Sound Generator.
	SNDG RecordTag = "SNDG"

	// Sound.
	SOUN RecordTag = "SOUN"

	// Spell.
	SPEL RecordTag = "SPEL"

	// Start Script.
	SSCR RecordTag = "SSCR"

	// Static.
	STAT RecordTag = "STAT"

	// Tes3 root.
	TES3 RecordTag = "TES3"

	// Weapon.
	WEAP RecordTag = "WEAP"
)

type SubrecordTag string

const (
	// Form. First part of TES3?
	FORM SubrecordTag = "FORM"
	/*
		Header
		    float32 - Version (1.2 for Morrowind, 1.3 for Bloodmoon and Tribunal)
		    uint32 - Flags

		        0x1 RecordTag = file should be treated as a master, regardless of the file extension

		    char[32] - Company name string
		    char[256] - File description
		    uint32 - Number of records after this one
	*/
	HEDR SubrecordTag = "HEDR"
	// Master filename (a null-terminated string).
	// Each pair of MAST/DATA subrecords represents a single master of the mod file.
	// Master files are listed in load order at the time the mod was saved.
	MAST SubrecordTag = "MAST"
	// Size of the previous master file in bytes (used for version tracking of plugin)
	DATA SubrecordTag = "DATA"

	LUAF SubrecordTag = "LUAF"
	LUAW SubrecordTag = "LUAW"
	LUAS SubrecordTag = "LUAS"
	LUAD SubrecordTag = "LUAD"
	LUAT SubrecordTag = "LUAT"
	LUAC SubrecordTag = "LUAC"
)
