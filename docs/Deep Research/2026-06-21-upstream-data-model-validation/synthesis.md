# SeekMoon 上游数据模型复核总汇报

日期：2026-06-21  
批次目录：`projects/seekmoon/docs/Deep Research/2026-06-21-upstream-data-model-validation/`  
输入材料：用户指定的 3 份 MoonBit 调查/SOP 文档、`002-raw消费者侧 MoonBit 包发现 CLI 设计.md` 聊天记录、5 份本轮 subagent 调查报告。

## 一、最终判断

这次复核后，SeekMoon 当前数据字典必须收束为“已确认上游事实 + 已确认本地事实 + SeekMoon 明确派生事实”三类，不能再把未来生态愿望或设计占位字段混进当前输出。

当前可以进入 v0 数据模型的上游事实源是：

- Mooncakes Modules API：`https://mooncakes.io/api/v0/modules`
- Mooncakes Statistics API：`https://mooncakes.io/api/v0/modules/statistics`
- Mooncakes Manifest API：`https://mooncakes.io/api/v0/manifest/<owner>/<module>`
- Mooncakes Skills API：`https://mooncakes.io/api/v0/skills`
- Mooncakes `module_index.json`
- Mooncakes `package_data.json`
- Mooncakes source zip，前提是当前模块/版本实际可下载
- 本地 `~/.moon/registry/index` JSON Lines
- 本地 `~/.moon/registry/symbols` JSON Lines
- 本地 `moon` 命令行为与缓存
- 当前 MoonBit 项目上下文
- GitHub 维护信号，前提是 manifest/modules 里有 repository，并且工具实际请求过
- SeekMoon 自己产生的 session、probe、record、report

当前不能进入 v0 当前输出的字段或能力是：

- `quality_score` / `Mooncake Score`
- `advisory` / `audit_status` / vulnerability status
- reverse dependents / dependents
- provenance / attestation / verified publisher / SBOM
- server-side search endpoint
- Web 搜索深链
- 无证据的 target support 布尔值
- `moon search` / `moon view` / `moon audit` / `moon outdated` 作为当前工具链能力
- 之前聊天中已经删除的公共面：`--why`、`--hints`、`guide`、顶层 `schema/shape/fields`

这些内容可以作为未来观察项或后续设计对象，但不能在当前 `search/view/api/source/skill/compare/probe/record/report` 输出 schema 中出现。没有当前来源的字段不显示，不写 `unknown`。只有当前 schema 已经定义、当前查询确实需要、但证据源不能回答时，才显示 `unknown`。当前源有字段位置但值为空、缺失或 null 时，才显示 `missing`。

## 二、公共 API 当前模型

`/api/v0/modules` 当前返回顶层数组，本轮 live fetch 为 1350 个 module。每个 item 观察到 8 个键：

- `name`
- `version`
- `description`
- `keywords`
- `repository`
- `license`
- `is_new`
- `created_at`

这些键在当前样本里都存在，但缺失内容用空字符串或空数组表达，不是字段缺失。具体口径：

- `description`：327 个为空字符串。
- `keywords`：417 个为空数组。
- `repository`：292 个为空字符串。
- `license`：14 个为空字符串。

所以 SeekMoon ingest 时要保留 raw 值，同时把空字符串/空数组转成当前命令 schema 里的 `missing` 状态。不能把它们当成字段不存在，也不能把空字符串显示成一个有效描述、有效 repository 或有效 license。

`/api/v0/modules?search=markdown` 与 `/api/v0/modules?search=cowsay` 当前仍返回完整 1350 项数组，解析后与 `/api/v0/modules` 相同，第一项仍是 `0Ayachi0/elk`。因此它不是服务端搜索。SeekMoon 的 `search` 必须明确为本地过滤：拉全量 modules，再按自己的搜索语义过滤。这个结论没有过期。

`/api/v0/modules/statistics` 当前返回对象：

- `total_modules`
- `total_packages`
- `total_lines`
- `total_downloads`

