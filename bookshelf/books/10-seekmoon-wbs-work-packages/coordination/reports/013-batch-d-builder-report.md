# Batch D Builder Report

## Summary

Implemented Batch D service behavior for:

- WP10: `doctor`, `sync`, `search`, `view`, `compare`
- WP11: `api`, `source`, `skill search`, `skill view`
- WP12: `probe`, `record`, `report`, `raw`

The runtime registry now wires concrete Batch D service flows instead of Batch C pending service placeholders. CLI syntax remains the Batch C/WP09 surface; service handlers now have behavior behind it. No commit was created.

## Files changed

Product implementation paths changed:

- `internal/app/runtime.go`
- `internal/model/output.go`
- `internal/output/error.go`
- `internal/output/pretty.go`
- `internal/source/mooncli.go`
- `internal/store/record_store.go`
- `internal/store/snapshot_store.go`
- `internal/service/api.go`
- `internal/service/batch_d_test.go`
- `internal/service/compare.go`
- `internal/service/doctor.go`
- `internal/service/helpers.go`
- `internal/service/probe.go`
- `internal/service/raw.go`
- `internal/service/record.go`
- `internal/service/report.go`
- `internal/service/search.go`
- `internal/service/skill.go`
- `internal/service/source.go`
- `internal/service/sync.go`
- `internal/service/view.go`

Removed/replaced:

- `internal/service/sync_test.go` was folded into the broader Batch D focused service test file.

Report path:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/013-batch-d-builder-report.md`

Pre-existing coordination artifact changes remain in the worktree and were not reverted.

## Reading evidence

Required WBS files read:

- `010-mandatory-global-context.adoc`
- `020-wbs-dependency-map.adoc`
- `120-wp10-discovery-and-profile-services.adoc`
- `130-wp11-inspection-and-source-services.adoc`
- `140-wp12-assessment-record-and-report-services.adoc`
- `900-source-include-map.adoc`
- `coordination/resource-map.md`
- `coordination/coordination-state.md`

Global full-file includes read from 9号书:

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

Batch D exact line includes read:

WP10:

- `parts/040-command-workbench/020-doctor-and-sync.adoc`, lines `4..43;47..83`
- `parts/040-command-workbench/030-search-and-view.adoc`, lines `4..44;48..74`
- `parts/040-command-workbench/060-compare-and-probe.adoc`, lines `4..28`
- `parts/055-go-implementation-architecture/050-data-flow.adoc`, lines `8..22;26..38;42..53;105..115`

WP11:

- `parts/040-command-workbench/040-api-and-source.adoc`, lines `4..33;37..78`
- `parts/040-command-workbench/050-skill.adoc`, lines `4..49`
- `parts/055-go-implementation-architecture/050-data-flow.adoc`, lines `57..68;72..87;91..101`
- `parts/060-journeys-and-acceptance/020-skill-discovery-journey.adoc`, lines `4..5;8..27;30..52;55..73`

WP12:

- `parts/040-command-workbench/060-compare-and-probe.adoc`, lines `32..58`
- `parts/040-command-workbench/070-record-report-and-raw.adoc`, lines `4..19;23..37;41..52`
- `parts/030-canonical-model/070-record-and-report.adoc`, lines `4..64;69..97`
- `parts/060-journeys-and-acceptance/040-failure-recovery-journey.adoc`, lines `63..91`

## Implementation notes by WP

### WP10

- Added `DoctorFlow`, which reports toolchain/command/path/network/project context status without creating snapshots or records.
- Extended `SyncService` to record Moon CLI update/version command evidence as raw snapshot metadata while preserving Modules API and Statistics API source state.
- Added `SearchFlow`, which reads latest snapshot when available, falls back to transient Modules API fetch, creates explainable match evidence, enriches visible results with Manifest API downloads/build status, keeps target support `unknown`, and writes session candidate mappings.
- Added `ViewFlow`, which reads Manifest API and module index asset summary without expanding package API details.
- Added `CompareFlow`, which aligns evidence fields across candidates and deliberately has no quality score or recommendation.

### WP11

- Added `APIFlow`, which resolves module/version input, loads module index before package data, validates package paths against the module index, derives package relpath via existing source/model logic, and returns a package-path error containing known packages.
- Added `SourceFlow`, which records `moon fetch`, source zip, local cache, core local source, and repository signal attempts separately, then selects a source deterministically from successful attempts.
- Added `SkillFlow`, which uses Skills API for search and view, writes skill candidate session mappings, fetches wasm/checksum asset states, and derives the pinned runwasm coordinate.

### WP12

- Added `ProbeFlow`, which creates isolated probe paths, records Moon command sequence results/log paths, maps target build failure to failed probe result, and does not write adoption records.
- Added `RecordFlow`, which validates the fixed conclusion enum, resolves candidate/session references, writes `AdoptionRecord`, and includes evidence refs and not-confirmed items.
- Added `ReportFlow`, which reads stored records and includes only evidence refs/sources actually present in records.
- Added `RawFlow`, which returns `RawEnvelope` values preserving source labels, URL/path, fetch state, raw ref, and upstream payload shape.

Shared support:

- Added model types for `EnvironmentStatus`, `Comparison`, and `RawEnvelope`.
- Added a model-level `SurfaceFailure` carrier so services can return structured error-surface data without importing `internal/output`.
- Added store helpers for latest snapshot and record listing.
- Added Moon CLI wrappers for `moon update` and `moon fetch`.
- Added pretty-text projections for Batch D canonical outputs.

## Tests and command evidence

### `go test ./...`

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

### `go test ./internal/service -run 'Test(Doctor|Sync|Search|View|Compare)'`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/service (cached)
```

