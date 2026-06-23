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
- Initial observed status: `## main...origin/main [ahead 3]`
- Initial observed commit: `206c830 (HEAD -> main) 000`
- Batch A promotion commit: `0f83682 Implement SeekMoon Batch A substrate`
- Post-Batch-A observed status: `## main...origin/main [ahead 4]`

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

- Role: pending Batch B builder
- Agent: `019ef59d-058d-78c1-bcb6-9cf1417d8b8c` (`Dewey`), intended persistent builder reuse
- Handoff: `coordination/handoffs/batch-b-builder.md`
- Prompt: `coordination/prompts/005-batch-b-builder.md`
- Write boundary: Batch B source readers, primitives, tests, and reports
- Report path: `coordination/reports/005-batch-b-builder-report.md`

## Batch B Review

- Builder report: `coordination/reports/005-batch-b-builder-report.md`
- Review package: `coordination/review-packages/006-batch-b-review.md`
- Reviewer: `019ef5b4-df23-7151-85dd-41239d63c743` (`Curie`)
- Review report: `coordination/reports/006-batch-b-review-report.md`
- Initial verdict: rejected
- Revision handoff: `coordination/handoffs/batch-b-revision-1-builder.md`
- Revision prompt: `coordination/prompts/007-batch-b-revision-1-builder.md`
- Revision report: `coordination/reports/007-batch-b-revision-1-builder-report.md`
- Re-review package: `coordination/review-packages/008-batch-b-re-review.md`
- Re-review prompt: `coordination/prompts/008-batch-b-re-reviewer.md`
- Blocking findings:
  - non-HTTP/project/repository source readers do not carry complete source envelopes;
  - project context parse failures are hidden by aggregate present status.
- Revision report: `coordination/reports/007-batch-b-revision-1-builder-report.md`
- Re-review report: `coordination/reports/008-batch-b-re-review-report.md`
- Final verdict: approved
- Batch B commit: `7bccfb1 Approve Batch B source-reader foundations`
- Promotion status: stable accepted belief for WP04-WP07 source-reader foundations

## Batch C Dispatch

- Status: ready to dispatch
- Work packages: WP08 Output Pipeline and Error Surface, then WP09 CLI Command Surface
- Intended builder: `019ef59d-058d-78c1-bcb6-9cf1417d8b8c` (`Dewey`)
- Intended reviewer: `019ef5b4-df23-7151-85dd-41239d63c743` (`Curie`)
- Handoff: `coordination/handoffs/batch-c-builder.md`
- Prompt: `coordination/prompts/009-batch-c-builder.md`
- Report path: `coordination/reports/009-batch-c-builder-report.md`
- Review package: `coordination/review-packages/010-batch-c-review.md`
- Review prompt: `coordination/prompts/010-batch-c-reviewer.md`
- Initial review report: `coordination/reports/010-batch-c-review-report.md`
- Initial verdict: rejected
- Revision handoff: `coordination/handoffs/batch-c-revision-1-builder.md`
- Revision prompt: `coordination/prompts/011-batch-c-revision-1-builder.md`
- Revision report: `coordination/reports/011-batch-c-revision-1-builder-report.md`
- Re-review package: `coordination/review-packages/012-batch-c-re-review.md`
- Re-review prompt: `coordination/prompts/012-batch-c-re-reviewer.md`
- Re-review report: `coordination/reports/012-batch-c-re-review-report.md`
- Final verdict: approved
- Batch C commit: `4a911ad Approve Batch C output and CLI surfaces`
- Promotion status: stable accepted belief for WP08-WP09 output and CLI surfaces
- Blocking findings:
  - `--shape`/`--schema` contract modes call normal services or require normal operands;
  - real `cmd/seekmoon` maps parse/flag failures to exit code `1` instead of `2`.

## Batch D Dispatch

