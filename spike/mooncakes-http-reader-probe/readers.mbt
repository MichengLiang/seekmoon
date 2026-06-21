///|
pub(all) struct HttpFetchResult {
  label : String
  url : String
  fetched_at_ms : UInt64
  status_code : Int
  reason : String
  header_count : Int
  body_bytes : Int
  json_state : String
  preview : String
}

///|
fn preview_text(text : String) -> String {
  if text.length() > 240 {
    text[:240].to_owned() + "..."
  } else {
    text
  }
}

///|
async fn fetch_url(
  label : String,
  url : String,
  expect_json? : Bool = true,
) -> HttpFetchResult {
  let fetched_at_ms = @env.now()
  let response = @mio.RequestClient::builder()
    .default_header("User-Agent", "seekmoon-mooncakes-http-reader-probe/0.1")
    .timeout(15000)
    .build()
    .get(url)
    .send()
  let bytes = response.binary()
  let text_preview = response.text() catch { _ => "<non-text response>" }
  let json_state = if expect_json {
    try {
      let json = response.json()
      let summary = Json::stringify(json, indent=0)
      "json-ok preview=\{preview_text(summary)}"
    } catch {
      err => "json-error \{err}"
    }
  } else {
    "json-skipped"
  }
  {
    label,
    url,
    fetched_at_ms,
    status_code: response.response.code,
    reason: response.response.reason,
    header_count: response.response.headers.length(),
    body_bytes: bytes.length(),
    json_state,
    preview: preview_text(text_preview),
  }
}

///|
fn redirect_location(headers : Map[String, String]) -> String? {
  match headers.get("Location") {
    Some(location) => Some(location)
    None => headers.get("location")
  }
}

///|
fn asset_file_url(
  base_url : String,
  owner : String,
  module_name : String,
  version : String,
  package_relpath : String,
  filename : String,
) -> String {
  let prefix = "\{base_url}/assets/\{owner}/\{module_name}@\{version}"
  if package_relpath.is_empty() {
    "\{prefix}/\{filename}"
  } else {
    "\{prefix}/\{package_relpath}/\{filename}"
  }
}

///|
async fn fetch_url_follow_redirect(
  label : String,
  url : String,
  expect_json? : Bool = true,
  redirects_left? : Int = 3,
) -> HttpFetchResult {
  let fetched_at_ms = @env.now()
  let response = @mio.RequestClient::builder()
    .default_header("User-Agent", "seekmoon-mooncakes-http-reader-probe/0.1")
    .timeout(15000)
    .build()
    .get(url)
    .send()
  if response.response.code >= 300 && response.response.code < 400 && redirects_left > 0 {
    match redirect_location(response.response.headers) {
      Some(location) => return fetch_url_follow_redirect(
        "\{label}:redirect",
        location,
        expect_json~,
        redirects_left=redirects_left - 1,
      )
      None => ()
    }
  }
  let bytes = response.binary()
  let text_preview = response.text() catch { _ => "<non-text response>" }
  let json_state = if expect_json {
    try {
      let json = response.json()
      let summary = Json::stringify(json, indent=0)
      "json-ok preview=\{preview_text(summary)}"
    } catch {
      err => "json-error \{err}"
    }
  } else {
    "json-skipped"
  }
  {
    label,
    url,
    fetched_at_ms,
    status_code: response.response.code,
    reason: response.response.reason,
    header_count: response.response.headers.length(),
    body_bytes: bytes.length(),
    json_state,
    preview: preview_text(text_preview),
  }
}

///|
pub(all) struct MooncakesApiReader {
  base_url : String
}

///|
pub fn MooncakesApiReader::new(base_url? : String = "https://mooncakes.io") -> MooncakesApiReader {
  { base_url, }
}

///|
pub async fn MooncakesApiReader::modules_statistics(self : MooncakesApiReader) -> HttpFetchResult {
  fetch_url("api:modules-statistics", "\{self.base_url}/api/v0/modules/statistics")
}

///|
pub async fn MooncakesApiReader::module_manifest(
  self : MooncakesApiReader,
  owner : String,
  module_name : String,
) -> HttpFetchResult {
  fetch_url("api:manifest", "\{self.base_url}/api/v0/manifest/\{owner}/\{module_name}")
}

///|
pub async fn MooncakesApiReader::skills(self : MooncakesApiReader) -> HttpFetchResult {
  fetch_url("api:skills", "\{self.base_url}/api/v0/skills")
}

///|
pub(all) struct MooncakesAssetReader {
  base_url : String
}

///|
pub fn MooncakesAssetReader::new(base_url? : String = "https://mooncakes.io") -> MooncakesAssetReader {
  { base_url, }
}

///|
pub async fn MooncakesAssetReader::module_index(
  self : MooncakesAssetReader,
  owner : String,
  module_name : String,
  version : String,
) -> HttpFetchResult {
  fetch_url(
    "asset:module-index",
    "\{self.base_url}/assets/\{owner}/\{module_name}@\{version}/module_index.json",
  )
}

///|
pub async fn MooncakesAssetReader::package_data(
  self : MooncakesAssetReader,
  owner : String,
  module_name : String,
  version : String,
  package_path : String,
) -> HttpFetchResult {
  fetch_url(
    "asset:package-data",
    asset_file_url(self.base_url, owner, module_name, version, package_path, "package_data.json"),
  )
}

///|
pub async fn MooncakesAssetReader::resource(
  self : MooncakesAssetReader,
  owner : String,
  module_name : String,
  version : String,
  package_path : String,
) -> HttpFetchResult {
  // A 404 is still useful reader data here: resource.json is optional per package.
  fetch_url(
    "asset:resource",
    asset_file_url(self.base_url, owner, module_name, version, package_path, "resource.json"),
  )
}

///|
pub(all) struct SourceZipReader {
  base_url : String
}

///|
pub fn SourceZipReader::new(base_url? : String = "https://download.mooncakes.io") -> SourceZipReader {
  { base_url, }
}

///|
pub async fn SourceZipReader::source_zip(
  self : SourceZipReader,
  owner : String,
  module_name : String,
  version : String,
) -> HttpFetchResult {
  fetch_url_follow_redirect(
    "source:zip",
    "\{self.base_url}/user/\{owner}/\{module_name}/\{version}.zip",
    expect_json=false,
  )
}

///|
pub async fn run_default_probe() -> Array[HttpFetchResult] {
  let api = MooncakesApiReader::new()
  let assets = MooncakesAssetReader::new()
  let source = SourceZipReader::new()
  [
    api.modules_statistics(),
    api.module_manifest("oboard", "mio"),
    api.skills(),
    assets.module_index("oboard", "mio", "0.5.2"),
    assets.package_data("oboard", "mio", "0.5.2", ""),
    assets.resource("oboard", "mio", "0.5.2", ""),
    source.source_zip("oboard", "mio", "0.5.2"),
  ]
}

///|
pub fn format_result(result : HttpFetchResult) -> String {
  "[\{result.label}]\n" +
  "url: \{result.url}\n" +
  "fetched_at_ms: \{result.fetched_at_ms}\n" +
  "status: \{result.status_code} \{result.reason}\n" +
  "headers: \{result.header_count}\n" +
  "body_bytes: \{result.body_bytes}\n" +
  "json: \{result.json_state}\n" +
  "preview: \{result.preview}\n"
}

///|
pub async fn run_cli(_argv : ArrayView[String]) -> Int {
  let results = run_default_probe()
  for result in results {
    println(format_result(result))
  }
  0
}
