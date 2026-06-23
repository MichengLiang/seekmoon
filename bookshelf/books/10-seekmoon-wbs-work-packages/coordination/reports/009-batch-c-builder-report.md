# Batch C Builder Report

## Summary

Batch C implements WP08 Output Pipeline and Error Surface, then WP09 CLI Command Surface, on top of the promoted Batch A and Batch B baseline.

The implementation adds:

- `internal/output` renderer package with pretty, JSON, jq, shape, schema, and error-surface projections.
- `model.OutputMode`, `model.SurfaceError`, and command input/request helper structs needed by output and CLI surfaces.
- Contract shape/schema entries for the currently exposed command output schema ids.
- Runtime renderer and fakeable service registry seams.
- Cobra command files for all WP09 command surfaces with shared output-mode flags, syntax validation, input parsing, thin service handlers, and process exit mapping helpers.
- Tests for output projections, jq failure mapping, command registration, parse failure exit code, service failure exit code, candidate parsing, and thin handler dispatch.

No concrete Batch D service behavior was implemented.

No commit was created.

## Files Changed

Implementation and test paths:

- `go.mod`
- `go.sum`
- `internal/app/runtime.go`
- `internal/cli/api.go`
- `internal/cli/compare.go`
- `internal/cli/doctor.go`
- `internal/cli/flags.go`
- `internal/cli/probe.go`
- `internal/cli/raw.go`
- `internal/cli/record.go`
- `internal/cli/report.go`
- `internal/cli/root.go`
- `internal/cli/root_test.go`
- `internal/cli/search.go`
- `internal/cli/skill.go`
- `internal/cli/source.go`
- `internal/cli/sync.go`
- `internal/cli/view.go`
- `internal/contract/schemas.go`
- `internal/contract/shapes.go`
- `internal/model/output.go`
- `internal/output/error.go`
- `internal/output/jq.go`
- `internal/output/json.go`
- `internal/output/pretty.go`
- `internal/output/render.go`
- `internal/output/render_test.go`
- `internal/output/schema.go`
- `internal/output/shape.go`
- `internal/service/registry.go`
- `internal/testutil/golden.go`

Report path:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/009-batch-c-builder-report.md`

Observed pre-existing coordination process changes left untouched:

- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/coordination-state.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/resource-map.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/role-registry.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/handoffs/batch-c-builder.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/prompts/009-batch-c-builder.md`

## Reading Evidence

Prompt and handoff read:

- `coordination/prompts/009-batch-c-builder.md`
- `coordination/handoffs/batch-c-builder.md`

Required WBS and coordination files read:

- `010-mandatory-global-context.adoc`
- `020-wbs-dependency-map.adoc`
- `100-wp08-output-pipeline-and-error-surface.adoc`
- `110-wp09-cli-command-surface.adoc`
- `900-source-include-map.adoc`
- `coordination/resource-map.md`
- `coordination/coordination-state.md`

Batch C exact include ranges read:

- WP08:
  - `parts/050-output-contracts/010-pretty-text.adoc`, lines `4..45`
  - `parts/050-output-contracts/020-json-output.adoc`, lines `4..58`
  - `parts/050-output-contracts/030-built-in-jq.adoc`, lines `4..36`
  - `parts/050-output-contracts/040-shape.adoc`, lines `4..47`
  - `parts/050-output-contracts/050-schema.adoc`, lines `4..32`
  - `parts/050-output-contracts/060-error-surface.adoc`, lines `4..59`
  - `parts/055-go-implementation-architecture/080-output-pipeline.adoc`, lines `4..51`
- WP09:
  - `parts/040-command-workbench/010-command-map.adoc`, lines `4..73`
  - `parts/055-go-implementation-architecture/060-control-flow.adoc`, lines `4..57`
  - `backmatter/appendix-c-command-reference.adoc`, lines `4..84`

Global full-file includes from the mandatory global context were also read:

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
- `backmatter/appendix-e-status-and-error-matrix.adoc`
- `parts/060-journeys-and-acceptance/010-library-discovery-journey.adoc`
- `parts/060-journeys-and-acceptance/020-skill-discovery-journey.adoc`
- `parts/060-journeys-and-acceptance/030-pipeline-journey.adoc`
- `parts/060-journeys-and-acceptance/040-failure-recovery-journey.adoc`
- `parts/060-journeys-and-acceptance/050-black-box-acceptance.adoc`

Protocol references read for role and promotion boundary:

- `principal-coordinator-protocol/SKILL.md`
- `principal-coordinator-protocol/references/index.md`

## Implementation Notes By WP

### WP08 Output Pipeline And Error Surface

`internal/output` now owns the renderer contract and projections:

- `render.go` defines `Renderer`, `Request`, `DefaultRenderer`, and output-mode dispatch.
- `json.go` produces canonical JSON projection with schema id; it marshals the canonical object and does not serialize pretty text.
- `jq.go` uses `github.com/itchyny/gojq@v0.12.19`; jq consumes the JSON projection and maps parse/eval failures to error surface.
- `shape.go` and `schema.go` read from `internal/contract`, not sample result values.
- `pretty.go` renders low-noise terminal projections for current canonical objects such as search, snapshot, manifest, package data, and probe results.
- `error.go` renders `model.SurfaceError` with command, object, source, state, meaning, recovery, and optional log path.

`internal/output` import inspection confirms it does not import `internal/source`.

`internal/contract` gained generic shape/schema entries for command schema ids beyond the already-defined search/adoption contracts. These are intentionally explicit contract definitions, not reflection-derived from sample output.

Projection-specific failures return an `output.SurfaceError` after writing the error surface, so CLI code can map them to process exit code `1`.

### WP09 CLI Command Surface

The Cobra surface now includes the WP09 files:

- `root.go`
- `doctor.go`
- `sync.go`
- `search.go`
- `view.go`
- `api.go`
- `source.go`
- `skill.go`
- `compare.go`
- `probe.go`
- `record.go`
- `report.go`
- `raw.go`
- `flags.go`

Shared output mode handling supports:

- `--json`
- `--jq <expr>`
- `--shape`
- `--schema`

`--jq` without an expression is mapped to a parse failure. Multiple output modes are rejected by shared validation.

Command handlers follow the required sequence:

```text
parse args -> build input -> context -> service -> render -> exit mapping
```

The handlers call fakeable service interfaces defined in `internal/service/registry.go`. Pending service implementations return explicit “outside Batch C” service errors. This creates the command/service seam for Batch D without implementing Batch D behavior.

Command input parsing distinguishes candidate numbers from module coordinates. Parameters use object names and candidate numbers; hidden implementation URLs are not exposed in command arguments.

`internal/cli` import inspection confirms it does not import `internal/source`.

## Tests And Command Evidence

Required verification commands:

```text
$ go test ./...
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
```

```text
$ go test ./internal/output ./internal/contract
ok  	github.com/yumiaura/seekmoon/internal/output	0.004s
ok  	github.com/yumiaura/seekmoon/internal/contract	(cached)
```

```text
$ go test ./internal/cli
ok  	github.com/yumiaura/seekmoon/internal/cli	0.004s
```

```text
$ go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/output ./internal/cli ./internal/service ./internal/app
github.com/yumiaura/seekmoon/internal/output context encoding/json errors fmt github.com/itchyny/gojq github.com/yumiaura/seekmoon/internal/contract github.com/yumiaura/seekmoon/internal/model io strings
github.com/yumiaura/seekmoon/internal/cli context fmt github.com/spf13/cobra github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/output github.com/yumiaura/seekmoon/internal/service io strconv strings
github.com/yumiaura/seekmoon/internal/service context fmt github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/source github.com/yumiaura/seekmoon/internal/store time
github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/output github.com/yumiaura/seekmoon/internal/platform github.com/yumiaura/seekmoon/internal/service github.com/yumiaura/seekmoon/internal/source github.com/yumiaura/seekmoon/internal/store net/http time
```

```text
$ just fmt-check
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal)"
```

```text
$ go mod tidy
```

Passed with no command output.

```text
$ go mod verify
all modules verified
```

```text
$ goreleaser check
• checking                                  path=.goreleaser.yaml
• 1 configuration file(s) validated
• thanks for using GoReleaser!
```

Additional command run:

```text
$ just mod-check
go mod tidy
git diff --exit-code -- go.mod go.sum
...
error: Recipe `mod-check` failed on line 37 with exit code 1
```

The failure is only the dirty-diff guard for intended Batch C dependency metadata:

- `github.com/itchyny/gojq v0.12.19`
- `github.com/itchyny/timefmt-go v0.1.8` indirect
- `golang.org/x/sys` upgraded from `v0.26.0` to `v0.38.0`

## Known Gaps Or Blocked Gates

No required Batch C verification command failed.

`just mod-check` fails only because the repository remains dirty with intended Batch C dependency metadata and the handoff says do not commit. `go mod tidy` and `go mod verify` both pass.

Concrete service behavior for doctor/search/view/api/source/skill/compare/probe/record/report/raw remains pending by design for Batch D. The command shell currently routes to fakeable pending service seams unless an actual service implementation already exists, such as the Batch B sync skeleton.

## Suggested Reviewer Focus

Review should focus on:

- whether `internal/output` projections are contract/canonical-object based and do not import `internal/source`;
- whether jq consumes JSON projection and maps projection failure to error surface;
- whether normal pretty text stays low-noise and avoids recovery/tutorial text;
- whether `internal/cli` command handlers stay thin and avoid source-reader imports;
- whether parse failures map to exit code `2` and service/projection failures map to exit code `1`;
- whether the fakeable service registry creates enough seam for Batch D without implementing Batch D behavior.
