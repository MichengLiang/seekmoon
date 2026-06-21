# MoonBit JSON 过滤 CLI / runwasm 方案调查报告

调查日期：2026-06-21  
工作目录：`/home/t103o/workbench/projects/seekmoon/spike/json-filter-wasm-cli-investigation`

## 1. 结论先行

如果目标是“用户不用自己通过系统包管理器安装 jq，也不用理解 Rust/WASM/Component Model，只通过 MoonBit 生态给出的一个命令就能得到类似 jq 的 JSON 过滤效果”，当前最优雅、最贴合 MoonBit 现状的路线是：

```bash
moon runwasm owner/jqx-skill@version -- '.foo' '{"foo":1}'
```

或者在未来补齐 stdin 后：

```bash
echo '{"foo":1}' | moon runwasm owner/jqx-skill@version -- '.foo'
```

这个方案的本体不是“内置系统 jq 二进制”，也不是“让用户的 MoonBit library 项目透明链接一个外部二进制工具”，而是“作者侧把 JSON 过滤器做成一个自足的 WebAssembly CLI artifact，并发布成 Mooncakes prebuilt wasm asset / skill；用户侧只依赖 `moon runwasm` 的公共契约”。这正好贴合你前面那份 `MoonBit runwasm CLI 场景确认.md` 里的判断：CLI/skill 场景可以把内部复杂性包起来，library dependency 场景不能假装复杂性不存在。

当前调查已经验证出一个可行的最小事实链：

1. `moon runwasm` 当前工具链支持本地 main package，也支持 Mooncakes 上的 prebuilt WebAssembly binary asset。当前本机版本是 `moon 0.1.20260608`、`moonc v0.10.0`、`moonrun 0.1.20260608`。
2. Mooncakes 当前 modules API 中已经有 jq 方向实现：`mizchi/jq`、`shina1024/jqx`、`bobzhang/moonjq`。
3. 当前 skills API 中没有发现这几个 jq 包已经发布为可直接 `moon runwasm` 的 skill/prebuilt wasm entry。
4. 实测 `moon runwasm mizchi/jq@0.2.2/cmd/moonjq -- --help` 和 `moon runwasm shina1024/jqx@0.2.0/cmd -- --help` 都失败，错误为 `Prebuilt wasm asset does not exist`。这说明“包存在”和“可一条命令远程 runwasm”是两层事实。
5. 拉取 `shina1024/jqx@0.2.0` 发布源码后，本地实测 `moon runwasm cmd -- '.foo' '{"foo":1}'` 输出 `1`。这说明 `jqx` 的 MoonBit CLI 在 wasm 目标下通过“第二个参数传 JSON”已经能跑。
6. 同一份 `jqx` 源码中，`cmd/stdin_wasm.mbt` 明确写着 Wasm fallback 下 stdin 不可用并返回空字符串；实测 `printf '{"foo":1}\n' | moon runwasm cmd -- '.foo'` 无输出。这说明管道式 stdin 不是当前可承诺能力，需要单独补。

因此，短期如果要做一个可以交付的“MoonBit 生态里的 jq-like JSON 过滤 CLI”，最稳的是选择 `shina1024/jqx` 作为基础，做一个 runwasm/skill 发布层。第一版公共契约可以先明确采用参数输入：

```bash
moon runwasm owner/jqx@version -- FILTER JSON_TEXT
```

然后把 stdin 管道能力列为后续工程任务，而不是在当前没有证据时承诺“完全 jq 同款管道体验”。

## 2. 问题对象与边界

这里的需求不是“MoonBit 项目里有 JSON parser”这么宽，也不是“把官方 jq 的所有功能复制一遍”。从你描述的使用场景看，核心摩擦是：用户想在 MoonBit 生态或 AI agent 场景中临时对一段 JSON 做过滤、提取、转换，但不希望自己再装系统 jq，不希望进入 npm/pnpm/cargo/apt/brew 等外部包管理链路，也不希望理解内部实现。

所以当前人工制品边界可以定义为：

```text
一个可通过 moon runwasm 调用的 JSON 过滤 CLI。
输入：过滤表达式和 JSON 文本。
输出：过滤后的 JSON 或 raw text。
消费者：终端用户、脚本、AI agent、MoonBit 相关自动化流程。
不负责：成为普通 MoonBit library 的二进制依赖系统；完整替代 jq 1.8.1；让用户项目透明链接外部 Rust/WASM artifact。
```

这个边界很重要。因为“内置 jq 二进制”“发布一个 MoonBit library”“发布一个 WASM CLI skill”“npm 包装一个 CLI”“Component Model 组合 Rust 与 MoonBit”看起来都能靠近问题，但它们的消费者契约完全不同。当前最省解释成本、最符合 `moon runwasm` 现状的，是最终 CLI artifact。

