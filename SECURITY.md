# Security Policy

SeekMoon is a local CLI and documentation project. It reads public package
metadata, local MoonBit toolchain state, local registry/cache state, repository
signals, and project records.

## Supported Versions

The default branch receives security fixes. Tagged release support is not
defined until the project publishes versioned release lines.

## Reporting a Vulnerability

Report security issues through GitHub's private vulnerability reporting for this
repository when available. If private reporting is unavailable, open a minimal
GitHub issue that does not include exploit details, secrets, tokens, private
repository data, or sensitive local paths.

Include:

- affected command or component
- expected behavior
- observed behavior
- local side effects, if any
- reproduction steps that avoid publishing secrets

## Security Scope

In scope:

- unsafe local file writes
- command execution behavior in probe or source flows
- credential or token exposure
- incorrect trust boundaries between upstream facts and local derived evidence
- dependency vulnerabilities reachable from the Go CLI

Out of scope:

- vulnerabilities in MoonBit packages discovered by SeekMoon
- vulnerabilities in third-party registries, repositories, or package sources
- package quality judgments based only on downloads, metadata, or local probe
  results

SeekMoon does not certify package safety. Security-related output is evidence for
consumer investigation, not a safety guarantee.
