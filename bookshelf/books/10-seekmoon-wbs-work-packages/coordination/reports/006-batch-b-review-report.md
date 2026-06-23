# Batch B Independent Review Report

## Verdict

Rejected.

The Batch B implementation has good coverage and the requested command gates pass, but it misses a core source-reader contract: several source result producers do not carry the complete source envelope required by the WP04 shared source layer and the 9号书 source-reader boundary. Because Batch B's purpose is to establish state-bearing source readers for downstream batches, this is blocking.

No commit was created.

## Review Basis Read

### Required WBS And Coordination Files

- `010-mandatory-global-context.adoc`
- `020-wbs-dependency-map.adoc`
- `060-wp04-public-http-sources-and-snapshot.adoc`
- `070-wp05-static-assets-and-api-profile.adoc`
- `080-wp06-local-toolchain-source-and-probe-primitives.adoc`
- `090-wp07-repository-and-skill-sources.adoc`
- `900-source-include-map.adoc`
- `coordination/resource-map.md`
- `coordination/handoffs/batch-b-builder.md`
- `coordination/reports/005-batch-b-builder-report.md`

### Batch B Exact Line Includes Read

WP04:

- `parts/020-evidence-boundary/010-public-http-sources.adoc`, lines `4..5`
- `parts/020-evidence-boundary/010-public-http-sources.adoc`, lines `9..46`
- `parts/020-evidence-boundary/010-public-http-sources.adoc`, lines `50..75`
- `parts/020-evidence-boundary/010-public-http-sources.adoc`, lines `79..124`
- `parts/055-go-implementation-architecture/070-source-reader-boundary.adoc`, lines `4..42`
- `parts/055-go-implementation-architecture/050-data-flow.adoc`, lines `8..22`
- `parts/055-go-implementation-architecture/050-data-flow.adoc`, lines `26..39`

WP05:

- `parts/020-evidence-boundary/020-static-assets.adoc`, lines `4..5`
- `parts/020-evidence-boundary/020-static-assets.adoc`, lines `9..42`
- `parts/020-evidence-boundary/020-static-assets.adoc`, lines `46..83`
- `parts/020-evidence-boundary/020-static-assets.adoc`, lines `87..95`
- `parts/020-evidence-boundary/020-static-assets.adoc`, lines `99..107`
- `parts/020-evidence-boundary/020-static-assets.adoc`, lines `111..127`
- `parts/030-canonical-model/040-api-profile.adoc`, lines `4..5`
- `parts/030-canonical-model/040-api-profile.adoc`, lines `9..60`
- `parts/030-canonical-model/040-api-profile.adoc`, lines `65..133`
- `parts/040-command-workbench/040-api-and-source.adoc`, lines `4..32`

WP06:

- `parts/020-evidence-boundary/030-local-toolchain-and-cache.adoc`, lines `4..5`
- `parts/020-evidence-boundary/030-local-toolchain-and-cache.adoc`, lines `9..43`
- `parts/020-evidence-boundary/030-local-toolchain-and-cache.adoc`, lines `47..56`
- `parts/020-evidence-boundary/030-local-toolchain-and-cache.adoc`, lines `59..68`
- `parts/020-evidence-boundary/040-project-context.adoc`, lines `4..43`
- `parts/030-canonical-model/060-source-resolution-and-probe-result.adoc`, lines `77..132`

WP07:

- `parts/020-evidence-boundary/050-repository-signals.adoc`, lines `4..46`
- `parts/020-evidence-boundary/010-public-http-sources.adoc`, lines `128..169`
- `parts/030-canonical-model/050-skill-profile.adoc`, lines `4..94`
- `backmatter/appendix-f-go-implementation-dependencies.adoc`, lines `29..37`
- `backmatter/appendix-f-go-implementation-dependencies.adoc`, lines `105..113`

Additional review inputs:

- `coordination/prompts/006-batch-b-reviewer.md`
- `coordination/review-packages/006-batch-b-review.md`
- principal coordinator protocol `SKILL.md`, `references/index.md`, and `references/30_coordination_and_runtime/review_evidence_promotion.md`

## Evidence Commands

`git status --short`

Result: Batch B implementation and coordination changes are present; no unrelated path outside the review object was observed.

```text
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md
 M go.mod
 M go.sum
 M internal/app/runtime.go
 M internal/app/runtime_test.go
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-b-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/005-batch-b-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/006-batch-b-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/005-batch-b-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/006-batch-b-review.md
?? internal/model/repository.go
?? internal/service/
?? internal/source/
?? internal/testutil/
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
ok  	github.com/yumiaura/seekmoon/internal/service	(cached)
ok  	github.com/yumiaura/seekmoon/internal/source	(cached)
ok  	github.com/yumiaura/seekmoon/internal/store	(cached)
?   	github.com/yumiaura/seekmoon/internal/testutil	[no test files]
```

`go test ./internal/source ./internal/service -run 'Test(Mooncakes|Sync)'`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/source	(cached)
ok  	github.com/yumiaura/seekmoon/internal/service	(cached)
```

`go test ./internal/source ./internal/model -run 'Test(Asset|ModuleIndex|PackageData|Relpath|SourceZip)'`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/source	(cached)
ok  	github.com/yumiaura/seekmoon/internal/model	(cached) [no tests to run]
```

