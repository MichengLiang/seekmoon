# SeekMoon 上游数据模型与输出边界

日期：2026-06-21  
对象：SeekMoon v0 对 Mooncakes、MoonBit 本地工具链、MoonBit 项目上下文和 SeekMoon 自有记录的读取、归一化、派生与输出边界。

本文档是独立规范。实现 SeekMoon v0 数据模型、命令输出、缓存记录和审计记录时，不依赖同目录下的分项文件。

## 1. 总边界

SeekMoon v0 只接收六类证据。

第一类是 Mooncakes 公共 HTTP 数据源：

- `https://mooncakes.io/api/v0/modules`
- `https://mooncakes.io/api/v0/modules/statistics`
- `https://mooncakes.io/api/v0/manifest/<owner>/<module>`
- `https://mooncakes.io/api/v0/skills`
- `https://mooncakes.io/api/v0/skills/<entry>`

第二类是 Mooncakes 静态资产：

- `https://mooncakes.io/assets/<module>@<version>/module_index.json`
- `https://mooncakes.io/assets/<module>@<version>/<package-relpath>/package_data.json`
- `https://mooncakes.io/assets/<module>@<version>/<package-relpath>/resource.json`
- `https://mooncakes.io/assets/<skill-entry>/SKILL.md`
- Skills API 返回的 `wasm_url`
- Skills API 返回的 `checksum_url`
- `https://download.mooncakes.io/user/<owner>/<module>/<version>.zip`

第三类是本地 MoonBit 工具链和缓存：

- `moon`
- `moonc`
- `moonrun`
- `mooncake`
- `~/.moon/registry/index`
- `~/.moon/registry/symbols`
- `~/.moon/registry/cache`
- `~/.moon/registry/cache/assets`
- `~/.moon/lib/core`

第四类是当前 MoonBit 项目上下文：

- `moon.mod`
- `moon.mod.json`
- `moon.pkg`
- `moon.pkg.json`
- 当前项目依赖、目标后端、构建与测试命令输出
- probe 临时项目里的命令输出和日志

第五类是 GitHub 维护信号。只有在 Mooncakes 数据给出 repository，且 SeekMoon 实际请求 GitHub 后，GitHub 字段才进入输出。GitHub 字段是外部维护信号，不是 Mooncakes 发布数据。

第六类是 SeekMoon 派生数据：

- 搜索 rank
- match fields
- 规范化字段
- 版本数量
- docs URL
- runwasm coordinate
- source resolution
- probe 汇总状态
- record/report 记录

SeekMoon v0 不接收没有当前来源的生态愿望字段。质量分、安全审计、漏洞状态、反向依赖、provenance、attestation、verified publisher、SBOM、server-side search、Web 搜索深链都不属于 v0 数据模型。

## 2. 状态词

所有可空、可失败、可派生的字段必须使用稳定状态词。

| 状态 | 含义 |
|---|---|
| `present` | 当前来源成功给出有效值。 |
| `missing` | 当前来源或 schema 有字段位置，但值为空字符串、空数组、null 或字段缺失。 |
| `unknown` | 当前问题需要判断，但已读取的来源不能回答。 |
| `failed` | HTTP 请求、命令或解析动作已经执行并失败。 |
| `unavailable` | 该对象没有这个可选来源，或可选资产不存在。 |
| `derived` | SeekMoon 从当前事实计算得到。 |
| `unsupported` | 当前没有证据来源；字段不进入默认输出和 v0 JSON。 |

`missing` 与 `unknown` 不能混用。空 description 是 `missing`。没有 target metadata 且没有 probe 时，target support 是 `unknown`。`resource.json` 404 是 `unavailable`。source zip 404 是 `unavailable`。symbols cache 未命中是 `unknown`，因为本地 symbols cache 不是全量索引。

默认 pretty text 只显示当前命令需要的字段。JSON 输出保留状态和值。raw 输出保留原始 payload。unsupported 字段不显示为 `unknown`，因为它们不是当前 schema 的问题。

## 3. 公共 API 模型

### 3.1 Modules API

来源：

```text
GET https://mooncakes.io/api/v0/modules
```

2026-06-21 观测到顶层 JSON 数组，长度为 1350。数组元素是 module summary object。每个元素有 8 个顶层键：

- `name`
- `version`
- `description`
- `keywords`
- `repository`
- `license`
- `is_new`
- `created_at`

这 8 个键在 1350 个对象中都存在。缺失内容用空字符串或空数组表达。

| 字段 | 当前值形态 |
|---|---|
| `name` | 1350 个非空字符串。 |
| `version` | 1350 个非空字符串。 |
| `description` | 1023 个非空字符串，327 个空字符串。 |
| `keywords` | 933 个非空数组，417 个空数组。 |
| `repository` | 1058 个非空字符串，292 个空字符串。 |
| `license` | 1336 个非空字符串，14 个空字符串。 |
| `is_new` | 1350 个布尔值。 |
| `created_at` | 1350 个非空字符串。 |

SeekMoon canonical `ModuleSummary` 使用以下映射：

| canonical 字段 | 上游字段 | 规则 |
|---|---|---|
| `module` | `name` | 必填；完整 owner/module 坐标。 |
| `version` | `version` | 必填；列表中的当前版本。 |
| `description` | `description` | 空字符串为 `missing`。 |
| `keywords` | `keywords` | 空数组为 `missing`。 |
| `repository` | `repository` | 空字符串为 `missing`。 |
| `license` | `license` | 空字符串为 `missing`。 |
| `is_new` | `is_new` | 列表层字段，不自动派生长期 freshness。 |
| `created_at` | `created_at` | 时间字符串；排序或展示时再解析。 |

