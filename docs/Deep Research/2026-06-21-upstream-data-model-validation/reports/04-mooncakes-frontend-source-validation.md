# Task 04: Mooncakes Frontend Source Validation

## Scope And Snapshot

- Source repo: `/home/t103o/workbench/external/mooncakes.io`
- Upstream: `git@github.com:moonbitlang/mooncakes.io.git`
- Update command: `git -C external/mooncakes.io pull --ff-only`
- Fetch/update time: `2026-06-21T20:10:24+08:00`
- Result: already up to date.
- Commit: `4d889ad9a49ed0e0b8b4d8dc58451490de4bb38b`
- Commit date: `2026-06-18 18:52:20 +0800`
- Commit subject: `Refine skills marketplace hero copy`

This report records frontend implementation facts from the Mooncakes source. It does not treat frontend behavior as a stable public API contract.

## API And Asset Call Sites

| Endpoint / asset | Frontend source references | Frontend role |
|---|---|---|
| `/api/v0/modules?recent=12&new=12&popular=12` | `src/page/home/state.mbt:122-130` | Home highlight lists: recent, new, popular. |
| `/api/v0/modules/statistics` | `src/page/home/state.mbt:130-132` | Home counters. `src/view/footer.mbt:19` links to raw statistics. |
| `/api/v0/modules` | `src/page/home/state.mbt:149-153`, loaded on non-empty search or “show more” at `src/page/home/state.mbt:186-215` | Full client-side module list and search index. |
| `/api/v0/manifest/{path}` | `src/page/docs/state.mbt:146-156` | Docs/detail page metadata and version resolution. |
| `/api/v0/builds/log?...` | `src/page/docs/state.mbt:258-267` | Build log panel when docs build failed. |
| `/assets/{asset_path}/resource.json` | `src/page/docs/state.mbt:116-137`, reload on package navigation at `src/page/docs/state.mbt:553-563` | README content and source file list. The code uses singular `resource.json`; the older SPA note in `src/index.html:33` mentions `resources.json`, so SeekMoon must treat the plural form as stale. |
| `/assets/{module}@{version}/module_index.json` | `src/page/docs/state.mbt:121-129` | Package tree, sidebar, docs search index. |
| `/assets/{asset_path}/package_data.json` | `src/page/docs/state.mbt:130-133`, `src/page/docs/state.mbt:553-563` | Per-package API docs. Loaded only when `manifest.has_package` permits it. |
| `/api/v0/skills` | `src/page/skills/state.mbt:70-74` | Skills marketplace list. |
| `/api/v0/skills/{entry}` | `src/page/skills/state.mbt:86-91` | Skill detail object. |
| `/assets/{entry}/SKILL.md` | `src/page/skills/state.mbt:86-93` | Skill detail README/SKILL content. |
| Skill wasm/checksum asset URLs | `src/page/skills/view.mbt:936-941` | Download links from API-provided `wasm_url` and `checksum_url`. |

## Frontend-Consumed Module Fields

The home module decoder consumes these fields from `/api/v0/modules` and highlight arrays:

| Field | Source references | Frontend use |
|---|---|---|
| `name` | Decoded as full path at `src/page/home/state.mbt:424-477`; split into `author` and short `name` at `src/page/home/state.mbt:455-464` | Display, docs links, user links, client search. |
| `version` | `src/page/home/state.mbt:430-431`, rendered at `src/page/home/view.mbt:56`, `src/page/home/view.mbt:109`, `src/page/home/view.mbt:166-170` | Display and versioned docs path, except core/empty version at `src/page/home/view.mbt:16-21`. |
| `keywords` | `src/page/home/state.mbt:432`, `src/page/home/state.mbt:439-446`; rendered at `src/page/home/view.mbt:68-79` | Display, client search, and filtering `mooncakes-test` modules at `src/page/home/state.mbt:524-530`. |
| `description` | `src/page/home/state.mbt:433`, `src/page/home/state.mbt:447`; rendered at `src/page/home/view.mbt:45-48`, `src/page/home/view.mbt:81-85`, `src/page/home/view.mbt:182-185` | Display only. Not in home search index. |
| `created_at` | `src/page/home/state.mbt:436`, `src/page/home/state.mbt:448`; rendered at `src/page/home/view.mbt:56-60`, `src/page/home/view.mbt:86-88`; recency penalty at `src/page/home/state.mbt:92-98`, `src/page/home/state.mbt:394-398` | Relative-time display and search ranking penalty. |
| `is_new` | `src/page/home/state.mbt:434-454`; rendered at `src/page/home/view.mbt:50-55` | “NEW” badge. |

Not consumed by the home module model: per-module `downloads`, `license`, `repository`, `build_status`, targets/target metadata, author avatar, quality score, dependents, advisory/provenance data.

## Frontend Search Behavior

