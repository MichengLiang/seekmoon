# Task 01: Public API Schema Validation

## Summary

Live fetch window: 2026-06-21T12:10:30Z to 2026-06-21T12:10:35Z.

Toolchain used:

- `uv 0.11.4`
- `Python 3.13.12`
- `httpx 0.28.1`
- Raw files and analysis artifacts: `/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-01/`

Main findings:

- `/api/v0/modules` returned a top-level array with 1350 module objects.
- `/api/v0/modules/statistics` returned a top-level object with four integer counters.
- `/api/v0/modules?search=markdown` and `/api/v0/modules?search=cowsay` both returned byte-equivalent parsed JSON to the unfiltered `/api/v0/modules` response: 1350 objects, first module `0Ayachi0/elk`.
- `/api/v0/skills` returned a top-level array with 70 skill objects.
- `/api/v0/modules` module fields are always present as keys, but `description`, `repository`, `license`, and `keywords` use empty string or empty array for missing content.
- Public API target metadata is not stable across endpoints. `/modules` and `/skills` expose no `targets`, `supported-targets`, `supported_targets`, or nested `metadata.*` target fields. In representative manifests, only `Yoorkin/cowsay` exposes `metadata["supported-targets"]`, as the string `"wasm"`.
- `/skills` fields `wasm_url`, `checksum_url`, `detail_url`, `package`, and `metadata` are present on all 70 objects. `package` is sometimes `""`; `metadata` is always an object, but its inner fields are sparse or empty.

## Fetch Log

| name | endpoint | UTC fetch time | status | shape |
|---|---|---|---:|---|
| `modules` | `/api/v0/modules` | 2026-06-21T12:10:30.406963Z to 12:10:31.732810Z | 200 | array, 1350 |
| `modules_statistics` | `/api/v0/modules/statistics` | 2026-06-21T12:10:31.740088Z to 12:10:32.017918Z | 200 | object, 4 keys |
| `modules_search_markdown` | `/api/v0/modules?search=markdown` | 2026-06-21T12:10:32.018239Z to 12:10:32.623202Z | 200 | array, 1350 |
| `modules_search_cowsay` | `/api/v0/modules?search=cowsay` | 2026-06-21T12:10:32.628536Z to 12:10:33.127746Z | 200 | array, 1350 |
| `skills` | `/api/v0/skills` | 2026-06-21T12:10:33.133065Z to 12:10:33.409934Z | 200 | array, 70 |
| manifests | 8 representative `/api/v0/manifest/<owner>/<module>` requests | 2026-06-21T12:10:33.410844Z to 12:10:35.729564Z | 200 | object, 9 keys each |

Representative manifests fetched:

- `moonbitlang/core`
- `mizchi/markdown`
- `moonbit-community/cmark`
- `Yoorkin/cowsay`
- `202306123/lib_mbt`
- `0Ayachi0/SwissTable`
- `bobzhang/testa`
- `AdUhTkJm/assembler`

## Endpoint Shapes

| endpoint | top-level shape | current keys or item shape |
|---|---|---|
| `/api/v0/modules` | array | each item is an object with `created_at`, `description`, `is_new`, `keywords`, `license`, `name`, `repository`, `version` |
| `/api/v0/modules/statistics` | object | `total_downloads`, `total_lines`, `total_modules`, `total_packages`; all integers |
| `/api/v0/manifest/<owner>/<module>` | object | `build_status`, `downloads`, `has_package`, `latest_version`, `metadata`, `module`, `name`, `version`, `versions` |
| `/api/v0/skills` | array | each item is an object with `author`, `author_avatar`, `checksum_url`, `created_at`, `detail_url`, `metadata`, `module`, `name`, `package`, `repository`, `version`, `wasm_url` |

Short evidence snippets:

```json
{
  "modules_type": "array",
  "modules_length": 1350,
  "modules_first": "0Ayachi0/elk"
}
```

