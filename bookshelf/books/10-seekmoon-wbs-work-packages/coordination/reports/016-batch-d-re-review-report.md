# Batch D Re-Review Report

## Verdict

Approved.

Batch D revision 1 fixes the two prior rejection findings from `014-batch-d-review-report.md`. `sync` now records local registry index source state and summary facts, and `raw` now returns upstream raw JSON shapes rather than normalized model values.

## Review Inputs

- `coordination/prompts/016-batch-d-re-reviewer.md`
- `coordination/review-packages/016-batch-d-re-review.md`
- `coordination/review-packages/014-batch-d-review.md`
- `coordination/reports/014-batch-d-review-report.md`
- `coordination/reports/015-batch-d-revision-1-builder-report.md`
- Prior Batch D WBS basis and exact include ranges listed in `coordination/reports/014-batch-d-review-report.md`
- Principal coordinator protocol `SKILL.md`, `references/index.md`, and `references/30_coordination_and_runtime/review_evidence_promotion.md`

## Re-Review Findings

No blocking findings.

## Prior Rejection Findings

### Sync Local Registry Index

Fixed.

Evidence:

- `internal/service/sync.go:33` through `:35` fetches Modules API, Statistics API, and local registry index.
- `internal/service/sync.go:39` through `:43` appends the local-index source result to `snapshot.Sources`.
- `internal/service/sync.go:44` through `:48` writes `snapshot.Raw["local_index"]`.
- `internal/service/sync.go:63` through `:80` records local-index status, parse state, path, raw ref, error, index HEAD, file count, record count, and malformed count when available.
- `internal/source/local_index.go:54` through `:87` preserves failed, unavailable, and present local-index file states.
- `internal/source/local_index.go:89` through `:117` reads directory roots recursively and summarizes `.index` files.
- `internal/store/paths.go:38` through `:58` resolves `MoonIndex` under `HOME/.moon/registry/index/user`.
- `internal/service/batch_d_test.go:45` through `:89` covers partial source failure with local index unavailable.
- `internal/service/batch_d_test.go:91` through `:120` covers present local index file count, record count, and malformed count.

### Raw Upstream Payload Shape

Fixed.

Evidence:

- `internal/service/raw.go:16` through `:36` calls only raw fetch methods for `modules`, `manifest`, `module-index`, `package-data`, and `skills`.
- `internal/source/fetch.go:127` through `:142` decodes the fetched body into `any` with JSON numbers preserved and keeps source metadata on the source result.
- `internal/source/mooncakes.go:65` through `:79` adds raw Modules API and Manifest API fetch paths.
- `internal/source/assets.go:54` through `:68` adds raw module-index and package-data asset fetch paths.
- `internal/source/skills.go:29` through `:32` adds raw Skills API fetch path.
- `internal/service/batch_d_test.go:308` through `:362` asserts upstream field names and shapes remain visible: `name` is not replaced by canonical `module`, `latest_version` is not replaced by derived manifest fields, `children` is not normalized into `childs`, and `signature` is not replaced by `plain_signature`.

## Positive Boundary Checks

Re-checked and passed:

- Services compose source results and stores; output rendering remains outside service flows.
- `doctor` reports environment/path/project state and does not create snapshots or adoption records.
- `sync` preserves partial source failures and does not do ranking or adoption.
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
- `raw` preserves upstream payload names and source status without becoming ordinary discovery flow.
- No WP13 black-box acceptance harness was added as part of Batch D.

The remaining OS-backed recursive local-index directory walk matches the current OS-backed store/runtime pattern and is not a promotion blocker for Batch D.

## Evidence Commands

`git status --short`

Result: only Batch D implementation and coordination paths were present before this report was written.

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
 M internal/source/assets.go
 M internal/source/fetch.go
 M internal/source/local_index.go
 M internal/source/mooncakes.go
 M internal/source/mooncli.go
 M internal/source/skills.go
 M internal/store/paths.go
 M internal/store/record_store.go
 M internal/store/snapshot_store.go
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-d-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-d-revision-1-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/013-batch-d-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/014-batch-d-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/015-batch-d-revision-1-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/016-batch-d-re-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/013-batch-d-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/014-batch-d-review-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/015-batch-d-revision-1-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/014-batch-d-review.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/016-batch-d-re-review.md
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

`go test ./internal/service -run 'TestSyncRecordsLocalIndexSummaryWhenPresent|TestSyncRecordsPartialSourceFailure|TestRawReturnsSourceStatusAndUpstreamPayloadShape'`

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

## Commit Scope

Approved commit paths:

- `internal/app/runtime.go`
- `internal/model/output.go`
- `internal/output/error.go`
- `internal/output/pretty.go`
- `internal/source/mooncli.go`
- `internal/source/assets.go`
- `internal/source/fetch.go`
- `internal/source/local_index.go`
- `internal/source/mooncakes.go`
- `internal/source/skills.go`
- `internal/store/paths.go`
- `internal/store/record_store.go`
- `internal/store/snapshot_store.go`
- `internal/service`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination`
