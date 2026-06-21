<Assistant id="1">
# MoonBit 包搜索发现 SOP 盲区核验与修订汇报

## 一、当前目光摆在哪里

我完整读完了你指定的五份资料：

- `notes/关于MoonBit/MoonBit 包发现闭环动向调查.md`
- `notes/关于MoonBit/MoonBit 语言、工具链与包生态调查报告.md`
- `projects/seekmoon/docs/关于moonbit/MoonBit 包生态定位论证汇报.md`
- `projects/seekmoon/docs/关于moonbit/001-raw.md`
- `notes/关于MoonBit/关于包搜索与发现/MoonBit包搜索发现与生态查重SOP.md`

同时按要求恢复了聊天记录里的上下文。当前目光不是作者侧健康门禁，不是 Mooncake Score，不是官方全生态治理，不是安全审计体系，也不是 Web 平台改造。当前对象非常明确：**消费者侧本地 CLI 包发现适配器**。

这个对象面对的真实摩擦是：MoonBit 生态里已经有 Mooncakes registry、modules API、manifest API、docs 页面、本地 registry index、`moon add`、`moon check`、`moon fetch`、`moon ide doc`、`moon runwasm`、skills API、assets 文档 JSON 等事实资源，但是这些资源没有被组织成一个顺滑的消费者工作流。用户为了找包、查包、看 API、看源码、比较候选、做本地验证，需要自己在 Web、curl、jq、本地 index、GitHub、Mooncakes docs、probe 项目之间来回拼接。

所以当前要确认的不是“MoonBit 有没有包管理”这种低层事实，也不是“官方有没有完整 discovery loop”这种上一轮已经确认过的问题，而是：**现有 SOP 是否准确覆盖了当前消费者侧可用事实源；是否漏掉了官方已经公开给 agent / CLI / 工具链使用的文档、源码、API、symbols、skills 入口；这些入口是否会改变消费者侧 CLI 工具的对象边界和数据源模型。**

结论先放前面：

1. SOP 的搜索主干没有过期。`/api/v0/modules?search=markdown` 仍然不执行服务端过滤，仍返回完整 modules 数组。消费者侧通用 package search 仍需要拉全量 `/api/v0/modules` 后本地过滤，或者依赖 Web 客户端搜索。
2. SOP 的候选下钻层存在实质盲区。`001-raw.md` 提醒的 Mooncakes docs HTML 注释不是废话，官方确实给了 `module_index.json`、`package_data.json`、`resource.json`、source zip、`moon ide` 等程序化入口。其中 `module_index.json` 和 `package_data.json` 已经能直接拿到 module 内 package、类型、函数、方法、docstring、源码位置。这些不应继续被 SOP 漏掉。
3. `moon ide doc` 已经是官方本地 API 文档入口，可以查 core / 已安装依赖 / symbols cache 中的导出 API，但它不是通用 registry 包搜索。它补的是“已知候选的 API 下钻”，不是“从不知道包名到发现包”。
4. `moon fetch` 是官方 CLI 直接获取发布源码到 `.repos` 的入口，实测成功。它比手写 source zip 更适合消费者侧工具封装。SOP 原来把 GitHub 放在“源码核验入口”的核心位置，这是不准确的。GitHub 应主要承担维护、协作、CI、issue、release、未发布源码核验；发布版源码应优先用 Mooncakes / `moon fetch` / registry cache。
5. Skills API 已经存在，当前返回 70 条可执行 Wasm / skill entry。它不能和 library package 混成同一类候选，但消费者侧 CLI 必须把它作为 `kind=skill` 或 `kind=executable` 的候选源。SOP 原来没收进来，是盲区。
6. 本地 `~/.moon/registry/symbols/**/*.symbols` 存在，但当前只覆盖少数模块，不是全 registry 符号索引。它是补充搜索源，不是通用搜索替代品。
7. 我已经直接修订了 SOP，新增了这些入口、命令、失败语义和注意事项。修改文件是：`notes/关于MoonBit/关于包搜索与发现/MoonBit包搜索发现与生态查重SOP.md`。

## 二、逐项实验事实

### 1. 本机工具链状态

我首先确认了当前工具链，不沿用旧报告里 2026-05-22 的工具链状态。

本机现在是：

```text
moon 0.1.20260608 (60bc8c3 2026-06-08)
moonc v0.10.0+e66899a54 (2026-06-09)
moonrun 0.1.20260608 (60bc8c3 2026-06-08)
Feature flags enabled: rr_moon_mod,rr_moon_pkg
```

`moon --help` 里已经有 `runwasm`，但仍然没有 `search`、`view`、`audit`、`outdated`。这点继续支持之前的结论：官方 CLI 已经有包消费、发布、运行、安装、依赖树等能力，但还没有通用 package discovery CLI 闭环。

`moon ide --help` 显示存在：

```text
moon ide doc               search exported APIs and documentation
```

`moon doc --help` 显示 `moon doc [SYMBOL]` 的符号查询已经 deprecated，并明确指向：

```text
Use `moon ide doc <SYMBOL>` instead.
```

这个事实对 SOP 很关键。`001-raw.md` 里提到 “Prefer moon ide” 不是虚假的方向，当前本地工具链确实有这个入口。

### 2. Mooncakes docs 页面是 SPA，直接 curl 只能拿到外壳

我请求：

```bash
curl -sSL 'https://mooncakes.io/docs/moonbitlang/core/argparse'
```

返回的是 SPA 外壳 HTML，其中注释明确写着：

```text
This app is a SPA. Documentation content is not available in index.html
or static JS bundles. Use the API/assets instead of scraping rendered pages.
```

它列出的入口包括：

```text
GET /api/v0/modules
GET /api/v0/user/<username>
GET /api/v0/manifest/<username>/<mod-relpath>
GET /assets/<username>/<mod-relpath>@<version>/module_index.json
GET /assets/<username>/<mod-relpath>@<version>/<pkg-relpath>/resource.json
GET /assets/<username>/<mod-relpath>@<version>/<pkg-relpath>/package_data.json
https://download.mooncakes.io/user/<username>/<mod-relpath>/<version>.zip
Prefer `moon ide` if the MoonBit toolchain is available.
```

这证明两件事：

第一，SOP 不能再把 “打开 docs 页面” 当成唯一文档核验路径。对 CLI / agent 来说，正确入口是 assets JSON 或 `moon ide doc`。

第二，直接爬 docs HTML 是错误对象。它不是文档内容，只是页面壳。原 SOP 没有明确这件事，会让后续自动化工具继续走错路。

### 3. manifest 可拿版本和 metadata

我请求：

```bash
curl -sSL 'https://mooncakes.io/api/v0/manifest/moonbitlang/core'
```

拿到当前 core 版本：

```text
0.1.20260609+84519ca0a
```

并且 manifest 里有：

```json
{
  "name": "moonbitlang/core",
  "module": "moonbitlang/core",
  "version": "0.1.20260609+84519ca0a",
  "latest_version": "0.1.20260609+84519ca0a",
  "versions_count": 153,
  "metadata": {
    "repository": "https://github.com/moonbitlang/core",
    "license": "Apache-2.0",
    "keywords": ["core", "standard library"],
    "checksum": "...",
    "created_at": "2026-06-09T05:26:09.780424+00:00"
  }
}
```

这说明 SOP 里 manifest 仍然是单模块详情的核心入口；这一点没错。但是 manifest 不是 API 结构的终点，它只是进入 module index / package data 的版本定位入口。

### 4. module_index.json 可用，而且非常强

我请求：

```bash
curl -sSL 'https://mooncakes.io/assets/moonbitlang/core@0.1.20260609+84519ca0a/module_index.json'
```

返回 200，大小约 443,995 bytes。

这个 JSON 是树结构，不是扁平列表。它能列出 module 内 package。对 `moonbitlang/core`，我解析后得到：

```text
module_index_packages = 61
first paths =
moonbitlang/core/abort
moonbitlang/core/argparse
moonbitlang/core/array
moonbitlang/core/bench
moonbitlang/core/bigint
moonbitlang/core/bool
...
```

并且 `moonbitlang/core/argparse` 的 package index 中已经包含：

- types：`ArgGroup`、`Command`、`FlagAction`、`FlagArg`、`Matches`、`OptionAction`、`OptionArg`、`PositionArg`、`ValueRange`、`ValueSource`
- methods：例如 `Command::new`、`Command::parse`、`Command::render_help`

这个事实非常重要。它说明官方已经提供了 module 内 package/API 的结构证据。SOP 原来只让人看 docs 页面、manifest、GitHub，会漏掉一个比页面更适合 CLI 消费的结构化索引。

这个入口直接改变消费者侧 CLI 的数据源模型：CLI 不应该只展示 module name/version/description/license/downloads，它可以在已知候选下钻时显示 module 内 packages、types、values、methods，让用户判断“这个候选到底暴露了什么 API”。

### 5. package_data.json 可用，能拿类型、函数、docstring、源码位置

我请求：

```bash
curl -sSL 'https://mooncakes.io/assets/moonbitlang/core@0.1.20260609+84519ca0a/argparse/package_data.json'
```

返回 200，大小约 18,579 bytes。

它的顶层 keys 是：

```text
errors
misc
name
traits
typealias
types
values
```

`types[0]` 中有：

- `name`: `ArgGroup`
- `docstring`: `Declarative argument group constructor.`
- `signature`: 带 HTML link 的签名
- `loc`: `path/file/line/column`
- `methods`: `new`，并且 method 也有 docstring、signature、loc

这说明 package_data 不是只有符号名，它已经包含足够强的 API 证据。对于消费者侧包发现工具，`view` 或 `inspect` 命令可以直接消费它，生成 pretty text 或 JSON 输出。

原 SOP 没有这一层，是明显盲区。

### 6. resource.json 不能机械相信

我按注释尝试：

```bash
https://mooncakes.io/assets/moonbitlang/core@0.1.20260609+84519ca0a/argparse/resource.json
```

结果返回 404。

但同一个 package 的 `package_data.json` 返回 200，`moon ide doc` 也能查 API，本地 `~/.moon/lib/core/argparse/README.mbt.md` 也存在。这说明：

```text
resource.json 404 != package 不存在
resource.json 404 != 没有 API 文档
resource.json 404 != 文档资产整体不可用
```

因此 `001-raw.md` 的大方向是对的：官方有结构化文档资产；但它里面“资源路径”那种表述不能无条件写进 SOP。SOP 应该写成：`resource.json` 是可尝试资源入口，404 时继续查 `package_data.json`、module index、Mooncakes docs、`moon ide doc`、发布源码。

我已经这样修了。

### 7. `moon ide doc` 实测可查 core package、类型、方法

我执行：

```bash
~/.moon/bin/moon ide doc '@moonbitlang/core/argparse'
```

它输出了 `argparse` package 的所有导出类型和方法，包括：

```text
package "moonbitlang/core/argparse"
struct ArgGroup
struct Command
  pub fn Command::new(...)
  pub fn Command::parse(...)
  pub fn Command::render_help(Self) -> String
...
```

我继续执行：

```bash
~/.moon/bin/moon ide doc '@moonbitlang/core/argparse.Command'
```

它输出了 `Command` 的说明和方法列表。

再执行：

```bash
~/.moon/bin/moon ide doc '@moonbitlang/core/argparse.Command::parse'
```

它输出了具体方法文档：

```text
Parse argv/environment according to this command spec.
Behavior:
- Help/version requests print output immediately and terminate with exit code `0`.
- Parse failures raise display-ready error text with full contextual help.
- Command-definition validation failures raise display-ready validation text.
Value precedence is `argv > env > default_values`.
```

这说明 `moon ide doc` 对“已知 package / 已知 symbol 的 API 下钻”已经非常有价值。

但是我也测试了：

```bash
~/.moon/bin/moon ide doc 'markdown'
```

结果是：

```text
No results found for query: 'markdown'
```

而 registry 里当然有 `mizchi/markdown` 等包。因此不能把 `moon ide doc` 解释成通用 package search。它不能替代 `moon search`，也不能替代 Modules API 本地过滤。它的对象是 API / symbol documentation，不是 registry discovery。

SOP 已经按这个边界修订：`moon ide doc` 是已知候选下钻入口，不是通用搜索入口；命中失败不证明 registry 没包。

### 8. 本地 core 源码已经存在于 `~/.moon/lib/core`

我检查了：

```text
~/.moon/lib/core/argparse/README.mbt.md
~/.moon/lib/core/argparse/command.mbt
~/.moon/lib/core/argparse/parser.mbt
~/.moon/lib/core/argparse/pkg.generated.mbti
...
```

`README.mbt.md` 中有 argparse 的 Quick Start、示例和使用说明。`command.mbt` 中有源码、doc comments、方法实现。

这说明对于 core / 标准库，发布版源码已经在本地工具链目录内，不需要 GitHub 才能阅读实现。

这对 SOP 的职责划分也很关键：GitHub 不应该继续被写成“源码核验的唯一入口”。至少对已发布版本，Moon 官方包体系、本地缓存、`moon fetch` 和 source zip 都是更接近发布事实的入口。

### 9. source zip 对普通包可用，但不能无条件写成全模块必通

我测试 `mizchi/markdown`：

```bash
curl -L --fail -o tmp/moonbit-doc-assets-check/markdown.zip \
  'https://download.mooncakes.io/user/mizchi/markdown/0.6.2.zip'
```

成功，zip 大小约 2.8 MB，里面有：

```text
moon.mod.json
README.md
LICENSE
src/*.mbt
src/api/moon.pkg
src/toc/moon.pkg
src/html_tests/*.mbt
benches/*.mbt
...
```

这证明官方 source zip 对已发布包确实可用，能拿到完整发布源码。

但我测试 `moonbitlang/core`：

```bash
curl -L --fail -o tmp/moonbit-doc-assets-check/core.zip \
  'https://download.mooncakes.io/user/moonbitlang/core/0.1.20260609+84519ca0a.zip'
```

返回 404。core 的源码在本地 `~/.moon/lib/core` 可见，但这个 source zip URL 对 core 最新版本失败。

所以 SOP 不能写成“所有模块都一定可以用同一个 source zip URL 下载”。准确写法是：source zip 是可用入口，需要跟随重定向；若失败，记录失败状态并回退到 `moon fetch`、本地 registry cache、`~/.moon/lib` 或 GitHub。

我已经按这个边界修订。

### 10. `moon fetch` 实测成功，是更适合工具封装的源码入口

我创建临时 probe 项目并执行：

```bash
~/.moon/bin/moon fetch mizchi/markdown@0.6.2
```

