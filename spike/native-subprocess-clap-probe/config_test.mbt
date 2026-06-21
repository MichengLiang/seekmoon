///|
test "parse check subcommand" {
  let config = parse_config(["check", "jq"])
  guard config.action is Check(name) else {
    fail("expected check action")
  }
  assert_eq(name, "jq")
}

///|
test "parse jq subcommand" {
  let payload = "{\"name\":\"moon\"}"
  let config = parse_config(["jq", ".name", payload])
  guard config.action is Jq(filter, json) else {
    fail("expected jq action")
  }
  assert_eq(filter, ".name")
  assert_eq(json, payload)
}

///|
test "parse help as normal action" {
  let config = parse_config(["--help"])
  guard config.action is Help(text) else {
    fail("expected help action")
  }
  assert_true(text.contains("Usage: native-subprocess-clap-probe"))
}
