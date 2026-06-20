# JS/TS 与 Rust 包发现和包管理机制研究报告

## 1. 范围和对象边界

本文只研究软件包复用生命周期中的两个相邻但不同的对象：包发现（package discovery）和包管理（package management）。包发现指消费者在 registry、搜索页、文档站、第三方索引或评分面上识别候选包、理解用途、维护状态、文档质量、下载量、依赖关系和风险信号的活动。包管理指工具在已有依赖声明之后执行解析、取回、校验、安装、缓存、锁定、发布、审计和权限控制的机制。二者都服务复用，但不能互相替代：搜索结果不能证明依赖树可复现，锁文件也不能证明包适合某个问题场景。

本文把 Manifest、Registry、Index、Resolver、Lockfile 和 Artifact 分开使用。Manifest 是包作者或项目作者写入的依赖和元数据表面；Registry 是发布、下载和查询包数据的服务；Index 是解析器可读取的包版本索引或元数据集合；Resolver 是从版本约束和来源约束中选定具体版本的算法或工具行为；Lockfile 是解析结果的持久化记录；Artifact 是实际被下载或安装的包内容，例如 npm tarball、Cargo `.crate` 文件、JSR 模块文件或 Yarn 缓存 zip。npm 文档明确把 `package.json` 中的 `name` 和 `version` 视为发布包时最重要且共同构成唯一标识的字段，Cargo 文档把 `Cargo.toml` 称为每个 package 的 manifest，并说明它包含编译所需元数据；这些都是本文区分 manifest 与管理产物的依据。来源：https://docs.npmjs.com/cli/v10/configuring-npm/package-json/ ，https://doc.rust-lang.org/cargo/reference/manifest.html

## 2. JS/TS 生命周期映射

JS/TS 生态的核心 canonical manifest 是 `package.json`。npm 文档规定 `package.json` 必须是实际 JSON；发布包时 `name` 与 `version` 是必需字段，且二者共同形成唯一标识。`description` 与 `keywords` 被 npm 文档明确说成会帮助 `npm search` 中的发现；`files` 控制安装为依赖时进入 tarball 的文件集合；`main`、`exports`、`bin` 等字段形成运行和命令暴露表面；`dependencies`、`devDependencies`、`peerDependencies`、`optionalDependencies` 等字段形成依赖声明表面。来源：https://docs.npmjs.com/cli/v10/configuring-npm/package-json/

npm registry 是包 metadata 与 artifact 的中心服务。npm 官方说明 public npm registry 是由 JavaScript package 及其 metadata 组成的数据库；npm CLI 默认使用 `https://registry.npmjs.org`，并通过 registry 按 name 和 version 解析包。npm 的 package spec 允许 `<name>`、`<name>@<tag>`、`<name>@<version>`、`<name>@<version range>`，也允许 tarball、git URL 和本地目录；npm 文档还说明 tarball 是上传到 registry 后包存在的格式。来源：https://docs.npmjs.com/about-the-public-npm-registry/ ，https://docs.npmjs.com/cli/v10/using-npm/registry/ ，https://docs.npmjs.com/cli/v10/using-npm/package-spec/

pnpm 和 Yarn 不是新的 manifest 体系，而是在 `package.json` 之上提供不同的 resolver、安装布局、缓存和锁定策略。pnpm 文档说明其安装分为依赖解析、目录结构计算、链接依赖三个阶段；它把包文件存入 content-addressable store，并从 store hard-link 到 `node_modules`，再通过 symlink 构造依赖图。Yarn 文档说明 Yarn manifest 仍是 `package.json`，包含 `name`、`version`、`dependencies`、`peerDependencies`、`resolutions` 等字段；现代 Yarn 默认安装策略是 Plug'n'Play，生成 `.pnp.cjs` loader，而不是传统 `node_modules`。来源：https://pnpm.io/motivation ，https://pnpm.io/symlinked-node-modules-structure ，https://yarnpkg.com/configuration/manifest ，https://yarnpkg.com/features/pnp

