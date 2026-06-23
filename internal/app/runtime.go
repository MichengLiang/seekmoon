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
		Mooncakes:  runtime.Sources.Mooncakes,
		MoonCLI:    runtime.Sources.MoonCLI,
		LocalIndex: runtime.Sources.LocalIndex,
		Snapshots:  runtime.Stores.Snapshots,
		Paths:      runtime.Paths,
		Now: func() time.Time {
			return runtime.Clock.Now()
		},
	}
	doctor := service.DoctorFlow{Project: runtime.Sources.Project, MoonCLI: runtime.Sources.MoonCLI, Paths: runtime.Paths}
	search := service.SearchFlow{Mooncakes: runtime.Sources.Mooncakes, Assets: runtime.Sources.Assets, Snapshots: runtime.Stores.Snapshots, Sessions: runtime.Stores.Sessions}
	view := service.ViewFlow{Mooncakes: runtime.Sources.Mooncakes, Assets: runtime.Sources.Assets, Sessions: runtime.Stores.Sessions}
	api := service.APIFlow{Mooncakes: runtime.Sources.Mooncakes, Assets: runtime.Sources.Assets, Sessions: runtime.Stores.Sessions}
	sourceFlow := service.SourceFlow{Mooncakes: runtime.Sources.Mooncakes, Assets: runtime.Sources.Assets, MoonCLI: runtime.Sources.MoonCLI, LocalCache: runtime.Sources.LocalCache, Sessions: runtime.Stores.Sessions, Paths: runtime.Paths}
	skill := service.SkillFlow{Skills: runtime.Sources.Skills, Sessions: runtime.Stores.Sessions}
	compare := service.CompareFlow{Mooncakes: runtime.Sources.Mooncakes, Sessions: runtime.Stores.Sessions}
	probe := service.ProbeFlow{MoonCLI: runtime.Sources.MoonCLI, Sessions: runtime.Stores.Sessions, Probes: runtime.Stores.Probes, Paths: runtime.Paths}
	record := service.RecordFlow{Sessions: runtime.Stores.Sessions, Records: runtime.Stores.Records, Project: runtime.Sources.Project, Paths: runtime.Paths, Now: func() time.Time { return runtime.Clock.Now() }}
	report := service.ReportFlow{Records: runtime.Stores.Records, Reports: runtime.Stores.Reports, Project: runtime.Sources.Project, Paths: runtime.Paths, Now: func() time.Time { return runtime.Clock.Now() }}
	raw := service.RawFlow{Mooncakes: runtime.Sources.Mooncakes, Assets: runtime.Sources.Assets, Skills: runtime.Sources.Skills}
	runtime.Services = ServiceRegistry{
		Sync: syncService,
	}
	runtime.Services.Registry = service.Registry{
		Doctor:  doctor,
		Sync:    syncService,
		Search:  search,
		View:    view,
		API:     api,
		Source:  sourceFlow,
		Skill:   skill,
		Compare: compare,
		Probe:   probe,
		Record:  record,
		Report:  report,
		Raw:     raw,
	}
	runtime.Renderer = output.DefaultRenderer{}
}
