# MoonBit Native Subprocess Clap Probe Report

## Summary

This spike demonstrates a native MoonBit CLI that parses subcommands with `TheWaWaR/clap@0.2.6` and executes host commands with `moonbitlang/async/process@0.19.4`.

The CLI has three subcommands:

| Command | Purpose |
|---|---|
| `check <command>` | Checks whether a command exists in `PATH`. |
| `jq <filter> <json>` | Runs host `jq -c <filter>` with inline JSON as stdin. |
| `missing-jq <filter> <json>` | Runs the same path with a deliberately absent executable name. |

## CLI Parser

`config.mbt` builds a `TheWaWaR/clap` parser with `Parser`, `SubCommand`, `Arg::positional`, and `Nargs::Fixed`.

The parser surface used by this probe is:

```moonbit
@clap.Parser::new(
  prog="native-subprocess-clap-probe",
  subcmds={
    "check": @clap.SubCommand::new(args={
      "command": @clap.Arg::positional(nargs=Fixed(1)),
    }),
    "jq": @clap.SubCommand::new(args={
      "args": @clap.Arg::positional(nargs=Fixed(2)),
    }),
  },
)
```

`parse_config` maps the parsed `SimpleValue` into the probe's internal `Action` enum. Help output is represented as `Help(String)` so `--help` stays on the normal success path.

## Subprocess Execution

`run.mbt` uses the official async process package:

```moonbit
@process.collect_output_merged(
  command,
  ["-c", filter],
  stdin=@process.redirect_from_file(input_path),
)
```

The inline JSON payload is written to `_build/native-subprocess-clap-probe-input.json` before process execution and removed afterward.

## Verification

The current probe passes:

```bash
moon check --target native
moon test --target native
moon build --target native
moon run cmd/main --target native -- --help
moon run cmd/main --target native -- check jq
moon run cmd/main --target native -- jq '.name' '{"name":"moon","n":1}'
moon run cmd/main --target native -- missing-jq '.name' '{"name":"moon"}'
```

Observed runtime behavior:

```text
Usage: native-subprocess-clap-probe [OPTIONS] <COMMAND>
found: jq
"moon"
Cannot run '__seekmoon_definitely_missing_jq__': command was not found in PATH.
Install it or choose another executable.
```

The example has no custom C subprocess bridge.
