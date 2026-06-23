# Batch C Revision 1 Builder Handoff

## Role

You are the persistent Batch C builder. Treat the independent review report as the controlling rejection basis.

You are not the principal coordinator and not the reviewer. Do not commit. Do not revert unrelated changes.

## Review Report To Address

Read:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/010-batch-c-review-report.md`

Original Batch C handoff remains binding:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-c-builder.md`

## Required Fixes

Fix only the Batch C rejection findings:

1. Contract projection modes
   - `--shape` and `--schema` must render command contract projections without ordinary command operands and without calling pending services.
   - This must work for output-capable command surfaces, including commands with normally required operands such as `api`, `view`, `source`, `probe`, `record`, `raw`, etc.
   - Contract projection modes should return success after rendering from contract definitions.
   - Add tests or behavior probes that would fail if `--shape`/`--schema` calls pending services or requires ordinary operands.

2. Real process exit mapping
   - The real `cmd/seekmoon` entrypoint must use the same exit mapping as the CLI execution path.
   - Required argument and flag parse failures must exit `2`.
   - Service/projection failures must exit `1`.
   - Success exits `0`.
   - Add tests that exercise the real command entrypoint behavior or an equivalent single execution path used by `cmd/seekmoon`.

## Verification Commands

Run and report:

- `go test ./...`
- `go test ./internal/output ./internal/contract`
- `go test ./internal/cli`
- `go run ./cmd/seekmoon search --shape`
- `go run ./cmd/seekmoon api --shape`
- a focused command proving `search` parse failure exits as `2` without relying on `go run`'s wrapper status ambiguity; use a small test/helper or shell command that reports the binary exit code directly
- `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/output ./internal/cli ./internal/service ./internal/app`
- `just fmt-check`
- `go mod tidy`
- `go mod verify`
- `goreleaser check`

## Report

Write the revision report at:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/011-batch-c-revision-1-builder-report.md`

Report format:

- Summary
- Review findings addressed
- Files changed
- Evidence commands
- Remaining risks or blocked gates
- Suggested reviewer focus

Final response should be short and include report path and changed paths.
