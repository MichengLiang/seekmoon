package platform

import "os"

// Env captures process environment values used to resolve runtime paths.
type Env struct {
	WorkingDir       string
	XDGCacheHome     string
	Home             string
	IntegrationTests bool
}

// ReadEnv reads the host environment into Env.
func ReadEnv() Env {
	wd, _ := os.Getwd()
	return Env{
		WorkingDir:       wd,
		XDGCacheHome:     os.Getenv("XDG_CACHE_HOME"),
		Home:             os.Getenv("HOME"),
		IntegrationTests: os.Getenv("SEEKMOON_INTEGRATION") != "",
	}
}
