# Batch B Re-Review Report

## Verdict

Approved.

The Batch B revision fixes both prior rejection findings. Local index, project context, repository, and pre-fetch validation source results now carry the required source envelope fields, and project context parse failures are visible at the aggregate source-result level while preserving useful partial context.

## Review Inputs

- `coordination/review-packages/008-batch-b-re-review.md`
- `coordination/review-packages/006-batch-b-review.md`
- `coordination/reports/006-batch-b-review-report.md`
- `coordination/reports/007-batch-b-revision-1-builder-report.md`
- Prior Batch B WBS basis and exact include ranges listed in `coordination/reports/006-batch-b-review-report.md`
- Principal coordinator protocol `SKILL.md`, `references/index.md`, and `references/30_coordination_and_runtime/review_evidence_promotion.md`

## Re-Review Findings

No blocking findings.

## Prior Rejection Findings

### Source Envelope Completeness

Fixed.

Evidence:

- `internal/source/local_index.go:46` through `internal/source/local_index.go:57` returns local index source results with `Source`, `Path`, `FetchedAt`, `Status`, `ParseState`, `RawRef`, `Error` on failure, and `Value` on success.
- `internal/source/project.go:62` through `internal/source/project.go:70` returns project context source results with `Source`, `Path`, `FetchedAt`, `Status`, `ParseState`, `RawRef`, `Error`, and `Value`.
- `internal/source/repository.go:108` through `internal/source/repository.go:117` constructs the repository source envelope with `Source`, `URL`, `FetchedAt`, `Status`, `ParseState`, `RawRef`, and `Error`; the signal method fills `Value` for unknown, failed, and present outcomes.
- `internal/source/assets.go:39` through `internal/source/assets.go:58` and `internal/source/mooncakes.go:44` through `internal/source/mooncakes.go:48` now build failed pre-fetch validation envelopes with timestamp and raw reference.
- `internal/service/sync.go:41` through `internal/service/sync.go:57` preserves envelope fields when erasing typed source values.

Test evidence:

- `TestLocalIndexSourceEnvelopeComplete`
- `TestProjectContextJSONAndTOML`
- `TestRepositorySignalMapping`

### Project Context Aggregate Failure

Fixed.

Evidence:

- `internal/source/project.go:53` through `internal/source/project.go:60` sets aggregate `Status` and `ParseState` to `failed` and records an error summary if any config parse action fails.
- `internal/source/project.go:32` through `internal/source/project.go:52` preserves successful module/package fields while copying failed config evidence into the corresponding nested field.
- `internal/source/local_test.go:91` through `internal/source/local_test.go:109` verifies malformed package config is observable at aggregate level while successful module context remains present.

## Positive Boundary Checks

Re-checked and passed:

- Stable WP04 labels remain `modules_api`, `statistics_api`, and `manifest_api`.
- Source readers return state-bearing source results and do not render output or make adoption/ranking decisions.
- Modules API reader does not read asset or repository data.
- Empty module summary fields map to `missing`.
- Manifest metadata remains open and raw metadata is preserved.
- Snapshot partial failure semantics keep successful source results.
- Asset URL construction uses manifest version and model-level package relpath derivation.
- `resource.json` 404 maps to `unavailable`.
- Raw signature is preserved while plain signature is derived.
- Source zip remains modeled as a `SourceAttempt`.
- Moon CLI command failure remains local evidence with command, cwd, exit code, status, and log path.
- Local index parser preserves raw line and handles malformed sparse input.
- Project context reader does not hide manifest mutation.
- Repository signal does not replace published source resolution.
- `go-github` and OAuth2 imports remain confined to `internal/source`.
- `SkillEntry` and `ModuleSummary` remain separate models.
- Checksum handling records asset state and does not claim provenance.
- Runtime registration adds source/service dependencies without command-surface behavior.

## Evidence Commands

`git status --short`

Result: only approved Batch B implementation and coordination paths were present in the worktree before this report was written.

```text
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md
 M go.mod
 M go.sum
 M internal/app/runtime.go
 M internal/app/runtime_test.go
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-b-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-b-revision-1-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/005-batch-b-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/006-batch-b-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/007-batch-b-revision-1-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/008-batch-b-re-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/005-batch-b-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/006-batch-b-review-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/007-batch-b-revision-1-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/006-batch-b-review.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/008-batch-b-re-review.md
?? internal/model/repository.go
?? internal/service/
?? internal/source/
?? internal/testutil/
```

`go test ./...`

Result: passed.

```text
?    github.com/yumiaura/seekmoon/cmd/seekmoon [no test files]
ok   github.com/yumiaura/seekmoon/internal/app (cached)
ok   github.com/yumiaura/seekmoon/internal/cli (cached)
ok   github.com/yumiaura/seekmoon/internal/contract (cached)
ok   github.com/yumiaura/seekmoon/internal/model (cached)
ok   github.com/yumiaura/seekmoon/internal/platform (cached)
ok   github.com/yumiaura/seekmoon/internal/service (cached)
ok   github.com/yumiaura/seekmoon/internal/source (cached)
ok   github.com/yumiaura/seekmoon/internal/store (cached)
?    github.com/yumiaura/seekmoon/internal/testutil [no test files]
```

`go test ./internal/source ./internal/service -run 'Test(Mooncakes|Sync)'`

Result: passed.

`go test ./internal/source ./internal/model -run 'Test(Asset|ModuleIndex|PackageData|Relpath|SourceZip)'`

Result: passed.

`go test ./internal/source ./internal/store -run 'Test(MoonCLI|LocalIndex|LocalCache|Project|ProbePath)'`

Result: passed.

`go test ./internal/source ./internal/model -run 'Test(Repository|Skill|Runwasm)'`

Result: passed.

`go test ./internal/source -run 'Test(LocalIndexSourceEnvelopeComplete|ProjectContextPartialParseFailureObservable|ProjectContextJSONAndTOML|RepositorySignalMapping)'`

Result: passed.

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

Additional import-boundary check:

`go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/source ./internal/service ./internal/app ./internal/model ./internal/platform ./internal/store`

Result: `github.com/google/go-github/v88/github` and `golang.org/x/oauth2` appear only in `internal/source`.

## Commit

Approved for commit.
