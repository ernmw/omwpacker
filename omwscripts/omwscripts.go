package omwscripts

import (
	"fmt"
	"strings"

	"github.com/ernmw/omwpacker/esm"
)

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
		luaf := &esm.LUAFdata{Targets: []string{}}
		for _, attach := range strings.Split(attachList, ",") {
			key := strings.ToUpper(strings.TrimSpace(attach))
			switch key {
			case "GLOBAL":
				luaf.Flags = luaf.Flags | 1<<0
			case "CUSTOM":
				key = "GLOBAL"
				luaf.Flags = luaf.Flags | 1<<1
			case "PLAYER":
				key = "NPC"
				luaf.Flags = luaf.Flags | 1<<2
			case "MENU":
				key = "GLOBAL"
				luaf.Flags = luaf.Flags | 1<<4
			}
			luaf.Targets = append(luaf.Targets, key)
		}
		luafRec, err := luaf.Marshal()
		if err != nil {
			return nil, fmt.Errorf("fail to marshal LUAF: %w", err)
		}
		luasRec, err := (&esm.LUASdata{Path: path}).Marshal()
		if err != nil {
			return nil, fmt.Errorf("fail to marshal LUAS for %q: %w", path, err)
		}
		out = append(out, luasRec, luafRec)
	}
	return out, nil
}
