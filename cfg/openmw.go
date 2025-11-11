package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type configContext struct {
	path          string
	baseDir       string
	dataDirs      []string
	dataLocal     []string
	userData      []string
	bsaArchives   []string
	content       []PluginEntry
	nestedConfig  []string
	replaceConfig bool
}

// OpenMWPlugins loads all plugin and data paths by reading openmw.cfg files
// recursively, respecting OpenMW’s replace/config semantics and token rules.
func OpenMWPlugins(cfgPath string) ([]PluginEntry, []string, error) {
	rootCfgPath, err := findRoot(cfgPath)
	if err != nil {
		return nil, nil, fmt.Errorf("find root cfg path %q: %w", cfgPath, err)
	}

	visited := map[string]bool{}
	var contexts []configContext

	if err := loadConfigRecursive(rootCfgPath, &contexts, visited); err != nil {
		return nil, nil, err
	}

	// OpenMW merges lowest → highest; child configs override parent.
	// We built in call order, so contexts is in ascending priority order already.

	var allPlugins []PluginEntry
	var bsaArchives, dataDirs, userData, dataLocal []string

	for _, ctx := range contexts {
		allPlugins = append(allPlugins, ctx.content...)
		bsaArchives = append(bsaArchives, ctx.bsaArchives...)
		dataDirs = append(dataDirs, ctx.dataDirs...)
		userData = append(userData, ctx.userData...)
		dataLocal = append(dataLocal, ctx.dataLocal...)
	}

	dataPaths := append([]string{}, bsaArchives...)
	dataPaths = append(dataPaths, dataDirs...)
	dataPaths = append(dataPaths, userData...)
	dataPaths = append(dataPaths, dataLocal...)

	return allPlugins, dataPaths, nil
}

// loadConfigRecursive recursively loads a config and any referenced sub-configs.
func loadConfigRecursive(cfgPath string, contexts *[]configContext, visited map[string]bool) error {
	cfgPath, _ = filepath.Abs(cfgPath)
	if visited[cfgPath] {
		return nil
	}
	visited[cfgPath] = true

	cfgDir := filepath.Dir(cfgPath)
	ctx := configContext{path: cfgPath, baseDir: cfgDir}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("read cfg: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := parseKV(line)
		if !ok {
			continue
		}
		val = expandTokens(strings.TrimSpace(val))

		switch key {
		case "config":
			p := verifyPath(cfgPath, val)
			nested := filepath.Join(p, "openmw.cfg")
			if _, err := os.Stat(nested); err == nil {
				ctx.nestedConfig = append(ctx.nestedConfig, nested)
			}

		case "replace":
			if strings.Contains(val, "config") {
				ctx.replaceConfig = true
			}

		case "data":
			ctx.dataDirs = append(ctx.dataDirs, verifyPath(cfgPath, val))

		case "data-local":
			ctx.dataLocal = append(ctx.dataLocal, verifyPath(cfgPath, val))

		case "user-data":
			ctx.userData = append(ctx.userData, verifyPath(cfgPath, val))

		case "fallback-archive":
			ctx.bsaArchives = append(ctx.bsaArchives, verifyPath(cfgPath, val))

		case "content":
			ext := strings.ToLower(filepath.Ext(val))
			if ext == ".esm" || ext == ".esp" || ext == ".omwaddon" {
				ctx.content = append(ctx.content,
					PluginEntry{Name: strings.ToLower(filepath.Base(val)), Path: val})
			}
		}
	}

	// Recurse into child configs — but if replace=config, drop earlier ones
	if ctx.replaceConfig && len(*contexts) > 0 {
		*contexts = (*contexts)[:0] // clear previous
	}
	*contexts = append(*contexts, ctx)

	for _, sub := range ctx.nestedConfig {
		if err := loadConfigRecursive(sub, contexts, visited); err != nil {
			return err
		}
	}
	return nil
}

// parseKV splits "key=value" respecting quoted right-hand values.
func parseKV(line string) (key, val string, ok bool) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key = strings.ToLower(strings.TrimSpace(parts[0]))
	val = strings.TrimSpace(parts[1])
	val = strings.Trim(val, "\"") // tolerate quoted form
	return key, val, true
}

func verifyPath(cfgPath, p string) string {
	p = strings.TrimSpace(strings.Trim(p, "\""))
	p = expandTokens(p)

	if strings.HasPrefix(p, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			p = filepath.Join(home, p[1:])
		}
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(filepath.Dir(cfgPath), p)
	}
	abs, _ := filepath.Abs(p)
	return abs
}

// expandTokens expands ?local?, ?userconfig?, ?userdata?, ?global?
func expandTokens(p string) string {
	home, _ := os.UserHomeDir()
	cfgHome := os.Getenv("XDG_CONFIG_HOME")
	dataHome := os.Getenv("XDG_DATA_HOME")

	if cfgHome == "" {
		cfgHome = filepath.Join(home, ".config")
	}
	if dataHome == "" {
		dataHome = filepath.Join(home, ".local", "share")
	}

	tokens := map[string]string{
		"?local?":      filepath.Dir(os.Args[0]),
		"?userconfig?": filepath.Join(cfgHome, "openmw"),
		"?userdata?":   filepath.Join(dataHome, "openmw"),
		"?global?": func() string {
			switch runtime.GOOS {
			case "darwin":
				return "/Library/Application Support/"
			case "windows":
				return "C:\\Program Files\\OpenMW"
			default:
				return "/usr/share/games"
			}
		}(),
	}

	for k, v := range tokens {
		if strings.Contains(p, k) {
			p = strings.ReplaceAll(p, k, v)
		}
	}
	return p
}
