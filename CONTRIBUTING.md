# Contributing to SeekMoon

SeekMoon contains two maintained surfaces:

- a Go CLI for MoonBit package discovery workbench behavior
- an AsciiDoc bookshelf for the package-reuse and SeekMoon specifications

Contributions should preserve the boundary between upstream facts, local derived
evidence, command contracts, output contracts, and research notes.

## Before Opening a Change

1. Read [README.md](README.md) or [README.zh-CN.md](README.zh-CN.md) for the
   project position and command surface.
2. Use the maintained bookshelf for specifications:
   - [Package Reuse Ecosystem](https://michengliang.github.io/seekmoon/books/08-package-reuse-ecosystem/book.html)
   - [SeekMoon: MoonBit Package Discovery Workbench](https://michengliang.github.io/seekmoon/books/09-seekmoon-cli-discovery-workbench/book.html)
3. Treat `docs/` as research notes, raw investigation material, and validation
   reports. Do not use raw notes as the public contract when a bookshelf section
   defines the same object.

## Change Types

| Area | Expected basis |
| --- | --- |
| CLI command behavior | Command contract, model type, service flow, output contract, and tests. |
| Output shape or schema | Contract definitions, JSON/schema tests, and compatibility impact. |
| Evidence source behavior | Source boundary, state vocabulary, and failure surface. |
| Help text or README | Object-language wording, command boundary, and current behavior. |
| Bookshelf content | Existing AsciiDoc structure, stable IDs, xrefs, and relation wording. |

## Development Setup

Install Go module dependencies:

```bash
go mod download
```

Build or run the CLI:

```bash
go run ./cmd/seekmoon --help
go build -o seekmoon ./cmd/seekmoon
```

Build the bookshelf:

```bash
cd bookshelf
pnpm install
pnpm run check
pnpm run build
```

## Quality Checks

Run the focused check for the changed area. For Go changes, run:

```bash
go test ./...
```

For the full local Go gate:

```bash
PATH="$(go env GOPATH)/bin:$PATH" just ci
```

For bookshelf changes:

```bash
cd bookshelf
pnpm run check
pnpm run build
```

## Pull Requests

Pull requests should include:

- a short description of the changed object
- the files or packages changed
- tests or checks run
- any contract, output, schema, or documentation impact

Keep generated artifacts out of commits unless a maintained build output is
explicitly part of the change.

## Commit Scope

Keep unrelated work out of the same commit. This repository contains source,
bookshelf, raw notes, and exploratory spikes; changes should stay scoped to the
object being modified.
