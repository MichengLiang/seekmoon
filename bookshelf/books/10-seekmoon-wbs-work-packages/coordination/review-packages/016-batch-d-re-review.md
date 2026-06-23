# Batch D Re-Review Package

## Role

You are the independent Batch D reviewer. Reuse your prior Batch D review context, but treat the current worktree and revision report as authoritative.

If approved, write the report and commit only approved Batch D and coordination paths. If rejected, write the report and do not commit.

## Review Inputs

- Prior review package: `coordination/review-packages/014-batch-d-review.md`
- Prior review report: `coordination/reports/014-batch-d-review-report.md`
- Revision builder report: `coordination/reports/015-batch-d-revision-1-builder-report.md`

## Required Re-Review Checks

Verify the rejected findings are fixed:

1. Sync local registry index
   - `sync` reads local registry index.
   - Local index source state appears in `snapshot.Sources`.
   - Local index summary facts appear in `snapshot.Raw["local_index"]`, including status/error and available file/record/malformed/head facts.
   - Partial failure/unavailable local index state is preserved.
   - Tests cover present and failed/unavailable behavior.

2. Raw upstream payload shape
   - `raw` uses raw fetch paths.
   - Upstream field names and object/array shape are preserved.
   - Normalized model-only fields are not substituted.
   - Source metadata remains attached.
   - Tests cover upstream field preservation.

Also re-check the positive Batch D service boundaries from the initial review before approving.

## Evidence Commands

Run and report:

- `git status --short`
- `go test ./...`
- `go test ./internal/service -run 'Test(Doctor|Sync|Search|View|Compare)'`
- `go test ./internal/service -run 'Test(API|Source|Skill)'`
- `go test ./internal/service -run 'Test(Probe|Record|Report|Raw)'`
- `go test ./internal/service -run 'TestSyncRecordsLocalIndexSummaryWhenPresent|TestSyncRecordsPartialSourceFailure|TestRawReturnsSourceStatusAndUpstreamPayloadShape'`
- `go test ./internal/cli`
- `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/service ./internal/cli ./internal/output ./internal/source`
- `just fmt-check`
- `go mod tidy`
- `go mod verify`
- `goreleaser check`

If approved and committed, run `just mod-check` after commit if practical and report it.

## Approval And Commit Rule

If approved:

1. Write report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/016-batch-d-re-review-report.md`
2. Stage new files if needed.
3. Commit only these paths:

```text
internal/app/runtime.go
internal/model/output.go
internal/output/error.go
internal/output/pretty.go
internal/source/mooncli.go
internal/source/assets.go
internal/source/fetch.go
internal/source/local_index.go
internal/source/mooncakes.go
internal/source/skills.go
internal/store/paths.go
internal/store/record_store.go
internal/store/snapshot_store.go
internal/service
bookshelf/books/10-seekmoon-wbs-work-packages/coordination
```

Use `git commit --only` with a detailed message naming WP10-WP12, sync local index, raw upstream payload shape, and evidence gates.

If rejected:

1. Write the same report path.
2. Do not commit.
3. Mark findings with severity and exact file/line evidence.

Final response: verdict, report path, and commit hash if created.
