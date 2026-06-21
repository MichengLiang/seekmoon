# 纯 MoonBit CLI 工具分发路径报告

调查日期：2026-06-21

## 结论

纯 MoonBit CLI 工具目前至少有四条分发路径：

1. `moon runwasm owner/module[/pkg][@version] -- args...`
2. `moon install owner/module/pkg[@version]`
3. npm 包分发
4. GitHub Release / 手工下载 native binary

它们不是互相替代的同一种东西。最优雅的路径取决于消费者动作。

如果用户只是“临时用一下工具，不想安装，不想管平台”，最优雅的是 `moon runwasm`。它像 npx 的地方是按坐标即时运行；不像 npx 的地方是它跑 Mooncakes 上的预构建 Wasm asset，不跑 npm package，不在用户机器上构建源码。

如果用户要“长期装到 PATH 里，当普通命令反复用”，MoonBit 自家路径是 `moon install`。它会从 registry / git / local path 选中 `is-main: true` 的 package，构建 native release 可执行文件，并安装到 `~/.moon/bin`。这条路径更像 `cargo install` 或 `go install`，不是 npx。

如果用户群已经在 JS/TS / Node 生态里，npm 依然是非常友好的分发渠道。它拥有成熟的 discovery、versioning、lockfile、CI、企业代理、缓存、跨平台 wrapper、`npx`/`pnpm dlx` 即时运行体验。对“开发工具”而言，npm 不是低级方案，反而是目前全球最顺手的工具分发基础设施之一。

所以我的推荐不是单选，而是分层：

```text
MoonBit 原生优先面：Mooncakes + runwasm + moon install
广域开发者面：npm
传统离线/平台面：GitHub Release native binaries
```

若只问“对 MoonBit 生态最优雅的主路径”，答案是：

```text
发布到 Mooncakes，提供 runwasm 即时运行；同时让同一个 main package 可 moon install 成 native CLI。
```

若问“对普通开发者最友好的公开分发”，答案通常是：

```text
再补一个 npm 包，提供 npx/pnpm dlx 入口。
```

## 路径一：`moon runwasm`

`moon runwasm` 的定位是预构建 Wasm CLI 资产运行器。用户不先安装工具，而是直接运行：

```bash
moon runwasm moonbitlang/parser/cmd/moonfmt -- --help
moon runwasm Betterlol/moon_zod@0.5.1/cmd/json2schema -- '{"foo":1}'
```

它的优点很清楚：

1. 用户只需要 MoonBit 工具链，不需要 npm、cargo、apt、brew。
2. 不需要在用户机器上编译。
3. 跨平台一致性比 native binary 更好。
4. Mooncakes 提供 `.wasm` 和 `.wasm.sha256`，`moon` 下载、校验并缓存。
5. 适合 formatter、schema converter、json filter、codegen、agent helper 等轻量 CLI。
6. 如果配 `SKILL.md`，还可以进入 MoonBit skill / agent 工具叙事。

它的短板也要正视：

1. 当前仍是 experimental 心智。
2. 不是每个 Mooncakes 包都有可运行 wasm asset。
3. 运行时能力受 `moonrun` host API 限制。
4. 缓存目前没有看到自动 GC。
5. unpinned 版本依赖本地 registry index 的 latest 解析，可复查脚本应 pin 版本。
6. 用户必须安装 MoonBit 工具链；对非 MoonBit 用户来说，这是额外门槛。

所以它非常适合 MoonBit 内部生态和 agent/tool 场景，但还不能替代 npm 在大众开发者工具分发里的地位。

## 路径二：`moon install`

`moon install` 的 binary installer 模式当前已经比较明确。命令帮助写的是：

```text
Install a binary package globally or install project dependencies (deprecated without args)
```

它支持：

```bash
moon install user/module/pkg
moon install user/module/pkg@1.2.3
moon install user/module/...
moon install ./local/pkg
moon install <git-url> [PATH_IN_REPO] --tag v1.2.3
```

它的语义是：选择一个 source，再选择其中的 `is-main: true` package，构建 native release binary，安装到默认 `~/.moon/bin`，也可以用 `--bin <DIR>` 指定目录。

实测 dry-run：

```bash
moon install --dry-run mizchi/jq/cmd/moonjq@0.2.2
```

输出显示会构建并安装：

```text
Dry-run: Would build `mizchi/jq/cmd/moonjq`
Dry-run: Would install `moonjq` to `/home/t103o/.moon/bin/moonjq`
```

