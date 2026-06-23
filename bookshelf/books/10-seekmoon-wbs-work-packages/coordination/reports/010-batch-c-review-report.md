# Batch C Independent Review Report

## Verdict

Rejected.

Batch C adds the requested output and CLI surfaces and the required package tests pass, but two WP08/WP09 command-surface contracts are not satisfied:

- contract projections such as `--shape` and `--schema` still require normal command arguments and/or call pending services;
- the real `cmd/seekmoon` entrypoint maps all command errors to process exit code `1`, so required argument and flag parse failures do not map to exit code `2`.

No commit was created.

## Review Basis Read

Required WBS and coordination files:

- `010-mandatory-global-context.adoc`
- `020-wbs-dependency-map.adoc`
- `100-wp08-output-pipeline-and-error-surface.adoc`
- `110-wp09-cli-command-surface.adoc`
- `900-source-include-map.adoc`
- `coordination/resource-map.md`
- `coordination/handoffs/batch-c-builder.md`
- `coordination/reports/009-batch-c-builder-report.md`

Batch C exact include ranges read:

WP08:

- `parts/050-output-contracts/010-pretty-text.adoc`, lines `4..45`
- `parts/050-output-contracts/020-json-output.adoc`, lines `4..58`
- `parts/050-output-contracts/030-built-in-jq.adoc`, lines `4..36`
- `parts/050-output-contracts/040-shape.adoc`, lines `4..47`
- `parts/050-output-contracts/050-schema.adoc`, lines `4..32`
- `parts/050-output-contracts/060-error-surface.adoc`, lines `4..59`
- `parts/055-go-implementation-architecture/080-output-pipeline.adoc`, lines `4..51`

WP09:

- `parts/040-command-workbench/010-command-map.adoc`, lines `4..73`
- `parts/055-go-implementation-architecture/060-control-flow.adoc`, lines `4..57`
- `backmatter/appendix-c-command-reference.adoc`, lines `4..84`

Additional review inputs:

- `coordination/prompts/010-batch-c-reviewer.md`
- `coordination/review-packages/010-batch-c-review.md`
- principal coordinator protocol `SKILL.md`, `references/index.md`, and `references/30_coordination_and_runtime/review_evidence_promotion.md`

## Evidence Commands

`git status --short`

Result: Batch C implementation and coordination changes are present; no unrelated path outside the review object was observed.

```text
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md
 M go.mod
 M go.sum
 M internal/app/runtime.go
 M internal/cli/root.go
 M internal/cli/root_test.go
 M internal/contract/schemas.go
 M internal/contract/shapes.go
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-c-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/009-batch-c-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/010-batch-c-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/009-batch-c-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/010-batch-c-review.md
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

```text
?    github.com/yumiaura/seekmoon/cmd/seekmoon [no test files]
ok   github.com/yumiaura/seekmoon/internal/app (cached)
ok   github.com/yumiaura/seekmoon/internal/cli (cached)
ok   github.com/yumiaura/seekmoon/internal/contract (cached)
ok   github.com/yumiaura/seekmoon/internal/model (cached)
ok   github.com/yumiaura/seekmoon/internal/output (cached)
ok   github.com/yumiaura/seekmoon/internal/platform (cached)
ok   github.com/yumiaura/seekmoon/internal/service (cached)
ok   github.com/yumiaura/seekmoon/internal/source (cached)
ok   github.com/yumiaura/seekmoon/internal/store (cached)
?    github.com/yumiaura/seekmoon/internal/testutil [no test files]
```

`go test ./internal/output ./internal/contract`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/output (cached)
ok   github.com/yumiaura/seekmoon/internal/contract (cached)
```

`go test ./internal/cli`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/cli (cached)
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

Additional focused behavior probes:

`go run ./cmd/seekmoon search --shape`

Result: failed by calling the pending search service instead of rendering the search contract shape.

```text
search service behavior is outside Batch C
command  seekmoon search
object   command
source   service
state    failed
meaning  search service behavior is outside Batch C
exit status 1
```

`go run ./cmd/seekmoon api --shape`

Result: failed argument validation before rendering the API contract shape.

```text
accepts 1 arg(s), received 0
exit status 1
```

`go run ./cmd/seekmoon search`

Result: parse failure exits as `1` from the actual command entrypoint, not `2`.

```text
search requires a query
exit status 1
```

## Findings

### High: Shape/schema contract projections are coupled to normal service execution or argument requirements

Evidence:

