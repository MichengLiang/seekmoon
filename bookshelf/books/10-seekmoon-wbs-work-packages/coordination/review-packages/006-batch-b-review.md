# Batch B Review Package

## Role

You are the independent reviewer for Batch B. Batch A is already approved and committed as `0f83682`. Review only the current Batch B changes and related coordination files.

You are not the builder and not the principal coordinator. If approved, write the review report and commit only approved Batch B and coordination paths. If rejected, write the report and do not commit.

## Review Object

Repository:

`/home/t103o/workbench/projects/seekmoon`

Builder report:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/005-batch-b-builder-report.md`

Expected changed paths:

- `go.mod`
- `go.sum`
- `internal/app/runtime.go`
- `internal/app/runtime_test.go`
- `internal/model/repository.go`
- `internal/service/*`
- `internal/source/*`
- `internal/testutil/*`
- Batch B coordination files under `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/`

## Review Basis

Read:

- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/060-wp04-public-http-sources-and-snapshot.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/070-wp05-static-assets-and-api-profile.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/080-wp06-local-toolchain-source-and-probe-primitives.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/090-wp07-repository-and-skill-sources.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/900-source-include-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-b-builder.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/005-batch-b-builder-report.md`

Include rule:

- if an include points to a whole file, read the whole file;
- if an include has `lines=`, read exactly those line ranges, no more and no less;
- if an include has several ranges separated by semicolons, read exactly every listed range.

Your report must list the files and include ranges used.

## Risk Focus

Review these risks directly:

- Source result envelopes carry stable source label, URL/path, fetched_at, status, parse_state, raw_ref, error, and normalized value.
- Stable labels are exactly `modules_api`, `statistics_api`, and `manifest_api` for WP04.
- Source readers do not render output and do not make adoption/ranking decisions.
- Modules API reader does not read asset or repository data.
- Empty module summary fields become `missing` evidence fields.
- Manifest metadata remains open and raw metadata is preserved.
- Snapshot partial failure semantics keep successful source results.
- Asset URL construction uses manifest version and derived package relpath.
- `resource.json` 404 maps to `unavailable`, not `failed`.
- Raw signature is preserved while plain signature is derived.
- Source zip is modeled as a source attempt and does not become the only published-source route.
- Moon CLI command failure is local evidence only and includes command/cwd/exit/status/log path.
- Local index parser preserves raw line and handles malformed sparse input.
- Project context reader does not hide manifest mutation.
- Repository signal does not replace published source resolution.
- `go-github` and OAuth2 imports stay in `internal/source`.
- SkillEntry and ModuleSummary remain separate models.
- Checksum state/wording does not imply provenance.
- Unrequested enrichment does not enter output objects.
- Runtime registration does not introduce dependency cycles or command-surface behavior.
- The `just mod-check` dirty-diff failure in the builder report is acceptable only if it is exactly caused by uncommitted Batch B dependency metadata and would pass after commit.

## Evidence Commands

Run and report:

- `git status --short`
- `go test ./...`
- `go test ./internal/source ./internal/service -run 'Test(Mooncakes|Sync)'`
- `go test ./internal/source ./internal/model -run 'Test(Asset|ModuleIndex|PackageData|Relpath|SourceZip)'`
- `go test ./internal/source ./internal/store -run 'Test(MoonCLI|LocalIndex|LocalCache|Project|ProbePath)'`
- `go test ./internal/source ./internal/model -run 'Test(Repository|Skill|Runwasm)'`
- `just fmt-check`
- `go mod tidy`
- `go mod verify`
- `goreleaser check`

If you approve and commit, run `just mod-check` after the commit if practical and report the result.

## Approval And Commit Rule

If approved:

1. Write the report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/006-batch-b-review-report.md`
2. Stage new files if needed.
3. Commit only these paths:

```text
go.mod
go.sum
internal/app/runtime.go
internal/app/runtime_test.go
internal/model/repository.go
internal/service
internal/source
internal/testutil
bookshelf/books/10-seekmoon-wbs-work-packages/coordination
```

Use `git commit --only` with a detailed message naming WP04-WP07 and the source-reader boundaries.

If rejected:

1. Write the same report path.
2. Do not commit.
3. Mark findings with severity and exact file/line evidence.
4. Include the required action before re-review.

Final response: verdict, report path, and commit hash if created.