```json
{
  "total_downloads": 4043104,
  "total_lines": 44597511,
  "total_modules": 1350,
  "total_packages": 12008
}
```

## `/api/v0/modules` Field Presence

All eight observed module fields are present as keys on all 1350 module objects. Missing content is encoded as empty string or empty array, not absent keys or JSON null.

| field | present | absent | observed value categories |
|---|---:|---:|---|
| `created_at` | 1350 | 0 | nonempty string: 1350 |
| `description` | 1350 | 0 | nonempty string: 1023; empty string: 327 |
| `is_new` | 1350 | 0 | boolean: 1350 |
| `keywords` | 1350 | 0 | nonempty array: 933; empty array: 417 |
| `license` | 1350 | 0 | nonempty string: 1336; empty string: 14 |
| `name` | 1350 | 0 | nonempty string: 1350 |
| `repository` | 1350 | 0 | nonempty string: 1058; empty string: 292 |
| `version` | 1350 | 0 | nonempty string: 1350 |

Examples of empty values:

```json
{
  "empty_description": "202306123/lib_mbt",
  "empty_repository": "0Ayachi0/SwissTable",
  "empty_license": "bobzhang/testa"
}
```

License distribution highlights:

| license | count |
|---|---:|
| `Apache-2.0` | 1037 |
| `MIT` | 253 |
| empty string | 14 |
| `GPL-3.0-or-later` | 8 |
| `ISC` | 7 |
| `AGPL-3.0` | 7 |
| `MPL-2.0` | 5 |
| `MulanPSL-2.0` | 4 |

## `/api/v0/modules?search=<query>`

Observed behavior: still unfiltered.

| request | length | same parsed JSON as `/modules` | first module |
|---|---:|---|---|
| `/api/v0/modules` | 1350 | baseline | `0Ayachi0/elk` |
| `/api/v0/modules?search=markdown` | 1350 | yes | `0Ayachi0/elk` |
| `/api/v0/modules?search=cowsay` | 1350 | yes | `0Ayachi0/elk` |

Assessment: SeekMoon must not treat `search` as a server-side filter. The current reliable model remains full-list fetch plus local filtering.

## Manifest Schema Observations

Across the eight representative manifest responses, these top-level fields were always present:

| field | present | absent | observed value categories |
|---|---:|---:|---|
| `build_status` | 8 | 0 | nonempty string: 6; null: 2 |
| `downloads` | 8 | 0 | integer: 8 |
| `has_package` | 8 | 0 | boolean: 8 |
| `latest_version` | 8 | 0 | nonempty string: 8 |
| `metadata` | 8 | 0 | nonempty object: 8 |
| `module` | 8 | 0 | nonempty string: 8 |
| `name` | 8 | 0 | nonempty string: 8 |
| `version` | 8 | 0 | nonempty string: 8 |
| `versions` | 8 | 0 | nonempty array: 8 |

`versions[]` items in all representative samples used:

```json
{
  "version": "0.6.2",
  "yanked": false
}
```

Observed `build_status` values in samples:

- `success`
- `legacy`
- `null`

`metadata` is flexible and not uniform. Across representative samples:

| metadata field | present | absent | observed value categories |
|---|---:|---:|---|
| `checksum` | 8 | 0 | nonempty string: 8 |
| `name` | 8 | 0 | nonempty string: 8 |
| `created_at` | 7 | 1 | nonempty string: 7 |
| `license` | 7 | 1 | nonempty string: 7 |
| `repository` | 7 | 1 | nonempty string: 4; empty string: 3 |
| `description` | 6 | 2 | nonempty string: 4; empty string: 2 |
| `keywords` | 6 | 2 | nonempty array: 3; empty array: 3 |
| `readme` | 6 | 2 | nonempty string: 6 |
| `source` | 5 | 3 | nonempty string: 5 |
| `deps` | 4 | 4 | nonempty object: 4 |
| `authors` | 1 | 7 | nonempty array: 1 |
| `dependencies` | 1 | 7 | empty string: 1 |
| `preferred-target` | 1 | 7 | nonempty string: 1 |
| `preferred_target` | 1 | 7 | nonempty string: 1 |
| `supported-targets` | 1 | 7 | nonempty string: 1 |
| `warn-list` | 1 | 7 | nonempty string: 1 |