本轮不同报告 fetch 时间略有差异，所以 downloads/lines 有小幅漂移，这是正常的易变计数。它只适合作为 snapshot 字段，不适合作为永久事实。

`/api/v0/manifest/<module>` 当前代表样本的顶层键是：

- `name`
- `module`
- `version`
- `latest_version`
- `downloads`
- `has_package`
- `build_status`
- `metadata`
- `versions`

`versions_count` 不是上游字段，只能由 `versions.length` 派生。`versions[]` 观察到 `version` 与 `yanked`。`build_status` 不是只有 `success`，样本里有 `success`、`legacy`、`null`。因此 build status 是一个 Mooncakes 构建/文档信号，不是“包可采用”结论，也不是 target build 成功结论。

manifest 的 `metadata` 是开放对象，不是稳定闭合 schema。当前可进入 raw metadata 的字段包括：

- `name`
- `version`
- `checksum`
- `created_at`
- `license`
- `repository`
- `description`
- `keywords`
- `readme`
- `source`
- `deps`
- `preferred-target`
- `preferred_target`
- `supported-targets`
- `warn-list`

其中 target 相关字段非常不稳定。`/modules` 和 `/skills` 都没有 top-level target 字段。代表 manifest 中，`mizchi/markdown` 有 `metadata["preferred-target"] = "js"`，`Yoorkin/cowsay` 有 `metadata.preferred_target = "wasm"` 和 `metadata["supported-targets"] = "wasm"`。这说明 SeekMoon 可以规范化 known spellings，但不能把 target compatibility 作为稳定上游字段。用户传 `--target js` 时，如果没有 manifest metadata 或 probe 证据，就只能显示 `target: unknown`，不能显示 supported/unsupported。

`/api/v0/skills` 当前返回顶层数组，本轮为 70 条 skill。每条 item 的 top-level key 是：

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

`wasm_url`、`checksum_url`、`detail_url` 当前 70 条都有非空值。`package` 字段存在但可能为空字符串，空字符串表示 root/default executable package 语义，不能丢。`metadata` 是对象，但内部字段稀疏，`metadata.description` 大量为空字符串，`metadata.name` 也经常为空或不存在。skill 必须作为 execution object，不能混入 library dependency 候选结论。

## 三、Assets 当前模型

`module_index.json` 是当前 API 下钻的一等事实源。它是树，不是扁平列表。当前实际 key 是：

- `name`
- `package`
- `childs`

注意 spelling 是 `childs`，不是 `children`。package node 中包含：

- `path`
- `traits`
- `errors`
- `types`
- `typealias`
- `values`
- `misc`

`package.path` 是完整包路径，例如 `moonbitlang/core/argparse`。构造 package data URL 时，不能直接把完整 path 拼进去，而要从 module path 中剥离 package relpath。例子：

- module path：`moonbitlang/core`
- package path：`moonbitlang/core/argparse`
- relpath：`argparse`
- package data URL：`/assets/moonbitlang/core@<version>/argparse/package_data.json`

root package 的 relpath 为空，例如 `mizchi/markdown` root package 的 package data URL 是 `/assets/mizchi/markdown@0.6.2/package_data.json`。因此 SeekMoon 的 `api` 命令必须先读 module index，不能靠猜 URL。

`package_data.json` 是单 package API 详情的一等事实源。当前 top-level key：

- `name`
- `traits`
- `errors`
- `types`
- `typealias`
- `values`
- `misc`

type entry 观察到：

- `name`
- `docstring`
- `signature`
- `loc`
- `methods`
- `impls`

value/method entry 观察到：

- `name`
- `docstring`
- `signature`
- `loc`

`loc` 包含：

- `path`
- `file`
- `line`
- `column`

`signature` 可能是带 HTML link 的签名字符串。SeekMoon 可以保留 raw signature，也可以派生 plain text signature，但 plain text 是 SeekMoon 派生投影，不是上游字段。

