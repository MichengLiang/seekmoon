// Package service contains command use-case skeletons.
package service

import (
	"context"
	"time"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

type SyncService struct {
	Mooncakes source.MooncakesClient
	Snapshots store.SnapshotStore
	Now       func() time.Time
}

func (s SyncService) Sync(ctx context.Context) (model.Snapshot, error) {
	now := time.Now()
	if s.Now != nil {
		now = s.Now()
	}
	modules := s.Mooncakes.FetchModules(ctx)
	stats := s.Mooncakes.FetchStatistics(ctx)
	snapshot := model.Snapshot{
		ID:        now.Format(time.RFC3339),
		CreatedAt: now.Format(time.RFC3339),
		Sources: []model.SourceResult[any]{
			eraseSourceValue(modules),
			eraseSourceValue(stats),
		},
	}
	if stats.Value != nil {
		snapshot.Statistics = *stats.Value
	}
	err := s.Snapshots.Write(ctx, snapshot)
	return snapshot, err
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
