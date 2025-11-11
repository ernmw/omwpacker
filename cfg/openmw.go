package cfg

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// OpenMWPlugins returns openmw plugins from openmw.cfg in the order they appear.
func OpenMWPlugins(cfgpath string) ([]PluginEntry, error) {
	var dataFolders []string
	var pendingContents []PluginEntry

	f, err := os.Open(cfgpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])
		if key == "data" {
			p := verifyPath(cfgpath, val)
			if p != "" {
				if info, err := os.Stat(p); err == nil && info.IsDir() {
					dataFolders = append(dataFolders, p)
				}
			}
		} else if key == "content" {
			ext := strings.ToLower(filepath.Ext(val))
			validExts := map[string]bool{".esm": true}
			validExts[".esp"] = true
			validExts[".omwaddon"] = true
			if validExts[ext] {
				pendingContents = append(pendingContents, PluginEntry{Name: strings.ToLower(filepath.Base(val)), Path: val})
			}
		}
	}

	// resolve pendingContents against dataFolders preserving order
	var out []PluginEntry
	for _, pc := range pendingContents {
		// search folders in order for first match
		found := false
		for _, dataPath := range dataFolders {
			candidate := filepath.Join(dataPath, pc.Path)
			if _, err := os.Stat(candidate); err == nil {
				out = append(out, PluginEntry{Name: strings.ToLower(filepath.Base(pc.Path)), Path: candidate})
				found = true
				break
			}
			// also check lowercase matches in directory
			items, _ := os.ReadDir(dataPath)
			for _, item := range items {
				if strings.EqualFold(item.Name(), pc.Path) {
					out = append(out, PluginEntry{Name: strings.ToLower(item.Name()), Path: filepath.Join(dataPath, item.Name())})
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		// If not found, skip it (mirrors original behavior)
	}

	if len(out) == 0 {
		return nil, nil
	}
	return out, nil
}

func verifyPath(baseCfgPath, s string) string {
	s = strings.Trim(s, "\" ")
	if s == "" {
		return ""
	}
	if filepath.IsAbs(s) {
		return s
	}
	absPath, err := filepath.Abs(filepath.Join(filepath.Dir(baseCfgPath), s))
	if err != nil {
		return ""
	}
	return absPath
}
