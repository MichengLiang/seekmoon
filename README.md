# SeekMoon

[![CI](https://github.com/MichengLiang/seekmoon/actions/workflows/ci.yml/badge.svg)](https://github.com/MichengLiang/seekmoon/actions/workflows/ci.yml)
[![Pages](https://github.com/MichengLiang/seekmoon/actions/workflows/pages.yml/badge.svg)](https://github.com/MichengLiang/seekmoon/actions/workflows/pages.yml)

SeekMoon is a MoonBit package discovery workbench and the public bookshelf that
defines its evidence model, command contracts, and implementation plan.

The Go CLI reads Mooncakes, MoonBit local toolchain state, project context, and
GitHub repository evidence. It keeps upstream facts separate from derived local
evidence so package discovery, inspection, comparison, probing, adoption
records, and reports can share the same source-state vocabulary.

The bookshelf publishes the surrounding research and design material for
《包复用生态：发现、管理与评价尺度》.

## Repository

- GitHub: <https://github.com/MichengLiang/seekmoon>
- Go module: `github.com/MichengLiang/seekmoon`
- Published bookshelf: <https://michengliang.github.io/seekmoon/>
- License: Apache-2.0

## Contents

```text
cmd/seekmoon/       CLI entrypoint
internal/           Go implementation packages
tests/              black-box, journey, acceptance, and opt-in integration tests
bookshelf/          AsciiDoc bookshelf source and build workspace
docs/               research notes and validation reports
spike/              exploratory MoonBit and CLI probes
justfile            local quality-gate entrypoints
```

## CLI

Run the CLI from source:

```bash
go run ./cmd/seekmoon --help
```

Build a local binary:

```bash
go build -o seekmoon ./cmd/seekmoon
./seekmoon --help
```

Available command groups:

```text
doctor       Check local MoonBit and project environment evidence
sync         Create a local source snapshot
search       Search library module candidates
view         View a library module profile
api          View a package API profile
source       Locate published source
skill        Search or view executable skill entries
compare      Compare candidate evidence
probe        Run local validation for a candidate
record       Save an adoption judgment
report       Render an investigation report
raw          Read an upstream source payload without normalization
```

Every public command supports the common output modes:

```text
--json       render machine-readable JSON
--jq <expr>  evaluate a jq expression against JSON output
--shape      render the stable field shape
--schema     render the JSON Schema contract
```

Examples:

```bash
go run ./cmd/seekmoon search argparse --json
go run ./cmd/seekmoon search --shape
go run ./cmd/seekmoon search --schema
```

Integration commands that use network, GitHub, or local Moon CLI mutation are
opt-in in the test suite. Default tests run without external credentials or
service availability.

## Development

SeekMoon uses Go for the CLI and pnpm for the bookshelf build.

Required local tools for the full Go gate:

- Go 1.26.x
- `just`
- `gofumpt`
- `golangci-lint`
- `gotestsum`
- `govulncheck`
- `goreleaser`

Install Go dependencies:

```bash
go mod download
```

Run the main checks:

```bash
just fmt-check
just lint
just test
just test-race
just cover
just vuln
just mod-check
just release-check
```

Run the full local gate:

```bash
PATH="$(go env GOPATH)/bin:$PATH" just ci
```

`just cover` writes `.artifacts/coverage.out`. The file is generated output and
is ignored by Git.

The GitHub `CI` workflow validates the bookshelf workspace. The Go quality gate
is represented by the local `just ci` target.

## Bookshelf

The source catalog is [bookshelf/catalog.adoc](bookshelf/catalog.adoc).

Build the bookshelf locally:

```bash
cd bookshelf
pnpm install
pnpm run check
pnpm run build
```

Generated HTML is written to `bookshelf/build/html/`.

GitHub Pages is built from GitHub Actions. The Pages workflow installs the
bookshelf dependencies, runs structural checks, builds the HTML output, and
deploys `bookshelf/build/html`.

## Release

Release configuration lives in [.goreleaser.yaml](.goreleaser.yaml).

Create a local snapshot release:

```bash
just release-snapshot
```

The release build targets Linux, macOS, and Windows on `amd64` and `arm64`.

## License

Apache-2.0. See [LICENSE](LICENSE).
