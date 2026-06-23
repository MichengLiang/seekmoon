# Batch C Review Package

## Role

You are the independent reviewer for Batch C. Batch A and Batch B are already approved and committed. Review only current Batch C changes and related coordination files.

You are not the builder and not the principal coordinator. If approved, write the report and commit only approved Batch C and coordination paths. If rejected, write the report and do not commit.

## Review Object

Repository:

`/home/t103o/workbench/projects/seekmoon`

Builder report:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/009-batch-c-builder-report.md`

Expected changed paths:

- `go.mod`
- `go.sum`
- `internal/app/runtime.go`
- `internal/cli/*`
- `internal/contract/schemas.go`
- `internal/contract/shapes.go`
- `internal/model/output.go`
- `internal/output/*`
- `internal/service/registry.go`
- `internal/testutil/golden.go`
- Batch C coordination files

## Review Basis

Read:

- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/100-wp08-output-pipeline-and-error-surface.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/110-wp09-cli-command-surface.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/900-source-include-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-c-builder.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/009-batch-c-builder-report.md`

Include rule:

- if an include points to a whole file, read the whole file;
- if an include has `lines=`, read exactly those line ranges, no more and no less;
- if an include has several ranges separated by semicolons, read exactly every listed range.

Your report must list the files and include ranges used.

## Risk Focus

Review these risks directly:

- `internal/output` does not import `internal/source`.
- Output projections consume canonical objects or contract definitions, not source readers.
- JSON projection includes schema id and canonical object data; it is not serialized pretty text.
- jq consumes JSON projection and maps parse/eval failures to error surface.
- Shape/schema read `internal/contract`, not sample result values.
- Error surface includes command, object, source, state, meaning, recovery, and log path when present.
- Normal pretty text does not include recovery/tutorial text.
- Schema ids are stable and explicit.
- `internal/cli` does not import `internal/source`.
- Cobra handlers stay thin: parse args -> build input -> context -> service -> render -> exit mapping.
- Help/syntax surfaces do not expose hidden implementation URLs.
- Shared output-mode handling is consistent; `--jq` without expression is parse failure.
- Required argument and flag parse failures map to exit code `2`.
- Service/projection failures map to exit code `1`.
- Process exit and evidence state remain separate objects.
- Service registry seams do not accidentally implement Batch D behavior.
- Pay special attention to whether `internal/service/registry.go` importing `internal/source` is acceptable for the Batch B sync seam or creates boundary coupling that should move elsewhere before Batch C promotion.
- `just mod-check` dirty-diff failure is acceptable only if exactly caused by uncommitted Batch C dependency metadata and would pass after commit.

## Evidence Commands

Run and report:

- `git status --short`
- `go test ./...`
- `go test ./internal/output ./internal/contract`
- `go test ./internal/cli`
- `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/output ./internal/cli ./internal/service ./internal/app`
- `just fmt-check`
- `go mod tidy`
- `go mod verify`
- `goreleaser check`

If approved and committed, run `just mod-check` after the commit if practical and report it.

## Approval And Commit Rule

If approved:

1. Write report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/010-batch-c-review-report.md`
2. Stage new files if needed.
3. Commit only these paths:

```text
go.mod
go.sum
internal/app/runtime.go
internal/cli
internal/contract/schemas.go
internal/contract/shapes.go
internal/model/output.go
internal/output
internal/service/registry.go
internal/testutil/golden.go
bookshelf/books/10-seekmoon-wbs-work-packages/coordination
```

Use `git commit --only` with a detailed message naming WP08-WP09, output projection boundaries, CLI thin-handler boundaries, and evidence gates.

If rejected:

1. Write the same report path.
2. Do not commit.
3. Mark findings with severity and exact file/line evidence.
4. Include required action before re-review.

Final response: verdict, report path, and commit hash if created.