这条路径对长期使用者更自然。装完以后用户可以直接：

```bash
moonjq '.foo'
```

不用每次写 `moon runwasm owner/module/pkg -- ...`。

但它的成本也更高：

1. 它在用户机器上构建 native binary。
2. 用户需要完整 MoonBit 工具链和 native 编译链。
3. 跨平台问题更直接，尤其 Windows native stub、C compiler、系统库等。
4. 安装结果进入 `~/.moon/bin`，有 PATH、命名冲突和升级管理问题。
5. 包路径命名会直接影响二进制名。

`shina1024/jqx` 就暴露了一个命名问题。实测：

```bash
moon install --dry-run shina1024/jqx/cmd@0.2.0
```

会安装到：

```text
~/.moon/bin/cmd
```

因为 binary name 默认取 package path 最后一段。这个名字对用户显然不友好。若要把它作为 MoonBit CLI 分发，最好把 package path 设计成 `cmd/jqx` 或 `jqx` 这类可安装名，而不是裸 `cmd`。

所以 `moon install` 很适合“我已经是 MoonBit 用户，我愿意长期安装一个 CLI”。它不适合“外部用户临时试一下”。

## 路径三：npm

npm 依然非常友好，尤其对开发工具而言。

原因不是 npm 技术上更“纯”，而是它的分发基础设施太成熟：

1. 用户基数大，Node/npm/pnpm/bun 已经在开发者机器上普遍存在。
2. `npx tool`、`pnpm dlx tool`、`bunx tool` 的心智非常成熟。
3. package discovery、README、版本、semver、deprecation、provenance、2FA、组织权限都成熟。
4. CI、企业镜像、缓存、lockfile、代理支持都成熟。
5. JS wrapper 可以统一包装不同平台 native binary、wasm、JS glue。
6. 对 VS Code、Web、Node 脚本、AI agent 工具链集成很顺。

对纯 MoonBit CLI 来说，npm 可以有几种形态：

1. 发布 JS target 编译产物。
2. 发布 Wasm + JS wrapper。
3. 发布各平台 native binary，并由 npm postinstall 或运行时 wrapper 选择。
4. 发布一个很薄的 bin wrapper，内部调用 `moon runwasm`。
5. 发布 TS/JS library API，同时附带 CLI。

`shina1024/jqx` 已经采用了多面分发：Mooncakes module、GitHub Release native CLI、npm packages。它的 workflow 里有 `release-cli.yml` 构建 Linux/macOS/Windows native binary，也有 `release-npm.yml` 发布 npm packages。这是一个现实信号：作者如果想服务 MoonBit 之外的用户，npm 和 GitHub Release 仍然有价值。

npm 的不足也很明确：

1. 它不是 MoonBit 原生注册区。
2. 用户如果没有 Node 工具链，就多一个依赖。
3. 如果用 postinstall 下载 native binary，会引入供应链和企业环境阻力。
4. 如果只是 wrapper 到 `moon runwasm`，那用户还得有 MoonBit 工具链，npm 只提供入口，不提供能力本身。

所以 npm 不是 MoonBit 原生最优雅路径，但它可能是大众开发者最友好的路径。

## 路径四：GitHub Release native binary

这是传统 CLI 分发路径：

```text
jqx-v0.2.0-linux.tar.gz
jqx-v0.2.0-macos.tar.gz
jqx-v0.2.0-windows.zip
SHA256SUMS
```

优点：

1. 用户不需要 MoonBit 工具链。
2. 启动速度和性能通常最好。
3. 可以进入 Homebrew、Scoop、Winget、AUR、Nix 等后续渠道。
4. 对企业离线环境友好。

缺点：

1. 发布矩阵复杂。
2. 每个平台都要构建、签名、校验。
3. 用户安装体验不如 `npx` / `moon runwasm` 一条命令。
4. 后续升级、PATH、权限、杀软、系统依赖都要处理。

这条路适合成熟 CLI，不适合早期实验作为唯一主路径。

## 当前 MoonBit CLI 分发生态友好吗？

我的判断是：对 MoonBit 圈内用户，已经“可用且正在变友好”；对广域开发者，还没有 npm 那么友好。

友好的部分：

