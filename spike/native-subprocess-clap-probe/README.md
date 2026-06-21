# Native Subprocess Clap Probe

This MoonBit spike combines a clap-style command-line parser with native host process execution.

It uses:

- `TheWaWaR/clap@0.2.6` for the CLI command tree, subcommands, positional arity, and help output.
- `moonbitlang/async/process@0.19.4` for native process execution and captured output.

The CLI exposes three subcommands:

```bash
moon run cmd/main --target native -- check jq
moon run cmd/main --target native -- jq '.name' '{"name":"moon"}'
moon run cmd/main --target native -- missing-jq '.name' '{"name":"moon"}'
```

`check` probes `PATH` through `sh -c 'command -v "$1"'`. `jq` delegates to the host `jq` executable and captures merged stdout/stderr. `missing-jq` uses a deliberately absent executable name so the missing-command path remains testable on machines where `jq` is installed.

Inline JSON is written to `_build/native-subprocess-clap-probe-input.json` and passed to `jq` with `@process.redirect_from_file`.