资源 asset 存在一个重要漂移点：前端源码当前请求 singular `resource.json`，旧 SPA 注释和本轮 Task 02 检查的是 plural `resources.json`。Task 02 对 297 个 plural `resources.json` 请求没有发现 200，代表样本全部 404。结论是：资源 asset 当前不能作为可靠输入。对于 v0，`resources/resource` 只能作为 optional source：200 才解析，404 只表示 resource asset unavailable，不表示 package 不存在，也不表示没有 API 文档，更不表示没有 README，因为 source zip 里可能有 README。这个点后续实现前建议单独验证 singular `resource.json` live 行为。

source zip 当前可用但不普遍。`mizchi/markdown`、`moonbit-community/cmark`、`jaredzhou/pony` 可下载，zip 内含 module config、package config、`.mbt`、README、LICENSE、tests/examples/benches 等不同组合。`moonbitlang/core` 最新版本 source zip 按同样 URL 404。结论：source zip 是“可用时的发布版源码证据”，不是全模块稳定来源。失败时应回退到 `moon fetch`、本地 cache、`~/.moon/lib/core` 或 GitHub repository；GitHub 是维护/协作/未发布状态证据，不是发布版源码的唯一事实源。

## 四、本地工具链与本地数据模型

当前本机 MoonBit 工具链已经是 2026-06-08 / v0.10.0 时代：

- `moon 0.1.20260608`
- `moonc v0.10.0+e66899a54`
- `moonrun 0.1.20260608`

旧报告里本机 v0.9.3 状态已经过期。当前 `moon --help` 有：

- `runwasm`
- `fetch`
- `ide doc`

仍没有：

- `search`
- `view`
- `audit`
- `outdated`

所以官方 CLI 仍没有通用 package discovery loop，但 runwasm/fetch/ide doc 已经是当前事实。

`moon update` 当前同时更新 registry index 和 symbols cache。registry index 位于 `~/.moon/registry/index/user/**/*.index`，是 JSON Lines，不是 JSON 数组。本轮统计为 1363 个 `.index` 文件，9952 条 version record。每条 line 是一个发布版本记录，不是一个 module 记录。常见字段包括：

- `name`
- `version`
- `checksum`
- `created_at`
- `description`
- `keywords`
- `license`
- `repository`
- `readme`
- `deps`
- `source`
- `preferred-target`
- `preferred_target`
- `supported-targets`
- `supported_targets`
- `targets`

这里比 public `/modules` 更丰富，但 schema 更散，有历史字段和拼写变体。SeekMoon 可以用它做 local/offline cross-check，但必须保留 raw，并做规范化派生，不要把 index 的所有历史字段暴露到默认输出。

symbols cache 位于 `~/.moon/registry/symbols/**/*.symbols`，也是 JSON Lines。当前只有 3 个模块：`moonbitlang/async`、`moonbitlang/quickcheck`、`moonbitlang/x`，共 1619 条。symbol record 有：

- `key`
- `pkg`
- `name`
- `parent`
- `kind`
- `sig_`
- `doc`
- `attrs`

`sig_` 可以规范化为 `signature`。但 symbols cache 覆盖很稀疏，命中失败只能是 unknown，不能说明 registry 没有相关包。

`moon ide doc` 当前搜索范围是 current module、`moonbitlang/core`、本地 registry symbol indexes。它不是通用 registry package search。`moon ide doc markdown` 无结果，但 registry 中有 Markdown 包。`moon ide doc sha256` 可以命中 `moonbitlang/x`，因为本地 symbols cache 有它。所以 SeekMoon 可把 `moon ide doc` 作为已知候选/API 下钻补充，不可作为 search 主入口。

`moon fetch mizchi/markdown@0.6.2` 在 temp 项目里实测成功，写入当前项目 `.repos/mizchi/markdown/0.6.2`，并使用全局 cache `~/.moon/registry/cache/mizchi/markdown/0.6.2.zip`。但 `fetch` help 标注 unstable，而且会修改当前项目目录。因此 SeekMoon 使用它时应默认放在显式 temp/probe 或用户认可的项目 context 内。

