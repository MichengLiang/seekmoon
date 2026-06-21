# Task 05: Canonical Dictionary And Drift Audit

Date: 2026-06-21  
Role: reviewer / reconnaissance agent  
Scope: existing SeekMoon design documents and current upstream/local data facts.

## Evidence Snapshot

Live spot checks were taken on 2026-06-21 from the current local toolchain and Mooncakes public endpoints.

| Evidence | Current observation |
|---|---|
| `https://mooncakes.io/api/v0/modules/statistics` | `total_modules=1350`, `total_packages=12008`, `total_lines=44597542`, `total_downloads=4043948` |
| `https://mooncakes.io/api/v0/modules` | array length `1350`; list item keys: `name`, `version`, `license`, `repository`, `keywords`, `description`, `is_new`, `created_at` |
| `https://mooncakes.io/api/v0/modules?search=markdown` | still returns full array length `1350`; first item `0Ayachi0/elk`; not a server-side search endpoint |
| `https://mooncakes.io/api/v0/manifest/mizchi/markdown` | keys: `name`, `module`, `version`, `latest_version`, `downloads`, `has_package`, `build_status`, `metadata`, `versions` |
| manifest metadata sample | `deps`, `readme`, `repository`, `license`, `keywords`, `description`, `source`, `preferred-target`, `checksum`, `created_at` |
| manifest versions sample | entries include `version`, `yanked` |
| module index asset | `module_index.json` is a tree with `name`, `package`, `childs`; package nodes contain `path`, `types`, `traits`, `errors`, `typealias`, `values`, `misc` |
| package data asset | `package_data.json` contains `name`, `types`, `values`, `traits`, `errors`, `typealias`, `misc`; type entries include `name`, `docstring`, `signature`, `loc`, `methods` |
| `https://mooncakes.io/api/v0/skills` | array length `70`; item keys include `module`, `author`, `author_avatar`, `version`, `package`, `name`, `detail_url`, `wasm_url`, `checksum_url`, `metadata`, `repository`, `created_at` |
| local `moon --version` | `moon 0.1.20260608`, `moonc v0.10.0+e66899a54`, `moonrun 0.1.20260608` |
| local `moon --help` | has `runwasm`, `fetch`, `ide doc`; no `search`, `view`, `audit`, `outdated` package discovery commands |
| local registry index | `1363` `.index` files; current git HEAD `7503ed87 update cybershang/agent-telemetry` |
| local index item sample | JSON Lines keys include `name`, `version`, `readme`, `repository`, `license`, `keywords`, `description`, `source`, `checksum`, `created_at`, sometimes `deps` |
| local symbols cache | `3` `.symbols` files; partial cache, not full registry coverage |
| `moon fetch --help` | current but explicitly unstable |
| `moon ide doc --help` | searches current module, core, and registry symbol indexes; not a registry package search command |

## Claimed Data Sources Audit

