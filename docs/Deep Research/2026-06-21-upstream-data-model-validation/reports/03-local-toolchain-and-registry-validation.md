# Task 03: Local Toolchain And Registry Validation

Date: 2026-06-21  
Probe root: `/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-03/`  
Report target: `/home/t103o/workbench/projects/seekmoon/docs/Deep Research/2026-06-21-upstream-data-model-validation/reports/03-local-toolchain-and-registry-validation.md`

## Summary

The local MoonBit toolchain is the v0.10.0-era 2026-06-08 toolchain:

```text
moon 0.1.20260608 (60bc8c3 2026-06-08)
moonc v0.10.0+e66899a54 (2026-06-09)
moonrun 0.1.20260608 (60bc8c3 2026-06-08)
Feature flags enabled: rr_moon_mod,rr_moon_pkg
```

`moon --help` now includes `runwasm` and `fetch`, and `moon ide --help` includes `ide doc`. It still does not include `search`, `view`, `audit`, or `outdated`.

`moon update` updates both local registry index and local symbols cache:

```text
Registry index updated successfully
Symbols updated successfully
```

The local registry index is JSON Lines under `~/.moon/registry/index/user/**/*.index`. Current count is 1,363 `.index` files and 9,952 JSONL records. The local symbols cache is JSON Lines under `~/.moon/registry/symbols/**/*.symbols`. Current count is 3 `.symbols` files and 1,619 JSONL records. Symbols cache is useful for cached API search but is not full-registry coverage.

`moon ide doc` searches exported APIs and documentation from the current module, `moonbitlang/core`, and local registry symbol indexes. It does not search all registry packages generally. For example, `moon ide doc markdown` returns no result even though the registry contains Markdown packages, while `moon ide doc sha256` finds `moonbitlang/x` because that module has a local symbol file.

`moon fetch mizchi/markdown@0.6.2` inside a temp project created `.repos/mizchi/markdown/0.6.2` and used the global zip cache at `~/.moon/registry/cache/mizchi/markdown/0.6.2.zip`.

`moon runwasm --help` documents local package inputs, accepted Mooncakes coordinate forms, latest-version resolution for unpinned coordinates, and the cache path `$MOON_HOME/registry/cache/assets`. The local cache currently contains `~/.moon/registry/cache/assets/Yoorkin/cowsay/0.1.0/cowsay.wasm`.

## Raw Probe Files

Raw command outputs and derived schema summaries were saved under:

```text
/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-03/outputs/
```

The temp MoonBit fetch probe project is:

```text
/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-03/probes/fetch_probe/
```

## Command Surface

### Versions

| Command | Result |
|---|---|
| `~/.moon/bin/moon --version` | `moon 0.1.20260608`; `moonc v0.10.0+e66899a54`; `moonrun 0.1.20260608`; feature flags `rr_moon_mod,rr_moon_pkg` |
| `~/.moon/bin/moonc --version` | Fails: `unknown option '--version'`; help says use `-v` |
| `~/.moon/bin/moonc -v` | `v0.10.0+e66899a54 (2026-06-09)` |
| `~/.moon/bin/moonrun --version` | `moonrun 0.1.20260608 (60bc8c3 2026-06-08)` |
| `~/.moon/bin/mooncake --version` | `mooncake-bin 0.1.20260521 (287155a 2026-05-21)` |

### `moon --help` Presence

| Command | Present in `moon --help` | Notes |
|---|---:|---|
| `search` | no | No package discovery command |
| `view` | no | No package detail command |
| `audit` | no | No local audit command |
| `outdated` | no | No outdated dependency command |
| `runwasm` | yes | "Run a local package as WebAssembly or a prebuilt WebAssembly binary" |
| `fetch` | yes | "Download a package to .repos directory (unstable)" |
| `ide doc` | yes | Under `moon ide`; "search exported APIs and documentation" |

`mooncake --help` remains narrow:

```text
login
register
publish
package
help
```

It does not expose search, view, audit, outdated, fetch, runwasm, or ide commands.

## `moon ide doc`

### Help Text Semantics

`moon ide doc --help` says:

```text
Search exported APIs and documentation from the current module, moonbitlang/core,
and registry symbol indexes.

Query examples:
  ''
  @json
  String
  String::length
  *parse*

Output: package lists, symbol summaries, or detailed documentation views.
```

Useful flags:

| Flag | Meaning |
|---|---|
| `--no-check` | Skip `moon check` |
| `--target` | Pass backend through to `moon check --target <backend>` |
| `--dump` | Dump symbol table |
| `--symbol-files` | Load additional symbol JSONL files; marked "for test only" |
| `--no-registry` | Exclude registry packages |

### Query Forms Observed

| Query | Result |
|---|---|
| `@moonbitlang/core/argparse` | Shows package-level exported types and methods |
| `@moonbitlang/core/argparse.Command` | Shows detailed `Command` type documentation and methods |
| `@moonbitlang/core/argparse.Command::parse` | Shows detailed method signature and docstring |
| `@json` | Shows `moonbitlang/core/json` package exports |
| `sha256` | Finds `moonbitlang/x/crypto` because `moonbitlang/x@0.4.45` is in symbols cache |
| `sha256 --no-registry` | No results |
| `@moonbitlang/x/crypto.sha256` | Finds cached registry symbol |
| `markdown` | No results |

### Behavior Conclusion

`moon ide doc` is an API/symbol documentation lookup tool over:

```text
current module
moonbitlang/core
registry symbol indexes already present under ~/.moon/registry/symbols
optional extra symbol files
```

It is not a general registry package search tool. A miss means only that the current module/core/symbol-cache search space did not contain the query. It does not prove the registry lacks a matching package.

For SeekMoon:

| Capability | Use `moon ide doc`? | Reason |
|---|---:|---|
| Known core package API | yes | Works directly |
| Known cached registry symbol | yes | Works if `.symbols` exists |
| General package discovery | no | `markdown` miss despite registry packages |
| Candidate API enrichment | yes, as supplemental source | Use after search/view when symbols are locally visible |
| Full registry API indexing | no | Must use Mooncakes assets or fetch/cache where available |

## Local Registry Index

### Paths And Counts

| Item | Value |
|---|---|
| Root | `~/.moon/registry/index` |
| User index root | `~/.moon/registry/index/user` |
| Git head after update | `7503ed87 (HEAD -> main, origin/main, origin/HEAD) update cybershang/agent-telemetry` |
| `.index` file count | 1,363 |
| JSONL record count | 9,952 |
| Parse errors | 0 |
| Directory size | 18M |

Each module has a path like:

```text
~/.moon/registry/index/user/<owner>/<module>.index
```

Each `.index` file is JSON Lines. Most files contain one line per published version, not one line per module.

Example files:

```text
~/.moon/registry/index/user/0Ayachi0/GB18030.index
~/.moon/registry/index/user/0xA672/moonbit_version_cli.index
```

### Registry Index Schema

The index has no single rigid schema. It is a JSONL stream of version records with common fields and historical spelling variants.

Most common schema shapes:

| Records | Keys |
|---:|---|
| 3,047 | `checksum, created_at, deps, description, keywords, license, name, readme, repository, version` |
| 1,181 | `checksum, created_at, deps, description, keywords, license, name, readme, repository, source, version` |
| 997 | `checksum, created_at, description, keywords, license, name, readme, repository, source, version` |
| 712 | `checksum, created_at, deps, description, keywords, license, name, preferred-target, readme, repository, source, version` |
| 698 | `checksum, created_at, description, keywords, license, name, readme, repository, version` |

Field occurrence and observed types:

| Field | Records | Observed types | Ingestion note |
|---|---:|---|---|
| `name` | 9,952 | string | Module coordinate |
| `version` | 9,952 | string | Version for this JSONL line |
| `checksum` | 9,952 | string | Published archive checksum |
| `license` | 9,853 | string | May be absent |
| `created_at` | 9,822 | string | Timestamp string |
| `repository` | 9,711 | string | May be empty string |
| `readme` | 9,316 | string | README path/name |
| `keywords` | 9,034 | array | May be absent |
| `description` | 8,955 | string | May be absent |
| `deps` | 7,061 | object | Dependency map |
| `source` | 4,262 | string | Source directory |
| `preferred-target` | 2,212 | string | Hyphen spelling |
| `exclude` | 619 | array | Packaging config |
| `warn-list` | 494 | string | Hyphen spelling |
| `preferred_target` | 475 | string | Underscore spelling |
| `include` | 314 | array | Packaging config |
| `--moonbit-unstable-prebuild` | 165 | string | Tooling/unstable flag-like field |
| `scripts` | 123 | object | Script map |
| `supported-targets` | 86 | string or array | Hyphen spelling; normalize |
| `bin-deps` | 70 | object | Binary dependencies |
| `homepage` | 65 | string | Optional |
| `alert-list` | 63 | string | Optional |
| `targets` | 43 | array | Target list |
| `supported_targets` | 37 | string or array | Underscore spelling; normalize |
| `authors` | 7 | array | Optional |
| `rule` | 6 | object | Optional |
| `import` | 5 | array | Optional |
| `author` | 3 | string | Optional |
| `preferred-backend` | 2 | string | Optional |
| `root-dir` | 2 | string | Optional |
| `keyword` | 2 | string | Singular spelling; normalize only if needed |
| `dependencies` | 1 | string | Legacy or anomalous; do not merge blindly with `deps` without shape check |
| `link` | 1 | object | Optional |
| `source-dir` | 1 | string | Optional |
| `deprecated` | 1 | string | Optional |

Example record:

```json
{
  "checksum": "6a801add78062f6c18a1c3801e9e8955236740099b2ae766c2aef28675bb875a",
  "created_at": "2025-08-29T19:55:38.462199+00:00",
  "description": "A comprehensive GB18030 encoding/decoding library for MoonBit",
  "keywords": ["gb18030", "encoding", "decoding", "chinese", "moonbit"],
  "license": "Apache-2.0",
  "name": "0Ayachi0/GB18030",
  "readme": "README.md",
  "repository": "https://github.com/0Ayachi0/GB18030",
  "source": "src",
  "version": "0.1.0"
}
```

### Ingestion Rules For SeekMoon

1. Treat `.index` files as JSONL, not JSON arrays.
2. Treat each line as one published version record.
3. Do not assume one schema.
4. Normalize hyphen/underscore variants into canonical SeekMoon fields while preserving raw keys.
5. Preserve unknown raw fields in a raw/debug projection, not default pretty output.
6. Treat empty strings as present-but-empty values, not equivalent to absent fields.
7. Keep `name` and `version` as required for an ingested version record; all other fields are optional in practice.
8. For current latest state, select the latest version by version semantics or registry/API latest source; do not assume last JSONL line is always sufficient without checking.

## Local Symbols Cache

### Paths And Counts

| Item | Value |
|---|---|
| Root | `~/.moon/registry/symbols` |
| File count | 3 |
| JSONL record count | 1,619 |
| Parse errors | 0 |
| Directory size | 580K |

Current files:

```text
~/.moon/registry/symbols/moonbitlang/async/0.19.1.symbols
~/.moon/registry/symbols/moonbitlang/quickcheck/0.14.0.symbols
~/.moon/registry/symbols/moonbitlang/x/0.4.45.symbols
```

### Symbols Schema

There are two record classes.

Meta record:

```json
{
  "kind": "meta",
  "module": "moonbitlang/async",
  "schema_version": 2,
  "version": "0.19.1"
}
```

Symbol record:

```json
{
  "attrs": [],
  "doc": "Compute the Sha256 digest in `Bytes` of some `data`. Note that Sha256 is big-endian.",
  "key": "moonbitlang/x/crypto.sha256",
  "kind": "function",
  "name": "sha256",
  "parent": "moonbitlang/x/crypto",
  "pkg": "moonbitlang/x/crypto",
  "sig_": "pub fn[Data : ByteSource] sha256(Data) -> FixedArray[Byte]"
}
```

Schema shapes:

| Records | Keys |
|---:|---|
| 1,616 | `attrs, doc, key, kind, name, parent, pkg, sig_` |
| 3 | `kind, module, schema_version, version` |

Field occurrence and observed types:

| Field | Records | Observed types | Ingestion note |
|---|---:|---|---|
| `kind` | 1,619 | string or array | `meta` record uses string; some symbol records use arrays for impl/alias forms |
| `key` | 1,616 | string | Stable symbol key |
| `pkg` | 1,616 | string | Package coordinate |
| `name` | 1,616 | string | Symbol name |
| `parent` | 1,616 | string or null | Parent type/package may be null |
| `sig_` | 1,616 | string or null | Signature; note trailing underscore in raw key |
| `doc` | 1,616 | string or null | Docstring |
| `attrs` | 1,616 | array | Attributes such as deprecated |
| `module` | 3 | string | Meta only |
| `version` | 3 | string | Meta only |
| `schema_version` | 3 | int | Meta only |

### Ingestion Rules For SeekMoon

1. Treat `.symbols` files as JSONL.
2. Read the first/meta record to bind module/version/schema version.
3. Treat symbol records as API/symbol search data.
4. Normalize `sig_` to canonical `signature` in SeekMoon output while preserving raw key in raw/debug output.
5. Accept `kind` as either string or array.
6. Do not interpret symbols cache as full registry coverage.

## `moon fetch`

### Probe Setup

Temp project:

```text
/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-03/probes/fetch_probe
```

Creation note:

```text
moon new fetch-probe
```

fails because project names may contain only alphanumeric characters and underscores. The successful project used:

```text
moon new fetch_probe
```

### Fetch Command

Run inside the temp project:

```bash
~/.moon/bin/moon fetch mizchi/markdown@0.6.2
```

Output summary:

```text
Registry index updated successfully
Symbols updated successfully
Fetching mizchi/markdown@0.6.2 to .../fetch_probe/.repos/mizchi/markdown/0.6.2
Using cached mizchi/markdown@0.6.2
Success: Successfully fetched mizchi/markdown@0.6.2 to .../fetch_probe/.repos/mizchi/markdown/0.6.2
```

Created path:

```text
/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-03/probes/fetch_probe/.repos/mizchi/markdown/0.6.2
```

Global cache used:

```text
~/.moon/registry/cache/mizchi/markdown/0.6.2.zip
```

Fetched source includes:

```text
CLAUDE.md
CONTRIBUTING.md
LICENSE
README.md
TODO.md
moon.mod.json
package.json
pnpm-lock.yaml
src/**/*.mbt
src/**/moon.pkg
src/**/pkg.generated.mbti
benches/
component/
docs/
e2e/
frontend/
js/
playground/
scripts/
worker/
```

Fetched `moon.mod.json`:

```json
{
  "name": "mizchi/markdown",
  "version": "0.6.2",
  "deps": {
    "moonbitlang/parser": "0.2.5",
    "moonbitlang/x": "0.4.41",
    "moonbitlang/yacc": "0.7.12",
    "mizchi/syntree": "0.2.3"
  },
  "readme": "README.md",
  "repository": "https://github.com/mizchi/markdown.mbt",
  "license": "MIT",
  "keywords": ["markdown", "parser", "cst", "incremental", "gfm"],
  "description": "Incremental Markdown parser and compiler",
  "source": "src",
  "preferred-target": "js"
}
```

### Behavior Semantics

