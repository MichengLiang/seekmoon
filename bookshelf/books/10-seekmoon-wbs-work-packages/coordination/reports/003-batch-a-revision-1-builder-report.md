# Batch A Revision 1 Builder Report

## Summary

Addressed all independent review rejection findings from `002-batch-a-review-report.md` without committing.

Revision results:

- Evidence wrappers now preserve `source` as a required nullable JSON field.
- Unknown/source-less evidence now encodes `"source": null`.
- Evidence schema now requires `source` and accepts `string|null`.
- `internal/platform` no longer imports `internal/model`.
- `platform.Runner` now returns platform-local `platform.RunResult`.
- Mapping from platform execution result to `model.CommandResult` / `model.State` moved to `internal/app`.
- Installed pinned `gofumpt@v0.10.0` and `goreleaser@v2.16.0` with `go install`.
- Exposed installed tools through `~/.local/bin`, which is already on this session PATH.
- Required verification commands all pass.

## Review Findings Addressed

### 1. Evidence Wrapper Source Preservation

Reviewer finding:

- `Evidence[T].Source` was `string` with `json:"source,omitempty"`.
- Source-less evidence omitted the `source` field instead of preserving `"source": null`.
- Contract schema required only `status` and `value`, not `source`.

Revision:

- Changed `internal/model/evidence.go`:
  - `Evidence[T].Source` is now `*string` with `json:"source"`.
  - Constructors that receive a source set a pointer for non-empty source labels.
  - `Unknown[T]()` leaves `Source` nil, which encodes as JSON null.
  - `error` remains optional through `json:"error,omitempty"`.
- Updated tests in `internal/model/evidence_test.go`:
  - Round-trip test checks non-null source pointer.
  - Added `TestEvidenceWrapperPreservesNullSource`, checking JSON contains `"source":null`.
  - Updated failed-evidence validation setup for pointer source.
- Updated `internal/contract/schemas.go`:
  - Evidence schema now has required fields `status`, `value`, and `source`.
  - `source` type is `["string", "null"]`.

### 2. Remove `internal/platform -> internal/model`

Reviewer finding:

- `internal/platform/exec.go` imported `internal/model`.
- `platform.Runner.Run` returned `model.CommandResult`.
- Platform mapped execution status directly into `model.State`.

Revision:

- Changed `internal/platform/exec.go`:
  - Removed `internal/model` import.
  - Added platform-local `RunResult`.
  - `Runner.Run(ctx, request)` now returns `(platform.RunResult, error)`.
  - `ExecRunner` only records host execution fields: command, cwd, exit code, and log path.
- Added `internal/app/command.go`:
  - `CommandResultFromRun(platform.RunResult, error) model.CommandResult`.
  - Maps nil error to `model.StatePresent` and non-nil error to `model.StateFailed`.
  - This package is allowed to import both `internal/platform` and `internal/model`.
- Added `internal/app/command_test.go`:
  - Covers successful and failed platform-to-model command result mapping.
- Verified import boundary with required `go list` command:
  - `internal/platform` imports only standard library packages.

### 3. External Quality Gates

Reviewer finding:

- `just fmt-check` was blocked because `gofumpt` was missing.
- `goreleaser check` was blocked because `goreleaser` was missing.

Revision:

- Ran:
  - `go install mvdan.cc/gofumpt@v0.10.0`
  - `go install github.com/goreleaser/goreleaser/v2@v2.16.0`
- The binaries were installed into `/home/t103o/go/bin`, which was not on the active PATH.
- Added symlinks into `/home/t103o/.local/bin`, which is on PATH:
  - `/home/t103o/.local/bin/gofumpt -> /home/t103o/go/bin/gofumpt`
  - `/home/t103o/.local/bin/goreleaser -> /home/t103o/go/bin/goreleaser`
- Ran `gofumpt -w cmd internal`.
- Adjusted `justfile` format recipes from repository-wide `.` to Batch A implementation paths:
  - `fmt`: `gofumpt -w cmd internal`
  - `fmt-check`: `test -z "$(gofumpt -l cmd internal)"`
