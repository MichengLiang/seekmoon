# Batch D Review Package

## Role

You are the independent reviewer for Batch D. Batches A-C are already approved and committed. Review only current Batch D changes and related coordination files.

You are not the builder and not the principal coordinator. If approved, write the report and commit only approved Batch D and coordination paths. If rejected, write the report and do not commit.

## Review Object

Repository:

`/home/t103o/workbench/projects/seekmoon`

Builder report:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/013-batch-d-builder-report.md`

Expected changed paths:

- `internal/app/runtime.go`
- `internal/model/output.go`
- `internal/output/error.go`
- `internal/output/pretty.go`
- `internal/source/mooncli.go`
- `internal/store/record_store.go`
- `internal/store/snapshot_store.go`
- `internal/service/*`
- Batch D coordination files

## Review Basis

Read:

- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/120-wp10-discovery-and-profile-services.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/130-wp11-inspection-and-source-services.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/140-wp12-assessment-record-and-report-services.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/900-source-include-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-d-builder.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/013-batch-d-builder-report.md`

Include rule:

- if an include points to a whole file, read the whole file;
- if an include has `lines=`, read exactly those line ranges, no more and no less;
- if an include has several ranges separated by semicolons, read exactly every listed range.

Your report must list the files and include ranges used.

## Risk Focus

Review these risks directly:

- Services compose source results and stores; they do not render directly.
- `doctor` does not create snapshots or records.
- `sync` preserves partial source failures and does not do ranking/adoption.
- `search` uses snapshot or transient modules fetch, writes session candidate mapping, enriches only visible manifest fields, and keeps target unknown before evidence.
- `view` reads manifest and module index summary without full API expansion.
- `compare` aligns evidence and has no quality score or recommendation.
- `api` does not re-run search; package path failures include known package paths from module index.
- `source` records all attempts and does not collapse failed/unavailable/unknown.
- Repository fallback is labeled as repository signal, not registry published source.
- `skill` uses Skills API and stays separate from library adoption.
- `probe` uses isolated path, records command results/logs, and does not mutate upstream facts or records.
- `record` validates only the fixed adoption conclusion enum.
- `report` lists only sources/evidence actually present in records.
- `raw` preserves upstream payload names and source status; it does not normalize or become normal discovery flow.
- Store listing helpers using `os.ReadDir` are acceptable or should be moved behind platform FS before promotion.
- No WP13 black-box acceptance harnesses were smuggled into Batch D.

## Evidence Commands

Run and report:

- `git status --short`
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

If approved and committed, run `just mod-check` after the commit if practical and report it.

## Approval And Commit Rule

If approved:

1. Write report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/014-batch-d-review-report.md`
2. Stage new files if needed.
3. Commit only these paths:

```text
internal/app/runtime.go
internal/model/output.go
internal/output/error.go
internal/output/pretty.go
internal/source/mooncli.go
internal/store/record_store.go
internal/store/snapshot_store.go
internal/service
bookshelf/books/10-seekmoon-wbs-work-packages/coordination
```

Use `git commit --only` with a detailed message naming WP10-WP12, service boundaries, evidence semantics, and gates.

If rejected:

1. Write the same report path.
2. Do not commit.
3. Mark findings with severity and exact file/line evidence.
4. Include required action before re-review.

Final response: verdict, report path, and commit hash if created.
