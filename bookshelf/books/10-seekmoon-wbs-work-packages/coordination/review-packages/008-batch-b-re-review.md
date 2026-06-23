# Batch B Re-Review Package

## Role

You are the independent Batch B reviewer. Reuse your prior Batch B review context, but treat the current worktree and revision report as authoritative.

If approved, write the re-review report and commit only approved Batch B and coordination paths. If rejected, write the report and do not commit.

## Review Inputs

- Prior review package: `coordination/review-packages/006-batch-b-review.md`
- Prior review report: `coordination/reports/006-batch-b-review-report.md`
- Revision builder report: `coordination/reports/007-batch-b-revision-1-builder-report.md`

## Required Re-Review Checks

Verify the rejection findings are fixed:

1. Source envelope completeness
   - Local index source results include source, path, fetched_at, status, parse_state, raw_ref, error on failure, and value on success.
   - Project context source results include source, path, fetched_at, status, parse_state, raw_ref, error on partial failure, and value.
   - Repository source results include source, URL, fetched_at, status, parse_state, raw_ref, error on failed/unknown outcomes, and value when normalized signal exists.
   - Pre-fetch validation failures in asset/mooncakes readers still produce enough envelope information for downstream report/error surfaces.

2. Project context aggregate failure
   - A malformed config parse action is visible at source-result aggregate level.
   - Successful partial context remains available where useful.
   - Tests cover this behavior.

Also re-check the positive Batch B boundaries from the initial review before approving.

## Evidence Commands

Run and report:

- `git status --short`
- `go test ./...`
- `go test ./internal/source ./internal/service -run 'Test(Mooncakes|Sync)'`
- `go test ./internal/source ./internal/model -run 'Test(Asset|ModuleIndex|PackageData|Relpath|SourceZip)'`
- `go test ./internal/source ./internal/store -run 'Test(MoonCLI|LocalIndex|LocalCache|Project|ProbePath)'`
- `go test ./internal/source ./internal/model -run 'Test(Repository|Skill|Runwasm)'`
- `go test ./internal/source -run 'Test(LocalIndexSourceEnvelopeComplete|ProjectContextPartialParseFailureObservable|ProjectContextJSONAndTOML|RepositorySignalMapping)'`
- `just fmt-check`
- `go mod tidy`
- `go mod verify`
- `goreleaser check`

If approved and committed, run `just mod-check` after the commit if practical and report it.

## Approval And Commit Rule

If approved:

1. Write report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/008-batch-b-re-review-report.md`
2. Stage new files if needed.
3. Commit only these paths:

```text
go.mod
go.sum
internal/app/runtime.go
internal/app/runtime_test.go
internal/model/repository.go
internal/service
internal/source
internal/testutil
bookshelf/books/10-seekmoon-wbs-work-packages/coordination
```

Use `git commit --only` with a detailed message naming WP04-WP07, source envelopes, project context partial failure semantics, and evidence gates.

If rejected:

1. Write the same report path.
2. Do not commit.
3. Mark findings with severity and exact file/line evidence.

Final response: verdict, report path, and commit hash if created.