## 3. 本地资料复核后的关键事实

你给的三份本地资料形成了很清楚的判断框架。

`MoonBit包搜索发现与生态查重SOP.md` 说明了 Mooncakes 搜索不能只看 Web 页面，可靠路径应使用 modules API、manifest API、module index、package data、Skills API、本地 fetch 与本地验证。这个调查按这个路径执行：拉取了 2026-06-21 的 modules API 与 skills API，检索 jq/json/query/filter/parser 等关键词，查看 manifest、module index、resource/package data，并实际 `moon fetch` 了候选源码。

`MoonBit runwasm CLI 场景确认.md` 的核心判断继续成立，并且现在有更强的命令帮助证据：`moon runwasm` 当前帮助明确说它可运行本地 package，也可运行发布为 Mooncakes asset 的 prebuilt wasm binary。未 pin 版本会从 registry index 解析 latest，并缓存到 `$MOON_HOME/registry/cache/assets`。这说明用户侧“只用一个 MoonBit 命令运行工具”的路径是官方命令面已经支持的。

`MoonBit与RustWASM互操作事实报告.adoc` 里 library 与 CLI 的差异仍然是本次方案判断的主轴。library 场景中，Rust WASM 或任意二进制能力会进入下游项目的 build/run 公共契约，不能隐藏。CLI 场景中，作者控制最终 artifact，可以把内部实现封装起来。JSON 过滤工具天然更适合 CLI 场景，因为消费者关心的是命令输入、输出、错误码，而不是把过滤器当成 MoonBit 源码 API 链接进自己的程序。

## 4. 当前 MoonBit / Mooncakes 事实快照

本机工具链：

```text
moon 0.1.20260608 (60bc8c3 2026-06-08)
moonc v0.10.0+e66899a54 (2026-06-09)
moonrun 0.1.20260608 (60bc8c3 2026-06-08)
```

`moon help runwasm` 的当前命令面说明：

```text
Run a local package as WebAssembly or a prebuilt WebAssembly binary published as a Mooncakes asset.
```

它接受的坐标形式包括：

```text
moon runwasm moonbitlang/parser/cmd/moonfmt@0.3.3
moon runwasm moonbitlang/parser@0.3.3/cmd/moonfmt
moon runwasm moonbitlang/parser/cmd/moonfmt
```

2026-06-21 拉取的 Mooncakes modules API 返回 1350 个 module。skills API 返回 70 个可执行 wasm entry。实测已有 skill 能运行，例如：

```bash
moon runwasm Yoorkin/cowsay@0.1.0 -- hello
```

能输出 cowsay 文本；又如：

```bash
moon runwasm Betterlol/moon_zod@0.5.1/cmd/json2schema -- '{"foo":1}'
```

能把 JSON 作为参数传入并输出 schema。这两个例子说明 `moon runwasm` 作为“轻量执行分发层”已经不是理论能力。

但 jq 方向候选虽然是 Mooncakes module，却还没有对应的 prebuilt wasm asset entry。这个空缺正是当前可做的工作空间。

## 5. jq / JSON 过滤候选对比

### 5.1 `shina1024/jqx`

manifest 快照：

```text
module: shina1024/jqx
version: 0.2.0
description: jq-compatible JSON processor written in MoonBit with a CLI and TypeScript bindings
keywords: jq,json,query,filter,moonbit
license: Apache-2.0
repository: https://github.com/shina1024/jqx
build_status: success
downloads: 6
versions_count: 4
```

GitHub 维护信号：

```text
repo: shina1024/jqx
archived: false
pushedAt: 2026-06-20T06:53:50Z
license: Apache-2.0
```

源码结构里有 `cmd` 包，并且 `cmd/moon.pkg` 设置：

```moonbit
options(
  "is-main": true,
  "native-stub": [ "windows_utf8_args.c" ],
  targets: {
    "args_native.mbt": [ "native", "llvm" ],
    "args_wasm.mbt": [ "wasm", "wasm-gc" ],
    "exit_native.mbt": [ "native", "llvm" ],
    "exit_wasm.mbt": [ "wasm", "wasm-gc" ],
    "stdin_native.mbt": [ "native", "llvm" ],
    "stdin_wasm.mbt": [ "wasm", "wasm-gc" ],
  },
)
```

CLI 参数解析支持：

```text
jqx [options] "<filter>" ["<json>"]
```

README 也明确给出：

```bash
jqx ".foo" '{"foo": 1}'
```