1. `moon runwasm` 已经能做到一条命令运行 Mooncakes 上的预构建 wasm。
2. Mooncakes skills API 已经有 70 个 wasm entry。
3. `moon install` 已经有明确 binary installer 语义，支持 registry / git / local。
4. Mooncakes 会为支持 Wasm target 的 package 构建 wasm asset。
5. 同一个 MoonBit main package 可以服务 runwasm 和 install 两个场景。

不够友好的部分：

1. 用户还不一定知道哪个 package 可 runwasm。
2. 普通 module 与可执行 skill / wasm asset 的发现界面还不如 npm 成熟。
3. 包作者需要理解 package path、`is-main`、target、SKILL.md、asset 命名等细节。
4. `moon install` 的二进制命名依赖 package path，容易出现 `cmd` 这种差体验。
5. `runwasm` experimental 心智还会让人不确定是否适合长期依赖。
6. 对非 MoonBit 用户，安装 MoonBit 工具链仍然是额外门槛。

所以，MoonBit 原生 CLI 分发已经有骨架，而且方向是对的，但生态习惯、发现入口、命名规范、发布模板还需要沉淀。

## 最优雅的分发组合

如果我是给一个纯 MoonBit CLI 工具设计分发，我会用三层结构。

第一层，MoonBit 原生即时运行：

```bash
moon runwasm owner/tool@1.2.3/cmd/tool -- args...
```

这是 MoonBit 生态里最优雅的“试用”和“agent 调用”入口。文档里应 pin version，保证可复查。

第二层，MoonBit 原生长期安装：

```bash
moon install owner/tool/cmd/tool@1.2.3
tool args...
```

这是 MoonBit 用户把工具装进 PATH 的入口。前提是 package path 的最后一段必须是好名字，例如 `cmd/tool`，不要叫 `cmd` 或 `main`。

第三层，广域开发者入口：

```bash
npx @owner/tool args...
pnpm dlx @owner/tool args...
```

或者：

```bash
npm install -g @owner/tool
tool args...
```

这服务不想安装 MoonBit 工具链的人。npm 包可以包装 JS target、Wasm + JS runner、native binary，或者在特定情况下包装 `moon runwasm`。但如果 wrapper 仍要求 MoonBit，文档必须明说。

成熟后再补 GitHub Release native binary，服务离线、CI、性能敏感和非 Node 用户。

## 对 JSON 过滤工具的具体建议

如果继续做类似 jq 的工具，我建议：

1. 以 Mooncakes + `runwasm` 作为 MoonBit 原生主入口。
2. package path 改成或封装成 `cmd/jqx`，避免 `moon install` 后叫 `cmd`。
3. `moon runwasm` 文档先承诺参数 JSON：

```bash
moon runwasm owner/jqx@0.2.1/cmd/jqx -- '.foo' '{"foo":1}'
```

4. `moon install` 文档提供长期安装：

```bash
moon install owner/jqx/cmd/jqx@0.2.1
jqx '.foo' '{"foo":1}'
```

5. npm 继续作为广域入口：

```bash
npx @owner/jqx '.foo' '{"foo":1}'
```

6. 如果补 stdin，优先接 `moonrun` 已有的 stdin host API，而不是为了 stdin 改走 npm。

这套组合中，Mooncakes 是原生能力中心，npm 是传播半径扩展器，GitHub Release 是成熟工具的保底分发层。

## 一句话判断

目前纯 MoonBit CLI 的原生分发已经有两条正经路径：`moon runwasm` 负责“不安装即时运行”，`moon install` 负责“构建并安装 native CLI”。它们对 MoonBit 用户是越来越友好的，但还没有 npm 那种全生态成熟度。

最优雅的 MoonBit 原生方案是：

```text
同一个 is-main package，同时支持 Mooncakes runwasm asset 和 moon install native binary。
```

最友好的公开开发者方案是：

```text
在 MoonBit 原生分发之外，再提供 npm bin 入口。
```

所以不是“Mooncakes vs npm”二选一。更合理的是：Mooncakes 定义 MoonBit 原生分发，npm 承担跨生态触达。

## 参考

- `moon help runwasm`
- `moon help install`
- `external/moon-src/moon/docs/dev/reference/runwasm.md`
- `external/moon-src/moon/docs/dev/reference/moon-install-binary.md`
- `external/moon-src/moon/crates/moon/src/cli/runwasm.rs`
- `external/moon-src/moon/crates/moonrun/src/main.rs`
- `shina1024/jqx` release workflows
