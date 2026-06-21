# Deep Research Task 01: JS/TS and Rust Package Ecosystems

## Objective

Investigate JS/TS and Rust package discovery and package management mechanisms as evidence for a cross-ecosystem model of software package reuse.

## Required Output

Write a rigorous Chinese report to:

`projects/seekmoon/docs/Deep Research/reports/01-js-ts-rust-report.md`

The report must be at least 2000 Chinese characters. More is acceptable when needed for precision.

## Scope

Cover these ecosystems:

- JS/TS: npm registry, npm CLI, package.json, pnpm, Yarn when relevant, JSR.
- Rust: Cargo, Cargo.toml, crates.io, docs.rs, lib.rs when relevant.

## Research Questions

1. What are the canonical package identity, metadata, dependency declaration, and artifact surfaces?
2. What registry or index infrastructure stores and exposes package metadata?
3. How does package discovery work through official and important third-party surfaces?
4. How does dependency resolution work?
5. What lockfile and reproducibility mechanisms exist?
6. What security, provenance, trusted publishing, audit, or permission mechanisms exist?
7. Which lifecycle stages are strong or weak in each ecosystem?
8. Which variables are measurable for a cross-ecosystem evaluation model?

## Evidence Priority

Use primary sources first:

- Official npm docs.
- Official JSR docs.
- Official pnpm / Yarn docs only when discussing their concrete mechanisms.
- Rust Cargo book/reference.
- crates.io and docs.rs official documentation where available.
- lib.rs public documentation or observable pages only when discussing discovery surface behavior.

Use secondary sources only to explain community practice, and label them as secondary.

## Required Report Structure

1. Scope and object boundary.
2. JS/TS lifecycle mapping.
3. Rust lifecycle mapping.
4. Discovery surfaces and data sources.
5. Management mechanisms: resolver, lockfile, cache/store, artifact integrity.
6. Security and provenance mechanisms.
7. Variables and indicators useful for the unified model.
8. Cross-ecosystem comparison: JS/TS vs Rust.
9. Open uncertainties and sources that require later checking.
10. References with exact URLs.

## Acceptance Criteria

- Distinguish Package Discovery from Package Management.
- Distinguish Manifest, Registry, Index, Resolver, Lockfile, and Artifact.
- Do not use unsupported claims such as “best”, “most advanced”, or “industry standard” without source or defined criterion.
- Every factual claim about a tool feature must have a cited source.
- The report must include concrete variables suitable for later scoring.

## Prohibited Content

- Do not write implementation plans or code.
- Do not recommend building a tool.
- Do not write motivational language.
- Do not include internal reasoning or agent process notes.