本地验证结果：

```bash
moon runwasm cmd -- '.foo' '{"foo":1}'
```

输出：

```text
1
```

这对当前需求非常关键。因为即使 stdin 管道暂时不可用，只要用户把 JSON 文本作为参数传入，`moon runwasm` 已经可以满足“接受一段 JSON 然后过滤”的最小动作轨迹。

限制也很明确。`cmd/stdin_wasm.mbt` 内容是：

```moonbit
///|
/// WASM fallback: stdin is not available, return empty input.
fn read_stdin_all_internal() -> String {
  ""
}
```

所以管道式输入当前不能承诺。实测：

```bash
printf '{"foo":1}\n' | moon runwasm cmd -- '.foo'
```

没有输出。这个限制不是 jqx 核心过滤能力的问题，而是 Wasm CLI stdin 接入问题。

### 5.2 `mizchi/jq`

manifest 快照：

```text
module: mizchi/jq
version: 0.2.2
description: A jq clone implemented in MoonBit
keywords: jq,json,moonbit
license: Apache-2.0
repository: https://github.com/mizchi/jq
build_status: success
downloads: 1942
versions_count: 5
```

README 声称兼容 jq 1.8.1 的 96.2% 已验证测试，并支持大量 jq 语法与 builtin。这个实现的过滤能力很强。

但它的 CLI package 配置显示：

```moonbit
targets: { "main.mbt": [ "native" ], "stub.mbt": [ "wasm-gc", "wasm", "js" ] }
```

`main.mbt` 是 native CLI，读取 stdin 使用 `moonbitlang/async/stdio`。在 wasm/js 目标下走 stub。因此它更像是一个强 library + native CLI，而不是当前可直接变成 runwasm JSON 参数 CLI 的最小基础。它当然可以改造，但改造量比 `jqx` 更大。

实测远程：

```bash
moon runwasm mizchi/jq@0.2.2/cmd/moonjq -- --help
```

失败：

```text
Error: Prebuilt wasm asset does not exist
```

### 5.3 `bobzhang/moonjq`

manifest 快照：

```text
module: bobzhang/moonjq
version: 0.1.0
description: A jq implementation in MoonBit
keywords: jq,json,query
license: Apache-2.0
repository: git@github.com:moonbit-community/moobit-jq.git
build_status: legacy
downloads: 65
versions_count: 1
```

它可以作为历史参考，但当前 build_status 是 legacy，版本也少。若以快速做出 runwasm CLI 为目标，它不是首选。

## 6. 为什么不建议“内置系统 jq 二进制”

“内置 jq”这句话有几种可能含义。

第一种是把官方 jq native 二进制打进某个发布包。这会立即遇到多平台问题：Linux、macOS、Windows、x86_64、arm64、glibc/musl、执行权限、杀软/下载信任、release checksum、PATH、更新策略。它能让用户不走系统包管理器，但会把分发复杂性转移给我们自己，而且不符合 `moon runwasm` 的跨平台初衷。

第二种是把 jq 编译成 WASI/WASM，然后用 `moon runwasm` 跑。这理论上可能，但它不是 MoonBit 内部能力；本质是 C/Rust/其他语言构建出来的 WASM CLI asset。它可以作为“最终 artifact”路线，但对 MoonBit 生态贡献较弱，也不复用当前已经存在的 MoonBit jq-like 实现。除非需要官方 jq 极高兼容性，否则不是最小可行路径。

第三种是让 MoonBit library 自动携带并调用某个外部 jq/WASM 二进制。这会回到你旧报告里的 library 二进制依赖缺口：文件能随 package 分发，不等于编译器、build system、runtime、host import、路径定位、版本关系都形成一等契约。当前没有证据说明 `moon add` 后可以透明获得这种能力。

所以，“内置系统 jq 二进制”不是最优雅方案。更清晰的做法是：把 JSON 过滤能力本身实现为 MoonBit CLI，并发布成 `runwasm` 能直接运行的预构建 asset。

## 7. 建议方案

### 方案 A：基于 `shina1024/jqx` 发布 runwasm skill

这是当前最推荐的方案。

用户契约：

```bash
moon runwasm owner/jqx@version -- FILTER JSON_TEXT
```

示例：

```bash
moon runwasm owner/jqx@version -- '.foo' '{"foo":1}'
```

预期输出：

```text
1
```

作者侧工作：