| Behavior | Observed |
|---|---|
| Destination | Current project `.repos/<owner>/<module>/<version>` |
| Global cache | `~/.moon/registry/cache/<owner>/<module>/<version>.zip` |
| Registry update | Runs registry/symbol update before fetch |
| Existing cache | Uses cached zip when present |
| Project mutation | Creates `.repos` under the current project |
| Command stability | Help marks `fetch` as unstable |

For SeekMoon, `moon fetch` is a first-class source acquisition mechanism, but it must be treated as a project-local mutation and should run only inside an explicit temp/probe or user-approved project context.

## `moon runwasm`

### Help Semantics

`moon runwasm --help` says the command runs either:

1. a local package as WebAssembly, handled like `moon run --target wasm`; or
2. a prebuilt WebAssembly binary published as a Mooncakes asset.

Local package forms:

```text
moon runwasm main
moon runwasm ./main
```

Accepted Mooncakes coordinate forms:

```text
moon runwasm moonbitlang/parser/cmd/moonfmt@0.3.3
moon runwasm moonbitlang/parser@0.3.3/cmd/moonfmt
moon runwasm moonbitlang/parser/cmd/moonfmt
```

Coordinate rule:

| Form | Meaning |
|---|---|
| `user/module/package@version` | Pinned version at the end |
| `user/module@version/package` | Pinned module version before package path |
| `user/module/package` | Unpinned; resolves latest version from registry index |

Cache rule from help:

```text
Fetched wasm files are cached under $MOON_HOME/registry/cache/assets and reused on later runs.
```

For the current environment, `$MOON_HOME` is effectively `~/.moon`.

Current assets cache sample:

```text
~/.moon/registry/cache/assets/Yoorkin/cowsay/0.1.0/.moon-lock
~/.moon/registry/cache/assets/Yoorkin/cowsay/0.1.0/cowsay.wasm
```

`cowsay.wasm` size:

```text
22534 bytes
```

### Failure Semantics

Invalid coordinate:

```bash
moon runwasm not-a-valid-coordinate
```

Output:

```text
Error: Invalid runwasm coordinate `not-a-valid-coordinate`

Caused by:
    must be in format `user/module/package`
```

Dry-run behavior for Mooncakes assets:

```bash
moon runwasm --dry-run moonbitlang/parser/cmd/moonfmt@0.3.3 -- --help
```

Output:

```text
Error: --dry-run is not supported for Mooncakes assets in `moon runwasm`
```

### SeekMoon Implications

1. `runwasm` candidates must be modeled as executable assets or skills, not normal library dependencies.
2. Pinned coordinate storage is required for reproducible records.
3. Unpinned coordinates imply registry latest resolution and should be treated as non-reproducible unless resolved and recorded.
4. Cache inspection belongs under `~/.moon/registry/cache/assets`.
5. `--dry-run` cannot be used to preview remote asset execution behavior.

## Cache Paths

| Path | Role | Current observed state |
|---|---|---|
| `~/.moon/registry/index` | Local registry index git repo | 18M; head `7503ed87` |
| `~/.moon/registry/index/user/**/*.index` | Module version JSONL records | 1,363 files; 9,952 records |
| `~/.moon/registry/symbols/**/*.symbols` | Local API/symbol JSONL cache | 3 files; 1,619 records |
| `~/.moon/registry/cache/<owner>/<module>/<version>.zip` | Downloaded source zip cache | Example: `mizchi/markdown/0.6.2.zip` |
| `~/.moon/registry/cache/assets/<owner>/<module>/<version>/...` | `runwasm` Wasm asset cache | Example: `Yoorkin/cowsay/0.1.0/cowsay.wasm` |
| `~/.moon/lib/core/**` | Local core library source | Used by `moon ide doc` core queries and source inspection |
| `<project>/.repos/<owner>/<module>/<version>` | `moon fetch` extracted source destination | Created in temp probe |