Home module search is client-side:

- Initial page load fetches highlight buckets and statistics, not the full module list: `src/page/home/state.mbt:122-133`.
- The full `/api/v0/modules` list is fetched only when a filter is entered or “show more modules” is opened: `src/page/home/state.mbt:186-215`.
- `build_search_index` indexes lowercased `path`, short `name`, `author`, `keywords`, and `recency_penalty`: `src/page/home/state.mbt:377-400`.
- It does not index `description`.
- Primary token matching checks path/name/author/keywords: `src/page/home/state.mbt:249-307`.
- Fuzzy scoring is applied to path/name/author/keywords, with keyword penalty `200`: `src/page/home/state.mbt:310-367`.
- Recency penalty is added to fuzzy score: `src/page/home/state.mbt:364-367`.
- Results are sorted by score and capped at `200`: `src/page/home/state.mbt:370-373`.

Docs page symbol search is also client-side after `module_index.json` loads:

- Search entries are built from packages, typealiases, traits, errors, types, methods, misc, and values: `src/page/docs/search.mbt:36-119`.
- Fuzzy matching scores `entry.fullname`: `src/page/docs/search.mbt:122-139`.

Skills marketplace filtering is client-side:

- `/api/v0/skills` returns all entries to the page: `src/page/skills/state.mbt:70-74`.
- Filter haystack is `entry.name`, `entry.module_`, `entry.package_path`, and `entry.description`: `src/page/skills/state.mbt:296-323`.

## Frontend-Consumed Manifest Fields

The docs page manifest decoder and metadata view consume these fields:

| Field | Source references | Frontend use |
|---|---|---|
| top-level `name` | `src/page/docs/module_index.mbt:247-264` | Module/package identity and install command via `MetaInfo.name`. |
| top-level `module` | `src/page/docs/module_index.mbt:247-264`; used at `src/page/docs/state.mbt:193-207`, `src/page/docs/state.mbt:414-430` | Resolve package path, versioned docs paths, build log query, docs search links. |
| top-level `version` | `src/page/docs/module_index.mbt:247-264`; `src/page/docs/view.mbt:335-339` | Versioned docs/assets, install command, source zip URL, display. |
| top-level `has_package` | `src/page/docs/module_index.mbt:247-264`; asset gating at `src/page/docs/state.mbt:130-136`, docs notice at `src/page/docs/view.mbt:400-407` | Controls whether `package_data.json` should load and whether missing package data is a docs failure. |
| top-level `downloads` | `src/page/docs/resource.mbt:202-210`; rendered at `src/page/docs/metainfo.mbt:154-158`, `src/page/docs/metainfo.mbt:272-276` | Display downloads, except hidden for `moonbitlang/core` or zero. |
| top-level `build_status` | `src/page/docs/resource.mbt:202-218`; docs state at `src/page/docs/view.mbt:377-411` | Docs availability state and build notices. |
| `versions[].version` + `versions[].build_status` | `src/page/docs/resource.mbt:131-151`, `src/page/docs/resource.mbt:211-218` | Prefer current-version build status over metadata/top-level status. |
| `metadata.name` | `src/page/docs/resource.mbt:190-230` | Display identity, author derivation, install command, zip URL. |
| `metadata.version` | `src/page/docs/resource.mbt:190-230` | Decoded into metadata, but view later overrides with top-level manifest version at `src/page/docs/view.mbt:335-339`. |
| `metadata.repository` | `src/page/docs/resource.mbt:190-230`; normalized at `src/page/docs/resource.mbt:44-51`, `src/page/docs/resource.mbt:172-179`; rendered at `src/page/docs/metainfo.mbt:140-152`, `src/page/docs/metainfo.mbt:256-270` | Repository link. |
| `metadata.license` | `src/page/docs/resource.mbt:190-230`; rendered at `src/page/docs/metainfo.mbt:132-139`, `src/page/docs/metainfo.mbt:248-255` | License display when present. |
| `metadata.description` | `src/page/docs/resource.mbt:190-230`; rendered at `src/page/docs/metainfo.mbt:99-105`, `src/page/docs/metainfo.mbt:197-205` | Description display. |
| `metadata.keywords` | `src/page/docs/resource.mbt:76-90`, `src/page/docs/resource.mbt:190-230`; rendered at `src/page/docs/metainfo.mbt:107-119`, `src/page/docs/metainfo.mbt:206-209` | Keyword chips. |
| `metadata.deps` | `src/page/docs/resource.mbt:94-112`, `src/page/docs/resource.mbt:190-230`; rendered at `src/page/docs/metainfo.mbt:163-187`, `src/page/docs/metainfo.mbt:281-305` | Dependency list with docs links. |
| `metadata.build_status` | `src/page/docs/resource.mbt:190-230` | Fallback build status if `versions` and top-level status do not supply it. |