JSR 是 JS/TS 生态中与 npm 相邻但对象边界不同的 registry。JSR 官方称 JSR 是 JavaScript 和 TypeScript 的 package registry，并支持 Node.js、Deno、Bun、浏览器等运行环境，同时通过 npm compatibility 与 npm 工具链兼容。JSR package 以 JavaScript 或 TypeScript 编写，并以 ESM module 发布；JSR 要求包包含 config file 以读取 `exports` 等 metadata，错误文档中明确 `missingConfigFile` 与 `configFileExportsInvalid` 这类约束。JSR 的 package identity 采用 `@scope/name` 与 SemVer version，版本发布后不可变。来源：https://jsr.io/docs ，https://jsr.io/docs/publishing-packages ，https://jsr.io/docs/troubleshooting ，https://jsr.io/docs/packages ，https://jsr.io/docs/immutability

## 3. Rust 生命周期映射

Rust 生态的 canonical manifest 是 `Cargo.toml`。Cargo Book 规定每个 package 的 `Cargo.toml` 是 manifest，使用 TOML 格式，包含编译 package 所需 metadata；`[package]` 中 `name` 是 Cargo 唯一要求的字段，而发布到 registry 时 registry 会要求更多字段，例如 version、description、license 或 license-file。`keywords` 和 `categories` 是 registry 搜索与分类表面，Cargo 文档明确 `keywords` 可以帮助在 registry 上搜索 package。来源：https://doc.rust-lang.org/cargo/reference/manifest.html

Rust 的依赖声明表面集中在 `[dependencies]`、`[dev-dependencies]`、`[build-dependencies]` 和 target-specific dependencies。Cargo 默认从 crates.io 查找依赖，只需要 name 和 version requirement；版本 requirement 可以是默认 requirement、caret、tilde、wildcard、比较表达式和组合范围。Cargo 文档说明普通字符串 `"0.1.12"` 表示一个可解析的版本范围，并举例其含义是 `>=0.1.12, <0.2.0`。来源：https://doc.rust-lang.org/cargo/reference/specifying-dependencies.html

Rust 的 registry 与 index 关系比 npm 更显式。Cargo Book 说明 Cargo 从 registry 安装 crate 和获取依赖，默认 registry 是 crates.io；registry 包含一个 index，index 包含可搜索的可用 crate 列表，也可以提供 web API 支持 `cargo publish`。Cargo registry index 文档说明 index 根目录包含 `config.json`，其中 `dl` 是下载 `.crate` 文件的 URL 模板，`api` 是 web API base URL；index 其余部分为每个 package 一个文件，每个版本一行 JSON。来源：https://doc.rust-lang.org/cargo/reference/registries.html ，https://doc.rust-lang.org/cargo/reference/registry-index.html

Rust 的主要 artifact 是 `.crate` 文件。Cargo registry index 文档说明 download endpoint 应返回请求 package 的 `.crate` 文件；Cargo manifest 的 include/exclude 规则与 `cargo package --list` 用于验证发布包内文件。docs.rs 是 Rust 发现生命周期中的文档投影面：docs.rs 官方说明它是 Rust crates 的开源文档托管站，所有发布到 crates.io 的 libraries 都会被 documented；docs.rs build 文档说明它会自动为发布到 crates.io 的 crates 构建文档，并在 sandbox 中使用 nightly Rust compiler。来源：https://doc.rust-lang.org/cargo/reference/registry-index.html ，https://doc.rust-lang.org/cargo/reference/manifest.html ，https://docs.rs/about ，https://docs.rs/about/builds

## 4. 发现表面和数据来源

npm 的官方发现表面是 npmjs.com 与 `npm search` 所依赖的 registry metadata。npm 文档说明 npm search bar 使用 package title、description、README 和 keywords 执行搜索，搜索由 OpenSearch 支撑，结果基于这些字段的 keyword matching；排序可按 keyword matching、download counts、most dependents 和 last published date 等方式进行；deprecated packages 被排除在搜索结果之外，新发布包可能最多两周后才出现。来源：https://docs.npmjs.com/searching-for-and-choosing-packages-to-download/

JSR 的发现表面包括 jsr.io package page、自动文档、score、badges 和 API。JSR API 文档说明 JSR registry API 用于下载模块、package version metadata 和 package metadata；JSR scoring 文档把 documentation、best practices 和 discoverability 列为评分维度，其中 documentation 包括 README、module documentation 和 public functions/types documentation，discoverability 要求 package 有 description。JSR badges 文档提供 total downloads、weekly downloads 和 score badge URL。来源：https://jsr.io/docs/api ，https://jsr.io/docs/scoring ，https://jsr.io/docs/badges

