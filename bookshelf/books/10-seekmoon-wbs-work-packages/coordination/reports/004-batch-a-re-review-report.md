# Batch A Re-Review Report

## Verdict

Approved.

The Revision 1 work fixes the two original rejection findings and clears the previously blocked quality gates. Batch A now satisfies the WP01-WP03 review basis for promotion.

Commit hash: pending at report write time; see final commit update below.

## Review Basis

This re-review reused the full review basis from `coordination/review-packages/002-batch-a-review.md` and the initial review report:

- WBS global context, dependency map, WP01, WP02, WP03, source include map.
- All global 9号书 full-file includes listed in the initial review report.
- All Batch A exact line includes listed in the initial review report.
- Coordination files: resource map, coordination state, handoffs, review packages, and reports.

Additional revision files read:

- `coordination/prompts/004-batch-a-re-reviewer.md`
- `coordination/review-packages/004-batch-a-re-review.md`
- `coordination/reports/001-batch-a-builder-report.md`
- `coordination/reports/002-batch-a-review-report.md`
- `coordination/reports/003-batch-a-revision-1-builder-report.md`

Revised implementation files inspected directly:

- `internal/model/evidence.go`
- `internal/model/evidence_test.go`
- `internal/contract/schemas.go`
- `internal/contract/contract_test.go`
- `internal/platform/exec.go`
- `internal/app/command.go`
- `internal/app/command_test.go`
- `internal/app/runtime.go`
- `internal/cli/root.go`
- `justfile`
- root engineering files and selected Batch A source files.

## Evidence Commands

`git status --short`

Result: expected Batch A implementation paths and coordination paths only.

```text
 M .gitignore
?? .golangci.yml
?? .goreleaser.yaml
?? bookshelf/books/10-seekmoon-wbs-work-packages/coordination/
?? cmd/
?? go.mod
?? go.sum
?? internal/
?? justfile
```

Expanded status confirmed the untracked files are Batch A implementation files and coordination process files, including the re-review package and revision reports.

`go test ./...`

Result: passed.

```text
?   	github.com/yumiaura/seekmoon/cmd/seekmoon	[no test files]
ok  	github.com/yumiaura/seekmoon/internal/app	(cached)
ok  	github.com/yumiaura/seekmoon/internal/cli	(cached)
ok  	github.com/yumiaura/seekmoon/internal/contract	(cached)
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
ok  	github.com/yumiaura/seekmoon/internal/platform	(cached)
ok  	github.com/yumiaura/seekmoon/internal/store	(cached)
```

