package service

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

// SourceFlow resolves source acquisition options for a candidate.
type SourceFlow struct {
	Mooncakes  source.MooncakesClient
	Assets     source.AssetClient
	MoonCLI    source.MoonCLI
	LocalCache source.LocalCacheReader
	Sessions   store.SessionStore
	Paths      store.Paths
}

// Source returns attempted source locations and the selected inspection route.
func (s SourceFlow) Source(ctx context.Context, input SourceInput) (model.SourceResolution, error) {
	candidate, err := candidateFromRequest(ctx, s.Sessions, input.Candidate)
	if err != nil {
		return model.SourceResolution{}, err
	}
	version := candidate.Version
	manifest := s.Mooncakes.FetchManifest(ctx, candidate.Module)
	if version == "" && manifest.Status == model.StatePresent && manifest.Value != nil {
		version = manifest.Value.Version
	}
	if version == "" {
		version = "latest"
	}
	moduleVersion := candidate.Module + "@" + version
	moonFetch := model.SourceAttempt{Status: model.StateUnknown}
	if s.MoonCLI.Runner != nil {
		result := s.MoonCLI.Fetch(ctx, s.Paths.ProjectRoot, moduleVersion)
		moonFetch = model.SourceAttempt{Status: result.Status, Path: filepath.Join(s.Paths.ProjectRoot, ".repos"), Error: commandError(result)}
	}
	sourceZip := s.Assets.FetchSourceZipAttempt(ctx, candidate.Module, version)
	localCache := s.LocalCache.CheckPath(ctx, model.SourceLocalCache, filepath.Join(s.Paths.Sources, store.SafeName(moduleVersion)))
	core := model.SourceAttempt{Status: model.StateUnavailable, Path: filepath.Join(s.Paths.ProjectRoot, ".moon", "lib", "core")}
	if candidate.Module == "moonbitlang/core" {
		core = s.LocalCache.CheckPath(ctx, model.SourceCoreLocalSource, core.Path)
	}
	repository := model.SourceAttempt{Status: model.StateUnknown}
	if manifest.Status == model.StatePresent && manifest.Value != nil && manifest.Value.Metadata.Repository.Value != nil {
		repository = model.SourceAttempt{Status: model.StatePresent, URL: *manifest.Value.Metadata.Repository.Value}
	}
	resolution := model.SourceResolution{
		Module:           candidate.Module,
		Version:          version,
		MoonFetch:        moonFetch,
		SourceZip:        sourceZip,
		LocalCache:       localCache,
		CoreLocalSource:  core,
		RepositorySource: repository,
		SelectedSource:   selectedSource(moonFetch, sourceZip, localCache, core, repository),
		FilesSummary:     filesSummaryFromAttempt(sourceZip),
	}
	return resolution, nil
}

func selectedSource(attempts ...model.SourceAttempt) model.SelectedSource {
	methods := []string{"moon_fetch", "source_zip", "local_cache", "core_local_source", "repository_signal"}
	for i, attempt := range attempts {
		if attempt.Status != model.StatePresent {
			continue
		}
		return model.SelectedSource{Method: methods[i], Path: attempt.Path, URL: attempt.URL}
	}
	return model.SelectedSource{Method: "none"}
}

func commandError(result model.CommandResult) string {
	if result.Status == model.StateFailed {
		return "command exited non-zero"
	}
	return ""
}

func filesSummaryFromAttempt(attempt model.SourceAttempt) model.FilesSummary {
	if attempt.Path == "" {
		return model.FilesSummary{}
	}
	var summary model.FilesSummary
	_ = json.Unmarshal([]byte(attempt.Path), &summary)
	return summary
}
