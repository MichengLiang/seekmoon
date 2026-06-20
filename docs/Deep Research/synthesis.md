# Package Discovery 与 Package Management 深度研究综合稿

## 1. 研究对象

当前研究对象是主流编程语言软件复用生态中的 Package Discovery 与 Package Management 机制。

Package Discovery 指消费者在引入依赖之前识别候选包、理解候选包、比较候选包并形成引入判断的机制。搜索是 Discovery 的入口之一，不等于 Discovery。Discovery 还包括分类浏览、标签、相似包、文档阅读、API 检查、依赖图观察、安全状态、维护状态、生态采用度和社区策展。

Package Management 指消费者项目声明、解析、获取、校验、安装、构建、锁定、升级、审计和移除依赖的机制。Package Manager 是这个机制的工具投影，例如 npm、pnpm、Yarn、Cargo、pip、uv、Maven、Gradle、NuGet、Conan、vcpkg、SwiftPM。

这两个机制共享同一条包复用生命周期，但职责不同。Discovery 处理引入前的判断，Management 处理引入后的工程成立。

## 2. 已完成调查范围

本轮调查拆成六个方向，均已生成独立报告：

- JS/TS 与 Rust：`reports/01-js-ts-rust-report.md`
- Python 与 Go：`reports/02-python-go-report.md`
- JVM/Kotlin 与 .NET：`reports/03-jvm-dotnet-report.md`
- C/C++ 与 Swift：`reports/04-cpp-swift-report.md`
- 安全、出处证明与治理：`reports/05-security-provenance-governance-report.md`
- 跨生态发现数据、图谱与索引表面：`reports/06-cross-ecosystem-discovery-data-report.md`

任务书保存在 `tasks/`。报告目录当前被仓库规则忽略，但文件已存在。

本轮使用的资料类型以官方文档、官方规范、语言项目文档、注册中心文档、安全基础设施文档为主；第三方发现页面只作为 Discovery surface 证据，不作为注册表事实源。

## 3. 生命周期模型

包复用生态可以建模为九个阶段。

1. Authoring：作者写源码、测试、文档、示例和构建脚本。
2. Declaration：作者通过 Manifest 声明身份、版本、依赖、许可证、平台、入口、构建后端、特性开关或目标环境。
3. Publication：作者把源码包、二进制构件或索引记录发布到 Registry、Repository 或 VCS-backed index。
4. Indexing：平台抽取元数据、依赖关系、版本、README、API 文档、下载数据、安全状态、兼容性状态。
5. Discovery：消费者通过搜索、分类、标签、文档、图谱、推荐、社区策展和 IDE/CLI 查询得到候选包。
6. Assessment：消费者判断候选包是否满足当前目标和约束。
7. Resolution：Package Manager 根据消费者项目声明、目标平台、版本约束、lockfile 和 Registry metadata 求出一致依赖图。
8. Realization：系统下载、校验、缓存、构建、链接、写入 lock，使候选依赖成为项目中的实际构件。
9. Maintenance：系统处理升级、漏洞、弃用、许可证变化、迁移、供应链风险和依赖图变化。

消费者项目发布为上层包时，生命周期回到 Authoring / Declaration。包生态不是单向管道，而是循环网络。

## 4. 核心对象

Package 是在某个语言或构建生态中被命名、版本化、可获取、可声明为依赖的软件复用单元。Package 可以是库、框架、CLI 工具、插件、类型声明、源码集合、二进制构件集合或系统库绑定。

Artifact 是可下载、可缓存、可校验、可构建或可链接的具体构件。一个 package 可以对应多个 artifact，例如不同平台、不同 ABI、不同 target framework、不同 feature 组合。

Manifest 是项目或包对机器声明自身身份、依赖、版本、平台、许可证、入口和构建约束的结构化表面。Manifest 让包成为可索引、可求解、可审计、可复现的对象。典型形式包括 `package.json`、`Cargo.toml`、`pyproject.toml`、`pom.xml`、`.nuspec`、`go.mod`、`Package.swift`、`conanfile.py`、`vcpkg.json`。

