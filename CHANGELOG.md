# Changelog

All notable changes to SeekMoon are recorded in this file.

## [0.1.0] - 2026-06-24

Initial public release of the SeekMoon MoonBit package discovery workbench.

### Added

- Go CLI entrypoint for `seekmoon` with discovery, inspection, validation, record, report, raw source, and skill workflows.
- Evidence model that separates upstream facts from local derived evidence across Mooncakes API data, assets, MoonBit toolchain state, repository signals, project context, snapshots, sessions, records, probes, and reports.
- Public commands for `doctor`, `sync`, `search`, `view`, `api`, `source`, `skill search`, `skill view`, `compare`, `probe`, `record`, `report`, and `raw`.
- Common output modes for pretty terminal output, JSON projections, jq evaluation, output shape inspection, JSON Schema inspection, and structured error rendering.
- Numbered candidate sessions so `search` and `skill search` results can be reused by later commands with inputs such as `1` or `2`.
- Maintained English CLI help with Chinese review material in the source tree.
- Acceptance, black-box, integration, journey, contract, service, output, source, store, model, platform, and CLI tests for the first public surface.
- AsciiDoc bookshelf source and GitHub Pages workflow for the project documentation site.
- Project governance and collaboration files, including contribution, security, support, citation, issue templates, pull request template, and dependabot configuration.
- GoReleaser configuration for Linux, macOS, and Windows artifacts on `amd64` and `arm64`.

### Release Artifacts

- `seekmoon` binary archives for Linux `amd64` and `arm64`.
- `seekmoon` binary archives for macOS `amd64` and `arm64`.
- `seekmoon` binary archives for Windows `amd64` and `arm64`.
- `checksums.txt` for release artifact verification.
