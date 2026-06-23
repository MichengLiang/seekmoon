# Batch D Review Report

## Verdict

Rejected.

Batch D adds concrete WP10-WP12 service flows and the required evidence commands pass, but the review object does not meet two promotion conditions from the Batch D review basis:

- `sync` does not read or record the local registry index summary required by WP10.
- `raw` returns normalized model payloads instead of upstream raw payload shapes required by WP12.

No commit was created.

## Review Inputs

- `coordination/prompts/014-batch-d-reviewer.md`
- `coordination/review-packages/014-batch-d-review.md`
- `010-mandatory-global-context.adoc`
- `020-wbs-dependency-map.adoc`
- `120-wp10-discovery-and-profile-services.adoc`
- `130-wp11-inspection-and-source-services.adoc`
- `140-wp12-assessment-record-and-report-services.adoc`
- `900-source-include-map.adoc`
- `coordination/resource-map.md`
- `coordination/handoffs/batch-d-builder.md`
- `coordination/reports/013-batch-d-builder-report.md`
- Principal coordinator protocol `SKILL.md`, `references/index.md`, and `references/30_coordination_and_runtime/review_evidence_promotion.md`

## Exact Include Ranges Used

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

Additional focused evidence ranges rechecked during review:

- `parts/040-command-workbench/020-doctor-and-sync.adoc`, lines `47..83`
- `parts/055-go-implementation-architecture/050-data-flow.adoc`, lines `8..22`
- `parts/040-command-workbench/070-record-report-and-raw.adoc`, lines `41..52`

## Findings

### High: `sync` omits the local registry index source and summary

`sync` cannot be promoted because the implemented flow does not read or record the local registry index, even though WP10 requires it as part of the evidence snapshot.

Review basis:

- `parts/040-command-workbench/020-doctor-and-sync.adoc:47` says `sync` creates a dated evidence snapshot fixing the current API, statistics, and local index view.
- `parts/040-command-workbench/020-doctor-and-sync.adoc:54` says `sync` executes `moon update`, reads Modules API, reads Statistics API, counts the local index, and saves a snapshot file.
- `parts/040-command-workbench/020-doctor-and-sync.adoc:66` through `:67` requires Local registry index data in the snapshot: index HEAD, file count, and record count.
- `parts/055-go-implementation-architecture/050-data-flow.adoc:10` through `:17` places `local registry index summary` between the source reads and `Snapshot`.
- `parts/040-command-workbench/020-doctor-and-sync.adoc:83` requires partial source failures to be recorded in snapshot state.

Implementation evidence:

- `internal/service/sync.go:13` through `:20` declares `LocalIndex source.LocalIndexReader`.
- `internal/service/sync.go:27` through `:34` runs Moon CLI update/version and fetches Modules API and Statistics API only.
- `internal/service/sync.go:38` through `:41` writes only Modules API and Statistics API into `snapshot.Sources`.
- `internal/service/sync.go:42` through `:52` writes only moon command results and optional modules/statistics values into `snapshot.Raw`.
- `internal/service/sync.go:56` writes the snapshot without any `s.LocalIndex.Read(...)` call or local-index source result.
- `internal/source/local_index.go:46` through `:57` already provides a `LocalIndexReader.Read` source result with status, parse state, raw ref, error, and value, but Batch D `sync` does not use it.

Impact:

The snapshot cannot express local index presence, failure, or unavailability. Downstream commands that rely on snapshot state cannot distinguish a missing local index from an implementation omission, and the risk focus requirement that `sync` preserves partial source failures is not satisfied for this source.

Required action before re-review:

Read the local registry index during `sync`, append its state-bearing source result to `snapshot.Sources`, and record the required local index summary facts in the snapshot payload, including status/error on failed or unavailable reads. Add tests covering present and failed/unavailable local index behavior.

### High: `raw` returns normalized model payloads instead of upstream payload shape

`raw` cannot be promoted because it routes through the normalized source-reader methods and places the normalized `Value` in `RawEnvelope.Payload`.

Review basis:

- `parts/040-command-workbench/070-record-report-and-raw.adoc:41` says `raw` exposes original payload for source audit, field recheck, and failure reproduction.
- `parts/040-command-workbench/070-record-report-and-raw.adoc:52` requires raw output to preserve upstream field names, not normalize fields, and not hide original shape.
- The Batch D risk focus requires `raw` to preserve upstream payload names and source status and not become normal discovery flow.

Implementation evidence:

- `internal/service/raw.go:16` through `:36` calls `FetchModules`, `FetchManifest`, `FetchModuleIndex`, `FetchPackageData`, and `FetchSkills`, then passes `result.Value` to `rawEnvelope`.
- `internal/source/mooncakes.go:19` through `:32` normalizes the modules response into `[]model.ModuleSummary`.
- `internal/source/mooncakes.go:44` through `:62` normalizes manifest responses into `model.ManifestProfile`.
- `internal/source/assets.go:23` through `:36` normalizes module index responses into `model.ModuleIndexTree`.
- `internal/source/assets.go:39` through `:51` normalizes package data responses into `model.PackageData`.
- `internal/source/fetch.go:23` through `:31` retains fetched `Body`, but `SourceResult` at `internal/source/fetch.go:113` through `:124` does not carry that raw body into the source result consumed by `RawFlow`.
- `internal/service/batch_d_test.go:271` through `:284` checks only that the raw envelope has source status and a non-nil payload; it does not fail if upstream field names are replaced by canonical model names.