## Assessment Of Existing SOP / Design Assumptions

| Assumption | Status | Correction |
|---|---|---|
| Local toolchain may still be 2026-05-22 v0.9.3 | stale | Current local toolchain is 2026-06-08 `moon` and 2026-06-09 `moonc v0.10.0+e66899a54` |
| `moon --help` lacks `search`, `view`, `audit`, `outdated` | current | Still absent |
| `moon --help` lacks `runwasm` | stale | `runwasm` is present |
| `moon --help` includes `fetch` | current | Present and marked unstable |
| `moon ide doc` exists and should replace deprecated `moon doc [SYMBOL]` for symbol lookup | current | `moon ide doc --help` confirms API/doc search |
| `moon ide doc` is a general registry package search | wrong | It searches current module, core, and local registry symbol indexes only |
| `moon ide doc` can enrich known cached package APIs | current | Works for `moonbitlang/x` symbols cache |
| `moon update` only updates registry index | incomplete | It also updates symbols cache |
| Local registry index is JSONL | current | Confirmed |
| Local registry index has one record per module | wrong | It has one record per version line within module `.index` files |
| Local symbols cache covers all registry packages | wrong | Current local cache has only 3 module/version files |
| `moon fetch` gets published source to `.repos` | current | Confirmed with `mizchi/markdown@0.6.2` |
| GitHub should be the primary published-source source | incomplete/wrong | Published source is available through `moon fetch` and cache/source zip; GitHub is maintenance/unpublished-source signal |
| `runwasm` cache path is `$MOON_HOME/registry/cache/assets` | current | Confirmed in help and local cache |
| `runwasm --dry-run` can preview remote asset execution | wrong | It errors for Mooncakes assets |

## Canonical Data Dictionary For SeekMoon Local Ingestion

### Module Version Record

Canonical object: `ModuleVersionRecord`

| Canonical field | Sources | Required | Notes |
|---|---|---:|---|
| `module` | index `name`, modules API `name`, manifest `module/name` | yes | Owner/module coordinate |
| `version` | index `version`, manifest `version` | yes | Version for this record |
| `checksum` | index `checksum`, manifest metadata checksum | no | Preserve as string |
| `created_at` | index `created_at`, manifest metadata created_at | no | Timestamp string |
| `description` | index/modules/manifest metadata | no | Missing is common |
| `keywords` | index/modules/manifest metadata | no | Array; normalize singular `keyword` only in raw-aware import |
| `license` | index/modules/manifest metadata | no | Empty string and missing are distinct raw states |
| `repository` | index/modules/manifest metadata | no | Empty string and missing are distinct raw states |
| `readme` | index `readme` | no | Path/name |
| `source_dir` | index `source`, `source-dir` | no | Normalize canonical value; preserve raw key |
| `deps` | index `deps`, fetched `moon.mod.json` deps | no | Prefer object shape |
| `preferred_target` | `preferred-target`, `preferred_target`, `preferred-backend` | no | Normalize spelling; preserve source |
| `supported_targets` | `supported-targets`, `supported_targets`, `targets` | no | Accept string or array; normalize to array plus raw |
| `include` | index `include` | no | Packaging config |
| `exclude` | index `exclude` | no | Packaging config |
| `scripts` | index `scripts` | no | Object |
| `raw` | index line/API payload | yes for debug | Raw source payload |
| `source_label` | ingestion layer | yes | `registry_index`, `modules_api`, `manifest_api`, `fetched_source` |
| `fetched_at` | ingestion layer | yes | Ingestion timestamp |

### Symbol Record

Canonical object: `SymbolRecord`

