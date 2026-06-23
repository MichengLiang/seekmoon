# Batch E Review Report

## Verdict

Rejected.

Batch E implements the A1-A16 black-box acceptance names, journey tests, opt-in integration boundary tests, and `justfile` quality-gate composition, but WP13 cannot be promoted because promotion-critical gates fail:

- `PATH="$(go env GOPATH)/bin:$PATH" just lint`
- `PATH="$(go env GOPATH)/bin:$PATH" just vuln`
- `PATH="$(go env GOPATH)/bin:$PATH" just ci`

No commit was created.

## Review Inputs

- `coordination/prompts/018-batch-e-reviewer.md`
- `coordination/review-packages/018-batch-e-review.md`
- `coordination/handoffs/batch-e-builder.md`
- `coordination/reports/017-batch-e-builder-report.md`
- `010-mandatory-global-context.adoc`
- `020-wbs-dependency-map.adoc`
- `150-wp13-black-box-acceptance-and-quality-gates.adoc`
- `900-source-include-map.adoc`
- `coordination/resource-map.md`
- Principal coordinator protocol `SKILL.md`, `references/index.md`, and `references/30_coordination_and_runtime/review_evidence_promotion.md`

## Exact Include Ranges Used

WP13:

- `parts/060-journeys-and-acceptance/050-black-box-acceptance.adoc`, lines `4..75`
- `parts/055-go-implementation-architecture/110-testing-and-tooling.adoc`, lines `4..53`
- `parts/055-go-implementation-architecture/120-engineering-quality-toolchain.adoc`, lines `59..159`

Global context was read through `010-mandatory-global-context.adoc`; it includes the full shared object, source, command/output, architecture, dependency, and journey context listed in that file.

## Findings

### High: WP13 quality gates do not pass

WP13 cannot be promoted while required quality gates fail.

Review basis:

- `020-wbs-dependency-map.adoc` says the common completion definition includes `just fmt-check`, `just lint`, `just test`, and related local tests.
- `150-wp13-black-box-acceptance-and-quality-gates.adoc` defines WP13 as black-box acceptance plus engineering quality gates.
- `150-wp13-black-box-acceptance-and-quality-gates.adoc` lists `just lint`, `just vuln`, and `just ci` in completion evidence.
- `parts/055-go-implementation-architecture/120-engineering-quality-toolchain.adoc:59` through `:73` lists the standard quality entrypoints, including `lint`, `vuln`, and `ci`.
- `parts/055-go-implementation-architecture/120-engineering-quality-toolchain.adoc:115` defines `ci` as `fmt-check lint test test-race cover vuln mod-check release-check`.
- The Batch E review package explicitly says to treat `just lint`, `just vuln`, and `just ci` as promotion-critical unless there is a stronger WP13-consistent reason to defer them.

Command evidence:

- `PATH="$(go env GOPATH)/bin:$PATH" just lint` failed with 68 issues.
- `PATH="$(go env GOPATH)/bin:$PATH" just vuln` failed with two reachable Go standard-library vulnerabilities under `go1.26.3`.
- `PATH="$(go env GOPATH)/bin:$PATH" just ci` failed at `lint`, so the composed WP13 gate does not pass.

Representative lint failures:

```text
internal/source/fetch.go:92:23: Error return value of `resp.Body.Close` is not checked (errcheck)
internal/platform/exec.go:36:9: G204: Subprocess launched with a potential tainted input or cmd arguments (gosec)
internal/source/local_index.go:101:16: G304: Potential file inclusion via variable (gosec)
internal/source/local_index.go:122:16: G304: Potential file inclusion via variable (gosec)
internal/service/raw.go:63:30: firstArg - fallback always receives "" (unparam)
internal/service/doctor.go:52:6: func environmentPathStatus is unused (unused)
internal/service/helpers.go:50:6: func evidenceIntValue is unused (unused)
tests/acceptance/harness.go:1:1: package-comments: should have a package comment (revive)
tests/acceptance/harness.go:17:6: exported: exported type Harness should have comment or be unexported (revive)
```