1. fork 或贡献上游 `shina1024/jqx`。
2. 确认 `cmd` 包作为 wasm main package 构建。
3. 发布时生成 `_build/wasm/release/build/cmd/cmd.wasm` 或等价 release wasm。
4. 让 Mooncakes 产生可被 `moon runwasm owner/module@version/cmd` 命中的 prebuilt wasm asset。
5. 在 README / SKILL.md 中明确当前输入方式是参数 JSON，不承诺 stdin。
6. 加 smoke test：`moon runwasm local-or-built -- '.foo' '{"foo":1}'` 必须输出 `1`。

这个方案的优势是：

1. 纯 MoonBit 实现，符合生态叙事。
2. 已经有 CLI、library、TS binding、CI、release workflow。
3. 本地 wasm runwasm 已验证。
4. 用户不需要 npm、apt、brew、cargo。
5. 失败面主要集中在发布 prebuilt wasm asset 与 stdin 支持，而不是重新发明 jq。

风险：

1. 当前 registry 版本没有 prebuilt wasm asset。
2. stdin 管道不可用。
3. `moon check -d` 可能因 deprecated `try?` 警告在严格模式下失败，需要清理或降低发布校验门槛。
4. 下载量低，说明还很新，长期稳定性要继续跟踪。

### 方案 B：基于 `mizchi/jq` 做 wasm CLI 改造

这个方案适合更看重 jq 兼容率时采用。`mizchi/jq` 的功能面和测试兼容报告更强，README 写明 96.2% jq 1.8.1 verified tests。

但它当前 CLI 的 main 只在 native target 下启用，wasm/js 下是 stub。要做成 runwasm skill，需要补一个 wasm-compatible main，处理 args、输入 JSON、输出结果。这个改造不是很大，但比 `jqx` 的现状多一步。

适用条件：

1. 我们愿意投入更多兼容性验证。
2. 目标是尽量贴近 jq，而不只是“类似 jq 的 JSON 过滤效果”。
3. 可以接受重写 CLI 输入层。

### 方案 C：直接做一个最小 JSON path/filter CLI

如果需求只是 `.foo`、`.items[]`、`.a.b`、数组取下标、简单 raw 输出，那么完全可以用 `moonbitlang/core/json` 自己写一个非常小的 JSON path 解释器。官方 core 已有 `@json.parse`、`Json::stringify`、`Json::object`、`Json::array` 等 API。

这个方案实现面最小、控制力最强，但能力边界必须诚实命名，不能叫 jq。它适合内部 agent 工具，不适合对外宣称 jq-compatible。

用户契约可以是：

```bash
moon runwasm owner/jsonpick@version -- '.foo.bar[0]' '{"foo":{"bar":[1]}}'
```

风险是后续需求很容易滑向 jq 全语法，最后重新发明 jq。除非明确只要 JSON path，不建议作为对外主线。

### 方案 D：Component Model / Rust WASM / 官方 jq wasm

这条路线技术上有探索价值，但不适合作为当前最短交付路径。

MoonBit 官方 Component Model 文档确实展示了 WIT、`wit-bindgen moonbit`、`wasm-tools component embed`、`wasm-tools component new`、Wasmtime 测试等工作流。Rust 也有 WASI Preview 2 / Component Model 方向。但这些工具解决的是跨语言 component 构建与组合，不会自动让 `moon runwasm` 用户更简单。

如果最终仍然发布一个 CLI artifact，Component Model 只是作者侧内部构建技术。它不应该进入用户契约。只有当我们确实需要复用 Rust/C jq 核心，并且 MoonBit 实现无法满足时，才值得走这条路。

## 8. 当前可以承诺的最小产品形态

建议把第一阶段的公共契约写得非常窄：

```text
命令：moon runwasm owner/jqx@version -- FILTER JSON_TEXT
输入：FILTER 是 jq-like 过滤表达式；JSON_TEXT 是单个 JSON 文本参数。
输出：每个结果一行，默认 JSON stringified；支持上游 CLI 已支持的 -r/-R/-n/-s/-e 等参数，以实际测试为准。
不承诺：stdin 管道、完全 jq 兼容、所有 jq module/import/input 能力、超大 JSON 流式处理、作为 MoonBit library 的二进制依赖。
```

这个契约能直接支撑用户动作：

```bash
moon runwasm owner/jqx@version -- '.users[].name' '{"users":[{"name":"alice"},{"name":"bob"}]}'
```

输出：

```text
"alice"
"bob"
```

对 AI agent 来说也足够实用，因为 agent 通常可以把 JSON 文本作为参数传入，而不是必须依赖 shell 管道。

## 9. stdin 管道能力怎么补

如果要接近 jq 的自然体验，stdin 是必须补的。当前 `jqx` 的 wasm target 已经把 stdin 层隔离在 `stdin_wasm.mbt`，这反而是好事：改造点集中。

需要确认的问题：

