# Batch E Builder Report

## Summary

Implemented Batch E / WP13 black-box acceptance and quality-gate evidence.

- Added executable A1-A16 black-box CLI acceptance tests with visible test names.
- Added library, skill, pipeline, and failure-recovery journey tests.
- Added opt-in integration boundary tests that skip by default unless explicit integration environment variables are set.
- Added a `tests/acceptance` CLI harness using real Cobra command wiring and the real output renderer with fake services, temp directories, and no external network or Moon CLI side effects.
- Updated `justfile` so format checks include `tests`, and the fuzz recipe includes `internal/source` as required by WP13.
- Installed missing pinned tools from Appendix G in `/home/t103o/go/bin`: `gotestsum v1.13.0`, `golangci-lint v2.12.2`, and `govulncheck v1.4.0`.

No commit was created.

## Files Changed

Product/test implementation:

- `justfile`
- `tests/acceptance/harness.go`
- `tests/blackbox/a_acceptance_test.go`
- `tests/journey/journey_test.go`
- `tests/integration/integration_test.go`

Coordination files present from Batch E dispatch:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-e-builder.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/017-batch-e-builder.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/017-batch-e-builder-report.md`

Pre-existing modified coordination state files were not reverted:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md`

Generated, untracked quality artifact:

- `.artifacts/coverage.out`

## Reading Evidence

Required Batch E prompt and handoff:

- `coordination/prompts/017-batch-e-builder.md`
- `coordination/handoffs/batch-e-builder.md`

Required WBS files:

- `010-mandatory-global-context.adoc`
- `020-wbs-dependency-map.adoc`
- `150-wp13-black-box-acceptance-and-quality-gates.adoc`
- `900-source-include-map.adoc`
- `coordination/resource-map.md`
- `coordination/coordination-state.md`

WP13 exact include ranges read:

- `09-seekmoon-cli-discovery-workbench/parts/060-journeys-and-acceptance/050-black-box-acceptance.adoc`, lines `4..75`
- `09-seekmoon-cli-discovery-workbench/parts/055-go-implementation-architecture/110-testing-and-tooling.adoc`, lines `4..53`
- `09-seekmoon-cli-discovery-workbench/parts/055-go-implementation-architecture/120-engineering-quality-toolchain.adoc`, lines `59..159`

Global context full-file includes read:

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

## A1-A16 Acceptance Mapping

All A1-A16 behaviors are covered in `tests/blackbox/a_acceptance_test.go` with visible names:

- A1: `TestA1SearchGeneratesCandidatesWithoutHandWrittenMooncakesURLs`
- A2: `TestA2SearchResultCanBeReferencedBySessionLocalNumber`
- A3: `TestA3LibraryModuleAndSkillEntryUseDifferentCommandSurfaces`
- A4: `TestA4ModuleProfileContainsManifestEvidenceAndPackageIndexState`
- A5: `TestA5PackageAPIProfileComesFromModuleIndexAndPackageData`
- A6: `TestA6PublishedSourceCanBeFetchedOrLocatedThroughSourceResolution`
- A7: `TestA7TargetSupportRemainsUnknownBeforeEvidence`
- A8: `TestA8ProbeProducesLocalDerivedEvidenceAndDoesNotMutateUpstreamFacts`
- A9: `TestA9AdoptionDecisionPersistsAsRecordWithEvidenceRefs`
- A10: `TestA10ReportListsOnlyActuallyUsedSources`
- A11: `TestA11PrettyTextIsLowNoiseAndNotParsingInterface`
- A12: `TestA12JSONOutputContainsSchemaIDAndEvidenceStates`
- A13: `TestA13BuiltInJQEvaluatesCommandJSONOutput`
- A14: `TestA14ShapeExplainsJSONFieldsWithoutRealQueryExecution`
- A15: `TestA15SchemaProvidesJSONSchemaForStrictConsumers`
- A16: `TestA16CommandFailureUsesErrorSurfaceWithSourceStateMeaningAndRecovery`

