# Batch A Builder Report

## Summary

Implemented Batch A substrate for WP01, WP02, and WP03 in `/home/t103o/workbench/projects/seekmoon`.

The implementation establishes:

- Go module metadata and root quality entrypoints.
- A thin `cmd/seekmoon` process entrypoint.
- A Cobra root command with placeholder subcommands only; no Batch B-D behavior is implemented.
- Explicit canonical model objects, evidence state vocabulary, evidence wrappers, schema IDs, and shape/schema contract definitions.
- Platform capabilities for clock, filesystem, HTTP client construction, command running, and environment summary.
- Store path construction and file mechanics for project storage, reusable cache, snapshots, sessions, records, reports, probes, sources, and logs.
- Runtime composition in `internal/app`.
- Focused unit tests for state parsing/semantics, evidence JSON, coordinate parsing, package relpath derivation, adoption conclusion validation, contract shape/schema compilation, path construction, atomic writes, session candidate mapping, sanitized record/report/log paths, and runtime construction.

The Go test gates passed. Some quality gates are blocked by missing external tools and are listed below.

## Files Changed

Created or modified:

- `.gitignore`
- `.golangci.yml`
- `.goreleaser.yaml`
- `go.mod`
- `go.sum`
- `justfile`
- `cmd/seekmoon/main.go`
- `internal/app/runtime.go`
- `internal/app/runtime_test.go`
- `internal/cli/root.go`
- `internal/cli/root_test.go`
- `internal/contract/shapes.go`
- `internal/contract/schemas.go`
- `internal/contract/contract_test.go`
- `internal/model/evidence.go`
- `internal/model/module.go`
- `internal/model/manifest.go`
- `internal/model/api.go`
- `internal/model/skill.go`
- `internal/model/source.go`
- `internal/model/probe.go`
- `internal/model/record.go`
- `internal/model/report.go`
- `internal/model/snapshot.go`
- `internal/model/project.go`
- `internal/model/schema.go`
- `internal/model/evidence_test.go`
- `internal/model/derivation_test.go`
- `internal/platform/clock.go`
- `internal/platform/fs.go`
- `internal/platform/http.go`
- `internal/platform/exec.go`
- `internal/platform/env.go`
- `internal/platform/platform_test.go`
- `internal/store/paths.go`
- `internal/store/atomic.go`
- `internal/store/snapshot_store.go`
- `internal/store/session_store.go`
- `internal/store/record_store.go`
- `internal/store/report_store.go`
- `internal/store/cache_store.go`
- `internal/store/probe_store.go`
- `internal/store/source_store.go`
- `internal/store/paths_test.go`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/001-batch-a-builder-report.md`

Existing untracked process files under `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/` were present before implementation; this report file was added inside that directory as required by the handoff.

## Reading Evidence

### Required WBS and Coordination Files Read

- `bookshelf/books/10-seekmoon-wbs-work-packages/010-mandatory-global-context.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/020-wbs-dependency-map.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/030-wp01-go-module-substrate.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/040-wp02-canonical-model-and-contracts.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/050-wp03-platform-runtime-and-storage.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/900-source-include-map.adoc`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-a-builder.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/001-batch-a-builder.md`

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

## Implementation Notes By WP

### WP01: Go Module Substrate

- Added module `github.com/yumiaura/seekmoon` with Go `1.26.3`.
- Added root engineering files:
  - `justfile` with exact WBS entry names: `fmt-check`, `lint`, `test`, `test-race`, `cover`, `fuzz`, `vuln`, `mod-check`, `release-check`, `ci`, plus `fmt` and `release-snapshot`.
  - `.golangci.yml` with the Appendix G lint surface.
  - `.goreleaser.yaml` with a checkable GoReleaser v2 release surface.
  - `.gitignore` entries for `.artifacts/`, local binary, and `dist/`.
- Added `cmd/seekmoon/main.go` as process startup only: construct runtime and execute CLI root.
- No source reader, service, output renderer, or business command behavior was added.

### WP02: Canonical Model And Contracts

- Added closed `model.State` vocabulary preserving exactly:
  - `present`
  - `missing`
  - `unknown`
  - `failed`
  - `unavailable`
  - `derived`
- Added generic `Evidence[T]` wrapper with `status`, `value`, `source`, and optional `error`.
- Added model files for module/search, manifest, API, skill, source resolution, probe result, adoption record, report, snapshot, project context, and schema IDs.
- Added explicit coordinate and relpath derivation helpers:
  - `ParseModuleCoordinate`
  - `PackageRelPath`
  - `RunwasmCoordinate`
- Added adoption conclusion enum and validation.
- Added `internal/contract` shape and schema objects with explicit schema IDs and no struct reflection.
- Added tests for state parsing, state semantic separation, evidence JSON round-trip, module coordinate parsing, package relpath derivation, adoption conclusion validation, runwasm derivation, fuzz targets, shape definitions, and schema compilation.
- Verified `internal/model` imports only Go standard library packages.

### WP03: Platform Runtime And Storage