- Status: ready to dispatch
- Work packages: WP10 Discovery and Profile Services, WP11 Inspection and Source Services, WP12 Assessment Record and Report Services
- Intended builder: `019ef59d-058d-78c1-bcb6-9cf1417d8b8c` (`Dewey`)
- Intended reviewer: `019ef5b4-df23-7151-85dd-41239d63c743` (`Curie`)
- Handoff: `coordination/handoffs/batch-d-builder.md`
- Prompt: `coordination/prompts/013-batch-d-builder.md`
- Report path: `coordination/reports/013-batch-d-builder-report.md`
- Review package: `coordination/review-packages/014-batch-d-review.md`
- Review prompt: `coordination/prompts/014-batch-d-reviewer.md`
- Initial review report: `coordination/reports/014-batch-d-review-report.md`
- Initial verdict: rejected
- Revision handoff: `coordination/handoffs/batch-d-revision-1-builder.md`
- Revision prompt: `coordination/prompts/015-batch-d-revision-1-builder.md`
- Revision report: `coordination/reports/015-batch-d-revision-1-builder-report.md`
- Re-review package: `coordination/review-packages/016-batch-d-re-review.md`
- Re-review prompt: `coordination/prompts/016-batch-d-re-reviewer.md`
- Re-review report: `coordination/reports/016-batch-d-re-review-report.md`
- Final verdict: approved
- Batch D commit: `059456c Approve Batch D WP10-WP12 service flows`
- Promotion status: stable accepted belief for WP10-WP12 service flows
- Blocking findings:
  - `sync` omits local registry index source and summary;
  - `raw` returns normalized model payloads instead of upstream raw payload shape.

## Batch E Dispatch

- Status: ready to dispatch
- Work package: WP13 Black-box Acceptance and Quality Gates
- Intended builder: `019ef59d-058d-78c1-bcb6-9cf1417d8b8c` (`Dewey`)
- Intended reviewer: `019ef5b4-df23-7151-85dd-41239d63c743` (`Curie`)
- Handoff: `coordination/handoffs/batch-e-builder.md`
- Prompt: `coordination/prompts/017-batch-e-builder.md`
- Report path: `coordination/reports/017-batch-e-builder-report.md`
- Review package: `coordination/review-packages/018-batch-e-review.md`
- Review prompt: `coordination/prompts/018-batch-e-reviewer.md`
- Builder-reported blockers:
  - `just lint` fails with broad lint findings;
  - `just vuln` fails due Go standard library vulnerabilities in local `go1.26.3`, fixed in `go1.26.4`;
  - `just ci` fails because it includes lint and vuln gates.
- Initial review report: `coordination/reports/018-batch-e-review-report.md`
- Initial verdict: rejected
- Revision handoff: `coordination/handoffs/batch-e-revision-1-builder.md`
- Revision prompt: `coordination/prompts/019-batch-e-revision-1-builder.md`
- Revision report: `coordination/reports/019-batch-e-revision-1-builder-report.md`
- Re-review package: `coordination/review-packages/020-batch-e-re-review.md`
- Re-review prompt: `coordination/prompts/020-batch-e-re-reviewer.md`
- Re-review status: ready for independent reviewer
- Revision evidence:
  - local Go toolchain updated to `go1.26.4`;
  - `PATH="$(go env GOPATH)/bin:$PATH" just lint` passes;
  - `PATH="$(go env GOPATH)/bin:$PATH" just vuln` passes;
  - `PATH="$(go env GOPATH)/bin:$PATH" just ci` passes.

## Review State

- Batch A initial builder report: `coordination/reports/001-batch-a-builder-report.md`
- Independent review report: `coordination/reports/002-batch-a-review-report.md`
- Revision builder report: `coordination/reports/003-batch-a-revision-1-builder-report.md`
- Re-review package: `coordination/review-packages/004-batch-a-re-review.md`
- Batch A final verdict: approved
- Batch A commit: `0f83682 Implement SeekMoon Batch A substrate`
- Promotion status: stable accepted belief for WP01-WP03 substrate

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