| Claimed source | Claim location | Status | Current judgment |
|---|---|---|---|
| Mooncakes Web home search | investigation report, SOP, CLI design | `confirmed-current` | Web search exists as user-facing discovery, but query state is not a stable API contract for SeekMoon ingestion. Use as human-facing reference, not canonical ingestion. |
| Modules API `/api/v0/modules` | all scoped docs | `confirmed-current` | Canonical v0 source for library module search corpus. Pull full JSON and filter locally. |
| Modules API query parameter `/api/v0/modules?search=<query>` | investigation report, SOP | `wrong` for server-side search | Current response is the full modules array. Must not be modeled as upstream search. |
| Statistics API `/api/v0/modules/statistics` | investigation report, SOP, CLI design | `confirmed-current` | Canonical source for registry snapshot counts. Counts are temporally unstable. |
| Manifest API `/api/v0/manifest/<owner>/<module>` | investigation report, SOP, CLI design | `confirmed-current` | Canonical source for per-module detail, downloads, build status, metadata, and version/yank list. |
| Mooncakes docs pages `/docs/...` | scoped docs | `confirmed-current` but not canonical machine content | Useful public URL. Direct HTML is SPA shell; API/docs content must come from assets or `moon ide doc`. |
| Module index asset `module_index.json` | SOP, CLI design | `confirmed-current` | Canonical source for module package/API structure index. It is tree-shaped. |
| Package data asset `package_data.json` | SOP, CLI design | `confirmed-current` | Canonical source for single package API details, signatures, docstrings, locations. |
| Package resources asset `resource.json` | official frontend source / follow-up | `missing-value` / `unavailable` possible | Valid asset class in current frontend source, but not guaranteed for each package. 404 means asset unavailable, not package absence. |
| Source zip `download.mooncakes.io/user/.../<version>.zip` | SOP, CLI design | `confirmed-current` with availability caveat | Valid source retrieval path for published source when available. Must model failures. |
| Skills API `/api/v0/skills` | SOP, CLI design, package closure investigation | `confirmed-current` | Canonical source for executable Wasm / skill entries. Do not merge with library dependency candidates. |
| Wasm asset URL from Skills API | SOP, CLI design | `confirmed-current` | Skill-only artifact path. Store as URL evidence; do not imply library package support. |
| Checksum asset URL from Skills API | SOP, CLI design | `confirmed-current` | Skill-only checksum evidence path. Availability still requires fetch/status modeling. |
| Local registry index `~/.moon/registry/index/user/**/*.index` | SOP, CLI design, investigation report | `confirmed-local` | Canonical local cross-check/offline source. JSON Lines; counts can differ from public API. |
| Local symbols cache `~/.moon/registry/symbols/**/*.symbols` | SOP, CLI design | `confirmed-local` | Partial local symbol source. Not full registry search. Missing hit means unknown, not absent. |
| Local Moon toolchain commands | all scoped docs | `confirmed-local` | Current `moon` supports package management, `runwasm`, unstable `fetch`, `ide doc`, local validation commands. |
| `moon ide doc` | SOP, CLI design | `confirmed-local` | API/symbol documentation source for current module, core, and symbol indexes. Not package discovery. |
| `moon fetch` | SOP, CLI design | `confirmed-local` | Valid current source retrieval command, but explicitly unstable. |
| Current project files `moon.mod(.json)`, `moon.pkg(.json)` | CLI design | `confirmed-local` | Context source for probe and project-aware reporting. |
| Local Moon cache `~/.moon/registry/cache`, `~/.moon/lib/core` | SOP, CLI design | `confirmed-local` | Valid local source/artifact fallback. Shape is implementation/cache state, not upstream schema. |
| GitHub repository from metadata | investigation report, SOP, CLI design | `derived` / external-maintenance source | Not Mooncakes upstream. Use for maintenance signals only: archived, pushed, CI, issues, releases. Not the canonical published source. |
| GitHub API repo status | investigation report, SOP | `derived` / external-maintenance source | Legitimate derived enrichment when repository is present. Must not be required for registry facts. |
| Probe project / `moon add`, `moon check`, `moon test`, target build | SOP, CLI design | `derived` | SeekMoon-computed validation evidence, not upstream metadata. |
| User records `.seekmoon/...` | CLI design | `derived` | SeekMoon-owned audit state. It records decisions and evidence, not upstream facts. |
| Web search deep links `/search?q=...` | investigation report | `unsupported-future` | Claimed as desired/observed gap. Not current upstream source. |
| Server-side search API `/api/v0/search` | investigation report | `unsupported-future` | Future watch item only. |
| Advisory / audit feed | investigation report, CLI design as future | `unsupported-future` | No current evidence source. Must not appear in v0 output. |
| Reverse dependency / dependents source | investigation report, CLI design as future | `unsupported-future` | No current evidence source. |
| Quality score / Mooncake Score | investigation report, CLI design as future | `unsupported-future` | Not current upstream field. SeekMoon must not emit it in v0. |
| Provenance / publisher identity / signatures / SBOM | investigation report, CLI design as future | `unsupported-future` | No current v0 source except checksum strings/URLs where explicitly present. |

## Claimed Fields Audit