Modules API 不提供 `downloads`、`build_status`、`has_package`、`versions`、`targets`、`supported-targets`、`supported_targets`、quality score、audit、reverse dependencies、provenance 或 SBOM。

### 3.2 Modules Search Query

来源形态：

```text
GET https://mooncakes.io/api/v0/modules?search=markdown
GET https://mooncakes.io/api/v0/modules?search=cowsay
```

2026-06-21 观测到这两个请求都返回完整 1350 项数组，解析后与 `/api/v0/modules` 相同，第一项仍是 `0Ayachi0/elk`。

`search` query 参数不是可依赖的服务端搜索。SeekMoon `search` 必须拉取完整 `/api/v0/modules`，再执行本地过滤与排序。若 SeekMoon 搜索 description，那是 SeekMoon 搜索语义；Mooncakes 官方前端首页搜索不把 description 放进模块搜索 index。

### 3.3 Statistics API

来源：

```text
GET https://mooncakes.io/api/v0/modules/statistics
```

顶层 JSON object 有 4 个整数键：

- `total_modules`
- `total_packages`
- `total_lines`
- `total_downloads`

2026-06-21 多个时间点的观测值出现小幅变化，例如 `total_modules=1350`、`total_packages=12008`，`total_lines` 与 `total_downloads` 随 registry 和下载计数变化。statistics 只能作为 snapshot 字段，不能作为永久常量。

### 3.4 Manifest API

来源：

```text
GET https://mooncakes.io/api/v0/manifest/<owner>/<module>
```

代表样本的顶层 object 有 9 个键：

- `name`
- `module`
- `version`
- `latest_version`
- `downloads`
- `has_package`
- `build_status`
- `metadata`
- `versions`

canonical `Manifest`：

| canonical 字段 | 上游字段 | 规则 |
|---|---|---|
| `module` | `module` / `name` | 必填；样本中二者相等。 |
| `version` | `version` | 必填；当前 manifest 版本。 |
| `latest_version` | `latest_version` | 必填；与 `version` 分开保存。 |
| `downloads` | `downloads` | 整数；detail 层字段。 |
| `has_package` | `has_package` | 布尔值；表示 Mooncakes docs/package asset 相关状态，不等于本地可采用。 |
| `build_status` | `build_status` | 字符串或 null；样本中有 `success`、`legacy`、null。 |
| `metadata` | `metadata` | 开放 object；必须保留 raw。 |
| `versions` | `versions` | 数组；样本元素包含 `version` 和 `yanked`。 |
| `versions_count` | derived | `versions.length`；不是上游字段。 |

`build_status` 是 Mooncakes 构建/文档信号。它不证明任意 target 构建成功，不证明包可采用，也不替代本地 probe。

`metadata` 是开放对象。代表样本中出现过：

- `checksum`
- `name`
- `version`
- `created_at`
- `license`
- `repository`
- `description`
- `keywords`
- `readme`
- `source`
- `deps`
- `authors`
- `dependencies`
- `preferred-target`
- `preferred_target`
- `supported-targets`
- `warn-list`

metadata 字段不能视为闭合 schema。SeekMoon 只规范化已知字段，同时保留 raw metadata。`dependencies` 在样本中出现为空字符串，不能无条件合并进 `deps`。`deps` 的有效形态是 object map。

target metadata 不稳定。2026-06-21 的观测事实：

- `/api/v0/modules` 没有 top-level target 字段。
- `/api/v0/skills` 没有 top-level target 字段。
- 代表 manifest 中，`mizchi/markdown` 有 `metadata["preferred-target"] = "js"`。
- 代表 manifest 中，`Yoorkin/cowsay` 有 `metadata.preferred_target = "wasm"` 和 `metadata["supported-targets"] = "wasm"`。
- 未观察到稳定的 `metadata.targets`、top-level `targets`、top-level `supported-targets` 或 top-level `supported_targets`。

SeekMoon 可以把 `preferred-target` 与 `preferred_target` 规范化为 `preferred_target`。SeekMoon 可以读取 `supported-targets`、`supported_targets`、`targets` 并规范化为 `supported_targets`。规范化结果必须带 source spelling。没有 metadata 或 probe 证据时，target support 是 `unknown`，不能输出 supported/unsupported 布尔值。

### 3.5 Skills API

来源：

```text
GET https://mooncakes.io/api/v0/skills
```

2026-06-21 观测到顶层数组，长度为 70。每个 skill object 有 12 个顶层键：

- `module`
- `author`
- `author_avatar`
- `version`
- `package`
- `name`
- `detail_url`
- `wasm_url`
- `checksum_url`
- `metadata`
- `repository`
- `created_at`

字段规则：

| 字段 | 当前值形态 |
|---|---|
| `module` | 70 个非空字符串。 |
| `author` | 70 个非空字符串。 |
| `author_avatar` | 非空或空字符串；空字符串为 `missing`。 |
| `version` | 70 个非空字符串。 |
| `package` | 非空或空字符串；空字符串表示 root/default executable package，不可删除。 |
| `name` | 70 个非空字符串。 |
| `detail_url` | 70 个非空相对 URL。 |
| `wasm_url` | 70 个非空相对 URL。 |
| `checksum_url` | 70 个非空相对 URL。 |
| `metadata` | 70 个 object；内部稀疏。 |
| `repository` | 非空或空字符串。 |
| `created_at` | 70 个非空字符串。 |