- Reason for scope: repository-wide `gofumpt -l .` reported unrelated cached/probe Go files under `spike/**/.mooncakes/**`, outside the Batch A write boundary. The revised recipe checks the actual Go implementation substrate paths owned by Batch A without rewriting unrelated spike cache material.
- `just fmt-check` now passes.
- `goreleaser check` now passes.

## Files Changed

Revision-specific changed files:

- `justfile`
- `internal/model/evidence.go`
- `internal/model/evidence_test.go`
- `internal/contract/schemas.go`
- `internal/platform/exec.go`
- `internal/app/command.go`
- `internal/app/command_test.go`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/003-batch-a-revision-1-builder-report.md`

Batch A changed path set still present in working tree:

- `.gitignore`
- `.golangci.yml`
- `.goreleaser.yaml`
- `go.mod`
- `go.sum`
- `justfile`
- `cmd/seekmoon/main.go`
- `internal/app/*`
- `internal/cli/*`
- `internal/contract/*`
- `internal/model/*`
- `internal/platform/*`
- `internal/store/*`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/001-batch-a-builder-report.md`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/003-batch-a-revision-1-builder-report.md`

Existing coordination files under `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/` remain untracked as process artifacts. This revision did not revert or modify unrelated work.

## Evidence Commands

### Tool Installation

`go install mvdan.cc/gofumpt@v0.10.0`

Result: passed.

Relevant output:

```text
go: downloading mvdan.cc/gofumpt v0.10.0
go: downloading golang.org/x/mod v0.35.0
go: downloading golang.org/x/tools v0.44.0
```

`go install github.com/goreleaser/goreleaser/v2@v2.16.0`

Result: passed.

Relevant output: large GoReleaser dependency download/compile completed with exit code 0.

Tool path/version confirmation:

```text
/home/t103o/.local/bin/gofumpt
v0.10.0 (go1.26.3)
/home/t103o/.local/bin/goreleaser
```

### Required Verification Commands

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

Result: passed.

```text
github.com/yumiaura/seekmoon/internal/model encoding/json errors fmt strings time
github.com/yumiaura/seekmoon/internal/contract github.com/yumiaura/seekmoon/internal/model
github.com/yumiaura/seekmoon/internal/platform bytes context net/http os os/exec path/filepath time
github.com/yumiaura/seekmoon/internal/store context encoding/json github.com/adrg/xdg github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/platform path/filepath strings
github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/platform github.com/yumiaura/seekmoon/internal/store net/http time
github.com/yumiaura/seekmoon/internal/cli context fmt github.com/spf13/cobra github.com/yumiaura/seekmoon/internal/app io
```

Boundary evidence:

- `internal/model` imports only standard library packages.
- `internal/platform` imports only standard library packages.
- `internal/app` now owns the platform-to-model command result mapping.

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

## Remaining Risks Or Blocked Gates

- No required revision gates are blocked.
- `just fmt-check` is intentionally scoped to `cmd internal` because unrelated spike/probe cache Go files under `spike/**/.mooncakes/**` are outside the Batch A write boundary and are not part of the Go implementation substrate.
- The full `coordination/` directory remains untracked process material from the coordinator flow; this revision added only the required revision report path.
- `go install` placed binaries under `/home/t103o/go/bin`; PATH exposure was provided by symlinks in `/home/t103o/.local/bin`.

## Suggested Reviewer Focus

- Confirm evidence JSON now always includes `source`, including source-less evidence as `null`.
- Confirm `internal/contract` evidence schema requires `source` and permits `string|null`.
- Confirm `internal/platform` no longer imports `internal/model`.
- Confirm `internal/app.CommandResultFromRun` is an acceptable Batch A mapping seam for platform execution results.
- Confirm `just fmt-check` scoped to `cmd internal` satisfies the WP01 quality gate for the SeekMoon Go implementation without touching unrelated spike cache files.
