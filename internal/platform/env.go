package platform

import "os"

type Env struct {
	WorkingDir       string
	XDGCacheHome     string
	Home             string
	IntegrationTests bool
}

func ReadEnv() Env {
	wd, _ := os.Getwd()
	return Env{
		WorkingDir:       wd,
		XDGCacheHome:     os.Getenv("XDG_CACHE_HOME"),
		Home:             os.Getenv("HOME"),
		IntegrationTests: os.Getenv("SEEKMOON_INTEGRATION") != "",
	}
}
