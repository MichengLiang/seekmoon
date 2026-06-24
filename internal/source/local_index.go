package source

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/platform"
)

// LocalIndexReader reads Moon registry local index files.
type LocalIndexReader struct {
	FS    platform.FS
	Clock platform.Clock
}

// LocalIndexSummary summarizes records parsed from local index files.
type LocalIndexSummary struct {
	Records     []model.ModuleSummary `json:"records"`
	FileCount   int                   `json:"file_count"`
	RecordCount int                   `json:"record_count"`
	Malformed   int                   `json:"malformed"`
	IndexHead   string                `json:"index_head,omitempty"`
}

// Parse parses newline-delimited local index records.
func (r LocalIndexReader) Parse(data []byte) LocalIndexSummary {
	var summary LocalIndexSummary
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		summary.RecordCount++
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
	cleanPath := filepath.Clean(path)
	if err := ctx.Err(); err != nil {
		return model.SourceResult[LocalIndexSummary]{Source: string(model.SourceLocalIndex), Path: cleanPath, FetchedAt: sourceNow(r.Clock), Status: model.StateFailed, ParseState: model.StateFailed, RawRef: fmt.Sprintf("file:%s", cleanPath), Error: err.Error()}
	}
	if _, ok := fs.(platform.OSFS); ok {
		stat, err := os.Stat(cleanPath)
		if err != nil {
			status := model.StateFailed
			if errors.Is(err, os.ErrNotExist) {
				status = model.StateUnavailable
			}
			return model.SourceResult[LocalIndexSummary]{Source: string(model.SourceLocalIndex), Path: cleanPath, FetchedAt: sourceNow(r.Clock), Status: status, ParseState: status, RawRef: fmt.Sprintf("file:%s", cleanPath), Error: err.Error()}
		}
		if stat.IsDir() {
			return r.readDirectory(ctx, cleanPath)
		}
	}
	data, err := fs.ReadFile(ctx, cleanPath)
	if err != nil {
		status := model.StateFailed
		if errors.Is(err, os.ErrNotExist) {
			status = model.StateUnavailable
		}
		return model.SourceResult[LocalIndexSummary]{Source: string(model.SourceLocalIndex), Path: cleanPath, FetchedAt: sourceNow(r.Clock), Status: status, ParseState: status, RawRef: fmt.Sprintf("file:%s", cleanPath), Error: err.Error()}
	}
	value := r.Parse(data)
	value.FileCount = 1
	return model.SourceResult[LocalIndexSummary]{Source: string(model.SourceLocalIndex), Path: cleanPath, FetchedAt: sourceNow(r.Clock), Status: model.StatePresent, ParseState: model.StatePresent, RawRef: fmt.Sprintf("file:%s", cleanPath), Value: &value}
}

func (r LocalIndexReader) readDirectory(ctx context.Context, root string) model.SourceResult[LocalIndexSummary] {
	var summary LocalIndexSummary
	cleanRoot := filepath.Clean(root)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".index" {
			return nil
		}
		cleanPath := filepath.Clean(path)
		rel, err := filepath.Rel(cleanRoot, cleanPath)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return fmt.Errorf("local index path escaped root: %s", path)
		}
		// #nosec G304,G122 -- WalkDir supplies paths under the configured local
		// index root; rel validation rejects callbacks outside that root before
		// reading discovered registry index files.
		data, err := os.ReadFile(cleanPath)
		if err != nil {
			return err
		}
		fileSummary := r.Parse(data)
		summary.FileCount++
		summary.RecordCount += fileSummary.RecordCount
		summary.Malformed += fileSummary.Malformed
		summary.Records = append(summary.Records, fileSummary.Records...)
		return nil
	})
	if err != nil {
		return model.SourceResult[LocalIndexSummary]{Source: string(model.SourceLocalIndex), Path: root, FetchedAt: sourceNow(r.Clock), Status: model.StateFailed, ParseState: model.StateFailed, RawRef: fmt.Sprintf("file:%s", root), Error: err.Error()}
	}
	summary.IndexHead = readIndexHead(root)
	return model.SourceResult[LocalIndexSummary]{Source: string(model.SourceLocalIndex), Path: root, FetchedAt: sourceNow(r.Clock), Status: model.StatePresent, ParseState: model.StatePresent, RawRef: fmt.Sprintf("file:%s", root), Value: &summary}
}

func readIndexHead(path string) string {
	for _, candidate := range []string{path, filepath.Dir(path), filepath.Dir(filepath.Dir(path))} {
		headPath := filepath.Join(candidate, ".git", "HEAD")
		// #nosec G304 -- HEAD probing is constrained to the configured local
		// registry index root and parents used by Moon's registry checkout.
		data, err := os.ReadFile(headPath)
		if err != nil {
			continue
		}
		head := strings.TrimSpace(string(data))
		if ref, ok := strings.CutPrefix(head, "ref: "); ok {
			ref = filepath.Clean(filepath.FromSlash(ref))
			if filepath.IsAbs(ref) || ref == ".." || strings.HasPrefix(ref, ".."+string(os.PathSeparator)) {
				continue
			}
			// #nosec G304 -- git ref path is resolved relative to the discovered
			// registry checkout metadata, not arbitrary command input.
			refData, err := os.ReadFile(filepath.Join(candidate, ".git", ref))
			if err == nil {
				return strings.TrimSpace(string(refData))
			}
		}
		return head
	}
	return ""
}
