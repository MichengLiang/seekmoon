# Batch C Re-Review Package

## Role

You are the independent Batch C reviewer. Reuse prior Batch C review context, but treat the current worktree and revision report as authoritative.

If approved, write the re-review report and commit only approved Batch C and coordination paths. If rejected, write the report and do not commit.

## Review Inputs

- Prior review package: `coordination/review-packages/010-batch-c-review.md`
- Prior review report: `coordination/reports/010-batch-c-review-report.md`
- Revision builder report: `coordination/reports/011-batch-c-revision-1-builder-report.md`

## Required Re-Review Checks

Verify the rejected findings are fixed:

1. Contract projections
   - `--shape` and `--schema` render without ordinary operands for output-capable commands.
   - They do not call pending services.
   - They read contract definitions and return success.
   - Verify at least `search --shape`, `api --shape`, and one other operand-required command.

2. Real process exit mapping
   - `cmd/seekmoon` uses the same exit mapping as the CLI execution path.
   - A directly built binary exits `2` for required argument/flag parse failure.
   - Service/projection failures still exit `1`.

Also re-check the positive Batch C boundaries from the initial review before approving.

## Evidence Commands

Run and report:

- `git status --short`
- `go test ./...`
- `go test ./internal/output ./internal/contract`
- `go test ./internal/cli`
- `go run ./cmd/seekmoon search --shape`
- `go run ./cmd/seekmoon api --shape`
- a direct built-binary exit-code check for parse failure
- `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/output ./internal/cli ./internal/service ./internal/app`
- `just fmt-check`
- `go mod tidy`
- `go mod verify`
- `goreleaser check`

If approved and committed, run `just mod-check` after the commit if practical and report it.

## Approval And Commit Rule

If approved:

1. Write report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/012-batch-c-re-review-report.md`
2. Stage new files if needed.
3. Commit only these paths:

```text
cmd/seekmoon/main.go
go.mod
go.sum
internal/app/runtime.go
internal/cli
internal/contract/schemas.go
internal/contract/shapes.go
internal/model/output.go
internal/model/schema.go
internal/output
internal/service/registry.go
internal/testutil/golden.go
bookshelf/books/10-seekmoon-wbs-work-packages/coordination
```

Use `git commit --only` with a detailed message naming WP08-WP09, contract projection bypass, real binary exit mapping, and evidence gates.

If rejected:

1. Write the same report path.
2. Do not commit.
3. Mark findings with severity and exact file/line evidence.

Final response: verdict, report path, and commit hash if created.
