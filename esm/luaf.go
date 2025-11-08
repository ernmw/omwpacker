package esm

import (
	"bytes"
	"strings"

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

// LUAFdata is the thing the script is attached to.
type LUAFdata struct {
	Target string
}

func (h *LUAFdata) Tag() tags.SubrecordTag {
	return tags.LUAF
}

func (h *LUAFdata) Unmarshal(sub *Subrecord) error {
	if h == nil || sub == nil {
		return ErrArgumentNil
	}
	h.Target = readPaddedString(sub.Data[0:4])
	return nil
}

func (h *LUAFdata) Marshal() (*Subrecord, error) {
	// make sure NPC gets written as NPC_.
	outTag := h.Target + strings.Repeat("_", 4-min(4, len(h.Target)))
	buff := new(bytes.Buffer)
	if err := writePaddedString(buff, []byte(outTag[:4]), 4); err != nil {
		return nil, err
	}
	return &Subrecord{Tag: h.Tag(), Data: buff.Bytes()}, nil
}
