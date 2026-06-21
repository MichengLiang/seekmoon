# MoonBit Native Subprocess CLI Probe Report

## Summary

This spike verifies that a MoonBit program can be compiled to a native Linux executable and can delegate JSON filtering to the host `jq` command through the official subprocess API.

The current implementation uses `moonbitlang/async/process@0.19.4`. The earlier hand-written C FFI bridge has been removed because it was a low-level feasibility probe and would be misleading as an example for normal MoonBit code.

On this machine, `jq` is available at `/usr/bin/jq` and reports `jq-1.7`. The module passes native check, test, build, and runtime jq passthrough verification.

## Implementation

Created module:

`projects/seekmoon/spike/native-subprocess-cli-probe`

File responsibilities:

- `moon.mod`: module metadata, `moonbitlang/async@0.19.4` dependency, native preferred target.
- `moon.pkg`: imports `argparse`, `test`, `moonbitlang/async/fs`, and `moonbitlang/async/process`.
- `cmd/main/main.mbt`: thin async executable entrypoint.
- `cmd/main/moon.pkg`: imports the root package, `moonbitlang/core/env`, and `moonbitlang/async`.
- `config.mbt`: CLI shape and argument parsing.
- `run.mbt`: command discovery, jq passthrough, and friendly missing-command errors.
- `config_test.mbt`: parser tests.
- `.gitignore`: ignores `_build/` and `.mooncakes/`.
- `docs/implementation-plan.md`: implementation plan and task log.
- `README.md`: quick usage notes.

There is no `ffi.mbt` and no `native_process.c` in the current version.

## CLI Behavior

The CLI exposes three paths:

```bash
moon run cmd/main --target native -- check jq
moon run cmd/main --target native -- jq '.name' '{"name":"moon","n":1}'
moon run cmd/main --target native -- missing-jq '.name' '{"name":"moon"}'
```

`check` probes `PATH` by running:

```text
sh -c 'command -v "$1" >/dev/null' sh <command>
```

`jq` writes the inline JSON payload to `_build/native-subprocess-cli-probe-input.json`, opens that file as process stdin with `@process.redirect_from_file`, and runs:

```moonbit
@process.collect_output_merged(
  "jq",
  ["-c", filter],
  stdin=input,
)
```

`missing-jq` intentionally routes the same friendly error path through `__seekmoon_definitely_missing_jq__`, so the missing-command branch can be tested on machines where real `jq` exists.

## Why Official `moonbitlang/async/process`

The Mooncakes package search found multiple process-related candidates:

| Module | Version | Role |
|---|---:|---|
| `moonbitlang/async` | `0.19.4` | Official async process, pipe, redirection, cwd/env, and cancellation API. |
| `trkbt10/subprocess` | `0.2.0` | Node.js child_process-like wrapper built on `moonbitlang/async/process`. |
| `sennenki/process` | `0.1.0` | Rust/Go-style process command builder with native stubs. |
| `FrenchPicnic/which` | `0.1.3` | Command discovery only. |
| `moonbit-community/pty` | `0.2.1` | PTY spawning for interactive terminal processes. |
| `mizchi/x` | `0.4.0` | Cross-target compatibility layer that includes process support. |

For this example, the official package is the healthiest choice. It avoids custom C glue, avoids manual `Bytes`/C-string conversion, avoids hand-written `popen`, and passes arguments as an array.

## Verification

Final verification command:

```bash
moon clean
moon check --target native
moon test --target native
moon build --target native
moon run cmd/main --target native -- check jq
moon run cmd/main --target native -- jq '.name' '{"name":"moon","n":1}'
moon run cmd/main --target native -- jq '{name,n}' '{"name":"moon","n":1}'
moon run cmd/main --target native -- missing-jq '.name' '{"name":"moon"}'
```

Observed output:

```text
Finished. moon: ran 24 tasks, now up to date
Total tests: 2, passed: 2, failed: 0.
Finished. moon: ran 58 tasks, now up to date
found: jq

"moon"

{"name":"moon","n":1}

Cannot run '__seekmoon_definitely_missing_jq__': command was not found in PATH.
Install it or choose another executable.
```

This confirms:

- native check succeeds;
- native test succeeds;
- native build succeeds;
- host `jq` can be found;
- host `jq` can be executed through `moonbitlang/async/process`;
- jq output can be captured;
- a deliberately missing command produces a friendly message.

## Notes

The current CLI still prints captured output with `println(result.output)`, so output that already ends with newline displays an extra blank line. That is acceptable for this spike. A polished CLI should write to stdout exactly through an async stdio writer.

This is still a small jq-shaped probe, not a general shell DSL like `zx`, `plumbum`, or `xshell`. The important point is that normal MoonBit code should start from `moonbitlang/async/process`, not custom C FFI, when it needs subprocess capability.