- `internal/cli/search.go:15` through `internal/cli/search.go:18` relaxes search argument validation for `--shape` or `--schema`, but `internal/cli/search.go:24` through `internal/cli/search.go:33` still resolves output mode, builds a search input, calls `rt.Services.Registry.Search.Search`, and only then renders.
- `internal/cli/api.go:12` through `internal/cli/api.go:15` uses `cobra.ExactArgs(1)` before output mode is resolved, so `seekmoon api --shape` cannot render the API contract without a module argument.
- The same ordinary-command pattern appears in handlers such as `view`, `source`, `skill view`, `compare`, `probe`, `record`, `report`, and `raw`: argument validation and service calls precede contract projection rendering.
- Runtime pending services intentionally return “outside Batch C” errors through `internal/service/registry.go:146` through `internal/service/registry.go:199`, which means shape/schema modes on pending command services produce error surfaces instead of contract projections.

Review basis:

- WP08 says shape output is a contract projection and “读取 shape 不需要执行真实查询.”
- WP08 says shape and schema are rendered from `internal/contract`, not sample result values.
- `parts/055-go-implementation-architecture/080-output-pipeline.adoc`, lines `4..51`, defines `--shape` and `--schema` as command contract projections, not service outputs.
- WP09 expects output mode options on output-capable commands.

Impact:

Users cannot inspect command output shapes/schemas for many Batch C command surfaces without providing ordinary command operands and triggering Batch D service seams. This breaks the Batch C purpose of establishing the output contract surface before concrete Batch D service behavior exists.

Required action before re-review:

Handle `--shape` and `--schema` before ordinary service invocation and before required data arguments that are unnecessary for contract projection. Contract projection modes should render from the command schema id and return success without calling pending services.

### High: Real process exit mapping does not return exit code 2 for parse and flag failures

Evidence:

- `internal/cli/root.go:60` through `internal/cli/root.go:74` implements `ExecuteWithCode`, which can distinguish usage failures as exit code `2`.
- `cmd/seekmoon/main.go:18` through `cmd/seekmoon/main.go:20` does not use `ExecuteWithCode`; it calls `cli.Execute`, prints any error, and exits `1`.
- Focused behavior probe `go run ./cmd/seekmoon search` produced `search requires a query` followed by `exit status 1`.
- Focused behavior probe `go run ./cmd/seekmoon search markdown --jq` produced a flag-argument error followed by `exit status 1`.

Review basis:

- WP09 requires required argument and flag parse failures to map to exit code `2`.
- `parts/055-go-implementation-architecture/060-control-flow.adoc`, lines `4..57`, distinguishes input parse failure as process exit code `2`, while core command object failure and jq parse/eval failure map to exit code `1`.
- The review package explicitly requires required argument and flag parse failures to map to exit code `2`.

Impact:

The package tests only verify `ExecuteWithCode`; the actual CLI binary surface collapses usage failures into generic command failure. This violates the public command contract and makes automation unable to distinguish syntax errors from service/projection failures.

Required action before re-review:

Wire `cmd/seekmoon/main.go` through the same exit mapping used by `ExecuteWithCode`, or expose a single execution path that returns the correct process status for parse failures, service/projection failures, and success.

## Positive Boundary Checks

These checks passed:

- `internal/output` does not import `internal/source`.
- `internal/cli` does not import `internal/source`.
- Output projections consume canonical objects or `internal/contract`, not source readers.
- JSON projection adds a stable schema id and projects canonical object data rather than pretty text.
- jq consumes JSON projection and maps parse/eval failures to error surface.
- Shape and schema functions read `internal/contract`.
- Error surface includes command, object, source, state, meaning, recovery when present, and log path when present.
- Normal pretty text tests check recovery/tutorial text is absent.
- Schema ids are stable constants from `internal/model`.
- Service registry seam in `internal/service/registry.go` does not import `internal/source`; the observed `internal/service` package import of `internal/source` comes from the pre-existing Batch B sync service.
- Pending service seams do not implement concrete Batch D behavior.

## Commit Hash If Approved

Not applicable. Rejected; no commit was created.

## Required Follow-Up Before Re-Review

1. Make `--shape` and `--schema` render contract projections without ordinary command operands or service calls.
2. Ensure the real `cmd/seekmoon` binary maps input parse and flag parse failures to process exit code `2`.
3. Add CLI tests that exercise the real command entrypoint behavior or an equivalent single execution path, not only `ExecuteWithCode`.
4. Re-run the full Batch C evidence command set.