`metadata.description` 在 70 条中都存在，但多数为空字符串。`metadata.name` 可能非空、空字符串或缺失。Skill record 属于 executable / Wasm / `moon runwasm` 对象，不属于 library dependency candidate。Search、view、record 可以支持 skill，但 skill 不得与 library module 共享采纳结论。

## 4. Mooncakes Assets 模型

### 4.1 Asset URL 构造

先读取 manifest，使用 manifest `version` 构造 asset 路径。

```text
module_asset_base = https://mooncakes.io/assets/<module>@<version>/
module_index_url = <module_asset_base>/module_index.json
package_data_url = <module_asset_base>/<package-relpath>/package_data.json
resource_url = <module_asset_base>/<package-relpath>/resource.json
source_zip_url = https://download.mooncakes.io/user/<module>/<version>.zip
```

`package-relpath` 不能靠文件系统猜测。必须先读 `module_index.json`，从 package node 的 `package.path` 派生。

派生规则：

| 值 | 示例 |
|---|---|
| module path | `moonbitlang/core` |
| package path | `moonbitlang/core/argparse` |
| package relpath | `argparse` |
| package data URL | `https://mooncakes.io/assets/moonbitlang/core@0.1.20260609+84519ca0a/argparse/package_data.json` |

root package 的 package path 等于 module path，package relpath 为空。`mizchi/markdown` root package 的 package data URL 是：

```text
https://mooncakes.io/assets/mizchi/markdown@0.6.2/package_data.json
```

### 4.2 `module_index.json`

`module_index.json` 是 module package tree 与紧凑 API index。顶层和子节点使用同一 node 结构：

- `name`
- `package`
- `childs`

当前实际拼写是 `childs`，不是 `children`。实现必须读取 `childs`。可兼容 `children`，但不能把 `children` 当作当前事实。

`package` 为 object 或 null。`package != null` 的 node 是 package summary node。package object 字段：

- `path`
- `traits`
- `errors`
- `types`
- `typealias`
- `values`
- `misc`

module index 中的 type summary 可包含：

- `name`
- `impls`
- `methods`

frontend 源码还消费 type reference / impl reference 的细分字段：

- `kind`
- `constructor`
- `arguments`
- `parameters`
- `return_type`
- `error_type`
- `is_async`
- `name`
- `path`
- `self`
- `trait`
- `methods`

SeekMoon `view` 使用 module index 计算 package count 和紧凑 API summary。SeekMoon `api` 使用 module index 解析 package id、package relpath 和 package data URL。

### 4.3 `package_data.json`

`package_data.json` 是单 package API 详情来源。顶层键：

- `name`
- `traits`
- `errors`
- `types`
- `typealias`
- `values`
- `misc`

type entry 字段：

- `name`
- `docstring`
- `signature`
- `loc`
- `methods`
- `impls`

value entry 与 method entry 字段：

- `name`
- `docstring`
- `signature`
- `loc`

`loc` 字段：

- `path`
- `file`
- `line`
- `column`

`signature` 是上游签名字符串，可包含 HTML link。SeekMoon 必须保留 raw signature。若输出 plain text signature，plain text 是 SeekMoon 派生投影。

`docstring` 可为空字符串或以换行开头。空 docstring 是 `missing`，不是解析失败。空数组表示该 package 在该类别没有符号，不能当成 package data 失败。

### 4.4 `resource.json` 与 `resources.json`

当前 Mooncakes 前端源码请求 singular：

```text
resource.json
```

旧 SPA 注释和旧版文档提到 plural：

```text
resources.json
```

plural `resources.json` 在 2026-06-21 的 297 次 bounded check 中未发现 200，代表样本全部 404。singular `resource.json` 是当前前端源码使用的 asset class，但仍是可选来源，不保证每个 package 都存在。

`resource.json` 语义：

- `readme_content`：README 正文。
- `source_files`：source file 链接列表。

状态规则：

```text
resource.json 200 -> resources.state = present
resource.json 404 -> resources.state = unavailable
resource.json non-2xx -> resources.state = failed
resources.json 404 -> 不影响 package/module/source 判断
```

resource asset 缺失不表示 package 缺失。package data、module index、source zip 或 `moon fetch` 仍可能可用。

### 4.5 Source Zip

来源：

```text
https://download.mooncakes.io/user/<owner>/<module>/<version>.zip
```

请求可能重定向到 CloudFront。2026-06-21 代表样本：

| module | version | source zip 状态 |
|---|---|---|
| `moonbitlang/core` | `0.1.20260609+84519ca0a` | 404 |
| `mizchi/markdown` | `0.6.2` | 200 |
| `moonbit-community/cmark` | `0.4.4` | 200 |
| `jaredzhou/pony` | `0.1.1` | 200 |

可下载 source zip 是发布版源码证据。观察到的 zip 内容包括 module config、package config、`.mbt` 源码、README、LICENSE、测试、示例、benchmark、生成的 `.mbti` 等不同组合。

source zip 不是全模块稳定来源。404 表示该 module/version 的 zip 在该路径不可用，不表示 module 不存在。fallback 顺序由 `source` 命令确定：`moon fetch`、本地 registry cache、`~/.moon/lib/core`、GitHub repository。