Rust 官方发现表面包括 crates.io、Cargo index metadata 和 docs.rs。crates.io 首页说明其可以发布和安装 crates，并使用 API 查找可用 crates 信息；Cargo manifest 中的 `description`、`documentation`、`readme`、`homepage`、`repository`、`license`、`keywords`、`categories` 是可被 registry 和文档站使用的 metadata。docs.rs 提供 API reference 与构建状态投影；`[package.metadata.docs.rs]` 可控制 docs.rs build 的 features、all-features、targets 和 rustdoc args。来源：https://crates.io/ ，https://doc.rust-lang.org/cargo/reference/manifest.html ，https://docs.rs/about/metadata

lib.rs 是 Rust 生态重要第三方发现面，但不是 Cargo 官方 resolver 或 crates.io registry。本文只把它作为 discovery surface，而不是管理机制；可观察页面显示 lib.rs 按 crate 呈现描述、分类、依赖和生态信号，适合纳入“第三方发现指标覆盖度”变量，但其排序和评分逻辑需要后续单独验证，不能在本文中等同于 crates.io 官方数据。来源：https://lib.rs/

## 5. 管理机制：resolver、lockfile、cache/store、artifact integrity

npm install 的解析与安装输出为 `node_modules` tree 和 `package-lock.json`。npm install 文档说明如果存在 `npm-shrinkwrap.json`、`package-lock.json` 或 `yarn.lock`，依赖安装会按该顺序受这些文件驱动；`package-lock=false` 会忽略并阻止写入 lockfile。npm 的 install algorithm 文档给出 hoisting 示例，并明确该 algorithm 是 deterministic，但如果以不同顺序请求安装依赖，可能产生不同 tree。来源：https://docs.npmjs.com/cli/v10/commands/npm-install/

npm lockfile 是可复现机制的中心对象。npm package-lock 文档说明 `package-lock.json` 描述已生成的 exact tree，使后续 install 能生成 identical tree；lockfile 中 `resolved` 表示包实际解析来源，registry 来源时是 tarball URL；`integrity` 是被解包 artifact 的 sha512 或 sha1 Standard Subresource Integrity 字符串。来源：https://docs.npmjs.com/cli/v10/configuring-npm/package-lock-json/

pnpm 的管理机制强调 content-addressable store、严格依赖可见性和 lockfile 校验。pnpm 文档说明所有 package 文件存入单一 store，安装时 hard-link 到项目；默认只有直接依赖 symlink 到 root `node_modules`，这种布局限制访问未声明依赖。`pnpm install --frozen-lockfile` 在 lockfile 与 manifest 不同步或需要更新时失败；`--offline` 只使用本地 store；自 v11.4.0 起 tarball integrity mismatch 默认是 hard failure，并报 `ERR_PNPM_TARBALL_INTEGRITY`，`--update-checksums` 是显式刷新 integrity 的窄口。来源：https://pnpm.io/motivation ，https://pnpm.io/symlinked-node-modules-structure ，https://pnpm.io/cli/install

Yarn 的管理机制将安装拆为 Resolution、Fetch、Link、Build 四步。Yarn install 文档说明 Resolution 选择依赖版本，Fetch 下载依赖并存入 cache，Link 把 dependency tree 交给内部插件写入磁盘，例如生成 `.pnp.cjs`，Build 按拓扑顺序运行 build scripts。`--immutable` 会在 lockfile 将被修改时失败，CI 默认启用；`--immutable-cache` 在 cache folder 会被修改时失败；`--check-cache` 会重新获取包并校验 checksum 与 lockfile、cache 文件一致。来源：https://yarnpkg.com/cli/install

Cargo 的 resolver 根据每个 package 的版本 requirement 选择具体版本，并把结果写入 `Cargo.lock`。Cargo resolver 文档说明 dependency resolution 由 resolver 执行，结果存储在 `Cargo.lock` 以锁定具体版本并随时间保持固定；Cargo 通常偏好当前可用的最高版本，并在 SemVer compatible 范围内统一依赖版本，冲突时会回溯或报错。Cargo.toml vs Cargo.lock 文档明确 `Cargo.toml` 描述 broad dependencies，由用户编写；`Cargo.lock` 包含 exact dependency 信息，由 Cargo 维护，不应手工编辑。来源：https://doc.rust-lang.org/cargo/reference/resolver.html ，https://doc.rust-lang.org/cargo/guide/cargo-toml-vs-cargo-lock.html

