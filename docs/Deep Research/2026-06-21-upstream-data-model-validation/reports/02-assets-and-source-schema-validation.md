# Task 02: Assets And Published Source Schema Validation

## Summary

Live fetch window: 2026-06-21T12:10:54Z to approximately 2026-06-21T12:14:08Z.

Raw files and analysis artifacts:

`/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-02/`

Main findings:

- `module_index.json` is available for all representative modules inspected.
- The tree key is spelled `childs`. No observed representative `module_index.json` used `children`.
- `module_index.json` package paths are full package paths such as `moonbitlang/core/argparse`; asset URLs require deriving the package relpath by removing the module path prefix.
- `package_data.json` is available for representative packages and uses top-level keys `name`, `traits`, `errors`, `types`, `typealias`, `values`, and `misc`.
- Type entries contain `name`, `docstring`, `signature`, `loc`, `methods`, and `impls`. Value entries contain `name`, `docstring`, `signature`, and `loc`. Method entries contain `name`, `docstring`, `signature`, and `loc`.
- `resources.json` returned 404 for every package checked in this task, including packages with valid `package_data.json` and source README files. A bounded scan made 297 `resources.json` attempts and found no 200 response. A 404 must mean only that this optional resource asset is absent, not that the package or docs are absent.
- This task intentionally validated the plural `resources.json` path. A later follow-up validated the singular `resource.json` path used by current frontend source and found it live for many modules. Keep the plural and singular paths separate in the canonical model.
- Source zip is available for `mizchi/markdown`, `moonbit-community/cmark`, and the sampled skill module `jaredzhou/pony`. It is not available for `moonbitlang/core` latest manifest version at the documented URL shape.
- Source zip contents are sufficient published-source evidence when available: observed zips include module config, package config, `.mbt` source files, README, tests/examples/benches in some modules, generated `.mbti`, and license files. They are not universal, so ingestion must support fallback paths.

Registry snapshot during fetch:

| source | value |
|---|---:|
| `/api/v0/modules` length | 1350 |
| `statistics.total_modules` | 1350 |
| `statistics.total_packages` | 12008 |
| `statistics.total_lines` | 44597542 |
| `statistics.total_downloads` | 4043944 |
| `/api/v0/skills` length | 70 |

## Representative Samples

| module | manifest version | sampled package | why sampled |
|---|---|---|---|
| `moonbitlang/core` | `0.1.20260609+84519ca0a` | `moonbitlang/core/argparse` | required core latest and argparse package data |
| `mizchi/markdown` | `0.6.2` | `mizchi/markdown` | required markdown module and at least one package |
| `moonbit-community/cmark` | `0.4.4` | `moonbit-community/cmark/char` | required cmark module and package |
| `jaredzhou/pony` | `0.1.1` | `jaredzhou/pony/cmd/main` | first observed Skills API entry with Wasm executable URLs |

## Asset URL Construction Rules

Use the manifest endpoint first:

```text
GET https://mooncakes.io/api/v0/manifest/<owner>/<module-relpath>
```

The manifest `version` is the version segment used by assets.

Module index:

```text
GET https://mooncakes.io/assets/<owner>/<module-relpath>@<version>/module_index.json
```

Package data:

```text
GET https://mooncakes.io/assets/<owner>/<module-relpath>@<version>/<pkg-relpath>/package_data.json
```

Package resources:

```text
GET https://mooncakes.io/assets/<owner>/<module-relpath>@<version>/<pkg-relpath>/resources.json
```

Source zip:

```text
GET https://download.mooncakes.io/user/<owner>/<module-relpath>/<version>.zip
```

The package path in `module_index.json` is not the same as the package asset relpath:

| value | example |
|---|---|
| module path | `moonbitlang/core` |
| package path from `module_index.json` | `moonbitlang/core/argparse` |
| package relpath for asset URL | `argparse` |
| package data URL | `https://mooncakes.io/assets/moonbitlang/core@0.1.20260609+84519ca0a/argparse/package_data.json` |

For a root package whose package path equals the module path, the package relpath is empty. The observed root package data URL for `mizchi/markdown` is:

```text
https://mooncakes.io/assets/mizchi/markdown@0.6.2/package_data.json
```

Do not derive package paths by guessing filesystem layout from source zip. Use `module_index.json` as the package list authority, then strip the module path prefix to construct the asset relpath.

## `module_index.json` Schema

Observed top-level keys:

| key | observed type | meaning |
|---|---|---|
| `name` | string | tree node display/name segment |
| `package` | object or null | package summary when this node is a package; null for grouping/root nodes |
| `childs` | array | child tree nodes |

Observed spelling:

| key | observed |
|---|---|
| `childs` | yes |
| `children` | no |

Observed package object keys:

| key | observed type | meaning |
|---|---|---|
| `path` | string | full package path, e.g. `moonbitlang/core/argparse` |
| `traits` | array | trait names or summaries |
| `errors` | array | error names or summaries |
| `types` | array | type summary objects |
| `typealias` | array | type alias names or summaries |
| `values` | array | value/function names or summaries |
| `misc` | array | miscellaneous exported API summaries |

Observed type summary object in module index:

```json
{
  "name": "Command",
  "impls": [],
  "methods": ["new", "parse", "render_help"]
}
```

Observed package summaries:

| module | package count | first package paths |
|---|---:|---|
| `moonbitlang/core` | 61 | `moonbitlang/core/abort`, `moonbitlang/core/argparse`, `moonbitlang/core/array`, `moonbitlang/core/bench`, `moonbitlang/core/bigint` |
| `mizchi/markdown` | 22 | `mizchi/markdown`, `mizchi/markdown/api`, `mizchi/markdown/experimental/crdt`, `mizchi/markdown/experimental/mdx`, `mizchi/markdown/experimental/multipass` |
| `moonbit-community/cmark` | 5 | `moonbit-community/cmark/char`, `moonbit-community/cmark/cmark`, `moonbit-community/cmark/cmark_base`, `moonbit-community/cmark/cmark_html`, `moonbit-community/cmark/cmark_renderer` |
| `jaredzhou/pony` | 3 | `jaredzhou/pony`, `jaredzhou/pony/cmd/main`, `jaredzhou/pony/mw` |

SeekMoon ingestion rule:

1. Traverse recursively through `childs`.
2. Treat nodes with non-null `package` as package summary nodes.
3. Use `package.path` as the canonical package id.
4. Derive the package asset relpath by removing `<module>/` from `package.path`; if equal to module path, use empty relpath.
5. Use package data for detailed API fields; module index is only a compact package/API index.

## `package_data.json` Schema

Observed top-level keys across representative packages:

| key | observed type | observed in samples | meaning |
|---|---|---|---|
| `name` | string | yes | package name/path label |
| `traits` | array | yes | detailed trait entries |
| `errors` | array | yes | detailed error entries |
| `types` | array | yes | detailed type entries |
| `typealias` | array | yes | detailed type alias entries |
| `values` | array | yes | detailed value/function entries |
| `misc` | array | yes | miscellaneous entries |

Observed counts:

| package | traits | errors | types | typealias | values | misc |
|---|---:|---:|---:|---:|---:|---:|
| `moonbitlang/core/argparse` | 0 | 0 | 10 | 0 | 0 | 0 |
| `mizchi/markdown` | 0 | 0 | 22 | 0 | 21 | 0 |
| `moonbit-community/cmark/char` | 0 | 0 | 0 | 0 | 22 | 0 |
| `jaredzhou/pony/cmd/main` | 0 | 0 | 0 | 0 | 0 | 0 |

Observed type entry fields:

| field | observed type | notes |
|---|---|---|
| `name` | string | type name |
| `docstring` | string | may be empty or begin with newline |
| `signature` | string | HTML-enhanced signature with links |
| `loc` | object | source location |
| `methods` | array | detailed method entries |
| `impls` | array | implementation summaries; empty in sampled first entries |

Observed method entry fields:

| field | observed type | notes |
|---|---|---|
| `name` | string | method name |
| `docstring` | string | may be empty or begin with newline |
| `signature` | string | HTML-enhanced signature |
| `loc` | object | source location |

Observed value entry fields:

| field | observed type | notes |
|---|---|---|
| `name` | string | function/value name |
| `docstring` | string | may be empty |
| `signature` | string | HTML-enhanced signature |
| `loc` | object | source location |

Observed `loc` fields:

| field | observed type | example |
|---|---|---|
| `path` | string | `moonbitlang/core/argparse` |
| `file` | string | `arg_group.mbt` |
| `line` | integer | `17` |
| `column` | integer | `12` |

Example type entry from `moonbitlang/core/argparse`:

```json
{
  "name": "ArgGroup",
  "docstring": "\n Declarative argument group constructor.",
  "signature": "pub struct ArgGroup { ... }",
  "loc": {
    "path": "moonbitlang/core/argparse",
    "file": "arg_group.mbt",
    "line": 17,
    "column": 12
  },
  "methods": [
    {
      "name": "new",
      "docstring": "\n Create an argument group...",
      "signature": "fn ArgGroup::new(...) -> ArgGroup",
      "loc": {
        "path": "moonbitlang/core/argparse",
        "file": "arg_group.mbt",
        "line": 54,
        "column": 18
      }
    }
  ],
  "impls": []
}
```

Example value entry from `moonbit-community/cmark/char`:

```json
{
  "name": "ascii_digit_to_int",
  "docstring": "",
  "signature": "fn ascii_digit_to_int(c : Char) -> Int",
  "loc": {
    "path": "moonbit-community/cmark/char",
    "file": "ascii.mbt",
    "line": 44,
    "column": 8
  }
}
```

SeekMoon ingestion rule:

- Preserve `signature` as an HTML-bearing source string unless a later renderer deliberately strips tags.
- Treat `docstring` as present even when empty; empty string means no doc text in that entry.
- Treat `loc` as best-effort source coordinate. It identifies package path and source file but does not by itself prove source zip availability.
- Use `package_data.json` as the authoritative API detail source for `view --api` / `api` ingestion.

## `resources.json` Availability And Failure Semantics

Observed result: no 200 response found.

| sample | package | URL relpath | status |
|---|---|---|---:|
| required core package | `moonbitlang/core/argparse` | `argparse/resources.json` | 404 |
| required markdown package | `mizchi/markdown` | `resources.json` | 404 |
| required cmark package | `moonbit-community/cmark/char` | `char/resources.json` | 404 |
| sampled skill package | `jaredzhou/pony/cmd/main` | `cmd/main/resources.json` | 404 |
| bounded scan | 297 package attempts | mixed modules/packages | 404 only |

Interpretation:

- `resources.json` is an optional asset path advertised by the docs shell, but it is not present for the representative packages checked here.
- 404 does not mean the package is absent. Every required 404 sample also had valid `module_index.json` and/or `package_data.json`.
- 404 does not mean README/source docs are absent. Source zips for `mizchi/markdown`, `moonbit-community/cmark`, and `jaredzhou/pony` contain README files even though `resources.json` returned 404.
- SeekMoon should model `resources_json` as `absent` or `unavailable`, not as package failure.

Ingestion rule:

```text
resources.json 200 -> parse as optional package resources.
resources.json 404 -> record optional_resource_absent; continue with package_data, module_index, source zip, moon fetch, or local cache.
other non-2xx -> record fetch error with status; do not collapse into package absence.
```

## `resource.json` Follow-up Note

Do not confuse the plural `resources.json` with the singular `resource.json` used by the Mooncakes frontend source. The singular path is a separate follow-up object and is not covered by the bounded scan above.

## Source Zip Availability And Semantics

| module | version | status | evidence |
|---|---|---:|---|
| `moonbitlang/core` | `0.1.20260609+84519ca0a` | 404 | CloudFront final URL returned 404 body; source zip unavailable at documented URL shape |
| `mizchi/markdown` | `0.6.2` | 200 | 372 zip entries; includes `moon.mod.json`, package config, `.mbt`, README, LICENSE |
| `moonbit-community/cmark` | `0.4.4` | 200 | 126 zip entries; includes `moon.mod.json`, `moon.pkg`, `.mbt`, README.mbt.md, LICENSE |
| `jaredzhou/pony` | `0.1.1` | 200 | 28 zip entries; includes `moon.mod`, `moon.pkg`, `.mbt`, README, LICENSE, `cmd/main/main.mbt` |

Observed source zip URL redirects to CloudFront:

```text
https://download.mooncakes.io/user/<module>/<version>.zip
-> https://d15l9c1mnzh3r.cloudfront.net/user/<module>/<version>.zip
```

Source zip is sufficient published-source evidence when status is 200 because it contains the published module package structure and source files. It should be preferred over GitHub for published-version source evidence. It is not universal: `moonbitlang/core` latest manifest source zip returned 404, so SeekMoon must support fallback to `moon fetch`, local toolchain/library cache, local registry cache, or repository evidence depending on module type.

Ingestion rule:

```text
source_zip 200 -> unzip/list and treat as published source evidence.
source_zip 404 -> source zip unavailable for that module/version; do not infer module absence.
source_zip non-2xx -> record source fetch failure and try fallback.
```

## Skill Module Asset Evidence

Skills API sample:

```json
{
  "module": "jaredzhou/pony",
  "author": "jaredzhou",
  "version": "0.1.1",
  "package": "cmd/main",
  "name": "jaredzhou/pony/cmd/main",
  "detail_url": "/skills/jaredzhou/pony@0.1.1/cmd/main",
  "wasm_url": "/assets/jaredzhou/pony@0.1.1/cmd/main/main.wasm",
  "checksum_url": "/assets/jaredzhou/pony@0.1.1/cmd/main/main.wasm.sha256",
  "repository": "github.com/jaredzhou/moonbase"
}
```

The sampled skill module also had:

- manifest 200
- module index 200
- package data 200 for `cmd/main`
- source zip 200
- `resources.json` 404

The `package` field in Skills API is already an asset relpath (`cmd/main`). The `name` field is the full package path (`jaredzhou/pony/cmd/main`). SeekMoon should keep both:

| Skills field | ingestion meaning |
|---|---|
| `module` | module coordinate |
| `version` | version for assets and runwasm pinning |
| `package` | package relpath for assets |
| `name` | full package path / skill package id |
| `wasm_url` | executable Wasm asset path |
| `checksum_url` | checksum asset path |

## Existing SOP / Design Assumption Assessment

| assumption | status | assessment |
|---|---|---|
| Use manifest version before fetching assets | correct | Required for all asset URLs. |
| `module_index.json` is a tree structure | correct | Observed key is `childs`; traversal is required. |
| Tree key may be `children` | wrong for current data | Current representative data uses `childs`, not `children`. |
| Package paths should be confirmed from module index | correct | Do not guess package relpaths. |
| Package relpath is package path after stripping module path prefix | correct/incomplete previously | Must handle root package as empty relpath. |
| `package_data.json` exposes API data | correct | It exposes detailed types, values, methods, docstrings, signatures, and source locations. |
| Top-level package data keys include `types`, `values`, `traits`, `errors`, `typealias`, `misc` | correct | All observed samples had these keys plus `name`. |
| `resources.json` can be used as package README/resource source | incomplete | It is optional. No 200 found in this task; 404 is common. |
| `resources.json` 404 means package missing | wrong | Package data and source can exist while resources returns 404. |
| Source zip URL shape is useful | correct | Works for ordinary representative modules and sampled skill module. |
| Source zip is available for `moonbitlang/core` | wrong currently | Latest core version returned 404 at source zip URL. |
| Source zip proves published source when available | correct | Observed zips contain published module/package config and source files. |
| GitHub should be the primary published-source evidence | incomplete/wrong | GitHub is maintenance/upstream repository evidence; source zip or `moon fetch` is closer to published package evidence. |
| Skills API entries are package-like library candidates | wrong | They are executable/skill entries and need separate ingestion semantics. |

## Canonical Data Dictionary For SeekMoon Asset Ingestion

### Module Asset Identity

| field | source | required | meaning |
|---|---|---|---|
| `module` | manifest / modules API / skills API | yes | full module coordinate, e.g. `mizchi/markdown` |
| `version` | manifest / skills API | yes | asset version segment |
| `module_asset_base` | derived | yes | `https://mooncakes.io/assets/<module>@<version>/` |
| `module_index_url` | derived | yes | `<module_asset_base>/module_index.json` |
| `source_zip_url` | derived | yes | `https://download.mooncakes.io/user/<module>/<version>.zip` |

### Module Index Node

| field | source | required | meaning |
|---|---|---|---|
| `node.name` | module index | yes | tree node name |
| `node.childs` | module index | yes | child nodes; current spelling is `childs` |
| `node.package` | module index | nullable | package summary object or null |

### Package Summary

| field | source | required | meaning |
|---|---|---|---|
| `package.path` | module index | yes | canonical full package path |
| `package.relpath` | derived | yes | asset relpath after stripping module prefix |
| `package_data_url` | derived | yes | package data asset URL |
| `resources_url` | derived | optional | package resources asset URL |
| `traits` | module index package | yes | compact trait summaries |
| `errors` | module index package | yes | compact error summaries |
| `types` | module index package | yes | compact type summaries |
| `typealias` | module index package | yes | compact alias summaries |
| `values` | module index package | yes | compact value summaries |
| `misc` | module index package | yes | compact miscellaneous summaries |

### Package Data Entry

