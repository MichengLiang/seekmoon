# Deep Research Task 03: JVM/Kotlin and .NET Package Ecosystems

## Objective

Investigate Java/JVM/Kotlin and .NET package discovery and package management mechanisms as evidence for a cross-ecosystem model of software package reuse.

## Required Output

Write a rigorous Chinese report to:

`projects/seekmoon/docs/Deep Research/reports/03-jvm-dotnet-report.md`

The report must be at least 2000 Chinese characters.

## Scope

Cover these ecosystems:

- Java/JVM/Kotlin: Maven POM, Maven Central, Gradle dependency management, Gradle Module Metadata, Gradle variant model, Kotlin Multiplatform publishing.
- .NET/C#: NuGet, .nuspec, PackageReference, packages.lock.json, target frameworks, NuGet audit.

## Research Questions

1. How do Maven coordinates and NuGet package IDs define package identity?
2. What metadata is declared in POM, Gradle metadata, and .nuspec?
3. How do Maven Central, MvnRepository, Gradle Plugin Portal, NuGet.org, and IDE package managers support discovery?
4. How do dependency resolution and conflict mediation work in Maven, Gradle, and NuGet?
5. How do variants, target frameworks, and Kotlin Multiplatform change the meaning of “package compatibility”?
6. What lockfile/reproducibility mechanisms exist?
7. What vulnerability/audit/security metadata is exposed?
8. Which variables should appear in a unified model?

## Evidence Priority

Use primary sources first:

- Maven official POM and dependency mechanism docs.
- Maven Central / Sonatype official docs where needed.
- Gradle official docs.
- Kotlin official docs.
- Microsoft NuGet official docs.

Use MvnRepository or FuGet only as discovery surface examples; label them as third-party surfaces.

## Required Report Structure

1. Scope and object boundary.
2. JVM/Kotlin lifecycle mapping.
3. .NET lifecycle mapping.
4. Metadata and identity surfaces.
5. Discovery surfaces.
6. Dependency resolution, variants, and compatibility.
7. Reproducibility, lock files, and audits.
8. Variables and indicators useful for the unified model.
9. Cross-ecosystem comparison: JVM/Kotlin vs .NET.
10. Open uncertainties and references.

## Acceptance Criteria

- Explain groupId/artifactId/version and how they differ from flat package names.
- Explain Gradle variants and why Kotlin Multiplatform requires variant-aware metadata.
- Explain NuGet target framework compatibility as a package evaluation variable.
- Cite exact URLs for claims.

## Prohibited Content

- Do not treat Maven, Gradle, and NuGet as interchangeable installers.
- Do not write implementation plans.
- Do not use unsourced community lore as fact.