Cargo registry index 在管理层提供 checksum 与 artifact 位置。index `config.json` 的 `dl` 字段定义 crate 下载 URL 模板，支持 `{crate}`、`{version}`、`{sha256-checksum}` 等 marker；index file 为每个 package 按版本记录 metadata。Cargo source replacement 文档还提供 registry mirroring、vendoring 等来源替换机制，适合离线和受控供应链场景，但替换源必须与原 source 表现为相同 package 集合，不能用于 patching。来源：https://doc.rust-lang.org/cargo/reference/registry-index.html ，https://doc.rust-lang.org/cargo/reference/source-replacement.html

## 6. 安全、来源证明和权限机制

npm 的安全机制覆盖 audit、registry signatures、provenance、trusted publishing、2FA 和 access tokens。`npm audit` 文档说明 audit 会把项目依赖描述提交到默认 registry 并请求已知漏洞报告；npm v7 起使用 Bulk Advisory endpoint，将 package name 和版本列表 POST 到 `/-/npm/v1/security/advisories/bulk`，并用 advisory objects 计算 vulnerabilities 和 meta-vulnerabilities。npm provenance 文档说明 provenance statement 可公开说明 package 在哪里构建、由谁发布；trusted publishing 文档说明受信发布使用 OIDC，可避免 CI 中长期 token，并可自动生成 provenance attestation。来源：https://docs.npmjs.com/cli/v10/commands/npm-audit/ ，https://docs.npmjs.com/generating-provenance-statements/ ，https://docs.npmjs.com/trusted-publishers/ ，https://docs.npmjs.com/viewing-package-provenance/

JSR 的安全机制以不可变版本、OIDC 发布、scope 权限和 provenance 为核心。JSR immutability 文档说明 package version 发布后不能改变；trust 文档说明从 GitHub Actions workflow 使用 `jsr publish` 或 `deno publish` 且使用 JSR + GitHub Actions 原生集成时，JSR 会自动创建 provenance statement；scope 文档说明 scope admin 可以要求所有 package versions 必须从 CI environment 发布，并使用 OIDC token。来源：https://jsr.io/docs/immutability ，https://jsr.io/docs/trust ，https://jsr.io/docs/scopes

pnpm 的管理安全主要体现在严格依赖可见性、store/cache 权限和 integrity failure。pnpm symlinked node_modules 文档说明只有真实在 dependencies 中的包可访问，flattened node_modules 则会使 hoisted packages 可访问；settings 文档说明 cacheDir 应只在互相信任的用户、作业和进程之间共享，并应使用文件系统权限保护；install 文档说明 tarball integrity mismatch 默认失败。来源：https://pnpm.io/symlinked-node-modules-structure ，https://pnpm.io/settings ，https://pnpm.io/cli/install

Yarn 的安全机制包括 audit、postinstall 控制、age gate 和 hardened mode。Yarn security 文档说明 Yarn 不在 install 时默认运行 audits，而是通过 `yarn npm audit` 执行；Yarn 4.14 起默认不运行 postinstall，除非全局启用 scripts 或在顶层 `package.json` 中逐包配置；Yarn 4.12 引入 `npmMinimalAgeGate`，默认要求包至少发布 1 天后才可安装；hardened mode 在公开 GitHub 仓库 PR 中默认启用，会自动打开 `--check-resolutions` 和 `--refresh-lockfile` 以对抗 lockfile poisoning。来源：https://yarnpkg.com/features/security ，https://yarnpkg.com/cli/npm/audit

Rust 官方 registry 安全正在形成 trusted publishing。crates.io trusted publishing 文档说明该机制允许从 GitHub Actions 和 GitLab CI/CD 等 CI/CD 平台发布 crate，而不需要手动管理 API tokens；Rust RFC 3691 说明 trusted publishing 通过 OIDC 发放 short-lived access tokens 以认证和授权 crates.io API actions。RustSec 是重要的社区安全数据面，而不是 Cargo 内建功能；RustSec 官方页说明它是 Rust ecosystem vulnerability database，`cargo-audit` 用于 audit `Cargo.lock` 中存在安全漏洞的 crates。来源：https://crates.io/docs/trusted-publishing ，https://rust-lang.github.io/rfcs/3691-trusted-publishing-cratesio.html ，https://rustsec.org/