成功输出：

```text
Fetching mizchi/markdown@0.6.2 to .../.repos/mizchi/markdown/0.6.2
Downloading mizchi/markdown@0.6.2
Success: Successfully fetched mizchi/markdown@0.6.2 to .../.repos/mizchi/markdown/0.6.2
```

`.repos/mizchi/markdown/0.6.2` 中有完整源码、README、LICENSE、moon.mod.json、tests、examples 等。

这个实验把 `001-raw.md` 里的“官方工具链 + Mooncakes 已经可拿源码/文档”的判断落到了命令层。对消费者侧 CLI 工具来说，`moon fetch` 是一个非常重要的适配入口。它不需要用户知道 source zip URL，也不需要用户去 GitHub clone。它直接以 registry module coordinate 为输入，获取发布版源码。

但 `moon fetch --help` 标明它是 unstable command，所以 SOP 不能把它写成长期稳定契约。正确定位是：适合当前源码核验自动化封装；不是依赖采纳；不是稳定公共 API。

我已经按这个定位写入 SOP。

### 11. Skills API 已经存在，当前 70 条

我请求：

```bash
curl -sSL 'https://mooncakes.io/api/v0/skills'
```

返回数组长度 70。第一条大概是：

```json
{
  "module": "jaredzhou/pony",
  "author": "jaredzhou",
  "version": "0.1.1",
  "package": "cmd/main",
  "name": "jaredzhou/pony/cmd/main",
  "detail_url": "/skills/jaredzhou/pony@0.1.1/cmd/main",
  "wasm_url": "/assets/jaredzhou/pony@0.1.1/cmd/main/main.wasm",
  "checksum_url": "/assets/jaredzhou/pony@0.1.1/cmd/main/main.wasm.sha256",
  "repository": "github.com/jaredzhou/moonbase",
  "created_at": "2026-06-21T03:05:54.897810"
}
```

同时 `moon runwasm --help` 明确了 coordinate 规则：

```text
moon runwasm moonbitlang/parser/cmd/moonfmt@0.3.3
moon runwasm moonbitlang/parser@0.3.3/cmd/moonfmt
moon runwasm moonbitlang/parser/cmd/moonfmt
```

并说明未 pin coordinate 会解析 latest，wasm 缓存在：

```text
$MOON_HOME/registry/cache/assets
```

这说明 skills/runwasm 是一个真实数据面。它不能替代 library package discovery，但消费者侧 CLI 如果目标是完整适配当前官方资源，就必须分出 `kind=skill` / `kind=executable`，不能只搜索 modules。

原 SOP 完全没覆盖 Skills API，这是盲区。我已经加入。

### 12. `/api/v0/modules?search=` 仍不可用作服务端搜索

我再次验证：

```bash
curl -sSL 'https://mooncakes.io/api/v0/modules?search=markdown'
```

结果：

```json
{
  "type": "array",
  "length": 1350,
  "first": "0Ayachi0/elk",
  "last": "zxch3n/expect"
}
```

普通 `/api/v0/modules` 也是：

```json
{
  "type": "array",
  "length": 1350,
  "first": "0Ayachi0/elk",
  "last": "zxch3n/expect"
}
```

所以 SOP 里“`/api/v0/modules?search=<query>` 当前不能作为服务端过滤使用，应拉全量 JSON 后本地过滤”仍然正确，没有过期。

这也是消费者侧 CLI 仍有必要存在的核心事实之一。

### 13. 本地 registry index 与 statistics 数字仍有口径差异

执行 `moon update` 后：

```text
Registry index updated successfully
Symbols updated successfully
```

本地 index：

```text
b09b7cdb update cybershang/agent-telemetry
find ~/.moon/registry/index/user -name '*.index' | wc -l => 1363
```

Statistics API：

```json
{
  "total_modules": 1350,
  "total_packages": 12008,
  "total_lines": 44597511,
  "total_downloads": 4043104
}
```

这个口径差异仍然存在。SOP 原来已经写了“API module count、本地 index 文件数、statistics 数字可能不同；记录调查日期和来源”，这点继续成立。

### 14. 本地 symbols cache 有价值，但不是全量

我查看：

```bash
find ~/.moon/registry/symbols -type f
```

当前只有：

```text
~/.moon/registry/symbols/moonbitlang/async/0.19.1.symbols
~/.moon/registry/symbols/moonbitlang/quickcheck/0.14.0.symbols
~/.moon/registry/symbols/moonbitlang/x/0.4.45.symbols
```

这些文件是 JSON Lines，包含：

```json
{"kind":"meta","module":"moonbitlang/x","version":"0.4.45","schema_version":2}
{"key":"moonbitlang/x/crypto.sha256_from_iter","pkg":"moonbitlang/x/crypto","name":"sha256_from_iter","kind":"function","sig_":"pub fn sha256_from_iter(Iter[Byte]) -> FixedArray[Byte]","doc":null}
```

这说明 symbols cache 是真实存在的 API / symbol 发现源，但当前不是全 registry 覆盖。SOP 应把它写成补充入口，不是通用搜索主入口。我已经这样写。

## 三、哪些说法被确认，哪些说法被证伪或修正

### 被确认的说法

1. Mooncakes docs 页面是 SPA，直接 curl docs URL 不能拿正文。
2. 官方页面注释确实给了 API/assets 入口。
3. `module_index.json` 真实可用，可拿 module 内 package/API 结构。
4. `package_data.json` 真实可用，可拿单 package 类型、函数、方法、docstring、源码位置。
5. `moon ide doc` 真实可用，可查本地可见 API 文档。
6. `moon fetch` 真实可用，能把发布版源码下载到 `.repos`。
7. `mizchi/markdown` 的 source zip 真实可下载，包含发布版完整源码。
8. Skills API 真实存在，当前有可执行 Wasm / skill entry 数据。
9. `/api/v0/modules?search=` 仍不执行服务端搜索。
10. `moon` CLI 仍没有通用 `search/view/audit/outdated`。

### 被修正的说法

1. “官方文档可以用资源路径拿到”必须改成：`resource.json` 是可能入口，不保证每个 package 都有；404 时不能推出 package 不存在。
2. “官方 source zip 可直接拿所有源码”必须改成：source zip 是可用入口，但需要跟随重定向；不同模块可能失败；失败时回退到 `moon fetch`、本地缓存、`~/.moon/lib` 或 GitHub。
3. “`moon ide doc` 是命令行查文档的完整方案”必须改成：它是已知 API / symbol 下钻方案，不是 registry package discovery。
4. “GitHub 是源码核验入口”必须改成：GitHub 是维护、协作、CI、issue、release、未发布源码核验入口；发布版源码优先看 Mooncakes assets / `moon fetch` / source zip / registry cache。
5. “当前 CLI 搜索只能 curl/jq 拼接”必须细化：通用 package discovery 仍要 curl/jq 或 Web；但已知候选的 API/文档/源码下钻已有官方 CLI/assets 入口。

## 四、SOP 原来的真实盲区

SOP 原来不是完全错。它在“通用搜索”层是对的：当前还是要拉 modules 全量 JSON，本地过滤，manifest 下钻，本地 probe 验证。

但它在“候选下钻”层明显不完整。具体盲区如下。

### 盲区一：把 docs 页面当主要文档入口，没有纳入官方 assets

原 SOP 说打开 docs 页面、读取 README/API，但没有说 SPA 页面背后的 `module_index.json` 和 `package_data.json`。这会导致 CLI 工具继续模拟人打开网页，而不是直接消费结构化 API。

修订后加入：

- Module index asset
- Package data asset
- Package resources asset
- 对 404 的解释规则
- package path 必须从 module index 或 docs 路径确认，不能凭 URL 猜

### 盲区二：没有把 `moon ide doc` 作为已知候选 API 下钻入口

原 SOP 没有这个入口。修订后加入 `moon ide doc`，并明确它不是通用 registry 搜索。

### 盲区三：没有把 `moon fetch` 作为发布源码获取入口

原 SOP 的源码核验基本指向 GitHub。实测 `moon fetch mizchi/markdown@0.6.2` 成功。这说明官方 CLI 已经能按 registry coordinate 获取发布版源码。

修订后加入 `moon fetch`，并说明它 unstable、适合源码核验、不替代依赖采纳。

### 盲区四：没有区分发布版源码与 GitHub 开发协作信息

GitHub 上的源码可能是 main 分支，不一定等于 registry 发布版本。消费者判断“我将要引入的依赖是什么”时，发布版源码更接近对象事实。GitHub 更适合看维护状态、未发布变化、CI、issues、PR、release。

修订后把 GitHub 的职责降回正确位置。

### 盲区五：没有纳入 Skills API / runwasm entry

当前官方已有 skills API 和 `moon runwasm`，这属于 execution / skill loop，不是 library dependency loop。但消费者侧工具如果目标是“发现可用 MoonBit 包生态对象”，不能忽略它。修订后加入 Skills API，并强调 skill 不能和 library 采纳结论混写。

### 盲区六：没有纳入 symbols cache

`moon update` 会更新 symbols cache，当前本地有少数模块的 `.symbols` JSONL。它不是全量搜索，但可以作为 API/symbol 发现补充入口。修订后加入。

### 盲区七：候选记录表缺少 API asset、源码获取、skill 相关字段

如果候选记录表不记录 `module_index`、`package_data`、`source_fetch`、`skill_entry`、`wasm_checksum`，后续报告仍会回到旧 SOP 的盲区。修订后已补入候选表。

## 五、我已经改了什么

修改文件：

```text
notes/关于MoonBit/关于包搜索与发现/MoonBit包搜索发现与生态查重SOP.md
```

本次 diff：

```text
1 file changed, 134 insertions(+), 12 deletions(-)
```

核心修改包括：

1. 候选记录表新增：

```text
module_index
package_data
source_fetch
skill_entry
wasm_checksum
```

2. “当前可用入口”新增：

```text
Skills API
Module index asset
Package data asset
Package resources asset
Source zip
moon ide doc
moon fetch
本地 symbols cache
```

3. 搜索/核验路径改为：

```text
搜索候选：Web、Modules API、本地 index、可用 symbols cache
核验候选：Manifest API、module index、package data、docs 页面、moon ide doc、moon fetch、source zip
GitHub：维护、协作、CI、issue、release、未发布源码
```

4. 快速路径新增 manifest + module index / package data / `moon ide doc` 下钻。

5. 可靠查重路径新增：

```text
读取 module index
读取 package data 或 moon ide doc
使用 moon fetch / source zip 获取发布版源码
读取 Skills API 和 runwasm coordinate
```

6. 查询命令新增五个小节：

```text
5.6 查询 module index 与 package data
5.7 使用 moon ide doc 查询本地 API 文档
5.8 获取发布版源码
5.9 查询 Skills API
5.10 搜索本地 symbols cache
```

7. 当前 Mooncakes 注意事项新增：

```text
Mooncakes docs 是 SPA，curl docs 只拿 HTML 外壳
package_data.json / module_index.json 是重要下钻入口
resource.json 可能不存在，404 不等于 package 不存在
moon ide doc 不是通用 registry 包搜索
moon fetch 可获取发布源码到 .repos
Skills API / runwasm 是执行对象，不和 library 候选混写
```

## 六、对消费者侧 CLI 对象的影响

这次验证把消费者侧 CLI 的数据源边界重新拉清楚了。

原来的消费者侧 CLI 模型如果只包含：

```text
modules API
manifest API
local registry index
GitHub
本地 probe
```

那是不完整的。

现在它至少应该包含：

```text
Modules API：生成 module 候选集合
Manifest API：单 module 元数据、版本、downloads、build_status、metadata
Module index asset：module 内 package/API 结构索引
Package data asset：单 package API、docstring、源码位置
Resources asset：可用时读取 README 等资源
moon ide doc：本地可见 API 文档下钻
moon fetch：发布版源码获取
Source zip：发布版源码下载补充入口
Skills API：可执行 Wasm / skill 候选源
runwasm help / coordinate：执行对象公共契约
Local symbols cache：可用时的符号搜索源
Local registry index：原始发布记录和交叉验证
GitHub：维护、协作、CI、issue、release、未发布源码
Local probe：add/check/target build/test 验证
```

这不意味着消费者侧 CLI 变成“全生态治理工具”。边界仍然是消费者侧发现适配器。它不生产上游 metadata，不修复作者声明，不替官方做 registry search API，不给安全背书。

但在当前对象边界内，它必须完整消费现有事实源。如果它漏掉 module_index、package_data、moon ide doc、moon fetch、skills API，就仍然会逼用户手工拼接这些官方已经给出来的入口。

更准确的 CLI 行动轨迹现在应是：

```text
search
  从 modules API / local index / symbols cache / skills API 生成候选

view
  用 manifest + module_index + package_data + docs URL + repository 展示证据画像

api
  用 package_data 或 moon ide doc 下钻类型、函数、方法、docstring

source
  用 moon fetch / source zip / local cache 获取发布源码

skill
  用 skills API 展示 wasm_url、checksum_url、runwasm coordinate

compare
  把多个候选的 metadata、API 结构、源码获取、docs、license、target、维护信号放在同一表面

probe
  回到本地项目做 moon add / moon check / target build

record
  记录采用、适配后采用、继续验证、贡献上游、fork、自建、暂不采用
```

这才是“语义化 CLI 工作流”应该覆盖的对象边界。

## 七、最终判断

SOP 不是整体过期，而是局部盲区严重。具体说：

- **搜索层没过期**：官方仍没有通用 `moon search`，`/api/v0/modules?search=` 仍不服务端过滤，Web 搜索仍无稳定深链。
- **下钻层已过期**：官方已经把文档/API/源码/skill 的结构化入口暴露出来，SOP 继续只写 Web/docs/GitHub/curl modules 会漏掉重要事实源。
- **`001-raw.md` 的警告是有效的**：Mooncakes docs HTML 里的 Agent note 确实揭示了官方给工具和 agent 准备的 API/assets 路径。旧 SOP 没吃进去，等于在候选核验层视而不见。
- **`001-raw.md` 自身也需要修正**：不能无条件相信资源路径，要改看 `resource.json`；不能把 source zip 写成所有模块必通，不能把 `moon ide doc` 当 package search。
- **GitHub 的位置要降级并精确化**：GitHub 是维护和协作事实源，不是发布版源码唯一事实源。
- **消费者侧 CLI 的当前对象更明确了**：它不是凭空造生态数据，而是把 Mooncakes/API/assets/工具链/本地缓存/GitHub/probe 这些已有事实源组织成搜索、下钻、比较、源码核验、skill 识别、本地验证、记录结论的命令环境。