No manifest target metadata is rendered or used in the docs page in this source snapshot, although `metadata.targets` / `supported-targets` remain useful candidates for live API validation because existing SeekMoon assumptions include target compatibility.

## Docs Asset Field Usage

### `resource.json`

Decoded at `src/page/docs/resource.mbt:257-268`:

- `readme_content`: rendered by the docs page at `src/page/docs/view.mbt:420-438`.
- `source_files`: rendered as source file links at `src/page/docs/sidebar.mbt:239-263`.

### `module_index.json`

Decoded at `src/page/docs/module_index.mbt:15-31`, `src/page/docs/module_index.mbt:439-493`.

Consumed fields:

- Root/tree: `name`, `package`, `childs`.
- Package index: `path`, `traits`, `errors`, `types`, `typealias`, `values`, `misc`.
- Type/error/misc index entries: `name`, `impls`, `methods`.
- Trait entries: `name`, `impls`.
- Impl entries: `self`, `trait`, `methods`.
- Type references: `kind`, `constructor`, `arguments`, `parameters`, `return_type`, `error_type`, `is_async`, `name`, `path`.

Frontend uses this asset for:

- Sidebar package/symbol tree: `src/page/docs/sidebar.mbt:32-71`, `src/page/docs/sidebar.mbt:97-172`.
- Breadcrumb validation: `src/page/docs/breadcrumbs.mbt:29-33`.
- Docs search entries: `src/page/docs/search.mbt:36-119`.
- Versioned docs links: `src/page/docs/module_index.mbt:185-244`.

### `package_data.json`

Decoded at `src/page/docs/package_data.mbt:16-24`, `src/page/docs/package_data.mbt:238-272`.

Consumed fields:

- Top-level: `name`, `types`, `traits`, `errors`, `typealias`, `values`, `misc`.
- Docs: `name`, `docstring`, `signature`, `loc`, `methods`, `impls`.
- `loc`: `path`, `file`, `line`, `column`.
- Impl docs: `self`, `trait`, `methods`.

Frontend uses this asset for rendered API docs through `src/page/docs/cards.mbt:466-507`, source-code links at `src/page/docs/cards.mbt:249-268`, and collapse state setup at `src/page/docs/state.mbt:336-355`.

## Skill Marketplace Field Usage

Skills list/detail model fields are defined at `src/page/skills/state.mbt:30-43` and decoded at `src/page/skills/state.mbt:327-377`.

| Skill field | Source references | Frontend use |
|---|---|---|
| `module` | decoded into `module_` at `src/page/skills/state.mbt:330`; used in paths at `src/page/skills/state.mbt:46-67` | Skill path, versioned runwasm command, asset paths. |
| `author` | `src/page/skills/state.mbt:331`; rendered at `src/page/skills/view.mbt:433-452`, `src/page/skills/view.mbt:858-869` | Avatar fallback and user link. |
| `author_avatar` | optional at `src/page/skills/state.mbt:332`, `src/page/skills/state.mbt:344-347` | Avatar image if present. |
| `version` | `src/page/skills/state.mbt:333`; run command at `src/page/skills/state.mbt:55-67`; metadata at `src/page/skills/view.mbt:922` | Versioned `moon runwasm` coordinate and detail metadata. |
| `package` | decoded as `package_path` at `src/page/skills/state.mbt:334`; used at `src/page/skills/state.mbt:46-67` | Skill path, run command, asset path, detail package display. |
| `name` | `src/page/skills/state.mbt:335`; rendered at `src/page/skills/view.mbt:460-461`, `src/page/skills/view.mbt:1003-1015` | Skill display name / wasm filename. |
| `detail_url` | `src/page/skills/state.mbt:336`; used at `src/page/skills/view.mbt:488` | Skill card link. |
| `wasm_url` | `src/page/skills/state.mbt:337`; download link at `src/page/skills/view.mbt:939` | Download wasm. |
| `checksum_url` | `src/page/skills/state.mbt:338`; download link at `src/page/skills/view.mbt:940` | Download checksum. |
| `repository` | optional at `src/page/skills/state.mbt:339`, normalized for link at `src/page/skills/view.mbt:899-920` | Repository metadata row if present. |
| `metadata.description` | `src/page/skills/state.mbt:340-358`; rendered at `src/page/skills/view.mbt:514-520`, `src/page/skills/view.mbt:1016-1022` | Card/detail description and skill filter haystack. |
| `created_at` | `src/page/skills/state.mbt:341`; rendered at `src/page/skills/view.mbt:501-504`, `src/page/skills/view.mbt:930-933` | Relative publish time. |
| `SKILL.md` asset text | `src/page/skills/state.mbt:86-93`; rendered at `src/page/skills/view.mbt:794-810`, copied at `src/page/skills/view.mbt:947-964` | Detail documentation. Frontmatter is stripped at `src/page/skills/view.mbt:744-790`. |

