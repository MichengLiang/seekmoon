package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

type APIFlow struct {
	Mooncakes source.MooncakesClient
	Assets    source.AssetClient
	Sessions  store.SessionStore
}

func (s APIFlow) API(ctx context.Context, input APIInput) (model.PackageData, error) {
	candidate, err := candidateFromRequest(ctx, s.Sessions, input.Candidate)
	if err != nil {
		return model.PackageData{}, err
	}
	version := candidate.Version
	if version == "" {
		manifest := s.Mooncakes.FetchManifest(ctx, candidate.Module)
		if manifest.Status != model.StatePresent || manifest.Value == nil {
			return model.PackageData{}, sourceFailure(candidate.Module, eraseSourceValue(manifest), "check the module coordinate or run view first")
		}
		version = manifest.Value.Version
	}
	index := s.Assets.FetchModuleIndex(ctx, candidate.Module, version)
	if index.Status != model.StatePresent || index.Value == nil {
		return model.PackageData{}, sourceFailure(candidate.Module, eraseSourceValue(index), "retry api after module index is available")
	}
	if !hasPackagePath(*index.Value, input.Package) {
		known := strings.Join(packagePaths(*index.Value), "\n  ")
		return model.PackageData{}, surfaceFailure(
			fmt.Sprintf("%s %s", candidate.Module, input.Package),
			string(model.SourceModuleIndex),
			model.StateFailed,
			"package path is not present in module index\n\nknown packages\n  "+known,
			"choose one of the known package paths",
		)
	}
	data := s.Assets.FetchPackageData(ctx, candidate.Module, version, input.Package)
	if data.Status != model.StatePresent || data.Value == nil {
		return model.PackageData{}, sourceFailure(input.Package, eraseSourceValue(data), "check package path and retry")
	}
	return *data.Value, nil
}
