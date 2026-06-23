# Batch C Re-Review Report

## Verdict

Approved.

Batch C revision 1 fixes the two prior rejection findings. Contract projection modes now render from command contract definitions without ordinary operands or pending service calls, and the real `cmd/seekmoon` binary now uses the CLI exit-code mapping for parse, service, and projection outcomes.

## Review Inputs

- `coordination/review-packages/012-batch-c-re-review.md`
- `coordination/review-packages/010-batch-c-review.md`
- `coordination/reports/010-batch-c-review-report.md`
- `coordination/reports/011-batch-c-revision-1-builder-report.md`
- Prior Batch C WBS basis and exact include ranges listed in `coordination/reports/010-batch-c-review-report.md`
- Principal coordinator protocol `SKILL.md`, `references/index.md`, and `references/30_coordination_and_runtime/review_evidence_promotion.md`

## Re-Review Findings

No blocking findings.

## Prior Rejection Findings

### Contract Projections

Fixed.

Evidence:

- `internal/cli/flags.go:66` through `internal/cli/flags.go:84` adds shared contract projection detection, positional-argument bypass, and immediate shape/schema rendering.
- `internal/cli/search.go:15` through `internal/cli/search.go:33` now bypasses ordinary search input and service calls for contract projection modes.
- `internal/cli/api.go:15` through `internal/cli/api.go:31` now lets `api --shape`/`api --schema` render before module and package argument validation.
- Operand-required handlers such as `view`, `source`, `probe`, `record`, `raw`, `skill search`, `skill view`, and `compare` use the same bypass pattern.
- `internal/cli/root_test.go:69` through `internal/cli/root_test.go:105` covers contract projection bypass for operand-required commands and checks pending-service/argument-validation messages do not appear.
- `internal/output/shape.go:31` through `internal/output/shape.go:59` and `internal/output/schema.go:22` through `internal/output/schema.go:50` read command contract definitions for all Batch C command schema IDs.

Behavior evidence:

- `go run ./cmd/seekmoon search --shape` returned `seekmoon.search-results.v1`.
- `go run ./cmd/seekmoon api --shape` returned `seekmoon.package-data.v1`.
- Built-binary probe `./tmp/batch-c-rereview/seekmoon source --schema` exited `0` and returned `seekmoon.source-resolution.v1`.

### Real Process Exit Mapping

Fixed.

Evidence:

- `cmd/seekmoon/main.go:18` now exits through `cli.ExecuteWithCode(...)`.
- `internal/cli/root.go:60` through `internal/cli/root.go:74` maps usage failures to `2`, service/projection failures to `1`, and success to `0`.
- Built-binary probe for `seekmoon search` exited `2`.
- Built-binary probe for `seekmoon search markdown --json` exited `1` with a service error surface.

## Positive Boundary Checks

Re-checked and passed:

- `internal/output` does not import `internal/source`.
- `internal/cli` does not import `internal/source`.
- Output projections consume canonical objects or `internal/contract`, not source readers.
- JSON projection includes schema id and canonical object data; it is not serialized pretty text.
- jq consumes JSON projection and maps parse/eval failures to error surface.
- Shape/schema renderers read `internal/contract`.
- Error surface includes command, object, source, state, meaning, recovery when present, and log path when present.
- Normal pretty text omits recovery/tutorial text.
- Schema ids are stable and explicit.
- Cobra handlers remain thin and route through service interfaces plus renderer.
- Help/syntax surfaces do not expose hidden implementation URLs.
- Shared output-mode handling is consistent; `--jq` without expression remains a parse failure.
- Required argument and flag parse failures map to exit code `2`.
- Service/projection failures map to exit code `1`.
- Process exit and evidence state remain separate objects.
- Pending service seams do not implement concrete Batch D behavior.
- The `internal/service` import of `internal/source` remains the pre-existing Batch B sync seam; `internal/service/registry.go` itself does not import `internal/source`.

## Evidence Commands

`git status --short`

Result: only Batch C implementation and coordination paths were present before this report was written.

```text
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md
 M cmd/seekmoon/main.go
 M go.mod
 M go.sum
 M internal/app/runtime.go
 M internal/cli/root.go
 M internal/cli/root_test.go
 M internal/contract/schemas.go
 M internal/contract/shapes.go
 M internal/model/schema.go
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-c-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-c-revision-1-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/009-batch-c-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/010-batch-c-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/011-batch-c-revision-1-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/012-batch-c-re-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/009-batch-c-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/010-batch-c-review-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/011-batch-c-revision-1-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/010-batch-c-review.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/012-batch-c-re-review.md
?? internal/cli/api.go
?? internal/cli/compare.go
?? internal/cli/doctor.go
?? internal/cli/flags.go
?? internal/cli/probe.go
?? internal/cli/raw.go
?? internal/cli/record.go
?? internal/cli/report.go
?? internal/cli/search.go
?? internal/cli/skill.go
?? internal/cli/source.go
?? internal/cli/sync.go
?? internal/cli/view.go
?? internal/model/output.go
?? internal/output/
?? internal/service/registry.go
?? internal/testutil/golden.go
```

`go test ./...`

Result: passed.

`go test ./internal/output ./internal/contract`

Result: passed.

`go test ./internal/cli`

Result: passed.

`go run ./cmd/seekmoon search --shape`

Result: passed and rendered `seekmoon.search-results.v1`.

`go run ./cmd/seekmoon api --shape`

Result: passed and rendered `seekmoon.package-data.v1`.

Direct built-binary exit-code check:

```text
parse_exit=2
service_exit=1
source_schema_exit=0
```

`go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/output ./internal/cli ./internal/service ./internal/app`

Result:

```text
github.com/yumiaura/seekmoon/internal/output context encoding/json errors fmt github.com/itchyny/gojq github.com/yumiaura/seekmoon/internal/contract github.com/yumiaura/seekmoon/internal/model io strings
github.com/yumiaura/seekmoon/internal/cli context fmt github.com/spf13/cobra github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/output github.com/yumiaura/seekmoon/internal/service io strconv strings
github.com/yumiaura/seekmoon/internal/service context fmt github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/source github.com/yumiaura/seekmoon/internal/store time
github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/output github.com/yumiaura/seekmoon/internal/platform github.com/yumiaura/seekmoon/internal/service github.com/yumiaura/seekmoon/internal/source github.com/yumiaura/seekmoon/internal/store net/http time
```

`just fmt-check`

Result: passed.

```text
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal)"
```

`go mod tidy`

Result: passed with no command output.

`go mod verify`

Result: passed.

```text
all modules verified
```

`goreleaser check`

Result: passed.

```text
• checking                                  path=.goreleaser.yaml
• 1 configuration file(s) validated
• thanks for using GoReleaser!
```

## Commit

Approved for commit.
