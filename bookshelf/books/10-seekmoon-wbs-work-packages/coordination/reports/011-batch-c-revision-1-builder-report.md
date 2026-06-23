# Batch C Revision 1 Builder Report

## Summary

Addressed the two independent review rejection findings for Batch C:

- Contract projection modes now render `--shape` and `--schema` directly from command contract definitions without requiring ordinary command operands and without invoking pending services.
- The real `cmd/seekmoon` process entrypoint now uses the same `ExecuteWithCode` exit mapping as the CLI package execution path.

No commit was created.

## Review findings addressed

### Contract projections no longer require operands or services

Added a shared CLI contract-mode path:

- `contractProjectionRequested` detects `--shape` / `--schema` during Cobra positional argument validation.
- `argsUnlessContract` bypasses ordinary positional operand validation for contract projection modes.
- `renderContractIfRequested` renders the command schema/shape immediately after output-mode resolution and returns before any service input parsing or service call.

The shared path is applied to all Batch C output-capable command surfaces:

- `doctor`
- `sync`
- `search`
- `view`
- `api`
- `source`
- `skill search`
- `skill view`
- `compare`
- `probe`
- `record`
- `report`
- `raw`

The previously failing probes now succeed:

```text
$ go run ./cmd/seekmoon search --shape
seekmoon.search-results.v1

schema: string
snapshot: object
  id: string
  sources: string[]
query: object
  text: string
  kind: library|skill
  target: string|null
results[]: array
  rank: int
  module: string
  version: string
  description: evidence<string>
  license: evidence<string>
  repository: evidence<string>
  target: evidence<object>
  match: object
  snapshot_id: string
```

```text
$ go run ./cmd/seekmoon api --shape
seekmoon.package-data.v1

schema: string
name: object|string
traits: object|string
errors: object|string
types: object|string
typealias: object|string
values: object|string
misc: object|string
```

Added `TestContractProjectionModesBypassOperandsAndPendingServices`, covering normally operand-required commands and checking that contract output contains the expected schema ID while excluding pending-service and argument-validation messages.

### Real process exit mapping now matches CLI mapping

`cmd/seekmoon/main.go` now exits with:

```go
os.Exit(cli.ExecuteWithCode(context.Background(), rt, cli.Options{Out: os.Stdout, Err: os.Stderr}, os.Args[1:]...))
```

Runtime construction failures still print to stderr and exit `1`. Normal command execution now uses the existing CLI mapper:

- success -> `0`
- service/projection failure -> `1`
- required argument / flag parse failure -> `2`

Focused binary probe:

```text
$ mkdir -p tmp/batch-c-revision && go build -o tmp/batch-c-revision/seekmoon ./cmd/seekmoon && set +e
$ ./tmp/batch-c-revision/seekmoon search >tmp/batch-c-revision/search-no-args.out 2>&1
$ code=$?
$ printf 'exit=%s\n' "$code"
exit=2
$ bat tmp/batch-c-revision/search-no-args.out
search requires a query
```

## Files changed

Revision-specific changed paths:

- `cmd/seekmoon/main.go`
- `internal/cli/api.go`
- `internal/cli/compare.go`
- `internal/cli/doctor.go`
- `internal/cli/flags.go`
- `internal/cli/probe.go`
- `internal/cli/raw.go`
- `internal/cli/record.go`
- `internal/cli/report.go`
- `internal/cli/root_test.go`
- `internal/cli/search.go`
- `internal/cli/skill.go`
- `internal/cli/source.go`
- `internal/cli/sync.go`
- `internal/cli/view.go`
- `internal/contract/schemas.go`
- `internal/contract/shapes.go`
- `internal/model/schema.go`
- `internal/output/schema.go`
- `internal/output/shape.go`
- `bookshelf/books/10-seekmoon-wbs-work-packages/coordination/reports/011-batch-c-revision-1-builder-report.md`

Existing Batch C implementation paths remain uncommitted from the prior builder pass; this revision did not revert them.

## Evidence commands

### `go test ./...`

Result: passed.

