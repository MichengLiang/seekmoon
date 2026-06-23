# Batch D Revision 1 Builder Report

## Summary

Addressed the two independent review rejection findings from `014-batch-d-review-report.md`.

- `sync` now reads the local registry index from the resolved Moon home registry path, appends its state-bearing source result to `snapshot.Sources`, and records status/error plus file, record, malformed, and HEAD summary facts in `snapshot.Raw["local_index"]`.
- `raw` now uses raw JSON fetch paths for Modules API, Manifest API, module index assets, package data assets, and Skills API. The raw envelope keeps source metadata and carries the upstream JSON object/array shape instead of normalized model values.

No commit was created.

## Review Findings Addressed

### Finding 1: `sync` omitted local registry index source and summary

Review finding:

- WP10 requires `sync` to execute `moon update`, read Modules API, read Statistics API, count the local index, and save a snapshot.
- The local registry index entry must preserve partial source failure state and include summary facts such as index HEAD, file count, and record count.

Revision actions:

- Added `store.Paths.MoonIndex`, resolved as `<HOME>/.moon/registry/index/user`, matching the WBS/source basis for `~/.moon/registry/index/user/**/*.index`.
- Updated `SyncService.Sync` to call `s.LocalIndex.Read(ctx, s.Paths.MoonIndex)`.
- Appended the erased local-index `SourceResult` to `snapshot.Sources`.
- Added `localIndexSummary` to store `status`, `parse_state`, `path`, `raw_ref`, `error`, `index_head`, `file_count`, `record_count`, and `malformed` under `snapshot.Raw["local_index"]`.
- Extended `source.LocalIndexSummary` with `FileCount`, `RecordCount`, and `IndexHead`.
- Extended `LocalIndexReader.Read` to handle directory roots, recurse over `.index` files, count files and JSONL records, and return `unavailable` when the configured index path is absent.
- Kept partial source failures in source-result state instead of converting absence into sync success.

Focused tests:

- `TestSyncRecordsPartialSourceFailure`
  - Confirms Modules API present, Statistics API failed, and local index unavailable are all present in snapshot source state.
  - Confirms `snapshot.Raw["local_index"]` carries unavailable status and error.
- `TestSyncRecordsLocalIndexSummaryWhenPresent`
  - Creates two local `.index` files under a temp Moon home.
  - Confirms local-index source status is present.
  - Confirms file count `2`, record count `3`, and malformed count `0`.

### Finding 2: `raw` returned normalized model payloads

Review finding:

- WP12 requires `raw` to expose original payloads for source audit, field recheck, and failure reproduction.
- `raw` must preserve upstream field names and original object/array shape, not normalized model payloads.

Revision actions:

- Added `source.RawJSONSourceResult`, which fetches bytes, decodes JSON into `any` with `UseNumber`, preserves source metadata, and records parse/failure state.
- Added raw source methods:
  - `MooncakesClient.FetchRawModules`
  - `MooncakesClient.FetchRawManifest`
  - `AssetClient.FetchRawModuleIndex`
  - `AssetClient.FetchRawPackageData`
  - `SkillsClient.FetchRawSkills`
- Updated `RawFlow.Raw` to call only raw fetch methods and unwrap the pointer value into `model.RawEnvelope.Payload`.
- Kept ordinary normalized source readers unchanged for `search`, `view`, `api`, and `skill`.

Focused test:

- `TestRawReturnsSourceStatusAndUpstreamPayloadShape`
  - Confirms Modules API raw payload keeps upstream `name` and does not substitute canonical `module`.
  - Confirms Manifest API raw payload keeps `latest_version` and does not add derived `docs_url` or `versions_count`.
  - Confirms module index raw payload keeps upstream `children` and is not converted to normalized `childs`.
  - Confirms package data raw payload keeps HTML `signature` and does not add canonical `plain_signature`.

## Files Changed

Revision-specific implementation paths:

- `internal/store/paths.go`
- `internal/source/local_index.go`
- `internal/source/fetch.go`
- `internal/source/mooncakes.go`
- `internal/source/assets.go`
- `internal/source/skills.go`
- `internal/service/sync.go`
- `internal/service/raw.go`
- `internal/service/doctor.go`
- `internal/service/batch_d_test.go`

Pre-existing Batch D paths still present in the worktree:

- `internal/app/runtime.go`
- `internal/model/output.go`
- `internal/output/error.go`
- `internal/output/pretty.go`
- `internal/source/mooncli.go`
- `internal/store/record_store.go`
- `internal/store/snapshot_store.go`
- `internal/service/api.go`
- `internal/service/compare.go`
- `internal/service/helpers.go`
- `internal/service/probe.go`
- `internal/service/record.go`
- `internal/service/report.go`
- `internal/service/search.go`
- `internal/service/skill.go`
- `internal/service/source.go`
- `internal/service/view.go`
- `internal/service/sync_test.go` deleted by the original Batch D builder because coverage moved into `batch_d_test.go`.

Coordination paths present from Batch D/revision flow:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-d-builder.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-d-revision-1-builder.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/013-batch-d-builder.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/014-batch-d-reviewer.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/015-batch-d-revision-1-builder.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/013-batch-d-builder-report.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/014-batch-d-review-report.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/015-batch-d-revision-1-builder-report.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/014-batch-d-review.md`

Pre-existing modified coordination state files were not reverted:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md`