`moon runwasm` 当前支持 Mooncakes coordinate，未 pin 版本会解析 latest。帮助文本给出的 cache 路径是 `$MOON_HOME/registry/cache/assets`，当前本地已有 `~/.moon/registry/cache/assets/Yoorkin/cowsay/0.1.0/cowsay.wasm`。`--dry-run` 不支持 Mooncakes assets。SeekMoon 的 skill 记录应保存 pinned coordinate 或解析后的 version，不能只保存 unpinned latest 输入。

## 五、官方前端源码给出的约束

Mooncakes 前端源码 commit 为 `4d889ad9a49ed0e0b8b4d8dc58451490de4bb38b`，2026-06-18，说明站点源仍活跃。

首页模块搜索是 client-side。首页初始请求 highlight buckets 与 statistics，完整 `/api/v0/modules` 在用户输入 filter 或 show more 时拉取。前端搜索 index 使用：

- path/full name
- short name
- author
- keywords
- recency penalty

它不把 `description` 放进 home search index。SeekMoon 如果默认搜索 `description`，这是 SeekMoon 自己的产品语义，不是 Mooncakes 前端 parity。这个选择可以做，但必须命名为 SeekMoon semantics。

首页模块卡片不使用 downloads、license、repository、build_status、targets。downloads/build_status 是 docs/manifest 层面的字段，不是 modules list 字段。Search 默认结果如果显示 downloads/build，需要额外拉 manifest 或明确是 enriched mode；否则默认 list-only search 不应伪装有这些字段。

docs 页消费 manifest、module_index、package_data。前端源码确认 `module_index.json` 使用 `childs`，`package_data.json` 使用 `name/types/traits/errors/typealias/values/misc` 这组模型。

skill marketplace 是单独页面，消费 `/api/v0/skills`，并展示 `moon runwasm`、`SKILL.md`、wasm、checksum。skill filter 使用 `name/module/package/description`，与 library module search 分离。这进一步支持 SeekMoon 必须把 `skill` 做成单独对象，不混进 library adoption。

## 六、当前命令字典收束

当前合法命令面：

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

每个输出型命令可有：

- 默认 pretty text
- `--json`
- 内置 `--jq <expr>`
- `--shape`
- `--schema`

不提供：

- `--why`
- `--hints`
- `guide`
- top-level `seekmoon schema search`
- top-level `seekmoon shape search`
- `fields`

原因不是“暂时不做”，而是消费者动作已被 README/help、`--shape`、`--schema`、`--jq`、默认输出、错误恢复分别覆盖；再增加这些面会重叠。

默认输出应低噪声。正常业务输出只投影当前动作需要的事实，不输出 Next/Notes。工作流牵引由命令结构、help、README、状态 session、编号引用和稳定输出形状完成。错误输出可以提供恢复动作，因为错误是工作流断点。

## 七、v0 数据模型状态规则

必须区分：

- `present`：字段有值，且当前 source 成功提供。
- `missing`：当前 source/schema 有字段位置，但值为空、null、空数组或缺失。
- `unknown`：当前问题需要判断，但现有证据源没有回答。
- `failed`：尝试请求/命令执行失败。
- `unavailable`：该对象当前没有这个来源，或该 optional source 不存在。
- `derived`：SeekMoon 自己计算得出，不是上游字段。
- `unsupported`：当前没有来源，不进入输出。

典型规则：

- description 空字符串：`missing`。
- repository 空字符串：`missing`。
- license 空字符串：`missing`。
- 用户指定 `--target js`，manifest 没 target metadata，也没 probe：`target unknown`。
- `resources.json` 404：`resources unavailable`，不是 package missing。
- symbols cache 查不到 markdown：`unknown`，不是 no package。
- source zip 404：`source_zip unavailable`，不是 module missing。
- `quality_score`：unsupported，不显示。
- `advisory`：unsupported，不显示。

