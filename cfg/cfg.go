package cfg

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type PluginEntry struct {
	Name string // lowercased basename
	Path string // full path
}

func findRoot(cfgPath string) (string, error) {

	pathsToCheck := []string{
		cfgPath,
	}
	if wd, err := os.Getwd(); err != nil {
		pathsToCheck = append(pathsToCheck, wd)
	}
	if exe, err := os.Executable(); err != nil {
		pathsToCheck = append(pathsToCheck, filepath.Dir(exe))
	}
	// now check best-known locations
	pathsToCheck = append(pathsToCheck,
		path.Join(os.ExpandEnv("$HOME"), ".config", "openmw"),
		path.Join(os.ExpandEnv("$HOME"), "Library", "Preferences", "openmw"),
		path.Join(os.ExpandEnv("$mydocuments"), "My Games", "OpenMW"),
	)

	for _, path := range pathsToCheck {
		info, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if info.IsDir() {
			subPath := filepath.Join(path, "openmw.cfg")
			if sub, err := os.Stat(subPath); errors.Is(err, os.ErrNotExist) {
				continue
			} else if sub.IsDir() {
				continue
			}
			return subPath, nil
		}
		return path, nil
	}

	return "", fmt.Errorf("resolve openmw.cfg: %q", cfgPath)
}