`go test ./internal/source ./internal/store -run 'Test(MoonCLI|LocalIndex|LocalCache|Project|ProbePath)'`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/source	(cached)
ok  	github.com/yumiaura/seekmoon/internal/store	(cached) [no tests to run]
```

`go test ./internal/source ./internal/model -run 'Test(Repository|Skill|Runwasm)'`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/source	(cached)
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
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

Additional import-boundary command:

`go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/source ./internal/service ./internal/app ./internal/model ./internal/platform ./internal/store`

Result: `github.com/google/go-github/v88/github` and `golang.org/x/oauth2` appear only in `internal/source`, as required.

## Findings

### High: Non-HTTP source readers and repository source results do not carry the complete source envelope

Evidence:

- `internal/model/evidence.go:155` through `internal/model/evidence.go:164` defines `SourceResult` fields `source`, `url`, `path`, `fetched_at`, `status`, `parse_state`, `raw_ref`, `error`, and `value`.
- `internal/source/local_index.go:44` through `internal/source/local_index.go:55` returns local index `SourceResult` values with `Source`, `Path`, `Status`, `ParseState`, `RawRef`, `Error`, and `Value`, but no `FetchedAt`.
- `internal/source/project.go:47` through `internal/source/project.go:52` returns project context `SourceResult` with only `Source`, `Path`, `Status`, and `Value`; it omits `FetchedAt`, `ParseState`, `RawRef`, and any error summary for partial config parse failures.
- `internal/source/repository.go:41` through `internal/source/repository.go:79` returns repository `SourceResult` values with `Source`, `URL`, `Status`, `Error`, and `Value`, but no `FetchedAt`, `ParseState`, or `RawRef`.
- `internal/model/evidence.go:157` through `internal/model/evidence.go:163` marks URL/path/fetched_at/parse_state/raw_ref/error as `omitempty`, so missing envelope fields disappear from JSON instead of remaining observable.

Review basis:

- WP04 defines the shared source result for Batch B: source readers convert HTTP results into `SourceResult` with source label, fetch time, status, parse state, raw reference, and normalized value.
- `parts/055-go-implementation-architecture/070-source-reader-boundary.adoc`, lines `4..42`, states all source results carry stable source label, URL/path, fetched_at, status, parse_state, raw_ref, error, and normalized value.
- The Batch B review risk explicitly requires source result envelopes to carry stable source label, URL/path, fetched_at, status, parse_state, raw_ref, error, and normalized value.

Impact:

Downstream services, raw output, error surface, and reports cannot uniformly audit when a local/project/repository source action happened or where its raw input is represented. This weakens Batch B as the common source-reader substrate and makes later output/report behavior depend on source-specific exceptions.

Required action:

Make source result construction uniform across HTTP, local index, project context, and repository readers. Every `model.SourceResult` producer should set:

- `Source`
- one of `URL` or `Path`
- `FetchedAt`
- `Status`
- `ParseState`
- `RawRef`
- `Error` when failed or partially failed
- `Value` when a normalized value exists

For local/project/repository readers, inject or use a clock consistently, set appropriate parse states, and keep a raw reference such as the file path, raw line reference, repository URL/API path, or project root/config path.

### Medium: Project context parse failures are hidden by the aggregate result status

Evidence:

- `internal/source/project.go:30` through `internal/source/project.go:45` reads module/package config and only assigns successful evidence fields.
- `internal/source/project.go:47` through `internal/source/project.go:52` always returns aggregate `Status: model.StatePresent` once the reader finishes, even if `readFirstConfig` returned `StateFailed` for malformed JSON/TOML.
- `internal/source/project.go:70` through `internal/source/project.go:75` creates failed evidence for malformed config files, but the aggregate source result does not surface the failed source action in `Status`, `ParseState`, or `Error`.

Review basis:

- WP06 requires project context reader behavior around JSON/TOML variants and says read failure affects contextual probe/report.
- The global review rules require `failed` to mean a request, command, or parsing action executed and failed, without collapsing that state into ordinary output.

Impact:

A malformed project file can exist as failed field evidence while the source result still appears wholly present. Downstream report/error behavior will not know that the project context source action partially failed unless it knows to inspect nested evidence fields.

Required action:

Represent project context as a partial source result when any config parse action fails: keep the normalized partial project context, but set source-level `ParseState` and/or `Status` plus `Error` so the failed action remains observable.

## Positive Boundary Checks

These checks passed:

- Stable WP04 labels are `modules_api`, `statistics_api`, and `manifest_api`.
- Modules API reader does not read asset or repository data.
- Empty module summary fields map to `missing`.
- Manifest metadata remains open and raw metadata is preserved.
- Snapshot partial failures keep successful source results.
- Asset URL construction uses manifest version and model-level package relpath derivation.
- `resource.json` 404 maps to `unavailable`.
- Raw signature is preserved while plain signature is derived.
- Source zip is modeled as a `SourceAttempt`, not a final source-resolution decision.
- Moon CLI command result includes command, cwd, exit code, status, and log path.
- Local index parser preserves raw line for valid records and handles malformed sparse input.
- Repository signal does not replace published source resolution in the implementation surface reviewed.
- `go-github` and OAuth2 imports stay in `internal/source`.
- `SkillEntry` and `ModuleSummary` remain separate models.
- Checksum handling records asset state and does not claim provenance.
- Runtime registration does not introduce command-surface behavior.

## Commit Hash If Approved

Not applicable. Rejected; no commit was created.

## Required Follow-Up Before Re-Review

1. Make all `SourceResult` producers carry a complete source envelope, including local/project/repository readers.
2. Preserve project context parse failures at the aggregate source-result level while still returning partial normalized context where useful.
3. Add tests that fail if `fetched_at`, `parse_state`, or `raw_ref` are omitted from local index, project context, and repository source results.
4. Re-run the full Batch B evidence command set.
