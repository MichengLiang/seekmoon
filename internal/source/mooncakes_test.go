package source

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

func TestMooncakesModulesAPIMissingEvidence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/modules" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"name":"mizchi/markdown","version":"0.6.2","description":"","keywords":[],"repository":"","license":"","is_new":true,"created_at":"2026-06-21"}]`))
	}))
	defer server.Close()

	client := MooncakesClient{
		BaseURL: server.URL,
		Fetcher: Fetcher{
			Client: server.Client(),
			Clock:  platform.FixedClock{Time: time.Date(2026, 6, 24, 1, 2, 3, 0, time.UTC)},
		},
	}
	result := client.FetchModules(context.Background())
	if result.Status != model.StatePresent || result.ParseState != model.StatePresent {
		t.Fatalf("source states = %s/%s error=%s", result.Status, result.ParseState, result.Error)
	}
	if result.Value == nil || len(*result.Value) != 1 {
		t.Fatalf("modules value = %#v", result.Value)
	}
	module := (*result.Value)[0]
	if module.Description.Status != model.StateMissing || module.Keywords.Status != model.StateMissing || module.Repository.Status != model.StateMissing || module.License.Status != model.StateMissing {
		t.Fatalf("empty module fields did not map to missing: %#v", module)
	}
	if result.Source != string(model.SourceModulesAPI) || result.URL == "" || result.FetchedAt == "" || result.RawRef == "" {
		t.Fatalf("source envelope incomplete: %#v", result)
	}
}

func TestMooncakesFetchNon2xxAndParseFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v0/modules":
			http.Error(w, "boom", http.StatusInternalServerError)
		case "/api/v0/modules/statistics":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"total_modules":`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := MooncakesClient{BaseURL: server.URL, Fetcher: Fetcher{Client: server.Client()}}
	modules := client.FetchModules(context.Background())
	if modules.Status != model.StateFailed || !strings.Contains(modules.Error, "500") {
		t.Fatalf("non-2xx result = %#v", modules)
	}
	stats := client.FetchStatistics(context.Background())
	if stats.Status != model.StateFailed || stats.ParseState != model.StateFailed {
		t.Fatalf("parse failure result = %#v", stats)
	}
}

func TestMooncakesManifestMismatchFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"module":"owner/a","name":"owner/b","version":"1.0.0","metadata":{}}`))
	}))
	defer server.Close()

	client := MooncakesClient{BaseURL: server.URL, Fetcher: Fetcher{Client: server.Client()}}
	result := client.FetchManifest(context.Background(), "owner/a")
	if result.Status != model.StateFailed || result.ParseState != model.StateFailed {
		t.Fatalf("manifest mismatch result = %#v", result)
	}
}
