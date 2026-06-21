# MoonBit `moon runwasm` 能力边界报告

调查日期：2026-06-21

## 结论

`moon runwasm` 可以理解成“Mooncakes 预构建 Wasm CLI 资产运行器 + 本地 MoonBit package 的 Wasm 运行入口”。它和 `npx` 有相似的用户体验：给一个坐标、传一些参数、马上运行，不需要用户先手动安装工具。但它不是 `npx` 的同构物。`npx` 面向 npm 包，会解析 npm registry、下载包、可能执行包内脚本、运行 JS/native wrapper；`moon runwasm` 当前只接受本地 package 或 Mooncakes 坐标，远程模式下载的是 Mooncakes 上已经构建好的 `.wasm` 文件，然后交给 `moonrun` 运行。

所以更准确的类比是：

```text
npx: package runner / installer-ish runner
moon runwasm: prebuilt Wasm asset runner
```

用户侧确实可以做到“想用就用”：

```bash
moon runwasm Yoorkin/cowsay -- hello
moon runwasm moonbitlang/parser/cmd/moonfmt -- --help
moon runwasm Betterlol/moon_zod@0.5.1/cmd/json2schema -- '{"foo":1}'
```

但它不是“任意东西都能运行”。它不能直接运行任意 URL，不能直接跑 npm 包，不能自动从 GitHub clone 项目后构建，也不能运行一个 registry module 里不存在 prebuilt asset 的 package。比如之前实测 `shina1024/jqx@0.2.0/cmd` 这个 module 和 package 存在，但 `moon runwasm shina1024/jqx@0.2.0/cmd -- --help` 失败，错误是 `Prebuilt wasm asset does not exist`。这说明“Mooncakes 上有包”和“能 runwasm 运行”之间还隔着一层预构建 Wasm 资产。

## 分发靠谁

远程 `runwasm` 的主分发层是 Mooncakes，不是 npm。

官方 MoonBit 更新说明里说得很直接：`moon runwasm` 和 SKILL marketplace 是实验性功能；可以直接运行发布在 mooncakes.io 且支持 Wasm target 的包；mooncakes.io 会为支持 Wasm 编译目标的 package 自动构建经过 wasm-opt 优化的 `*.wasm` 文件，并把同目录的 `SKILL.md` 显示到 skills 页面。

这意味着作者侧理想流程不是让用户去 npm 装一个 wrapper，而是：

1. 写一个 MoonBit package。
2. 在 `moon.pkg` 里把 CLI package 标成 `"is-main": true`。
3. 确保该 package 支持 `wasm` 或 `wasm-gc` 编译目标。
4. 发布到 Mooncakes。
5. Mooncakes 构建出对应 `.wasm` asset 和 `.wasm.sha256`。
6. 用户用 `moon runwasm owner/module[/package][@version] -- args...` 运行。

npm 不是这条链路的必要组成部分。npm 当然仍可以作为另一个分发面，比如给 JS/TS 用户发包、给 Node 环境发 wrapper、或者包装额外资源。但如果讨论的是 MoonBit 自家注册区里的 CLI 工具，主路径就是 Mooncakes asset + `moon runwasm`。

## 坐标与版本

`moon runwasm` 远程模式接受 Mooncakes 坐标，不接受任意 URL。当前命令帮助列出的形式包括：

```bash
moon runwasm moonbitlang/parser/cmd/moonfmt@0.3.3
moon runwasm moonbitlang/parser@0.3.3/cmd/moonfmt
moon runwasm moonbitlang/parser/cmd/moonfmt
```

源码里的坐标解析规则大致是：

```text
user/module/package@version
user/module@version/package
user/module/package
user/module@version
```

如果 package path 为空，二进制名取 module 最后一段；如果 package path 非空，二进制名取 package path 最后一段。资产 URL 会按这个规则拼出来：

```text
https://mooncakes.io/assets/<user>/<module>@<version>/<package>/<binary>.wasm
https://mooncakes.io/assets/<user>/<module>@<version>/<package>/<binary>.wasm.sha256
```

例如：

```text
https://mooncakes.io/assets/moonbitlang/parser@0.3.5/cmd/moonfmt/moonfmt.wasm
```

这也解释了为什么 `cmd/moonfmt` 这类 package 可以跑，而有些普通 library package 不行。`runwasm` 不会自己猜“哪个包应该是 CLI”，它就是按坐标去找那个预构建 wasm 文件。

版本策略也不像 `npx` 那样总是去问远端最新版本。源码里的开发参考文档把策略写得很清楚：

1. pinned coordinate 有显式版本，直接使用这个版本，不需要 registry index 做版本解析。
2. unpinned coordinate 没有版本时，先读本地 registry index。
3. 如果本地 index 里已有 module 的可用版本元数据，就直接用，不更新 registry。
4. 如果本地 index 里没有可用版本元数据，才执行一次 registry update，然后重试解析。
5. 解析出版本以后，资产缓存和版本解析是分开的；资产 cache miss 不会再触发 registry update。

