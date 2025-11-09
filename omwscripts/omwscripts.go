// Package omwscripts handles conversion of .omwscript files to subrecords.
// Basically a port of https://github.com/OpenMW/openmw/blob/39d117e362808dc13cd411debcb48e363e11639c/components/lua/configuration.cpp
package omwscripts

import (
	"fmt"
	"strings"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/record/lua"
)

var flagsByName = map[string]uint32{
	"GLOBAL": 1 << 0,
	"CUSTOM": 1 << 1,
	"PLAYER": 1 << 2,
	"MENU":   1 << 4,
}

var tagsByName = map[string]esm.RecordTag{
	"ACTIVATOR":  "ACTI",
	"ARMOR":      "ARMO",
	"BOOK":       "BOOK",
	"CLOTHING":   "CLOT",
	"CONTAINER":  "CONT",
	"CREATURE":   "CREA",
	"DOOR":       "DOOR",
	"INGREDIENT": "INGR",
	"LIGHT":      "LIGH",
	"MISC_ITEM":  "MISC",
	"NPC":        "NPC_",
	"POTION":     "ALCH",
	"WEAPON":     "WEAP",
	"APPARATUS":  "APPA",
	"LOCKPICK":   "LOCK",
	"PROBE":      "PROB",
	"REPAIR":     "REPA",
}

func Package(content string) ([]*esm.Subrecord, error) {
	lines := strings.Split(string(content), "\n")

	out := []*esm.Subrecord{}
	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		// Expect "ATTACH: path"
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line %d: %q (expected 'ATTACH: path')", i+1, line)
		}
		attachList := strings.TrimSpace(parts[0])
		path := strings.TrimSpace(parts[1])
		if attachList == "" || path == "" {
			return nil, fmt.Errorf("invalid line %d: %q (empty attach or path)", i+1, line)
		}
		luaf := &lua.LUAFdata{Targets: []string{}}
		for _, attach := range strings.Split(attachList, ",") {
			key := strings.ToUpper(strings.TrimSpace(attach))
			if flag, ok := flagsByName[key]; ok {
				luaf.Flags = luaf.Flags | flag
			} else if target, ok := tagsByName[key]; ok {
				luaf.Targets = append(luaf.Targets, string(target))
			} else {
				return nil, fmt.Errorf("unknown attach key %q", attach)
			}
		}
		luafRec, err := luaf.Marshal()
		if err != nil {
			return nil, fmt.Errorf("fail to marshal LUAF: %w", err)
		}
		luasRec, err := (&lua.LUASdata{Value: path}).Marshal()
		if err != nil {
			return nil, fmt.Errorf("fail to marshal LUAS for %q: %w", path, err)
		}
		out = append(out, luasRec, luafRec)
		// after LUAF, there's LUAR* and then LUAI*
	}
	return out, nil
}
