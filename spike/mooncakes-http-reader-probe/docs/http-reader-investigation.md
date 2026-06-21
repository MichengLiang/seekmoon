# Mooncakes HTTP Reader Investigation

调查日期：2026-06-22。

目标：判断 MoonBit 生态里是否存在适合作为 HTTP client 基础的包，并验证它能否支撑 `MooncakesApiReader`、`MooncakesAssetReader`、`SourceZipReader` 这类 typed reader 对象。

## 候选

| module | version | 用途 | 记录 |
|---|---:|---|---|
| `oboard/mio` | `0.5.2` | MoonBit native/js HTTP client | Apache-2.0；Mooncakes build success；README 声明 async request、GET/POST/PUT/DELETE/PATCH/OPTIONS/HEAD、binary/text/json response helper、native streaming HTTP/1.1、gzip/deflate/br/zstd decode、timeout/proxy/cert verify builder。 |
| `moonbitlang/async` | `0.19.4` | 官方 async 基础库 | 当前 `moon ide doc @moonbitlang/async/http` 只显示 `pub let unimplemented : Unit`，不构成可用 HTTP client。 |
| `mizchi/x` | `0.4.0` | compatibility layer | 包含 http/tls/socket 方向能力，但不是面向 Mooncakes reader 的清晰 high-level client。 |
| `mizchi/crater-browser-http` | `0.18.0` | crater/browser 场景 HTTP helper | 偏浏览器/应用沙箱，不是 native reader 首选。 |
| `bikallem/webapi` | `0.5.0` | browser fetch/WebAPI binding | 适合 browser target，不覆盖 native source zip reader。 |

结论：`oboard/mio` 是当前最接近 Python `httpx` / Rust `reqwest` 角色的 MoonBit 候选；还不是“事实标准”级别，因为 registry 下载量、生态引用和长期维护信号都有限，但它的 API 面已经足够支撑本 spike。

## 本地验证

本 spike 使用：

- `MooncakesApiReader`：`/api/v0/modules/statistics`、`/api/v0/manifest/oboard/mio`、`/api/v0/skills`
- `MooncakesAssetReader`：`/assets/oboard/mio@0.5.2/module_index.json`、`package_data.json`、`resource.json`
- `SourceZipReader`：`https://download.mooncakes.io/user/oboard/mio/0.5.2.zip`

每个 reader 返回 `HttpFetchResult`，分开记录 transport 是否成功、HTTP status、body bytes、JSON 解析结果。HTTP 200 不等于 Mooncakes 业务数据可用；`resource.json` 的 404 也只是“该资源文件不存在”，不能解释为 package 不存在。

验证结果：

| endpoint | status | body bytes | JSON |
|---|---:|---:|---|
| `/api/v0/modules/statistics` | 200 | 94 | ok |
| `/api/v0/manifest/oboard/mio` | 200 | 1267 | ok |
| `/api/v0/skills` | 200 | 41468 | ok |
| `/assets/oboard/mio@0.5.2/module_index.json` | 200 | 7300 | ok |
| `/assets/oboard/mio@0.5.2/package_data.json` | 200 | 42359 | ok |
| `/assets/oboard/mio@0.5.2/resource.json` | 200 | 5663 | ok |
| `https://download.mooncakes.io/user/oboard/mio/0.5.2.zip` | 302 then 200 | 97834 | skipped |

实现发现：

- `oboard/mio` 不自动跟随 source zip 的 302；`SourceZipReader` 需要读取 `Location` 头并显式跟随。
- `package_data.json` 与 `resource.json` 的 URL 使用 asset-relative package path。根包路径是空串，对应 `/assets/<owner>/<module>@<version>/package_data.json`，不能把 module index 里的 full package path `oboard/mio` 原样拼成 `/oboard/mio/package_data.json`。
- `ResponseBody::binary()`、`ResponseBody::text()`、`ResponseBody::json()` 能覆盖 body bytes、文本预览和 JSON parse outcome；`Response` 暴露 `code`、`reason`、`headers`、`cookies`。

本地验证命令：

```bash
moon check --target native
moon build --target native
moon run cmd/main --target native
```

三项均已通过。