这意味着 unpinned 坐标的“latest”是基于本地 registry index 的 latest。为了可复查、可复现，正式文档或脚本里应优先 pin 版本：

```bash
moon runwasm owner/tool@1.2.3/cmd -- args...
```

而不是长期依赖：

```bash
moon runwasm owner/tool/cmd -- args...
```

## 缓存行为

远程资产会缓存到：

```text
$MOON_HOME/registry/cache/assets
```

本机实测后出现了：

```text
~/.moon/registry/cache/assets/Betterlol/moon_zod/0.5.1/cmd/json2schema/json2schema.wasm
~/.moon/registry/cache/assets/Yoorkin/cowsay/0.1.0/cowsay.wasm
~/.moon/registry/cache/assets/moonbitlang/parser/0.3.5/cmd/moonfmt/moonfmt.wasm
~/.moon/registry/cache/assets/moonbitlang/parser/0.3.5/cmd/wasm/mq/mq.wasm
```

缓存命中时直接运行缓存的 wasm，不重新下载。缓存 miss 时，`moon` 会先下载 `.wasm.sha256`，再下载 `.wasm`，计算本地 SHA-256，比对通过后用原子写入方式放入 cache。实测 `moonfmt.wasm` 的远端 sha256 和本地缓存文件对得上。

我没有在当前命令帮助或 `runwasm` 实现中看到自动垃圾回收策略。也就是说，它更像“下载到 Moon 自己的 cache 后长期复用”，不是“npx 一次性临时目录跑完自动删”。如果用户想清理，应该手动清理 `$MOON_HOME/registry/cache/assets` 下的内容，或者等以后工具链提供专门 cache clean 命令。当前不要对用户承诺“会自动清理”。

## 运行时能力

远程资产下载后由 `moonrun` 运行。当前 `moonrun` 是基于 V8 的 WebAssembly runner，并给 guest 提供了一组 host API。源码里能确认的能力包括：

1. argv：参数会传进 guest。`argv[0]` 是 wasm 文件路径，后面是用户在 `--` 后传入的参数。
2. exit code：guest 调用 exit 后，`moon runwasm` 会把退出码传回宿主进程。测试里有 `exit-7` 用例。
3. stdout/stderr：`println`、WASI `fd_write`、内部 IO host API 都能输出到宿主 stdout/stderr。
4. stdin：`moonrun` 有 `__moonbit_io_unstable.read_bytes_from_stdin`，也有 WASI `fd_read` 对 stdin 的实现。源码测试里有读取中文和 emoji stdin 后输出字节长度的用例。所以 `runwasm` 运行时并非不能读 stdin。之前 `jqx` 读不了管道，是 `jqx` 自己的 `stdin_wasm.mbt` fallback 返回空字符串，没有接这层 host API。
5. env：`moonrun` 会暴露宿主环境变量，也有读取、设置、删除环境变量的 host API。
6. filesystem：`moonrun` 实现了一批 WASI path/file API，并把当前工作目录作为 preopen `.` 暴露给 guest。实现里有路径边界检查，防止从 preopen 根逃逸到外部路径。这个能力适合 CLI 工具读写当前目录内文件，但不是任意系统访问。
7. current directory / preopen：WASI 侧 preopen 名称是 `.`，host path 是当前工作目录。

我没有看到 `moonrun` 给 guest 提供任意网络访问、执行宿主 shell 命令、动态安装外部依赖这类能力。它是受 host API 限制的 Wasm 运行，而不是普通 npm package 的安装脚本/Node 进程。

## 和 `npx` 的相似点

相似点主要在用户体验层：

1. 都能用一个命令运行远程生态里的工具。
2. 都支持“没显式安装到 PATH 也能运行”的使用方式。
3. 都有缓存，第二次运行更快。
4. 都支持把命令行参数转交给被运行程序。
5. 都适合工具型 CLI，比如 formatter、schema converter、codegen、lint/search 小工具。

如果只从用户感受讲，可以说它有一点“MoonBit 生态里的 npx 味道”。

## 和 `npx` 的关键不同

但工程边界差异很大：

1. `npx` 跑的是 npm package entry；`moon runwasm` 跑的是 Mooncakes prebuilt `.wasm` asset。
2. `npx` 可能涉及 npm dependency resolution、install scripts、JS wrapper、native optional deps；`moon runwasm` 不在用户机器上构建包，也不执行包安装脚本。
3. `npx` 可以从 npm registry 拉很多形态的包；`moon runwasm` 远程模式只认 Mooncakes 坐标和按规则存在的 wasm asset。
4. `npx` 的运行环境是 Node/系统进程；`moon runwasm` 的运行环境是 `moonrun` 提供的 Wasm host。
5. `npx` 对 package 的权限边界更接近普通本机程序；`moon runwasm` 的 guest 能力取决于 `moonrun` 暴露的 host API，天然更窄。
6. `npx` 常见心智是“临时安装再执行”；`moon runwasm` 更像“下载预构建 artifact，校验，缓存，执行”。

