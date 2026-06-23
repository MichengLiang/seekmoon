# Batch E Review Package

## Role

You are the independent reviewer for Batch E/WP13. Batches A-D are already approved and committed. Review current Batch E acceptance/quality-gate work and related coordination files.

You are not the builder and not the principal coordinator. If approved, write the report and commit only approved Batch E and coordination paths. If rejected, write the report and do not commit.

## Review Object

Repository:

`/home/t103o/workbench/projects/seekmoon`

Builder report:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/017-batch-e-builder-report.md`

Expected changed paths:

- `justfile`
- `tests/acceptance/*`
- `tests/blackbox/*`
- `tests/journey/*`
- `tests/integration/*`
- Batch E coordination files

Generated artifact:

- `.artifacts/coverage.out` may exist from coverage; it should not be committed unless the reviewer finds a policy reason.

## Review Basis

Read:

- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/150-wp13-black-box-acceptance-and-quality-gates.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/900-source-include-map.adoc`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-e-builder.md`
- `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/017-batch-e-builder-report.md`

Include rule:

- if an include points to a whole file, read the whole file;
- if an include has `lines=`, read exactly those line ranges, no more and no less;
- if an include has several ranges separated by semicolons, read exactly every listed range.

Your report must list the files and include ranges used.

## Risk Focus

Review these risks directly:

- A1-A16 all map to executable test names.
- Default tests do not depend on network, real Moon CLI mutation, GitHub token, or external services.
- Integration tests are opt-in, visibly named, and skip by default.
- Journey tests cover library, skill, pipeline, and failure recovery behavior without becoming WP13-only product behavior.
- Golden/schema/failure tests check public behavior, not implementation internals only.
- `justfile` composes the final quality gates and includes tests in relevant commands.
- `just ci` uses the same gates as local development.
- The generated `.artifacts/coverage.out` is not accidentally committed if it is only a generated quality output.
- Treat `just lint`, `just vuln`, and `just ci` as promotion-critical gates unless there is a stronger WP13-consistent reason to defer them. The builder report and coordinator probe show `lint` and `vuln` currently fail.

## Evidence Commands

Run and report:

- `git status --short`
- `go test ./...`
- `go test ./tests/blackbox ./tests/journey`
- `go test ./tests/integration`
- `PATH="$(go env GOPATH)/bin:$PATH" just fmt-check`
- `PATH="$(go env GOPATH)/bin:$PATH" just lint`
- `PATH="$(go env GOPATH)/bin:$PATH" just test`
- `PATH="$(go env GOPATH)/bin:$PATH" just test-race`
- `PATH="$(go env GOPATH)/bin:$PATH" just cover`
- `PATH="$(go env GOPATH)/bin:$PATH" just vuln`
- `PATH="$(go env GOPATH)/bin:$PATH" just mod-check`
- `PATH="$(go env GOPATH)/bin:$PATH" just release-check`
- `PATH="$(go env GOPATH)/bin:$PATH" just ci`

If rejecting because quality gates fail, include enough representative findings for the builder to repair the gate.

## Approval And Commit Rule

If approved:

1. Write report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/018-batch-e-review-report.md`
2. Stage new files if needed.
3. Commit only these paths:

```text
justfile
tests
bookshelf/books/10-seekmoon-wbs-work-packages/coordination
```

Use `git commit --only` with a detailed message naming WP13, A1-A16 acceptance, integration boundaries, and quality gates.

If rejected:

1. Write the same report path.
2. Do not commit.
3. Mark findings with severity and exact file/line or command evidence.
4. Include required action before re-review.

Final response: verdict, report path, and commit hash if created.
