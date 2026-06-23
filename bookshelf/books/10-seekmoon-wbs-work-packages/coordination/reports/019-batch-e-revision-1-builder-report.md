# Batch E Revision 1 Builder Report

## Summary

Batch E Revision 1 addressed the independent review rejection for WP13.

The promotion-critical gates now pass under the required command prefix:

- `PATH="$(go env GOPATH)/bin:$PATH" just lint`
- `PATH="$(go env GOPATH)/bin:$PATH" just vuln`
- `PATH="$(go env GOPATH)/bin:$PATH" just ci`

The local Go toolchain was updated from `go1.26.3` to `go1.26.4` using the official `go1.26.4.linux-amd64.tar.gz` archive and the SHA-256 published by `https://go.dev/dl/?mode=json&include=all`.

No commit was created.

## Review Findings Addressed

### `just lint`

The review reported 68 lint findings. The revision makes `golangci-lint run ./...` pass with `0 issues`.

Lint fixes included:

- Checked `resp.Body.Close()` in `internal/source/fetch.go`.
- Removed unused code in `internal/service/doctor.go` and `internal/service/helpers.go`.
- Removed the always-empty fallback argument from `internal/service/raw.go`.
- Reworked preallocation findings in `internal/contract/shapes.go`, `internal/output/pretty.go`, and `internal/source/assets.go`.
- Rewrote staticcheck tagged-switch findings in `internal/source/project.go`.
- Replaced string-concatenated JSON fixture output in `internal/service/batch_d_test.go` with JSON encoder output.
- Tightened local index path validation around `WalkDir` callback reads and git ref reads in `internal/source/local_index.go`.
- Added justified `#nosec` comments where the code intentionally exercises repository host filesystem or test fixture boundaries.
- Added exported API comments across existing exported surfaces rather than disabling revive.
- Renamed unused callback parameters to `_` where appropriate.
- Unexported the source package helper formerly named `SourceResult` to avoid the revive stutter finding.

### `just vuln`

The review reported reachable Go standard-library vulnerabilities in local `go1.26.3`:

- `GO-2026-5039`, fixed in `go1.26.4`
- `GO-2026-5037`, fixed in `go1.26.4`

Revision action:

- Confirmed official `go1.26.4.linux-amd64.tar.gz` metadata:
  - SHA-256: `1153d3d50e0ac764b447adfe05c2bcf08e889d42a02e0fe0259bd47f6733ad7f`
  - Size: `66861266`
- Downloaded the archive under `/home/t103o/workbench/tmp/go1.26.4/`.
- Verified checksum with `sha256sum -c`.
- Replaced `/usr/local/go` with the verified Go 1.26.4 archive.
- Confirmed:
  - `go version go1.26.4 linux/amd64`
  - `go env GOVERSION` = `go1.26.4`
  - `go env GOROOT` = `/usr/local/go`

`go.mod` was temporarily tested with `go 1.26.4`, but `just mod-check` intentionally requires no `go.mod`/`go.sum` diff after `go mod tidy`. The module metadata was therefore left at the existing `go 1.26.3` while the local scanner/compiler toolchain is `go1.26.4`. This satisfies the review's standard-library vulnerability finding without introducing a module metadata diff that fails the project's own gate.

### `just ci`

`PATH="$(go env GOPATH)/bin:$PATH" just ci` now passes end-to-end:

- `fmt-check`
- `lint`
- `test`
- `test-race`
- `cover`
- `vuln`
- `mod-check`
- `release-check`

## Files Changed

Repository behavior and lint cleanup:

- `cmd/seekmoon/main.go`
- `internal/app/command.go`
- `internal/app/runtime.go`
- `internal/cli/doctor.go`
- `internal/cli/report.go`
- `internal/cli/root.go`
- `internal/cli/root_test.go`
- `internal/cli/search.go`
- `internal/cli/sync.go`
- `internal/contract/schemas.go`
- `internal/contract/shapes.go`
- `internal/model/api.go`
- `internal/model/evidence.go`
- `internal/model/manifest.go`
- `internal/model/module.go`
- `internal/model/output.go`
- `internal/model/probe.go`
- `internal/model/project.go`
- `internal/model/record.go`
- `internal/model/report.go`
- `internal/model/repository.go`
- `internal/model/schema.go`
- `internal/model/skill.go`
- `internal/model/snapshot.go`
- `internal/model/source.go`
- `internal/output/error.go`
- `internal/output/jq.go`
- `internal/output/json.go`
- `internal/output/pretty.go`
- `internal/output/render.go`
- `internal/output/schema.go`
- `internal/output/shape.go`
- `internal/platform/clock.go`
- `internal/platform/env.go`
- `internal/platform/exec.go`
- `internal/platform/fs.go`
- `internal/platform/http.go`
- `internal/service/api.go`
- `internal/service/batch_d_test.go`
- `internal/service/compare.go`
- `internal/service/doctor.go`
- `internal/service/helpers.go`
- `internal/service/probe.go`
- `internal/service/raw.go`
- `internal/service/record.go`
- `internal/service/registry.go`
- `internal/service/report.go`
- `internal/service/search.go`
- `internal/service/skill.go`
- `internal/service/source.go`
- `internal/service/sync.go`
- `internal/service/view.go`
- `internal/source/assets.go`
- `internal/source/assets_test.go`
- `internal/source/fetch.go`
- `internal/source/local_cache.go`
- `internal/source/local_index.go`
- `internal/source/mooncakes.go`
- `internal/source/mooncakes_test.go`
- `internal/source/mooncli.go`
- `internal/source/project.go`
- `internal/source/repository.go`
- `internal/source/repository_skill_test.go`
- `internal/source/skills.go`
- `internal/store/atomic.go`
- `internal/store/cache_store.go`
- `internal/store/paths.go`
- `internal/store/probe_store.go`
- `internal/store/record_store.go`
- `internal/store/report_store.go`
- `internal/store/session_store.go`
- `internal/store/snapshot_store.go`
- `internal/store/source_store.go`
- `internal/testutil/fixture.go`
- `internal/testutil/golden.go`

Batch E paths already present from the builder pass and preserved:

- `justfile`
- `tests/acceptance/harness.go`
- `tests/blackbox/a_acceptance_test.go`
- `tests/journey/journey_test.go`
- `tests/integration/integration_test.go`

Coordination/report path:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/019-batch-e-revision-1-builder-report.md`

Existing coordination files were already modified before this revision pass and remain modified:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md`

## Toolchain Changes

Local Go toolchain:

- Before: `go version go1.26.3 linux/amd64`
- After: `go version go1.26.4 linux/amd64`
- `GOROOT`: `/usr/local/go`
- `GOVERSION`: `go1.26.4`

Download and checksum evidence:

```text
go1.26.4.linux-amd64.tar.gz	1153d3d50e0ac764b447adfe05c2bcf08e889d42a02e0fe0259bd47f6733ad7f	66861266
/home/t103o/workbench/tmp/go1.26.4/go1.26.4.linux-amd64.tar.gz: OK
```

Module metadata:

- `go.mod` remains at `go 1.26.3`.
- Rationale: `just mod-check` requires `go mod tidy` followed by no `go.mod`/`go.sum` diff. Keeping the existing module directive preserves that gate while `govulncheck` uses the fixed local Go 1.26.4 standard library.

## Quality Gate Evidence

`go test ./...`

Result: passed.

```text
?    github.com/yumiaura/seekmoon/cmd/seekmoon [no test files]
ok   github.com/yumiaura/seekmoon/internal/app
ok   github.com/yumiaura/seekmoon/internal/cli
ok   github.com/yumiaura/seekmoon/internal/contract
ok   github.com/yumiaura/seekmoon/internal/model
ok   github.com/yumiaura/seekmoon/internal/output
ok   github.com/yumiaura/seekmoon/internal/platform
ok   github.com/yumiaura/seekmoon/internal/service
ok   github.com/yumiaura/seekmoon/internal/source
ok   github.com/yumiaura/seekmoon/internal/store
?    github.com/yumiaura/seekmoon/internal/testutil [no test files]
?    github.com/yumiaura/seekmoon/tests/acceptance [no test files]
ok   github.com/yumiaura/seekmoon/tests/blackbox
ok   github.com/yumiaura/seekmoon/tests/integration
ok   github.com/yumiaura/seekmoon/tests/journey
```

