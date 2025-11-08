package lua

import (
	"bytes"

	"github.com/ernmw/omwpacker/esm"
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

// LUASdata is the script path.
type LUASdata struct {
	Path string
}

func (h *LUASdata) Tag() esm.SubrecordTag {
	return LUAS
}

func (h *LUASdata) Unmarshal(sub *esm.Subrecord) error {
	if h == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	h.Path = string(sub.Data)
	return nil
}

func (h *LUASdata) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)
	if _, err := buff.WriteString(h.Path); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: h.Tag(), Data: buff.Bytes()}, nil
}