Registry 是接收、保存、分发和索引 package metadata 与 artifact 的基础设施。Index 是可搜索、可解析、可用于发现或依赖求解的数据视图。Repository 可以是源码仓库，也可以是构件仓库；它不必然等同于 Registry。

Resolver 是把版本范围、依赖关系、平台约束、features、extras、target frameworks、lockfile 和 registry metadata 合成一致依赖图的求解器。Resolver 不是下载器。

Lockfile 是一次依赖求解结果的固化记录。它服务协作一致性、可复现构建和后续审计。

Provenance 描述 artifact 从什么源码、什么构建过程、什么发布身份产生。Attestation 是对这些事实的可验证声明。它们属于供应链信任机制。

## 5. 证据源分类

评价模型应区分六类证据源。

Declared Evidence 来自 Manifest 和作者声明。它表达作者主张。字段包括 name、version、description、keywords、license、dependencies、repository、engines、platforms、targets、features。

Structural Evidence 来自源码、AST、类型签名、导出符号、文档注释、测试结构、示例、deprecated 标记。它表达包本体结构。

Graph Evidence 来自依赖图和被依赖图。字段包括 direct dependencies、transitive dependencies、reverse dependencies、依赖深度、重复版本、生态中心性。

Operational Evidence 来自 CI、构建矩阵、平台二进制、wheel/tag、target framework、ABI、triplet、checksum、构建结果。它表达包能否在目标环境中成立。

Social Evidence 来自下载量、stars、forks、issues、PR、讨论、Awesome list、企业案例。它表达注意力、经验和采用情况，但容易被历史惯性和刷量污染。

Security and Governance Evidence 来自 CVE/GHSA/OSV、OpenSSF Scorecard、SLSA、Sigstore、trusted publishing、provenance、2FA、owner policy、撤回规则、命名空间策略。它表达信任边界。

## 6. 变量、指标、权重、分数

维度是评价空间中的抽象方向，例如安全、维护、兼容、文档、依赖图、发现性。

变量是某个维度下可取值的观察对象，例如最近 release 距今天数、直接依赖数量、critical 漏洞数量、支持的 target framework 数量。

指标是变量的具体测量口径，例如 `days_since_last_release`、`direct_dependency_count`、`critical_vulnerability_count`、`reverse_dependency_count_log`。

权重是某个指标或维度在给定使用场景中的重要程度。企业生产、个人原型、嵌入式、移动端、科研脚本、基础设施库的权重不同。

分数是指标归一化和加权后的评价输出。分数不是事实本身，而是模型投影。

## 7. 评价函数

包级评价函数：

```text
Score_package(package, query, consumer_context) -> ranked suitability
```

该函数评价某个候选包对某个查询和上下文的适配度。输入包括候选包、消费者功能意图、语言、运行时、平台、许可证策略、安全阈值、生产/原型场景、组织策略。输出是可解释排序，不是世界唯一最优。

生态级评价函数：

```text
Score_ecosystem(ecosystem, lifecycle_stage, criterion_set) -> capability profile
```

该函数评价某个语言生态在生命周期某阶段的机制成熟度。例如 Rust 的文档投影机制、Python 的二进制兼容机制、Java/Kotlin 的 variant metadata、Go 的 checksum/proxy 机制、C++ 的 ABI/平台变体管理、JS/TS 的供应链治理。

这两类评价不能混写。Package-level 关注“这个包是否适合当前任务”，Ecosystem-level 关注“这个生态是否提供了让包复用成立的机制”。

## 8. 主流生态结论

JS/TS 的 manifest 和 registry 表面极丰富。`package.json` 承载身份、入口、依赖、环境、脚本和发现字段。npm 提供大规模 registry、搜索、audit、provenance 和 trusted publishing。pnpm 强调 content-addressable store 与严格依赖布局。Yarn 强调 Plug'n'Play、immutable install、hardened mode 等。JSR 强调 TypeScript-first、score、文档、跨运行时和 provenance。JS/TS 的主要风险是生态规模带来的噪声、供应链风险、postinstall/脚本风险、工具分叉和依赖树膨胀。