| Field or field group | Current source | Status | Canonical treatment |
|---|---|---|---|
| `module` / module coordinate | Modules API item `name`; manifest `module`/`name`; skills `module` | `confirmed-current` | Canonical identity for library and skill records. For modules list, normalize `name` to SeekMoon `module`. |
| `version` | Modules API, manifest, skills, local index | `confirmed-current` | Current selected/listed version. |
| `latest_version` | manifest | `confirmed-current` | Detail-only field. Do not invent in search list if manifest not fetched. |
| `versions` | manifest | `confirmed-current` | Detail field. Entries include `version`, `yanked`. |
| `versions_count` | computed from manifest `versions` | `derived` | Legitimate SeekMoon-computed count. |
| `yanked` | manifest `versions[].yanked` | `confirmed-current` | Version-entry field; package-level current yanked status must be derived only from selected version. |
| `description` | Modules API, manifest metadata, skills metadata | `confirmed-current` / `missing-value` | Current field. In current modules sample, no null descriptions observed, but records can be empty strings or absent across sources; model `missing` when required by command schema. |
| `keywords` | Modules API, manifest metadata | `confirmed-current` / `missing-value` | Current field. Empty array is common. |
| `repository` | Modules API, manifest metadata, skills | `confirmed-current` / `missing-value` | Current registry field. GitHub checks are derived enrichment. |
| `license` | Modules API, manifest metadata, local index | `confirmed-current` / `missing-value` | Current field. Required for adoption judgment; model missing distinctly. |
| `is_new` | Modules API | `confirmed-current` | List-only upstream field. Include only where list/search snapshot needs it; do not derive from it without defining semantics. |
| `created_at` | Modules API, manifest metadata, local index, skills | `confirmed-current` | Timestamp source. Temporal and source-specific. |
| `downloads` | manifest | `confirmed-current` | Detail field. Search output may include only after manifest enrichment or if separately fetched; not present in modules list. |
| `build_status` | manifest | `confirmed-current` | Mooncakes build signal. Does not prove every target works. |
| `has_package` | manifest | `confirmed-current` | Detail field. Interpret as Mooncakes package/documentation availability signal, not local install success. |
| `metadata.name` | manifest metadata, skills metadata | `confirmed-current` | Detail metadata. For manifest usually repeats module name. |
| `metadata.version` | manifest metadata | `confirmed-current` | Detail metadata. Usually repeats selected version. |
| `deps` | manifest metadata, local index sometimes | `confirmed-current` | Dependency metadata from registry. Shape is object map in observed manifest. |
| `readme` | manifest metadata, local index | `confirmed-current` | Metadata path hint. It is not README body. |
| `source` | manifest metadata, local index | `confirmed-current` | Source directory/path hint. Not fetched source content. |
| `checksum` | manifest metadata, local index | `confirmed-current` | Registry checksum string. Do not conflate with skill wasm checksum URL. |
| `preferred-target` | manifest metadata | `confirmed-current` | Current metadata key. Use exact upstream key in raw; normalize to `preferred_target` in SeekMoon JSON if schema defines it. |
| `targets` | manifest metadata claim in docs | `unknown-evidence` | Current sample did not include generic `targets`; SOP query may check it. Treat absent as unknown when target compatibility is requested. |
| `supported-targets` | manifest metadata claim in docs | `unknown-evidence` | Current sample did not include it. Treat as unknown when target compatibility is requested. |
| `target compatibility` | manifest metadata, local validation | `derived` / `unknown-evidence` | Not a single upstream truth. SeekMoon may compute status from metadata and probe results. Unknown when no target evidence. |
| `module_index` availability | asset fetch result | `derived` from current asset source | Store asset status and parsed summary. |
| `package path` | module index package node `path` | `confirmed-current` | Canonical package coordinate inside module index. |
| API symbol names | module index and package data | `confirmed-current` | Types/values/traits/errors/typealias/methods/misc are current API data groups. |
| API `docstring` | package data | `confirmed-current` / `missing-value` | Current field; docstring may be null/empty. |
| API `signature` | package data | `confirmed-current` | Current field; contains rendered HTML links in observed source, so consumers may need cleaned text as derived projection. |
| API `loc` | package data | `confirmed-current` / `missing-value` | Current source location field. |
| `resources` / README resources | resources asset | `missing-value` / `unavailable` | Valid asset class but may be 404. Must model unavailable. |
| `source_fetch` | `moon fetch`, source zip, cache fallback | `derived` | SeekMoon-computed retrieval result. Must include method/status/path/error. |
| `skill_entry` | Skills API | `confirmed-current` | Skill-only field group. |
| skill `author` | Skills API | `confirmed-current` | Skill-only. |
| skill `author_avatar` | Skills API | `confirmed-current` | Skill-only; display-only enrichment. |
| skill `package` | Skills API | `confirmed-current` | Executable package path. |
| skill `name` | Skills API | `confirmed-current` | Skill entry name / coordinate-like field. |
| skill `detail_url` | Skills API | `confirmed-current` | Relative URL. |
| skill `wasm_url` | Skills API | `confirmed-current` | Relative artifact URL. |
| skill `checksum_url` | Skills API | `confirmed-current` | Relative checksum URL. |
| skill `metadata.description` / `metadata.name` | Skills API | `confirmed-current` / `missing-value` | Current field; observed empty strings. |
| `runwasm_coordinate` | derived from skill/module/version/package and `moon runwasm` coordinate rules | `derived` | Legitimate SeekMoon output if generated from current skill/toolchain data. |
| `local_add`, `local_check`, `local_test`, `local_target_build` | local probe commands | `derived` | Current validation result, not upstream metadata. |
| `github_archived`, `github_pushed_at`, `github_tests_or_ci` | GitHub | `derived` | External maintenance signals. Must be absent unless GitHub checked. |
| `risk` | human/SeekMoon record | `derived` | Record/report field only. |
| `conclusion` | user/SeekMoon record | `derived` | Record/report field only; use controlled enum. |
| `match.fields` | SeekMoon search algorithm | `derived` | Machine-reproducibility field for JSON/record. Not upstream. |
| `rank` | SeekMoon search algorithm | `derived` | Search output ordering field. Not upstream. |
| `score` / `quality_score` / `Mooncake Score` | prior proposals | `unsupported-future` | Not current. Must not appear in v0 output. |
| `advisory`, `audit`, vulnerability status | prior proposals | `unsupported-future` | Not current. Must not appear in v0 output. |
| `outdated` | proposed command/field | `unsupported-future` as upstream field; `derived` only after dependency context exists | No current upstream command/source. Do not expose as v0 package field. |
| `dependents`, `reverse_dependencies` | prior proposals | `unsupported-future` | No current evidence source. |
| `provenance`, `signature`, `SBOM`, `publisher_identity` | prior proposals | `unsupported-future` | No current source. Checksum fields are not equivalent. |
| `docs_url` | constructed from module coordinate | `derived` | Legitimate derived URL, but docs content comes from assets/API. |
| `api_index_available` | asset fetch result | `derived` | Status of fetched module index. |

