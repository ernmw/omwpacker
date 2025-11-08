package lua

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
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
// https://gitlab.com/OpenMW/openmw/-/blob/master/components/esm/luascripts.cpp#L62
type LUAFdata struct {
	// Flags cover types of script attachment which don't logically correlate to a type of gameobject -> Player, Global, Custom, Menu, etc.
	// Player is one of them, since they are an NPC.
	//
	// sGlobal = 1ull << 0; // start as a global script
	//
	// sCustom = 1ull << 1; // local; can be attached/detached by a global script
	//
	// sPlayer = 1ull << 2; // auto attach to players
	//
	// sMerge = 1ull << 3; // merge with configuration from previous content files
	//
	// sMenu = 1ull << 4; // start as a menu script
	Flags uint32
	// Targets is a list of 4-byte strings which map to ccfour constants.
	Targets []string
}

func (h *LUAFdata) Tag() esm.SubrecordTag {
	return LUAF
}

func (h *LUAFdata) Unmarshal(sub *esm.Subrecord) error {
	if h == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	h.Flags = binary.LittleEndian.Uint32(sub.Data[0:4])

	rawTargets := util.ReadPaddedString(sub.Data[4:])
	h.Targets = make([]string, len(rawTargets)/4)
	for i := 0; i < len(rawTargets); i = i + 4 {
		h.Targets[i/4] = rawTargets[i : i+4]
	}

	return nil
}

func (h *LUAFdata) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.LittleEndian, h.Flags); err != nil {
		return nil, err
	}
	for _, target := range h.Targets {
		outTag := target + strings.Repeat("_", 4-min(4, len(target)))
		if err := util.WritePaddedString(buff, []byte(outTag), 4); err != nil {
			return nil, err
		}
	}
	return &esm.Subrecord{Tag: h.Tag(), Data: buff.Bytes()}, nil
}
