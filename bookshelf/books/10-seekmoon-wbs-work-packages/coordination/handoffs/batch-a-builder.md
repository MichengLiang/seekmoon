# Batch A Builder Handoff

## Role

You are a specialized executor in builder role. You are not the principal coordinator, not the reviewer, and not the user-facing interface.

You are not alone in this repository. Do not revert other people's work. If you encounter unrelated changes, leave them alone. If changes affect your task, adapt to the current state and report the interaction.

## Goal

Implement Batch A from the SeekMoon WBS:

1. WP01: Go Module Substrate
2. WP02: Canonical Model and Contracts
3. WP03: Platform Runtime and Storage

The requested end state is not a toy skeleton. Batch A should establish a compilable Go module, canonical object language, explicit contract definitions, platform/store boundaries, runtime composition, focused tests, and quality command entrypoints aligned with 9号书 and 10号书.

## Required Reading

Read these files first:

- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/030-wp01-go-module-substrate.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/040-wp02-canonical-model-and-contracts.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/050-wp03-platform-runtime-and-storage.adoc`
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

- `go.mod`
- `go.sum`
- `justfile`
- `.golangci.yml`
- `.goreleaser.yaml`
- `.gitignore`
- `.github/workflows/*`
- `cmd/seekmoon/main.go`
- `internal/app/*`
- `internal/cli/*`
- `internal/model/*`
- `internal/contract/*`
- `internal/platform/*`
- `internal/store/*`
- `internal/testutil/*`
- tests under those package directories
- your report file: `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/001-batch-a-builder-report.md`

Do not implement Batch B-D behavior except minimal compile seams required by Batch A. Do not add source readers, service flows, output rendering behavior, or complete CLI commands beyond minimal placeholders needed for the process entrypoint to compile.

## Architecture Constraints

- `cmd/seekmoon` only owns process startup: construct runtime and execute CLI root.
- `internal/model` imports only Go standard library packages.
- `internal/contract` may depend on model constants but must not auto-generate public schema from struct reflection.
- `internal/platform` owns host primitives: clock, filesystem, HTTP client, runner, env.
- `internal/store` owns path construction and read/write mechanics, not business decisions.
- `internal/app` owns runtime construction and registry scaffolding.
- `internal/cli` may use Cobra for the root command, but business behavior stays out of CLI handlers.
- `context.Context` must be present in APIs that can block, cancel, or touch host environment.
- State vocabulary must preserve exactly: `present`, `missing`, `unknown`, `failed`, `unavailable`, `derived`.
- Do not collapse `missing`, `unknown`, `failed`, or `unavailable` into generic errors.

## Implementation Expectations

WP01:

- create the Go module and module layout from 9号书;
- add `justfile` entries named exactly as the WBS requires;
- add `.golangci.yml`, `.goreleaser.yaml`, and `.artifacts/` handling consistent with generated-output boundaries;
- make `go test ./...` pass from repository root.

WP02:

- define model files listed in the WBS with explicit JSON shapes and state/evidence wrappers;
- add tests for state parsing, state semantics, evidence wrapper JSON round-trip, module coordinate parsing, package relpath derivation, adoption conclusion validation, and contract shape/schema definitions;
- keep schema IDs explicit and stable.

WP03:

- define platform capabilities and concrete OS/system implementations;
- define store path construction, atomic write behavior where required, project storage and cache paths, session candidate mapping, record/report path sanitization, and log path construction;
- define runtime composition in `internal/app` with construction order aligned to the spec;
- add tests using `t.TempDir()` and fake platform implementations where useful.

## Verification Commands

Run what is available and report exact results:

- `go test ./...`
- `go test ./internal/model ./internal/contract`
- `go test ./internal/platform ./internal/store ./internal/app`
- `just fmt-check`
- `just mod-check`
- `goreleaser check`

If a tool is not installed, report that specifically. Do not treat missing tools as passing.

## Git And Report

Do not commit unless all Batch A evidence passes or the coordinator explicitly asks for an interim commit. The reviewer is expected to commit after approval.

Write a long report at:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/001-batch-a-builder-report.md`

Report format:

- Summary
- Files changed
- Reading evidence
- Implementation notes by WP
- Tests and command evidence
- Known gaps or blocked gates
- Suggested reviewer focus

Your final chat response should be short and point to the report path.
