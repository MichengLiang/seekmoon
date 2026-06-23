# Batch E Re-Review Package

## Role

You are the independent Batch E reviewer. Reuse your prior Batch E review context, but treat the current worktree, Revision 1 builder report, and this package as authoritative.

You are not the builder and not the principal coordinator. If approved, write the report and commit only approved Batch E, quality-gate remediation, and coordination paths. If rejected, write the report and do not commit.

## Review Inputs

- Prior review package: `coordination/review-packages/018-batch-e-review.md`
- Prior review report: `coordination/reports/018-batch-e-review-report.md`
- Revision builder report: `coordination/reports/019-batch-e-revision-1-builder-report.md`
- Revision handoff: `coordination/handoffs/batch-e-revision-1-builder.md`

## Review Object

Repository:

`/home/t103o/workbench/projects/seekmoon`

Review the current worktree changes for WP13 after Revision 1. The revised object intentionally includes both the original WP13 acceptance/test work and the broader lint/security remediation needed to make the WP13 quality gates pass.

Expected changed path groups:

- `cmd/seekmoon/main.go`
- `internal/app/*`
- `internal/cli/*`
- `internal/contract/*`
- `internal/model/*`
- `internal/output/*`
- `internal/platform/*`
- `internal/service/*`
- `internal/source/*`
- `internal/store/*`
- `internal/testutil/*`
- `justfile`
- `tests/acceptance/*`
- `tests/blackbox/*`
- `tests/journey/*`
- `tests/integration/*`
- Batch E coordination files

Generated artifact:

- `.artifacts/coverage.out` may exist from coverage; it should not be committed unless you find a policy reason.

## Include Reading Rule

You must preserve the WBS include-reading contract:

- if an include points to a whole file, read that entire file;
- if an include has `lines=`, read exactly those line ranges, no more and no less;
- if an include has multiple ranges separated by semicolons, read each listed range exactly.

Your report must list the files and include ranges used.

## Required Re-Review Checks

Verify the rejected findings are fixed:

1. `just lint`
   - `PATH="$(go env GOPATH)/bin:$PATH" just lint` passes.
   - Lint remediation does not introduce unrelated behavior drift.
   - Exported-surface comments add useful contract or API context where possible and do not hide design churn.
   - `#nosec` comments are contract-bound and justified, especially in:
     - `internal/platform/exec.go`
     - `internal/platform/fs.go`
     - `internal/source/local_index.go`
     - `internal/testutil/golden.go`

2. `just vuln`
   - `PATH="$(go env GOPATH)/bin:$PATH" just vuln` passes.
   - Local toolchain evidence shows `go1.26.4` or later.
   - Keeping `go.mod` at the existing directive is acceptable only if `just mod-check` passes and no module metadata diff is required by the project gates.

3. `just ci`
   - `PATH="$(go env GOPATH)/bin:$PATH" just ci` passes end-to-end.
   - `ci` still composes the WP13 quality gates rather than skipping them.

Also re-check the positive Batch E acceptance boundaries from the initial review before approving:

- A1-A16 all map to visible executable test names.
- Default tests do not depend on network, real Moon CLI mutation, GitHub token, or external services.
- Integration tests are opt-in, visibly named, and skip by default.
- Journey tests cover library, skill, pipeline, and failure recovery behavior without becoming WP13-only product behavior.
- Golden/schema/failure tests check public behavior, not implementation internals only.
- `justfile` includes tests in relevant formatting and quality commands.
- `.artifacts/coverage.out` is not accidentally staged or committed if generated.

## Evidence Commands

Run and report:

- `git status --short`
- `go version`
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

If approved and committed, run `git status --short` and `PATH="$(go env GOPATH)/bin:$PATH" just ci` after commit if practical and report them.

## Approval And Commit Rule

If approved:

1. Write report at:
   `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/020-batch-e-re-review-report.md`
2. Stage new files if needed.
3. Commit only these paths:

```text
cmd/seekmoon/main.go
internal
justfile
tests
bookshelf/books/10-seekmoon-wbs-work-packages/coordination
```

Use `git commit --only` with a detailed message naming WP13, A1-A16 acceptance, opt-in integration boundaries, lint/vulnerability remediation, and final quality gates.

If rejected:

1. Write the same report path.
2. Do not commit.
3. Mark findings with severity and exact file/line or command evidence.
4. Include required action before any next re-review.

Final response: verdict, report path, and commit hash if created.