我已经把 SOP 按这个结论改掉了。
</Assistant>

<User id="1">
好的啊，那么再把这个愿望 再次捋一捋啊，也就是希望要有一款工具。一款开发工具，一款CLI工具。工具本身的设计就代表着环境。一个丝滑的工作流，工具本身就可以代表着一套丝滑的工作流，那么什么工作流呢？是不是就是Package Discover？ 那这里是一部分的啊。那么需要明确100%当下所有的数据来源，也就是所有的输入。那我们的这个CLI是做什么？做工作流，做聚合输出。做格式的整理，做Pretty text ，--json方便Code Agent与人查看，也方便他们做管道机械处理。在这里，我希望啊，Json的这里还希望加一个模式，能够去看一下这个接口的Schema。我觉得这个在JQ这里还是非常常用的一个功能，但是我目前可能是我信息闭塞吧，我没有看到什么比较成熟的工具会去做这个东西。还是说只是我没有见过吗？也就是现代CLI工具，对吧普遍都是默认pretty text的，然后呢所有命令那都有JSON，都能拿到机械格式。但是呢，拿到机械格式下一步呢？是不是？通常是搭配ZQ。那使用者还要得去看一下，还尝试真实的请求一次才能拿到Schema吗？是不是体验就是不好？是吧？这里有个小摩擦啊。我是。嗯，通常都有看到，我今天把它说出来一下，我也不知道你觉得是不是真摩擦？你觉得这是摩擦吗？是吧？可以纳入设计之中吗？
嗯，那么总之呢需要隐藏复杂性。也就是把业务暴露给使用者。使用者带着他自己的意图、认知活动，用我们这款工具。我们这款工具呢，哎，它可以有一份入口指南，然后呢，哎，每一个命令呢，都可以有详细的help。开发者呢，是吧？不管是开发者还是CodeAgent都可以快速的上手。不要去了解太多没有必要的信息，占据注意力，占据上下文，浪费token，是不是？
就像之前描述的，做这件事情还要看好几百行的SOP才能操作，这本身就是一个非常奇怪的、扭曲的一些平替的出现，是不是？

那这里有二八法则对不对？把常用的，诶，这边直接固化成清晰的命令、子命令，它就是环境里的东西。使用者都能看懂这些命令，它就是它工作流的一部分，是高频的，经常会用到的。那剩下的那些呢？诶，它们可以用JQ自己去补上，按需索要，是吧？ 那么这里涉及到了DX开发体验、交互逻辑、元因工程的一些东西，你可以去调动相关的知识或者参考来思考一下。

然后呢之前的讨论里面有提到什么未来的字段，我们只是标准化接口、标准化架构、标准化这个管径，没有的字段不是显unknown，是直接就是什么呢？就直接不显示，因为没有这个东西。啊，如果说是现有的字段以及进入设计字段，但是这个字段值本身是没有的，那我们才会显示unknown，不是什么虚空一个未来的本体字段，这个东西有没有分清楚？

好的啊，请你自己思考这些问题，想一想啊，然后呢来跟我讨论，把你的一整套设计告诉我，你自己不要留有不确定性。如果有，你自己想清楚，想明白再来找我，和我沟通不少于6000字或者更多。从头到尾，从头发丝到脚趾头，全部给我讲明白。然后呢，是吧？在这里，明确明确一定要明确，我们这一个聚合工具是吧？我们这个Adapter，我们这个开发工具，它的数据来源于哪里？不要讲可能大概也许，直接给我明确出来，精确的啊。然后呢，诶，这个人因工程交互逻辑怎么设计？二八法则如何分配？ 等等啊。然后呢，我希望你最后还能给我来一个完整的用户旅程、用户故事。不是随便挑几个基本流，也不是随便挑几个分支，就是完整的。啊，你可以。借用一下离散事件仿真的思想，去定义一下事件，定义一下时序。比如说一个开发者是吧，从生命周期的0开始启动，触发事件，比如说他比如说做什么工作呀，是吧？然后呢，在这里他运行什么命令啊，看到什么东西啊，就把这种期望的这种状态直接描绘出来，高仿真的仿真出来，直接把这这个开发者打了什么东西啊？输出了什么东西啊？这种描绘出来，直接写出来。啊？就好像我能直接看到我们一起来看它的终端执行记录一样，我们好像也能看到开发者手打的命令一样，哦，就啊仿真的故事啊。那讲故事的过程中，你自己就可以去推导、反向推理、论证、推导合理性，是否合理。同时呢，故事呢，也方便我看，因为这个故事呢，只是现在是故事，是吧？那下一秒可能就变现实了呀，是不是？你懂我意思吗？只要故事，我们从逻辑上去推理、推导、自洽，诶，那其实它就可以作为一种变相的黑盒验收测试，是不是？好好想一想吧，想好之后呢，准备好之后啊，用弹窗跟我沟通，开始吧。

Tip:
完成后记得使用弹窗与我沟通。
</User>

<Assistant id="2">
# 消费者侧 MoonBit 包发现 CLI 设计

## 一、先回答 schema/shape 是不是真摩擦

这是一个真摩擦，可以进入设计，而且它不是边角功能。

原因很直接：现代 CLI 常见模式是默认给人看 pretty text，再用 `--json` 给机器看；但用户一旦拿到 JSON，下一步就会问“这个 JSON 到底有哪些字段、哪些字段稳定、哪些字段可能为空、哪些字段是数组、哪些字段来自哪个来源”。如果工具只说“你可以 `--json | jq`”，那只是把输出格式交给了用户，没有把机器消费的公共契约交给用户。用户仍然要先跑一次真实请求，观察样例，再试 jq，再被字段缺失、嵌套结构、数组路径、null 值、命名不稳定折腾。

这不是你信息闭塞。成熟工具确实各自做了一部分，但没有形成每个 CLI 都自带“输出 schema 浏览”的普遍纪律。

GitHub CLI 是一个参考点。它的官方手册写得很清楚：`--json` 需要字段名；如果省略字段参数，CLI 会显示可用 JSON 字段；并且提供 `--jq` 和 `--template` 做格式化。这个设计解决了一部分问题：用户不用猜字段名，也不用额外安装 jq 才能做基础过滤。但它主要暴露的是可选字段列表，不是完整 JSON Schema，也不告诉你字段来源、缺失状态、稳定性等级。来源：GitHub CLI formatting 文档：https://cli.github.com/manual/gh_help_formatting

Kubernetes 是另一个参考点。它的 API 有 OpenAPI schema，`kubectl explain` 可以让用户在命令行里查看资源字段文档；Kubernetes 文档也明确 OpenAPI v3 是访问 API schema 的首选方式。这个方向证明“命令行里看对象 schema / field documentation”本身是成熟需求，不是多余设计。来源：Kubernetes API 文档：https://kubernetes.io/docs/concepts/overview/kubernetes-api/

OpenAPI 本身的定义也说明了同一件事：一个标准接口让人和机器在不访问源码、额外文档或抓包的情况下理解服务能力。来源：OpenAPI Specification：https://swagger.io/specification/

所以这里的结论是：**schema/shape 是消费者侧 CLI 的公共契约投影，不是锦上添花。**

但我们要把它设计准。它不是“显示未来可能有的字段”。它只说明当前工具已经承诺的输出对象。未来扩展点存在于内部架构和版本演进规则中，不出现在当前 pretty text，也不出现在当前 JSON schema 的 required/visible 字段里。

## 二、当前人工制品定义

这个工具的对象定义是：

**一个本地 MoonBit 包发现适配器。它从 Mooncakes 官方 API、Mooncakes assets、本地 Moon registry、Moon 工具链、本地项目上下文、GitHub 维护信号和用户调查记录中读取已存在事实，把 MoonBit 包发现、候选下钻、API 阅读、源码核验、skill 识别、本地验证和采纳记录组织成语义化 CLI 工作流，并向人和 Code Agent 输出 pretty text、JSON、schema/shape 和报告。**

这个定义里的每个词都有责任：

- “本地”：工具在用户机器上运行，能读当前项目、Moon 缓存、registry index、symbols cache、`.repos`、`moon` 命令输出。
- “MoonBit 包发现”：对象是消费者侧 Package Discovery，不是作者侧发布健康检查。
- “适配器”：工具不替官方 registry 发明事实，不替作者补 metadata，不替安全系统宣布无漏洞。它消费已有事实，整理并投影。
- “语义化 CLI 工作流”：用户不再手动拼 curl、jq、GitHub、docs URL、probe 项目；用户沿命令走。
- “pretty text、JSON、schema/shape”：同一内部事实面向不同消费者有不同投影。

这款工具不是：

- 不是 Mooncakes 替代站。
- 不是官方 registry。
- 不是包作者健康门禁。
- 不是漏洞审计系统。
- 不是通用搜索引擎。
- 不是只包一层 curl 的脚本合集。

它的边界内必须完整解决消费者侧发现工作流；边界外只读取事实或标明不可用，不冒充。

## 三、明确数据来源：当前 100% 输入清单

下面是当前工具必须认识的数据来源。这里不写“可能”，只写当前已确认的输入类型。每个输入进入工具时都要有 source label、fetch time、failure state。

### 1. Mooncakes Modules API

```text
GET https://mooncakes.io/api/v0/modules
```

用途：生成 module 候选集合。

当前事实：`/api/v0/modules?search=<query>` 不执行服务端过滤，所以工具必须拉全量 modules 后本地过滤。

消费字段：

```text
name
version
description
keywords
repository
license
```

这个数据源服务 `search`。

### 2. Mooncakes Statistics API

```text
GET https://mooncakes.io/api/v0/modules/statistics
```

用途：记录搜索快照的 registry 规模，生成报告时说明调查基准。

消费字段：

```text
total_modules
total_packages
total_lines
total_downloads
```

这个数据源服务 `sync`、`report`、`doctor`。

### 3. Mooncakes Manifest API

```text
GET https://mooncakes.io/api/v0/manifest/<owner>/<module>
```

用途：单 module 下钻。

消费字段：

```text
name
module
version
latest_version
downloads
has_package
build_status
metadata.repository
metadata.license
metadata.keywords
metadata.description
metadata.checksum
metadata.targets
metadata.supported-targets
versions
versions_count
```

这个数据源服务 `view`、`compare`、`probe` 前置判断。

### 4. Mooncakes Module Index Asset

```text
GET https://mooncakes.io/assets/<owner>/<module>@<version>/module_index.json
```

用途：读取 module 内 package/API 结构索引。

消费对象：

```text
package path
types
traits
errors
typealias
values
methods
misc
child package tree
```

这个数据源服务 `view`、`api`、`compare`。

关键解释：这是树结构，不是扁平列表。工具内部必须解析树，生成 package path 列表和 API 摘要。

### 5. Mooncakes Package Data Asset

```text
GET https://mooncakes.io/assets/<owner>/<module>@<version>/<pkg-relpath>/package_data.json
```

用途：读取单 package 的 API 详情。

消费对象：

```text
name
types[].name
types[].docstring
types[].signature
types[].loc
types[].methods[].name
types[].methods[].docstring
types[].methods[].signature
values[].name
values[].docstring
values[].signature
traits
errors
typealias
misc
```

这个数据源服务 `api`、`view --api`、`compare --api`。

### 6. Mooncakes Package Resources Asset

```text
GET https://mooncakes.io/assets/<owner>/<module>@<version>/<pkg-relpath>/resource.json
```

用途：读取 README 等资源。

失败语义：404 只表示这个资源文件不可用，不表示 package 不存在，不表示没有 API 文档。

这个数据源服务 `view --readme`、`report`。

### 7. Mooncakes Source Zip

```text
GET https://download.mooncakes.io/user/<owner>/<module>/<version>.zip
```

用途：下载 registry 发布版源码。

消费对象：

```text
moon.mod.json / moon.mod
moon.pkg.json / moon.pkg
README
LICENSE
*.mbt
examples
tests
bench files
CI/config files included in package
```

执行要求：必须跟随重定向。

失败语义：下载失败记录为 source_zip_unavailable；工具回退到 `moon fetch`、本地 registry cache、`~/.moon/lib` 或 GitHub。

### 8. Mooncakes Skills API

```text
GET https://mooncakes.io/api/v0/skills
```

用途：发现可执行 Wasm / skill entry。

消费字段：

```text
module
author
version
package
name
detail_url
wasm_url
checksum_url
metadata.description
metadata.name
repository
created_at
```

这个数据源服务 `skill search`、`skill view`、`runwasm` 提示。

关键边界：skill 是执行对象，不是 library 依赖。它进入 `kind=skill`，不和 `kind=library` 混排成同一采纳结论。

### 9. Wasm Asset 与 Checksum Asset

来自 Skills API：

```text
wasm_url
checksum_url
```

用途：记录可执行 artifact、校验入口、runwasm 复现信息。

这个数据源服务 `skill view`、`runwasm`、`record`。

### 10. 本地 Moon Registry Index

```text
~/.moon/registry/index/user/**/*.index
```

用途：交叉验证原始发布记录，补充 API 结果，离线搜索。

格式：JSON Lines。

这个数据源服务 `search --local`、`view --index`、`doctor`。

### 11. 本地 Symbols Cache

```text
~/.moon/registry/symbols/**/*.symbols
```

用途：已缓存模块的 API/symbol 搜索。

格式：JSON Lines。

消费字段：

```text
module
version
schema_version
key
pkg
name
parent
kind
sig_
doc
attrs
```

这个数据源服务 `api search`、`search --symbol`。

边界：当前不是全 registry 覆盖，命中失败不证明 registry 没包。

### 12. 本地 Moon 工具链命令

工具调用本机 `moon`，但每个命令的职责不同。

```text
moon --version
moon update
moon ide doc <query>
moon fetch <module[@version]>
moon add <module>
moon check
moon test
moon build --target <target>
moon runwasm <coordinate>
moon tree
moon install
```

用途：

```text
版本记录
更新 registry / symbols
API 文档下钻
发布源码获取
依赖引入验证
target 验证
skill 运行提示或执行
依赖树观察
```

边界：工具封装命令，但不篡改 `moon` 的语义。

### 13. 本地 Moon 缓存

```text
~/.moon/registry/cache/**
~/.moon/registry/cache/assets/**
~/.moon/lib/core/**
```

用途：读取已经下载的包源码、wasm asset、core 标准库源码。

这个数据源服务 `source`、`api`、`doctor`、`offline`。

### 14. 当前 MoonBit 项目上下文

运行目录中的：

