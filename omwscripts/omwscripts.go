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
		for _, attach := range strings.Split(attachList, ",") {
			luafRec, err := (&esm.LUAFdata{Target: strings.TrimSpace(attach)}).Marshal()
			if err != nil {
				return nil, fmt.Errorf("fail to marshal LUAF for %q: %w", attach, err)
			}
			luasRec, err := (&esm.LUASdata{Path: path}).Marshal()
			if err != nil {
				return nil, fmt.Errorf("fail to marshal LUAS for %q: %w", path, err)
			}
			out = append(out, luafRec, luasRec)
		}
	}
	return out, nil
}
