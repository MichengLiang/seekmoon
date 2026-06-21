# Native Subprocess CLI Probe Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use inline execution for this spike. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a MoonBit native CLI experiment that detects host commands and delegates JSON filtering to the system `jq` executable through a child process.

**Architecture:** The root package owns argument parsing, user-facing errors, and subprocess result shaping. The official `moonbitlang/async/process` package performs command execution and output capture; `cmd/main` only reads argv and prints the result.

**Tech Stack:** MoonBit native target, `moonbitlang/core/argparse`, `moonbitlang/async/process`, `moonbitlang/async/fs`, jq.

---

### Task 1: Establish The Probe Module

**Files:**
- Create: `projects/seekmoon/spike/native-subprocess-cli-probe/moon.mod`
- Create: `projects/seekmoon/spike/native-subprocess-cli-probe/moon.pkg`
- Create: `projects/seekmoon/spike/native-subprocess-cli-probe/cmd/main/moon.pkg`
- Create: `projects/seekmoon/spike/native-subprocess-cli-probe/cmd/main/main.mbt`

- [x] **Step 1: Declare a native-preferred module**

Use `preferred_target = "native"` so normal `moon run`, `moon build`, and `moon test` exercise the Linux executable path by default.

- [x] **Step 2: Keep the main package thin**

Import the root package as `@app`, import `moonbitlang/core/env`, and call a single library function from `main`.

### Task 2: Implement Command Discovery And jq Delegation

**Files:**
- Create: `projects/seekmoon/spike/native-subprocess-cli-probe/config.mbt`
- Create: `projects/seekmoon/spike/native-subprocess-cli-probe/run.mbt`

- [x] **Step 1: Use `argparse` for a friendly CLI**

Expose `check <command>` for command discovery and `jq <filter> <json>` for host jq passthrough.

- [x] **Step 2: Use the official async process API**

Use `@process.collect_output_merged` for command execution and output capture. Use `@process.redirect_from_file` for jq stdin so the probe does not need hand-written C glue.

- [x] **Step 3: Return structured MoonBit results**

Turn process exit code and output into `ProbeResult` and friendly error text.

### Task 3: Verify The Spike

**Files:**
- Create: `projects/seekmoon/spike/native-subprocess-cli-probe/README.md`
- Create: `projects/seekmoon/spike/native-subprocess-cli-probe/report.md`

- [x] **Step 1: Check and build native**

Run `moon check --target native` and `moon build --target native`.

- [x] **Step 2: Exercise present and missing command paths**

Run the CLI against `jq`, a deliberately missing command, a jq filter, and a deliberately missing passthrough command.

- [x] **Step 3: Record adoption evidence**

Summarize the SOP-driven package search, local docs evidence, and exact verification commands in the report.
