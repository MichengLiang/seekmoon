# Batch E Re-Review Report

## Verdict

Approved.

Batch E Revision 1 fixes the prior WP13 rejection findings. The lint gate, vulnerability gate, and full `just ci` composition now pass under the required command prefix. The A1-A16 acceptance mapping, offline default tests, opt-in integration boundaries, and final quality-gate composition remain intact.

## Review Inputs

- `coordination/prompts/020-batch-e-re-reviewer.md`
- `coordination/review-packages/020-batch-e-re-review.md`
- `coordination/review-packages/018-batch-e-review.md`
- `coordination/reports/018-batch-e-review-report.md`
- `coordination/reports/019-batch-e-revision-1-builder-report.md`
- `coordination/handoffs/batch-e-revision-1-builder.md`
- `010-mandatory-global-context.adoc`
- `020-wbs-dependency-map.adoc`
- `150-wp13-black-box-acceptance-and-quality-gates.adoc`
- `900-source-include-map.adoc`
- Principal coordinator protocol `SKILL.md`, `references/index.md`, and `references/30_coordination_and_runtime/review_evidence_promotion.md`

## Exact Include Ranges Used

WP13:

- `parts/060-journeys-and-acceptance/050-black-box-acceptance.adoc`, lines `4..75`
- `parts/055-go-implementation-architecture/110-testing-and-tooling.adoc`, lines `4..53`
- `parts/055-go-implementation-architecture/120-engineering-quality-toolchain.adoc`, lines `59..159`

Global context was read through `010-mandatory-global-context.adoc`; it includes the full shared object, source, command/output, architecture, dependency, and journey context listed in that file.

## Re-Review Findings

No blocking findings.

## Prior Rejection Findings

### `just lint`

Fixed.

Evidence:

- `PATH="$(go env GOPATH)/bin:$PATH" just lint` passed with `0 issues`.
- `internal/platform/exec.go:41` through `:43` keeps the command-runner `#nosec G204` suppression tied to argv-based SeekMoon service command construction rather than shell string execution.
- `internal/platform/fs.go:33` through `:35` keeps the filesystem read suppression tied to the explicit `platform.FS` host filesystem capability and caller-owned path confinement.
- `internal/source/local_index.go:105` through `:113` validates walked local-index paths remain under the configured root before reading `.index` files.
- `internal/source/local_index.go:134` through `:148` constrains HEAD/ref reads to the configured local registry checkout metadata.
- `internal/testutil/golden.go:18` through `:19` confines golden reads to test fixture paths.
- Exported-surface comments were added across `cmd/`, `internal/`, and `tests/acceptance`; the reviewed examples describe contracts or object roles rather than hiding behavior changes.

### `just vuln`

Fixed.

Evidence:

- `go version` reports `go version go1.26.4 linux/amd64`.
- `PATH="$(go env GOPATH)/bin:$PATH" just vuln` passed with `No vulnerabilities found.`
- `PATH="$(go env GOPATH)/bin:$PATH" just mod-check` passed, so keeping `go.mod` at the existing directive did not leave a module metadata diff.

### `just ci`

Fixed.

Evidence:

- `PATH="$(go env GOPATH)/bin:$PATH" just ci` passed end-to-end.
- `justfile:48` still composes `ci` as `fmt-check lint test test-race cover vuln mod-check release-check`; the composition was not weakened to bypass gates.

## Positive Boundary Checks

Re-checked and passed:

- A1-A16 all map to visible executable test names in `tests/blackbox/a_acceptance_test.go:12` through `:209`.
- Default black-box and journey tests use `tests/acceptance/harness.go`, real Cobra command wiring, the real renderer, fake services, and temp runtime paths.
- Default tests do not depend on network, real Moon CLI mutation, GitHub token, or external services.
- `tests/integration/integration_test.go:8` through `:20` keeps the three integration tests visibly named.
- `tests/integration/integration_test.go:23` through `:27` skips integration tests unless their explicit environment variables are set.
- Journey tests still cover library, skill, pipeline, and failure recovery behavior.
- Shape/schema/jq/error-surface checks remain public command behavior checks rather than implementation-only assertions.
- `justfile:3` through `:9` includes `tests` in formatting commands.
- `.artifacts/coverage.out` exists from coverage runs but is not listed by `git status --short --untracked-files=all`; no generated coverage artifact is committed.

## Evidence Commands

`git status --short`

Result: approved Batch E, remediation, and coordination paths were present before this report was written.