## Current Canonical Dictionary By Command Surface

The dictionary below is the v0 ingestion/output boundary. Fields outside it are unsupported for current output unless a command explicitly defines them later.

### `sync`

Purpose: create a dated local evidence snapshot.

| Field | Status | Source |
|---|---|---|
| `snapshot.fetched_at` | `derived` | SeekMoon runtime clock |
| `statistics.total_modules` | `confirmed-current` | Statistics API |
| `statistics.total_packages` | `confirmed-current` | Statistics API |
| `statistics.total_lines` | `confirmed-current` | Statistics API |
| `statistics.total_downloads` | `confirmed-current` | Statistics API |
| `modules_api.length` | `derived` | count of Modules API array |
| `local_index.file_count` | `confirmed-local` | local registry index |
| `local_index.head` | `confirmed-local` | local registry index git |
| `symbols.file_count` | `confirmed-local` | local symbols cache |
| `toolchain.moon_version` | `confirmed-local` | `moon --version` |
| `toolchain.moonc_version` | `confirmed-local` | `moon --version` |
| `toolchain.moonrun_version` | `confirmed-local` | `moon --version` |

### `search`

Purpose: produce candidate modules or skills from a query.

Library search fields:

| Field | Status | Source |
|---|---|---|
| `schema` | `derived` | SeekMoon output contract |
| `snapshot_id` / `snapshot.fetched_at` | `derived` | current or cached sync |
| `query.text` | `derived` | user input |
| `query.kind` | `derived` | command option, default `library` unless specified |
| `query.target` | `derived` | command option or project context when defined |
| `results[].rank` | `derived` | SeekMoon search algorithm |
| `results[].module` | `confirmed-current` | Modules API `name` |
| `results[].version` | `confirmed-current` | Modules API |
| `results[].description.status/value/source` | `confirmed-current` / `missing-value` | Modules API |
| `results[].keywords.status/value/source` | `confirmed-current` / `missing-value` | Modules API |
| `results[].license.status/value/source` | `confirmed-current` / `missing-value` | Modules API |
| `results[].repository.status/value/source` | `confirmed-current` / `missing-value` | Modules API |
| `results[].is_new` | `confirmed-current` | Modules API |
| `results[].created_at` | `confirmed-current` | Modules API |
| `results[].target.status/value/source` | `unknown-evidence` / `derived` | only when target requested; manifest/local validation needed |
| `results[].match.fields` | `derived` | SeekMoon matching algorithm |

