// Package store owns SeekMoon persistence paths and file mechanics.
package store

import (
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/yumiaura/seekmoon/internal/platform"
)

// Paths groups all project and cache paths used by SeekMoon.
type Paths struct {
	ProjectRoot string
	ProjectDir  string
	Snapshots   string
	Sessions    string
	Records     string
	Reports     string
	Probes      string
	Sources     string
	Logs        string
	CacheRoot   string
	Mooncakes   string
	Assets      string
	GitHub      string
	MoonIndex   string
}

// ResolvePaths derives project and cache paths from environment values.
func ResolvePaths(env platform.Env) Paths {
	root := env.WorkingDir
	if root == "" {
		root = "."
	}
	cacheRoot := filepath.Join(xdg.CacheHome, "seekmoon")
	if env.XDGCacheHome != "" {
		cacheRoot = filepath.Join(env.XDGCacheHome, "seekmoon")
	}
	home := env.Home
	if home == "" {
		home = root
	}
	projectDir := filepath.Join(root, ".seekmoon")
	return Paths{
		ProjectRoot: root,
		ProjectDir:  projectDir,
		Snapshots:   filepath.Join(projectDir, "snapshots"),
		Sessions:    filepath.Join(projectDir, "sessions"),
		Records:     filepath.Join(projectDir, "records"),
		Reports:     filepath.Join(projectDir, "reports"),
		Probes:      filepath.Join(projectDir, "probes"),
		Sources:     filepath.Join(projectDir, "sources"),
		Logs:        filepath.Join(projectDir, "logs"),
		CacheRoot:   cacheRoot,
		Mooncakes:   filepath.Join(cacheRoot, "mooncakes"),
		Assets:      filepath.Join(cacheRoot, "assets"),
		GitHub:      filepath.Join(cacheRoot, "github"),
		MoonIndex:   filepath.Join(home, ".moon", "registry", "index", "user"),
	}
}

// SafeName converts arbitrary identifiers into filesystem-safe names.
func SafeName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unnamed"
	}
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_' {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if r == '-' && !lastDash {
			b.WriteByte('-')
			lastDash = true
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "unnamed"
	}
	return out
}
