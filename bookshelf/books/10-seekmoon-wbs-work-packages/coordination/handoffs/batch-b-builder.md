# Batch B Builder Handoff

## Role

You are the persistent specialized builder for SeekMoon. Batch A has been approved and committed as `0f83682 Implement SeekMoon Batch A substrate`; treat that commit as the stable substrate.

You are not the principal coordinator and not the reviewer. Do not commit. Do not promote your own work. Do not revert unrelated changes.

## Goal

Implement Batch B from the WBS:

- WP04 Public HTTP Sources and Snapshot
- WP05 Static Assets and API Profile
- WP06 Local Toolchain Source and Probe Primitives
- WP07 Repository and Skill Sources

Batch B establishes source readers and local primitives that return state-bearing source results. It must not implement output renderers, full CLI command behavior, ranking/adoption decisions, or Batch D service flows.

## Required Reading

Read these files:

- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/060-wp04-public-http-sources-and-snapshot.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/070-wp05-static-assets-and-api-profile.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/080-wp06-local-toolchain-source-and-probe-primitives.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/090-wp07-repository-and-skill-sources.adoc`
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

- `internal/source/*`
- `internal/service/sync.go` and focused sync service tests only
- `internal/model/*` only for source-result, repository signal, local index/project context, API/skill/source fields required by WP04-WP07
- `internal/store/*` only when Batch B needs cache/source/probe/log mechanics not already present
- `internal/app/runtime.go` only to register Batch B source/service primitives
- `internal/testutil/*`
- `go.mod` and `go.sum` for WP07 dependencies
- tests for the packages above
- your report file: `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/005-batch-b-builder-report.md`

Do not modify `internal/cli` except if a compile-only type dependency requires it; do not implement command behavior. Do not implement `internal/output`.

## Implementation Expectations

Shared source layer:

- Add `internal/source/fetch.go` with shared HTTP fetch behavior: context-aware requests, timeout via provided client, redirect support through `net/http`, limited retry policy if using `backoff`, status classification, JSON decoding, parse state, raw reference/envelope, and stable source labels.
- Source results must carry source label, URL/path, fetched_at, status, parse_state, raw_ref, error, and normalized value.
- Keep source readers separate from output rendering and adoption judgment.

WP04:

- Implement Mooncakes Modules API, Statistics API, and Manifest API readers.
- Stable source labels: `modules_api`, `statistics_api`, `manifest_api`.
- Empty description/keywords/repository/license become `missing` evidence fields.
- Manifest metadata stays open; normalize only known fields and preserve raw metadata.
- Add sync service skeleton that composes source results and writes snapshot state without ranking/adoption behavior.
- Snapshot partial failures must not erase successful source results.

WP05:

- Implement asset readers for `module_index.json`, `package_data.json`, optional `resource.json`, and source zip attempts.
- Parse `childs`; optionally accept `children` as compatibility input while preserving canonical `childs`.
- Use model-level package relpath derivation; do not handwrite relpath logic in commands.
- `resource.json` 404 maps to `unavailable`.
- Source zip redirect/download/read produces a `SourceAttempt` and file summary, not source-resolution final judgment.
- Preserve raw signature and derive plain signature without replacing raw signature.

WP06:

- Implement Moon CLI runner wrapper, local registry JSONL parser, local cache/core source readers, and project context reader.
- Command results must record command, cwd, exit code, status, and log path.
- Command failure is local evidence, not upstream metadata mutation.
- Project manifest mutation must not be hidden in primitives.
- Probe/source/log paths must stay project-bounded.

WP07:

- Implement repository signal reader and skill source readers.
- Keep `go-github` and `oauth2` imports in `internal/source`.
- Repository signal must not replace published source resolution.
- SkillEntry and ModuleSummary remain separate models.
- Empty skill package maps to root marker.
- Checksum wording and state must not imply provenance.
- Unrequested enrichment must not enter output objects.

## Verification Commands

Run and report:

- `go test ./...`
- `go test ./internal/source ./internal/service -run 'Test(Mooncakes|Sync)'`
- `go test ./internal/source ./internal/model -run 'Test(Asset|ModuleIndex|PackageData|Relpath|SourceZip)'`
- `go test ./internal/source ./internal/store -run 'Test(MoonCLI|LocalIndex|LocalCache|Project|ProbePath)'`
- `go test ./internal/source ./internal/model -run 'Test(Repository|Skill|Runwasm)'`
- `just fmt-check`
- `just mod-check`

Run `goreleaser check` unless your changes make it clearly irrelevant; if skipped, explain why.

## Report

Write a long report at:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/005-batch-b-builder-report.md`

Report format:

- Summary
- Files changed
- Reading evidence
- Implementation notes by WP
- Tests and command evidence
- Known gaps or blocked gates
- Suggested reviewer focus

Final chat response should be short and include the report path and changed paths.