## 八、对旧资料的校正

旧资料中正确保留的判断：

- 官方仍无通用 `moon search/view/audit/outdated`。
- `/api/v0/modules?search=` 仍不能当服务端搜索。
- Mooncakes API/assets/toolchain 已经足够支撑消费者侧 adapter。
- `moon ide doc` 是 API 下钻，不是 package discovery。
- `moon fetch` 是发布源码获取入口，但不等于依赖采纳。
- skill/runwasm 是 execution loop，不是 library dependency loop。
- future field 不应显示；当前 schema 内缺值才 missing/unknown。

需要修正的判断：

- 本机工具链版本不是旧的 v0.9.3，而是 v0.10.0 时代。
- search 默认不能显示 downloads/build，除非 search 实现明确做 manifest enrichment。
- target compatibility 不能作为稳定 list 字段。
- `resources.json` 不能当可靠 README/resources source，且 singular/plural 命名存在前端源码与旧注释不一致。
- GitHub 不能作为发布版源码唯一事实源。
- `module_index` 的 tree key 是 `childs`。
- registry index 是 version-record JSONL，不是 module list。

## 九、最终可实施数据字典

`sync`：

- snapshot time
- statistics counters
- modules API length
- local index count/head
- symbols count
- moon/moonc/moonrun versions

`search` library：

- rank derived
- module from modules.name
- version
- description status/value
- keywords status/value
- license status/value
- repository status/value
- is_new
- created_at
- target status only when target requested
- match fields only in JSON/record, not default pretty text

`view`：

- module/version/latest_version/downloads/has_package/build_status
- metadata description/keywords/repository/license/deps/readme/source/checksum/preferred target
- versions and yanked entries
- versions_count derived
- docs_url derived
- module_index status/package_count
- target evidence status

`api`：

- module/version/package_path
- module_index status
- package_data status
- types/values/traits/errors/typealias/misc
- docstring/signature/loc/methods/impls
- optional `moon ide doc` output as local supplemental evidence

`source`：

- moon_fetch status/path/error
- source_zip status/path/error
- local_cache status/path
- core_local_source status/path
- selected_source method/path derived
- file summary derived

`skill`：

- module/author/author_avatar/version/package/name/detail_url
- wasm_url/checksum_url
- metadata description/name status
- repository status
- created_at
- runwasm coordinate derived
- wasm/checksum asset status only if checked

`compare`：

- selected candidates
- manifest/list/API/source/probe/GitHub evidence actually loaded
- no unsupported future score fields

`probe`：

- module/version/target/probe path
- moon add/check/test/build statuses
- logs
- aggregate status

`record/report`：

- controlled conclusion enum
- evidence references
- cannot-confirm list
- data sources actually used

`raw`：

- raw modules/statistics/manifest/module-index/package-data/resources/skills/local-index/symbols payloads

## 十、最终结论

SeekMoon v0 的数据模型现在可以稳定推进，但必须严格执行三个边界。

第一，当前事实边界：公共 API、assets、本地 registry/toolchain、GitHub enrichment、probe/record 是事实源；质量评分、安全审计、反向依赖、provenance、server-side search 是未来愿望，不进入当前 schema。

第二，对象边界：library module discovery 与 skill/runwasm execution discovery 是两个对象，不能混成一个采纳结论。skill 可以搜索、查看、运行、记录，但不是 library dependency。

第三，投影边界：默认 pretty text 低噪声，JSON 完整表达当前命令对象，`--shape/--schema` 表达契约，`--jq` 做机械处理，`raw` 做审计逃生口。正常业务输出不做教程，README/help 做入口学习，错误输出做恢复。

这次复核的结果不是“字段大概没问题”，而是明确了哪些字段能进、哪些字段不能进、哪些字段必须派生、哪些字段只能 unknown、哪些字段是 missing、哪些来源失败不能被误读为对象不存在。后续实现 SeekMoon 时，可以直接按本汇报和 5 份报告建 v0 ingestion schema。
