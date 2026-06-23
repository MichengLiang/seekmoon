# Batch A Revision 1 Builder Handoff

## Role

You are the same specialized builder for Batch A. Keep your local Batch A context, but treat the independent reviewer report as the controlling rejection basis for this revision.

You are not the principal coordinator and not the reviewer. Do not promote your own work. Do not commit unless explicitly told to do so.

## Review Report To Address

Read and address:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/002-batch-a-review-report.md`

Original builder handoff remains binding:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-a-builder.md`

## Required Fixes

Fix only the rejected Batch A issues and required gates:

1. Evidence wrapper source preservation
   - `source` must be preserved in JSON as a required nullable field.
   - Unknown or source-less evidence should encode `"source": null`, not omit the field.
   - Keep `error` optional.
   - Update constructors, schema definitions, and tests.

2. Remove `internal/platform -> internal/model`
   - `internal/platform` must not import `internal/model`.
   - `platform.Runner` should return a platform-local execution result and error.
   - Move mapping to `model.CommandResult` / `model.State` into a package allowed to import both model and platform, preferably `internal/store` only if the mapping is actually store-related, or `internal/app`/future source seam if appropriate for Batch A. Keep the abstraction clean and minimal.
   - Update tests and import-boundary evidence.

3. External quality gates
   - Install or otherwise provide the pinned tools if missing:
     - `go install mvdan.cc/gofumpt@v0.10.0`
     - `go install github.com/goreleaser/goreleaser/v2@v2.16.0`
   - Run:
     - `just fmt-check`
     - `goreleaser check`
   - If tool installation fails for an environmental reason, report the exact failure.

## Verification Commands

Run and report:

- `go test ./...`
- `go test ./internal/model ./internal/contract`
- `go test ./internal/platform ./internal/store ./internal/app`
- `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/model ./internal/contract ./internal/platform ./internal/store ./internal/app ./internal/cli`
- `just fmt-check`
- `just mod-check`
- `goreleaser check`

## Report

Write the revision report at:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/003-batch-a-revision-1-builder-report.md`

Report format:

- Summary
- Review findings addressed
- Files changed
- Evidence commands
- Remaining risks or blocked gates
- Suggested reviewer focus

Your final chat response should be short and include the report path and changed paths.