1. `moon runwasm` 运行的 wasm 当前暴露给程序的 stdin 是什么 host 能力。
2. MoonBit core 或 x/async 是否已有 WASI stdin 可用 API。
3. 是否可以通过 `moonbitlang/async/stdio` 在 wasm target 下读取 stdin，还是该包当前只有 native 可用。
4. 如果没有现成 API，是否需要写 WASI import 绑定，或等待/贡献 MoonBit 标准库层支持。

本次 `moon ide doc '@moonbitlang/async/stdio'` 只显示 `pub let unimplemented : Unit`，不能作为可用 stdin API 证据。`moonbitlang/core/env` 则明确有 `args()`，这解释了为什么 `jqx` 的 wasm 参数读取可用。

所以 stdin 的后续路径应该是工程验证，而不是口头承诺：

```text
先做一个最小 wasm main，只读取 stdin 并 println 长度。
用 moon runwasm 本地跑管道。
如果能读，再接入 jqx 的 stdin_wasm.mbt。
如果不能读，保持 JSON 参数输入作为稳定契约。
```

## 10. 发布与生态建议

最好的生态路线不是在我们这边长期 fork 出一个“又一个 jq”，而是：

1. 向 `shina1024/jqx` 提 issue 或 PR：增加 Mooncakes prebuilt wasm / skill 发布。
2. 在 release workflow 中增加 wasm build smoke。
3. 在 README 增加 `moon runwasm` 用法。
4. 如果上游愿意，直接由上游发布 `shina1024/jqx@next/cmd` 的 prebuilt wasm asset。
5. 如果上游暂时不处理，可以在 seekmoon 下做薄封装模块，但要在文档里声明它是 `jqx` 的 runwasm 分发层。

这个选择能避免生态分裂，也让维护成本落在真实 jqx 实现上，而不是我们维护一份复制逻辑。

如果要做自有封装，建议名字避免叫 `jq`，可以叫：

```text
seekmoon/json-filter
seekmoon/jqx-runwasm
seekmoon/jsonq
```

但如果只是发布 `jqx` 的 CLI artifact，最好名称里保留 jqx 来源，避免用户误以为它是官方 jq 完整替代。

## 11. 本次产生的本地材料

本次 spike 目录包含：

```text
projects/seekmoon/spike/json-filter-wasm-cli-investigation/
├── .repos/shina1024/jqx/0.2.0
├── .repos/mizchi/jq/0.2.2
├── moon-commands-current.md
├── skills-2026-06-21.tsv
└── report.md
```

`.repos` 是通过 `moon fetch` 拉取的发布版源码，用于复核发布包内容。`skills-2026-06-21.tsv` 是当天 skills API 的简表快照。`moon-commands-current.md` 是从 Moon 官方仓库拉取的当前命令手册。

## 12. 最终判断

当前可行的优雅解决方案是：

```text
基于现有 MoonBit jq-compatible 实现，发布一个自足的 Mooncakes prebuilt wasm CLI / skill，让用户通过 moon runwasm 直接执行。
```

首选基础是 `shina1024/jqx`，因为它已经有 wasm target 下可运行的 CLI 参数路径，并且实测成功。当前不应承诺 stdin 管道；可以先承诺 `FILTER JSON_TEXT` 参数输入。若必须支持 `echo JSON | ...`，需要单独做 WASI/stdin probe 并修改 `stdin_wasm.mbt`。

不建议把“内置 jq”理解成塞系统 native 二进制，也不建议把这个问题做成普通 MoonBit library 的透明二进制依赖。那些路线都会扩大消费者契约，增加解释与维护成本。

最短下一步：

1. 给 `shina1024/jqx` 做一个 release wasm asset / skill 发布 probe。
2. 本地验证 `moon runwasm <published-coordinate> -- '.foo' '{"foo":1}'`。
3. 写清楚公共契约：参数 JSON 可用，stdin 暂不承诺。
4. 另开一个 stdin probe，验证 `moon runwasm` 是否能给 MoonBit wasm main 暴露 stdin。

## 参考链接

- Moon commands manual: https://github.com/moonbitlang/moon/blob/main/docs/manual/src/commands.md
- MoonBit Component Model tutorial: https://docs.moonbitlang.com/en/latest/toolchain/wasm/component-model-tutorial.html
- Mooncakes modules API: https://mooncakes.io/api/v0/modules
- Mooncakes skills API: https://mooncakes.io/api/v0/skills
- `shina1024/jqx`: https://github.com/shina1024/jqx
- `mizchi/jq`: https://github.com/mizchi/jq
- MoonBit skills repository: https://github.com/moonbitlang/skills
