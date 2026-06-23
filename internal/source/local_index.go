package source

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

type LocalIndexReader struct {
	FS    platform.FS
	Clock platform.Clock
}

type LocalIndexSummary struct {
	Records   []model.ModuleSummary `json:"records"`
	Malformed int                   `json:"malformed"`
}

func (r LocalIndexReader) Parse(data []byte) LocalIndexSummary {
	var summary LocalIndexSummary
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Bytes()
		var item moduleItem
		if err := json.Unmarshal(line, &item); err != nil {
			summary.Malformed++
			continue
		}
		if item.Name == "" {
			summary.Malformed++
			continue
		}
		normalized := normalizeModule(item)
		normalized.Raw["raw_line"] = string(line)
		summary.Records = append(summary.Records, normalized)
	}
	return summary
}

func (r LocalIndexReader) Read(ctx context.Context, path string) model.SourceResult[LocalIndexSummary] {
	fs := r.FS
	if fs == nil {
		fs = platform.OSFS{}
	}
	data, err := fs.ReadFile(ctx, path)
	cleanPath := filepath.Clean(path)
	if err != nil {
		return model.SourceResult[LocalIndexSummary]{Source: string(model.SourceLocalIndex), Path: cleanPath, FetchedAt: sourceNow(r.Clock), Status: model.StateFailed, ParseState: model.StateFailed, RawRef: fmt.Sprintf("file:%s", cleanPath), Error: err.Error()}
	}
	value := r.Parse(data)
	return model.SourceResult[LocalIndexSummary]{Source: string(model.SourceLocalIndex), Path: cleanPath, FetchedAt: sourceNow(r.Clock), Status: model.StatePresent, ParseState: model.StatePresent, RawRef: fmt.Sprintf("file:%s", cleanPath), Value: &value}
}
