# Batch A Independent Review Report

## Verdict

Rejected.

Batch A has useful substrate work and the Go test gates pass, but it does not meet the review basis closely enough for approval. Two implementation findings block promotion:

- Evidence wrappers do not preserve the `source` field in JSON when the source is unknown/null.
- `internal/platform` imports `internal/model`, which violates the package-boundary direction in the 9号书 architecture basis.

Two required external gates are also blocked by missing tools and were not counted as passing: `just fmt-check` requires `gofumpt`, and `goreleaser check` requires `goreleaser`.

No commit was created.

## Review Basis Read

### Required WBS And Coordination Files

- `bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/030-wp01-go-module-substrate.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/040-wp02-canonical-model-and-contracts.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/050-wp03-platform-runtime-and-storage.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/900-source-include-map.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-a-builder.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/001-batch-a-builder-report.md`

### Global Full-File Includes Read From 9号书

- `parts/010-artifact-and-workbench/020-seekmoon-object-definition.adoc`
- `parts/010-artifact-and-workbench/030-consumer-actions.adoc`
- `parts/010-artifact-and-workbench/040-cli-workbench-environment.adoc`
- `parts/020-evidence-boundary/010-public-http-sources.adoc`
- `parts/020-evidence-boundary/020-static-assets.adoc`
- `parts/020-evidence-boundary/030-local-toolchain-and-cache.adoc`
- `parts/020-evidence-boundary/040-project-context.adoc`
- `parts/020-evidence-boundary/050-repository-signals.adoc`
- `parts/030-canonical-model/010-evidence-state.adoc`
- `parts/030-canonical-model/020-module-summary-and-search-result.adoc`
- `parts/030-canonical-model/030-manifest-profile.adoc`
- `parts/030-canonical-model/040-api-profile.adoc`
- `parts/030-canonical-model/050-skill-profile.adoc`
- `parts/030-canonical-model/060-source-resolution-and-probe-result.adoc`
- `parts/030-canonical-model/070-record-and-report.adoc`
- `parts/040-command-workbench/010-command-map.adoc`
- `parts/050-output-contracts/010-pretty-text.adoc`
- `parts/050-output-contracts/020-json-output.adoc`
- `parts/050-output-contracts/030-built-in-jq.adoc`
- `parts/050-output-contracts/040-shape.adoc`
- `parts/050-output-contracts/050-schema.adoc`
- `parts/050-output-contracts/060-error-surface.adoc`
- `parts/055-go-implementation-architecture/010-architecture-premises.adoc`
- `parts/055-go-implementation-architecture/020-module-layout.adoc`
- `parts/055-go-implementation-architecture/030-package-boundaries.adoc`
- `parts/055-go-implementation-architecture/040-runtime-composition.adoc`
- `parts/055-go-implementation-architecture/050-data-flow.adoc`
- `parts/055-go-implementation-architecture/060-control-flow.adoc`
- `parts/055-go-implementation-architecture/070-source-reader-boundary.adoc`
- `parts/055-go-implementation-architecture/080-output-pipeline.adoc`
- `parts/055-go-implementation-architecture/090-storage-and-side-effects.adoc`
- `parts/055-go-implementation-architecture/100-error-and-state-flow.adoc`
- `parts/055-go-implementation-architecture/110-testing-and-tooling.adoc`
- `parts/055-go-implementation-architecture/120-engineering-quality-toolchain.adoc`
- `backmatter/appendix-f-go-implementation-dependencies.adoc`
- `backmatter/appendix-g-go-engineering-toolchain.adoc`
- `backmatter/appendix-e-status-and-error-matrix.adoc`
- `parts/060-journeys-and-acceptance/010-library-discovery-journey.adoc`
- `parts/060-journeys-and-acceptance/020-skill-discovery-journey.adoc`
- `parts/060-journeys-and-acceptance/030-pipeline-journey.adoc`
- `parts/060-journeys-and-acceptance/040-failure-recovery-journey.adoc`
- `parts/060-journeys-and-acceptance/050-black-box-acceptance.adoc`

### Batch A Exact Line Includes Read

WP01:

- `parts/055-go-implementation-architecture/020-module-layout.adoc`, lines `4..174`
- `parts/055-go-implementation-architecture/120-engineering-quality-toolchain.adoc`, lines `4..73`
- `parts/055-go-implementation-architecture/120-engineering-quality-toolchain.adoc`, lines `118..159`
- `backmatter/appendix-g-go-engineering-toolchain.adoc`, lines `8..51`
- `backmatter/appendix-g-go-engineering-toolchain.adoc`, lines `55..67`
- `backmatter/appendix-g-go-engineering-toolchain.adoc`, lines `131..152`
- `backmatter/appendix-g-go-engineering-toolchain.adoc`, lines `156..167`

WP02:

- `parts/030-canonical-model/010-evidence-state.adoc`, lines `4..62`
- `parts/030-canonical-model/020-module-summary-and-search-result.adoc`, lines `4..55`
- `parts/030-canonical-model/020-module-summary-and-search-result.adoc`, lines `59..118`
- `parts/030-canonical-model/030-manifest-profile.adoc`, lines `4..100`
- `parts/030-canonical-model/040-api-profile.adoc`, lines `4..5`
- `parts/030-canonical-model/040-api-profile.adoc`, lines `9..60`
- `parts/030-canonical-model/040-api-profile.adoc`, lines `65..133`
- `parts/030-canonical-model/050-skill-profile.adoc`, lines `4..94`
- `parts/030-canonical-model/060-source-resolution-and-probe-result.adoc`, lines `4..5`
- `parts/030-canonical-model/060-source-resolution-and-probe-result.adoc`, lines `8..72`
- `parts/030-canonical-model/060-source-resolution-and-probe-result.adoc`, lines `77..132`
- `parts/030-canonical-model/070-record-and-report.adoc`, lines `4..64`
- `parts/030-canonical-model/070-record-and-report.adoc`, lines `69..97`
- `backmatter/appendix-b-data-dictionary.adoc`, lines `6..27`
- `backmatter/appendix-b-data-dictionary.adoc`, lines `31..50`
- `backmatter/appendix-b-data-dictionary.adoc`, lines `54..93`
- `backmatter/appendix-b-data-dictionary.adoc`, lines `97..144`
- `backmatter/appendix-b-data-dictionary.adoc`, lines `148..191`
- `backmatter/appendix-b-data-dictionary.adoc`, lines `195..226`
- `backmatter/appendix-b-data-dictionary.adoc`, lines `230..269`
- `backmatter/appendix-b-data-dictionary.adoc`, lines `273..282`

WP03:

- `parts/055-go-implementation-architecture/040-runtime-composition.adoc`, lines `4..54`
- `parts/055-go-implementation-architecture/090-storage-and-side-effects.adoc`, lines `4..60`
- `parts/020-evidence-boundary/040-project-context.adoc`, lines `4..43`

## Evidence Commands

`git status --short`

Result: working tree contains the expected Batch A and coordination changes, with no unrelated path outside the review object observed.

```text
 M .gitignore
?? .golangci.yml
?? .goreleaser.yaml
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/
?? cmd/
?? go.mod
?? go.sum
?? internal/
?? justfile
```

Expanded untracked status showed only Batch A implementation paths and coordination paths:

```text
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-a-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/001-batch-a-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/002-batch-a-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/001-batch-a-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/002-batch-a-review-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/002-batch-a-review.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md
```

`go test ./...`

Result: passed.

```text
?   	github.com/yumiaura/seekmoon/cmd/seekmoon	[no test files]
ok  	github.com/yumiaura/seekmoon/internal/app	(cached)
ok  	github.com/yumiaura/seekmoon/internal/cli	(cached)
ok  	github.com/yumiaura/seekmoon/internal/contract	(cached)
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
ok  	github.com/yumiaura/seekmoon/internal/platform	(cached)
ok  	github.com/yumiaura/seekmoon/internal/store	(cached)
```

`go test ./internal/model ./internal/contract`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
ok  	github.com/yumiaura/seekmoon/internal/contract	(cached)
```

`go test ./internal/platform ./internal/store ./internal/app`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/platform	(cached)
ok  	github.com/yumiaura/seekmoon/internal/store	(cached)
ok  	github.com/yumiaura/seekmoon/internal/app	(cached)
```

`just fmt-check`

Result: blocked by missing external tool `gofumpt`; not counted as passing.

```text
command -v gofumpt >/dev/null
error: Recipe `fmt-check` failed on line 8 with exit code 127
```

`just mod-check`

Result: passed.

```text
go mod tidy
git diff --exit-code -- go.mod go.sum
go mod verify
all modules verified
```

`goreleaser check`

Result: blocked by missing external tool `goreleaser`; not counted as passing.

```text
/bin/bash: line 1: goreleaser: command not found
```

Additional focused command:

`go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/model ./internal/contract ./internal/platform ./internal/store ./internal/app ./internal/cli`

Result: confirms `internal/model` uses only standard library imports; also confirms the boundary violation `internal/platform -> internal/model`.

```text
github.com/yumiaura/seekmoon/internal/model encoding/json errors fmt strings time
github.com/yumiaura/seekmoon/internal/contract github.com/yumiaura/seekmoon/internal/model
github.com/yumiaura/seekmoon/internal/platform bytes context github.com/yumiaura/seekmoon/internal/model net/http os os/exec path/filepath time
github.com/yumiaura/seekmoon/internal/store context encoding/json github.com/adrg/xdg github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/platform path/filepath strings
github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/platform github.com/yumiaura/seekmoon/internal/store net/http time
github.com/yumiaura/seekmoon/internal/cli context fmt github.com/spf13/cobra github.com/yumiaura/seekmoon/internal/app io
```

## Findings

### High: Evidence wrappers omit `source`, so they do not preserve the required wrapper contract