```text
?    github.com/yumiaura/seekmoon/cmd/seekmoon [no test files]
ok   github.com/yumiaura/seekmoon/internal/app 0.004s
ok   github.com/yumiaura/seekmoon/internal/cli (cached)
ok   github.com/yumiaura/seekmoon/internal/contract (cached)
ok   github.com/yumiaura/seekmoon/internal/model (cached)
ok   github.com/yumiaura/seekmoon/internal/output (cached)
ok   github.com/yumiaura/seekmoon/internal/platform (cached)
ok   github.com/yumiaura/seekmoon/internal/service (cached)
ok   github.com/yumiaura/seekmoon/internal/source (cached)
ok   github.com/yumiaura/seekmoon/internal/store (cached)
?    github.com/yumiaura/seekmoon/internal/testutil [no test files]
```

### `go test ./internal/output ./internal/contract`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/output 0.005s
ok   github.com/yumiaura/seekmoon/internal/contract 0.009s
```

### `go test ./internal/cli`

Result: passed.

```text
ok   github.com/yumiaura/seekmoon/internal/cli (cached)
```

### `go run ./cmd/seekmoon search --shape`

Result: passed; rendered contract shape and did not call the pending search service.

### `go run ./cmd/seekmoon api --shape`

Result: passed; rendered API contract shape without requiring the normal module operand or `--package`.

### Focused binary exit-code probe

Command:

```bash
mkdir -p tmp/batch-c-revision && go build -o tmp/batch-c-revision/seekmoon ./cmd/seekmoon && set +e
./tmp/batch-c-revision/seekmoon search >tmp/batch-c-revision/search-no-args.out 2>&1
code=$?
printf 'exit=%s\n' "$code"
bat tmp/batch-c-revision/search-no-args.out
```

Result:

```text
exit=2
search requires a query
```

### Import boundary listing

Command:

```bash
go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/output ./internal/cli ./internal/service ./internal/app
```

Result:

```text
github.com/yumiaura/seekmoon/internal/output context encoding/json errors fmt github.com/itchyny/gojq github.com/yumiaura/seekmoon/internal/contract github.com/yumiaura/seekmoon/internal/model io strings
github.com/yumiaura/seekmoon/internal/cli context fmt github.com/spf13/cobra github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/output github.com/yumiaura/seekmoon/internal/service io strconv strings
github.com/yumiaura/seekmoon/internal/service context fmt github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/source github.com/yumiaura/seekmoon/internal/store time
github.com/yumiaura/seekmoon/internal/app github.com/yumiaura/seekmoon/internal/model github.com/yumiaura/seekmoon/internal/output github.com/yumiaura/seekmoon/internal/platform github.com/yumiaura/seekmoon/internal/service github.com/yumiaura/seekmoon/internal/source github.com/yumiaura/seekmoon/internal/store net/http time
```

### `just fmt-check`

Result: passed.

```text
command -v gofumpt >/dev/null
test -z "$(gofumpt -l cmd internal)"
```

### `go mod tidy`

Result: passed with no command output.

### `go mod verify`

Result: passed.

```text
all modules verified
```

### `goreleaser check`

Result: passed.

```text
• checking                                  path=.goreleaser.yaml
• 1 configuration file(s) validated
• thanks for using GoReleaser!
```

## Remaining risks or blocked gates

No blocked verification gates remain for the two review findings.

Residual scope notes:

- The new `doctor`, `sync`, `compare`, and `raw` schema IDs use generic Batch C contract definitions. This is intentional because concrete Batch D service behavior and full final object semantics remain outside Batch C.
- `internal/service` still imports `internal/source` through the pre-existing Batch B sync runner path. This is unchanged from the prior Batch C review and is visible in the import boundary listing.
- The worktree contains broader uncommitted Batch C and coordination files from the active multi-executor run. This revision did not revert or commit them.

## Suggested reviewer focus

- Confirm `--shape` and `--schema` bypass ordinary operands and pending services for all output-capable commands, especially `api`, `view`, `source`, `probe`, `record`, and `raw`.
- Confirm `cmd/seekmoon` now uses `ExecuteWithCode` and that a directly built binary exits `2` for parse/argument failures.
- Confirm the generic schema additions for `doctor`, `sync`, `compare`, and `raw` are acceptable as Batch C contract placeholders pending Batch D service implementation.