Representative manifest summaries:

| module | version | downloads | has_package | build_status | notable metadata |
|---|---:|---:|---|---|---|
| `moonbitlang/core` | `0.1.20260609+84519ca0a` | 270 | false | `success` | `warn-list`, no `description` |
| `mizchi/markdown` | `0.6.2` | 7587 | true | `success` | `preferred-target: "js"` |
| `moonbit-community/cmark` | `0.4.4` | 3230 | false | `success` | `deps` object |
| `Yoorkin/cowsay` | `0.1.0` | 2 | true | `success` | `preferred_target: "wasm"`, `supported-targets: "wasm"` |
| `202306123/lib_mbt` | `0.1.0` | 3 | false | null | empty `description`, empty `repository` |
| `0Ayachi0/SwissTable` | `0.1.0` | 5 | true | `legacy` | `authors`, `dependencies: ""` |
| `bobzhang/testa` | `0.1.0` | 1 | false | null | minimal metadata: `checksum`, `name`, `version` |
| `AdUhTkJm/assembler` | `0.1.1` | 1 | false | `legacy` | unusual license `GPL-3.0-or-later` |

Short evidence snippets:

```json
{
  "module": "mizchi/markdown",
  "version": "0.6.2",
  "build_status": "success",
  "metadata": {
    "description": "Incremental Markdown parser and compiler",
    "preferred-target": "js",
    "source": "src"
  }
}
```

```json
{
  "module": "Yoorkin/cowsay",
  "metadata": {
    "preferred_target": "wasm",
    "supported-targets": "wasm"
  }
}
```

## Target Field Validation

Observed target spelling presence:

| endpoint group | `targets` | `supported-targets` | `supported_targets` | `metadata.targets` | `metadata.supported-targets` | `metadata.supported_targets` |
|---|---:|---:|---:|---:|---:|---:|
| `/modules`, 1350 objects | 0 | 0 | 0 | 0 | 0 | 0 |
| `/skills`, 70 objects | 0 | 0 | 0 | 0 | 0 | 0 |
| representative manifests, 8 objects | 0 | 0 | 0 | 0 | 1 | 0 |

Only `Yoorkin/cowsay` among representative manifests had target metadata:

```json
{
  "metadata": {
    "preferred_target": "wasm",
    "supported-targets": "wasm"
  }
}
```

Important spelling drift:

- `mizchi/markdown` uses `metadata["preferred-target"]`.
- `Yoorkin/cowsay` uses `metadata.preferred_target`.
- Only `Yoorkin/cowsay` sample uses `metadata["supported-targets"]`.
- No sample used `metadata.targets`, top-level `targets`, top-level `supported-targets`, or `supported_targets`.

Assessment: target compatibility cannot be treated as a stable module-list field. Manifest ingestion should preserve raw `metadata`, normalize known spelling variants into a derived target view, and keep unknown when no evidence exists.

## `/api/v0/skills` Field Presence

All 12 top-level skill fields are present on all 70 skill objects. Missing content is represented by empty strings or sparse `metadata`, not absent top-level keys.

