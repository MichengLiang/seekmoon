# Batch E Revision 1 Builder Handoff

## Role

You are the persistent Batch E builder. Treat the independent review report as the controlling rejection basis.

You are not the principal coordinator and not the reviewer. Do not commit. Do not revert unrelated changes.

## Review Report To Address

Read:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/018-batch-e-review-report.md`

Original Batch E handoff remains binding:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-e-builder.md`

## Required Fixes

Fix the WP13 promotion-critical gates:

1. `just lint`
   - Make `PATH="$(go env GOPATH)/bin:$PATH" just lint` pass.
   - Address the reported lint findings without weakening the public behavior or hiding real defects.
   - Prefer small, principled fixes over disabling linters. If a suppression is truly warranted, it must include a concise reason tied to the contract or test fixture.

2. `just vuln`
   - Make `PATH="$(go env GOPATH)/bin:$PATH" just vuln` pass.
   - The reviewer identified Go standard-library vulnerabilities in local `go1.26.3`, fixed in `go1.26.4`.
   - Update the Go toolchain to a fixed 1.26.x version if available, and align module/toolchain metadata if needed.
   - If a toolchain update is impossible in the current environment, report the exact blocker and do not claim the gate passes.

3. `just ci`
   - Make `PATH="$(go env GOPATH)/bin:$PATH" just ci` pass end-to-end.

Preserve the positive Batch E checks:

- A1-A16 executable test names remain present.
- Default tests remain offline and side-effect free.
- Integration tests remain opt-in and skip by default.
- `fmt-check`, `test`, `test-race`, `cover`, `mod-check`, and `release-check` remain passing.

## Verification Commands

Run and report:

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

## Report

Write:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/019-batch-e-revision-1-builder-report.md`

Report format:

- Summary
- Review findings addressed
- Files changed
- Toolchain changes, if any
- Quality gate evidence
- Acceptance regression evidence
- Remaining risks or blocked gates
- Suggested reviewer focus

Final response should be short and include report path and changed paths.
