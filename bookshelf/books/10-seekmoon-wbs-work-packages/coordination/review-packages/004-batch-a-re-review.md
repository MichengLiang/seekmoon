# Batch A Re-Review Package

## Role

You are the independent Batch A reviewer. Reuse your review context from the first review, but treat the current worktree and the revision report as authoritative current state.

You are not the principal coordinator and not the builder. If approved, write the re-review report and commit only the approved Batch A and coordination paths. If rejected, write the re-review report and do not commit.

## Review Object

Review the current uncommitted Batch A work in:

`/home/t103o/workbench/projects/seekmoon`

Required reports:

- Initial builder report: `coordination/reports/001-batch-a-builder-report.md`
- Initial review report: `coordination/reports/002-batch-a-review-report.md`
- Revision builder report: `coordination/reports/003-batch-a-revision-1-builder-report.md`

## Review Basis

Use the same review basis as:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/002-batch-a-review.md`

You do not need to duplicate the full basis list in this package, but your report must state whether you re-used the prior basis and which additional revision files you read.

## Required Re-Review Checks

Verify the original rejection findings are fixed:

1. Evidence wrapper source preservation
   - JSON always includes `source`.
   - Source-less evidence encodes `"source": null`.
   - `error` remains optional.
   - Contract schema requires `source` and accepts `string|null`.
   - Tests cover these cases.

2. Package-boundary direction
   - `internal/platform` does not import `internal/model`.
   - `platform.Runner` returns a platform-local result.
   - Mapping to `model.CommandResult` is in an acceptable higher-level package.

3. Quality gates
   - `just fmt-check` passes.
   - `goreleaser check` passes.
   - Judge whether the revised `just fmt-check` scope of `cmd internal` is acceptable for WP01, given unrelated spike cache files outside Batch A.

Also re-check the positive Batch A boundaries from your initial review before approving.

## Evidence Commands

Run and report:

- `git status --short`
- `go test ./...`
- `go test ./internal/model ./internal/contract`
- `go test ./internal/platform ./internal/store ./internal/app`
- `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/model ./internal/contract ./internal/platform ./internal/store ./internal/app ./internal/cli`
- `just fmt-check`
- `just mod-check`
- `goreleaser check`

## Approval And Commit Rule

If approved:

1. Write report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/004-batch-a-re-review-report.md`
2. Stage any new files if needed.
3. Commit only these paths:

```text
.gitignore
.golangci.yml
.goreleaser.yaml
go.mod
go.sum
justfile
cmd
internal
bookshelf/books/10-seekmoon-wbs-work-packages/coordination
```

Use `git commit --only` and a detailed message naming WP01, WP02, WP03, the boundary fixes, and the evidence gates.

If rejected:

1. Write the same report path.
2. Do not commit.
3. Mark findings with severity and exact file/line evidence.

Final response: verdict, report path, and commit hash if created.