| field | present | absent | observed value categories |
|---|---:|---:|---|
| `author` | 70 | 0 | nonempty string: 70 |
| `author_avatar` | 70 | 0 | nonempty string: 54; empty string: 16 |
| `checksum_url` | 70 | 0 | nonempty string: 70 |
| `created_at` | 70 | 0 | nonempty string: 70 |
| `detail_url` | 70 | 0 | nonempty string: 70 |
| `metadata` | 70 | 0 | nonempty object: 70 |
| `module` | 70 | 0 | nonempty string: 70 |
| `name` | 70 | 0 | nonempty string: 70 |
| `package` | 70 | 0 | nonempty string: 62; empty string: 8 |
| `repository` | 70 | 0 | nonempty string: 68; empty string: 2 |
| `version` | 70 | 0 | nonempty string: 70 |
| `wasm_url` | 70 | 0 | nonempty string: 70 |

`metadata` inner fields:

| metadata field | count / shape |
|---|---|
| `description` | present on all 70, but empty string for 65 and nonempty for 5 |
| `name` | nonempty string for 4, empty string for 37, absent for 29 |

First skill shape:

```json
{
  "module": "jaredzhou/pony",
  "package": "cmd/main",
  "detail_url": "/skills/jaredzhou/pony@0.1.1/cmd/main",
  "wasm_url": "/assets/jaredzhou/pony@0.1.1/cmd/main/main.wasm",
  "checksum_url": "/assets/jaredzhou/pony@0.1.1/cmd/main/main.wasm.sha256"
}
```

`Yoorkin/cowsay` skill was present:

```json
{
  "module": "Yoorkin/cowsay",
  "version": "0.1.0",
  "package": "",
  "detail_url": "/skills/Yoorkin/cowsay@0.1.0",
  "wasm_url": "/assets/Yoorkin/cowsay@0.1.0/cowsay.wasm",
  "checksum_url": "/assets/Yoorkin/cowsay@0.1.0/cowsay.wasm.sha256"
}
```

Assessment for required fields:

- `wasm_url`: stable top-level key, always nonempty string.
- `checksum_url`: stable top-level key, always nonempty string.
- `detail_url`: stable top-level key, always nonempty string.
- `package`: stable top-level key, but may be empty string for root executable packages.
- `metadata`: stable top-level key and always object, but inner schema is sparse and not reliable beyond observed `description` and `name`.

## Existing Design / SOP Field Assessment

| claimed or implied field/source | assessment | current finding |
|---|---|---|
| `/api/v0/modules` as full module list | Correct | Top-level array, 1350 objects. |
| `/api/v0/modules?search=<query>` as server-side search | Wrong / stale | Returns the same full array as `/modules`; do not use as search filter. |
| module `name` | Correct | Present and nonempty on all module objects. |
| module `version` | Correct | Present and nonempty on all module objects. |
| module `description` | Incomplete | Key is always present, but 327 values are `""`; not absent or null. |
| module `keywords` | Incomplete | Key is always present, but 417 values are `[]`. |
| module `repository` | Incomplete | Key is always present, but 292 values are `""`. |
| module `license` | Incomplete | Key is always present, but 14 values are `""`; unusual licenses exist. |
| module `created_at` | Correct | Present and nonempty on all module objects. |
| module `is_new` | Correct | Present boolean on all module objects. |
| manifest top-level `downloads` | Correct | Present integer in all representative manifests. |
| manifest top-level `build_status` | Incomplete | Present in all representative manifests, but may be `null`, `legacy`, or `success`. |
| manifest `has_package` | Correct | Present boolean in all representative manifests. |
| manifest `versions` / `versions_count` | Correct with derivation | `versions` present; `versions_count` is not a field and must be derived from array length. |
| manifest `metadata.repository`, `metadata.license`, `metadata.description`, `metadata.keywords` | Incomplete | These are common but not stable across manifests; some are absent or empty. |
| manifest `metadata.deps` | Incomplete | Present in 4 of 8 samples; `dependencies` also appears as an empty string in one sample. |
| `metadata.targets` | Wrong for current samples | Not present in representative manifests. |
| `metadata.supported-targets` | Incomplete | Present only on `Yoorkin/cowsay`, as string `"wasm"`. |
| `targets` top-level | Wrong for current public API | Not present in modules, skills, or representative manifests. |
| `supported_targets` spelling | Wrong for current samples | Not present in modules, skills, or representative manifests. |
| `preferred-target` | Incomplete | Present in `mizchi/markdown`, spelling uses hyphen. |
| `preferred_target` | Incomplete | Present in `Yoorkin/cowsay`, spelling uses underscore. |
| `/api/v0/skills` | Correct | Top-level array, 70 skill objects. |
| skill `wasm_url`, `checksum_url`, `detail_url` | Correct | Present and nonempty on all 70 skill objects. |
| skill `package` | Incomplete | Present on all skill objects, but empty string for 8. |
| skill `metadata` | Incomplete | Present object on all skill objects, but inner fields are sparse or empty. |
| `docs_url` in API | Wrong as raw field | Not present in public API responses; can be derived from module name. |
| `module_index`, `package_data`, `source_fetch` | Correct as downstream sources, not public API fields | Not part of the four public API responses validated in Task 01. |

