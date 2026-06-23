package service

import (
	"context"
	"fmt"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

type ViewFlow struct {
	Mooncakes source.MooncakesClient
	Assets    source.AssetClient
	Sessions  store.SessionStore
}

func (s ViewFlow) View(ctx context.Context, input ViewInput) (model.ManifestProfile, error) {
	candidate, err := candidateFromRequest(ctx, s.Sessions, input.Candidate)
	if err != nil {
		return model.ManifestProfile{}, err
	}
	manifest := s.Mooncakes.FetchManifest(ctx, candidate.Module)
	if manifest.Status != model.StatePresent || manifest.Value == nil {
		return model.ManifestProfile{}, sourceFailure(candidate.Module, eraseSourceValue(manifest), "check the module coordinate or retry view")
	}
	profile := *manifest.Value
	index := s.Assets.FetchModuleIndex(ctx, profile.Module, profile.Version)
	if index.Status == model.StatePresent && index.Value != nil {
		profile.Metadata.Raw["package_count"] = len(packagePaths(*index.Value))
		profile.Metadata.Raw["module_index_state"] = string(index.Status)
	} else {
		profile.Metadata.Raw["module_index_state"] = string(index.Status)
		if index.Error != "" {
			profile.Metadata.Raw["module_index_error"] = index.Error
		}
	}
	if profile.Module == "" {
		return model.ManifestProfile{}, fmt.Errorf("manifest profile missing module")
	}
	return profile, nil
}