## 7. 统一模型可用变量和指标

可用于跨生态评分的变量应按对象层位拆分，而不是把“生态好坏”写成单一印象分。包身份变量包括：identity tuple 是否清楚、scope/namespace 是否存在、version 是否 SemVer、name 规则是否可机器校验、同名冲突或保留名规则是否公开。Manifest 变量包括：依赖分组粒度、入口点声明、发布文件白名单/黑名单、license 表达、repository/documentation/readme 字段、workspace 支持、环境约束字段。

Registry/Index 变量包括：metadata API 是否公开、index 是否可镜像、artifact URL 是否可由 index 推导、版本 metadata 是否包含 checksum、是否支持私有 registry 或 alternate registry、是否支持 scoped registry、是否有 immutable version 规则。Discovery 变量包括：搜索字段覆盖范围、排序维度、下载量窗口、dependent count、文档生成状态、README 渲染、category/keyword、deprecated/archived 状态、provenance badge、score badge、最近发布时间、维护者/owner 显示、第三方索引覆盖。

Resolver/Lockfile 变量包括：版本选择策略、peer dependency 处理、SemVer 兼容解释、是否偏好最高版本、是否可配置最低/时间门槛、lockfile 是否记录 artifact source、checksum、完整 transitive graph、workspace 信息、是否支持 frozen/immutable install、lockfile 与 manifest 不一致时是否失败。Cache/Store 变量包括：content-addressed 程度、offline install 能力、cache 可校验性、cache 是否可提交、store 是否跨项目共享、cache 权限要求。Security 变量包括：audit 数据源、audit 默认执行时机、provenance 支持、OIDC trusted publishing、2FA/token scope、postinstall 默认策略、age gate、hardened lockfile verification、registry signature 或 attestation verification、malware report 通道。

## 8. 跨生态比较：JS/TS vs Rust

JS/TS 生态的强项是工具多样性和发现表面密度。npm 提供大规模 registry、搜索、download/dependent 排序、audit endpoint 和 provenance；pnpm 提供 content-addressable store 与严格依赖布局；Yarn 提供 PnP、zero-install、hardened mode 和 age gate；JSR 对 TypeScript、自动文档、score 和 provenance 进行更强约束。弱点是管理机制碎片化：同一 `package.json` 可被 npm、pnpm、Yarn、Bun、Deno 或 JSR compatibility 层解释，`node_modules`、PnP、symlink store、npm compatibility registry 的行为差异需要在模型中作为变量，而不能用“JS 包管理”一概而论。

Rust 生态的强项是 Cargo、Cargo.toml、Cargo.lock、crates.io index、`.crate` artifact 和 docs.rs 之间对象边界稳定。Cargo resolver 与 lockfile 关系清楚，registry index 可作为解析器数据源，docs.rs 自动把 crates.io release 投影成 API documentation。弱点是发现和安全评分的官方表面较少：crates.io 与 docs.rs 提供基础 metadata、下载和文档投影，但 RustSec、cargo-audit、cargo-vet、lib.rs 等重要判断面多在社区或第三方工具中；它们对实际复用判断有价值，但在统一模型里必须标为“非 Cargo 内建管理机制”或“第三方发现/审计面”。

从 lifecycle stage 看，JS/TS 在 discovery、registry metadata、发布来源证明和多工具安装策略上表面丰富，但 resolver 与安装布局存在工具分叉；Rust 在 manifest、resolver、lockfile、artifact index 和 documentation hosting 上更集中，但安全审计和质量评分更多依赖外部生态。跨生态模型不能把“官方集中度”直接等同于“风险低”，也不能把“工具多样性”直接等同于“能力强”；可测变量必须落在 identity、metadata、resolution、reproducibility、artifact integrity、provenance、permissions 和 discovery signals 上。

## 9. 开放不确定性和后续需核查来源

第一，npm registry 的底层 CouchDB mirror 与当前 npm public registry API 的完整字段契约需要后续使用 npm registry API 文档或实际 endpoint 样本核查；本文只依据 npm CLI registry 文档确认默认 registry、scope registry 和 CouchDB public mirror 的存在。第二，pnpm lockfile 格式细节目前官方用户文档覆盖有限，完整结构可能需要以后查阅 pnpm/spec lockfile 文档；本文只使用 pnpm install、motivation、settings 和 node_modules structure 文档确认 resolver、store、integrity 与 frozen lockfile 行为。第三，Yarn lockfile checksum 的精确格式和 PnP loader 内部 schema 需要后续阅读 Yarn specification 或源码文档；本文只依据 Yarn install 和 PnP 文档确认 install steps、immutable/cache/check-cache 与 `.pnp.cjs` 行为。

