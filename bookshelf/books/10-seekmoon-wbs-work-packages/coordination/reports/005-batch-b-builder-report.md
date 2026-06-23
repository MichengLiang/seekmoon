# Batch B Builder Report

## Summary

Batch B implements source readers and local probe primitives for WP04 through WP07 on top of the Batch A substrate commit `0f83682 Implement SeekMoon Batch A substrate`.

The implementation adds:

- shared HTTP fetch behavior with context-aware requests, source envelopes, JSON parse state, limited retry, and evidence status classification;
- Mooncakes modules/statistics/manifest readers and a sync service skeleton that writes snapshot source state without ranking or adoption behavior;
- static asset readers for `module_index.json`, `package_data.json`, optional `resource.json`, and source zip attempts;
- local Moon CLI, local registry JSONL, local cache/core source, and project context readers;
- GitHub repository signal and skill source readers;
- runtime registration for Batch B source/service primitives.

The work intentionally does not implement output renderers, full CLI behavior, ranking/adoption decisions, or Batch D service flows.

## Files Changed

Implementation files:

- `go.mod`
- `go.sum`
- `internal/app/runtime.go`
- `internal/app/runtime_test.go`
- `internal/model/repository.go`
- `internal/service/sync.go`
- `internal/service/sync_test.go`
- `internal/source/assets.go`
- `internal/source/assets_test.go`
- `internal/source/fetch.go`
- `internal/source/local_cache.go`
- `internal/source/local_index.go`
- `internal/source/local_test.go`
- `internal/source/mooncakes.go`
- `internal/source/mooncakes_test.go`
- `internal/source/mooncli.go`
- `internal/source/normalize.go`
- `internal/source/project.go`
- `internal/source/repository.go`
- `internal/source/repository_skill_test.go`
- `internal/source/skills.go`
- `internal/testutil/fixture.go`

Report file:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/005-batch-b-builder-report.md`

Pre-existing coordination changes observed and left untouched:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-b-builder.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/005-batch-b-builder.md`

## Reading Evidence

Prompt and handoff read:

- `coordination/prompts/005-batch-b-builder.md`
- `coordination/handoffs/batch-b-builder.md`

Required WBS and coordination files read:

- `010-mandatory-global-context.adoc`
- `020-wbs-dependency-map.adoc`
- `060-wp04-public-http-sources-and-snapshot.adoc`
- `070-wp05-static-assets-and-api-profile.adoc`
- `080-wp06-local-toolchain-source-and-probe-primitives.adoc`
- `090-wp07-repository-and-skill-sources.adoc`
- `900-source-include-map.adoc`
- `coordination/resource-map.md`
- `coordination/coordination-state.md`

Exact include ranges read for Batch B:

- WP04:
  - `parts/020-evidence-boundary/010-public-http-sources.adoc` lines `4..5;9..46;50..75;79..124`
  - `parts/055-go-implementation-architecture/070-source-reader-boundary.adoc` lines `4..42`
  - `parts/055-go-implementation-architecture/050-data-flow.adoc` lines `8..22;26..39`
- WP05:
  - `parts/020-evidence-boundary/020-static-assets.adoc` lines `4..5;9..42;46..83;87..95;99..107;111..127`
  - `parts/030-canonical-model/040-api-profile.adoc` lines `4..5;9..60;65..133`
  - `parts/040-command-workbench/040-api-and-source.adoc` lines `4..32`
- WP06:
  - `parts/020-evidence-boundary/030-local-toolchain-and-cache.adoc` lines `4..5;9..43;47..56;59..68`
  - `parts/020-evidence-boundary/040-project-context.adoc` lines `4..43`
  - `parts/030-canonical-model/060-source-resolution-and-probe-result.adoc` lines `77..132`
- WP07:
  - `parts/020-evidence-boundary/050-repository-signals.adoc` lines `4..46`
  - `parts/020-evidence-boundary/010-public-http-sources.adoc` lines `128..169`
  - `parts/030-canonical-model/050-skill-profile.adoc` lines `4..94`
  - `backmatter/appendix-f-go-implementation-dependencies.adoc` lines `29..37;105..113`

