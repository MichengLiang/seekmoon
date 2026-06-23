# Batch A Review Package

## Role

You are a specialized executor in reviewer role. You are independent from the Batch A builder and must not treat the builder report as proof by itself.

You are not the principal coordinator and do not directly address the user. You decide whether the Batch A implementation meets the review basis, then write your report. If it passes, you commit the approved change set. If it fails, do not commit; write findings that can be routed back to the builder.

## Review Object

Review the current uncommitted Batch A work in:

`/home/t103o/workbench/projects/seekmoon`

Builder report:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/001-batch-a-builder-report.md`

Expected changed areas:

- `.gitignore`
- `.golangci.yml`
- `.goreleaser.yaml`
- `go.mod`
- `go.sum`
- `justfile`
- `cmd/seekmoon/main.go`
- `internal/app/*`
- `internal/cli/*`
- `internal/contract/*`
- `internal/model/*`
- `internal/platform/*`
- `internal/store/*`
- coordination files under `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/`

## Review Basis

Read these first:

- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/030-wp01-go-module-substrate.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/040-wp02-canonical-model-and-contracts.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/050-wp03-platform-runtime-and-storage.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/900-source-include-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-a-builder.md`

Include rule:

- if an include points to a whole file, read the whole file;
- if an include has `lines=`, read exactly those line ranges, no more and no less;
- if an include has several ranges separated by semicolons, read exactly every listed range.

Your report must list the files and line ranges you used as review basis.

## Risk Focus

Check these risks directly:

- `internal/model` imports only Go standard library packages.
- State vocabulary is exactly `present`, `missing`, `unknown`, `failed`, `unavailable`, `derived`.
- Helpers do not collapse `missing`, `unknown`, `failed`, or `unavailable`.
- Evidence wrappers preserve `status`, `value`, `source`, and optional `error`.
- `internal/contract` schemas are explicit public contract objects, not reflection-derived from structs.
- `cmd/seekmoon` remains process startup only.
- `internal/cli` contains no source/service/output business behavior beyond Batch A compile seams.
- `internal/platform` and `internal/store` boundaries match WP03. Pay special attention to whether `internal/platform` importing `internal/model` violates the package-boundary diagram.
- Stores own mechanics and path construction only; they must not own adoption conclusions, source priority, or output shape decisions.
- Project storage paths stay under `.seekmoon/`; reusable cache paths stay under `$XDG_CACHE_HOME/seekmoon/`.
- Quality command names in `justfile` match the WBS exactly.
- Tests cover WP01-WP03 required responsibilities well enough for Batch A promotion.
- The repo has no unrelated staged or unstaged changes mixed into the review object.

## Evidence Commands

Run what is available and report exact results:

- `git status --short`
- `go test ./...`
- `go test ./internal/model ./internal/contract`
- `go test ./internal/platform ./internal/store ./internal/app`
- `just fmt-check`
- `just mod-check`
- `goreleaser check`

If a required external tool is missing, report it as a blocked gate. Do not count it as passing.

You may run additional focused commands if they materially support a finding.

## Approval And Commit Rule

If Batch A passes review:

1. Write the report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/002-batch-a-review-report.md`
2. Stage the approved new files if needed.
3. Commit only the Batch A and coordination paths using `git commit --only`.
4. Use a detailed commit message that names WP01-WP03 and explains the boundary decisions and remaining externally blocked gates.

Suggested path list for commit:

```text
.gitignore
.golangci.yml
.goreleaser.yaml
go.mod
go.sum
justfile
cmd
internal
bookshelf/books/10-seekmoon-wbs-work-packages/coordination
```

If Batch A fails review:

1. Write the report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/002-batch-a-review-report.md`
2. Do not commit.
3. Mark each finding with severity and exact file/line evidence where possible.
4. Include the action needed before re-review.

## Report Format

- Verdict: approved or rejected
- Review basis read
- Evidence commands
- Findings
- Boundary judgment
- Commit hash if approved
- Required follow-up if rejected or externally blocked

Your final chat response should be short and include the report path, verdict, and commit hash if one was created.