```text
moon.mod / moon.mod.json
moon.pkg / moon.pkg.json
当前 workspace
当前 deps
当前 target 需求
```

用途：给 `probe`、`add-check`、`contextual search` 提供上下文。

例如用户在一个目标 target 是 `js` 的项目里搜索 `markdown`，工具应默认把 `target=js` 作为可见上下文提示；但不能默认过滤掉 target unknown 包，除非用户要求硬过滤。

### 15. GitHub Repository

来源：manifest 或 modules API 中的 `repository` 字段。

用途：维护信号和协作信号。

消费字段：

```text
repo existence
archived
pushed_at
updated_at
licenseInfo
stars
forks
issues
PRs
releases
CI workflows
default branch
README
tests/examples not included in release package
unpublished code state
```

边界：GitHub 不是 registry 发布版源码的唯一事实源。它补充维护和未发布状态。

### 16. 用户本地调查记录

工具自己的记录目录，例如：

```text
.seekmoon/sessions/*.json
.seekmoon/records/*.json
.seekmoon/reports/*.md
```

用途：保存一次 search/view/compare/probe/record 的可复查轨迹。

消费对象：

```text
query
filters
snapshot time
candidate ids
selected candidates
probe results
conclusion
notes
```

这个数据源服务 `compare --from-search`、`report`、`record update`。

## 四、字段显示纪律：未来字段、unknown、missing 必须分清

这里必须严格。

### 1. 当前命令 schema 外的字段，不显示

如果某个字段只是未来可能接入，例如：

```text
advisory
provenance
reverse_dependents
quality_score
publisher_identity
```

而当前工具没有这个数据源、没有这个命令 schema、没有这个字段契约，那么 pretty text 不显示，JSON 也不输出这个字段。

它只存在于内部架构的扩展空间，不进入当前公共投影。

### 2. 当前命令 schema 内的字段，查了但没有值，才显示 missing / unknown

例如 `view mizchi/markdown` 的当前 schema 包含 `license`。如果 manifest 没给 license，JSON 输出：

```json
"license": {"status": "missing", "source": "manifest_api"}
```

pretty text 可以显示：

```text
license: missing
```

因为 license 是当前采用判断的硬字段。

再比如 `--target js` 时，target compatibility 是当前查询上下文字段。如果 manifest 和 index 都没有 target 信息，pretty text 应显示：

```text
target: unknown for js
```

因为这会影响用户下一步行动。

### 3. 当前命令没有请求的深层字段，不主动显示

用户执行：

```bash
seekmoon search markdown
```

默认搜索结果不应展开所有 API types、methods、source files。它只显示候选摘要。

用户执行：

```bash
seekmoon api mizchi/markdown
```

才显示 API 结构。

用户执行：

```bash
seekmoon source mizchi/markdown
```

才显示源码获取状态。

这就是二八原则：高频路径保持清楚，深层信息按意图下钻。

### 4. JSON 保留结构状态，pretty text 只显示行动相关状态

JSON 是机器投影，保留 status/source/fetched_at/error。

Pretty text 是人类投影，只显示影响当前动作的状态。

例如 source zip 失败但 `moon fetch` 成功：

JSON：

```json
"source": {
  "moon_fetch": {"status": "ok", "path": ".repos/mizchi/markdown/0.6.2"},
  "source_zip": {"status": "failed", "error": "404"}
}
```

Pretty text：

```text
source: fetched via moon fetch (.repos/mizchi/markdown/0.6.2)
```

只有用户执行 `--verbose` 或 `source --details` 时才显示 source zip 失败。

## 五、输出设计：pretty text、JSON、schema/shape

### 1. 默认 pretty text

默认输出服务人。它要短、稳定、可扫读。

`search` 默认输出：

```text
query: markdown       snapshot: 2026-06-21 22:10 CST       modules: 1350

#  module                         version   license   build    downloads   why
1  mizchi/markdown                 0.6.2     MIT       success  7567        name, keywords, description
2  moonbit-community/cmark         0.4.4     Apache?   success  ...         name, description
3  shiguri-01/markdown             0.1.1     ...       ...      ...         name

Next:
  seekmoon view 1
  seekmoon compare 1 2 3
  seekmoon probe 1 --target js
```

它不让用户看 API 地址，不让用户写 jq，不把所有字段摊开。

### 2. `--json`

`--json` 输出当前命令的规范对象，不是屏幕文本的 JSON 化。

```bash
seekmoon search markdown --json
```

输出：

```json
{
  "schema": "seekmoon.search-results.v1",
  "snapshot": {
    "fetched_at": "2026-06-21T22:10:00+08:00",
    "modules_count": 1350,
    "sources": ["modules_api", "statistics_api"]
  },
  "query": {
    "text": "markdown",
    "kind": "library",
    "target": null
  },
  "results": [
    {
      "rank": 1,
      "module": "mizchi/markdown",
      "version": "0.6.2",
      "description": "Incremental Markdown parser and compiler",
      "license": {"value": "MIT", "source": "modules_api"},
      "repository": {"value": "https://github.com/mizchi/markdown.mbt", "source": "modules_api"},
      "match": {
        "fields": ["name", "description", "keywords"],
        "reason": "query terms matched module name, description, and keywords"
      }
    }
  ]
}
```

JSON 输出包含 schema id，让 agent 和脚本知道自己处理的对象版本。

### 3. `--schema`

`--schema` 输出当前命令 JSON 的 JSON Schema，不发真实网络请求，不依赖样例数据。

```bash
seekmoon search --schema
seekmoon view --schema
seekmoon api --schema
seekmoon skill search --schema
seekmoon probe --schema
```

输出对象是工具承诺的机器契约。例如：

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "seekmoon.search-results.v1",
  "type": "object",
  "required": ["schema", "snapshot", "query", "results"],
  "properties": {
    "schema": {"const": "seekmoon.search-results.v1"},
    "snapshot": {"$ref": "#/$defs/snapshot"},
    "query": {"$ref": "#/$defs/searchQuery"},
    "results": {
      "type": "array",
      "items": {"$ref": "#/$defs/searchResult"}
    }
  }
}
```

这解决的是：用户不用先跑一次真实查询才能知道 jq 怎么写。

### 4. `--shape`

`--shape` 是给人看的 schema 摘要。它不是 JSON Schema，而是快速字段树。

```bash
seekmoon search --shape
```

输出：

```text
seekmoon.search-results.v1

schema: string
snapshot:
  fetched_at: datetime
  modules_count: int
  sources: string[]
query:
  text: string
  kind: library|skill|all
  target: string|null
results[]:
  rank: int
  module: string
  version: string
  description: string|null
  license:
    value: string|null
    source: string
    status: present|missing|unknown
  repository:
    value: string|null
    source: string
    status: present|missing|unknown
  match:
    fields: string[]
    reason: string
```

这个命令直接回应你说的 jq 前置摩擦。用户看一次 shape，就知道 `.results[].module`、`.results[].license.value` 怎么取。

### 5. `--fields`

`--fields` 给出当前命令可选择的字段列表，参考 GitHub CLI 的字段发现思想。

```bash
seekmoon search --fields
```

输出：

```text
rank
module
version
description
keywords
license
repository
build_status
downloads
match.fields
match.reason
```

配合：

```bash
seekmoon search markdown --json module,version,license,repository
```

这让 80% 的脚本用户不用处理完整 JSON。

### 6. `--jq`

工具可以内置 jq 风格过滤，也可以只输出 JSON 给系统 jq。这里的设计结论是：提供 `--jq`，因为它减少 agent 和人类的管道摩擦。

```bash
seekmoon search markdown --json --jq '.results[].module'
```

如果实现层不想内置 jq 解释器，也可以调用系统 jq；但公共契约上要明确：`--jq` 消费的是工具 JSON 输出，不是 pretty text。

## 六、命令设计：高频路径和低频路径

### 高频命令：二八中的 80%

这部分命令是普通开发者和 Code Agent 每天会用的环境动作。

#### 1. `sync`

同步数据快照。

```bash
seekmoon sync
```

职责：

```text
调用 moon update
拉取 modules API
拉取 statistics API
记录 snapshot
检查本地 registry index
检查 symbols cache
```

输出：

```text
Mooncakes snapshot updated
  modules API: 1350 modules
  statistics: 12008 packages, 4043104 downloads
  local index: 1363 .index files
  symbols cache: 3 modules
  moon: 0.1.20260608
```

#### 2. `search`

从意图生成候选集合。

```bash
seekmoon search markdown
seekmoon search json --target js
seekmoon search argparse --kind library
seekmoon search cowsay --kind skill
```

职责：

```text
读取 modules API / cache
按 name, description, keywords, repository 搜索
按 kind/target/license/build_status 过滤或标记
输出候选摘要
保存本次 search session
```

#### 3. `view`

查看单个候选证据画像。

```bash
seekmoon view mizchi/markdown
seekmoon view 1
```

职责：

```text
读取 manifest
读取 module_index
显示 package 列表
显示 docs URL
显示 repository
显示 license/build/download/version
显示缺失和风险
给出下一步命令
```

#### 4. `api`

下钻 API。

```bash
seekmoon api mizchi/markdown
seekmoon api moonbitlang/core/argparse.Command::parse
seekmoon api 1 --package src/api
```

职责：

```text
优先 package_data
可用时调用 moon ide doc
显示 types/functions/traits/docstrings/source loc
```

#### 5. `compare`

比较多个候选。

```bash
seekmoon compare 1 2 3
seekmoon compare mizchi/markdown moonbit-community/cmark shiguri-01/markdown
```

职责：

```text
横向展示 metadata
横向展示 API/package 结构摘要
横向展示 docs/source/license/target/build/download/versions/repo status
不替用户裁决，只降低表格劳动
```

#### 6. `source`

获取或定位发布版源码。

```bash
seekmoon source mizchi/markdown@0.6.2
seekmoon source 1
```

职责：

```text
优先 moon fetch
失败则尝试 source zip
再查 registry cache
core 包查 ~/.moon/lib/core
需要时显示 GitHub repository
```

#### 7. `probe`

本地验证候选。

```bash
seekmoon probe mizchi/markdown --target js
seekmoon probe 1 --target wasm-gc --keep
```

职责：

```text
创建 probe 项目
moon add
moon check
moon test
moon check/build --target
记录结果
```

#### 8. `record`

记录采纳结论。

```bash
seekmoon record 1 --conclusion continue --note 'needs GFM coverage check'
seekmoon record mizchi/markdown --conclusion adopt
```

固定结论枚举：

```text
adopt
adopt-with-adapter
continue-verification
contribute-upstream
fork
build-own
reject-for-now
```

中文 pretty text 可以显示：采用、适配后采用、继续验证、贡献上游、fork、自建、暂不采用。

#### 9. `report`

生成调查报告。

```bash
seekmoon report --format markdown
seekmoon report --format json
```

职责：把 search session、view、compare、probe、record 变成可复查报告。

#### 10. `schema` / `shape` / `fields`

输出机器契约。

```bash
seekmoon schema search
seekmoon shape search
seekmoon fields search
```

这不是低频调试命令，而是 agent 和脚本用户高频需要的入口。

### 低频命令：二八中的 20%

这些命令服务深层诊断，不挤占主工作流。

```bash
seekmoon raw modules
seekmoon raw manifest mizchi/markdown
seekmoon raw module-index mizchi/markdown
seekmoon raw package-data mizchi/markdown src/api
seekmoon cache list
seekmoon cache clear
seekmoon doctor
seekmoon sources mizchi/markdown
seekmoon explain statuses
```

它们让专家能看底层，但普通用户不需要先学这些。

## 七、人因工程与交互逻辑

### 1. 用户只面对业务意图，不面对底层地址

用户说：

```bash
seekmoon search markdown
```

不需要知道：

```text
/api/v0/modules
jq test("markdown|commonmark")
manifest path
module_index path
package_data path
```

工具内部知道这些。复杂性留在工具里，不占用户注意力。

### 2. 输出引导下一步动作

每个 pretty output 末尾给下一步命令，不讲教程。

例如 search 后：

```text
Next:
  seekmoon view 1
  seekmoon compare 1 2 3
  seekmoon probe 1 --target js
```

view 后：

```text
Next:
  seekmoon api 1
  seekmoon source 1
  seekmoon probe 1 --target js
```

probe 后：

```text
Next:
  seekmoon record 1 --conclusion continue-verification
  seekmoon report --format markdown
```

这就是环境设计：让用户自然往下走。

### 3. 编号引用降低复制成本

搜索结果编号可以贯穿当前 session。

```bash
seekmoon view 1
seekmoon compare 1 2 4
seekmoon probe 1
seekmoon record 1 --conclusion adopt
```

用户不需要重复复制 `mizchi/markdown`。

### 4. 缺失只在影响行动时打扰人

默认 search 不显示所有缺失字段。

但这些情况必须显示：

```text
license missing：影响长期采用
repository missing：影响维护核验
target unknown：当用户指定 --target 时影响采用
build failed：影响可验证性
source unavailable：当用户执行 source 时影响源码核验
checksum missing：当用户执行 skill/runwasm 时影响执行复现
```

### 5. 错误信息必须带恢复动作

坏错误：

```text
404
```

正确错误：

```text
package resources not found
  tried: /assets/mizchi/markdown@0.6.2/src/resource.json
  meaning: resource.json is unavailable; package may still have API data
  next: seekmoon api mizchi/markdown --package src
```

### 6. Agent 友好不等于人不友好

Agent 需要 JSON、schema、stable ids、sources、status。

人需要 pretty text、编号、下一步命令、少量风险提示。

两者使用同一内部对象，不互相污染。

## 八、完整离散事件用户旅程

下面用离散事件仿真的方式写完整故事。角色是一个普通 MoonBit 开发者，目标是在项目里寻找 Markdown parser。这个故事也是黑盒验收测试的雏形。

### 初始状态 S0

```text
Actor: developer
Project: /home/user/projects/doc-tool
Need: parse Markdown and render / inspect AST
Target: js
Knowledge: does not know Mooncakes API, does not know candidate package names
Tool state: seekmoon installed, no fresh snapshot
Moon state: moon 0.1.20260608 installed
```

### Event E1：进入项目

用户动作：

```bash
cd /home/user/projects/doc-tool
```

状态变化：

```text
current workspace detected
moon.mod.json readable
project target context: js inferred or manually supplied later
```

### Event E2：检查工具环境

用户动作：

```bash
seekmoon doctor
```

终端输出：

```text
SeekMoon doctor