Protocol references used for role/report boundary:

- `principal-coordinator-protocol/SKILL.md`
- `principal-coordinator-protocol/references/index.md`
- `principal-coordinator-protocol/references/30_coordination_and_runtime/principal_coordinator_protocol.md`
- `principal-coordinator-protocol/references/30_coordination_and_runtime/review_evidence_promotion.md`

## Implementation Notes By WP

### Shared Source Layer

`internal/source/fetch.go` provides the shared fetch primitive:

- context-aware `GET` requests through a supplied or default `http.Client`;
- timeout inherited from the supplied client or `platform.NewHTTPClient`;
- default redirect behavior through `net/http`;
- limited retry with `github.com/cenkalti/backoff/v5`;
- status mapping for present, failed, and unavailable HTTP evidence;
- JSON decoding with parse state;
- stable source result envelope fields: source label, URL, fetched_at, status, parse_state, raw_ref, error, and value.

The source layer stays separate from output rendering and adoption judgment.

### WP04 Public HTTP Sources And Snapshot

`internal/source/mooncakes.go` implements:

- `modules_api` reader for `/api/v0/modules`;
- `statistics_api` reader for `/api/v0/modules/statistics`;
- `manifest_api` reader for `/api/v0/manifest/{owner}/{name}`.

Mooncakes module summaries map empty description, keywords, repository, and license into `missing` evidence rather than silently dropping the field. Manifest metadata preserves the raw metadata map and normalizes only known fields. Manifest module/name mismatch fails the parse state instead of creating a merged profile.

`internal/service/sync.go` adds a sync skeleton that composes modules/statistics source results and writes snapshot state. Partial source failures remain recorded in the snapshot source list and do not erase successful source results.

### WP05 Static Assets And API Profile

`internal/source/assets.go` implements:

- `module_index.json` reader;
- `package_data.json` reader;
- optional `resource.json` reader where HTTP 404 maps to `unavailable`;
- source zip attempt reader that returns a `SourceAttempt` and a file summary, not a final source-resolution judgment.

The module index parser accepts canonical `childs` and compatibility `children`, while output remains canonical `Childs`. Package relpath derivation uses `model.PackageRelPath`. Asset URL construction now preserves nested relpath path segments while escaping each segment.

Package data keeps raw HTML signatures and derives plain signatures separately. The raw signature is not replaced by the derived plain form.

### WP06 Local Toolchain Source And Probe Primitives

`internal/source/mooncli.go` wraps a runner and maps command observations into `model.CommandResult`, recording command, cwd, exit code, status, and log path. Command failure is recorded as local evidence only.

`internal/source/local_index.go` parses local registry JSONL with sparse/malformed-line handling and raw line preservation.

`internal/source/local_cache.go` checks local cache/core source candidate paths and returns source attempts without mutating upstream metadata.

`internal/source/project.go` reads project context from `moon.mod.json`/`moon.mod` and `moon.pkg.json`/`moon.pkg`, parses JSON or TOML, records evidence-bearing project context, and does not hide manifest mutation inside primitives.

### WP07 Repository And Skill Sources

`internal/source/repository.go` implements GitHub repository signal reading through `go-github`. Unsupported repository hosts map to `unknown`; unreachable GitHub calls map to `failed`. Missing `pushed_at` now remains `unknown` instead of formatting the Go zero time as a present value.

`go-github` and `oauth2` imports are confined to `internal/source`; import inspection confirmed they do not enter `internal/app`, `internal/service`, `internal/model`, `internal/platform`, or `internal/store`.

`internal/source/skills.go` implements skill list/detail/asset readers. Empty skill package remains an empty root marker. Skill entries and module summaries stay separate model objects. Asset/checksum fetches record state without implying checksum provenance.

## Tests And Command Evidence

Required verification commands:

```text
$ go test ./...
?   	github.com/yumiaura/seekmoon/cmd/seekmoon	[no test files]
ok  	github.com/yumiaura/seekmoon/internal/app	0.003s
ok  	github.com/yumiaura/seekmoon/internal/cli	0.003s
ok  	github.com/yumiaura/seekmoon/internal/contract	(cached)
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
ok  	github.com/yumiaura/seekmoon/internal/platform	(cached)
ok  	github.com/yumiaura/seekmoon/internal/service	0.062s
ok  	github.com/yumiaura/seekmoon/internal/source	0.203s
ok  	github.com/yumiaura/seekmoon/internal/store	(cached)
?   	github.com/yumiaura/seekmoon/internal/testutil	[no test files]
```

```text
$ go test ./internal/source ./internal/service -run 'Test(Mooncakes|Sync)'
ok  	github.com/yumiaura/seekmoon/internal/source	0.078s
ok  	github.com/yumiaura/seekmoon/internal/service	0.059s
```

```text
$ go test ./internal/source ./internal/model -run 'Test(Asset|ModuleIndex|PackageData|Relpath|SourceZip)'
ok  	github.com/yumiaura/seekmoon/internal/source	0.133s
ok  	github.com/yumiaura/seekmoon/internal/model	(cached) [no tests to run]
```

```text
$ go test ./internal/source ./internal/store -run 'Test(MoonCLI|LocalIndex|LocalCache|Project|ProbePath)'
ok  	github.com/yumiaura/seekmoon/internal/source	0.006s
ok  	github.com/yumiaura/seekmoon/internal/store	(cached) [no tests to run]
```

```text
$ go test ./internal/source ./internal/model -run 'Test(Repository|Skill|Runwasm)'
ok  	github.com/yumiaura/seekmoon/internal/source	0.114s
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
```

```text
$ just fmt-check
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal)"
```

```text
$ goreleaser check
• checking                                  path=.goreleaser.yaml
• 1 configuration file(s) validated
• thanks for using GoReleaser!
```

Module verification:

```text
$ go mod verify
all modules verified
```

Import boundary inspection:

```text
$ go list -f '{{.ImportPath}} {{join .Imports "\n"}}' ./internal/source ./internal/service ./internal/app ./internal/model ./internal/platform ./internal/store
```

Relevant result: `github.com/google/go-github/v88/github` and `golang.org/x/oauth2` appear only in `github.com/yumiaura/seekmoon/internal/source`.

## Known Gaps Or Blocked Gates

`just mod-check` failed at its dirty-diff guard:

```text
$ just mod-check
go mod tidy
git diff --exit-code -- go.mod go.sum
...
error: Recipe `mod-check` failed on line 37 with exit code 1
```

The diff shown by that guard is the intended Batch B dependency metadata change:

- `github.com/cenkalti/backoff/v5 v5.0.3`
- `github.com/google/go-github/v88 v88.0.0`
- `github.com/pelletier/go-toml/v2 v2.4.0`
- `golang.org/x/oauth2 v0.36.0`
- indirect `github.com/google/go-querystring v1.2.0`

Because the handoff explicitly says do not commit, the repository must remain dirty after adding these dependencies. `go mod tidy` was run and `go mod verify` passed, so the failure is the recipe's expected "dependency metadata differs from HEAD" gate rather than a module integrity failure.

No Batch B output renderer, full CLI behavior, ranking/adoption decision, or Batch D service flow was implemented.

## Suggested Reviewer Focus

Review should focus on:

- whether all source result statuses match the WBS evidence vocabulary, especially `missing`, `unknown`, `failed`, and `unavailable`;
- whether `SyncService` preserves partial source results without implying ranking or adoption;
- whether asset URL construction and package relpath derivation match real Mooncakes asset layouts for nested packages;
- whether source zip summary in `SourceAttempt.Path` is acceptable as a lightweight observable attempt detail, or whether a model field should be added in a later batch;
- whether repository signal scope remains bounded and does not replace published source resolution;
- whether runtime registration introduces any unexpected dependency coupling outside `internal/source`.