### 4.6 Skill Assets

Skills API 给出：

- `detail_url`
- `wasm_url`
- `checksum_url`

skill detail 页还读取：

```text
GET /assets/<entry>/SKILL.md
```

Skill asset 字段只属于 skill/executable 对象。`wasm_url` 与 `checksum_url` 不证明 library package 可作为依赖使用。

## 5. 本地工具链与缓存模型

### 5.1 Toolchain Surface

2026-06-21 本地 MoonBit 工具链：

```text
moon 0.1.20260608 (60bc8c3 2026-06-08)
moonc v0.10.0+e66899a54 (2026-06-09)
moonrun 0.1.20260608 (60bc8c3 2026-06-08)
Feature flags enabled: rr_moon_mod,rr_moon_pkg
mooncake-bin 0.1.20260521 (287155a 2026-05-21)
```

`moon --help` 中存在：

- `runwasm`
- `fetch`
- `ide doc`

`moon --help` 中不存在：

- `search`
- `view`
- `audit`
- `outdated`

因此官方 CLI 仍没有通用 package discovery 命令。SeekMoon search/view/audit-like 能力是 SeekMoon 自己的产品面，不是当前 `moon` 命令转发。

`moon update` 输出同时包含：

```text
Registry index updated successfully
Symbols updated successfully
```

`moon update` 更新本地 registry index 与 symbols cache。

### 5.2 Local Registry Index

路径：

```text
~/.moon/registry/index/user/**/*.index
```

当前形态：

- `.index` 文件：1363 个。
- JSONL records：9952 条。
- parse errors：0。
- registry index 是 git repo；2026-06-21 观测 head 为 `7503ed87 update cybershang/agent-telemetry`。

`.index` 文件是 JSON Lines。每一行是一个发布版本记录，不是一个 module 记录。不能把 `.index` 当 JSON 数组读取，不能假设一个文件只有一个版本。

local `ModuleVersionRecord` 原始字段包括：

- `name`
- `version`
- `checksum`
- `license`
- `created_at`
- `repository`
- `readme`
- `keywords`
- `description`
- `deps`
- `source`
- `preferred-target`
- `exclude`
- `warn-list`
- `preferred_target`
- `include`
- `--moonbit-unstable-prebuild`
- `scripts`
- `supported-targets`
- `bin-deps`
- `homepage`
- `alert-list`
- `targets`
- `supported_targets`
- `authors`
- `rule`
- `import`
- `author`
- `preferred-backend`
- `root-dir`
- `keyword`
- `dependencies`
- `link`
- `source-dir`
- `deprecated`

local registry index 比 public `/api/v0/modules` 更丰富，也更松散。SeekMoon 可以用它做 offline/cross-check 来源。默认输出不能倾倒所有历史字段。raw/debug 输出必须保留原始 line。

规范化规则：

| canonical 字段 | 原始字段 |
|---|---|
| `module` | `name` |
| `version` | `version` |
| `source_dir` | `source`, `source-dir` |
| `deps` | `deps`；`dependencies` 只在 shape 合法时另行处理 |
| `preferred_target` | `preferred-target`, `preferred_target`, `preferred-backend` |
| `supported_targets` | `supported-targets`, `supported_targets`, `targets` |
| `keywords` | `keywords`; `keyword` 只作为 raw-aware 兼容输入 |

### 5.3 Local Symbols Cache

路径：

```text
~/.moon/registry/symbols/**/*.symbols
```

2026-06-21 本地状态：

- `.symbols` 文件：3 个。
- JSONL records：1619 条。
- 模块：`moonbitlang/async`、`moonbitlang/quickcheck`、`moonbitlang/x`。
- parse errors：0。

`.symbols` 是 JSON Lines。有两类 record。

meta record：

```json
{
  "kind": "meta",
  "module": "moonbitlang/async",
  "schema_version": 2,
  "version": "0.19.1"
}
```

symbol record：

```json
{
  "attrs": [],
  "doc": "Compute the Sha256 digest in `Bytes` of some `data`. Note that Sha256 is big-endian.",
  "key": "moonbitlang/x/crypto.sha256",
  "kind": "function",
  "name": "sha256",
  "parent": "moonbitlang/x/crypto",
  "pkg": "moonbitlang/x/crypto",
  "sig_": "pub fn[Data : ByteSource] sha256(Data) -> FixedArray[Byte]"
}
```

canonical `SymbolRecord`：

| canonical 字段 | 原始字段 |
|---|---|
| `module` | meta `module` 或文件路径 |
| `version` | meta `version` 或文件路径 |
| `schema_version` | meta `schema_version` |
| `key` | `key` |
| `package` | `pkg` |
| `name` | `name` |
| `parent` | `parent` |
| `kind` | `kind`；允许 string 或 array |
| `signature` | `sig_` |
| `doc` | `doc` |
| `attrs` | `attrs` |

symbols cache 覆盖很小。`moon ide doc sha256` 能命中 `moonbitlang/x`，因为本地有该 symbols 文件；`moon ide doc markdown` 无结果，不能推出 registry 没有 Markdown 包。

### 5.4 `moon ide doc`

`moon ide doc --help` 的语义是搜索：

- current module
- `moonbitlang/core`
- local registry symbol indexes
- optional extra symbol files

