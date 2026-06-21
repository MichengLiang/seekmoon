# Native Subprocess CLI Probe

This MoonBit spike checks whether a native Linux executable can delegate work to a host command through a subprocess boundary.

It uses the official `moonbitlang/async/process@0.19.4` package. There is no custom C subprocess bridge in this probe.

The CLI has three paths:

```bash
moon run cmd/main check jq
moon run cmd/main jq '.name' '{"name":"moon"}'
moon run cmd/main missing-jq '.name' '{"name":"moon"}'
moon build --target native
```

`check` probes `PATH` through `sh -c 'command -v "$1"'`. `jq` delegates to host `jq` via `@process.collect_output_merged`. `missing-jq` intentionally uses a command name that should not exist so the friendly error path stays testable on machines where `jq` is installed.

Inline JSON is written to an `_build/` scratch file and passed to `jq` with `@process.redirect_from_file`, because the official process API already gives us file redirection and captured output without hand-written FFI.