The black-box tests drive `cli.ExecuteWithCode` through the real Cobra command tree and real `output.DefaultRenderer`. Services are fake in-process services that return canonical model objects and record inputs. The harness uses `t.TempDir()` for runtime paths and does not use real network, real Moon CLI, GitHub API, source zip downloads, or probe mutation.

## Journey And Integration Boundary Notes

Journey tests in `tests/journey/journey_test.go`:

- `TestLibraryDiscoveryJourneySearchViewAPIProbeRecordReport`
- `TestSkillDiscoveryJourneyUsesSkillSurfaceAndRunwasmProfile`
- `TestPipelineJourneyCombinesJSONShapeSchemaAndJQ`
- `TestFailureRecoveryJourneySurfacesActionableAPIPathFailure`

Integration boundary tests in `tests/integration/integration_test.go`:

- `TestIntegrationRealNetworkIsOptIn` requires `SEEKMOON_INTEGRATION_NETWORK`
- `TestIntegrationMoonCLIProbeMutationIsOptIn` requires `SEEKMOON_INTEGRATION_MOONCLI`
- `TestIntegrationGitHubAPIIsOptIn` requires `SEEKMOON_INTEGRATION_GITHUB`

Default run evidence confirms all three integration tests skip when env vars are unset.

## Quality Gate Evidence

Tool installation/availability:

Before installation, `gotestsum`, `golangci-lint`, and `govulncheck` were not on PATH. Installed pinned tools from Appendix G with:

```text
go install gotest.tools/gotestsum@v1.13.0
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
go install golang.org/x/vuln/cmd/govulncheck@v1.4.0
```

Installed binary paths:

```text
/home/t103o/go/bin/gotestsum
/home/t103o/go/bin/golangci-lint
/home/t103o/go/bin/govulncheck
```

Version checks with `PATH="$(go env GOPATH)/bin:$PATH"`:

```text
gotestsum version v1.13.0
golangci-lint has version 2.12.2 built with go1.26.3
Go: go1.26.3
Scanner: govulncheck@v1.4.0
```

`go test ./...`

Result: passed.

```text
?   	github.com/yumiaura/seekmoon/cmd/seekmoon	[no test files]
ok  	github.com/yumiaura/seekmoon/internal/app	(cached)
ok  	github.com/yumiaura/seekmoon/internal/cli	(cached)
ok  	github.com/yumiaura/seekmoon/internal/contract	(cached)
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
ok  	github.com/yumiaura/seekmoon/internal/output	(cached)
ok  	github.com/yumiaura/seekmoon/internal/platform	(cached)
ok  	github.com/yumiaura/seekmoon/internal/service	(cached)
ok  	github.com/yumiaura/seekmoon/internal/source	(cached)
ok  	github.com/yumiaura/seekmoon/internal/store	(cached)
?   	github.com/yumiaura/seekmoon/internal/testutil	[no test files]
?   	github.com/yumiaura/seekmoon/tests/acceptance	[no test files]
ok  	github.com/yumiaura/seekmoon/tests/blackbox	0.010s
ok  	github.com/yumiaura/seekmoon/tests/integration	(cached)
ok  	github.com/yumiaura/seekmoon/tests/journey	(cached)
```

