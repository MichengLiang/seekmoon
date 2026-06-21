# Task 02: Assets And Published Source Schema Validation

## Role

You are a reconnaissance agent validating Mooncakes asset endpoints and published source retrieval paths for SeekMoon.

You are not designing the CLI. You are validating upstream facts.

## Scope

Validate current shape, availability, and failure semantics for:

- `module_index.json`
- `package_data.json`
- `resources.json`
- Mooncakes source zip downloads

## Required Background Reading

Read these files before investigating:

- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 包发现闭环动向调查.md`
- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 语言、工具链与包生态调查报告.md`
- `/home/t103o/workbench/notes/关于MoonBit/关于包搜索与发现/MoonBit包搜索发现与生态查重SOP.md`
- `/home/t103o/workbench/projects/seekmoon/docs/关于moonbit/002-raw消费者侧 MoonBit 包发现 CLI 设计.md`

## Questions To Answer

1. What is the exact current schema of `module_index.json` for representative modules?
2. Is the tree key spelled `childs`, `children`, or something else?
3. How should package paths be derived from `module_index.json`?
4. What is the exact current schema of `package_data.json`?
5. Which keys are present in package data: `types`, `values`, `traits`, `errors`, `typealias`, `misc`, or others?
6. What fields exist inside type/value/method entries: `name`, `docstring`, `signature`, `loc`, `methods`, etc.?
7. When does `resources.json` exist? What does it contain when present? What does 404 mean?
8. Is source zip available for representative modules? Is it available for `moonbitlang/core`?
9. Are source zip contents sufficient as published-source evidence?
10. Which asset/source assumptions in the existing SOP/design are current, stale, incomplete, or wrong?

## Representative Samples

At minimum inspect:

- `moonbitlang/core` latest manifest version.
- `moonbitlang/core/argparse` package data.
- `mizchi/markdown` latest manifest version.
- At least one package under `mizchi/markdown`.
- `moonbit-community/cmark` latest manifest version.
- One skill module that has a Wasm executable entry.
- At least one module where `resources.json` returns 404.
- At least one module where `resources.json` returns 200, if you can find one without excessive search.

## Method

Use structured JSON parsing. Keep temporary files under:

`/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-02/`

Do not scan outside the workspace except for the explicit HTTP endpoints above.

## Required Report

Write the final report to:

`/home/t103o/workbench/projects/seekmoon/docs/Deep Research/2026-06-21-upstream-data-model-validation/reports/02-assets-and-source-schema-validation.md`

The report must contain:

- Fetch time.
- Asset URL construction rules.
- Schema tables for `module_index.json` and `package_data.json`.
- Availability/failure table for `resources.json`.
- Availability/failure table for source zip.
- Correct/stale/incomplete/wrong assessment of existing assumptions.
- Canonical data dictionary for SeekMoon asset ingestion.

## Done Criteria

The report is complete only when SeekMoon can implement `view`, `api`, and `source` ingestion without guessing asset field names.