Skill search/filter fields are exactly `name`, `module`, `package`, and `description`: `src/page/skills/state.mbt:296-323`.

## Existing Design Assumptions Assessment

### Correct

- `/api/v0/modules?search=<query>` should not be treated as a server-side search contract. The frontend fetches full `/api/v0/modules` and filters locally: `src/page/home/state.mbt:149-153`, `src/page/home/state.mbt:186-215`, `src/page/home/state.mbt:231-374`.
- Home package search should index `path/name`, author, and keywords; it should not assume description is part of the official frontend search index.
- Recency affects frontend ranking through `created_at`, but only as a penalty in local fuzzy scoring: `src/page/home/state.mbt:92-98`, `src/page/home/state.mbt:394-398`.
- Skills are a separate execution object from library module discovery. The frontend has a separate `/skills` model and uses `moon runwasm` coordinates: `src/page/skills/state.mbt:55-67`, `src/page/skills/view.mbt:694-719`.
- Docs automation should consume assets rather than scrape rendered SPA HTML. `src/index.html:21-37` explicitly lists API/assets for agents, and frontend source consumes those assets directly.

### Stale Or Incomplete

- If SeekMoon’s module search model includes `description` as a default search field, that is broader than the official frontend search behavior. It may be a SeekMoon product choice, but it should be labeled as SeekMoon search semantics, not Mooncakes frontend parity.
- If SeekMoon assumes target metadata is frontend-visible, this source snapshot does not support that. Target metadata still requires live API/source validation outside frontend usage.
- If SeekMoon assumes per-module downloads are available in `/api/v0/modules`, this frontend does not decode them. Downloads are displayed on docs pages from manifest, while “Most downloaded” is a server-provided highlight bucket.
- If SeekMoon assumes build status is part of module-list cards/search, this frontend does not use it there. Build status is docs-page metadata and build-log gating.
- If SeekMoon assumes `resources.json`, that is wrong for current frontend source. The current frontend source requests singular `resource.json`; the SPA note in `src/index.html:33` is stale and refers to the plural form.

### Wrong If Treated As Frontend Facts

- “Description participates in official frontend module search” is wrong for this snapshot.
- “Search is server-side” is wrong for module search, docs symbol search, and skills filtering.
- “Frontend list search uses target metadata” is wrong for this snapshot.
- “Frontend list cards use license/repository/build status/downloads” is wrong for this snapshot.
- “Skill marketplace entries are ordinary package modules” is wrong. The frontend treats them as Wasm executable entries with `wasm_url`, `checksum_url`, `SKILL.md`, and `moon runwasm`.

## Live API Facts Vs Frontend Implementation Facts

Frontend implementation facts from this source:

- Which fields are decoded and rendered by Mooncakes frontend.
- Which fields enter frontend search/filter indexes.
- Which endpoints/assets the SPA requests.
- Which missing fields are tolerated by frontend decoders.
- How docs assets are linked from frontend state.

Facts still requiring live API validation:

- Current JSON shape for `/api/v0/modules`, `/api/v0/modules/statistics`, `/api/v0/manifest/*`, `/api/v0/skills`, `/api/v0/skills/*`.
- Whether singular `resource.json` is available for current deployed docs assets, and how often it falls back to source zip or other sources. The plural `resources.json` should be treated as stale unless a future release reintroduces it.
- Whether target metadata appears in manifest/list/index data even though frontend does not consume it.
- Whether optional manifest fields are consistently present or nullable across real modules.
- Whether live skills entries always include `wasm_url`, `checksum_url`, `created_at`, and `metadata.description`.
- Whether source zip URLs match the frontend-constructed `https://download.mooncakes.io/user/{name}/{version}.zip` for all module classes.

## SeekMoon Ingestion Alignment

SeekMoon should align its base ingestion model with these frontend-consumed objects:

1. `ModuleSummary`: full module path from `name`, derived author/short name, version, keywords, description, created_at, is_new.
2. `ModuleStatistics`: total modules, packages, lines, downloads.
3. `Manifest`: top-level name/module/version/has_package/downloads/build_status/versions plus metadata name/version/repository/license/description/keywords/deps/build_status.
4. `ResourceAsset`: readme_content and source_files.
5. `ModuleIndexAsset`: package tree and symbol summaries.
6. `PackageDataAsset`: API docs with docstrings, signatures, locs, methods, impls.
7. `SkillEntry`: module, author, author_avatar, version, package, name, detail_url, wasm_url, checksum_url, repository, metadata.description, created_at, SKILL.md text.

SeekMoon may intentionally provide richer behavior than Mooncakes frontend, but those choices should be named as SeekMoon semantics. Frontend source supports trusting the above fields as actively consumed by the official site, not as a complete or stable public API contract.
