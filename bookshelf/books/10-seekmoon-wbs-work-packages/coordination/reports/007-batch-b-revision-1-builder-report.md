# Batch B Revision 1 Builder Report

## Summary

Batch B revision 1 addresses the independent review rejection in `006-batch-b-review-report.md`.

The revision completes source envelopes for the rejected non-HTTP source result producers and makes project-context parse failures observable at the aggregate `SourceResult` level while preserving partial normalized context.

No commit was created.

## Review Findings Addressed

### Finding 1: Complete source envelopes

Addressed.

Changes made:

- Added `sourceNow` helper in `internal/source/fetch.go` so local/project/repository readers can set `FetchedAt` consistently from an injected `platform.Clock`.
- Added `Path` to `FetchResult` and propagated it through the generic `SourceResult` helper.
- Added `Clock` to `LocalIndexReader`, `ProjectReader`, and `RepositoryReader`.
- Updated runtime registration in `internal/app/runtime.go` so those readers receive the runtime clock.
- Updated local index source results to include `Source`, `Path`, `FetchedAt`, `Status`, `ParseState`, `RawRef`, `Error` on failure, and `Value` on success.
- Updated project context source results to include `Source`, `Path`, `FetchedAt`, `Status`, `ParseState`, `RawRef`, `Error` on partial failure, and `Value`.
- Updated repository source results to include `Source`, `URL`, `FetchedAt`, `Status`, `ParseState`, `RawRef`, `Error` on failed/unknown outcomes, and `Value` when a signal object exists.
- Updated direct `FetchResult` error paths for manifest/package/resource relpath validation so the generated `SourceResult` envelopes include `FetchedAt` and `RawRef`.

Focused tests added or strengthened:

- `TestLocalIndexSourceEnvelopeComplete`
- `TestProjectContextJSONAndTOML`
- `TestRepositorySignalMapping`

These tests assert that `fetched_at`, `parse_state`, and `raw_ref` are populated for local index, project context, and repository source results.

### Finding 2: Project context partial failure

Addressed.

Changes made:

- `ProjectReader.Read` now preserves successful fields when one config parses successfully and another config parse fails.
- Failed module/package config parse evidence is copied into the relevant nested evidence field.
- Aggregate project context result now returns `Status: failed`, `ParseState: failed`, and an `Error` summary when any config parse action fails.
- The partial `model.ProjectContext` value is still returned in `Value` so downstream code can use successful normalized fields without losing the failed parse evidence.

Focused test added:

- `TestProjectContextPartialParseFailureObservable`

This test verifies the aggregate source result is failed/observable while the successful module context remains present.

## Files Changed

Implementation changes:

- `internal/app/runtime.go`
- `internal/source/fetch.go`
- `internal/source/assets.go`
- `internal/source/mooncakes.go`
- `internal/source/local_index.go`
- `internal/source/project.go`
- `internal/source/repository.go`

Test changes:

- `internal/source/local_test.go`
- `internal/source/repository_skill_test.go`

Report:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/007-batch-b-revision-1-builder-report.md`

Batch B files still present from the original implementation:

- `go.mod`
- `go.sum`
- `internal/app/runtime_test.go`
- `internal/model/repository.go`
- `internal/service/`
- `internal/source/`
- `internal/testutil/`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/005-batch-b-builder-report.md`

Pre-existing coordination changes remain untouched:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md`
- Batch B handoff/prompt/review package/report files already present in the worktree.

## Evidence Commands

```text
$ go test ./...
?   	github.com/yumiaura/seekmoon/cmd/seekmoon	[no test files]
ok  	github.com/yumiaura/seekmoon/internal/app	0.003s
ok  	github.com/yumiaura/seekmoon/internal/cli	0.003s
ok  	github.com/yumiaura/seekmoon/internal/contract	(cached)
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
ok  	github.com/yumiaura/seekmoon/internal/platform	(cached)
ok  	github.com/yumiaura/seekmoon/internal/service	0.078s
ok  	github.com/yumiaura/seekmoon/internal/source	0.181s
ok  	github.com/yumiaura/seekmoon/internal/store	(cached)
?   	github.com/yumiaura/seekmoon/internal/testutil	[no test files]
```

```text
$ go test ./internal/source ./internal/service -run 'Test(Mooncakes|Sync)'
ok  	github.com/yumiaura/seekmoon/internal/source	0.093s
ok  	github.com/yumiaura/seekmoon/internal/service	0.061s
```

```text
$ go test ./internal/source ./internal/model -run 'Test(Asset|ModuleIndex|PackageData|Relpath|SourceZip)'
ok  	github.com/yumiaura/seekmoon/internal/source	0.069s
ok  	github.com/yumiaura/seekmoon/internal/model	(cached) [no tests to run]
```

```text
$ go test ./internal/source ./internal/store -run 'Test(MoonCLI|LocalIndex|LocalCache|Project|ProbePath)'
ok  	github.com/yumiaura/seekmoon/internal/source	(cached)
ok  	github.com/yumiaura/seekmoon/internal/store	(cached) [no tests to run]
```

```text
$ go test ./internal/source ./internal/model -run 'Test(Repository|Skill|Runwasm)'
ok  	github.com/yumiaura/seekmoon/internal/source	(cached)
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
```

Focused revision tests:

```text
$ go test ./internal/source -run 'Test(LocalIndexSourceEnvelopeComplete|ProjectContextPartialParseFailureObservable|ProjectContextJSONAndTOML|RepositorySignalMapping)'
ok  	github.com/yumiaura/seekmoon/internal/source	0.113s
```

```text
$ just fmt-check
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal)"
```

```text
$ go mod tidy
```

Passed with no command output.

```text
$ go mod verify
all modules verified
```

```text
$ goreleaser check
• checking                                  path=.goreleaser.yaml
• 1 configuration file(s) validated
• thanks for using GoReleaser!
```

Import boundary check:

```text
$ go list -f '{{.ImportPath}} {{join .Imports "\n"}}' ./internal/source ./internal/service ./internal/app ./internal/model ./internal/platform ./internal/store
```

Relevant result: `github.com/google/go-github/v88/github` and `golang.org/x/oauth2` remain confined to `internal/source`.

## Remaining Risks Or Blocked Gates

No required revision verification command failed.

`just mod-check` was not requested by the revision handoff. `go mod tidy` and `go mod verify` both passed.

The existing Batch B dependency metadata changes remain uncommitted because the handoff says do not commit.

## Suggested Reviewer Focus

Re-review should focus on:

- whether the local index, project context, and repository `SourceResult` envelopes now satisfy the full source-reader contract;
- whether project-context aggregate failure semantics are acceptable when a partial value is still returned;
- whether direct pre-fetch validation failures now carry enough envelope information for downstream report/error surfaces;
- whether runtime clock injection for source readers is sufficient for deterministic source timing in future tests.
