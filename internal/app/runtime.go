// Package app composes SeekMoon runtime dependencies for one process execution.
package app

import (
	"net/http"
	"time"

	"github.com/yumiaura/seekmoon/internal/platform"
	"github.com/yumiaura/seekmoon/internal/store"
)

type Runtime struct {
	Env    platform.Env
	Clock  platform.Clock
	FS     platform.FS
	HTTP   *http.Client
	Runner platform.Runner
	Paths  store.Paths
	Stores store.Registry
}

type Option func(*Runtime)

func NewRuntime(options ...Option) (*Runtime, error) {
	env := platform.ReadEnv()
	paths := store.ResolvePaths(env)
	fs := platform.OSFS{}
	runtime := &Runtime{
		Env:   env,
		Clock: platform.SystemClock{},
		FS:    fs,
		HTTP:  platform.NewHTTPClient(30 * time.Second),
		Paths: paths,
	}
	runtime.Runner = platform.ExecRunner{FS: runtime.FS}
	runtime.Stores = store.NewRegistry(runtime.FS, runtime.Paths)
	for _, option := range options {
		option(runtime)
	}
	return runtime, nil
}

func WithEnv(env platform.Env) Option {
	return func(runtime *Runtime) {
		runtime.Env = env
		runtime.Paths = store.ResolvePaths(env)
		runtime.Stores = store.NewRegistry(runtime.FS, runtime.Paths)
	}
}

func WithFS(fs platform.FS) Option {
	return func(runtime *Runtime) {
		runtime.FS = fs
		runtime.Runner = platform.ExecRunner{FS: fs}
		runtime.Stores = store.NewRegistry(fs, runtime.Paths)
	}
}

func WithClock(clock platform.Clock) Option {
	return func(runtime *Runtime) {
		runtime.Clock = clock
	}
}