`go test ./internal/model ./internal/contract`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/model	(cached)
ok  	github.com/yumiaura/seekmoon/internal/contract	(cached)
```

`go test ./internal/platform ./internal/store ./internal/app`

Result: passed.

```text
ok  	github.com/yumiaura/seekmoon/internal/platform	(cached)
ok  	github.com/yumiaura/seekmoon/internal/store	(cached)
ok  	github.com/yumiaura/seekmoon/internal/app	(cached)
```

`go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/model ./internal/contract ./internal/platform ./internal/store ./internal/app ./internal/cli`

Result: passed. `internal/model` imports only standard library packages. `internal/platform` imports only standard library packages and no longer imports `internal/model`. `internal/app` now owns the platform-to-model command result mapping.

```text
github.com/yumiaura/seekmoon/internal/model encoding/json errors fmt strings time
github.com/yumiaura/seekmoon/internal/contract github.com/yumiaura/seekmoon/internal/model
github.com/yumiaura/seekmoon/internal/platform bytes context net/http os os/exec path/filepath time
github.com/yumiaura/seekmoon/internal/store context encoding/json github.com/adrg/xdg github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/platform path/filepath strings
github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/platform github.com/yumiaura/seekmoon/internal/store net/http time
github.com/yumiaura/seekmoon/internal/cli context fmt github.com/spf13/cobra github.com/yumiaura/seekmoon/internal/app io
```

`just fmt-check`

Result: passed.

```text
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal)"
```

`just mod-check`

Result: passed.

```text
go mod tidy
git diff --exit-code -- go.mod go.sum
go mod verify
all modules verified
```

`goreleaser check`

Result: passed.

```text
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
```

## Original Findings Re-Checked

### Evidence wrapper source preservation

Status: fixed.

Evidence:

- `internal/model/evidence.go:77` through `internal/model/evidence.go:82` defines `Evidence[T]` with `Source *string `json:"source"`` and optional `Error string `json:"error,omitempty"``.
- `internal/model/evidence.go:92` through `internal/model/evidence.go:93` leaves source nil for `Unknown[T]()`, which encodes as `"source": null`.
- `internal/model/evidence.go:121` through `internal/model/evidence.go:126` returns nil for empty source labels.
- `internal/contract/schemas.go:53` through `internal/contract/schemas.go:61` requires `status`, `value`, and `source`, and permits `source` as `string|null`.
- `internal/model/evidence_test.go:57` through `internal/model/evidence_test.go:72` covers source-less evidence encoding `"source":null`.

Judgment:

The evidence wrapper now preserves `status`, `value`, `source`, and optional `error` as required by WP02 and the data dictionary. The JSON contract can represent explicit source absence.

### Package-boundary direction

Status: fixed.

Evidence:

- `internal/platform/exec.go:3` through `internal/platform/exec.go:8` imports only standard library packages.
- `internal/platform/exec.go:17` through `internal/platform/exec.go:26` defines a platform-local `Runner` and `RunResult`.
- `internal/app/command.go:8` through `internal/app/command.go:20` maps `platform.RunResult` plus error to `model.CommandResult` and `model.State`.
- Required `go list` evidence confirms `internal/platform` no longer imports `internal/model`.

Judgment:

The platform package now owns host execution mechanics only. The canonical command-result mapping has moved to `internal/app`, which is an acceptable Batch A seam because it can depend on both host capability objects and canonical model objects.

### Quality gates

Status: fixed.

Evidence:

- `just fmt-check` passed.
- `just mod-check` passed.
- `goreleaser check` passed.

The revised `just fmt-check` scope is `cmd internal`. This is acceptable for Batch A because WP01's Go substrate lives in root engineering files, `cmd`, and `internal`, while the broader repository contains unrelated spike/probe cache material outside the Batch A write boundary. The scope still checks the actual Go implementation source introduced by Batch A.

## Positive Batch A Boundary Re-Check

- `cmd/seekmoon/main.go` remains process startup only.
- `internal/cli/root.go` contains placeholder command seams only and no Batch B-D business behavior.
- `internal/model` imports only Go standard library packages.
- State vocabulary remains exactly `present`, `missing`, `unknown`, `failed`, `unavailable`, and `derived`.
- Helpers preserve distinct meanings for `missing`, `unknown`, `failed`, and `unavailable`.
- Evidence wrappers preserve `status`, `value`, `source`, and optional `error`.
- `internal/contract` schemas are explicit public contract objects, not reflection-derived from structs.
- `internal/store` owns mechanics and path construction only.
- Project storage paths remain under `.seekmoon/`.
- Reusable cache paths remain under `$XDG_CACHE_HOME/seekmoon/`.
- Runtime construction provides Batch A platform/store substrate without implementing source readers, services, renderers, or public command behavior.
- Tests cover the required WP01-WP03 responsibilities well enough for Batch A promotion.

## Commit Hash

Pending before commit. The approval commit is created after this report is written.

## Required Follow-Up

No Batch A blocking follow-up remains.

Non-blocking notes for later batches:

- `just fmt-check` is intentionally scoped to Batch A Go implementation paths; if later work brings more Go source under new stable paths, update the formatting scope.
- `internal/app.CommandResultFromRun` is a Batch A compile seam. Later source/service work should decide whether Moon CLI command-result mapping belongs in `internal/source` or service-level command flow once those packages are introduced.