Skill search fields:

| Field | Status | Source |
|---|---|---|
| `results[].module` | `confirmed-current` | Skills API |
| `results[].version` | `confirmed-current` | Skills API |
| `results[].package` | `confirmed-current` | Skills API |
| `results[].name` | `confirmed-current` | Skills API |
| `results[].author` | `confirmed-current` | Skills API |
| `results[].author_avatar` | `confirmed-current` | Skills API |
| `results[].description.status/value/source` | `confirmed-current` / `missing-value` | Skills API `metadata.description` |
| `results[].repository.status/value/source` | `confirmed-current` / `missing-value` | Skills API |
| `results[].created_at` | `confirmed-current` | Skills API |

### `view`

Purpose: show a module evidence profile.

| Field | Status | Source |
|---|---|---|
| `module` | `confirmed-current` | manifest `module` / `name` |
| `version` | `confirmed-current` | manifest |
| `latest_version` | `confirmed-current` | manifest |
| `downloads` | `confirmed-current` | manifest |
| `has_package` | `confirmed-current` | manifest |
| `build_status` | `confirmed-current` | manifest |
| `metadata.description.status/value/source` | `confirmed-current` / `missing-value` | manifest metadata |
| `metadata.keywords.status/value/source` | `confirmed-current` / `missing-value` | manifest metadata |
| `metadata.repository.status/value/source` | `confirmed-current` / `missing-value` | manifest metadata |
| `metadata.license.status/value/source` | `confirmed-current` / `missing-value` | manifest metadata |
| `metadata.deps` | `confirmed-current` | manifest metadata |
| `metadata.readme` | `confirmed-current` / `missing-value` | manifest metadata |
| `metadata.source` | `confirmed-current` / `missing-value` | manifest metadata |
| `metadata.checksum` | `confirmed-current` / `missing-value` | manifest metadata |
| `metadata.preferred_target` | `confirmed-current` / `missing-value` | manifest metadata `preferred-target` |
| `metadata.created_at` | `confirmed-current` | manifest metadata |
| `versions[].version` | `confirmed-current` | manifest |
| `versions[].yanked` | `confirmed-current` | manifest |
| `versions_count` | `derived` | count of manifest `versions` |
| `docs_url` | `derived` | module coordinate |
| `api_index.status` | `present` / `failed` / `unavailable` | module index asset fetch |
| `api_index.package_count` | `derived` | parsed module index |
| `target.status/value/source` | `unknown-evidence` / `derived` | manifest metadata and/or probe evidence |

### `api`

Purpose: inspect package API structure.

| Field | Status | Source |
|---|---|---|
| `module` | `confirmed-current` | selected module |
| `version` | `confirmed-current` | manifest |
| `package_path` | `confirmed-current` | module index |
| `module_index.status` | `present` / `failed` / `unavailable` | module index asset |
| `package_data.status` | `present` / `failed` / `unavailable` | package data asset |
| `types[].name` | `confirmed-current` | package data |
| `types[].docstring.status/value` | `confirmed-current` / `missing-value` | package data |
| `types[].signature` | `confirmed-current` | package data |
| `types[].loc.status/value` | `confirmed-current` / `missing-value` | package data |
| `types[].methods[]` | `confirmed-current` | package data |
| `values[]` | `confirmed-current` | package data |
| `traits[]` | `confirmed-current` | package data |
| `errors[]` | `confirmed-current` | package data |
| `typealias[]` | `confirmed-current` | package data |
| `misc[]` | `confirmed-current` | package data |
| `ide_doc.status/output` | `confirmed-local` / `failed` / `unavailable` | `moon ide doc` when used |

### `source`

Purpose: fetch or locate published source.