| field | source | required | meaning |
|---|---|---|---|
| `entry.kind` | derived from containing array | yes | `type`, `value`, `trait`, `error`, `typealias`, or `misc` |
| `entry.name` | package data | yes for observed type/value/method entries | symbol name |
| `entry.docstring` | package data | yes for observed type/value/method entries | documentation string; may be empty |
| `entry.signature` | package data | yes for observed type/value/method entries | HTML-enhanced signature |
| `entry.loc.path` | package data | yes when loc present | package path |
| `entry.loc.file` | package data | yes when loc present | source file path relative to package/source root context |
| `entry.loc.line` | package data | yes when loc present | source line |
| `entry.loc.column` | package data | yes when loc present | source column |
| `entry.methods` | package data type entry | yes for type entries | detailed method entries |
| `entry.impls` | package data type entry | yes for type entries | implementation summaries |

### Resource Status

| field | source | meaning |
|---|---|---|
| `resources.status` | HTTP status | `200`, `404`, or other |
| `resources.state` | derived | `available`, `optional_absent`, or `fetch_error` |
| `resources.data` | resources JSON | only when 200 and parse succeeds |

### Published Source Status

| field | source | meaning |
|---|---|---|
| `source_zip.status` | HTTP status | zip availability |
| `source_zip.final_url` | HTTP client | redirect target |
| `source_zip.file_count` | zip listing | number of entries when available |
| `source_zip.has_moon_mod` | zip listing | module config evidence |
| `source_zip.has_moon_pkg` | zip listing | package config evidence |
| `source_zip.has_source_mbt` | zip listing | MoonBit source evidence |
| `source_zip.has_readme` | zip listing | README evidence |
| `source_zip.has_license` | zip listing | license file evidence |
| `source_zip.state` | derived | `available`, `unavailable`, or `fetch_error` |

### Skill Entry

| field | source | meaning |
|---|---|---|
| `skill.module` | Skills API | module coordinate |
| `skill.version` | Skills API | published version |
| `skill.package_relpath` | Skills API `package` | package relpath |
| `skill.full_package_path` | Skills API `name` | full package path |
| `skill.detail_url` | Skills API | web skill page path |
| `skill.wasm_url` | Skills API | Wasm asset path |
| `skill.checksum_url` | Skills API | checksum asset path |
| `skill.repository` | Skills API | repository string |

## Implementation Consequences For `view`, `api`, And `source`

`view` ingestion:

- Fetch manifest.
- Fetch module index.
- Traverse `childs`.
- Display package count and compact API summaries from module index.
- Do not require `resources.json`.
- Show source zip availability only if explicitly fetched or cached; do not assume.

`api` ingestion:

- Resolve package id from module index.
- Derive package relpath.
- Fetch package data.
- Read top-level arrays `types`, `values`, `traits`, `errors`, `typealias`, `misc`.
- Preserve docstring/signature/loc/methods/impls.
- Treat empty arrays as valid no-symbol results.

`source` ingestion:

- Construct source zip URL from manifest version.
- Follow redirect.
- If 200, inspect zip as published source evidence.
- If 404, record source zip unavailable and use fallback path. For `moonbitlang/core`, expect local toolchain/library source to be a likely fallback.
- Do not substitute GitHub HEAD for published-source evidence without marking it as repository evidence, not registry source evidence.

## Raw Evidence Files

Important saved files:

| file | content |
|---|---|
| `summary.json` | consolidated schema and availability summary |
| `fetch_log.json` | first-pass endpoint fetch log |
| `modules.json` | live `/api/v0/modules` response |
| `statistics.json` | live statistics response |
| `skills.json` | live Skills API response |
| `modules/moonbitlang__core/module_index.json` | core module index |
| `modules/moonbitlang__core/argparse_package_data.json` | core argparse package data |
| `modules/mizchi__markdown/module_index.json` | markdown module index |
| `modules/mizchi__markdown/_root_package_data.json` | markdown root package data |
| `modules/moonbit-community__cmark/module_index.json` | cmark module index |
| `modules/moonbit-community__cmark/char_package_data.json` | cmark char package data |
| `source_zips/*.zip` | fetched source zip responses |

## Final Answer

SeekMoon can implement current `view`, `api`, and `source` ingestion without guessing asset field names if it uses:

1. manifest version for all asset paths;
2. `module_index.json` with `childs` traversal for package discovery;
3. derived package relpaths for package assets;
4. `package_data.json` for detailed API symbols;
5. optional `resources.json` semantics where 404 is not failure;
6. source zip as preferred published-source evidence when available, with fallback required for modules like `moonbitlang/core`;
7. separate skill ingestion for Skills API executable entries.
