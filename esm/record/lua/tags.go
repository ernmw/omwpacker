// Package lua handles LUAL and LUAM records, which are
// specific to openmw.
//
//go:generate go run ../generator/gen.go subrecords.json
package lua

import "github.com/ernmw/omwpacker/esm"

const (
	// LUAL - LuaScriptsCfg - list of all scripts (in content files)
	LUAL esm.RecordTag = "LUAL"
	// LUAM - MWLua::LuaManager (in saves)
	LUAM esm.RecordTag = "LUAM"
)

const (
	// LUAF - LuaScriptCfg::mFlags and ESM::RecNameInts list
	LUAF esm.SubrecordTag = "LUAF"
	// LUAW - Simulation time and last generated RefNum
	LUAW esm.SubrecordTag = "LUAW"
	// LUAE - Start of MWLua::LocalEvent or MWLua::GlobalEvent (eventName)
	LUAE esm.SubrecordTag = "LUAE"
	// LUAD - Serialized Lua variable
	LUAD esm.SubrecordTag = "LUAD"
	// LUAT - MWLua::ScriptsContainer::Timer
	LUAT esm.SubrecordTag = "LUAT"
	// LUAC - Name of a timer callback (string)
	LUAC esm.SubrecordTag = "LUAC"
	// LUAR - Attach script to a specific record (LuaScriptCfg::PerRecordCfg)
	LUAR esm.SubrecordTag = "LUAR"
	// LUAI - Attach script to a specific instance (LuaScriptCfg::PerRefCfg)
	LUAI esm.SubrecordTag = "LUAI"
)