## Canonical Data Dictionary For SeekMoon Public API Ingestion

### Common Representation Rules

- Preserve raw upstream objects for audit.
- Normalize empty string and empty array as missing content at the ingestion layer, but retain original raw values.
- Do not invent fields not present in current command schema.
- For target compatibility, store both raw metadata and a derived normalized target view with source evidence.
- Treat derived fields separately from upstream fields.

### Module List Entity

Source: `/api/v0/modules`

| canonical field | upstream field | type | nullability / missing rule | note |
|---|---|---|---|---|
| `module` | `name` | string | required, nonempty | Installation and manifest coordinate. |
| `version` | `version` | string | required, nonempty | Latest list version. |
| `description` | `description` | string or missing-content | key required; `""` means missing content | Do not treat empty string as useful description. |
| `keywords` | `keywords` | string array or missing-content | key required; `[]` means missing content | Search should handle empty array. |
| `repository` | `repository` | string or missing-content | key required; `""` means missing content | URL may still require validation separately. |
| `license` | `license` | string or missing-content | key required; `""` means missing content | Preserve exact spelling. |
| `created_at` | `created_at` | datetime string | required, nonempty | Parse cautiously; upstream returns string. |
| `is_new` | `is_new` | boolean | required | UI/list freshness signal. |
| `targets` | none | derived | unknown unless manifest/source evidence exists | Do not read from module list. |

### Statistics Entity

Source: `/api/v0/modules/statistics`

| canonical field | upstream field | type | note |
|---|---|---|---|
| `total_modules` | `total_modules` | integer | Matched `/modules` length during this fetch. |
| `total_packages` | `total_packages` | integer | Registry counter. |
| `total_lines` | `total_lines` | integer | Registry counter. |
| `total_downloads` | `total_downloads` | integer | Registry counter. |

### Manifest Entity

Source: `/api/v0/manifest/<owner>/<module>`

| canonical field | upstream field | type | nullability / missing rule | note |
|---|---|---|---|---|
| `module` | `module` / `name` | string | required, nonempty | Both observed and equal in samples. |
| `version` | `version` | string | required, nonempty | Current manifest version. |
| `latest_version` | `latest_version` | string | required, nonempty | Same as `version` in samples, but store separately. |
| `downloads` | `downloads` | integer | required | Current download counter. |
| `has_package` | `has_package` | boolean | required | Does not imply docs/API assets are sufficient alone. |
| `build_status` | `build_status` | string or null | required key; value may be null | Observed `success`, `legacy`, null. |
| `metadata` | `metadata` | object | required object | Flexible raw metadata bag. |
| `versions` | `versions` | array | required, nonempty in samples | Items observed with `version`, `yanked`. |
| `versions_count` | derived from `versions.length` | integer | derived | Do not expect upstream field. |

Manifest metadata canonical handling:

| canonical field | upstream spelling(s) | type | rule |
|---|---|---|---|
| `metadata_name` | `metadata.name` | string | Required in representative samples. |
| `metadata_version` | `metadata.version` | string | Present in all representative samples. |
| `checksum` | `metadata.checksum` | string | Present in all representative samples. |
| `description` | `metadata.description` | string or missing-content | May be absent or `""`. |
| `keywords` | `metadata.keywords` | string array or missing-content | May be absent or `[]`. |
| `repository` | `metadata.repository` | string or missing-content | May be absent or `""`. |
| `license` | `metadata.license` | string or missing-content | May be absent. |
| `readme` | `metadata.readme` | string | May be absent. |
| `source` | `metadata.source` | string | May be absent. |
| `deps` | `metadata.deps` | object | May be absent. |
| `dependencies` | `metadata.dependencies` | unknown/string | Observed as empty string in one sample; keep raw. |
| `preferred_target` | `metadata["preferred-target"]`, `metadata.preferred_target` | string | Normalize spelling, preserve source spelling. |
| `supported_targets` | `metadata["supported-targets"]` | string or future array | Observed as string `"wasm"` once; do not assume array. |
| `raw_metadata` | `metadata` | object | Always store raw bag for drift handling. |

### Skill Entity

Source: `/api/v0/skills`

| canonical field | upstream field | type | nullability / missing rule | note |
|---|---|---|---|---|
| `module` | `module` | string | required, nonempty | Module coordinate. |
| `version` | `version` | string | required, nonempty | Skill artifact version. |
| `author` | `author` | string | required, nonempty | Author name. |
| `author_avatar` | `author_avatar` | string or missing-content | `""` means missing content | UI-only signal. |
| `package` | `package` | string or root marker | `""` means root/default package | Do not drop empty package; it affects URL/coordinate formatting. |
| `name` | `name` | string | required, nonempty | Skill display/coordinate name. |
| `detail_url` | `detail_url` | string | required, nonempty | Relative URL. |
| `wasm_url` | `wasm_url` | string | required, nonempty | Relative artifact URL. |
| `checksum_url` | `checksum_url` | string | required, nonempty | Relative checksum URL. |
| `repository` | `repository` | string or missing-content | `""` means missing content | May not be full URL in all cases; validate separately. |
| `created_at` | `created_at` | datetime string | required, nonempty | Upstream string. |
| `metadata` | `metadata` | object | required object | Inner fields sparse. |
| `metadata.description` | `metadata.description` | string or missing-content | present in all 70, but often `""` | Do not assume useful description. |
| `metadata.name` | `metadata.name` | string or missing-content | absent or empty often | Use top-level `name` as stable field. |

## Ingestion Recommendations

1. Make `/api/v0/modules` the module candidate source, not `/modules?search=...`.
2. In the module-list model, treat empty strings and empty arrays as missing content, but do not report them as absent keys.
3. In the manifest model, keep `metadata` as an open raw object. Normalize only known fields into derived fields.
4. Do not expose target compatibility as a stable module-list field. Use manifest-derived evidence with source spelling, and mark unknown when no target metadata exists.
5. Accept both `preferred-target` and `preferred_target`; prefer a normalized derived field while preserving raw metadata.
6. Accept `metadata["supported-targets"]` as current observed spelling; store value shape as string-or-list-compatible because current value is a string.
7. Treat `build_status` as nullable and not limited to `success`.
8. For skills, keep root-package `package: ""` as meaningful rather than dropping it.
9. For skills, rely on top-level `name`, `module`, `wasm_url`, `checksum_url`, and `detail_url`; treat `metadata` as supplementary.
10. Keep derived fields (`docs_url`, `versions_count`, normalized targets) separate from upstream fields.

## Files Produced

- Raw endpoint snapshots: `/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-01/*.json`
- Fetch metadata: `/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-01/*.meta.json`
- Aggregated fetch log: `/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-01/fetch_log.json`
- Analysis artifact: `/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-01/analysis.json`