| Field | Status | Source |
|---|---|---|
| `module` | `confirmed-current` | selected module |
| `version` | `confirmed-current` | manifest or user coordinate |
| `moon_fetch.status/path/error` | `confirmed-local` / `failed` / `unavailable` | `moon fetch` |
| `source_zip.status/path/error` | `confirmed-current` / `failed` / `unavailable` | source zip endpoint |
| `local_cache.status/path` | `confirmed-local` / `unavailable` | local cache |
| `core_local_source.status/path` | `confirmed-local` / `unavailable` | `~/.moon/lib/core` for core |
| `selected_source.method/path` | `derived` | SeekMoon source resolution |
| `files.summary` | `derived` | fetched/local source directory listing |

### `skill`

Purpose: discover and inspect executable Wasm / skill entries.

| Field | Status | Source |
|---|---|---|
| `module` | `confirmed-current` | Skills API |
| `author` | `confirmed-current` | Skills API |
| `author_avatar` | `confirmed-current` | Skills API |
| `version` | `confirmed-current` | Skills API |
| `package` | `confirmed-current` | Skills API |
| `name` | `confirmed-current` | Skills API |
| `detail_url` | `confirmed-current` | Skills API |
| `wasm_url` | `confirmed-current` | Skills API |
| `checksum_url` | `confirmed-current` | Skills API |
| `metadata.name.status/value` | `confirmed-current` / `missing-value` | Skills API |
| `metadata.description.status/value` | `confirmed-current` / `missing-value` | Skills API |
| `repository.status/value` | `confirmed-current` / `missing-value` | Skills API |
| `created_at` | `confirmed-current` | Skills API |
| `runwasm_coordinate` | `derived` | Skills API + `moon runwasm` coordinate rules |
| `wasm_asset.status` | `present` / `failed` / `unavailable` | asset fetch/check when performed |
| `checksum_asset.status` | `present` / `failed` / `unavailable` | asset fetch/check when performed |

### `compare`

Purpose: put multiple candidates on one evidence surface.

| Field | Status | Source |
|---|---|---|
| `candidates[]` | `derived` | selected search/session/modules |
| `module`, `version`, `description`, `license`, `repository` | mixed | modules/manifest |
| `downloads`, `build_status`, `versions_count`, `has_package` | `confirmed-current` / `derived` | manifest |
| `package_count`, `api_index.status` | `derived` | module index |
| `source.status` | `derived` | source command evidence |
| `target.status` | `unknown-evidence` / `derived` | metadata/probe |
| `probe.status` | `derived` | probe evidence |
| `github.status` | `derived` / `unavailable` | GitHub enrichment when performed |

### `probe`

Purpose: validate a candidate in local project/probe context.

| Field | Status | Source |
|---|---|---|
| `module` | `confirmed-current` | selected candidate |
| `version` | `confirmed-current` / `derived` | manifest/user coordinate |
| `target` | `derived` | user option/project context |
| `probe_path` | `derived` | SeekMoon |
| `moon_new.status` | `derived` | local command |
| `moon_add.status` | `derived` | local command |
| `moon_check.status` | `derived` | local command |
| `moon_test.status` | `derived` | local command |
| `moon_check_target.status` | `derived` | local command |
| `moon_build_target.status` | `derived` | local command |
| `result.status` | `derived` | SeekMoon aggregation |
| `logs.path` | `derived` | SeekMoon |

### `record`

Purpose: store adoption decision and evidence references.

| Field | Status | Source |
|---|---|---|
| `record_id` | `derived` | SeekMoon |
| `created_at` | `derived` | SeekMoon runtime clock |
| `module` | `confirmed-current` | selected candidate |
| `version` | `confirmed-current` | selected candidate |
| `conclusion` | `derived` | user/SeekMoon controlled enum |
| `note` | `derived` | user input |
| `evidence_refs[]` | `derived` | snapshot/session/probe/source ids |
| `not_confirmed[]` | `derived` | user/SeekMoon audit text |

Allowed conclusion values should remain the SOP set, normalized if needed: `adopt`, `adopt-with-adapter`, `continue-verification`, `contribute-upstream`, `fork`, `build-own`, `reject-for-now`.

### `report`

Purpose: render a reproducible audit record.

| Field | Status | Source |
|---|---|---|
| `goal` | `derived` | user/report input |
| `date` | `derived` | SeekMoon/runtime |
| `toolchain` | `confirmed-local` | sync/doctor |
| `data_sources[]` | mixed | sources actually used |
| `query` | `derived` | search/session |
| `candidates[]` | mixed | search/view/compare |
| `local_validation[]` | `derived` | probe |
| `cannot_confirm[]` | `derived` | audit judgment |
| `conclusion` | `derived` | record/report |

