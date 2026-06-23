# SeekMoon WBS Coordination State

Updated: 2026-06-24

## Coordinator Role Boundary

The principal coordinator owns process runtime objects, not implementation. Direct coordinator responsibilities are:

- maintain batch order, resource map, role registry, handoff packages, review packages, promotion boundary, and evidence routes;
- dispatch specialized executors for implementation and independent review;
- receive reports, route findings, and decide whether a work package can be promoted.

The coordinator must not become the primary builder or the primary reviewer. Any coordinator file edits in this directory are process artifacts, not product implementation.

## Current Baseline

- Repository: `/home/t103o/workbench/projects/seekmoon`
- Git root: `/home/t103o/workbench/projects/seekmoon`
- Branch: `main`
- Baseline command run by coordinator: `git -C /home/t103o/workbench/projects/seekmoon status --short --branch`
- Observed status: `## main...origin/main [ahead 3]`
- Latest commit observed: `206c830 (HEAD -> main) 000`
- Implementation state observed: no `go.mod`, no Go source, no `justfile`.

This baseline is process-side recorded state. Promotion requires executor evidence plus independent review.

## Reading Contract

Every executor and reviewer must read:

1. `010-mandatory-global-context.adoc`
2. `020-wbs-dependency-map.adoc`
3. the assigned work-package chapter
4. `900-source-include-map.adoc`
5. every include target required by those files

Include rule:

- if an include points to a whole file, read that entire file;
- if an include has `lines=`, read exactly those line ranges, no more and no less;
- if an include has multiple ranges separated by semicolons, read each listed range exactly.

Reports must state which include files and ranges were read.

## Batch Plan

Batch A is the current critical path:

1. WP01 Go Module Substrate
2. WP02 Canonical Model and Contracts
3. WP03 Platform Runtime and Storage

Batch B waits for Batch A promotion:

- WP04 Public HTTP Sources and Snapshot
- WP05 Static Assets and API Profile
- WP06 Local Toolchain Source and Probe Primitives
- WP07 Repository and Skill Sources

Batch C waits for relevant upstream objects:

- WP08 Output Pipeline and Error Surface
- WP09 CLI Command Surface

Batch D waits for source readers, stores, output pipeline, and CLI surface:

- WP10 Discovery and Profile Services
- WP11 Inspection and Source Services
- WP12 Assessment, Record and Report Services

Batch E closes acceptance:

- WP13 Black-box Acceptance and Quality Gates

## Active Dispatch

- Role: builder revision
- Agent: `019ef59d-058d-78c1-bcb6-9cf1417d8b8c` (`Dewey`)
- Handoff: `coordination/handoffs/batch-a-revision-1-builder.md`
- Prompt: `coordination/prompts/003-batch-a-revision-1-builder.md`
- Write boundary: SeekMoon Go implementation root files and `internal/` Go packages listed in Batch A handoff.
- Report path: `coordination/reports/003-batch-a-revision-1-builder-report.md`

## Review State

- Batch A initial builder report: `coordination/reports/001-batch-a-builder-report.md`
- Independent review report: `coordination/reports/002-batch-a-review-report.md`
- Revision builder report: `coordination/reports/003-batch-a-revision-1-builder-report.md`
- Re-review package: `coordination/review-packages/004-batch-a-re-review.md`
- Current review verdict: pending re-review
- Blocking findings:
  - evidence wrapper JSON omits nullable `source`;
  - `internal/platform` imports `internal/model`;
  - `just fmt-check` and `goreleaser check` blocked by missing external tools.

## Promotion Boundary

Batch A cannot be promoted until all of the following are true:

- executor reports concrete changed paths and evidence commands;
- WP01, WP02, and WP03 completion evidence has been attempted;
- `go test ./...` passes;
- relevant package tests for `internal/model`, `internal/contract`, `internal/platform`, `internal/store`, and `internal/app` pass or any missing tool is explicitly reported;
- no package imports violate the 9号书 package-boundary direction;
- independent reviewer reads the same work-package basis and produces findings;
- reviewer either commits the approved changes or returns a rejection report with findings.

## Open Risks

- The repo currently lacks a Go module; dependency and tool installation may be needed.
- External tools such as `goreleaser`, `golangci-lint`, `gofumpt`, `gotestsum`, and `govulncheck` may not be installed even though the WBS requires their gates.
- WP02 is broad. The builder must avoid under-modeling state semantics or inventing public schema from reflection.
- WP03 must keep stores mechanical and avoid business decisions in path or persistence packages.
