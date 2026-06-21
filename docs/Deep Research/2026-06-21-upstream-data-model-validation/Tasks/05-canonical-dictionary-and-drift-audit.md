# Task 05: Canonical Dictionary And Drift Audit

## Role

You are a reviewer/reconnaissance agent auditing the existing SeekMoon design documents against current upstream data facts.

You are not designing new features. You are identifying which previously proposed fields and data sources are correct, stale, incomplete, unsupported, or outside the current object boundary.

## Scope

Audit the data-source and field claims in:

- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 包发现闭环动向调查.md`
- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 语言、工具链与包生态调查报告.md`
- `/home/t103o/workbench/notes/关于MoonBit/关于包搜索与发现/MoonBit包搜索发现与生态查重SOP.md`
- `/home/t103o/workbench/projects/seekmoon/docs/关于moonbit/002-raw消费者侧 MoonBit 包发现 CLI 设计.md`

Use live spot checks where necessary, but your primary output is the audit map and canonical dictionary.

## Questions To Answer

1. What data sources are explicitly claimed across the existing docs and chat record?
2. Which claimed sources are current and confirmed?
3. Which claimed sources are stale, incomplete, ambiguous, or unsupported?
4. Which proposed fields belong to current upstream data?
5. Which proposed fields are derived fields SeekMoon may compute?
6. Which proposed fields are future/unsupported and must not appear in current output?
7. Which fields require status modeling (`present`, `missing`, `unknown`, `failed`, `unavailable`)?
8. What is the canonical current data dictionary for SeekMoon v0 ingestion?

## Classification Rules

Use these statuses:

- `confirmed-current`: observed current upstream field/source.
- `confirmed-local`: observed current local toolchain/cache field/source.
- `derived`: not upstream, but legitimately computed by SeekMoon from current sources.
- `missing-value`: current schema field exists, but value may be absent.
- `unknown-evidence`: current query asks for it, but no current evidence source can determine it.
- `unsupported-future`: not a current field; must not appear in current output.
- `stale`: previously true or claimed, but no longer current.
- `wrong`: contradicted by current evidence.

## Required Report

Write the final report to:

`/home/t103o/workbench/projects/seekmoon/docs/Deep Research/2026-06-21-upstream-data-model-validation/reports/05-canonical-dictionary-and-drift-audit.md`

The report must contain:

- Audit table of all claimed data sources.
- Audit table of all important claimed fields.
- Current canonical dictionary grouped by command surface: `sync`, `search`, `view`, `api`, `source`, `skill`, `compare`, `probe`, `record`, `report`, `raw`.
- Explicit list of unsupported future fields that must not appear in current output.
- Explicit list of fields that require `missing` vs `unknown` distinction.
- Notes on what must be revalidated later because it is temporally unstable.

## Done Criteria

The report is complete only when the principal coordinator can synthesize a final user-facing answer and update SeekMoon's data model without re-reading all prior chat history.