### `raw`

Purpose: expose original source payloads for audit escape hatches.

Allowed raw surfaces:

| Surface | Status | Source |
|---|---|---|
| `raw modules` | `confirmed-current` | Modules API |
| `raw statistics` | `confirmed-current` | Statistics API |
| `raw manifest <module>` | `confirmed-current` | Manifest API |
| `raw module-index <module@version>` | `confirmed-current` / `failed` / `unavailable` | module index asset |
| `raw package-data <module@version> <package>` | `confirmed-current` / `failed` / `unavailable` | package data asset |
| `raw resources <module@version> <package>` | `missing-value` / `unavailable` | resources asset |
| `raw skills` | `confirmed-current` | Skills API |
| `raw local-index <module>` | `confirmed-local` | local registry index |
| `raw symbols <query/module>` | `confirmed-local` / `unknown-evidence` | local symbols cache |

## Unsupported Future Fields That Must Not Appear In Current Output

These fields may be future watch items, but they are not current v0 output fields:

| Unsupported field/source | Reason |
|---|---|
| `quality_score`, `mooncake_score`, `score` as quality metric | No current upstream field or defined SeekMoon scoring contract. Search `rank` is allowed only as search algorithm output. |
| `advisory`, `audit_status`, `vulnerabilities`, `security_score` | No current advisory/audit source. |
| `outdated` package status | No current upstream command/source. It may later be derived from project deps + manifest, but is not v0 canonical ingestion. |
| `dependents`, `reverse_dependencies`, `reverse_dependency_count` | No current source. |
| `provenance`, `attestation`, `signature`, `publisher_identity`, `verified_publisher` | No current source. Registry checksums and skill checksum URLs are not provenance. |
| `sbom` | No current source. |
| `ci_status` as registry fact | Only GitHub-derived if checked; not upstream Mooncakes field. |
| `tests_present`, `examples_present` as registry fact | Only source/GitHub-derived if checked; not upstream field. |
| `docs_build_status` separate from `build_status` | Current manifest has `build_status`; do not invent separate docs build status without evidence. |
| `server_search_url`, `search_deeplink`, `/search?q=` | Not current upstream source. |
| `target_supported=true/false` without evidence | Only `unknown`, metadata-derived, or probe-derived status is allowed. |
| `resources_available=false` as package absence | `resource.json` 404 only means resource asset unavailable. |
| `moon search`, `moon view`, `moon audit`, `moon outdated` as current toolchain commands | Current `moon` help does not expose them. |
| `--why`, `--hints`, `guide`, top-level `schema/shape/fields` command surfaces | Latest design explicitly deleted them as unsupported surfaces. Do not reintroduce into v0 dictionary. |

## Fields Requiring `missing` vs `unknown` Distinction

Use `missing` when the current schema/source has the field position and the value is absent, empty, or null. Use `unknown` when the current question asks for a fact that current evidence cannot determine. Use `failed` when an attempted fetch/command failed. Use `unavailable` when a source class is not available for the selected object or environment.