## Evidence Commands

Required reading for this revision:

- `coordination/prompts/015-batch-d-revision-1-builder.md`
- `coordination/handoffs/batch-d-revision-1-builder.md`
- `coordination/reports/014-batch-d-review-report.md`
- `coordination/handoffs/batch-d-builder.md`
- `coordination/review-packages/014-batch-d-review.md`

Focused WBS/source ranges rechecked for the rejected findings:

- `120-wp10-discovery-and-profile-services.adoc`
- `140-wp12-assessment-record-and-report-services.adoc`
- `900-source-include-map.adoc`
- `09-seekmoon-cli-discovery-workbench/parts/040-command-workbench/020-doctor-and-sync.adoc`, lines `47..83`
- `09-seekmoon-cli-discovery-workbench/parts/055-go-implementation-architecture/050-data-flow.adoc`, lines `8..22`
- `09-seekmoon-cli-discovery-workbench/parts/040-command-workbench/070-record-report-and-raw.adoc`, lines `41..52`
- `09-seekmoon-cli-discovery-workbench/parts/020-evidence-boundary/030-local-toolchain-and-cache.adoc`, lines `40..56`

`go test ./...`

Result: passed.

```text
?   	github.com/yumiaura/seekmoon/cmd/seekmoon	[no test files]
ok  	github.com/yumiaura/seekmoon/internal/app	0.003s
ok  	github.com/yumiaura/seekmoon/internal/cli	0.005s
ok  	github.com/yumiaura/seekmoon/internal/contract	(cached)
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
ok  	github.com/yumiaura/seekmoon/internal/output	(cached)
ok  	github.com/yumiaura/seekmoon/internal/platform	(cached)
ok  	github.com/yumiaura/seekmoon/internal/service	0.237s
ok  	github.com/yumiaura/seekmoon/internal/source	0.253s
ok  	github.com/yumiaura/seekmoon/internal/store	0.005s
?   	github.com/yumiaura/seekmoon/internal/testutil	[no test files]
```

`go test ./internal/service -run 'Test(Doctor|Sync|Search|View|Compare)'`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/service	0.087s
```

`go test ./internal/service -run 'Test(API|Source|Skill)'`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/service	0.062s
```

`go test ./internal/service -run 'Test(Probe|Record|Report|Raw)'`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/service	0.052s
```

Focused tests added for sync local index and raw upstream payload shape:

`go test ./internal/service -run 'TestSyncRecordsLocalIndexSummaryWhenPresent|TestSyncRecordsPartialSourceFailure|TestRawReturnsSourceStatusAndUpstreamPayloadShape'`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/service	0.088s
```

Additional focused source-reader check run during implementation:

`go test ./internal/source -run 'TestLocalIndex|TestAsset|TestMooncakes|TestSkills'`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/source	0.128s
```

`go test ./internal/cli`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/cli	(cached)
```

`go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/service ./internal/cli ./internal/output ./internal/source`

Result:

```text
github.com/yumiaura/seekmoon/internal/service context encoding/json errors fmt github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/source github.com/yumiaura/seekmoon/internal/store os path/filepath sort strings time
github.com/yumiaura/seekmoon/internal/cli context fmt github.com/spf13/cobra github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/output github.com/yumiaura/seekmoon/internal/service io strconv strings
github.com/yumiaura/seekmoon/internal/output context encoding/json errors fmt github.com/itchyny/gojq github.com/yumiaura/seekmoon/internal/contract github.com/yumiaura/seekmoon/internal/model io strings
github.com/yumiaura/seekmoon/internal/source archive/zip bufio bytes context encoding/json errors fmt github.com/cenkalti/backoff/v5 github.com/google/go-github/v88/github github.com/pelletier/go-toml/v2 github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/platform github.com/yumiaura/seekmoon/internal/store golang.org/x/oauth2 io io/fs net/http net/url os path/filepath regexp strings time
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

## Remaining Risks Or Blocked Gates

- No blocked gates. All revision handoff verification commands passed.
- `LocalIndexReader` directory traversal is currently OS-backed for recursive directory reads. This matches the existing OS-backed store/runtime pattern noted in the Batch D review, but a later cleanup could move recursive file listing behind `platform.FS` if the principal coordinator chooses to generalize filesystem abstraction.
- The local index HEAD summary is best-effort. It is present when a `.git/HEAD` can be resolved from the configured index path or its parent directories, and empty otherwise. The status/error and file/record counts remain recorded independently.

## Suggested Reviewer Focus

- Confirm `sync` snapshots include three source results for Modules API, Statistics API, and Local registry index, including failed/unavailable local index state.
- Confirm `snapshot.Raw["local_index"]` includes the required summary facts and does not treat omission as success.
- Confirm `raw` uses `FetchRaw*` paths and does not call normalized source readers.
- Confirm raw tests would fail if upstream field names were replaced with canonical-only fields.
- Confirm the added `store.Paths.MoonIndex` path matches the local registry source basis `~/.moon/registry/index/user/**/*.index`.
