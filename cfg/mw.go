package cfg

import (
	"os"
	"path/filepath"
	"strings"
)

// MWPlugins returns plugins listed in morrowind.ini in the order they appear.
func MWPlugins(iniPath string) ([]PluginEntry, error) {
	masters := []PluginEntry{}
	plugins := []PluginEntry{}

	dataDir := filepath.Join(filepath.Dir(iniPath), "Data Files")
	b, err := os.ReadFile(iniPath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])
		if strings.HasPrefix(key, "gamefile") {
			pluginPath := filepath.Join(dataDir, val)
			if _, err := os.Stat(pluginPath); err == nil {
				ext := strings.ToLower(filepath.Ext(val))
				name := strings.ToLower(val)
				if ext == ".esm" {
					masters = append(masters, PluginEntry{Name: name, Path: pluginPath})
				} else if ext == ".esp" {
					plugins = append(plugins, PluginEntry{Name: name, Path: pluginPath})
				}
			}
		}
	}

	// combine masters then plugins (same semantic as original)
	out := make([]PluginEntry, 0, len(masters)+len(plugins))
	out = append(out, masters...)
	out = append(out, plugins...)
	if len(out) == 0 {
		return nil, nil
	}
	return out, nil
}
