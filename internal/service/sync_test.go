package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

func TestSyncRecordsPartialSourceFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v0/modules":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[{"name":"mizchi/markdown","version":"0.6.2"}]`))
		case "/api/v0/modules/statistics":
			http.Error(w, "down", http.StatusBadGateway)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	service := SyncService{
		Mooncakes: source.MooncakesClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		Snapshots: store.SnapshotStore{FS: platform.OSFS{}, Paths: paths},
		Now:       func() time.Time { return time.Date(2026, 6, 24, 1, 2, 3, 0, time.UTC) },
	}
	snapshot, err := service.Sync(context.Background())
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if len(snapshot.Sources) != 2 {
		t.Fatalf("sources = %#v", snapshot.Sources)
	}
	if snapshot.Sources[0].Status != model.StatePresent || snapshot.Sources[1].Status != model.StateFailed {
		t.Fatalf("partial source states = %#v", snapshot.Sources)
	}
	read, err := service.Snapshots.Read(context.Background(), snapshot.ID)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	if len(read.Sources) != 2 {
		t.Fatalf("stored sources = %#v", read.Sources)
	}
}
