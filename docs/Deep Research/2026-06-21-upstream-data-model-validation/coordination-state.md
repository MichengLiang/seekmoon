# SeekMoon Upstream Data Model Validation Coordination State

## Objective

Validate the current upstream data model for SeekMoon package discovery inputs, especially whether fields proposed in prior SOP/design documents are current, stale, complete, or unsupported.

## Batch

- Date: 2026-06-21
- Host: `/home/t103o/workbench/projects/seekmoon/docs/Deep Research/2026-06-21-upstream-data-model-validation`
- Tasks: `Tasks/`
- Reports: `reports/`

## Required Source Materials

- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 包发现闭环动向调查.md`
- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 语言、工具链与包生态调查报告.md`
- `/home/t103o/workbench/notes/关于MoonBit/关于包搜索与发现/MoonBit包搜索发现与生态查重SOP.md`
- `/home/t103o/workbench/projects/seekmoon/docs/关于moonbit/002-raw消费者侧 MoonBit 包发现 CLI 设计.md`

## Resource Map

| Resource | Role |
|---|---|
| Mooncakes `/api/v0/modules` | Module candidate list source |
| Mooncakes `/api/v0/modules/statistics` | Registry scale snapshot source |
| Mooncakes `/api/v0/manifest/<module>` | Single module metadata source |
| Mooncakes `/api/v0/skills` | Wasm/skill entry source |
| Mooncakes `assets/.../module_index.json` | Module package/API index source |
| Mooncakes `assets/.../package_data.json` | Package API data source |
| Mooncakes `assets/.../resource.json` | Package resource source when available |
| Mooncakes source zip | Published source retrieval source |
| Local `~/.moon/registry/index` | Local registry index and raw publication source |
| Local `~/.moon/registry/symbols` | Local symbols cache source |
| Local `moon` commands | Toolchain behavior source |
| `external/mooncakes.io` | Official frontend implementation evidence |
| Existing SOP/design docs | Claimed field/data-source inventory |

## Role Registry

| Agent | Task | Expected Report |
|---|---|---|
| Nietzsche | Task 01 public API schema validation | completed: `reports/01-public-api-schema-validation.md` |
| Noether | Task 02 assets and source schema validation | completed: `reports/02-assets-and-source-schema-validation.md` |
| Popper | Task 03 local toolchain and registry validation | completed: `reports/03-local-toolchain-and-registry-validation.md` |
| Mencius | Task 04 frontend source validation | completed: `reports/04-mooncakes-frontend-source-validation.md` |
| Russell | Task 05 canonical dictionary and drift audit | completed: `reports/05-canonical-dictionary-and-drift-audit.md` |

## Promotion Boundary

A field or source may be treated as current accepted input only when at least one report provides direct evidence and the coordinator can map it to a current SeekMoon command/data object.

Unsupported future fields must not enter current output schemas.

Fields with current schema membership but absent values require explicit status modeling.

## Report Review Checklist

- Does the report include fetch time or source snapshot?
- Does it distinguish live API facts from local facts and frontend implementation facts?
- Does it give enough schema detail to implement ingestion without guessing field names?
- Does it classify existing assumptions as current/stale/incomplete/wrong or equivalent?
- Does it identify failure semantics, not only success examples?

## Coordinator Review Result

All five reports were completed and read by the coordinator.

Accepted current inputs:

- `/api/v0/modules`
- `/api/v0/modules/statistics`
- `/api/v0/manifest/<owner>/<module>`
- `/api/v0/skills`
- `module_index.json`
- `package_data.json`
- source zip when available
- local registry index JSONL
- local symbols JSONL
- local `moon` command surfaces: `update`, `ide doc`, `fetch`, `add`, `check`, `test`, `build --target`, `runwasm`
- local Moon cache and current project context
- optional GitHub maintenance enrichment
- SeekMoon probe/session/record/report artifacts

Rejected as current output fields:

- quality score / Mooncake Score
- advisory / audit / vulnerability status
- reverse dependents
- provenance / attestation / publisher identity / SBOM
- server-side search endpoint
- target support booleans without metadata or probe evidence
- old deleted CLI surfaces: `--why`, `--hints`, `guide`, top-level `schema/shape/fields`

Important unresolved drift:

- Frontend source currently requests singular `resource.json`, while older SPA note and Task 02 checked plural `resources.json`. Task 02 found plural `resources.json` unavailable in all checked cases; follow-up checks found singular `resource.json` live on sampled modules. Treat plural as stale and singular as the current optional resource path.
