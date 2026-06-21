# MoonBit CLI Package Snapshot

Investigation date: 2026-06-22.

Toolchain:

```text
moon 0.1.20260608 (60bc8c3 2026-06-08)
moonc v0.10.0+e66899a54 (2026-06-09)
moonrun 0.1.20260608 (60bc8c3 2026-06-08)
```

Registry snapshot:

```text
modules: 1352
packages: 12017
downloads: 4047641
```

## CLI Parser

This probe uses `TheWaWaR/clap@0.2.6`.

The package provides a Rust-like parser surface for MoonBit:

- `Parser` defines the program and top-level command table.
- `SubCommand` defines a command branch.
- `Arg` defines flag, named, and positional arguments.
- `Nargs` defines arity constraints such as `Fixed(1)` and `Fixed(2)`.
- `SimpleValue` stores parsed results for simple command shapes.
- `Value` supports custom typed parse-result sinks.
- `BasicValue` supports typed conversion for scalar values.

The published source and tests cover flags, named arguments, positionals, nested subcommands, global arguments, environment defaults, choices, typed conversion, custom result sinks, and generated help.

## Process Execution

This probe uses `moonbitlang/async/process@0.19.4`.

The package provides native process execution, file redirection, merged output capture, environment/cwd controls, and cancellation support. This probe uses captured output and file redirection to pass inline JSON into host `jq`.
