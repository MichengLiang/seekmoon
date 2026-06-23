# Batch E Builder Handoff

## Role

You are the persistent specialized builder for SeekMoon. Batches A-D have been independently reviewed and committed:

- Batch A: `0f83682 Implement SeekMoon Batch A substrate`
- Batch B: `7bccfb1 Approve Batch B source-reader foundations`
- Batch C: `4a911ad Approve Batch C output and CLI surfaces`
- Batch D: `059456c Approve Batch D WP10-WP12 service flows`

You are not the principal coordinator and not the reviewer. Do not commit. Do not revert unrelated changes.

## Goal

Implement Batch E / WP13: black-box acceptance and quality gates.

Batch E must provide executable evidence for A1-A16, journey behavior, schema/golden/failure behavior, opt-in integration boundaries, and final quality gate composition. It should not redefine public behavior or add new product features except small testability seams needed for acceptance.

## Required Reading

Read:

- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/150-wp13-black-box-acceptance-and-quality-gates.adoc`
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

- `tests/blackbox/*`
- `tests/journey/*`
- `tests/integration/*`
- `internal/testutil/*`
- `internal/output/*` or `internal/cli/*` only for small testability seams required by black-box acceptance
- `justfile` for final WP13 quality and integration commands
- `.github/workflows/*` if needed to expose the same gates as local commands
- documentation inside coordination reports only
- report file: `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/017-batch-e-builder-report.md`

Do not make broad service or source-reader rewrites unless an acceptance test exposes a defect that must be fixed for A1-A16.

## Required Acceptance Mapping

Create executable tests whose names map visibly to A1-A16:

- A1 search generates candidates without hand-written Mooncakes URLs.
- A2 search result can be referenced by session-local number.
- A3 library module and skill entry use different command surfaces.
- A4 module profile contains manifest evidence and package index state.
- A5 package API profile comes from module index and package data.
- A6 published source can be fetched or located through source resolution.
- A7 target support remains `unknown` before evidence.
- A8 probe produces local derived evidence and does not mutate upstream facts.
- A9 adoption decision persists as record with evidence refs.
- A10 report lists only actually used sources.
- A11 pretty text is low-noise and is not a parsing interface.
- A12 JSON output contains schema id and evidence states.
- A13 built-in jq evaluates command JSON output.
- A14 shape explains JSON fields without real query execution.
- A15 schema provides JSON Schema for strict consumers.
- A16 command failure uses error surface with source, state, meaning, and recovery.

Default acceptance tests must be offline, repeatable, and free of external side effects. Use fake source servers, temp directories, fake command runners, and local fixtures.

## Journey And Integration Boundaries

- Add library, skill, pipeline, and failure-recovery journey tests where they add coverage beyond A1-A16.
- Add integration tests only under opt-in environment variables. They must be clearly named and skipped by default.
- Tests that use real network, Moon CLI, GitHub API, source zip downloads, or probe mutation must require explicit integration env vars.

## Quality Gates

Ensure `justfile` exposes and runs the final gates from WP13:

- `just fmt-check`
- `just lint`
- `just test`
- `just test-race`
- `just cover`
- `just vuln`
- `just mod-check`
- `just release-check`
- `just ci`

If any required external tool is missing, install the pinned tool from the toolchain appendix when possible and report exact evidence. If installation is impossible, report the blocker precisely.

## Verification Commands

Run and report:

- `go test ./...`
- `go test ./tests/blackbox ./tests/journey`
- `just fmt-check`
- `just lint`
- `just test`
- `just test-race`
- `just cover`
- `just vuln`
- `just mod-check`
- `just release-check`
- `just ci`

If integration tests exist, run their default no-env form and report the skip behavior.

## Report

Write:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/017-batch-e-builder-report.md`

Report format:

- Summary
- Files changed
- Reading evidence
- A1-A16 acceptance mapping
- Journey and integration boundary notes
- Quality gate evidence
- Known gaps or blocked gates
- Suggested reviewer focus

Final chat response should be short and include report path and changed paths.