`go test ./tests/blackbox ./tests/journey`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/tests/blackbox	(cached)
ok  	github.com/yumiaura/seekmoon/tests/journey	(cached)
```

Default integration skip evidence:

`go test ./tests/integration`

Result: passed with skips; `just test` and `just test-race` also reported the three explicit skips.

```text
ok  	github.com/yumiaura/seekmoon/tests/integration	(cached)
```

`just fmt-check`

Result: passed.

```text
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal tests)"
```

`PATH="$(go env GOPATH)/bin:$PATH" just lint`

Result: failed. The pinned lint tool now runs, but the repository has broad lint findings outside the Batch E write boundary.

Representative findings:

```text
internal/source/fetch.go:92:23: Error return value of `resp.Body.Close` is not checked (errcheck)
internal/platform/exec.go:36:9: G204: Subprocess launched with a potential tainted input or cmd arguments (gosec)
internal/platform/fs.go:29:9: G304: Potential file inclusion via variable (gosec)
internal/service/batch_d_test.go:384:18: G705: XSS via taint analysis (gosec)
internal/model/api.go:8:6: exported type ModuleIndexTree should have comment or be unexported (revive)
internal/service/doctor.go:52:6: func environmentPathStatus is unused (unused)
```

Final lint summary from the first run:

```text
69 issues:
* errcheck: 1
* gosec: 8
* prealloc: 3
* revive: 50
* staticcheck: 3
* unparam: 2
* unused: 2
```

After removing one new-test `unparam` issue, `just ci` still failed at lint with:

```text
68 issues:
* errcheck: 1
* gosec: 8
* prealloc: 3
* revive: 50
* staticcheck: 3
* unparam: 1
* unused: 2
```

`PATH="$(go env GOPATH)/bin:$PATH" just test`

Result: passed.

```text
DONE 113 tests, 3 skipped in 0.619s
```

`PATH="$(go env GOPATH)/bin:$PATH" just test-race`

Result: passed.

```text
DONE 113 tests, 3 skipped in 1.761s
```

`PATH="$(go env GOPATH)/bin:$PATH" just cover`

Result: passed. Coverage summary ended with:

```text
total:									(statements)			54.3%
```

`PATH="$(go env GOPATH)/bin:$PATH" just vuln`

Result: failed due to Go standard library vulnerabilities in local `go1.26.3`; fixed in `go1.26.4`.

```text
Vulnerability #1: GO-2026-5039
    Arbitrary inputs are included in errors without any escaping in
    net/textproto
  Standard library
    Found in: net/textproto@go1.26.3
    Fixed in: net/textproto@go1.26.4

Vulnerability #2: GO-2026-5037
    Inefficient candidate hostname parsing in crypto/x509
  Standard library
    Found in: crypto/x509@go1.26.3
    Fixed in: crypto/x509@go1.26.4

Your code is affected by 2 vulnerabilities from the Go standard library.
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

Result: failed at `lint`, so later composed gates were not reached in the composed command.

```text
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal tests)"
command -v golangci-lint >/dev/null
golangci-lint run ./...
...
error: Recipe `lint` failed on line 13 with exit code 1
```

## Known Gaps Or Blocked Gates

- `just lint` is blocked by existing repository-wide lint findings across `internal/model`, `internal/output`, `internal/platform`, `internal/source`, `internal/service`, and `internal/testutil`. Resolving those findings would require broad edits outside the Batch E write boundary.
- `just vuln` is blocked by local Go toolchain `go1.26.3` standard-library vulnerabilities fixed in `go1.26.4`.
- `just ci` is blocked because it composes `lint` first; it stops at the lint failure. Even after lint cleanup, `vuln` would still fail until the Go toolchain is upgraded to a fixed `go1.26.4` build.
- `.artifacts/coverage.out` was generated by `just cover` and is not a runtime source artifact.

## Suggested Reviewer Focus

- Confirm all A1-A16 tests are present as executable test names in `tests/blackbox/a_acceptance_test.go`.
- Confirm default acceptance and journey tests are offline and use fake services/temp directories rather than real network or Moon CLI.
- Confirm integration tests are visibly opt-in and skip by default.
- Confirm `justfile` includes `tests` in format checks and includes `internal/source` in the fuzz recipe.
- Confirm the blocked gates are not caused by Batch E acceptance tests: `go test`, `tests/blackbox`, `tests/journey`, `fmt-check`, `test`, `test-race`, `cover`, `mod-check`, and `release-check` pass; `lint`, `vuln`, and therefore `ci` are blocked by the findings listed above.
