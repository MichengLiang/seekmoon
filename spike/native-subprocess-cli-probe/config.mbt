///|
pub(all) enum Action {
  Check(String)
  Jq(String, String)
  MissingJq(String, String)
}

///|
pub(all) struct Config {
  action : Action
}

///|
pub fn command() -> @argparse.Command {
  @argparse.Command::new(
    "native-subprocess-cli-probe",
    about="Probe host command discovery and jq passthrough from a MoonBit native binary",
    subcommands=[
      @argparse.Command::new(
        "check",
        about="Check whether a command exists in PATH",
        positionals=[
          @argparse.PositionArg::new(
            "command",
            about="Command name to find through PATH",
            num_args=@argparse.ValueRange::single(),
          ),
        ],
      ),
      @argparse.Command::new(
        "jq",
        about="Run host jq with a filter and inline JSON payload",
        positionals=[
          @argparse.PositionArg::new(
            "filter",
            about="jq filter expression",
            num_args=@argparse.ValueRange::single(),
          ),
          @argparse.PositionArg::new(
            "json",
            about="JSON text passed to jq",
            num_args=@argparse.ValueRange::single(),
          ),
        ],
      ),
      @argparse.Command::new(
        "missing-jq",
        about="Run the same passthrough path with a deliberately missing command",
        positionals=[
          @argparse.PositionArg::new(
            "filter",
            about="filter expression",
            num_args=@argparse.ValueRange::single(),
          ),
          @argparse.PositionArg::new(
            "json",
            about="JSON text",
            num_args=@argparse.ValueRange::single(),
          ),
        ],
      ),
    ],
    subcommand_required=true,
    arg_required_else_help=true,
  )
}

///|
pub fn parse_config(argv : ArrayView[String]) -> Config raise {
  let matches = @argparse.parse(command(), argv~)
  match matches.subcommand {
    Some(("check", sub)) => {
      guard sub.values is { "command": [name], .. } else {
        fail("missing command name")
      }
      { action: Check(name) }
    }
    Some(("jq", sub)) => {
      guard sub.values is { "filter": [filter], "json": [json], .. } else {
        fail("missing jq filter or json payload")
      }
      { action: Jq(filter, json) }
    }
    Some(("missing-jq", sub)) => {
      guard sub.values is { "filter": [filter], "json": [json], .. } else {
        fail("missing filter or json payload")
      }
      { action: MissingJq(filter, json) }
    }
    _ => fail("missing subcommand")
  }
}