- Added platform capabilities:
  - `Clock` / `SystemClock`
  - `FS` / `OSFS`
  - HTTP client construction with timeout
  - `Runner` / `ExecRunner`
  - `Env` summary and integration flag access
- Added store path construction:
  - project `.seekmoon/` directories: snapshots, sessions, records, reports, probes, sources, logs
  - reusable `$XDG_CACHE_HOME/seekmoon/` directories: mooncakes, assets, github
- Added atomic write mechanics using temp file plus rename.
- Added snapshot, session, record, report, cache, probe, and source stores.
- Added runtime construction in `internal/app` aligned to:
  - process environment
  - path resolution
  - platform clients
  - stores
  - runtime value
  - CLI command tree
- Source readers, services, and renderers are intentionally absent because the handoff write boundary excludes Batch B-D behavior. Runtime composition currently exposes the platform and store substrate those later packages will consume.
- Noted for reviewer: `internal/platform` imports `internal/model` for `model.CommandResult` mapping in `ExecRunner`, matching WP03’s command result mapping expectation. The 9号书 dependency diagram says `internal/platform` owns host primitives and is below business packages, so reviewer should decide whether this cross-edge is acceptable in Batch A or should be split into a platform-local command result later.

## Tests And Command Evidence

### Required Commands

`go test ./...`

Result: passed.

Output summary:

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

Output summary:

```text
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
ok  	github.com/yumiaura/seekmoon/internal/contract	0.005s
```

`go test ./internal/platform ./internal/store ./internal/app`

Result: passed.

Output summary:

```text
ok  	github.com/yumiaura/seekmoon/internal/platform	(cached)
ok  	github.com/yumiaura/seekmoon/internal/store	(cached)
ok  	github.com/yumiaura/seekmoon/internal/app	(cached)
```

`just fmt-check`

Result: blocked by missing tool, correctly failed after recipe guard was added.

Output summary:

```text
command -v gofumpt >/dev/null
error: Recipe `fmt-check` failed on line 8 with exit code 127
```

`just mod-check`

Result: passed.

Output summary:

```text
go mod tidy
git diff --exit-code -- go.mod go.sum
go mod verify
all modules verified
```

`goreleaser check`

Result: blocked by missing tool.

Output summary:

```text
/bin/bash: line 1: goreleaser: command not found
```

### Additional Evidence

`go version`

```text
go version go1.26.3 linux/amd64
```

Tool availability probe:

```text
just found at /home/t103o/.cargo/bin/just
gofumpt not found
gotestsum not found
golangci-lint not found
goreleaser not found
```

Import boundary check:

```text
github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/platform github.com/yumiaura/seekmoon/internal/store net/http time
github.com/yumiaura/seekmoon/internal/cli context fmt github.com/spf13/cobra github.com/yumiaura/seekmoon/internal/app io
github.com/yumiaura/seekmoon/internal/contract github.com/yumiaura/seekmoon/internal/model
github.com/yumiaura/seekmoon/internal/model encoding/json errors fmt strings time
github.com/yumiaura/seekmoon/internal/platform bytes context github.com/yumiaura/seekmoon/internal/model net/http os os/exec path/filepath time
github.com/yumiaura/seekmoon/internal/store context encoding/json github.com/adrg/xdg github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/platform path/filepath strings
```

## Known Gaps Or Blocked Gates

- `just fmt-check` cannot pass until `gofumpt` is installed. Source was formatted with `gofmt`, not `gofumpt`.
- `goreleaser check` cannot run until `goreleaser` is installed.
- `just test`, `just test-race`, `just lint`, and `just vuln` were not part of the explicit handoff verification list, but their required tools are currently missing:
  - `gotestsum`
  - `golangci-lint`
  - `govulncheck`
- Batch A intentionally does not include full source readers, service flows, output rendering, or real CLI command behavior. Placeholder subcommands return an error that behavior is outside Batch A.
- Runtime registry scaffolding covers stores and host capabilities. Source/service/renderer registries are not added because the handoff write boundary excludes `internal/source`, `internal/service`, and `internal/output`.
- Reviewer should evaluate whether `internal/platform -> internal/model` for `CommandResult` is acceptable. If the dependency diagram is treated strictly, a follow-up should move runner raw result into `platform` and map to `model.CommandResult` in a higher package.

## Suggested Reviewer Focus

- Confirm the state vocabulary remains exact and that `missing`, `unknown`, `failed`, and `unavailable` are not collapsed by helpers or stores.
- Confirm `internal/model` has only standard-library imports.
- Confirm schemas are explicit contract objects rather than reflection output.
- Review whether `internal/platform` should import `internal/model` for command result status mapping in Batch A.
- Confirm placeholder CLI commands are acceptable compile seams and do not implement Batch B-D behavior.
- Confirm store paths match `.seekmoon/` and `$XDG_CACHE_HOME/seekmoon/` contracts and that stores do not make adoption/source-priority decisions.
- Install or provide `gofumpt` and `goreleaser` before promotion if the promotion boundary requires all quality gates passing locally rather than reported as missing.