Impact:

Auditors using `seekmoon raw` cannot inspect the exact upstream API shape. For example, manifest metadata and package/index structures are already projected into canonical model structs before the raw command returns them, so the command no longer satisfies the raw audit and field-recheck role.

Required action before re-review:

Expose source-specific raw fetch paths for `raw` that preserve upstream JSON object/array field names and original shape while still attaching source URL/path, fetch status, timestamp, raw ref, and parse/failure state. Add tests that assert upstream field names remain visible and canonical-only fields are not substituted for the raw payload.

## Positive Checks

The following Batch D boundaries were checked and did not produce blocking findings:

- Services compose source results and stores; output rendering remains outside service flows.
- `doctor` reports environment/path/project state and does not create snapshots or adoption records.
- `search` uses latest snapshot when available, falls back to transient Modules API fetch, writes session candidate mapping, enriches only visible manifest fields, and keeps target support unknown before evidence.
- `view` reads manifest and module-index summary without expanding package API details.
- `compare` aligns evidence fields and does not introduce a quality score or recommendation.
- `api` loads module index before package data and includes known package paths on package-path failures.
- `source` records separate attempts and does not collapse failed, unavailable, and unknown states.
- Repository fallback is represented as a repository signal, not as a registry published source.
- `skill` uses Skills API and remains separate from library adoption flow.
- `probe` uses isolated paths, records command results/logs, and does not write adoption records.
- `record` validates only the fixed adoption conclusion enum.
- `report` lists only evidence sources actually present in stored records.
- No WP13 black-box acceptance harness was added as part of Batch D.

Store listing helpers using `os.ReadDir` were reviewed as an explicit risk focus. I do not treat them as a current rejection finding because the present store implementations are already OS-backed and the Batch D failures above are sufficient to block promotion. The principal coordinator should decide whether moving listing behind platform FS is required before a later promotion boundary.

## Evidence Commands

`git status --short`

Result: Batch D implementation and coordination paths were present before this report was written.

```text
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md
 M internal/app/runtime.go
 M internal/model/output.go
 M internal/output/error.go
 M internal/output/pretty.go
 M internal/service/sync.go
 D internal/service/sync_test.go
 M internal/source/mooncli.go
 M internal/store/record_store.go
 M internal/store/snapshot_store.go
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-d-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/013-batch-d-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/014-batch-d-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/013-batch-d-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/014-batch-d-review.md
?? internal/service/api.go
?? internal/service/batch_d_test.go
?? internal/service/compare.go
?? internal/service/doctor.go
?? internal/service/helpers.go
?? internal/service/probe.go
?? internal/service/raw.go
?? internal/service/record.go
?? internal/service/report.go
?? internal/service/search.go
?? internal/service/skill.go
?? internal/service/source.go
?? internal/service/view.go
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

`go test ./internal/service -run 'Test(Doctor|Sync|Search|View|Compare)'`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/service (cached)
```

`go test ./internal/service -run 'Test(API|Source|Skill)'`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/service (cached)
```

`go test ./internal/service -run 'Test(Probe|Record|Report|Raw)'`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/service (cached)
```

`go test ./internal/cli`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/cli (cached)
```

`go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/service ./internal/cli ./internal/output ./internal/source`

Result:

```text
github.com/yumiaura/seekmoon/internal/service context encoding/json errors fmt github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/source github.com/yumiaura/seekmoon/internal/store os path/filepath sort strings time
github.com/yumiaura/seekmoon/internal/cli context fmt github.com/spf13/cobra github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/output github.com/yumiaura/seekmoon/internal/service io strconv strings
github.com/yumiaura/seekmoon/internal/output context encoding/json errors fmt github.com/itchyny/gojq github.com/yumiaura/seekmoon/internal/contract github.com/yumiaura/seekmoon/internal/model io strings
github.com/yumiaura/seekmoon/internal/source archive/zip bufio bytes context encoding/json errors fmt github.com/cenkalti/backoff/v5 github.com/google/go-github/v88/github github.com/pelletier/go-toml/v2 github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/platform github.com/yumiaura/seekmoon/internal/store golang.org/x/oauth2 io net/http net/url path/filepath regexp strings time
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

## Required Actions Before Re-Review

1. Fix `sync` so local registry index source state and summary facts are recorded in the snapshot, including partial failure state.
2. Fix `raw` so it preserves upstream payload field names and raw shape rather than returning normalized canonical model values.
3. Add tests that fail on both regressions.

After those changes, rerun the Batch D review package gates.