所以对用户宣传时可以类比 `npx`，但设计和文档里最好不要直接说“它就是 npx”。更严谨的说法是：

```text
moon runwasm provides an npx-like invocation experience for Mooncakes-hosted prebuilt Wasm CLI assets.
```

中文就是：

```text
moon runwasm 给 Mooncakes 上的预构建 Wasm CLI 资产提供了类似 npx 的即时调用体验。
```

## 发布侧回答

目前看，MoonBit 自家工具链已经有比较完整的自家分发路径，不必须靠 npm。

官方更新说明说 mooncakes.io 会自动为支持 Wasm 编译目标的 package 构建优化后的 wasm 文件，并展示同目录的 `SKILL.md`。这说明“自家 registry 的 CLI 工具分发”是 Mooncakes 负责，不是 npm 负责。

但作者仍然要让 package 符合可运行条件。一个普通 library 即使能 `moon add`，也不一定有 runwasm 资产。要能被 `runwasm` 用起来，至少需要：

1. 有 main package，也就是 `moon.pkg` 里 `"is-main": true`。
2. 支持 wasm 或 wasm-gc 目标，不要把 main 限死为 native。
3. Mooncakes 服务端构建成功，并产生对应 `.wasm` asset。
4. 如果想进入 skills 页面，目录里还应有 `SKILL.md`。

这里要注意，`moon publish` 命令本身的帮助只写“发布当前 module”，没有暴露一个“上传 wasm asset”的本地参数。结合官方更新说明，更合理的理解是：作者发布 MoonBit module 源码/包元数据到 Mooncakes，Mooncakes 服务端根据 package target 自动构建 wasm asset。也就是说，分发的权责是：

```text
作者：写可 wasm 编译的 main package，并发布到 Mooncakes。
Mooncakes：构建并托管 wasm asset 和 checksum。
用户：moon runwasm 坐标运行。
```

npm 可以作为平行分发渠道，但不是这条主链路。

## 对 JSON / jq 工具的含义

回到前面的 JSON 过滤工具，最合适的对象形态就是 Mooncakes runwasm CLI：

```bash
moon runwasm owner/json-filter@version -- '.foo' '{"foo":1}'
```

如果基于 `shina1024/jqx`，当前本地已经验证 wasm 参数输入能跑。它没有出现在 skills API，不是因为 `runwasm` 做不到，而是因为当前 registry 版本没有可命中的 prebuilt wasm asset，或者发布侧没有把 `cmd` 这个 main package 变成 runwasm 可用资产。

stdin 也类似。`runwasm` runtime 有 stdin 能力；`jqx` 当前 wasm 代码没接。要补的是 package 代码，不是换成 npm。

## 最终判断

`moon runwasm` 目前是一个已经可用但仍标为 experimental 的 MoonBit 自家 CLI 分发/运行通道。它的强项是：

1. 用户无需 npm、apt、brew、cargo。
2. 不需要本地构建。
3. 从 Mooncakes 拉预构建 wasm。
4. 校验 SHA-256。
5. 缓存到 `$MOON_HOME/registry/cache/assets` 后复用。
6. 参数、stdout/stderr、exit code、stdin、env、当前目录文件访问等 CLI 基础能力具备。

它的边界是：

1. 只跑本地 package 或 Mooncakes prebuilt wasm asset。
2. 不跑任意 URL。
3. 不跑 npm 包。
4. 不自动构建远程源码。
5. 不保证每个 Mooncakes module 都有 wasm asset。
6. 没看到自动清理缓存的机制。
7. unpinned 坐标依赖本地 registry index 的 latest 解析，不适合可复查脚本。

所以我的一句话判断是：

```text
它像 npx 的地方是“按坐标即时运行工具”；不像 npx 的地方是“它运行的是 Mooncakes 预构建 Wasm 资产，而不是临时安装一个完整包生态里的任意 package”。
```

对 seekmoon 这类工具实验来说，这反而是好消息：如果我们想给用户一个“不额外安装 jq”的 JSON 过滤能力，优先做 Mooncakes runwasm asset，而不是绕到 npm。

## 参考链接

- Moon Commands manual: https://moonbitlang.github.io/moon/commands.html
- MoonBit Updates, 2026-06-08 toolchain update: https://www.moonbitlang.com/updates/
- Mooncakes skills API: https://mooncakes.io/api/v0/skills
- Moon 源码 `runwasm.rs`: `external/moon-src/moon/crates/moon/src/cli/runwasm.rs`
- Moon 源码 runwasm reference: `external/moon-src/moon/docs/dev/reference/runwasm.md`
- Moon 源码 `moonrun`: `external/moon-src/moon/crates/moonrun/src/main.rs`
