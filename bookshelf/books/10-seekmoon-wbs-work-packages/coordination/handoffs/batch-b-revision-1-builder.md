# Batch B Revision 1 Builder Handoff

## Role

You are the persistent Batch B builder. Treat the independent review report as the controlling rejection basis.

You are not the principal coordinator and not the reviewer. Do not commit. Do not revert unrelated changes.

## Review Report To Address

Read:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/006-batch-b-review-report.md`

Original Batch B handoff remains binding:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-b-builder.md`

## Required Fixes

Fix only the Batch B rejection findings:

1. Complete source envelopes
   - Every `model.SourceResult` producer must set:
     - `Source`
     - one of `URL` or `Path`
     - `FetchedAt`
     - `Status`
     - `ParseState`
     - `RawRef`
     - `Error` when failed or partially failed
     - `Value` when normalized value exists
   - Apply this to HTTP, local index, local cache/core source where applicable, project context, repository, and skills readers.
   - Add tests that fail if `fetched_at`, `parse_state`, or `raw_ref` are omitted from local index, project context, and repository source results.

2. Project context partial failure
   - If any project config parse action fails, return a partial source result rather than a wholly present result.
   - Preserve successful normalized fields where useful.
   - Set source-level `Status`, `ParseState`, and/or `Error` so the failed action is observable without requiring downstream code to inspect every nested evidence field.

## Verification Commands

Run and report:

- `go test ./...`
- `go test ./internal/source ./internal/service -run 'Test(Mooncakes|Sync)'`
- `go test ./internal/source ./internal/model -run 'Test(Asset|ModuleIndex|PackageData|Relpath|SourceZip)'`
- `go test ./internal/source ./internal/store -run 'Test(MoonCLI|LocalIndex|LocalCache|Project|ProbePath)'`
- `go test ./internal/source ./internal/model -run 'Test(Repository|Skill|Runwasm)'`
- `just fmt-check`
- `go mod tidy`
- `go mod verify`
- `goreleaser check`

Also run any focused tests you add for envelope completeness.

## Report

Write the revision report at:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/007-batch-b-revision-1-builder-report.md`

Report format:

- Summary
- Review findings addressed
- Files changed
- Evidence commands
- Remaining risks or blocked gates
- Suggested reviewer focus

Final response should be short and include report path and changed paths.
