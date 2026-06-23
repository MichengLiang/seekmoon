# Batch D Builder Handoff

## Role

You are the persistent specialized builder for SeekMoon. Batches A, B, and C have been independently reviewed and committed:

- Batch A: `0f83682 Implement SeekMoon Batch A substrate`
- Batch B: `7bccfb1 Approve Batch B source-reader foundations`
- Batch C: `4a911ad Approve Batch C output and CLI surfaces`

You are not the principal coordinator and not the reviewer. Do not commit. Do not revert unrelated changes.

## Goal

Implement Batch D from the WBS:

- WP10 Discovery and Profile Services
- WP11 Inspection and Source Services
- WP12 Assessment, Record and Report Services

Batch D turns the already-established source readers, stores, output pipeline, and CLI command surface into service behavior. Do not implement WP13 black-box acceptance as a separate batch; keep journey/acceptance work to focused service/CLI tests needed for WP10-WP12.

## Required Reading

Read:

- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/120-wp10-discovery-and-profile-services.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/130-wp11-inspection-and-source-services.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/140-wp12-assessment-record-and-report-services.adoc`
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

- `internal/service/*`
- `internal/app/runtime.go` and tests to register concrete services
- `internal/cli/*` only where service integration requires thin handler adjustment, not new syntax beyond WP09
- `internal/model/*` only for comparison, environment status, raw envelope, service outputs, or evidence refs needed by WP10-WP12
- `internal/output/*` only where rendering existing canonical service outputs requires a projection seam
- `internal/store/*` only where service persistence mechanics are missing
- `internal/testutil/*`
- tests for these packages
- report file: `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/013-batch-d-builder-report.md`

Do not implement WP13 acceptance harnesses or unrelated broad integration suites.

## WP10 Expectations

- Implement `doctor`, `sync`, `search`, `view`, and `compare` service flows.
- `doctor` reports environment/toolchain/path/network/project context status without creating snapshots or records.
- `sync` records `moon update`, Modules API, Statistics API, local registry summary, and toolchain versions where available; partial failures remain in snapshot state.
- `search` uses latest snapshot or transient modules fetch, ranks by explainable match evidence, enriches only visible-window manifest fields, and writes session candidate mappings.
- `view` reads manifest and module index summary without expanding full API symbols.
- `compare` aligns evidence fields and must not generate a quality score or recommendation.
- Downloads/build/repository presence must not become quality proof.
- Target support remains `unknown` before metadata/probe evidence.

## WP11 Expectations

- Implement `api`, `source`, `skill search`, and `skill view` service flows.
- `api` consumes resolved module/package input, loads module index before package data, derives package relpath, and reports package path failures with known package paths.
- `source` records all attempts: `moon fetch`, source zip, local cache, core local source, repository fallback signal. Do not fold failed/unavailable/unknown attempts together.
- Repository fallback must be marked as repository signal, not registry published source.
- `skill search` and `skill view` use Skills API and skill assets, not Modules API.
- Skill profile does not enter library adoption conclusion.

## WP12 Expectations

- Implement `probe`, `record`, `report`, and `raw` service flows.
- `probe` uses isolated probe path and records command sequence/results/logs. Probe failure must not mutate upstream facts or create adoption records.
- `record` validates fixed adoption conclusion enum and persists candidate/version/snapshot/evidence refs/not-confirmed/note.
- `report` includes only records and evidence refs actually used. Do not invent repository activity or validation status.
- `raw` returns upstream field names and source status, bypassing ordinary normalization without becoming normal discovery flow.

## Verification Commands

Run and report:

- `go test ./...`
- `go test ./internal/service -run 'Test(Doctor|Sync|Search|View|Compare)'`
- `go test ./internal/service -run 'Test(API|Source|Skill)'`
- `go test ./internal/service -run 'Test(Probe|Record|Report|Raw)'`
- `go test ./internal/cli`
- `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/service ./internal/cli ./internal/output ./internal/source`
- `just fmt-check`
- `go mod tidy`
- `go mod verify`
- `goreleaser check`

If `just mod-check` fails only because the repository is dirty with approved Batch D metadata, report that precisely and include `go mod tidy` plus `go mod verify` evidence.

## Report

Write:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/013-batch-d-builder-report.md`

Report format:

- Summary
- Files changed
- Reading evidence
- Implementation notes by WP
- Tests and command evidence
- Known gaps or blocked gates
- Suggested reviewer focus

Final chat response should be short and include report path and changed paths.
