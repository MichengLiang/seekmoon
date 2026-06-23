# SeekMoon WBS Resource Map

## Source Of Truth

- WBS book root: `/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages`
- Specification book root: `/home/t103o/workbench/projects/seekmoon/bookshelf/books/09-seekmoon-cli-discovery-workbench`
- Implementation repository root: `/home/t103o/workbench/projects/seekmoon`

## Required WBS Files

- Global context: `010-mandatory-global-context.adoc`
- Dependency map: `020-wbs-dependency-map.adoc`
- WP01: `030-wp01-go-module-substrate.adoc`
- WP02: `040-wp02-canonical-model-and-contracts.adoc`
- WP03: `050-wp03-platform-runtime-and-storage.adoc`
- Source include map: `900-source-include-map.adoc`

## Global Full-File Includes From 9号书

The global context uses full-file includes for these sources. Executors and reviewers must read them completely when a work package consumes global context.

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

## Batch A Exact Line Includes

WP01:

- `parts/055-go-implementation-architecture/020-module-layout.adoc`, lines `4..174`
- `parts/055-go-implementation-architecture/120-engineering-quality-toolchain.adoc`, lines `4..73;118..159`
- `backmatter/appendix-g-go-engineering-toolchain.adoc`, lines `8..51;55..67;131..152;156..167`

WP02:

- `parts/030-canonical-model/010-evidence-state.adoc`, lines `4..62`
- `parts/030-canonical-model/020-module-summary-and-search-result.adoc`, lines `4..55;59..118`
- `parts/030-canonical-model/030-manifest-profile.adoc`, lines `4..100`
- `parts/030-canonical-model/040-api-profile.adoc`, lines `4..5;9..60;65..133`
- `parts/030-canonical-model/050-skill-profile.adoc`, lines `4..94`
- `parts/030-canonical-model/060-source-resolution-and-probe-result.adoc`, lines `4..5;8..72;77..132`
- `parts/030-canonical-model/070-record-and-report.adoc`, lines `4..64;69..97`
- `backmatter/appendix-b-data-dictionary.adoc`, lines `6..27;31..50;54..93;97..144;148..191;195..226;230..269;273..282`

WP03:

- `parts/055-go-implementation-architecture/040-runtime-composition.adoc`, lines `4..54`
- `parts/055-go-implementation-architecture/090-storage-and-side-effects.adoc`, lines `4..60`
- `parts/020-evidence-boundary/040-project-context.adoc`, lines `4..43`

## Build Objects

Batch A may create or modify:

- `go.mod`
- `go.sum`
- `justfile`
- `.golangci.yml`
- `.goreleaser.yaml`
- `.gitignore`
- `.github/workflows/*`
- `cmd/seekmoon/main.go`
- `internal/app/*`
- `internal/cli/*`
- `internal/model/*`
- `internal/contract/*`
- `internal/platform/*`
- `internal/store/*`
- `internal/testutil/*`
- tests under those package directories

Batch A must not implement full HTTP source readers, services, output renderers, or command behavior beyond minimal wiring needed to compile the substrate.

## Evidence Routes

- Builder report: `coordination/reports/001-batch-a-builder-report.md`
- Reviewer report: `coordination/reports/002-batch-a-review-report.md`
- Git evidence: commits in `/home/t103o/workbench/projects/seekmoon`
- Command evidence: output summaries recorded in reports; raw generated artifacts under `.artifacts/` when produced

## Batch B Exact Line Includes

WP04:

- `parts/020-evidence-boundary/010-public-http-sources.adoc`, lines `4..5;9..46;50..75;79..124`
- `parts/055-go-implementation-architecture/070-source-reader-boundary.adoc`, lines `4..42`
- `parts/055-go-implementation-architecture/050-data-flow.adoc`, lines `8..22;26..39`

WP05:

- `parts/020-evidence-boundary/020-static-assets.adoc`, lines `4..5;9..42;46..83;87..95;99..107;111..127`
- `parts/030-canonical-model/040-api-profile.adoc`, lines `4..5;9..60;65..133`
- `parts/040-command-workbench/040-api-and-source.adoc`, lines `4..32`

WP06:

- `parts/020-evidence-boundary/030-local-toolchain-and-cache.adoc`, lines `4..5;9..43;47..56;59..68`
- `parts/020-evidence-boundary/040-project-context.adoc`, lines `4..43`
- `parts/030-canonical-model/060-source-resolution-and-probe-result.adoc`, lines `77..132`

WP07:

- `parts/020-evidence-boundary/050-repository-signals.adoc`, lines `4..46`
- `parts/020-evidence-boundary/010-public-http-sources.adoc`, lines `128..169`
- `parts/030-canonical-model/050-skill-profile.adoc`, lines `4..94`
- `backmatter/appendix-f-go-implementation-dependencies.adoc`, lines `29..37;105..113`
