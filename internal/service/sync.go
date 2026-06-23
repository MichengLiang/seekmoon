// Package service contains command use-case skeletons.
package service

import (
	"context"
	"time"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

// SyncService builds and stores source snapshots.
type SyncService struct {
	Mooncakes  source.MooncakesClient
	MoonCLI    source.MoonCLI
	LocalIndex source.LocalIndexReader
	Snapshots  store.SnapshotStore
	Paths      store.Paths
	Now        func() time.Time
}

// Sync fetches current source evidence and writes a snapshot.
func (s SyncService) Sync(ctx context.Context) (model.Snapshot, error) {
	now := time.Now()
	if s.Now != nil {
		now = s.Now()
	}
	moonUpdate := unknownCommand("moon", "update")
	moonVersion := unknownCommand("moon", "--version")
	if s.MoonCLI.Runner != nil {
		moonUpdate = s.MoonCLI.Update(ctx, s.Paths.ProjectRoot)
		moonVersion = s.MoonCLI.Version(ctx, s.Paths.ProjectRoot)
	}
	modules := s.Mooncakes.FetchModules(ctx)
	stats := s.Mooncakes.FetchStatistics(ctx)
	localIndex := s.LocalIndex.Read(ctx, s.Paths.MoonIndex)
	snapshot := model.Snapshot{
		ID:        now.Format(time.RFC3339),
		CreatedAt: now.Format(time.RFC3339),
		Sources: []model.SourceResult[any]{
			eraseSourceValue(modules),
			eraseSourceValue(stats),
			eraseSourceValue(localIndex),
		},
		Raw: map[string]any{
			"moon_update":  moonUpdate,
			"moon_version": moonVersion,
			"local_index":  localIndexSummary(localIndex),
		},
	}
	if modules.Value != nil {
		snapshot.Raw["modules"] = modules.Value
	}
	if stats.Value != nil {
		snapshot.Raw["statistics"] = stats.Value
	}
	if stats.Value != nil {
		snapshot.Statistics = *stats.Value
	}
	err := s.Snapshots.Write(ctx, snapshot)
	return snapshot, err
}

func localIndexSummary(result model.SourceResult[source.LocalIndexSummary]) map[string]any {
	summary := map[string]any{
		"status":      result.Status,
		"parse_state": result.ParseState,
		"path":        result.Path,
		"raw_ref":     result.RawRef,
	}
	if result.Error != "" {
		summary["error"] = result.Error
	}
	if result.Value != nil {
		summary["index_head"] = result.Value.IndexHead
		summary["file_count"] = result.Value.FileCount
		summary["record_count"] = result.Value.RecordCount
		summary["malformed"] = result.Value.Malformed
	}
	return summary
}

func unknownCommand(command ...string) model.CommandResult {
	return model.CommandResult{Command: command, Status: model.StateUnknown}
}

func eraseSourceValue[T any](input model.SourceResult[T]) model.SourceResult[any] {
	var value *any
	if input.Value != nil {
		v := any(*input.Value)
		value = &v
	}
	return model.SourceResult[any]{
		Source:     input.Source,
		URL:        input.URL,
		Path:       input.Path,
		FetchedAt:  input.FetchedAt,
		Status:     input.Status,
		ParseState: input.ParseState,
		RawRef:     input.RawRef,
		Error:      input.Error,
		Value:      value,
	}
}
