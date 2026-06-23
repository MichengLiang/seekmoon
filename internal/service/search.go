package service

import (
	"context"
	"os"
	"sort"
	"strings"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

type SearchFlow struct {
	Mooncakes source.MooncakesClient
	Assets    source.AssetClient
	Snapshots store.SnapshotStore
	Sessions  store.SessionStore
}

func (s SearchFlow) Search(ctx context.Context, input SearchInput) (model.SearchOutput, error) {
	snapshot, modules, err := snapshotModules(ctx, s.Snapshots)
	if err != nil {
		if err != os.ErrNotExist {
			return model.SearchOutput{}, err
		}
		fetched := s.Mooncakes.FetchModules(ctx)
		if fetched.Status != model.StatePresent || fetched.Value == nil {
			return model.SearchOutput{}, sourceFailure("search modules", eraseSourceValue(fetched), "run seekmoon sync or retry search")
		}
		modules = *fetched.Value
		snapshot = model.Snapshot{
			ID:        "transient",
			CreatedAt: fetched.FetchedAt,
			Sources:   []model.SourceResult[any]{eraseSourceValue(fetched)},
		}
	}
	query := strings.ToLower(strings.TrimSpace(input.Query))
	matches := make([]model.SearchResult, 0, len(modules))
	for _, module := range modules {
		fields := matchFields(module, query)
		if len(fields) == 0 {
			continue
		}
		result := model.SearchResult{
			Module:      module.Module,
			Version:     module.Version,
			Description: module.Description,
			Keywords:    module.Keywords,
			License:     module.License,
			Repository:  module.Repository,
			Target:      model.Unknown[map[string]any](),
			Match:       model.MatchEvidence{Fields: fields, Token: input.Query},
			SnapshotID:  snapshot.ID,
		}
		if input.Target != "" {
			result.Target = model.Unknown[map[string]any]()
		}
		matches = append(matches, result)
	}
	sort.SliceStable(matches, func(i, j int) bool {
		if len(matches[i].Match.Fields) == len(matches[j].Match.Fields) {
			return matches[i].Module < matches[j].Module
		}
		return len(matches[i].Match.Fields) > len(matches[j].Match.Fields)
	})
	for i := range matches {
		matches[i].Rank = i + 1
		enrichment := s.Mooncakes.FetchManifest(ctx, matches[i].Module)
		if enrichment.Status == model.StatePresent && enrichment.Value != nil {
			matches[i].Downloads = model.Present(enrichment.Value.Downloads, string(model.SourceManifestAPI))
			matches[i].BuildStatus = enrichment.Value.BuildStatus
			continue
		}
		matches[i].Downloads = model.Failed[int](string(model.SourceManifestAPI), firstNonEmpty(enrichment.Error, "manifest enrichment failed"))
		matches[i].BuildStatus = model.Failed[string](string(model.SourceManifestAPI), firstNonEmpty(enrichment.Error, "manifest enrichment failed"))
	}
	candidates := store.SessionCandidates{SessionID: defaultSessionID, SnapshotID: snapshot.ID, Candidates: map[int]model.CandidateRef{}}
	for _, result := range matches {
		candidates.Candidates[result.Rank] = model.CandidateRef{Kind: "library", Module: result.Module, Version: result.Version}
	}
	if err := s.Sessions.WriteCandidates(ctx, candidates); err != nil {
		return model.SearchOutput{}, err
	}
	return model.SearchOutput{
		Schema:   model.SchemaSearchResultsV1,
		Snapshot: model.SnapshotRef{ID: snapshot.ID, Sources: sourceLabels(snapshot.Sources)},
		Query:    model.SearchQuery{Text: input.Query, Kind: "library", Target: input.Target},
		Results:  matches,
	}, nil
}

func matchFields(module model.ModuleSummary, query string) []string {
	var fields []string
	if query == "" {
		return fields
	}
	if strings.Contains(strings.ToLower(module.Module), query) {
		fields = append(fields, "module")
	}
	if module.Description.Value != nil && strings.Contains(strings.ToLower(*module.Description.Value), query) {
		fields = append(fields, "description")
	}
	if module.Keywords.Value != nil {
		for _, keyword := range *module.Keywords.Value {
			if strings.Contains(strings.ToLower(keyword), query) {
				fields = append(fields, "keywords")
				break
			}
		}
	}
	if module.Repository.Value != nil && strings.Contains(strings.ToLower(*module.Repository.Value), query) {
		fields = append(fields, "repository")
	}
	return fields
}

func sourceLabels(sources []model.SourceResult[any]) []string {
	labels := make([]string, 0, len(sources))
	for _, source := range sources {
		labels = append(labels, source.Source)
	}
	return labels
}