moon:        0.1.20260608
moonc:       v0.10.0+e66899a54
registry:    present
symbols:     3 modules cached
network:     mooncakes.io reachable
project:     MoonBit module detected

Data sources:
  modules API        ok
  statistics API     ok
  manifest API       ok
  assets API         ok
  skills API         ok
  GitHub CLI         available

Next:
  seekmoon sync
  seekmoon search <query>
```

状态变化：

```text
environment_validated = true
```

### Event E3：同步快照

用户动作：

```bash
seekmoon sync
```

终端输出：

```text
Mooncakes snapshot updated

remote:
  modules:    1350
  packages:   12008
  downloads:  4043104

local:
  registry index: 1363 files
  symbols cache: 3 modules

snapshot: .seekmoon/snapshots/2026-06-21T22-35-10.json
```

状态变化：

```text
snapshot_id = 2026-06-21T22-35-10
modules cached
statistics cached
local index checked
```

### Event E4：第一次宽搜索

用户动作：

```bash
seekmoon search markdown --target js
```

终端输出：

```text
Search: markdown       target: js       kind: library       snapshot: 2026-06-21T22:35

#  module                         version   license   build    downloads   target       why
1  mizchi/markdown                 0.6.2     MIT       success  7567        unknown     name, description, keywords
2  moonbit-community/cmark         0.4.4     Apache-2  success  3120        unknown     name, description
3  shiguri-01/markdown             0.1.1     MIT       success  91          unknown     name, description
4  tonyfettes/tree_sitter_markdown 0.1.26    Apache-2  success  540         unknown     name
5  hustcer/html2md                 0.1.0     MIT       success  80          unknown     description

Notes:
  target is unknown for listed packages because no target matrix was found in current sources.
  target unknown is not treated as supported.

Next:
  seekmoon view 1
  seekmoon compare 1 2 3
  seekmoon probe 1 --target js
```

状态变化：

```text
candidate_set created
search session stores #1..#5
```

解释：这里没有显示未来 advisory、provenance、score。因为这些字段不属于当前 `search` schema。

### Event E5：查看第一候选

用户动作：

```bash
seekmoon view 1
```

终端输出：

```text
mizchi/markdown 0.6.2

Purpose:
  Incremental Markdown parser and compiler

Registry:
  downloads:       7567
  build_status:    success
  license:         MIT
  repository:      https://github.com/mizchi/markdown.mbt
  docs:            https://mooncakes.io/docs/mizchi/markdown

Packages:
  src
  src/api
  src/toc
  src/slug
  src/frontmatter
  src/highlight_json
  src/highlight_css
  src/highlight_html
  src/highlight_moonbit
  ... 12 more

API summary:
  module_index: ok
  package_data: not loaded in this view

Source:
  not fetched

Target:
  declared target: unknown
  local verification: not run

Next:
  seekmoon api 1 --package src/api
  seekmoon source 1
  seekmoon probe 1 --target js
```

状态变化：

```text
manifest fetched
module_index fetched
candidate #1 enriched
```

### Event E6：查看 API

用户动作：

```bash
seekmoon api 1 --package src/api
```

终端输出：

```text
API: mizchi/markdown src/api

Types:
  Document
  Node
  ParseOptions
  RenderOptions

Functions:
  parse(input : String) -> Document
  render(doc : Document) -> String
  to_json(doc : Document) -> Json

Docs:
  package_data: ok
  docstrings: 18 entries
  source locations: available

Next:
  seekmoon source 1
  seekmoon compare 1 2 3 --api-summary
```

状态变化：

```text
package_data fetched
API evidence added
```

如果 package path 写错，输出应是：

```text
Package not found in module index: src/apix

Known packages:
  src
  src/api
  src/toc
  src/slug

Next:
  seekmoon view 1 --packages
```

### Event E7：获取源码

用户动作：

```bash
seekmoon source 1
```

终端输出：

```text
Source: mizchi/markdown@0.6.2

moon fetch: ok
  path: .repos/mizchi/markdown/0.6.2

