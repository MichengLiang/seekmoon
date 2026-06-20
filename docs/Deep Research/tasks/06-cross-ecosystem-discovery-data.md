# Deep Research Task 06: Cross-Ecosystem Discovery Data, Graphs, and Index Surfaces

## Objective

Investigate cross-ecosystem discovery surfaces and data providers that aggregate package metadata, dependency graphs, advisories, licenses, and ecosystem signals.

## Required Output

Write a rigorous Chinese report to:

`projects/seekmoon/docs/Deep Research/reports/06-cross-ecosystem-discovery-data-report.md`

The report must be at least 2000 Chinese characters.

## Scope

Cover discovery and data surfaces:

- deps.dev / Open Source Insights.
- Libraries.io.
- ecosyste.ms if useful.
- lib.rs as Rust discovery augmentation.
- Swift Package Index.
- MvnRepository.
- FuGet.
- pkg.go.dev.
- GitHub dependency graph and advisory surfaces where relevant.

## Research Questions

1. Which surfaces are official, third-party, registry-backed, source-backed, or graph-backed?
2. What data does each surface expose?
3. Which surfaces support search, comparison, dependency graph, reverse dependencies, licenses, vulnerabilities, advisories, maintenance data, docs, or API inspection?
4. Which data can support a cross-ecosystem evaluation model?
5. What are known limitations or source-of-truth problems?
6. How should official registry data and third-party discovery data be weighted differently?

## Evidence Priority

Use primary docs for each platform when available.
Use observable public pages only when documentation is unavailable; clearly mark observation-based claims.

## Required Report Structure

1. Scope and object boundary.
2. Classification of discovery/data surfaces.
3. Data fields and evidence types by surface.
4. Dependency graph and reverse dependency support.
5. Security/license/advisory support.
6. Documentation/API inspection support.
7. Variables and indicators useful for the unified model.
8. Source-of-truth and reliability issues.
9. Open uncertainties and references.

## Acceptance Criteria

- Distinguish official source-of-truth data from third-party aggregation.
- Explain which discovery signals are comparable across ecosystems and which are not.
- Cite exact URLs for claims.

## Prohibited Content

- Do not rank websites by taste.
- Do not make unsourced claims about user popularity.
- Do not write implementation plans.
