# Task 04: Mooncakes Frontend Source Validation

## Role

You are a reconnaissance agent validating the official Mooncakes frontend source as an additional upstream evidence source.

You are not designing the CLI. You are verifying how the official site itself consumes API fields and asset structures.

## Scope

Inspect the current `moonbitlang/mooncakes.io` source. Use an external shallow clone or existing clone under:

`/home/t103o/workbench/external/mooncakes.io`

If the repo already exists, update it safely without deleting user changes. If cloning is needed, use `git clone --depth 1`.

## Required Background Reading

Read these files before investigating:

- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 包发现闭环动向调查.md`
- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 语言、工具链与包生态调查报告.md`
- `/home/t103o/workbench/notes/关于MoonBit/关于包搜索与发现/MoonBit包搜索发现与生态查重SOP.md`
- `/home/t103o/workbench/projects/seekmoon/docs/关于moonbit/002-raw消费者侧 MoonBit 包发现 CLI 设计.md`

## Questions To Answer

1. Which frontend files request `/api/v0/modules`, `/api/v0/modules/statistics`, `/api/v0/manifest`, `/api/v0/skills`, and asset JSON?
2. What fields does the frontend use from module objects?
3. What fields does the frontend use from manifest objects?
4. What fields does the frontend use from skill objects?
5. How does frontend search work? Is it client-side? Which fields enter the search index?
6. Does frontend code use `description`, `keywords`, `author`, recency, downloads, pinned modules, build status, or target metadata?
7. How are docs assets loaded or referenced?
8. Which source-code-derived assumptions should SeekMoon trust, and which still require live API validation?

## Method

Use `mcp__docutouch.search_text` or `rg` scoped to `/home/t103o/workbench/external/mooncakes.io`.

Do not edit the external repo.

## Required Report

Write the final report to:

`/home/t103o/workbench/projects/seekmoon/docs/Deep Research/2026-06-21-upstream-data-model-validation/reports/04-mooncakes-frontend-source-validation.md`

The report must contain:

- Repo commit hash and fetch/update time.
- File/line references for API calls and field usage.
- Frontend-consumed field dictionary.
- Search algorithm summary.
- Skill marketplace field usage summary.
- Correct/stale/incomplete/wrong assessment of existing design assumptions based on frontend source.
- Distinction between live API facts and frontend implementation facts.

## Done Criteria

The report is complete only when SeekMoon can align its ingestion model with official frontend field usage without mistaking frontend behavior for a stable public API contract.
