///|
pub(all) struct ProbeResult {
  code : Int
  output : String
}

///|
async fn command_exists(name : String) -> Bool {
  let (code, _) = @process.collect_output_merged("sh", ["-c", "command -v \"$1\" >/dev/null", "sh", name])
  code == 0
}

///|
async fn run_capture(command : String, filter : String, json : String) -> ProbeResult {
  let input_path = "_build/native-subprocess-clap-probe-input.json"
  @fs.write_file(input_path, json, create_mode=CreateOrTruncate)
  let input = @process.redirect_from_file(input_path)
  let (code, output) = @process.collect_output_merged(
    command,
    ["-c", filter],
    stdin=input,
  )
  @fs.remove(input_path) catch { _ => () }
  { code, output: output.text() }
}

///|
async fn run_jq(command : String, filter : String, json : String) -> ProbeResult {
  if !command_exists(command) {
    {
      code: 127,
      output: "Cannot run '\{command}': command was not found in PATH.\nInstall it or choose another executable.\n",
    }
  } else {
    run_capture(command, filter, json)
  }
}

///|
pub async fn run(config : Config) -> ProbeResult {
  match config.action {
    Check(name) =>
      if command_exists(name) {
        { code: 0, output: "found: \{name}\n" }
      } else {
        { code: 1, output: "missing: \{name}\n" }
      }
    Jq(filter, json) => {
      run_jq("jq", filter, json)
    }
    MissingJq(filter, json) =>
      run_jq("__seekmoon_definitely_missing_jq__", filter, json)
    Help(text) => { code: 0, output: text }
  }
}

///|
pub async fn run_cli(argv : ArrayView[String]) -> Int {
  let config = try {
    parse_config(argv)
  } catch {
    err => {
      println("error: \{err}")
      println(command().gen_help_message([command().prog], {}))
      return 2
    }
  }
  let result = run(config)
  println(result.output)
  result.code
}
