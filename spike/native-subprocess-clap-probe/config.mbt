///|
pub(all) enum Action {
  Check(String)
  Jq(String, String)
  MissingJq(String, String)
  Help(String)
}

///|
pub(all) struct Config {
  action : Action
}

///|
pub fn command() -> @clap.Parser {
  @clap.Parser::new(
    prog="native-subprocess-clap-probe",
    description="Probe host subprocess execution and clap-style CLI parsing",
    subcmds={
      "check": @clap.SubCommand::new(
        args={
          "command": @clap.Arg::positional(nargs=Fixed(1), help="Command name to find through PATH"),
        },
        help="Check whether a command exists in PATH",
      ),
      "jq": @clap.SubCommand::new(
        args={
          "args": @clap.Arg::positional(nargs=Fixed(2), help="jq filter followed by JSON text"),
        },
        help="Run host jq with a filter and inline JSON payload",
      ),
      "missing-jq": @clap.SubCommand::new(
        args={
          "args": @clap.Arg::positional(nargs=Fixed(2), help="jq filter followed by JSON text"),
        },
        help="Run the same passthrough path with a deliberately missing command",
      ),
    },
  )
}

///|
pub fn parse_config(argv : ArrayView[String]) -> Config raise {
  let parser = command()
  let value = @clap.SimpleValue::new(parser.prog)
  let help = parser.parse(value, argv)
  if help is Some(text) {
    { action: Help(text.to_string()) }
  } else {
    guard value.subcmd is Some(subcmd) else {
      fail(parser.gen_help_message([parser.prog], {}))
    }
    match subcmd.name {
      "check" => {
        guard subcmd.positional_args is [name] else {
          fail("missing command name")
        }
        { action: Check(name) }
      }
      "jq" => {
        guard subcmd.positional_args is [filter, json] else {
          fail("missing jq filter or json payload")
        }
        { action: Jq(filter, json) }
      }
      "missing-jq" => {
        guard subcmd.positional_args is [filter, json] else {
          fail("missing filter or json payload")
        }
        { action: MissingJq(filter, json) }
      }
      name => fail("unknown subcommand: \{name}")
    }
  }
}