第四，crates.io API 和 data access policy 在本次访问中部分页面返回限制信息，后续应在遵守 crates.io data access policy 的前提下使用官方 API 文档或低频请求验证字段。第五，lib.rs 的排序、分类和评分规则需要后续读取 lib.rs 公共文档或观察多个页面，不应把单一页面观察扩展为完整机制声明。第六，Rust trusted publishing 已有 crates.io 官方文档和 Rust RFC，但其支持的 CI provider、UI 约束和 token 权限细节会随 crates.io 运行策略变化，后续评分前应重新核对官方页面。

## 10. 参考资料

- npm package.json: https://docs.npmjs.com/cli/v10/configuring-npm/package-json/
- npm package-lock.json: https://docs.npmjs.com/cli/v10/configuring-npm/package-lock-json/
- npm install: https://docs.npmjs.com/cli/v10/commands/npm-install/
- npm registry: https://docs.npmjs.com/cli/v10/using-npm/registry/
- npm package spec: https://docs.npmjs.com/cli/v10/using-npm/package-spec/
- npm public registry: https://docs.npmjs.com/about-the-public-npm-registry/
- npm package discovery: https://docs.npmjs.com/searching-for-and-choosing-packages-to-download/
- npm audit: https://docs.npmjs.com/cli/v10/commands/npm-audit/
- npm provenance: https://docs.npmjs.com/generating-provenance-statements/
- npm trusted publishers: https://docs.npmjs.com/trusted-publishers/
- npm provenance viewing: https://docs.npmjs.com/viewing-package-provenance/
- pnpm motivation/store/install stages: https://pnpm.io/motivation
- pnpm node_modules layout: https://pnpm.io/symlinked-node-modules-structure
- pnpm install: https://pnpm.io/cli/install
- pnpm settings: https://pnpm.io/settings
- Yarn manifest: https://yarnpkg.com/configuration/manifest
- Yarn Plug'n'Play: https://yarnpkg.com/features/pnp
- Yarn install: https://yarnpkg.com/cli/install
- Yarn security: https://yarnpkg.com/features/security
- Yarn npm audit: https://yarnpkg.com/cli/npm/audit
- JSR introduction: https://jsr.io/docs
- JSR publishing packages: https://jsr.io/docs/publishing-packages
- JSR packages: https://jsr.io/docs/packages
- JSR API: https://jsr.io/docs/api
- JSR scoring: https://jsr.io/docs/scoring
- JSR npm compatibility: https://jsr.io/docs/npm-compatibility
- JSR provenance and trust: https://jsr.io/docs/trust
- JSR scopes: https://jsr.io/docs/scopes
- JSR immutability: https://jsr.io/docs/immutability
- JSR badges: https://jsr.io/docs/badges
- JSR troubleshooting: https://jsr.io/docs/troubleshooting
- Cargo manifest: https://doc.rust-lang.org/cargo/reference/manifest.html
- Cargo specifying dependencies: https://doc.rust-lang.org/cargo/reference/specifying-dependencies.html
- Cargo resolver: https://doc.rust-lang.org/cargo/reference/resolver.html
- Cargo.toml vs Cargo.lock: https://doc.rust-lang.org/cargo/guide/cargo-toml-vs-cargo-lock.html
- Cargo registries: https://doc.rust-lang.org/cargo/reference/registries.html
- Cargo registry index: https://doc.rust-lang.org/cargo/reference/registry-index.html
- Cargo source replacement: https://doc.rust-lang.org/cargo/reference/source-replacement.html
- crates.io: https://crates.io/
- crates.io trusted publishing: https://crates.io/docs/trusted-publishing
- Rust RFC 3691 trusted publishing for crates.io: https://rust-lang.github.io/rfcs/3691-trusted-publishing-cratesio.html
- docs.rs about: https://docs.rs/about
- docs.rs builds: https://docs.rs/about/builds
- docs.rs metadata: https://docs.rs/about/metadata
- RustSec: https://rustsec.org/
- lib.rs: https://lib.rs/
