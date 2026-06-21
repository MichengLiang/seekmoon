# Task 01: Public API Schema Validation

## Role

You are a reconnaissance agent validating the current public Mooncakes HTTP API data model for SeekMoon.

You are not designing the CLI. You are validating upstream facts.

## Scope

Validate the current schema, field presence, nullability, and observed drift for these endpoints:

- `https://mooncakes.io/api/v0/modules`
- `https://mooncakes.io/api/v0/modules/statistics`
- `https://mooncakes.io/api/v0/manifest/<owner>/<module>`
- `https://mooncakes.io/api/v0/skills`

Use current live responses on 2026-06-21. Record exact fetch time.

## Required Background Reading

Read these files before investigating:

- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 包发现闭环动向调查.md`
- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 语言、工具链与包生态调查报告.md`
- `/home/t103o/workbench/notes/关于MoonBit/关于包搜索与发现/MoonBit包搜索发现与生态查重SOP.md`
- `/home/t103o/workbench/projects/seekmoon/docs/关于moonbit/002-raw消费者侧 MoonBit 包发现 CLI 设计.md`

## Questions To Answer

1. What is the exact top-level shape of each endpoint response?
2. For `/api/v0/modules`, what fields appear across all module objects? Which are always present, sometimes absent, null, empty string, empty array, or nested?
3. Does `/api/v0/modules?search=<query>` still return the unfiltered full list?
4. For `/api/v0/manifest/<owner>/<module>`, what fields appear for representative modules?
5. Are `metadata.targets`, `metadata.supported-targets`, `targets`, or `supported_targets` actually present in current manifest responses? If yes, in what spelling and shape?
6. For `/api/v0/skills`, what fields appear across skill objects? Are `wasm_url`, `checksum_url`, `detail_url`, `package`, and `metadata` stable?
7. Which fields claimed in the existing SOP/design are current, stale, incomplete, or wrong?

## Representative Samples

At minimum inspect:

- `moonbitlang/core`
- `mizchi/markdown`
- `moonbit-community/cmark`
- `Yoorkin/cowsay` if present
- One module with missing `description`
- One module with missing `repository`
- One module with missing or unusual `license`
- The first 20 skills and any skill for `Yoorkin/cowsay` if present

You may choose additional samples when they reveal schema variation.

## Method

Use `uv run python <<'EOF' ... EOF` for structured analysis, or `curl` + `jq` when sufficient. Keep temporary files under:

`/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-01/`

Do not scan outside the workspace except for the explicit HTTP endpoints above.

## Required Report

Write the final report to:

`/home/t103o/workbench/projects/seekmoon/docs/Deep Research/2026-06-21-upstream-data-model-validation/reports/01-public-api-schema-validation.md`

The report must contain:

- Fetch time and toolchain used.
- Endpoint-by-endpoint schema tables.
- Field presence statistics for `/modules` and `/skills`.
- Manifest schema observations from representative samples.
- A `Correct / Stale / Incomplete / Wrong` assessment for fields mentioned in the existing design.
- A canonical data dictionary for SeekMoon public API ingestion.
- Evidence snippets limited to short JSON excerpts, not full dumps.

## Done Criteria

The report is complete only when a maintainer can update SeekMoon's current data dictionary without re-running your investigation.
