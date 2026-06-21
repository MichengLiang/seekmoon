# Deep Research Task 02: Python and Go Package Ecosystems

## Objective

Investigate Python and Go package discovery and package management mechanisms as evidence for a cross-ecosystem model of software package reuse.

## Required Output

Write a rigorous Chinese report to:

`projects/seekmoon/docs/Deep Research/reports/02-python-go-report.md`

The report must be at least 2000 Chinese characters.

## Scope

Cover these ecosystems:

- Python: PyPI, pip, uv, pyproject.toml, project metadata, wheel tags, Simple Repository API, trusted publishing where relevant.
- Go: go.mod, Go modules, module proxy, checksum database, Minimal Version Selection, pkg.go.dev.

## Research Questions

1. How does each ecosystem define package/module identity?
2. Which metadata is declared locally, and which metadata is inferred by registry/index infrastructure?
3. How does package discovery happen through PyPI, pkg.go.dev, search, docs, classifiers, module paths, or external indexes?
4. How does dependency resolution work, especially Python resolver behavior and Go Minimal Version Selection?
5. How do lock/reproducibility mechanisms differ?
6. How do Python wheels and Go module checksums affect artifact integrity and compatibility?
7. What security, trusted publishing, vulnerability, or provenance mechanisms exist?
8. Which variables should appear in a unified model?

## Evidence Priority

Use primary sources first:

- PyPA specifications and packaging.python.org.
- PyPI official docs where relevant.
- uv official docs.
- Go official module reference.
- pkg.go.dev official/about pages.

Use secondary sources only for ecosystem practice, and label them as secondary.

## Required Report Structure

1. Scope and object boundary.
2. Python lifecycle mapping.
3. Go lifecycle mapping.
4. Metadata and declared surfaces.
5. Discovery surfaces and data sources.
6. Dependency resolution and reproducibility.
7. Security, integrity, and provenance.
8. Variables and indicators useful for the unified model.
9. Cross-ecosystem comparison: Python vs Go.
10. Open uncertainties and exact references.

## Acceptance Criteria

- Explain why Python binary compatibility differs from Go module compatibility.
- Explain why Go has less traditional descriptive metadata and how pkg.go.dev compensates.
- Separate package-level evaluation variables from ecosystem-level mechanism variables.
- Cite exact URLs for claims.

## Prohibited Content

- Do not write code or tool design.
- Do not present memory-based claims without source.
- Do not use broad value judgments without variables.