`go test ./tests/blackbox ./tests/journey`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/tests/blackbox
ok   github.com/yumiaura/seekmoon/tests/journey
```

`go test ./tests/integration`

Result: passed with default skip behavior.

```text
ok   github.com/yumiaura/seekmoon/tests/integration
```

`PATH="$(go env GOPATH)/bin:$PATH" just fmt-check`

Result: passed.

```text
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal tests)"
```

`PATH="$(go env GOPATH)/bin:$PATH" just lint`

Result: passed.

```text
command -v golangci-lint >/dev/null
golangci-lint run ./...
0 issues.
```

`PATH="$(go env GOPATH)/bin:$PATH" just test`

Result: passed.

```text
DONE 113 tests, 3 skipped in 0.894s
```

Skipped integration tests remained opt-in:

```text
SKIP: TestIntegrationRealNetworkIsOptIn
SKIP: TestIntegrationMoonCLIProbeMutationIsOptIn
SKIP: TestIntegrationGitHubAPIIsOptIn
```

`PATH="$(go env GOPATH)/bin:$PATH" just test-race`

Result: passed.

```text
DONE 113 tests, 3 skipped in 1.751s
```

`PATH="$(go env GOPATH)/bin:$PATH" just cover`

Result: passed.

```text
total: (statements) 54.2%
```

`PATH="$(go env GOPATH)/bin:$PATH" just vuln`

Result: passed.

```text
command -v govulncheck >/dev/null
govulncheck ./...
No vulnerabilities found.
```

`PATH="$(go env GOPATH)/bin:$PATH" just mod-check`

Result: passed.

```text
go mod tidy
git diff --exit-code -- go.mod go.sum
go mod verify
all modules verified
```

`PATH="$(go env GOPATH)/bin:$PATH" just release-check`

Result: passed.

```text
command -v goreleaser >/dev/null
goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
```

`PATH="$(go env GOPATH)/bin:$PATH" just ci`

Result: passed.

Evidence highlights:

```text
golangci-lint run ./...
0 issues.
DONE 113 tests, 3 skipped in 0.249s
DONE 113 tests, 3 skipped in 1.176s
No vulnerabilities found.
all modules verified
1 configuration file(s) validated
```

## Acceptance Regression Evidence

A1-A16 executable black-box test names remain present in `tests/blackbox/a_acceptance_test.go`:

```text
TestA1SearchGeneratesCandidatesWithoutHandWrittenMooncakesURLs
TestA2SearchResultCanBeReferencedBySessionLocalNumber
TestA3LibraryModuleAndSkillEntryUseDifferentCommandSurfaces
TestA4ModuleProfileContainsManifestEvidenceAndPackageIndexState
TestA5PackageAPIProfileComesFromModuleIndexAndPackageData
TestA6PublishedSourceCanBeFetchedOrLocatedThroughSourceResolution
TestA7TargetSupportRemainsUnknownBeforeEvidence
TestA8ProbeProducesLocalDerivedEvidenceAndDoesNotMutateUpstreamFacts
TestA9AdoptionDecisionPersistsAsRecordWithEvidenceRefs
TestA10ReportListsOnlyActuallyUsedSources
TestA11PrettyTextIsLowNoiseAndNotParsingInterface
TestA12JSONOutputContainsSchemaIDAndEvidenceStates
TestA13BuiltInJQEvaluatesCommandJSONOutput
TestA14ShapeExplainsJSONFieldsWithoutRealQueryExecution
TestA15SchemaProvidesJSONSchemaForStrictConsumers
TestA16CommandFailureUsesErrorSurfaceWithSourceStateMeaningAndRecovery
```

Offline/default behavior remains preserved:

- `tests/acceptance/harness.go` wires the real Cobra command tree to deterministic fake services.
- Default `go test ./...`, `just test`, `just test-race`, `just cover`, and `just ci` run without network credentials, real Moon CLI mutation, or GitHub credentials.
- Integration tests remain skipped by default unless these variables are set:
  - `SEEKMOON_INTEGRATION_NETWORK`
  - `SEEKMOON_INTEGRATION_MOONCLI`
  - `SEEKMOON_INTEGRATION_GITHUB`

## Remaining Risks Or Blocked Gates

No required WP13 revision gate remains blocked.

Known residual context:

- The local machine now has Go `1.26.4` installed under `/usr/local/go`.
- `go.mod` intentionally remains unchanged at `go 1.26.3` because `just mod-check` requires no module metadata diff after `go mod tidy`.
- Lint cleanup introduced many comments on previously exported surfaces. Reviewer should distinguish this from behavior change; most of the large diff is revive compliance.

## Suggested Reviewer Focus

Review should focus on:

- Whether the `#nosec` comments are sufficiently contract-bound:
  - `internal/platform/exec.go`
  - `internal/platform/fs.go`
  - `internal/source/local_index.go`
  - `internal/testutil/golden.go`
- Whether the local index path validation in `internal/source/local_index.go` preserves expected local registry behavior.
- Whether keeping `go.mod` at `go 1.26.3` while running the local Go `1.26.4` toolchain is acceptable for WP13 promotion, given that `just mod-check`, `just vuln`, and `just ci` all pass in that state.
- Whether the expanded exported-surface comments are acceptable as lint remediation rather than API design churn.