Rust 的 manifest、registry、resolver、lockfile 和文档投影关系清楚。`Cargo.toml` 声明 package 与 dependencies，Cargo resolver 生成 `Cargo.lock`，crates.io index 支持解析，`.crate` 是 artifact，docs.rs 自动为 crates.io crate 构建文档。Rust 的官方链路集中，第三方发现和安全审计补充明显，包括 lib.rs、RustSec、cargo-audit 等。

Python 的核心事实是 distribution package 与 import package 不同。`pyproject.toml` 与 Core Metadata 逐步提供静态声明表面，Simple Repository API 提供索引和文件级数据，wheel tag 表达 Python 实现、ABI 与平台兼容性。pip 使用回溯解析，uv 使用跨平台锁定和更明确的索引优先策略。Python 的难点是历史动态元数据、二进制 wheel 兼容、多个索引、环境 marker、extras 和供应链风险。

Go 的模式不是传统上传式 registry。`go.mod` 声明 module path、require、replace、exclude、retract。Go 使用 Minimal Version Selection，build list 不保存为普通 lockfile，`go.sum` 与 checksum database 提供内容认证。pkg.go.dev 从 module proxy 和源码生成文档、搜索、符号、imported-by 和漏洞投影。Go 的描述性 manifest 字段少，但 module path、源码文档、checksum database 和 pkg.go.dev 形成另一种稳定模式。

JVM/Kotlin 的身份结构是坐标化的。Maven 使用 groupId、artifactId、version，POM 同时是构建描述和依赖元数据。Gradle 在 Maven/Ivy 之上引入 variant-aware resolution、attributes、capabilities 和 Gradle Module Metadata。Kotlin Multiplatform 要求 target-aware metadata，因为同一 library 可能发布 JVM、JS、Native、Android 等变体。该生态强在企业治理和变体表达，复杂性来自构建系统与多平台元数据。

.NET/NuGet 使用 package ID + version 标识包，`.nuspec` 和 SDK project properties 提供 metadata，PackageReference 进入 MSBuild 项目，target framework 是兼容性核心变量。NuGet 有 lock file、contentHash、locked mode、restore audit 和 vulnerability info API。NuGet 的强项是 IDE/CLI 集成、TFM 过滤和恢复期审计。

C/C++ 的核心差异是二进制身份不是名称+版本。Conan 的 package_id 由 settings、options、依赖版本等配置参与计算；settings 包括 OS、arch、compiler、libcxx、cppstd、runtime、build_type 等。vcpkg 使用 manifest、baseline、overrides、registry、triplet 和 binary caching。C/C++ 的评价模型必须把 ABI、compiler、platform、build type、link mode、feature/options 作为一等变量。

Swift 使用 `Package.swift` 声明 package、products、targets、dependencies、platforms 和 tools version。SwiftPM 负责解析并写入 `Package.resolved`。Swift Package Registry Service 提供 registry protocol，Swift Package Index 是第三方发现和兼容性评估表面，尤其提供 Swift/compiler/platform 构建矩阵。Swift 的主要约束来自平台、tools version、products/targets 和源码/二进制 target。

## 9. 安全与治理结论

安全不能压成一个单一分数。应拆分为漏洞状态、发布者身份、构建来源证明、制品完整性和项目治理。

漏洞状态描述某包版本是否被 OSV、GHSA、NVD、RustSec、Go vuln DB、NuGet audit 等公告覆盖。恶意包 advisory 与普通漏洞 advisory 的处理语义不同。

发布者身份描述谁被 registry 接受为发布主体。npm trusted publishing 和 PyPI Trusted Publishing 使用 OIDC 减少长期 token 风险。

构建来源证明描述 artifact 从什么源码、什么构建过程、什么构建平台产生。SLSA 和 Sigstore 提供统一的 provenance / attestation 语言和签名透明日志机制。

制品完整性描述下载内容是否匹配预期摘要、签名、attestation 或透明日志材料。该机制通常是文件级或 artifact 级，不等同于项目安全。

项目治理描述命名空间、所有权、维护状态、权限、2FA、token、撤回、yank、deprecation、name transfer 等注册表制度。

