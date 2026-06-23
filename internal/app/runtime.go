// Package app composes SeekMoon runtime dependencies for one process execution.
package app

import (
	"net/http"
	"time"

	"github.com/yumiaura/seekmoon/internal/output"
	"github.com/yumiaura/seekmoon/internal/platform"
	"github.com/yumiaura/seekmoon/internal/service"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

type Runtime struct {
	Env      platform.Env
	Clock    platform.Clock
	FS       platform.FS
	HTTP     *http.Client
	Runner   platform.Runner
	Paths    store.Paths
	Sources  SourceRegistry
	Stores   store.Registry
	Services ServiceRegistry
	Renderer output.Renderer
}

type SourceRegistry struct {
	Mooncakes  source.MooncakesClient
	Assets     source.AssetClient
	Skills     source.SkillsClient
	LocalIndex source.LocalIndexReader
	LocalCache source.LocalCacheReader
	Project    source.ProjectReader
	MoonCLI    source.MoonCLI
	Repository source.RepositoryReader
}

type ServiceRegistry struct {
	Registry service.Registry
	Sync     service.SyncService
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
	runtime.registerBatchB()
	for _, option := range options {
		option(runtime)
	}
	runtime.registerBatchB()
	return runtime, nil
}

func WithEnv(env platform.Env) Option {
	return func(runtime *Runtime) {
		runtime.Env = env
		runtime.Paths = store.ResolvePaths(env)
		runtime.Stores = store.NewRegistry(runtime.FS, runtime.Paths)
		runtime.registerBatchB()
	}
}

func WithFS(fs platform.FS) Option {
	return func(runtime *Runtime) {
		runtime.FS = fs
		runtime.Runner = platform.ExecRunner{FS: fs}
		runtime.Stores = store.NewRegistry(fs, runtime.Paths)
		runtime.registerBatchB()
	}
}

func WithClock(clock platform.Clock) Option {
	return func(runtime *Runtime) {
		runtime.Clock = clock
	}
}

func (runtime *Runtime) registerBatchB() {
	fetcher := source.Fetcher{Client: runtime.HTTP, Clock: runtime.Clock}
	runtime.Sources = SourceRegistry{
		Mooncakes:  source.MooncakesClient{Fetcher: fetcher},
		Assets:     source.AssetClient{Fetcher: fetcher},
		Skills:     source.SkillsClient{Fetcher: fetcher},
		LocalIndex: source.LocalIndexReader{FS: runtime.FS, Clock: runtime.Clock},
		LocalCache: source.LocalCacheReader{FS: runtime.FS},
		Project:    source.ProjectReader{FS: runtime.FS, Clock: runtime.Clock},
		MoonCLI:    source.MoonCLI{Runner: runtime.Runner, Paths: runtime.Paths},
		Repository: source.RepositoryReader{Clock: runtime.Clock},
	}
	syncService := service.SyncService{
		Mooncakes: runtime.Sources.Mooncakes,
		Snapshots: runtime.Stores.Snapshots,
		Now: func() time.Time {
			return runtime.Clock.Now()
		},
	}
	runtime.Services = ServiceRegistry{
		Sync: syncService,
	}
	runtime.Services.Registry = service.NewPendingRegistry(syncService)
	runtime.Renderer = output.DefaultRenderer{}
}
