# Batch C Builder Handoff

## Role

You are the persistent specialized builder for SeekMoon. Batch A and Batch B have both been approved and committed:

- Batch A: `0f83682 Implement SeekMoon Batch A substrate`
- Batch B: `7bccfb1 Approve Batch B source-reader foundations`

You are not the principal coordinator and not the reviewer. Do not commit. Do not revert unrelated changes.

## Goal

Implement Batch C from the WBS:

1. WP08 Output Pipeline and Error Surface
2. WP09 CLI Command Surface

The internal order matters: establish output projections first, then wire the Cobra command surface to runtime services and renderer seams. Batch C must not implement Batch D service behavior. Concrete source reading and business service composition remain outside this batch except for already existing Batch B sync skeletons and fake service seams needed for CLI tests.

## Required Reading

Read these files:

- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/100-wp08-output-pipeline-and-error-surface.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/110-wp09-cli-command-surface.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/900-source-include-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md`

Include rule:

- if an include points to a whole file, read the whole file;
- if an include has `lines=`, read exactly those line ranges, no more and no less;
- if an include has several ranges separated by semicolons, read exactly every listed range.

Your report must list the WBS files and include ranges you read.

## Write Boundary

You may create or modify:

- `internal/output/*`
- `internal/cli/*`
- `internal/app/runtime.go` and tests only to register renderer/service seams needed by WP08-WP09
- `internal/service/*` only for interface definitions or fakeable registry seams required by CLI command handlers
- `internal/model/*` only for output request/error/input structs required by WP08-WP09
- `internal/contract/*` only if existing shape/schema definitions must be extended for rendering
- `internal/testutil/*`
- `go.mod` and `go.sum` for WP08 dependencies
- tests for the packages above
- report file: `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/009-batch-c-builder-report.md`

Do not implement concrete Batch D command services such as full search/view/api/source/probe/record/report behavior. Use fakeable interfaces and placeholder service seams when necessary.

## WP08 Expectations

- Add `internal/output/render.go` with renderer interface, output request, output mode dispatch, and writer handling.
- Add JSON projection that emits schema id and canonical object data; do not serialize pretty text as JSON.
- Add jq projection using `github.com/itchyny/gojq@v0.12.19`; jq consumes JSON projection, never pretty text.
- Add shape and schema projections from `internal/contract`, not sample result values.
- Add error surface projection with command, object, source, state, meaning, recovery, and log path when present.
- Pretty text should be low-noise terminal projection for current canonical objects. It can cover the current available/fake canonical object set, but must not contain recovery tutorial text in normal output.
- `internal/output` must not import `internal/source`.
- Projection-specific failures must map to error surface, not panic or stack traces.

## WP09 Expectations

- Implement Cobra command files listed by WP09:
  - `root.go`
  - `doctor.go`
  - `sync.go`
  - `search.go`
  - `view.go`
  - `api.go`
  - `source.go`
  - `skill.go`
  - `compare.go`
  - `probe.go`
  - `record.go`
  - `report.go`
  - `raw.go`
  - `flags.go`
- Define root command, global flags, output mode options, command registration, syntax validation, input parsing, thin handlers, and process exit mapping.
- Command handlers follow: parse args -> build input -> context -> service -> render -> exit mapping.
- Cobra types remain in `internal/cli`; `internal/cli` must not import `internal/source`.
- Command parameters use object names and candidate numbers, not hidden implementation URLs.
- Required argument and flag parse failures map to exit code `2`; service/projection failures map to exit code `1`.
- Output mode handling must be shared and consistent; `--jq` without an expression is a parse failure.

## Verification Commands

Run and report:

- `go test ./...`
- `go test ./internal/output ./internal/contract`
- `go test ./internal/cli`
- `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/output ./internal/cli ./internal/service ./internal/app`
- `just fmt-check`
- `go mod tidy`
- `go mod verify`
- `goreleaser check`

If `just mod-check` fails only because the repository is dirty with approved Batch C dependency metadata, report that precisely and include `go mod tidy` plus `go mod verify` evidence.

## Report

Write a long report at:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/009-batch-c-builder-report.md`

Report format:

- Summary
- Files changed
- Reading evidence
- Implementation notes by WP
- Tests and command evidence
- Known gaps or blocked gates
- Suggested reviewer focus

Final chat response should be short and include the report path and changed paths.