Evidence:

- `internal/model/evidence.go:77` defines `Evidence[T]`.
- `internal/model/evidence.go:80` declares `Source string `json:"source,omitempty"``.
- `internal/model/evidence.go:92` through `internal/model/evidence.go:93` constructs `unknown` evidence with no source value.
- `internal/contract/schemas.go:53` through `internal/contract/schemas.go:61` defines evidence schema with required fields `status` and `value`, but not `source`.

Review basis:

- `backmatter/appendix-b-data-dictionary.adoc`, lines `6..27`, defines the evidence wrapper shared fields as `status`, `value`, `source`, and optional `error`, with `source` typed as `string/null`.
- `parts/050-output-contracts/020-json-output.adoc` shows `target` evidence with `"source": null` for `unknown`.
- WP02 review risk explicitly requires evidence wrappers to preserve `status`, `value`, `source`, and optional `error`.

Impact:

Current JSON omits `source` entirely when it is empty, so downstream contract readers and pipeline consumers cannot distinguish an explicitly null source from an absent field. This breaks the public JSON contract surface for states such as `unknown`, and the contract schema does not catch it.

Required action:

Represent evidence source as a nullable JSON field that is always emitted, for example `*string `json:"source"``, and update constructors, schema, and tests so `source` is preserved as either a string or null. Keep `error` optional.

### High: `internal/platform` imports `internal/model`, reversing the package-boundary direction

Evidence:

- `internal/platform/exec.go:3` through `internal/platform/exec.go:10` imports `github.com/yumiaura/seekmoon/internal/model`.
- `internal/platform/exec.go:19` through `internal/platform/exec.go:21` exposes `Runner.Run` returning `model.CommandResult`.
- `internal/platform/exec.go:27` through `internal/platform/exec.go:67` maps host execution directly into `model.CommandResult` and `model.State`.
- Focused import command confirms `github.com/yumiaura/seekmoon/internal/platform ... github.com/yumiaura/seekmoon/internal/model ...`.

Review basis:

- `parts/055-go-implementation-architecture/030-package-boundaries.adoc` defines `internal/platform` as host primitives below SeekMoon business packages, and does not list `internal/platform -> internal/model`.
- The same package-boundary source says `internal/model` imports only Go standard library packages, and packages such as `internal/source`, `internal/store`, and `internal/service` depend on both model and platform.
- `parts/055-go-implementation-architecture/070-source-reader-boundary.adoc` assigns Moon CLI reader behavior to source readers using `platform.Runner`; command result belongs to local command evidence.

Impact:

The platform package now knows the canonical evidence vocabulary and `CommandResult` object. That makes a host primitive package depend on the model layer, rather than returning a platform-local process result that a source/service layer maps into `model.CommandResult`.

Required action:

Move the canonical command-result mapping out of `internal/platform`. `platform.Runner` should return a platform-local execution result and error. A higher package that is allowed to depend on both `platform` and `model` should map exit code/log path/error into `model.CommandResult` and `model.State`.

### Medium: Required quality gates are externally blocked

Evidence:

- `just fmt-check` failed because `gofumpt` is not installed.
- `goreleaser check` failed because `goreleaser` is not installed.

Review basis:

- WP01 completion evidence includes `just fmt-check`, `just mod-check`, and `goreleaser check`.
- The review package says missing required external tools must be reported as blocked gates and must not be counted as passing.

Required action:

Install the pinned tools from `backmatter/appendix-g-go-engineering-toolchain.adoc` and re-run the gates:

```bash
go install mvdan.cc/gofumpt@v0.10.0
go install github.com/goreleaser/goreleaser/v2@v2.16.0
just fmt-check
goreleaser check
```

## Boundary Judgment

Positive boundary checks:

- `cmd/seekmoon/main.go` remains process startup only.
- `internal/cli/root.go` contains placeholder command seams and no source/service/output business behavior.
- `internal/model` imports only Go standard library packages.
- State vocabulary is exactly `present`, `missing`, `unknown`, `failed`, `unavailable`, and `derived`.
- `internal/contract` schemas are explicit objects, not reflection-derived from structs.
- Store paths stay under `.seekmoon/` for project storage and `$XDG_CACHE_HOME/seekmoon/` for reusable cache.
- Stores mostly own path construction and read/write mechanics, not adoption conclusions, source priority, or output shape decisions.

Blocking boundary checks:

- Evidence wrapper JSON does not preserve `source` as a required nullable field.
- `internal/platform` imports `internal/model`, which violates the architecture package-boundary direction.

## Commit Hash If Approved

Not applicable. Rejected; no commit was created.

## Required Follow-Up

Before re-review:

1. Fix evidence wrapper source preservation and update schema/tests.
2. Remove the `internal/platform -> internal/model` dependency by moving command-result mapping above platform.
3. Install or otherwise provide the required external tools and re-run `just fmt-check` and `goreleaser check`.
4. Re-run the full review evidence command set.