它不是 registry package search。SeekMoon 可以把 `moon ide doc` 用作已知 package/API 的补充证据，不能把它作为 `search` 主来源。

有效用途：

- core package API lookup
- 本地 symbols cache 中已有模块的 API lookup
- 用户已选候选后的 API 下钻补充

无效用途：

- 通用 package discovery
- 全 registry API indexing
- registry 缺包判断

### 5.5 `moon fetch`

`moon fetch <module>@<version>` 是当前命令，但 help 标注 unstable。它会在当前项目写入：

```text
<project>/.repos/<owner>/<module>/<version>
```

并使用全局 zip cache：

```text
~/.moon/registry/cache/<owner>/<module>/<version>.zip
```

2026-06-21，`moon fetch mizchi/markdown@0.6.2` 在 temp project 中成功，目标目录为 `.repos/mizchi/markdown/0.6.2`，并使用 `~/.moon/registry/cache/mizchi/markdown/0.6.2.zip`。下载源码包含 `moon.mod.json`、`moon.pkg`、`.mbt`、README、LICENSE、tests/examples/benches 等内容。

SeekMoon 使用 `moon fetch` 时必须明确 mutation boundary。默认运行位置应是 SeekMoon 创建的 temp/probe project，或用户明确认可的项目目录。

### 5.6 `moon runwasm`

`moon runwasm` 支持本地 package 和 Mooncakes Wasm asset。远程 coordinate 形态：

```text
user/module/package@version
user/module@version/package
user/module/package
```

未 pin 版本的 coordinate 会解析 latest。SeekMoon 记录必须保存 pinned coordinate 或解析后的 version，不能只保存未 pin 输入。

cache 路径：

```text
~/.moon/registry/cache/assets/<owner>/<module>/<version>/...
```

当前本地样本：

```text
~/.moon/registry/cache/assets/Yoorkin/cowsay/0.1.0/cowsay.wasm
```

`moon runwasm --dry-run` 不支持 Mooncakes assets。SeekMoon 不得依赖 dry run 来预览远程 asset execution。

## 6. Mooncakes 前端源码约束

Mooncakes 前端源码 snapshot：

- repo：`external/mooncakes.io`
- commit：`4d889ad9a49ed0e0b8b4d8dc58451490de4bb38b`
- commit date：`2026-06-18 18:52:20 +0800`
- subject：`Refine skills marketplace hero copy`

前端源码不是稳定公共 API 合同，但它说明官方页面实际消费哪些字段。

首页模块搜索是 client-side。首页初始请求 highlight buckets 与 statistics；完整 `/api/v0/modules` 在用户输入 filter 或展开更多模块时加载。首页 search index 使用：

- full path / module name
- short name
- author
- keywords
- recency penalty from `created_at`

首页模块 search index 不使用 `description`。如果 SeekMoon 使用 description 搜索，这是 SeekMoon 自己的搜索语义。

首页模块卡片不使用：

- downloads
- license
- repository
- build_status
- target metadata

downloads 与 build_status 来自 docs/manifest 层。Search 默认 list-only 输出不能假装已有这些字段。若 search 输出 downloads/build_status，必须做 manifest enrichment 或把字段状态标为未加载。

Docs 页面消费：

- manifest
- `resource.json`
- `module_index.json`
- `package_data.json`
- build log API

Docs symbol search 在 `module_index.json` 加载后本地执行。package/API 展示由 `package_data.json` 驱动。

Skills marketplace 消费：

- `/api/v0/skills`
- `/api/v0/skills/<entry>`
- `/assets/<entry>/SKILL.md`
- skill `wasm_url`
- skill `checksum_url`

Skill filter 使用：

- `name`
- `module`
- `package`
- `description`

Skills 页面展示 `moon runwasm`，说明 skill 是 execution object，不是 library module。

## 7. Canonical 命令面

SeekMoon v0 命令面：

- `doctor`
- `sync`
- `search`
- `view`
- `api`
- `source`
- `skill`
- `compare`
- `probe`
- `record`
- `report`
- `raw`

输出型命令可支持：

- 默认 pretty text
- `--json`
- `--jq <expr>`
- `--shape`
- `--schema`

不提供：

- `--why`
- `--hints`
- `guide`
- top-level `seekmoon schema search`
- top-level `seekmoon shape search`
- `fields`

这些删除项与当前对象边界重叠。README/help 负责入口学习，`--shape` 和 `--schema` 负责契约，`--jq` 负责机器处理，raw 负责审计逃生口，错误输出负责恢复动作。

### 7.1 `doctor`

目的：检查本地环境是否具备 SeekMoon 所需来源。

字段：

- `toolchain.moon_version`
- `toolchain.moonc_version`
- `toolchain.moonrun_version`
- `toolchain.mooncake_version`
- `commands.runwasm`
- `commands.fetch`
- `commands.ide_doc`
- `commands.search`
- `commands.view`
- `commands.audit`
- `commands.outdated`
- `paths.registry_index`
- `paths.symbols_cache`
- `paths.registry_cache`
- `paths.assets_cache`
- `paths.core_source`
- `network.mooncakes_api`
- `network.mooncakes_assets`

`search/view/audit/outdated` 在当前 `moon` 中应显示为 absent，不是 SeekMoon 失败。

### 7.2 `sync`

目的：创建 dated evidence snapshot。

字段：

