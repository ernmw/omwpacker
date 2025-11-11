package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// OpenMWPlugins now returns resolved plugin absolute paths (searching data dirs),
// plus dataPaths (BSAs first, then folders).
func OpenMWPlugins(path string) ([]string, []string, error) {
	cfgPath, err := findRoot(path)
	if err != nil {
		return nil, nil, fmt.Errorf("find root openmw.cfg file %q: %w", path, err)
	}

	visited := map[string]bool{}
	var contexts []configContext

	if err := loadConfigRecursive(cfgPath, &contexts, visited); err != nil {
		return nil, nil, err
	}

	// Build ordered data directories list (lowest -> highest priority),
	// and list of BSAs (lowest priority).
	var bsaArchives []string
	var dataDirs []string
	var userData []string
	var dataLocal []string

	for _, ctx := range contexts {
		bsaArchives = append(bsaArchives, ctx.bsaArchives...)
		dataDirs = append(dataDirs, ctx.dataDirs...)
		userData = append(userData, ctx.userData...)
		dataLocal = append(dataLocal, ctx.dataLocal...)
	}

	// Compose dataPaths: BSAs first, then folders, then userdata, then data-local
	dataPaths := make([]string, 0, len(bsaArchives)+len(dataDirs)+len(userData)+len(dataLocal))
	dataPaths = append(dataPaths, bsaArchives...)
	dataPaths = append(dataPaths, dataDirs...)
	dataPaths = append(dataPaths, userData...)
	dataPaths = append(dataPaths, dataLocal...)

	// Resolve plugin names to absolute paths by searching dataDirs in order
	pluginPaths := resolvePluginNames(contexts, dataDirs)

	return pluginPaths, dataPaths, nil
}

// resolvePluginNames resolves plugin names declared in contexts into absolute
// file paths by searching the provided dataDirs in order (lowest -> highest).
func resolvePluginNames(contexts []configContext, dataDirs []string) []string {
	var resolved []string

	// iterate contexts in order (lowest -> highest priority) and collect plugin names in that order
	for _, ctx := range contexts {
		for _, pluginName := range ctx.pluginNames {
			// If pluginName already contains a path separator or is absolute, resolve it directly
			cleanName := strings.TrimSpace(pluginName)
			cleanName = strings.Trim(cleanName, "\"") // remove quotes if present
			cleanName = expandTokens(cleanName)

			if filepath.IsAbs(cleanName) || strings.ContainsAny(cleanName, string(os.PathSeparator)+"/") {
				// Absolute or path: return absolute path (resolve relative to config if not absolute)
				abs := cleanName
				if !filepath.IsAbs(abs) {
					abs = filepath.Join(ctx.baseDir, abs)
				}
				if p, err := filepath.Abs(abs); err == nil {
					resolved = append(resolved, p)
					continue
				}
				// fallback: use as-is
				resolved = append(resolved, abs)
				continue
			}

			// Otherwise search each dataDir for the file
			found := ""
			for _, d := range dataDirs {
				candidate := filepath.Join(d, cleanName)
				if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
					if p, err := filepath.Abs(candidate); err == nil {
						found = p
					} else {
						found = candidate
					}
					break
				}
				// case-insensitive fallback: check dir entries
				if dirEntries, err := os.ReadDir(d); err == nil {
					lowerWant := strings.ToLower(cleanName)
					for _, e := range dirEntries {
						if strings.ToLower(e.Name()) == lowerWant {
							candidate = filepath.Join(d, e.Name())
							if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
								if p, err := filepath.Abs(candidate); err == nil {
									found = p
								} else {
									found = candidate
								}
								break
							}
						}
					}
					if found != "" {
						break
					}
				}
			}

			if found != "" {
				resolved = append(resolved, found)
				continue
			}

			// Not found in dataDirs: last resort - resolve relative to the context's cfg directory.
			backupCandidate := filepath.Join(ctx.baseDir, cleanName)
			if fi, err := os.Stat(backupCandidate); err == nil && !fi.IsDir() {
				if p, err := filepath.Abs(backupCandidate); err == nil {
					resolved = append(resolved, p)
				} else {
					resolved = append(resolved, backupCandidate)
				}
			} else {
				// Not found anywhere: still append the plugin name (unresolved)
				resolved = append(resolved, cleanName)
			}
		}
	}

	return resolved
}

type configContext struct {
	path          string
	baseDir       string
	dataDirs      []string
	dataLocal     []string
	userData      []string
	bsaArchives   []string
	pluginNames   []string // store plugin *names* as declared (not resolved)
	nestedConfigs []string
	replaceConfig bool
}

// loadConfigRecursive recursively loads an openmw.cfg and any referenced sub-configs.
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
				ctx.nestedConfigs = append(ctx.nestedConfigs, nested)
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
				// keep the raw name — we'll resolve later using data dirs
				ctx.pluginNames = append(ctx.pluginNames, val)
			}
		}
	}

	// Apply replace=config semantics — drop earlier configs if needed
	if ctx.replaceConfig && len(*contexts) > 0 {
		*contexts = (*contexts)[:0]
	}
	*contexts = append(*contexts, ctx)

	// Recursively load nested configs (higher priority)
	for _, sub := range ctx.nestedConfigs {
		if err := loadConfigRecursive(sub, contexts, visited); err != nil {
			return err
		}
	}

	return nil
}

// parseKV splits "key=value" and trims quotes.
func parseKV(line string) (key, val string, ok bool) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key = strings.ToLower(strings.TrimSpace(parts[0]))
	val = strings.TrimSpace(parts[1])
	val = strings.Trim(val, "\"")
	return key, val, true
}

// verifyPath expands relative paths and tokens into an absolute path.
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

// expandTokens replaces OpenMW-style tokens (?local?, ?userconfig?, etc.)
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