Files:
  moon.mod.json
  README.md
  LICENSE
  src/*.mbt
  src/api/moon.pkg
  src/html_tests/*.mbt
  benches/*.mbt

Next:
  seekmoon probe 1 --target js
  seekmoon source 1 --open README.md
```

状态变化：

```text
source_fetch.status = ok
source path recorded
```

### Event E8：比较候选

用户动作：

```bash
seekmoon compare 1 2 3
```

终端输出：

```text
Compare: markdown candidates

field             #1 mizchi/markdown        #2 moonbit-community/cmark     #3 shiguri-01/markdown
version           0.6.2                     0.4.4                          0.1.1
license           MIT                       Apache-2                       MIT
build             success                   success                        success
repository        present                   present                        present
downloads         7567                      3120                           91
packages          18                        4                              3
api_data          loaded                    not loaded                     not loaded
source            fetched                   not fetched                    not fetched
target js         unverified                unverified                     unverified

Next:
  seekmoon api 2
  seekmoon probe 1 --target js
  seekmoon probe 2 --target js
```

状态变化：

```text
comparison table stored
```

### Event E9：本地验证

用户动作：

```bash
seekmoon probe 1 --target js --keep
```

终端输出：

```text
Probe: mizchi/markdown@0.6.2       target: js

workspace:
  .seekmoon/probes/mizchi-markdown-0.6.2-js

steps:
  moon new probe                  ok
  moon add mizchi/markdown        ok
  moon check                      ok
  moon test                       ok
  moon check --target js          ok
  moon build --target js          ok

Result:
  local_add:          ok
  local_check:        ok
  local_target_build: ok
  status:             verifiable for target js in probe context

Next:
  seekmoon record 1 --conclusion continue-verification
  seekmoon report --format markdown
```

状态变化：

```text
candidate #1 verification[js] = ok
probe path stored
```

### Event E10：记录结论

用户动作：

```bash
seekmoon record 1 --conclusion continue-verification --note 'Need to confirm GFM table and task-list coverage before adoption.'
```

终端输出：

```text
Record saved

module:      mizchi/markdown
version:     0.6.2
conclusion:  continue-verification
note:        Need to confirm GFM table and task-list coverage before adoption.

Evidence:
  manifest:       ok
  module_index:   ok
  package_data:   ok
  source_fetch:   ok
  probe js:       ok

Not confirmed:
  GFM coverage
  project-specific AST compatibility
```

状态变化：

```text
adoption record stored
```

### Event E11：生成报告

用户动作：

```bash
seekmoon report --format markdown
```

终端输出：

```text
Report written:
  .seekmoon/reports/markdown-2026-06-21.md
```

报告摘要：

```text
目标：寻找 Markdown parser。
调查日期：2026-06-21。
工具链：moon 0.1.20260608。
数据来源：modules API、statistics API、manifest API、module_index、package_data、moon fetch、本地 probe、GitHub repository metadata。
候选数量：5。
高相关候选：mizchi/markdown、moonbit-community/cmark、shiguri-01/markdown。
本地验证：mizchi/markdown 通过 js target probe。
不能确认：GFM 覆盖、项目特定 AST 适配。
结论：mizchi/markdown 进入继续验证。
```

状态变化：

```text
human-readable report created
machine-readable session remains available
```

### Event E12：Code Agent 获取 JSON 和 schema

另一个消费者是 Code Agent。

Agent 动作：

```bash
seekmoon search markdown --target js --json > search.json
seekmoon schema search > search.schema.json
seekmoon shape search
```

终端输出：

```text
seekmoon.search-results.v1

schema: string
snapshot:
  fetched_at: datetime
  modules_count: int
query:
  text: string
  kind: library|skill|all
  target: string|null
results[]:
  rank: int
  module: string
  version: string
  description: string|null
  license: evidence<string>
  repository: evidence<string>
  match:
    fields: string[]
    reason: string
```

Agent 下一步：

```bash
jq -r '.results[] | select(.license.value == "MIT") | .module' search.json
```

这个事件证明 schema/shape 不是装饰。它减少了 agent 上下文消耗，减少了试错请求，也减少了用户复制 SOP 的 token 成本。

## 九、黑盒验收条件

这个工具成立，不看代码有多少，看这些行为是否成立。

1. 用户不手写 Mooncakes API URL，也能搜索包。
2. 用户不写 jq，也能看到候选摘要。
3. 用户能用 `--json` 拿机器稳定输出。
4. 用户能用 `schema/shape/fields` 在不发真实查询的情况下理解 JSON 输出。
5. 用户能从 search 编号直接 view/compare/probe/record。
6. 用户能读取 manifest、module_index、package_data，而不用知道 assets URL。
7. 用户能用 `moon ide doc` 或 package_data 下钻 API。
8. 用户能用 `moon fetch` 或 source zip 获取发布版源码。
9. 用户能把 skill/executable 和 library dependency 分开看。
10. 用户能看到影响行动的 missing/unknown，但不会看到虚空未来字段。
11. 用户能本地 probe，并记录 target 验证结果。
12. 用户能生成可复查报告。
13. Code Agent 能拿 JSON、schema 和 shape，直接做管道处理。
14. 工具能说明每个字段来源，不把 GitHub、Mooncakes、local cache 混成同一种事实。
15. 工具失败时给恢复动作，而不是裸错误码。

## 十、最终压缩定义

这款工具的核心不是“搜索命令”。它是一个 **Package Discovery 工作流环境**。

它把当前分散在 Mooncakes Web、Modules API、Manifest API、assets JSON、本地 registry index、symbols cache、`moon ide doc`、`moon fetch`、source zip、Skills API、GitHub、probe 项目里的事实，组织成开发者可以自然执行的命令序列：

```text
sync
→ search
→ view
→ api / source / skill
→ compare
→ probe
→ record
→ report
```

它默认给人 pretty text；给 agent 和脚本 JSON；给 JSON 消费者 schema、shape、fields；给深度用户 raw 和 cache 命令；给普通用户下一步动作。

它隐藏的是 API、路径、缓存、JSON 树、工具链细节；暴露的是开发者真正关心的业务对象：候选包、包用途、API、源码、license、target、build、skill、验证结果、采纳结论。

这就是当前消费者侧 MoonBit 包发现 CLI 的完整对象。
</Assistant>

<User id="2">
好的，首先第一点啊，JQ可以内置。 然后呢再就是关于这个。

“……Event E4：第一次宽搜索
用户动作：

seekmoon search markdown --target js
终端输出：

Search: markdown       target: js       kind: library       snapshot: 2026-06-21T22:35

#  module                         version   license   build    downloads   target       why
1  mizchi/markdown                 0.6.2     MIT       success  7567        unknown     name, description, keywords
2  moonbit-community/cmark         0.4.4     Apache-2  success  3120        unknown     name, description
3  shiguri-01/markdown             0.1.1     MIT       success  91          unknown     name, description
4  tonyfettes/tree_sitter_markdown 0.1.26    Apache-2  success  540         unknown     name
5  hustcer/html2md                 0.1.0     MIT       success  80          unknown     description

Notes:
  target is unknown for listed packages because no target matrix was found in current sources.
  target unknown is not treated as supported.

Next:
  seekmoon view 1
  seekmoon compare 1 2 3
  seekmoon probe 1 --target js……”
你确定有任何现代的CLI工具是这么设计的吗？ 每次搜索都有这个什么next，都有这个什么note吗？那code agent呢？比如说GPT Codex呢？那每次执行的一大堆啊，上下文啊，全都是这些东西，重复的噪声一大堆垒在一起吗？ 你自己想一想啊，想一想。而且Code Agent是什么东西呀？自回归Transformer架构的模型是什么东西？它是什么东西啊？它是Lexer还是Parser呢？什么东西？你把它想象成什么鬼呀？你把它想象成什么东西啊？人家不也是上下文吗？人家不也是注意力吗？ Code agent呢？因为看贼json么看个屁呀，要机器处理的时候才用JSON，人家自己终端，你自己不也有终端吗？你要不然用JQ，要不然呢？你就直接看pretty text，你以为你你自己不会看配置的吗？

这款工具，是吧？肯定是会有README，会有指南嘛。指南里面就要写一些东西嘛，比如说第一次使用这款工具，是吧？但比如说人，或者说Code Agent，无所谓，我管你是什么。都是消费者，都是一样的。你人不可以消费Json吗？当然可以了，干嘛不行？你想什么鬼呀？也就是这个README里面可以有一句话是这么说的，第一次使用工具，干嘛？看help嘛，这不是很自然的操作吗？是不是？然后呢，help里面要写什么？就是说每一个命令、每个子命令第一次用全部都要看help，是不是？那help不就可以把这个什么note呀、把这个next之类的机制讲好？那后面不管是人还是Code Agent，他不都不都知道要要可以可以怎么用吗？他就知道这里面的每一个数字、每一个字符、每个符号、每个东西可以怎么看、可以怎么用，一次性不就讲清楚了吗？你能想象得出来吗？

啊，这里刚好对应到这个Code Agent的每次刚好new chat新开个话题，它就是个新人，第一次使用这个工具，是不是？来看一下稍微的这个指南啊，那不就可以用了，那不就连上去了吗？ 我发现你好像脑子都不会转起来。我前面确实有提到过环境自然而然地往下走。你以为就是直接给出提示，把人家文字贴到人家脸上吗？我说的环境是什么意思？是你有这个工具，有这个名字叫search，这就是环境啊，这就是环境啊。你以为是什么东西呀？你在沙发上工作跟专业工作台工作能一样吗？这不就是环境吗？还是说我这个environment我讲的不清楚啊？我给你看篇小文献吧。啊。\\wsl.localhost\Ubuntu\home\t103o\workbench\notes\0-我\环境决定论的微观与宏观谱系：媒介生态、空间可供性与行为场景的理论综述.adoc 好好读一下这篇文档吧。然后思考一下，你之前所心里所想的那个什么环境，跟我脑子里想的是同一个东西吗？是一样的吗？100%一样吗？

然后呢，CLI工具，我前面已经明确提出了人因工程，没错，就是那个人因工程，你以为是大白话吗？也就是啊，你所有的这个命令，你给出的东西都要能够去支撑下一步的决策。最理想的东西呢，就是啊你这个命令拿到的东西可以支撑下一步动作，可能是某个命令的输入，也可能是他其他的操作活动，是不是？就是在整个工作流里面，就是顺的，从上到下，能够丝滑地往下走。这就很考验你的设计功底了，是不是？

好的啊，你本次总体而言还是可以的，但是有一些瑕疵和一些奇怪的地方。好好想一想，看一下你有什么要对我补充说明的，有什么想对我说的吗？好好想，想好之后呢，弹窗再次向我汇报，不少于3000字或者更多。开始。

Tip:
完成后记得使用弹窗与我沟通。
</User>

<Assistant id="3">
# 对消费者侧 CLI 设计的修正说明

## 一、先承认上一版的偏差

你指出的问题成立。我上一版把“环境自然而然牵引人往下走”误投影成了每次命令输出末尾都出现 `Next:`、`Notes:` 这类显式文字提示。这是一个设计层位错误。

环境不是把说明贴到用户脸上。环境是由工具的命名、命令层级、默认参数、输出形状、错误语义、help 结构、可组合数据接口、缓存机制和工作目录约定共同构成的行为场景。一个专业工作台不是在桌面上每隔十厘米贴一张“下一步请拿起螺丝刀”的便签，而是通过工具的位置、物体形状、操作路径、照明、约束和可供性，让人在场景中自然形成动作序列。

CLI 也是媒介环境。`search`、`view`、`api`、`source`、`compare`、`probe`、`record` 这些命令名本身就是可供性。它们形成行为场景中的工具陈设。用户第一次看 README 和 help 后，已经知道这些动作如何衔接；后续每次命令输出不应该重复教学。重复教学会污染工作流，消耗注意力，也会消耗 Code Agent 上下文。

所以需要修正：

- 默认输出不放 `Next:`。
- 默认输出不放长 `Notes:`。
- 提示、解释、工作流说明放在 README、`seekmoon help`、子命令 help、`seekmoon guide`、`seekmoon explain` 里。
- 命令输出只交付当前动作需要的事实。
- 如果用户想要引导，可以显式开 `--hints` 或 `--guide`，但默认关闭。
- 错误输出可以给恢复动作，因为错误中断了工作流；正常输出不需要每次教学。

这就是从“文本提示牵引”修正为“媒介环境牵引”。

## 二、读完环境文献后的校正

你给的《环境决定论的微观与宏观谱系》里有几个点对 CLI 设计非常直接。

第一，媒介不是中立管道。CLI 不只是把 Mooncakes API 搬到终端，它会重组用户的认知路径。现在用户面对的是 Web、curl、jq、GitHub、assets、moon 命令、probe 项目这些碎片环境。这个环境本身在制造认知成本。用户不是在判断包，而是在抵抗环境。CLI 的目标是重建媒介环境，让包发现动作自然发生。

第二，可供性不是说明书。一个按钮、一把手、一张工作台通过形状发出行动邀请。CLI 里的可供性就是命令名、参数名、输出列、编号引用、稳定 JSON 字段、help 示例。`seekmoon view 1` 比“请复制这个模块名再请求 manifest API”更像一个真实工具。`seekmoon api mizchi/markdown` 比“去找 package_data.json 的 URL 拼一下”更像专业工作台上的物件。

第三，行为场景依赖同形性。工具结构必须和用户的真实行为模式同形。用户真实行为不是“看一堆说明文字”，而是：搜索候选，查看详情，读 API，拿源码，比较候选，本地验证，记录结论。命令层级如果正好对应这些动作，环境就会自己牵引行为。反过来，如果每个命令都输出大段解释，就会制造语境坍塌：工具既像执行器，又像教程，又像报告，又像聊天助手，用户反而要重新分辨当前该看什么。

第四，惯习来自反复实践。一个好 CLI 会让人形成肌肉记忆：`search -> view -> api/source -> compare -> probe -> record -> report`。这个路径靠稳定命令结构形成，不靠每次输出“下一步”。输出越稳定、越低噪声，越容易形成惯习。

所以我需要把上一版中的“Next/Notes 默认输出”删掉，改成“结构化环境设计”。

## 三、人和 Code Agent 都是消费者，不应被错误二分

你说得对，我上一版里隐含了一个不严谨二分：人看 pretty text，Code Agent 看 JSON。这个二分不成立。

更准确的对象划分是：

- **终端阅读投影**：给任何正在看终端输出的消费者，包括人，也包括 Code Agent。Code Agent 同样会读 pretty text，同样受上下文长度限制，同样会被重复噪声污染。
- **机器管道投影**：给 jq、脚本、批处理、报告生成、自动化工具，也可以给人手动处理。这个投影是 JSON。
- **契约说明投影**：给第一次使用工具或写自动化的人/agent。这个投影是 help、schema、shape、fields、README。

消费者不是按“人/AI”分，而是按当前动作分：

```text
阅读判断 -> pretty text
机械处理 -> JSON / --jq
理解契约 -> help / schema / shape / fields
学习工作流 -> README / guide
排错恢复 -> error + recoverable action
```

Code Agent 不是 lexer，也不是 parser；它是带上下文窗口和注意力分配限制的自回归模型。它看终端输出时和人一样需要低噪声。每次 search 都把 `Next:`、长 `Notes:` 打出来，就是把上下文塞满重复提示。这个设计不应该进默认输出。

所以修正后的原则是：**默认输出面向终端消费者，不区分人和 agent；JSON 面向机械处理，不等于 agent 专属；help/schema/shape 面向契约学习，不混入每次业务输出。**

## 四、默认输出纪律

默认 pretty text 应该只回答当前命令的对象问题。

### `search` 的对象问题

`search` 回答：有哪些候选，为什么命中，关键采用前字段是什么。

默认输出应该像这样：

```text
Search: markdown    target: js    kind: library    snapshot: 2026-06-21T22:35

#  module                         version   license   build    downloads   target
1  mizchi/markdown                 0.6.2     MIT       success  7567        unknown
2  moonbit-community/cmark         0.4.4     Apache-2  success  3120        unknown
3  shiguri-01/markdown             0.1.1     MIT       success  91          unknown
4  tonyfettes/tree_sitter_markdown 0.1.26    Apache-2  success  540         unknown
5  hustcer/html2md                 0.1.0     MIT       success  80          unknown
```

如果 `why` 很重要，可以通过 `--why` 显示，不默认占列：

```bash
seekmoon search markdown --why
```

输出增加：

```text
why
name, description, keywords
name, description
name
```

### `view` 的对象问题

`view` 回答：这个候选的证据画像是什么。

默认输出：

```text
mizchi/markdown 0.6.2

description  Incremental Markdown parser and compiler
license      MIT
repository   https://github.com/mizchi/markdown.mbt
downloads    7567
build        success
docs         https://mooncakes.io/docs/mizchi/markdown

packages     18
api index    available
source       not fetched
js target    unknown
```

不输出教程。

### `api` 的对象问题

`api` 回答：这个 package 暴露什么 API。

```text
API: mizchi/markdown src/api

types
  Document
  Node
  ParseOptions
  RenderOptions

functions
  parse(input : String) -> Document
  render(doc : Document) -> String
  to_json(doc : Document) -> Json

source loc  available
```

### `probe` 的对象问题

`probe` 回答：候选在本地验证中是否通过。

```text
Probe: mizchi/markdown@0.6.2    target: js

moon add                 ok
moon check               ok
moon test                ok
moon check --target js   ok
moon build --target js   ok

status  verifiable
path    .seekmoon/probes/mizchi-markdown-0.6.2-js
```

### 错误输出例外

错误输出允许恢复动作，因为错误不是正常结果，是工作流断点。

```text
package_data unavailable

module   moonbitlang/core
package  argparse
status   404 resource.json, 200 package_data.json

meaning  resource.json is unavailable; API data still exists.
recover  seekmoon api moonbitlang/core/argparse --package-data
```

这种恢复动作不是日常噪声，而是故障处理。

## 五、README、guide、help 的职责

真正的引导放在三个位置。

### README

README 说明工具是什么、解决什么摩擦、最短路径是什么。

README 的首屏应该给这类内容：

```text
SeekMoon is a local MoonBit package discovery adapter.

Common workflow:
  seekmoon sync
  seekmoon search markdown --target js
  seekmoon view 1
  seekmoon api 1
  seekmoon source 1
  seekmoon probe 1 --target js
  seekmoon record 1 --conclusion continue-verification
  seekmoon report --format markdown

Output modes:
  default       pretty text for terminal reading
  --json        machine-readable output
  --jq EXPR     run jq against JSON output
  --shape       human-readable output shape
  --schema      JSON Schema for this command
  --fields      selectable fields
```

这里一次性讲清楚。

### `seekmoon help`

顶层 help 说明命令地图：

```text
Discovery
  search      find candidate modules or skills
  view        show one candidate's evidence profile
  compare     compare candidates from names or a search session

Inspection
  api         inspect package API from package_data or moon ide doc
  source      fetch or locate published source
  skill       inspect runwasm skill entries

Verification
  probe       run moon add/check/test/build in a probe project
  record      save an adoption decision
  report      render a review report

Data
  sync        update Mooncakes and local registry snapshots
  raw         print raw upstream data
  schema      print JSON Schema
  shape       print output field tree
  fields      list selectable JSON fields
  doctor      inspect local environment
```

这就是环境布局。

### 子命令 help

每个子命令 help 说明输入、输出、常用例子、JSON/schema/shape。

例如：

```bash
seekmoon search --help
```

输出应包含：

```text
Usage:
  seekmoon search <query> [--target <target>] [--kind library|skill|all]

Default searches module name, description, keywords, repository.
It does not call /api/v0/modules?search= because Mooncakes does not currently
filter that endpoint.

Examples:
  seekmoon search markdown
  seekmoon search json --target js
  seekmoon search cowsay --kind skill
  seekmoon search markdown --json
  seekmoon search markdown --jq '.results[].module'

Inspect output contract:
  seekmoon search --shape
  seekmoon search --schema
  seekmoon search --fields
```

help 是教学场景；业务输出不是教学场景。

## 六、JQ 内置的定位

你说 JQ 可以内置，这个决定是对的。

原因不是“用户不会 jq”，而是工具既然提供 JSON 公共契约，就应该提供最短机械处理路径。内置 jq 让人和 agent 可以直接写：

```bash
seekmoon search markdown --jq '.results[].module'
seekmoon view mizchi/markdown --jq '.manifest.metadata.license'
seekmoon api 1 --jq '.types[].name'
```

这比：

```bash
seekmoon search markdown --json | jq '.results[].module'
```

少一个管道，也减少 shell quoting 的局部摩擦。两种方式都应该支持。

`--jq` 的公共契约是：

```text
先生成该命令的 JSON 输出，再对该 JSON 执行 jq 表达式。
```

它不处理 pretty text，不改变命令行为，不改变数据源。

还可以支持：

```bash
seekmoon search markdown --json module,version,license
```

但这个和 `--jq` 是两个对象：

- `--json fields` 是字段选择，适合常用窄输出。
- `--jq expr` 是结构查询，适合表达式处理。

## 七、schema / shape 的修正位置

schema/shape 仍然成立，但它们不应该污染每次输出。

它们属于契约查看命令：

```bash
seekmoon search --schema
seekmoon search --shape
seekmoon search --fields
seekmoon view --schema
seekmoon api --shape
```

也可以有顶层形式：

```bash
seekmoon schema search
seekmoon shape search
seekmoon fields search
```

我更倾向两种都支持：

- `seekmoon search --shape` 对第一次用某个命令的人最自然。
- `seekmoon shape search` 对工具化脚本更统一。

但默认 `search markdown` 不显示 schema/shape。

## 八、字段显示纪律再校准

你强调的这一点必须写进设计原则。

### 未来字段不显示

没有进入当前命令 schema 的字段，完全不显示。

例如当前没有 advisory 数据源，那么默认输出和 JSON 都不应该有：

```json
"advisory": "unknown"
```

因为这是虚空字段。它会让用户误以为工具有 advisory 子系统，只是当前没查到。错误。

### 当前 schema 字段没有值，才显示 missing / unknown

例如 `search --target js` 当前 schema 包含 target 状态。target 数据源查不到，显示：

```text
target unknown
```

因为 target 是当前用户显式带入的判断上下文。

例如 `license` 是长期采用判断中的当前字段。manifest 没有 license，显示：

```text
license missing
```

### 没请求的深字段不显示

`search` 不显示 API types。

`view` 可以显示 `api index available`，不展开 API。

`api` 才展开 API。

`source` 才展开源码路径。

`skill` 才展开 wasm/checksum。

这就是环境里的抽屉，不是把所有工具摊满桌面。

## 九、命令结构如何形成环境

环境不是提示语，环境是对象排列。

这款工具的场景结构应该是：

```text
入口层
  sync
  doctor
  help

发现层
  search
  skill search

下钻层
  view
  api
  source
  skill view

判断层
  compare
  probe

记录层
  record
  report

机器契约层
  json
  jq
  schema
  shape
  fields

底层逃生口
  raw
  cache
```

这个结构本身就是工作流。用户不需要每次被提醒下一步，因为命令地图已经表达了行动空间。

## 十、重新写一个低噪声用户旅程

这次不再把 `Next` 放进业务输出。

### S0：新 chat / 新用户第一次使用

用户或 Code Agent 进入项目：

```bash
cd ~/workbench/projects/my-moonbit-tool
seekmoon help
```

输出命令地图。用户看到 `search/view/api/source/compare/probe/record/report`。

然后看 search help：

```bash
seekmoon search --help
```

知道 search 支持：

```text
--target
--kind
--json
--jq
--shape
--schema
--fields
```

这一次性建立环境认知。

### S1：同步数据

```bash
seekmoon sync
```

输出：

```text
snapshot  2026-06-21T22:35:10+08:00
modules   1350
packages  12008
index     1363 files
symbols   3 modules
moon      0.1.20260608
```

没有解释段落。

### S2：搜索候选

```bash
seekmoon search markdown --target js
```

输出：

```text
Search: markdown    target: js    kind: library    snapshot: 2026-06-21T22:35

#  module                         version   license   build    downloads   target
1  mizchi/markdown                 0.6.2     MIT       success  7567        unknown
2  moonbit-community/cmark         0.4.4     Apache-2  success  3120        unknown
3  shiguri-01/markdown             0.1.1     MIT       success  91          unknown
4  tonyfettes/tree_sitter_markdown 0.1.26    Apache-2  success  540         unknown
5  hustcer/html2md                 0.1.0     MIT       success  80          unknown
```

这里 `target unknown` 是当前 schema 字段，因为用户显式传了 `--target js`。没有 `Notes`。

### S3：查看字段来源

用户想知道为什么 target 是 unknown：

```bash
seekmoon search markdown --target js --why
```

输出增加 why/source：

```text
#  module             target   target_source   match
1  mizchi/markdown    unknown  no target field name,description,keywords
```

这是显式请求解释，不默认噪声。

### S4：查看候选详情

```bash
seekmoon view 1
```

输出：

```text
mizchi/markdown 0.6.2

description  Incremental Markdown parser and compiler
license      MIT
repository   https://github.com/mizchi/markdown.mbt
downloads    7567
build        success
docs         https://mooncakes.io/docs/mizchi/markdown

packages     18
api index    available
source       not fetched
target js    unknown
```

### S5：下钻 API

```bash
seekmoon api 1 --package src/api
```

输出：

```text
API: mizchi/markdown src/api

types
  Document
  Node
  ParseOptions
  RenderOptions

functions
  parse(input : String) -> Document
  render(doc : Document) -> String
  to_json(doc : Document) -> Json

source loc  available
```

### S6：拿源码

```bash
seekmoon source 1
```

输出：

```text
source  mizchi/markdown@0.6.2
status  fetched
method  moon fetch
path    .repos/mizchi/markdown/0.6.2
files   moon.mod.json, README.md, LICENSE, src/, benches/
```

### S7：比较候选

```bash
seekmoon compare 1 2 3
```

输出：

```text
field       mizchi/markdown   moonbit-community/cmark   shiguri-01/markdown
version     0.6.2             0.4.4                     0.1.1
license     MIT               Apache-2                  MIT
build       success           success                   success
repo        present           present                   present
packages    18                4                         3
api         loaded            indexed                   indexed
source      fetched           not fetched               not fetched
target js   unverified        unverified                unverified
```

### S8：本地验证

```bash
seekmoon probe 1 --target js --keep
```

输出：

```text
probe   mizchi/markdown@0.6.2
target  js
path    .seekmoon/probes/mizchi-markdown-0.6.2-js

moon new                 ok
moon add                 ok
moon check               ok
moon test                ok
moon check --target js   ok
moon build --target js   ok

status  verifiable
```

### S9：记录结论

```bash
seekmoon record 1 --conclusion continue-verification --note 'Need GFM coverage check.'
```

输出：

```text
recorded  mizchi/markdown@0.6.2
status    continue-verification
file      .seekmoon/records/mizchi-markdown-0.6.2.json
```

### S10：生成报告

```bash
seekmoon report --format markdown
```

输出：

```text
report  .seekmoon/reports/markdown-2026-06-21.md
```

这个旅程才是低噪声环境。工作流靠命令结构和状态文件贯穿，不靠每次输出教学。

## 十一、对上一版设计的保留与删除

保留：

- 数据源清单。
- pretty text / JSON / schema / shape / fields 的分层。
- `--jq`，并明确内置。
- search/view/api/source/compare/probe/record/report 的工作流。
- future field 不显示，当前 schema 字段缺值才 unknown。
- skill 和 library 分开。

删除或改为可选：

- 默认 `Next:`。
- 默认长 `Notes:`。
- 把 Code Agent 特殊化成 JSON 消费者。
- 在普通输出里重复教学。

新增校正：

- README / guide / help 承担学习场景。
- 正常命令输出承担业务事实投影。
- 错误输出承担恢复动作。
- `--hints` 可以作为显式教学模式，但默认关闭。

## 十二、最终补充说明

你说的 environment 我现在校准为：**命令结构、输出契约、参数可供性、help 文档、状态缓存、错误恢复和可组合数据形态共同形成的行为场景。**

这和我上一版默认在输出里塞 `Next` 不是同一个东西。上一版更像“说明书环境”；你要的是“工作台环境”。工作台环境不靠反复说教，它靠工具陈设、动作路径和对象形状让人自然行动。

所以这款 CLI 的设计基准应改成：

```text
默认输出低噪声。
命令命名承担可供性。
help 承担一次性学习。
schema/shape/fields 承担机器契约学习。
--jq 承担内置机械处理。
错误输出承担恢复提示。
状态文件承担跨命令连续性。
```

这样它才是开发工具，不是会说话的 SOP。
</Assistant>

<User id="3">
？？？？？？？？
“……如果 why 很重要，可以通过 --why 显示，不默认占列：……”什么叫如果很重要？到底是谁要啊？你形式论证上、论证推导上你都通过了吗？你证明了谁要这个东西吗？消费者是谁？是谁？你留在那里拉屎还是拉坑啊？什么鬼呀？你讲什么鬼话？你是说你自己不确定你还告诉我吗？我前面老早就讲过，讲过五六千万次了。你自己不确定的不要来找我，你想干什么呀？而且我就很莫名其妙了。你自己不是在模拟吗？你自己不在看吗？你自己，你自己呀，你自己不就在看这些东西吗？你就问自己，问你自己， yourself就是你， you you you就是你呀。你问你自己要不要吗？你自己在做这件事情，你想象一下嘛，你要吗？你要不要它嘛？你不也是一个消费者吗？你不也可以是一个普通开发者吗？

“……次看 README 和 help 后，已经知道这些动作如何衔接；后续每次命令输出不应该重复教学。重复教学会污染工作流，消耗注意力，也会消耗 Code Agent 上下文。

所以需要修正：

默认输出不放 Next:。
默认输出不放长 Notes:。
提示、解释、工作流说明放在 README、seekmoon help、子命令 help、seekmoon guide、seekmoon explain 里。
命令输出只交付当前动作需要的事实。
如果用户想要引导，可以显式开 --hints 或 --guide，但默认关闭。
错误输出可以给恢复动作，因为错误中断了工作流；正常输出不需要每次教学。……”
什么默认不放？默认不放，这不恶心啊？是为什么要放呢？我们是从哪个地方？是你从哪个地方？严谨推导说是谁？是甲方爸爸还是某个还是某个欧洲的什么政府条约条款在限制你一定要做这个呢？什么叫用户想要引导？哪个用户？告诉我哪个用户？还用户想要引导？啊？你设计的东西每一个都要有无可替代的价值，一点都不能重合，完全正交，你有做到吗？那你在这个README里面，大help、小help、子命令help都引导过去了，那你还要这个模式干什么？你就告诉我到底是给谁的？你能证明出来就留下，证明不出来就滚蛋。你是怎么设计东西的？

你是在设计东西吗？一坨屎啊。一下这些可能，一下那些可能。你全程没有出现可能，但是我感觉你从头到尾的语气就是，你根本就没有告诉我你的设计是什么，一坨屎。啊，这边可能带个药，那边可能带个药，要个屁呀。老板要喝水，给他一杯水吗？你给他一个火山干什么呢？老板可能想喝火山，是吗？你的意思是说默认展示一个火山，默认不展示火山。来，我把你的句子换成这样子，你能看得到了吗？啊，一模一样，同构的。“……默认输出不放 Next:。
默认输出不放长 Notes:。
提示、解释、工作流说明放在 R……”
换成火山这个词好听吗？

“……也可以有顶层形式：

seekmoon schema search
seekmoon shape search
seekmoon fields search
我更倾向两种都支持：

seekmoon search --shape 对第一次用某个命令的人最自然。
seekmoon shape search 对工具化脚本更统一。……”你什么判断啊？什么叫你倾向于呀？你是论证出来的，你的论证合理，我无条件同意。我根本就不看你这些狗屎词，你讲这些狗屎词，我看都不看一眼，直接否定。你不觉得你很奇怪吗？

好好想一想，想好之后，使用弹窗再来回答我吧。

Tip:
完成后记得使用弹窗与我沟通。
</User>

<Assistant id="4">
# SeekMoon CLI 设计收束版

## 一、前面的问题在哪里

你刚才指出的是设计资格问题，不是措辞问题。

我上一版里保留了几个没有完成准入审查的东西：`--why`、`--hints`、`--guide`、`seekmoon schema search` 与 `seekmoon search --schema` 双入口、`fields`。这些东西的问题不是“不一定要做”，而是我没有证明它们分别服务哪个消费者动作、改变哪个判断、减少哪项成本、为什么不能被已有结构吸收。没有完成这个证明，它们就不能进当前对象。

我现在把它们重新审查一遍，结论如下：

- 删除 `--why`。搜索命中原因是排序内部证据，不是默认终端动作所必需的表面单位。需要机器复查时，`--json` 里保留 `match` 结构；需要人工看结果时，默认 pretty text 不展示它。没有单独 `--why` 命令面。
- 删除 `--hints`。README、顶层 help、子命令 help 已经承担学习职责；正常业务输出不承担教学职责；错误输出承担恢复职责。`--hints` 没有独立消费者动作。
- 删除 `seekmoon guide`。README 和 `seekmoon help` 足够形成入口指南；再加一个 guide 会和 help/README 重叠。
- 删除顶层 `seekmoon schema search` / `seekmoon shape search`。输出契约属于具体命令，所以公共面只保留 `seekmoon search --schema`、`seekmoon search --shape` 这种局部入口。
- 删除 `--fields`。`--shape` 已经给字段树和类型，`--jq` 已经提供机械选择能力；单独 fields 与 shape 重叠。

保留：

- 默认 pretty text。
- `--json`。
- 内置 `--jq`。
- 每个输出型命令的 `--shape`。
- 每个输出型命令的 `--schema`。
- README。
- 顶层 help。
- 子命令 help。
- 错误恢复信息。

这次设计不再保留“也许有人要”的坑。每个留下来的公共面都有明确消费者动作。

## 二、当前工具的对象定义

**SeekMoon 是一个本地 MoonBit 包发现适配器。它读取 Mooncakes、Moon 工具链、本地 registry、本地项目和 GitHub 维护信号中的已存在事实，把消费者侧包发现组织为一组低噪声、可组合、可复查的 CLI 动作。**

它解决的摩擦是：消费者为了找包、看包、比较包、读 API、拿发布源码、识别 skill、做本地验证和记录结论，不再手工拼 API、jq、docs URL、GitHub 和 probe 项目。

它不解决的对象：

- 不生成作者侧 metadata。
- 不做发布门禁。
- 不替官方做 registry 治理。
- 不宣称包安全。
- 不替用户做业务选型。
- 不把 future fields 当当前事实展示。

它的完成判断是：用户能从“我要找一个包”自然走到“我有可复查采纳结论”，且过程中不需要阅读几百行 SOP 或手工拼接底层入口。

## 三、消费者是谁

这里的消费者按动作分类，不按人类/AI 分类。

### 1. 终端阅读消费者

动作：在终端读结果，做下一步判断。

可以是人，也可以是 Code Agent。两者都受注意力和上下文限制。这个消费者需要 pretty text。

### 2. 管道处理消费者

动作：把结果交给 jq、脚本、报告生成、批处理或自动化流程。

可以是人写管道，也可以是 agent 写管道。这个消费者需要 `--json` 和 `--jq`。

### 3. 契约学习消费者

动作：第一次使用命令、写自动化、确认 JSON 字段结构。

这个消费者需要 README、help、`--shape`、`--schema`。

### 4. 故障恢复消费者

动作：命令失败后判断如何恢复。

这个消费者需要错误信息中的原因和恢复动作。

这个划分替代“人看 pretty text、agent 看 JSON”的错误二分。Code Agent 看终端时也是终端阅读消费者；人写管道时也是管道处理消费者。

## 四、环境是什么

这里的环境不是输出里的提示文字。

环境由这些结构形成：

```text
命令名
子命令层级
参数名
默认输出列
编号引用
状态缓存
错误恢复
README
help
JSON 契约
shape/schema
本地记录文件
```

这些结构共同形成行为场景。用户第一次通过 README / help 建立工具地图，后续靠命令结构和稳定输出形成惯习。

一个专业工作台不靠反复贴提示条让人工作；它靠工具的位置、形状和约束使行为自然发生。SeekMoon 也是这样。`search`、`view`、`api`、`source`、`compare`、`probe`、`record`、`report` 这些命令就是工作台上的工具位置。

因此默认输出不出现教学文本。不是“默认不放 Next”，而是：**正常结果的对象不是教学，正常结果只投影当前动作需要的事实。**

## 五、当前 100% 数据来源

SeekMoon 当前明确消费以下数据源。

### 1. Modules API

```text
https://mooncakes.io/api/v0/modules
```

用途：生成 module 候选集合。

字段：

```text
name
version
description
keywords
repository
license
```

`/api/v0/modules?search=<query>` 当前不服务端过滤，所以工具本地过滤。

### 2. Statistics API

```text
https://mooncakes.io/api/v0/modules/statistics
```

用途：记录快照规模。

字段：

```text
total_modules
total_packages
total_lines
total_downloads
```

### 3. Manifest API

```text
https://mooncakes.io/api/v0/manifest/<owner>/<module>
```

用途：单 module 详情。

字段：

```text
name
module
version
latest_version
downloads
has_package
build_status
metadata.repository
metadata.license
metadata.keywords
metadata.description
metadata.checksum
metadata.targets
metadata.supported-targets
versions
```

### 4. Module Index Asset

```text
https://mooncakes.io/assets/<owner>/<module>@<version>/module_index.json
```

用途：module 内 package/API 索引。

对象：

```text
package path
types
traits
errors
typealias
values
methods
misc
child tree
```

### 5. Package Data Asset

```text
https://mooncakes.io/assets/<owner>/<module>@<version>/<pkg-relpath>/package_data.json
```

用途：单 package API 详情。

对象：

```text
types[].name
types[].docstring
types[].signature
types[].loc
types[].methods
values
traits
errors
typealias
misc
```

### 6. Package Resources Asset

```text
https://mooncakes.io/assets/<owner>/<module>@<version>/<pkg-relpath>/resource.json
```

用途：README 等资源。

状态规则：404 只表示 resources asset 不存在，不表示 package 不存在。

### 7. Source Zip

```text
https://download.mooncakes.io/user/<owner>/<module>/<version>.zip
```

用途：下载发布版源码。

状态规则：失败记录为 source_zip_failed，回退到 `moon fetch`、本地 cache、`~/.moon/lib` 或 GitHub。

### 8. Skills API

```text
https://mooncakes.io/api/v0/skills
```

用途：发现可执行 Wasm / skill entry。

字段：

```text
module
author
version
package
name
detail_url
wasm_url
checksum_url
metadata.description
metadata.name
repository
created_at
```

### 9. Wasm / Checksum Asset

来自 Skills API 的：

```text
wasm_url
checksum_url
```

用途：记录 skill 执行 artifact 与校验入口。

### 10. 本地 Registry Index

```text
~/.moon/registry/index/user/**/*.index
```

用途：交叉验证原始发布记录、离线搜索。

格式：JSON Lines。

### 11. 本地 Symbols Cache

```text
~/.moon/registry/symbols/**/*.symbols
```

用途：已缓存模块的 API / symbol 搜索。

格式：JSON Lines。

字段：

```text
module
version
key
pkg
name
parent
kind
sig_
doc
attrs
```

状态规则：不是全 registry 覆盖，命中失败不证明无包。

### 12. 本地 Moon 命令

```text
moon --version
moon update
moon ide doc
moon fetch
moon add
moon check
moon test
moon build --target
moon runwasm
moon tree
```

用途：工具链版本、registry 更新、API 文档、源码获取、依赖验证、target 验证、skill 执行、依赖树观察。

### 13. 本地 Moon 缓存

```text
~/.moon/registry/cache/**
~/.moon/registry/cache/assets/**
~/.moon/lib/core/**
```

用途：已下载源码、Wasm asset、core 源码。

### 14. 当前项目上下文

```text
moon.mod / moon.mod.json
moon.pkg / moon.pkg.json
workspace
当前 deps
当前 target
```

用途：contextual search、probe、record。

### 15. GitHub Repository

来源：manifest / modules API 中的 repository。

用途：维护和协作信号。

字段：

```text
archived
pushed_at
updated_at
licenseInfo
issues
PRs
releases
CI
README
default branch
unpublished code
```

边界：GitHub 不是发布版源码唯一事实源。

### 16. SeekMoon 本地记录

```text
.seekmoon/snapshots/*.json
.seekmoon/sessions/*.json
.seekmoon/records/*.json
.seekmoon/reports/*.md
```

用途：跨命令引用、复查、报告。

## 六、字段显示规则

### 1. 当前 schema 外字段不显示

没有进入当前命令 schema 的字段不显示。

当前没有 advisory 数据源，输出里就没有 advisory。不是 `advisory: unknown`，而是没有这个字段。

### 2. 当前 schema 内字段缺值才显示 missing / unknown

`license` 是当前 schema 字段；manifest 无 license，显示 missing。

用户传 `--target js`，target 是当前查询上下文字段；无 target evidence，显示 unknown。

### 3. 没有请求的深层对象不显示

`search` 不展开 API。

`view` 显示 API index 是否可用，不展开所有函数。

`api` 展开 API。

`source` 展示源码获取。

`skill` 展示 wasm/checksum。

## 七、命令集：保留项

### 1. `doctor`

消费者动作：确认本地环境是否能运行工具。

```bash
seekmoon doctor
```

输出：工具链、网络、registry、symbols、GitHub CLI、当前项目状态。

保留理由：没有它，第一次使用或失败排查要手动检查多个外部条件。

### 2. `sync`

消费者动作：创建可复查数据快照。

```bash
seekmoon sync
```

输出：snapshot id、modules count、packages count、本地 index count、symbols count、moon version。

保留理由：没有它，报告和搜索结果无法绑定调查时间和数据口径。

### 3. `search`

消费者动作：从功能意图生成候选集合。

```bash
seekmoon search markdown --target js
seekmoon search cowsay --kind skill
```

保留理由：这是消费者侧发现入口的核心动作。

### 4. `view`

消费者动作：查看单候选证据画像。

```bash
seekmoon view 1
seekmoon view mizchi/markdown
```

保留理由：search 结果只能初筛，采纳判断需要 manifest、module index、docs、license、repository、package summary。

### 5. `api`

消费者动作：查看 package API。

```bash
seekmoon api 1 --package src/api
```

保留理由：基础库能否复用取决于 API，不是只看 description。

### 6. `source`

消费者动作：获取或定位发布版源码。

```bash
seekmoon source 1
```

保留理由：README/API 不足时需要源码、tests、examples、moon.pkg；发布源码不应让用户手工拼 zip 或 GitHub。

### 7. `skill`

消费者动作：发现和查看可执行 Wasm / skill entry。

```bash
seekmoon skill search cowsay
seekmoon skill view Yoorkin/cowsay
```

保留理由：skill 是 execution object，不应混入 library search 的默认语义；单独命令避免对象污染。

### 8. `compare`

消费者动作：横向比较候选。

```bash
seekmoon compare 1 2 3
```

保留理由：生态查重和包选择需要同屏比较；否则用户回到手工表格。

### 9. `probe`

消费者动作：本地验证候选。

```bash
seekmoon probe 1 --target js
```

保留理由：存在、相关、可理解不等于可验证；本地项目或 probe 验证是采纳前动作。

### 10. `record`

消费者动作：保存采纳结论。

```bash
seekmoon record 1 --conclusion continue-verification
```

保留理由：没有记录，搜索工作只停留在屏幕输出，不能复查。

### 11. `report`

消费者动作：生成调查报告。

```bash
seekmoon report --format markdown
```

保留理由：长期依赖采纳需要可交付、可复查文本。

### 12. `raw`

消费者动作：检查上游原始数据。

```bash
seekmoon raw manifest mizchi/markdown
seekmoon raw module-index mizchi/markdown
```

保留理由：当 pretty text 或 JSON projection 可疑时，需要回到底层事实源。它是逃生口，不是主流程。

## 八、删除项

### 删除 `--why`

消费者动作不成立。

搜索默认已经通过 rank 和 query 生成候选。人工扫描时不需要每行看 match fields；机器复查时 JSON 中有 `match` 结构。单独 `--why` 与 JSON 的 match 投影重叠。

### 删除 `--hints`

消费者动作不成立。

学习由 README/help 承担。正常输出不教学。错误输出才给恢复动作。`--hints` 夹在业务输出和教程之间，没有独立对象。

### 删除 `guide`

消费者动作不成立。

README 和 `seekmoon help` 已经是入口指南。`guide` 与它们重叠。

### 删除顶层 `schema/shape` 命令

消费者动作被局部入口覆盖。

输出契约属于具体命令，所以：

```bash
seekmoon search --shape
seekmoon search --schema
seekmoon view --shape
seekmoon view --schema
```

不提供：

```bash
seekmoon shape search
seekmoon schema search
```

### 删除 `fields`

消费者动作被 `--shape` 和 `--jq` 覆盖。

字段发现由 `--shape` 完成；字段选择由 `--jq` 完成。单独 fields 没有不可替代职责。

## 九、输出模式

每个输出型命令有四种公共投影。

### 1. 默认 pretty text

面向终端阅读消费者。

原则：低噪声、稳定列、只显示当前动作需要的事实。

### 2. `--json`

面向机器管道消费者。

原则：输出完整当前 schema 对象，包含 source/status/fetched_at/error 等机器字段。

### 3. `--jq <expr>`

面向机械筛选消费者。

原则：先生成该命令 JSON，再执行内置 jq 表达式。

示例：

```bash
seekmoon search markdown --jq '.results[].module'
seekmoon view 1 --jq '.manifest.metadata.license'
seekmoon api 1 --jq '.types[].name'
```

### 4. `--shape`

面向契约学习消费者。

原则：显示人可读字段树、类型、状态枚举。

### 5. `--schema`

面向代码生成和严格校验消费者。

原则：输出 JSON Schema。

## 十、默认输出范式

### `search`

```bash
seekmoon search markdown --target js
```

```text
Search: markdown    target: js    kind: library    snapshot: 2026-06-21T22:35

#  module                         version   license   build    downloads   target
1  mizchi/markdown                 0.6.2     MIT       success  7567        unknown
2  moonbit-community/cmark         0.4.4     Apache-2  success  3120        unknown
3  shiguri-01/markdown             0.1.1     MIT       success  91          unknown
4  tonyfettes/tree_sitter_markdown 0.1.26    Apache-2  success  540         unknown
5  hustcer/html2md                 0.1.0     MIT       success  80          unknown
```

这里没有 `Next`，没有 `Notes`，没有 `why`。

`target` 显示是因为用户传了 `--target js`。如果用户没传 target，不显示 target 列。

### `view`

```bash
seekmoon view 1
```

```text
mizchi/markdown 0.6.2

description  Incremental Markdown parser and compiler
license      MIT
repository   https://github.com/mizchi/markdown.mbt
downloads    7567
build        success
docs         https://mooncakes.io/docs/mizchi/markdown

packages     18
api index    available
source       not fetched
target js    unknown
```

### `api`

```bash
seekmoon api 1 --package src/api
```

```text
API: mizchi/markdown src/api

types
  Document
  Node
  ParseOptions
  RenderOptions

functions
  parse(input : String) -> Document
  render(doc : Document) -> String
  to_json(doc : Document) -> Json

source loc  available
```

### `source`

```bash
seekmoon source 1
```

```text
source  mizchi/markdown@0.6.2
status  fetched
method  moon fetch
path    .repos/mizchi/markdown/0.6.2
files   moon.mod.json, README.md, LICENSE, src/, benches/
```

### `compare`

```bash
seekmoon compare 1 2 3
```

```text
field       mizchi/markdown   moonbit-community/cmark   shiguri-01/markdown
version     0.6.2             0.4.4                     0.1.1
license     MIT               Apache-2                  MIT
build       success           success                   success
repo        present           present                   present
packages    18                4                         3
api         loaded            indexed                   indexed
source      fetched           not fetched               not fetched
target js   unverified        unverified                unverified
```

### `probe`

```bash
seekmoon probe 1 --target js
```

```text
probe   mizchi/markdown@0.6.2
target  js
path    .seekmoon/probes/mizchi-markdown-0.6.2-js

moon new                 ok
moon add                 ok
moon check               ok
moon test                ok
moon check --target js   ok
moon build --target js   ok

status  verifiable
```

## 十一、完整用户旅程

### E0：第一次进入环境

```bash
seekmoon help
```

用户看到命令地图。这里是学习场景，不是业务输出。

```text
Discovery
  search      find library module candidates
  skill       find executable Wasm skill entries
  view        show one candidate evidence profile

Inspection
  api         inspect package API
  source      fetch or locate published source
  compare     compare candidates

Verification
  probe       verify a candidate in a local probe project
  record      save an adoption decision
  report      render a review report

Data
  sync        update snapshots
  doctor      inspect local environment
  raw         print upstream raw data

Output
  --json      machine-readable output
  --jq EXPR   run jq against JSON output
  --shape     print output shape for this command
  --schema    print JSON Schema for this command
```

### E1：检查环境

```bash
seekmoon doctor
```

```text
moon        0.1.20260608
moonc       v0.10.0+e66899a54
registry    present
network     mooncakes.io ok
project     MoonBit module detected
github      gh available
```

### E2：同步快照

```bash
seekmoon sync
```

```text
snapshot  2026-06-21T22:35:10+08:00
modules   1350
packages  12008
index     1363 files
symbols   3 modules
```

### E3：搜索

```bash
seekmoon search markdown --target js
```

```text
Search: markdown    target: js    kind: library    snapshot: 2026-06-21T22:35

#  module                         version   license   build    downloads   target
1  mizchi/markdown                 0.6.2     MIT       success  7567        unknown
2  moonbit-community/cmark         0.4.4     Apache-2  success  3120        unknown
3  shiguri-01/markdown             0.1.1     MIT       success  91          unknown
```

### E4：查看详情

```bash
seekmoon view 1
```

```text
mizchi/markdown 0.6.2

description  Incremental Markdown parser and compiler
license      MIT
repository   https://github.com/mizchi/markdown.mbt
downloads    7567
build        success
docs         https://mooncakes.io/docs/mizchi/markdown

packages     18
api index    available
source       not fetched
target js    unknown
```

### E5：看 API

```bash
seekmoon api 1 --package src/api
```

```text
API: mizchi/markdown src/api

types
  Document
  Node
  ParseOptions
  RenderOptions

functions
  parse(input : String) -> Document
  render(doc : Document) -> String
  to_json(doc : Document) -> Json

source loc  available
```

### E6：获取发布源码

```bash
seekmoon source 1
```

```text
source  mizchi/markdown@0.6.2
status  fetched
method  moon fetch
path    .repos/mizchi/markdown/0.6.2
files   moon.mod.json, README.md, LICENSE, src/, benches/
```

### E7：比较候选

```bash
seekmoon compare 1 2 3
```

```text
field       mizchi/markdown   moonbit-community/cmark   shiguri-01/markdown
version     0.6.2             0.4.4                     0.1.1
license     MIT               Apache-2                  MIT
build       success           success                   success
repo        present           present                   present
packages    18                4                         3
api         loaded            indexed                   indexed
source      fetched           not fetched               not fetched
target js   unverified        unverified                unverified
```

### E8：本地验证

```bash
seekmoon probe 1 --target js
```

```text
probe   mizchi/markdown@0.6.2
target  js
path    .seekmoon/probes/mizchi-markdown-0.6.2-js

moon new                 ok
moon add                 ok
moon check               ok
moon test                ok
moon check --target js   ok
moon build --target js   ok

status  verifiable
```

### E9：记录结论

```bash
seekmoon record 1 --conclusion continue-verification --note 'Need GFM coverage check.'
```

```text
recorded  mizchi/markdown@0.6.2
status    continue-verification
file      .seekmoon/records/mizchi-markdown-0.6.2.json
```

### E10：生成报告

```bash
seekmoon report --format markdown
```

```text
report  .seekmoon/reports/markdown-2026-06-21.md
```

### E11：机器处理

```bash
seekmoon search markdown --json --jq '.results[].module'
```

```text
mizchi/markdown
moonbit-community/cmark
shiguri-01/markdown
```

### E12：查看 JSON 结构

```bash
seekmoon search --shape
```

```text
seekmoon.search-results.v1

schema: string
snapshot:
  fetched_at: datetime
  modules_count: int
query:
  text: string
  kind: library|skill
  target: string|null
results[]:
  rank: int
  module: string
  version: string
  description: string|null
  license:
    status: present|missing
    value: string|null
    source: string
  repository:
    status: present|missing
    value: string|null
    source: string
  match:
    fields: string[]
```

这里的 `match` 留在 JSON/shape 中，因为它是搜索结果对象的机器可复查部分；默认 pretty text 不展示它，因为默认扫描动作不需要它。

## 十二、最终设计规则

SeekMoon 的设计规则压缩如下：

1. 命令结构承担环境牵引，不靠业务输出反复教学。
2. README 和 help 承担学习职责。
3. 默认 pretty text 只输出当前动作需要的事实。
4. JSON 输出当前命令完整对象。
5. 内置 jq 服务机械处理。
6. shape/schema 服务契约学习。
7. 错误输出服务恢复动作。
8. future fields 不显示。
9. 当前 schema 字段缺值才 missing / unknown。
10. library 与 skill 分开。
11. GitHub 是维护信号源，不是发布版源码唯一事实源。
12. `moon fetch`、module index、package data、`moon ide doc` 是已知候选下钻的一等入口。

这版里不再保留没有消费者动作授权的分支。
</Assistant>