- `snapshot.fetched_at`
- `statistics.total_modules`
- `statistics.total_packages`
- `statistics.total_lines`
- `statistics.total_downloads`
- `modules_api.length`
- `skills_api.length`
- `local_index.file_count`
- `local_index.record_count`
- `local_index.head`
- `symbols.file_count`
- `symbols.record_count`
- `toolchain.moon_version`
- `toolchain.moonc_version`
- `toolchain.moonrun_version`

statistics counters 是 snapshot 值。它们随时间变化。

### 7.3 `search`

目的：从 query 产生 library 或 skill 候选。

library search 默认来源是 `/api/v0/modules`。字段：

- `schema`
- `snapshot_id`
- `query.text`
- `query.kind`
- `query.target`
- `results[].rank`
- `results[].module`
- `results[].version`
- `results[].description.status`
- `results[].description.value`
- `results[].keywords.status`
- `results[].keywords.value`
- `results[].license.status`
- `results[].license.value`
- `results[].repository.status`
- `results[].repository.value`
- `results[].is_new`
- `results[].created_at`
- `results[].target.status`，仅在 target 被请求时出现
- `results[].target.value`，仅在 target 有证据时出现
- `results[].match.fields`，默认 pretty text 不显示，JSON/record 可显示

list-only search 不显示 downloads、build_status、has_package、versions_count。manifest-enriched search 可以显示这些字段，但必须记录 enrichment 来源和失败状态。

skill search 来源是 `/api/v0/skills`。字段：

- `results[].module`
- `results[].version`
- `results[].package`
- `results[].name`
- `results[].author`
- `results[].author_avatar.status`
- `results[].author_avatar.value`
- `results[].description.status`
- `results[].description.value`
- `results[].repository.status`
- `results[].repository.value`
- `results[].created_at`

### 7.4 `view`

目的：展示单个 module evidence profile。

字段：

- `module`
- `version`
- `latest_version`
- `downloads`
- `has_package`
- `build_status`
- `metadata.description.status/value/source`
- `metadata.keywords.status/value/source`
- `metadata.repository.status/value/source`
- `metadata.license.status/value/source`
- `metadata.deps`
- `metadata.readme.status/value`
- `metadata.source.status/value`
- `metadata.checksum.status/value`
- `metadata.preferred_target.status/value/source_key`
- `metadata.supported_targets.status/value/source_key`
- `versions[].version`
- `versions[].yanked`
- `versions_count`
- `docs_url`
- `api_index.status`
- `api_index.package_count`
- `target.status/value/source`

`view` 不输出 unsupported future fields。`target` 没证据时为 `unknown`。

### 7.5 `api`

目的：检查 package API 结构。

字段：

- `module`
- `version`
- `package_path`
- `package_relpath`
- `module_index.status`
- `module_index.url`
- `package_data.status`
- `package_data.url`
- `types[]`
- `values[]`
- `traits[]`
- `errors[]`
- `typealias[]`
- `misc[]`
- `docstring.status/value`
- `signature`
- `loc.status/value`
- `methods[]`
- `impls[]`
- `ide_doc.status/output`，仅在使用 `moon ide doc` 时出现

`api` 必须先读 module index，再构造 package data URL。没有 package data 时不能伪造空 API。

### 7.6 `source`

目的：获取或定位发布版源码。

字段：

- `module`
- `version`
- `moon_fetch.status/path/error`
- `source_zip.status/path/final_url/error`
- `local_cache.status/path`
- `core_local_source.status/path`
- `github_source.status/path/url`，仅在使用 GitHub 时出现
- `selected_source.method`
- `selected_source.path`
- `files.summary`

source resolution 优先级必须区分发布版源码与维护仓库源码。GitHub repository 是外部维护信号，不自动等同于当前发布版源码。

### 7.7 `skill`

目的：发现、查看和记录 executable Wasm / skill 条目。

字段：

- `module`
- `author`
- `author_avatar.status/value`
- `version`
- `package`
- `name`
- `detail_url`
- `wasm_url`
- `checksum_url`
- `metadata.name.status/value`
- `metadata.description.status/value`
- `repository.status/value`
- `created_at`
- `skill_md.status/value`
- `runwasm_coordinate`
- `wasm_asset.status`
- `checksum_asset.status`

`runwasm_coordinate` 是 derived。未 pin 输入必须解析并记录 version。

### 7.8 `compare`

目的：把多个候选放在同一证据表面。

字段来源只包括已经加载的 evidence：

- modules list 字段
- manifest 字段
- module index/package data 字段
- source 字段
- probe 字段
- GitHub enrichment 字段

`compare` 不引入 quality score、security score、reverse dependency count 或 provenance。排序可有 search rank，但 rank 是 SeekMoon 算法输出，不是质量指标。

### 7.9 `probe`

目的：在本地项目或 temp project 中验证候选。

字段：

- `module`
- `version`
- `target`
- `probe_path`
- `moon_new.status`
- `moon_add.status`
- `moon_check.status`
- `moon_test.status`
- `moon_check_target.status`
- `moon_build_target.status`
- `result.status`
- `logs.path`

probe 字段都是 local derived evidence。probe 成功不改变上游事实；probe 失败也不证明 registry 元数据错误，只证明该环境、该命令、该 target 下失败。

### 7.10 `record`

目的：存储采用判断和证据引用。

字段：

- `record_id`
- `created_at`
- `module`
- `version`
- `conclusion`
- `note`
- `evidence_refs[]`
- `not_confirmed[]`

