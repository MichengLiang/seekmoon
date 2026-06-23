# Batch D Revision 1 Builder Handoff

## Role

You are the persistent Batch D builder. Treat the independent review report as the controlling rejection basis.

You are not the principal coordinator and not the reviewer. Do not commit. Do not revert unrelated changes.

## Review Report To Address

Read:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/014-batch-d-review-report.md`

Original Batch D handoff remains binding:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-d-builder.md`

## Required Fixes

Fix only the Batch D rejection findings:

1. Sync local registry index source
   - `sync` must read the local registry index.
   - Append the state-bearing local index `SourceResult` to `snapshot.Sources`.
   - Record local index summary facts required by WP10: status/error and available summary facts such as file count and record count.
   - Preserve partial source failure state rather than treating local index omission as success.
   - Add tests for present and failed/unavailable local index behavior.

2. Raw upstream payload shape
   - `raw` must preserve upstream JSON object/array field names and original shape.
   - Do not return normalized model payloads for raw API/asset sources.
   - Still include source label, URL/path, fetch status, timestamp, raw ref, and parse/failure state.
   - Add tests that assert upstream field names remain visible and canonical-only fields are not substituted.

## Verification Commands

Run and report:

- `go test ./...`
- `go test ./internal/service -run 'Test(Doctor|Sync|Search|View|Compare)'`
- `go test ./internal/service -run 'Test(API|Source|Skill)'`
- `go test ./internal/service -run 'Test(Probe|Record|Report|Raw)'`
- focused tests added for sync local index and raw upstream payload shape
- `go test ./internal/cli`
- `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/service ./internal/cli ./internal/output ./internal/source`
- `just fmt-check`
- `go mod tidy`
- `go mod verify`
- `goreleaser check`

## Report

Write:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/015-batch-d-revision-1-builder-report.md`

Report format:

- Summary
- Review findings addressed
- Files changed
- Evidence commands
- Remaining risks or blocked gates
- Suggested reviewer focus

Final response should be short and include report path and changed paths.