### `go test ./internal/service -run 'Test(API|Source|Skill)'`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/service (cached)
```

### `go test ./internal/service -run 'Test(Probe|Record|Report|Raw)'`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/service (cached)
```

### `go test ./internal/cli`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/cli (cached)
```

### Import boundary listing

Command:

```bash
go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/service ./internal/cli ./internal/output ./internal/source
```

Result:

```text
github.com/yumiaura/seekmoon/internal/service context encoding/json errors fmt github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/source github.com/yumiaura/seekmoon/internal/store os path/filepath sort strings time
github.com/yumiaura/seekmoon/internal/cli context fmt github.com/spf13/cobra github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/output github.com/yumiaura/seekmoon/internal/service io strconv strings
github.com/yumiaura/seekmoon/internal/output context encoding/json errors fmt github.com/itchyny/gojq github.com/yumiaura/seekmoon/internal/contract github.com/yumiaura/seekmoon/internal/model io strings
github.com/yumiaura/seekmoon/internal/source archive/zip bufio bytes context encoding/json errors fmt github.com/cenkalti/backoff/v5 github.com/google/go-github/v88/github github.com/pelletier/go-toml/v2 github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/platform github.com/yumiaura/seekmoon/internal/store golang.org/x/oauth2 io net/http net/url path/filepath regexp strings time
```

### `just fmt-check`

Result: passed.

```text
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal)"
```

### `go mod tidy`

Result: passed with no command output.

### `go mod verify`

Result: passed.

```text
all modules verified
```

### `goreleaser check`

Result: passed.

```text
• checking                                  path=.goreleaser.yaml
• 1 configuration file(s) validated
• thanks for using GoReleaser!
```

## Known gaps or blocked gates

No required verification gate is blocked.

Known scope notes:

- Batch D uses focused fake HTTP readers and fake command runner tests, not WP13 black-box acceptance harnesses.
- Store listing helpers currently use `os.ReadDir` for OS-backed store directories; this matches current store mechanics but is less abstract than the `platform.FS` interface.
- `sync` records local index capability as a service field but does not yet aggregate a full local registry index count into `Snapshot.Raw`; Modules API and Statistics API source states remain covered.
- `record` derives conservative evidence refs and not-confirmed items from current candidate input rather than a richer user-specified evidence-ref flag surface, because the handoff disallowed new WP09 syntax expansion.

## Suggested reviewer focus

- Confirm services compose existing source readers and stores and do not render directly.
- Confirm `search` session candidate mapping, target `unknown` behavior, and visible Manifest API enrichment.
- Confirm `api` package-path failure reports known package paths from module index.
- Confirm `source` keeps failed/unavailable/unknown attempts distinct and treats repository as repository signal.
- Confirm `compare` output has no quality score or recommendation.
- Confirm `probe` does not mutate records and `report` does not invent repository/probe evidence.