| Canonical field | Source field | Required | Notes |
|---|---|---:|---|
| `module` | meta `module` or file path | yes | Module coordinate |
| `version` | meta `version` or file path | yes | Symbol cache version |
| `schema_version` | meta `schema_version` | no | Usually 2 |
| `key` | `key` | yes for symbol | Stable symbol key |
| `package` | `pkg` | yes for symbol | Package coordinate |
| `name` | `name` | yes for symbol | Symbol name |
| `parent` | `parent` | no | May be null |
| `kind` | `kind` | yes | String or array in raw; canonical should allow list |
| `signature` | `sig_` | no | Normalize name |
| `doc` | `doc` | no | String or null |
| `attrs` | `attrs` | no | Array |
| `source_label` | ingestion layer | yes | `symbols_cache` |

### Fetched Source Record

Canonical object: `FetchedSourceRecord`

| Canonical field | Source | Required | Notes |
|---|---|---:|---|
| `module` | fetch coordinate | yes | Owner/module |
| `version` | fetch coordinate | yes | Pinned version |
| `project_path` | runtime | yes | Project where `.repos` was created |
| `dest_path` | `moon fetch` output | yes | `.repos/<owner>/<module>/<version>` |
| `cache_zip` | registry cache | no | `~/.moon/registry/cache/<owner>/<module>/<version>.zip` |
| `status` | command result | yes | `ok` or error |
| `files_sample` | filesystem scan | no | For report/debug |
| `moon_mod` | fetched `moon.mod.json` | no | Parsed when present |
| `command_output` | raw output | yes for debug | Preserve |

### Runwasm Asset Record

Canonical object: `RunwasmAssetRecord`

| Canonical field | Source | Required | Notes |
|---|---|---:|---|
| `coordinate` | user input / skills API | yes | Preserve original |
| `module` | parsed coordinate / skills API | yes | Owner/module |
| `version` | parsed or resolved | yes for reproducible record | Unpinned input must resolve and store version |
| `package` | parsed coordinate / skills API | yes | Package path |
| `wasm_cache_path` | cache inspection | no | Under `~/.moon/registry/cache/assets` |
| `wasm_url` | skills API | no | Remote asset URL |
| `checksum_url` | skills API | no | Remote checksum URL |
| `checksum` | checksum asset/cache | no | Fetch separately if needed |
| `status` | command/cache result | yes | `cached`, `not_cached`, `error`, etc. |
| `coordinate_form` | parser | yes | pinned-at-end, pinned-module, unpinned |

## Required Support Mapping

This is the minimum local model SeekMoon needs for the named support surfaces.

| SeekMoon surface | Required local facts |
|---|---|
| `doctor` | Tool versions; command presence; registry/symbol/cache paths; path existence; counts; `moon update` behavior |
| `sync` | `moon update`; registry index git head/counts; symbols count; cache path inventory |
| `api` | `moon ide doc`; symbols JSONL; Mooncakes package data assets when available; fetched source fallback |
| `source` | `moon fetch`; `.repos` destination; registry zip cache; source zip fallback; core source path |
| `probe` | Temp project creation constraints; `moon fetch` mutation boundary; `moon add/check/test/build` context |
| `skill` | `moon runwasm` coordinate parser; assets cache path; Skills API/wasm/checksum fields; dry-run limitation |

## Final Conclusions

1. The local toolchain has moved to the 2026-06-08 / v0.10.0 command surface.
2. Official CLI package discovery is still absent: no `moon search`, `moon view`, `moon audit`, or `moon outdated`.
3. `runwasm`, `fetch`, and `ide doc` are real local surfaces and must be modeled by SeekMoon.
4. Local registry index is a JSONL version-record store, not a simple module list.
5. Local symbols cache is a JSONL API/symbol store, but current local coverage is only 3 modules.
6. `moon ide doc` is an API documentation lookup over local/core/cached symbols, not a general package search.
7. `moon fetch` is the right local command for published source acquisition, provided SeekMoon confines it to explicit temp/probe/project destinations.
8. `moon runwasm` should be modeled as an executable asset workflow with pinned coordinates, cache inspection under `~/.moon/registry/cache/assets`, and no reliance on `--dry-run` for Mooncakes assets.
