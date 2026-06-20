# Deep Research Task 04: C/C++ and Swift Package Ecosystems

## Objective

Investigate C/C++ and Swift package discovery and package management mechanisms as evidence for a cross-ecosystem model of software package reuse.

## Required Output

Write a rigorous Chinese report to:

`projects/seekmoon/docs/Deep Research/reports/04-cpp-swift-report.md`

The report must be at least 2000 Chinese characters.

## Scope

Cover these ecosystems:

- C/C++: Conan 2, vcpkg manifest mode, ConanCenter, vcpkg registries, lockfiles, profiles, settings/options, triplets.
- Swift: Swift Package Manager, Package.swift, PackageDescription, Swift Package Registry Service, Swift Package Index.

## Research Questions

1. How does C/C++ package identity differ from source-only language package identity?
2. Which variables capture compiler, ABI, platform, build type, link mode, and options?
3. How do Conan and vcpkg declare dependencies and variants?
4. How do Conan/vcpkg discovery surfaces work?
5. How does SwiftPM declare packages, products, targets, dependencies, and platform constraints?
6. What role does Swift Package Index play in discovery and compatibility assessment?
7. What registry, lock, reproducibility, and security mechanisms exist?
8. Which variables should appear in a unified model?

## Evidence Priority

Use primary sources first:

- Conan 2 official docs.
- Microsoft vcpkg official docs.
- Swift official docs and Swift evolution proposal for package registry service.
- Swift Package Index official pages for discovery surface behavior.

Use secondary sources only when clearly labeled.

## Required Report Structure

1. Scope and object boundary.
2. C/C++ lifecycle mapping.
3. Swift lifecycle mapping.
4. Metadata, variants, and compatibility surfaces.
5. Discovery surfaces and data sources.
6. Dependency resolution, lockfiles, and artifact realization.
7. Variables and indicators useful for the unified model.
8. Cross-ecosystem comparison: C/C++ vs Swift.
9. Open uncertainties and references.

## Acceptance Criteria

- Explain why ABI/platform/compiler variables are first-class in C/C++ package evaluation.
- Explain Swift products/targets/platforms in Package.swift.
- Separate package manager mechanisms from third-party discovery indexes.
- Cite exact URLs for claims.

## Prohibited Content

- Do not generalize C/C++ behavior from npm/PyPI-style registries.
- Do not write code or build instructions.
- Do not make unsourced claims about ecosystem popularity.
