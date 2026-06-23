package source

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v88/github"
	"golang.org/x/oauth2"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

// RepositoryReader reads repository evidence from GitHub.
type RepositoryReader struct {
	Client *github.Client
	Token  string
	Clock  platform.Clock
}

// RepositoryCoordinate identifies a GitHub owner/repository pair.
type RepositoryCoordinate struct {
	Owner string
	Name  string
}

// ParseGitHubRepository parses supported GitHub repository URLs.
func ParseGitHubRepository(rawURL string) (RepositoryCoordinate, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return RepositoryCoordinate{}, err
	}
	if parsed.Host != "github.com" {
		return RepositoryCoordinate{}, fmt.Errorf("unsupported repository host %q", parsed.Host)
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return RepositoryCoordinate{}, fmt.Errorf("repository URL must contain owner/name")
	}
	return RepositoryCoordinate{Owner: parts[0], Name: strings.TrimSuffix(parts[1], ".git")}, nil
}

// Signal reads repository maintenance signals.
func (r RepositoryReader) Signal(ctx context.Context, rawURL string) model.SourceResult[model.RepositorySignal] {
	envelope := r.envelope(rawURL, model.StateUnknown, model.StateUnknown, "")
	coord, err := ParseGitHubRepository(rawURL)
	if err != nil {
		value := model.RepositorySignal{URL: rawURL, Status: model.StateUnknown, Error: err.Error()}
		envelope.Error = err.Error()
		envelope.Value = &value
		return envelope
	}
	envelope.RawRef = fmt.Sprintf("github:%s/%s", coord.Owner, coord.Name)
	client := r.Client
	if client == nil {
		httpClient := http.DefaultClient
		if r.Token != "" {
			httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: r.Token}))
		}
		client, err = github.NewClient(github.WithHTTPClient(httpClient))
		if err != nil {
			value := model.RepositorySignal{URL: rawURL, Status: model.StateFailed, Error: err.Error()}
			envelope.Status = model.StateFailed
			envelope.ParseState = model.StateFailed
			envelope.Error = err.Error()
			envelope.Value = &value
			return envelope
		}
	}
	repo, _, err := client.Repositories.Get(ctx, coord.Owner, coord.Name)
	if err != nil {
		value := model.RepositorySignal{URL: rawURL, Status: model.StateFailed, Error: err.Error()}
		envelope.Status = model.StateFailed
		envelope.ParseState = model.StateFailed
		envelope.Error = err.Error()
		envelope.Value = &value
		return envelope
	}
	value := model.RepositorySignal{
		URL:           rawURL,
		Status:        model.StatePresent,
		IsArchived:    model.Present(repo.GetArchived(), string(model.SourceRepositoryAPI)),
		PushedAt:      repositoryPushedAt(repo),
		License:       evidenceString(repo.GetLicense().GetSPDXID(), string(model.SourceRepositoryAPI)),
		HasReleases:   model.Unknown[bool](),
		OpenIssues:    model.Present(repo.GetOpenIssuesCount(), string(model.SourceRepositoryAPI)),
		OpenPulls:     model.Unknown[int](),
		HasWorkflows:  model.Unknown[bool](),
		HasReadme:     model.Unknown[bool](),
		HasTests:      model.Unknown[bool](),
		HasExamples:   model.Unknown[bool](),
		DefaultBranch: evidenceString(repo.GetDefaultBranch(), string(model.SourceRepositoryAPI)),
	}
	envelope.Status = model.StatePresent
	envelope.ParseState = model.StatePresent
	envelope.Value = &value
	return envelope
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"

func repositoryPushedAt(repo *github.Repository) model.EvidenceString {
	if repo.PushedAt == nil {
		return model.Unknown[string]()
	}
	return model.Present(repo.GetPushedAt().Format(timeRFC3339), string(model.SourceRepositoryAPI))
}

func (r RepositoryReader) envelope(rawURL string, status, parseState model.State, err string) model.SourceResult[model.RepositorySignal] {
	return model.SourceResult[model.RepositorySignal]{
		Source:     string(model.SourceRepositoryAPI),
		URL:        rawURL,
		FetchedAt:  sourceNow(r.Clock),
		Status:     status,
		ParseState: parseState,
		RawRef:     "repository:" + rawURL,
		Error:      err,
	}
}
