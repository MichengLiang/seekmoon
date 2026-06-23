package service

import (
	"context"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

// CompareFlow compares evidence for selected candidates.
type CompareFlow struct {
	Mooncakes source.MooncakesClient
	Sessions  store.SessionStore
}

// Compare returns field-by-field evidence for candidate comparison.
func (s CompareFlow) Compare(ctx context.Context, input CompareInput) (any, error) {
	refs := make([]model.CandidateRef, 0, len(input.Candidates))
	fields := map[string]map[string]string{
		"version": {},
		"license": {},
		"build":   {},
		"repo":    {},
		"target":  {},
	}
	for _, request := range input.Candidates {
		ref, err := candidateFromRequest(ctx, s.Sessions, request)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
		manifest := s.Mooncakes.FetchManifest(ctx, ref.Module)
		key := ref.Module
		if manifest.Status == model.StatePresent && manifest.Value != nil {
			fields["version"][key] = manifest.Value.Version
			fields["license"][key] = evidenceStringValue(manifest.Value.Metadata.License)
			fields["build"][key] = evidenceStringValue(manifest.Value.BuildStatus)
			fields["repo"][key] = evidenceStringValue(manifest.Value.Metadata.Repository)
		} else {
			fields["version"][key] = string(manifest.Status)
			fields["license"][key] = string(manifest.Status)
			fields["build"][key] = string(manifest.Status)
			fields["repo"][key] = string(manifest.Status)
		}
		fields["target"][key] = "unknown"
	}
	return model.Comparison{
		Schema:     model.SchemaComparisonV1,
		Candidates: refs,
		Fields: []model.ComparisonField{
			{Name: "version", Values: fields["version"]},
			{Name: "license", Values: fields["license"]},
			{Name: "build", Values: fields["build"]},
			{Name: "repo", Values: fields["repo"]},
			{Name: "target", Values: fields["target"]},
		},
	}, nil
}
