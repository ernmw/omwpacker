package esm

import (
	"bytes"

	"github.com/ernmw/omwpacker/esm/tags"
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

func (h *LUASdata) Tag() tags.SubrecordTag {
	return tags.LUAS
}

func (h *LUASdata) Unmarshal(sub *Subrecord) error {
	if h == nil || sub == nil {
		return ErrArgumentNil
	}
	h.Path = string(sub.Data)
	return nil
}

func (h *LUASdata) Marshal() (*Subrecord, error) {
	buff := new(bytes.Buffer)
	if _, err := buff.WriteString(h.Path); err != nil {
		return nil, err
	}
	return &Subrecord{Tag: h.Tag(), Data: buff.Bytes()}, nil
}