Provenance 只能说明“这个 artifact 来自哪里、如何构建”，不能说明“这个 artifact 无恶意”。Scorecard 只能作为项目治理与开发流程证据，不能作为安全保证。没有公告不等于没有漏洞；有公告也需要区分版本范围、修复状态、撤回状态和可达性。

## 10. 跨生态发现数据结论

官方 registry 数据和第三方聚合数据必须分层。

官方事实源包括 npm registry、crates.io、PyPI、Go proxy/checksum/pkg.go.dev、Maven Central、NuGet.org、ConanCenter/vcpkg registries、SwiftPM registry protocol 等。它们保存或发布包的身份、版本、构件、元数据或解析所需索引。

第三方或增强发现面包括 deps.dev、Libraries.io、ecosyste.ms、lib.rs、Swift Package Index、MvnRepository、FuGet。它们能提供依赖图、反向依赖、安全公告、许可证、API 检视、文档、构建矩阵、维护信号和聚合搜索，但不能替代官方发布事实。

反向依赖尤其需要保留语义来源。GitHub dependents 是 GitHub 可见公共仓库引用，MvnRepository Used By 是 Maven artifact 层观察，pkg.go.dev imported-by 是 Go package import 层，Libraries.io dependent repositories 是聚合器观察结果，Swift Package Index 的依赖来自 Package.swift 和构建语境。这些不能合并成一个无来源的 dependency_count。

许可证也需要来源字段。registry_declared、repository_detected、aggregator_normalized、advisory_or_compliance_override 可能冲突。冲突时应标记冲突，不应自动选择一个第三方归一化值。

## 11. 可进入正式文档的变量组

身份变量：生态、registry、namespace/scope/group、package id、artifact id、module path、version、variant identity、source repository、owner/publisher。

声明变量：description、keywords、license、dependencies、optional dependencies、peer dependencies、dev dependencies、extras、features、targets、platforms、engines、toolchain version、repository URL、homepage、bugs/issues。

解析变量：version range、resolver family、conflict rule、dependency mediation、SemVer policy、MVS、backtracking、variant selection、target framework selection、feature resolution、lockfile presence、frozen/immutable mode。

构件变量：artifact type、tarball/wheel/crate/jar/nupkg/source zip/binary package、hash、signature、checksum database、provenance file、attestation predicate、binary cache key。

兼容变量：Python tag/ABI/platform、Go module major suffix、JVM bytecode/Gradle attributes/Kotlin target、NuGet TFM/RID、Conan settings/options/package_id、vcpkg triplet/features、Swift platforms/tools version/products/targets。

图谱变量：direct dependencies、transitive dependencies、dependency depth、duplicate versions、reverse dependencies、source importers、repository dependents、ecosystem centrality、important dependents。

维护变量：last release age、release cadence、issue response、PR merge, maintainers count、bus factor、deprecated/yanked/retracted/unpublished status、security policy。

安全变量：known vulnerabilities、malware advisory、fixed version、advisory source、review status、Scorecard checks、trusted publishing、2FA policy、token policy、SLSA level、Sigstore/Rekor material、provenance availability、artifact integrity.

发现变量：search fields、README availability、API docs, examples, type signatures, categories/classifiers/tags, downloads, stars, forks, SourceRank, score badges, docs build status, compatibility build matrix.

## 12. 当前剩余边界

评分函数可以设计，但不能先于变量字典完成。每个变量必须有数据来源、口径、缺失值处理、归一化方式和失真风险说明。

跨生态评分必须区分 package-level score 和 ecosystem-level score。包级分数回答候选包对具体查询和上下文的适配度；生态级分数回答某语言生态在生命周期某阶段的机制能力。

有些变量不适合跨生态直接比较。下载量、stars、Used By、imported-by、compatibility matrix、API surface size、SourceRank、Scorecard 总分都需要保留来源、覆盖范围和计算条件。

最终文档应先写术语表、生命周期模型、证据源分类、变量字典和生态事实表，再写评价函数。评价函数应给出可解释排序，不应输出没有来源说明的总分。

