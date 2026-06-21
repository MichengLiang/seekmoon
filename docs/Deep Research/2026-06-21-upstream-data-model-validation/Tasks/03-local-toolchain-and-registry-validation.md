# Task 03: Local Toolchain And Registry Validation

## Role

You are a reconnaissance agent validating local MoonBit toolchain behavior and local registry/cache data models for SeekMoon.

You are not designing the CLI. You are validating local facts and command behavior.

## Scope

Validate:

- Current `moon`, `moonc`, `moonrun`, and `mooncake` versions and command surfaces.
- `moon update` behavior.
- Local registry index structure under `~/.moon/registry/index`.
- Local symbols cache structure under `~/.moon/registry/symbols`.
- `moon ide doc` behavior.
- `moon fetch` behavior.
- `moon runwasm` command help and coordinate rules.
- Relevant cache paths.

## Required Background Reading

Read these files before investigating:

- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 包发现闭环动向调查.md`
- `/home/t103o/workbench/notes/关于MoonBit/MoonBit 语言、工具链与包生态调查报告.md`
- `/home/t103o/workbench/notes/关于MoonBit/关于包搜索与发现/MoonBit包搜索发现与生态查重SOP.md`
- `/home/t103o/workbench/projects/seekmoon/docs/关于moonbit/002-raw消费者侧 MoonBit 包发现 CLI 设计.md`

## Questions To Answer

1. What is the exact current local MoonBit toolchain version?
2. Does `moon --help` include `search`, `view`, `audit`, `outdated`, `runwasm`, `fetch`, and `ide doc`?
3. What does `moon ide doc --help` say? What query forms work?
4. Does `moon ide doc` search registry packages generally, or only locally visible symbols/packages?
5. What is the current local registry index file count and JSON Lines schema?
6. What is the current local symbols cache file count and JSON Lines schema?
7. What does `moon fetch <module>@<version>` create, and where?
8. What are `moon runwasm` coordinate rules and cache paths?
9. Which local-toolchain assumptions in the existing SOP/design are current, stale, incomplete, or wrong?

## Required Experiments

Run non-destructive commands only. Use temp/probe directories under:

`/home/t103o/workbench/tmp/seekmoon-upstream-schema-2026-06-21/task-03/`

At minimum test:

- `~/.moon/bin/moon --version`
- `~/.moon/bin/moon --help`
- `~/.moon/bin/moon ide --help`
- `~/.moon/bin/moon ide doc --help`
- `~/.moon/bin/moon ide doc '@moonbitlang/core/argparse'`
- `~/.moon/bin/moon ide doc 'markdown'`
- `~/.moon/bin/moon fetch mizchi/markdown@0.6.2` inside a temp project
- `~/.moon/bin/moon runwasm --help`
- registry index JSON Lines samples
- symbols cache JSON Lines samples

Do not modify user projects.

## Required Report

Write the final report to:

`/home/t103o/workbench/projects/seekmoon/docs/Deep Research/2026-06-21-upstream-data-model-validation/reports/03-local-toolchain-and-registry-validation.md`

The report must contain:

- Command output summaries.
- Local data schema tables.
- Exact paths and cache locations.
- Behavior/failure semantics for `moon ide doc`, `moon fetch`, and `moon runwasm`.
- Correct/stale/incomplete/wrong assessment.
- Canonical data dictionary for SeekMoon local ingestion.

## Done Criteria

The report is complete only when SeekMoon can implement local `doctor`, `sync`, `api`, `source`, `probe`, and `skill` support without guessing local file shapes.