```text
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md
 M cmd/seekmoon/main.go
 M internal/app/command.go
 M internal/app/runtime.go
 M internal/cli/doctor.go
 M internal/cli/report.go
 M internal/cli/root.go
 M internal/cli/root_test.go
 M internal/cli/search.go
 M internal/cli/sync.go
 M internal/contract/schemas.go
 M internal/contract/shapes.go
 M internal/model/api.go
 M internal/model/evidence.go
 M internal/model/manifest.go
 M internal/model/module.go
 M internal/model/output.go
 M internal/model/probe.go
 M internal/model/project.go
 M internal/model/record.go
 M internal/model/report.go
 M internal/model/repository.go
 M internal/model/schema.go
 M internal/model/skill.go
 M internal/model/snapshot.go
 M internal/model/source.go
 M internal/output/error.go
 M internal/output/jq.go
 M internal/output/json.go
 M internal/output/pretty.go
 M internal/output/render.go
 M internal/output/schema.go
 M internal/output/shape.go
 M internal/platform/clock.go
 M internal/platform/env.go
 M internal/platform/exec.go
 M internal/platform/fs.go
 M internal/platform/http.go
 M internal/service/api.go
 M internal/service/batch_d_test.go
 M internal/service/compare.go
 M internal/service/doctor.go
 M internal/service/helpers.go
 M internal/service/probe.go
 M internal/service/raw.go
 M internal/service/record.go
 M internal/service/registry.go
 M internal/service/report.go
 M internal/service/search.go
 M internal/service/skill.go
 M internal/service/source.go
 M internal/service/sync.go
 M internal/service/view.go
 M internal/source/assets.go
 M internal/source/assets_test.go
 M internal/source/fetch.go
 M internal/source/local_cache.go
 M internal/source/local_index.go
 M internal/source/mooncakes.go
 M internal/source/mooncakes_test.go
 M internal/source/mooncli.go
 M internal/source/project.go
 M internal/source/repository.go
 M internal/source/repository_skill_test.go
 M internal/source/skills.go
 M internal/store/atomic.go
 M internal/store/cache_store.go
 M internal/store/paths.go
 M internal/store/probe_store.go
 M internal/store/record_store.go
 M internal/store/report_store.go
 M internal/store/session_store.go
 M internal/store/snapshot_store.go
 M internal/store/source_store.go
 M internal/testutil/fixture.go
 M internal/testutil/golden.go
 M justfile
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-e-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-e-revision-1-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/017-batch-e-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/018-batch-e-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/019-batch-e-revision-1-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/020-batch-e-re-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/017-batch-e-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/018-batch-e-review-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/019-batch-e-revision-1-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/018-batch-e-review.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/020-batch-e-re-review.md
?? tests/acceptance/harness.go
?? tests/blackbox/a_acceptance_test.go
?? tests/integration/integration_test.go
?? tests/journey/journey_test.go
```

`go version`

Result:

```text
go version go1.26.4 linux/amd64
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
?    github.com/yumiaura/seekmoon/tests/acceptance [no test files]
ok   github.com/yumiaura/seekmoon/tests/blackbox (cached)
ok   github.com/yumiaura/seekmoon/tests/integration (cached)
ok   github.com/yumiaura/seekmoon/tests/journey (cached)
```

`go test ./tests/blackbox ./tests/journey`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/tests/blackbox (cached)
ok   github.com/yumiaura/seekmoon/tests/journey (cached)
```

`go test ./tests/integration`

Result: passed with default skips.

```text
ok   github.com/yumiaura/seekmoon/tests/integration (cached)
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
DONE 113 tests, 3 skipped in 0.712s
```

`PATH="$(go env GOPATH)/bin:$PATH" just test-race`

Result: passed.

```text
DONE 113 tests, 3 skipped in 1.876s
```

`PATH="$(go env GOPATH)/bin:$PATH" just cover`

Result: passed.

```text
total:                                  (statements)            54.2%
```

`PATH="$(go env GOPATH)/bin:$PATH" just vuln`

Result: passed.

```text
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

```text
golangci-lint run ./...
0 issues.
DONE 113 tests, 3 skipped in 0.548s
DONE 113 tests, 3 skipped in 1.704s
No vulnerabilities found.
all modules verified
1 configuration file(s) validated
```

## Commit Scope

Approved commit paths:

- `cmd/seekmoon/main.go`
- `internal`
- `justfile`
- `tests`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination`