Representative vulnerability failures:

```text
Vulnerability #1: GO-2026-5039
    Found in: net/textproto@go1.26.3
    Fixed in: net/textproto@go1.26.4
    Example trace: internal/source/fetch.go:134:26

Vulnerability #2: GO-2026-5037
    Found in: crypto/x509@go1.26.3
    Fixed in: crypto/x509@go1.26.4
    Example traces: internal/output/error.go:52:27, internal/output/jq.go:32:66
```

Impact:

WP13 is the work package that establishes final acceptance and quality-gate closure. Passing behavior tests while the static quality gate, vulnerability gate, and composed `ci` gate fail is not enough to promote the acceptance package.

Required action before re-review:

Make `PATH="$(go env GOPATH)/bin:$PATH" just lint`, `PATH="$(go env GOPATH)/bin:$PATH" just vuln`, and `PATH="$(go env GOPATH)/bin:$PATH" just ci` pass. That likely requires both repository lint cleanup and a Go toolchain update to a version containing the fixed standard library, such as `go1.26.4` or later.

## Positive Checks

The following checks passed and should be preserved during revision:

- A1-A16 all map to visible executable test names in `tests/blackbox/a_acceptance_test.go`.
- `tests/acceptance/harness.go` uses real Cobra command wiring and the real renderer with fake services and temp runtime paths.
- Default black-box and journey tests do not require network, real Moon CLI mutation, GitHub credentials, or external services.
- Journey tests cover library, skill, pipeline, and failure recovery behavior.
- Integration tests are visibly named and skip by default behind explicit environment variables:
  - `SEEKMOON_INTEGRATION_NETWORK`
  - `SEEKMOON_INTEGRATION_MOONCLI`
  - `SEEKMOON_INTEGRATION_GITHUB`
- Golden/schema/failure checks exercise public command behavior through output modes and error surface.
- `justfile` includes `tests` in `fmt` and `fmt-check`.
- `justfile` composes `ci` from the WP13 gate set.
- `.artifacts/coverage.out` exists from `just cover` but is not shown in `git status --short --untracked-files=all`, so it is not accidentally staged or visible as a commit candidate.

## Evidence Commands

`git status --short`

Result: Batch E paths and coordination files were present before this report was written.

```text
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md
 M bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md
 M justfile
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-e-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/017-batch-e-builder.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/018-batch-e-reviewer.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/017-batch-e-builder-report.md
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/018-batch-e-review.md
?? tests/
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

Result: failed.

```text
68 issues:
* errcheck: 1
* gosec: 8
* prealloc: 3
* revive: 50
* staticcheck: 3
* unparam: 1
* unused: 2
error: Recipe `lint` failed on line 13 with exit code 1
```

`PATH="$(go env GOPATH)/bin:$PATH" just test`

Result: passed.

```text
DONE 113 tests, 3 skipped in 0.550s
```

`PATH="$(go env GOPATH)/bin:$PATH" just test-race`

Result: passed.

```text
DONE 113 tests, 3 skipped in 1.793s
```

`PATH="$(go env GOPATH)/bin:$PATH" just cover`

Result: passed and generated `.artifacts/coverage.out`.

```text
total:                                  (statements)            54.3%
```

`PATH="$(go env GOPATH)/bin:$PATH" just vuln`

Result: failed.

```text
Your code is affected by 2 vulnerabilities from the Go standard library.
error: Recipe `vuln` failed on line 33 with exit code 3
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

Result: failed at `lint`.

```text
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal tests)"
command -v golangci-lint >/dev/null
golangci-lint run ./...
...
error: Recipe `lint` failed on line 13 with exit code 1
```

## Required Actions Before Re-Review

1. Make `just lint` pass under `PATH="$(go env GOPATH)/bin:$PATH"`.
2. Make `just vuln` pass under `PATH="$(go env GOPATH)/bin:$PATH"`.
3. Make `just ci` pass under `PATH="$(go env GOPATH)/bin:$PATH"`.
4. Preserve the current A1-A16, journey, and opt-in integration coverage while repairing the gates.
