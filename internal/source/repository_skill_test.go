package source

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/platform"
	"github.com/google/go-github/v88/github"
)

func TestRepositoryURLParsingAndUnsupported(t *testing.T) {
	coord, err := ParseGitHubRepository("https://github.com/moonbitlang/core.git")
	if err != nil {
		t.Fatalf("ParseGitHubRepository: %v", err)
	}
	if coord.Owner != "moonbitlang" || coord.Name != "core" {
		t.Fatalf("coord = %#v", coord)
	}
	result := RepositoryReader{}.Signal(context.Background(), "https://example.com/owner/repo")
	if result.Status != model.StateUnknown || result.Value == nil || result.Value.Status != model.StateUnknown {
		t.Fatalf("unsupported repository result = %#v", result)
	}
}

func TestRepositorySignalMapping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"archived":true,"pushed_at":"2026-06-21T00:00:00Z","license":{"spdx_id":"MIT"},"open_issues_count":3,"default_branch":"main"}`))
	}))
	defer server.Close()
	base := server.URL + "/"
	client, err := github.NewClient(github.WithHTTPClient(server.Client()), github.WithURLs(&base, &base))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	result := RepositoryReader{Client: client, Clock: platform.FixedClock{Time: time.Date(2026, 6, 24, 1, 2, 3, 0, time.UTC)}}.Signal(context.Background(), "https://github.com/owner/repo")
	if result.Status != model.StatePresent || result.Value == nil {
		t.Fatalf("repository result = %#v", result)
	}
	if result.FetchedAt != "2026-06-24T01:02:03Z" || result.ParseState != model.StatePresent || result.RawRef == "" {
		t.Fatalf("repository envelope = %#v", result)
	}
	if result.Value.IsArchived.Value == nil || !*result.Value.IsArchived.Value || result.Value.License.Value == nil || *result.Value.License.Value != "MIT" {
		t.Fatalf("repository signal = %#v", result.Value)
	}
}

func TestRepositorySignalMissingPushedAtUnknown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"archived":false,"open_issues_count":0,"default_branch":"main"}`))
	}))
	defer server.Close()
	base := server.URL + "/"
	client, err := github.NewClient(github.WithHTTPClient(server.Client()), github.WithURLs(&base, &base))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	result := RepositoryReader{Client: client}.Signal(context.Background(), "https://github.com/owner/repo")
	if result.Status != model.StatePresent || result.Value == nil {
		t.Fatalf("repository result = %#v", result)
	}
	if result.Value.PushedAt.Status != model.StateUnknown {
		t.Fatalf("pushed_at evidence = %#v", result.Value.PushedAt)
	}
}

func TestSkillParsingRootMarkerAndRunwasm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"module":"Yoorkin/cowsay","author":"Yoorkin","author_avatar":"","version":"0.1.0","package":"","name":"cowsay","wasm_url":"https://asset/wasm","checksum_url":"https://asset/checksum","repository":"","created_at":"2026"}]`))
	}))
	defer server.Close()
	client := SkillsClient{BaseURL: server.URL, Fetcher: Fetcher{Client: server.Client()}}
	result := client.FetchSkills(context.Background())
	if result.Status != model.StatePresent || result.Value == nil || len(*result.Value) != 1 {
		t.Fatalf("skill result = %#v", result)
	}
	entry := (*result.Value)[0]
	if entry.Package != "" {
		t.Fatalf("root marker package = %q", entry.Package)
	}
	if got := model.RunwasmCoordinate(entry); got != "Yoorkin/cowsay@0.1.0" {
		t.Fatalf("runwasm coordinate = %q", got)
	}
	if entry.AuthorAvatar.Status != model.StateMissing || entry.Repository.Status != model.StateMissing {
		t.Fatalf("skill missing evidence = %#v", entry)
	}
}

func TestSkillAssetFetchStates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte("asset"))
	}))
	defer server.Close()
	client := SkillsClient{Fetcher: Fetcher{Client: server.Client()}}
	if got := client.FetchSkillAsset(context.Background(), server.URL+"/asset"); got.Status != model.StatePresent {
		t.Fatalf("asset state = %#v", got)
	}
	if got := client.FetchSkillAsset(context.Background(), server.URL+"/missing"); got.Status != model.StateUnavailable {
		t.Fatalf("missing asset state = %#v", got)
	}
}