`conclusion` 枚举：

- `adopt`
- `adopt-with-adapter`
- `continue-verification`
- `contribute-upstream`
- `fork`
- `build-own`
- `reject-for-now`

record 是 SeekMoon 自有记录。它不改变上游字段。

### 7.11 `report`

目的：输出可复现审计记录。

字段：

- `goal`
- `date`
- `toolchain`
- `data_sources[]`
- `query`
- `candidates[]`
- `local_validation[]`
- `cannot_confirm[]`
- `conclusion`

`data_sources[]` 只列实际使用的来源。没有使用 GitHub 就不列 GitHub。没有执行 probe 就不列 local validation。

### 7.12 `raw`

目的：暴露原始 payload。

允许 raw surface：

- `raw modules`
- `raw statistics`
- `raw manifest <module>`
- `raw module-index <module@version>`
- `raw package-data <module@version> <package>`
- `raw resource <module@version> <package>`
- `raw skills`
- `raw skill <entry>`
- `raw local-index <module>`
- `raw symbols <module>`

raw 输出不能改写字段名。raw 输出可以附带 fetch status 和 source path。

## 8. Unsupported Fields

以下字段或能力没有当前 v0 来源，不进入当前输出：

| 字段或能力 | 边界 |
|---|---|
| `quality_score`, `Mooncake Score`, quality `score` | 没有上游字段；search rank 不是质量分。 |
| `advisory`, `audit_status`, vulnerability status | 没有当前 advisory/audit 来源。 |
| `outdated` package status | 没有当前上游命令或字段；未来可由项目依赖和 manifest 派生。 |
| `dependents`, `reverse_dependencies`, `reverse_dependency_count` | 没有当前来源。 |
| `provenance`, `attestation`, `signature`, `publisher_identity`, `verified_publisher` | 当前 checksum 字符串或 checksum URL 不等于 provenance。 |
| `sbom` | 没有当前来源。 |
| `ci_status` 作为 registry 字段 | 只能作为 GitHub enrichment，且必须实际请求 GitHub。 |
| `tests_present`, `examples_present` 作为 registry 字段 | 只能由 source/GitHub 检查派生。 |
| `docs_build_status` 独立字段 | 当前 manifest 是 `build_status`。 |
| `server_search_url`, `search_deeplink`, `/search?q=` | 没有当前稳定来源。 |
| `target_supported=true/false` 无证据布尔值 | 只能输出 metadata-derived、probe-derived 或 `unknown`。 |
| `moon search`, `moon view`, `moon audit`, `moon outdated` | 当前 `moon --help` 不存在。 |
| `--why`, `--hints`, `guide`, top-level `schema/shape/fields` | 不属于 v0 命令面。 |

## 9. 归一化规则

### 9.1 空值

空字符串、空数组和 null 不能直接当作有效业务值。

| 来源字段 | 空值状态 |
|---|---|
| modules `description` | `missing` |
| modules `keywords` | `missing` |
| modules `repository` | `missing` |
| modules `license` | `missing` |
| skills `author_avatar` | `missing` |
| skills `package` | root/default package marker，不是 missing |
| skills `metadata.description` | `missing` |
| manifest `build_status: null` | `missing` 或 nullable present，由 schema 决定；不得当 success。 |

### 9.2 Target

target 字段没有单一上游真相。读取顺序：

1. manifest metadata exact raw keys。
2. local registry index raw keys。
3. source config。
4. local probe。

规范化输入：

- `preferred-target`
- `preferred_target`
- `preferred-backend`
- `supported-targets`
- `supported_targets`
- `targets`

规范化输出：

- `preferred_target`
- `supported_targets`
- `target.status`
- `target.source`

`supported-targets` 当前可能是 string，也可能在历史数据中是 array。canonical 存储必须接受 string 或 array，并可派生 array projection。

### 9.3 Source

发布版源码优先级：

1. source zip 200。
2. `moon fetch` 成功。
3. local registry cache zip。
4. `~/.moon/lib/core`，仅适用于 core。
5. GitHub repository，标记为 external maintenance source。

source zip 404、resource 404、symbols miss 都不是 module absence。

### 9.4 Signature 与 Docstring

`package_data.json` 的 `signature` 保留 raw。HTML stripping 是派生投影。`docstring` 空字符串是 `missing`，但 entry 本身仍存在。

### 9.5 Version

Manifest `version`、`latest_version`、Skills API `version`、local registry index `version` 都必须保存来源。`versions_count` 只能由 manifest `versions.length` 派生。

runwasm 未 pin coordinate 必须解析 latest 并记录解析后的 version。

## 10. 实现不变量

1. 公共 API payload 和 asset payload 必须保留 raw。
2. 默认输出不得显示没有来源的未来字段。
3. list-only search 不得显示 manifest-only 字段。
4. skill 与 library module 必须分开建模。
5. `module_index.json` traversal 使用 `childs`。
6. `package_data.json` URL 必须由 module index package path 派生。
7. `resource.json` 是可选 README/source-files asset；404 不代表 package 缺失。
8. source zip 是可用时的发布版源码证据；404 不代表 module 缺失。
9. `moon fetch` 会修改当前项目；默认只在 temp/probe boundary 内执行。
10. `moon ide doc` 不是 package discovery。
11. local registry index 与 symbols cache 都是 JSONL。
12. target support 没有证据时只能是 `unknown`。
13. GitHub 字段只在实际请求 GitHub 后出现。
14. Search rank 是 SeekMoon 算法输出，不是质量评分。
15. `missing`、`unknown`、`failed`、`unavailable` 必须按状态词定义使用。