| Field | `missing` case | `unknown` case | Other status cases |
|---|---|---|---|
| `description` | field absent, null, or empty in current source | not applicable after source fetched | source fetch `failed` if source not retrieved |
| `keywords` | absent or empty array | not applicable after source fetched |  |
| `repository` | absent or empty | GitHub maintenance state unknown until checked | GitHub fetch `failed` / GitHub unavailable |
| `license` | absent or empty in modules/manifest/local index | license compatibility unknown until policy applies |  |
| `readme` metadata path | absent/empty in metadata | README body unknown until resources/source fetched | resources/source `failed` / `unavailable` |
| `checksum` registry metadata | absent/empty in metadata | integrity interpretation unknown without verification policy |  |
| `preferred_target` | absent from metadata | support for requested target remains unknown |  |
| `targets` / `supported_targets` | field absent from metadata | requested target compatibility unknown |  |
| `target.status` | not a missing value unless target field exists but empty | no current metadata/probe determines support | local target probe can be `failed` |
| `downloads` | not available until manifest fetched | unknown before manifest fetch | manifest fetch `failed` |
| `build_status` | absent in manifest payload | unknown before manifest fetch | manifest fetch `failed` |
| `versions` | absent in manifest payload | unknown before manifest fetch | manifest fetch `failed` |
| API `docstring` | null/empty in package data | unknown before package data / ide doc | package data fetch `failed` / `unavailable` |
| API `loc` | absent/null in package data | unknown before package data | package data fetch `failed` |
| `module_index.status` | not applicable | unknown before attempted fetch | `present`, `failed`, `unavailable` |
| `package_data.status` | not applicable | unknown before attempted fetch | `present`, `failed`, `unavailable` |
| `resources.status` | resource file absent/404 | unknown before attempted fetch | `failed`, `unavailable` |
| `source_fetch.status` | not applicable | unknown before source command/fetch | `present`, `failed`, `unavailable` |
| skill `metadata.description` | empty string in Skills API | not applicable after source fetched |  |
| skill `wasm_asset.status` | not applicable | unknown before asset check | `present`, `failed`, `unavailable` |
| skill `checksum_asset.status` | not applicable | unknown before asset check | `present`, `failed`, `unavailable` |
| GitHub maintenance fields | repository URL missing means GitHub enrichment unavailable | repo status unknown before GitHub check | GitHub check `failed` / repo `unavailable` |
| local symbols result | not applicable | no hit means unknown because cache is partial | symbols cache unavailable |

## Drift Findings

1. The modules list schema currently includes `is_new` and `created_at`; earlier dictionary sketches did not consistently include them.
2. The current modules list does not include `downloads` or `build_status`; those are manifest-level fields. Search output that displays them must either fetch manifests or mark them outside list-only mode.
3. Manifest metadata currently includes `preferred-target` in the checked package, not a proven general `supported-targets` matrix. Target compatibility remains unknown unless metadata or probe evidence exists.
4. `moon ide doc` is current, but only for current module/core/registry symbol indexes. It must not be treated as package discovery.
5. `moon fetch` is current but explicitly unstable. SeekMoon may use it as a source retrieval method with status modeling, not as a stable upstream contract.
6. Skills API is current and belongs to an execution-object surface. It must remain separate from library package adoption.
7. `resource.json` is a current asset class but not guaranteed. 404 must be `unavailable`, not package absence.
8. The design document contains superseded proposals (`--why`, `--hints`, top-level schema/shape commands, `fields`). The final design deletes them; the canonical dictionary follows the deletion.
9. Prior notes mention quality score, audit, outdated, reverse dependents, provenance, and search deep links as desired ecosystem capabilities. None are current v0 ingestion fields.

## Temporally Unstable Items To Revalidate Later

Revalidate these before final implementation or public reporting:

| Item | Why unstable |
|---|---|
| Modules API item keys and counts | Registry changes continuously. |
| Statistics API counts | Counts change with publishes/downloads. |
| Manifest metadata keys | Package authors and registry schema may add/remove metadata. |
| Skills API shape and count | Skills marketplace is active and may still be experimental. |
| `moon --help` command surface | Toolchain releases may add package discovery commands or change `fetch` stability. |
| Local registry index count and HEAD | Depends on latest `moon update`. |
| Local symbols cache coverage | Depends on toolchain/cache update behavior. |
| Source zip availability | Varies by module/version. |
| `resource.json` availability | Varies by package/version. |
| GitHub maintenance fields | Repository state, CI, issues, archived state, pushed time are external and mutable. |
| Target metadata | `preferred-target`, `targets`, and `supported-targets` may evolve as Mooncakes target support matures. |

## Final Canonical Boundary

SeekMoon v0 ingestion should consume current Mooncakes registry/API/assets, local Moon toolchain/cache/index, current project context, optional GitHub maintenance enrichment, local probe results, and SeekMoon records.

SeekMoon v0 output must separate:

- upstream facts: current API/assets/manifest/skills fields;
- local facts: toolchain, registry index, symbols cache, probe commands;
- derived facts: rank, match fields, docs URLs, source resolution, version counts, runwasm coordinates, records;
- unknown facts: questions current evidence cannot answer;
- missing values: current schema/source fields whose values are absent;
- unsupported future fields: desired ecosystem capabilities with no current source.

The canonical command surfaces are `sync`, `search`, `view`, `api`, `source`, `skill`, `compare`, `probe`, `record`, `report`, and `raw`. No other public field or command surface is required for the current data model.
