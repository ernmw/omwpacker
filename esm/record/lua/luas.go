package lua

import (
	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/record"
)

// https://gitlab.com/OpenMW/openmw/-/blob/master/components/lua/serialization.cpp

/*
 * // Records:
 // LUAL - LuaScriptsCfg - list of all scripts (in content files)        - added in !1271
 // LUAM - MWLua::LuaManager (in saves)
 //
 // Subrecords:
 // LUAF - LuaScriptCfg::mFlags                                          - added in !1271
 // LUAW - Start of MWLua::WorldView data
 // LUAE - Start of MWLua::LocalEvent or MWLua::GlobalEvent (eventName)
 // LUAS - VFS path to a Lua script
 // LUAD - Serialized Lua variable
 // LUAT - MWLua::ScriptsContainer::Timer
 // LUAC - Name of a timer callback (string)
*/

type luasTagger struct{}

func (t *luasTagger) Tag() esm.SubrecordTag { return "LUAS" }

// LUASdata is the script path.
type LUASdata = record.ZstringSubrecord[*luasTagger]