## 11. v0 可实施数据字典

### `ModuleSummary`

| 字段 | 类型 | 来源 |
|---|---|---|
| `module` | string | modules `name` |
| `version` | string | modules `version` |
| `description` | status/value | modules `description` |
| `keywords` | status/value | modules `keywords` |
| `repository` | status/value | modules `repository` |
| `license` | status/value | modules `license` |
| `is_new` | boolean | modules `is_new` |
| `created_at` | string | modules `created_at` |
| `raw` | object | modules item |

### `ManifestProfile`

| 字段 | 类型 | 来源 |
|---|---|---|
| `module` | string | manifest `module` / `name` |
| `version` | string | manifest `version` |
| `latest_version` | string | manifest `latest_version` |
| `downloads` | integer | manifest `downloads` |
| `has_package` | boolean | manifest `has_package` |
| `build_status` | string/null | manifest `build_status` |
| `metadata` | raw object + normalized projection | manifest `metadata` |
| `versions` | array | manifest `versions` |
| `versions_count` | integer | derived |

### `ModuleIndexTree`

| 字段 | 类型 | 来源 |
|---|---|---|
| `name` | string | module index node |
| `package` | object/null | module index node |
| `childs` | array | module index node |
| `package.path` | string | package node |
| `package.relpath` | string | derived |
| `traits/errors/types/typealias/values/misc` | arrays | package node |

### `PackageData`

| 字段 | 类型 | 来源 |
|---|---|---|
| `name` | string | package data |
| `traits` | array | package data |
| `errors` | array | package data |
| `types` | array | package data |
| `typealias` | array | package data |
| `values` | array | package data |
| `misc` | array | package data |
| `entry.name` | string | package data entry |
| `entry.docstring` | status/value | package data entry |
| `entry.signature` | string | package data entry |
| `entry.loc` | status/value | package data entry |
| `entry.methods` | array | type entry |
| `entry.impls` | array | type entry |

### `SkillEntry`

| 字段 | 类型 | 来源 |
|---|---|---|
| `module` | string | Skills API |
| `author` | string | Skills API |
| `author_avatar` | status/value | Skills API |
| `version` | string | Skills API |
| `package` | string/root marker | Skills API |
| `name` | string | Skills API |
| `detail_url` | string | Skills API |
| `wasm_url` | string | Skills API |
| `checksum_url` | string | Skills API |
| `metadata.name` | status/value | Skills API |
| `metadata.description` | status/value | Skills API |
| `repository` | status/value | Skills API |
| `created_at` | string | Skills API |
| `runwasm_coordinate` | string | derived |

### `LocalModuleVersionRecord`

| 字段 | 类型 | 来源 |
|---|---|---|
| `module` | string | index `name` |
| `version` | string | index `version` |
| `checksum` | status/value | index `checksum` |
| `created_at` | status/value | index `created_at` |
| `description` | status/value | index `description` |
| `keywords` | status/value | index `keywords` |
| `license` | status/value | index `license` |
| `repository` | status/value | index `repository` |
| `readme` | status/value | index `readme` |
| `source_dir` | status/value | index `source` / `source-dir` |
| `deps` | status/value | index `deps` |
| `preferred_target` | status/value/source_key | target spelling variants |
| `supported_targets` | status/value/source_key | target spelling variants |
| `raw` | object | index JSONL line |

### `SymbolRecord`

| 字段 | 类型 | 来源 |
|---|---|---|
| `module` | string | meta/file path |
| `version` | string | meta/file path |
| `schema_version` | integer | meta |
| `key` | string | symbol record |
| `package` | string | `pkg` |
| `name` | string | symbol record |
| `parent` | string/null | symbol record |
| `kind` | string/array | symbol record |
| `signature` | string/null | `sig_` |
| `doc` | string/null | symbol record |
| `attrs` | array | symbol record |

## 12. 结论

SeekMoon v0 的数据模型由当前 Mooncakes API、Mooncakes assets、本地 MoonBit 工具链、本地 registry/cache、当前项目上下文、可选 GitHub enrichment 和 SeekMoon 自有记录构成。

library module discovery 与 skill/runwasm execution discovery 是两个对象。library 候选来自 modules、manifest、assets、本地 registry、source 和 probe。skill 候选来自 Skills API、SKILL.md、Wasm/checksum asset 和 `moon runwasm` coordinate。

默认输出是低噪声事实投影。JSON 输出承载完整状态。`--shape` 与 `--schema` 承载契约。`raw` 承载原始来源。错误输出承载恢复动作。

当前 v0 不输出没有来源的愿望字段。当前 schema 内的空值是 `missing`，当前证据不能回答的问题是 `unknown`，请求或命令失败是 `failed`，可选来源不存在是 `unavailable`，SeekMoon 计算值是 `derived`。

按这些边界实现时，SeekMoon 可以稳定支持 `doctor`、`sync`、`search`、`view`、`api`、`source`、`skill`、`compare`、`probe`、`record`、`report` 和 `raw`，同时避免把未来生态能力、官方前端行为、MoonBit 本地命令、发布版源码、维护仓库源码和 SeekMoon 自有判断混成同一个事实层。
