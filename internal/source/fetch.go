// Package source reads external and local evidence into state-bearing results.
package source

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

type Fetcher struct {
	Client *http.Client
	Clock  platform.Clock
}

type FetchResult struct {
	URL        string
	Path       string
	FetchedAt  string
	Status     model.State
	ParseState model.State
	RawRef     string
	Body       []byte
	Error      string
}

func (f Fetcher) Fetch(ctx context.Context, url string) FetchResult {
	var out FetchResult
	operation := func() (FetchResult, error) {
		out = f.fetchOnce(ctx, url)
		if out.Status == model.StateFailed {
			return out, fmt.Errorf("%s", out.Error)
		}
		return out, nil
	}
	result, _ := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewConstantBackOff(10*time.Millisecond)), backoff.WithMaxTries(2))
	if result.URL != "" {
		return result
	}
	return out
}

func (f Fetcher) FetchJSON(ctx context.Context, url string, target any) FetchResult {
	result := f.Fetch(ctx, url)
	if result.Status != model.StatePresent {
		return result
	}
	if err := json.NewDecoder(bytes.NewReader(result.Body)).Decode(target); err != nil {
		result.ParseState = model.StateFailed
		result.Status = model.StateFailed
		result.Error = err.Error()
		return result
	}
	result.ParseState = model.StatePresent
	return result
}

func (f Fetcher) fetchOnce(ctx context.Context, url string) FetchResult {
	now := time.Now()
	if f.Clock != nil {
		now = f.Clock.Now()
	}
	result := FetchResult{
		URL:        url,
		FetchedAt:  now.Format(time.RFC3339),
		ParseState: model.StateUnknown,
		RawRef:     "inline:" + url,
	}
	client := f.Client
	if client == nil {
		client = platform.NewHTTPClient(30 * time.Second)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		result.Status = model.StateFailed
		result.Error = err.Error()
		return result
	}
	resp, err := client.Do(req)
	if err != nil {
		result.Status = model.StateFailed
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	result.Body = body
	if readErr != nil {
		result.Status = model.StateFailed
		result.Error = readErr.Error()
		return result
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusNotFound {
			result.Status = model.StateUnavailable
		} else {
			result.Status = model.StateFailed
		}
		result.Error = fmt.Sprintf("http status %d", resp.StatusCode)
		return result
	}
	result.Status = model.StatePresent
	return result
}

func SourceResult[T any](label string, fetch FetchResult, value *T) model.SourceResult[T] {
	return model.SourceResult[T]{
		Source:     label,
		URL:        fetch.URL,
		Path:       fetch.Path,
		FetchedAt:  fetch.FetchedAt,
		Status:     fetch.Status,
		ParseState: fetch.ParseState,
		RawRef:     fetch.RawRef,
		Error:      fetch.Error,
		Value:      value,
	}
}

func RawJSONSourceResult(label string, fetch FetchResult) model.SourceResult[any] {
	if fetch.Status != model.StatePresent {
		return SourceResult[any](label, fetch, nil)
	}
	var raw any
	decoder := json.NewDecoder(bytes.NewReader(fetch.Body))
	decoder.UseNumber()
	if err := decoder.Decode(&raw); err != nil {
		fetch.Status = model.StateFailed
		fetch.ParseState = model.StateFailed
		fetch.Error = err.Error()
		return SourceResult[any](label, fetch, nil)
	}
	fetch.ParseState = model.StatePresent
	return SourceResult(label, fetch, &raw)
}

func sourceNow(clock platform.Clock) string {
	now := time.Now()
	if clock != nil {
		now = clock.Now()
	}
	return now.Format(time.RFC3339)
}
