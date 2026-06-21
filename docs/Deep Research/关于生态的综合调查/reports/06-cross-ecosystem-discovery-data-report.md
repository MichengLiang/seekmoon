# Deep Research Task 06 报告：跨生态发现数据、图谱与索引表面

## 1. 范围与对象边界

本文研究对象不是各语言包管理器本身，而是围绕包发现、元数据聚合、依赖图、反向依赖、安全公告、许可证、文档和 API 检视形成的“发现/数据表面”。这些表面可分为四类：一是官方注册表或官方生态入口，例如 NuGet.org API、Maven Central / Central Search、pkg.go.dev、GitHub Dependency Graph；二是跨生态第三方聚合器，例如 deps.dev、Libraries.io、ecosyste.ms；三是单生态增强索引，例如 lib.rs、Swift Package Index、MvnRepository、FuGet；四是安全公告和依赖风险表面，例如 GitHub Advisory Database、OSV 及 Go Vulnerability Database 的下游展示。

本报告只评价数据来源、字段能力和可靠性边界，不评价网站审美，不声称用户流行度，不给出实现计划。报告中的“官方”指由生态、注册表、包托管服务或语言项目维护者提供的事实源或正式接口；“第三方”指不拥有注册表写入权、通过抓取、索引、解析、克隆源码或调用官方 API 聚合数据的表面；“图谱支持”指能表达依赖边或反向依赖边，而不是仅显示包名称；“观察性声明”指公开页面可见但缺少正式接口文档支持的能力。

## 2. 发现/数据表面的分类

| 表面 | 类型 | 数据基础 | 事实源地位 | 主要用途 |
|---|---|---|---|---|
| deps.dev / Open Source Insights | 跨生态、第三方聚合、图谱支持 | npm、Maven、PyPI、Cargo、NuGet、Go、RubyGems、GitHub/GitLab/Bitbucket、OSV 等 | Google 托管的聚合与派生数据，不是各注册表的写入事实源 | 包检索、版本比较、依赖图、安全公告、许可证、项目到包映射 |
| Libraries.io | 跨生态、第三方聚合 | 包管理器、源码仓库、公开网页抓取 | 明确声明数据来自互联网抓取且不保证准确 | 搜索、API、依赖、反向依赖、SourceRank、仓库信号 |
| ecosyste.ms | 跨生态开放 API 套件 | 包注册表、仓库、归档、安全与许可证服务 | 第三方开放数据基础设施 | 包/版本/依赖元数据、依赖解析、仓库、许可证、SBOM、归档内容分析 |
| lib.rs | Rust 单生态第三方发现增强 | crates.io、源码仓库、Rust 生态页面的二次索引 | 非 crates.io 官方事实源 | Rust 包发现、分类、依赖和维护信号补充 |
| Swift Package Index | Swift 单生态第三方/社区索引 | PackageList、源码克隆、Package.swift、Git 历史、GitHub 元数据、构建结果 | 非 SwiftPM 注册表；是 Swift 包搜索和兼容性索引 | Swift 包搜索、平台/Swift 版本兼容性、构建状态、DocC 文档 |
| MvnRepository | JVM/Maven 单生态第三方发现页 | Maven 坐标、POM、仓库元数据、外部安全数据 | 非 Maven Central 官方入口 | 坐标检索、依赖片段、版本、反向使用、漏洞提示 |
| FuGet | NuGet 单生态第三方包/API 浏览器 | NuGet.org API、包内 XML 文档和程序集元数据 | 非 NuGet.org 官方入口 | NuGet 包浏览、目标框架、依赖、API/类型/成员检视、版本 diff |
| pkg.go.dev | Go 官方发现/文档入口 | Go 模块代理、源码、Go 文档生成、Go 漏洞数据库 | Go 官方包发现与文档表面 | 模块/包文档、搜索、符号、版本、导入者、漏洞 |
| GitHub Dependency Graph / Advisory Database | GitHub 官方仓库图谱与安全公告 | 仓库 manifest/lockfile、Dependency Submission API、GitHub Advisory Database、NVD、OSV 格式公告 | GitHub 仓库层事实源；不是各注册表的唯一事实源 | 仓库依赖图、反向依赖、SBOM、Dependabot alerts、GHSA 安全公告 |

deps.dev 的官方文档说明 Open Source Insights 是 Google 开发和托管的服务，会检查包、构造完整依赖图并定期更新数据，覆盖 npm、Maven、PyPI、Cargo、NuGet、Go、RubyGems、项目托管源和 OSV 安全公告；API 文档列出 GetPackage、GetVersion、GetRequirements、GetDependencies、GetProject、GetAdvisory 等接口。来源：`https://docs.deps.dev/`，`https://docs.deps.dev/api/v3/`。

Libraries.io 在数据页中明确说自己收集公开开源包信息，并说明免费数据“scraped from the internet and not validated, corrected, or curated for accuracy”。其 API 页要求 `api_key`，并列出 platforms、project、project dependencies、project dependents、dependent repositories、contributors、SourceRank、search、repository dependencies 等端点。来源：`https://libraries.io/data`，`https://libraries.io/api`。

ecosyste.ms 的 packages 仓库 README 将其定义为“open API service providing package, version and dependency metadata of many open source software ecosystems and registries”，并指向 `https://packages.ecosyste.ms/docs`；resolve 仓库定义为跨生态依赖树解析服务，CLI 仓库还列出 advisories、archives、licenses、repos、resolve、sbom 等 API 命令组。来源：`https://github.com/ecosyste-ms/packages`，`https://github.com/ecosyste-ms/resolve`，`https://github.com/ecosyste-ms/ecosyste_ms_cli`。

Swift Package Index FAQ 说明索引数据来自包仓库列表；加入列表后会完整克隆包源码，从源码、Git 历史和 GitHub 托管仓库抽取元数据，并每隔数小时轮询变化。PackageList 仓库则给出入选约束：公开可访问、根目录有有效 `Package.swift`、Swift 5.0 或更高、语义化版本 tag、`swift package dump-package` 能输出有效 JSON、可编译等。来源：`https://swiftpackageindex.com/faq`，`https://github.com/SwiftPackageIndex/PackageList`。

FuGetGallery 源码仓库说明 FuGet 是 NuGet package browser 与 API browser 的组合，包浏览器使用 nuget.org API 展示索引包，API 浏览器结合 XML 文档和程序集元数据帮助浏览 API，并支持通过 code tab 和版本 diff 查看变化。来源：`https://github.com/praeclarum/FuGetGallery`。

pkg.go.dev 是 Go 官方包发现与文档表面。Go 博客在 2026 年 5 月 21 日介绍 pkg.go.dev API，说明 pkg.go.dev 是 Go 社区包文档和发现的主要资源，并提供 `v1beta` GET-only API，包含 package、module、versions、packages、search、symbols 等端点；Google Open Source Blog 后续说明还列出 imported-by 与 vulns 端点。来源：`https://go.dev/blog/pkgsite-api`，`https://opensource.googleblog.com/2026/06/a-new-pkggodev-api-for-go.html`，`https://pkg.go.dev/api`。

GitHub Dependency Graph 官方文档定义它是仓库 manifest 与 lock files，以及 Dependency Submission API 提交依赖的摘要；它展示 dependencies 与 dependents，并对每个依赖显示版本、许可证、manifest 文件和已知漏洞。来源：`https://docs.github.com/en/code-security/concepts/supply-chain-security/dependency-graph`。

## 3. 各表面的字段与证据类型

deps.dev 的核心字段围绕“包、版本、项目、公告、解析图”组织。GetPackage 返回包及版本列表；GetVersion 返回指定版本信息、许可证和直接影响该版本的安全公告；GetRequirements 返回版本声明的依赖要求，并按生态保留 npm devDependencies / peerDependencies、RubyGems runtime/dev、Go direct/indirect、PyPI environment marker 等生态差异；GetDependencies 返回已解析依赖图；GetProject 返回 GitHub、GitLab 或 Bitbucket 项目信息；GetProjectPackageVersions 建立项目与包版本映射；GetAdvisory 返回 OSV 托管公告。它适合作为跨生态评估模型中的“规范化包版本节点 + 依赖边 + 许可证 + 安全公告”数据源，但它仍然是聚合与派生层，不能替代注册表原始元数据。来源：`https://docs.deps.dev/api/v3/`。

Libraries.io 暴露字段更偏发现和维护信号。API 示例中可见 package name、platform、description、homepage、repository_url、normalized_licenses、latest_release_number、latest_release_published_at、stars、forks、rank、dependent repositories、contributors、SourceRank 等字段。其优势是覆盖包管理器多、对仓库和包之间的关系做了归并；其弱点是官方明示数据未验证，许可证和维护字段必须作为提示而非事实裁决。来源：`https://libraries.io/api`，`https://libraries.io/data`。

ecosyste.ms 的 packages 服务适合读取 package、version、dependency metadata；resolve 服务适合依赖树解析；archives 服务适合查看包归档内容、README/changelog 抽取和文件内容；CLI README 显示其 API 族还覆盖 advisories、licenses、repos、sbom、sponsors、timeline 等。因为本次研究中 `packages.ecosyste.ms` 直接 HTTPS 请求失败，本文对 ecosyste.ms 的具体端点能力主要依据其 GitHub README 和 OpenAPI 文件说明，不把单次网络失败解释为服务不可用。来源：`https://github.com/ecosyste-ms/packages`，`https://github.com/ecosyste-ms/packages/blob/main/openapi/api/v1/openapi.yaml`，`https://github.com/ecosyste-ms/archives`，`https://github.com/ecosyste-ms/ecosyste_ms_cli`。

lib.rs 是 Rust 发现增强面，不是 crates.io 官方事实源。它适合补充 crates.io 官方页面之外的分类、可读性排序、依赖展示、替代包发现和维护线索。由于 lib.rs 页面在本次命令行访问中触发 Cloudflare 挑战，报告只将 `https://lib.rs/crates/serde` 这类页面作为可观察表面，不把它作为可自动化依赖图事实源。Rust 的权威包发布事实仍应回到 crates.io；crates.io 的数据访问政策也说明其公共 API 使用受策略约束，不能无约束抓取。来源：`https://crates.io/data-access`，观察页面：`https://lib.rs/crates/serde`。

Swift Package Index 的字段来自源码与构建结果，而不是中心注册表。它可显示包、作者/仓库信息、平台兼容性、Swift 版本兼容性、构建状态、文档生成状态、release/tag、README、许可证和依赖。该表面的特别价值是“能否在特定 Swift / 平台矩阵下构建”和“DocC 文档是否可生成”，这些变量在 Swift 生态中比单纯下载量更接近可用性。来源：`https://swiftpackageindex.com/faq`，`https://github.com/SwiftPackageIndex/SwiftPackageIndex-Server`，`https://github.com/SwiftPackageIndex/PackageList`。

MvnRepository 是 JVM/Maven 第三方发现页。以 `https://mvnrepository.com/artifact/org.springframework/spring-core` 这类页面为例，可观察字段包括版本列表、Maven/Gradle/SBT/Ivy/Leiningen 依赖片段、licenses、categories、tags、compile/runtime/test dependencies、used by 和 vulnerabilities。它的价值是快速发现和复制坐标，但 Maven Central 的官方事实源是 Central Repository / Central Search；Sonatype 文档说明 Central Search 官方入口是 `https://central.sonatype.com`，并提供 REST API 搜索能力。来源：`https://mvnrepository.com/artifact/org.springframework/spring-core`，`https://central.sonatype.org/`，`https://central.sonatype.org/search/`，`https://central.sonatype.org/search/rest-api-guide/`。

FuGet 的字段更偏 .NET API 检视。它基于 NuGet.org API 显示包索引和依赖，同时解析包内 XML 文档与程序集元数据，让使用者浏览命名空间、类型、成员、目标框架和版本差异。它不是 NuGet.org 权威源，包发布、版本、元数据、弃用、漏洞等仍应回到 NuGet.org API。NuGet 官方 API 文档说明 service index、catalog、registration、package metadata 等资源；catalog 是 append-only API，可查看发布、修改、删除历史。来源：`https://github.com/praeclarum/FuGetGallery`，`https://learn.microsoft.com/en-us/nuget/api/overview`，`https://github.com/NuGet/docs.microsoft.com-nuget/blob/main/docs/guides/api/query-for-all-published-packages.md`。

pkg.go.dev 暴露的是 Go 模块/包文档和结构化元数据。官方 v1beta API 覆盖 package、module、versions、module packages、search、symbols、imported-by、vulns 等。Go Vulnerability Management 文档说明 Go 团队维护漏洞报告流水线，报告存储在 Go vulnerability database 中，并集成到 Go 包发现站点和 govulncheck。来源：`https://go.dev/blog/pkgsite-api`，`https://go.dev/doc/security/vuln/`。

GitHub Dependency Graph 暴露仓库依赖视角：依赖所属生态、包名、版本、许可证、manifest 文件、已知漏洞、传递路径、dependents、“Used by”、SBOM。GitHub Advisory Database 暴露 GHSA ID、CVE、生态、包名、受影响版本、修复版本、严重性、CVSS、EPSS、参考链接、workaround、credit 等。来源：`https://docs.github.com/en/code-security/concepts/supply-chain-security/dependency-graph`，`https://docs.github.com/en/code-security/concepts/vulnerability-reporting-and-management/github-advisory-database`，`https://docs.github.com/en/rest/dependency-graph/sboms`。

## 4. 依赖图与反向依赖支持

跨生态依赖图的最大问题不是“是否能显示依赖”，而是“依赖边语义是否可比”。npm 的 peerDependencies、optionalDependencies、bundledDependencies，Go 的 direct/indirect，Maven 的 scope 与 optional，NuGet 的 target framework 分组，PyPI 的 environment marker，Cargo 的 feature 与 target 条件，并不是同一种关系。deps.dev 的 GetRequirements 保留生态特定结构，GetDependencies 提供解析后的图，因此更适合作为跨生态“同名字段规范化”与“保留原始语义”之间的中间层。来源：`https://docs.deps.dev/api/v3/`。

Libraries.io 支持 project dependencies、project dependents 和 dependent repositories，适合做“哪些包/仓库依赖此包”的发现线索。但其反向依赖受抓取覆盖、仓库解析、生态支持和数据新鲜度影响，不能视为完整图谱。来源：`https://libraries.io/api`，`https://libraries.io/data`。

ecosyste.ms packages 暴露 package/version/dependency metadata，resolve 服务专门用于解析跨生态 dependency trees；其 issue 讨论中也出现 dependent_packages 等端点，说明它以跨生态依赖关系为对象。来源：`https://github.com/ecosyste-ms/packages`，`https://github.com/ecosyste-ms/resolve`，`https://github.com/ecosyste-ms/packages/issues/407`。

GitHub Dependency Graph 的依赖图是仓库视角，而不是注册表全量视角。它从仓库中的 manifest/lockfile 和 Dependency Submission API 构造图；public repository 可以显示 dependents，但 private repository 不报告 dependents；“Used by”还要求包发布到支持生态、源码指向公开仓库、超过 100 个仓库依赖等条件。这个边界意味着 GitHub reverse dependency 数值是“GitHub 可见公共引用”，不是全互联网下载或全注册表依赖。来源：`https://docs.github.com/en/code-security/concepts/supply-chain-security/dependency-graph`。

MvnRepository、Swift Package Index、pkg.go.dev、FuGet 和 lib.rs 都支持不同程度的反向或邻接关系观察，但不能直接等价。MvnRepository 的 Used By 是 Maven artifact 层观察；pkg.go.dev 的 imported-by 是 Go package import 层；Swift Package Index 的依赖与兼容性来自 SwiftPM manifest 和构建；FuGet 主要服务 NuGet 包/程序集 API 检视；lib.rs 是 Rust 发现增强。跨生态模型应将这些字段拆成“声明依赖”“解析依赖”“源码导入”“仓库使用”“公开引用”五类，不能合并成一个无差别 dependency_count。

## 5. 安全、许可证与公告支持

安全公告的权威层次需要分开。GitHub Advisory Database 包含 GitHub-reviewed advisories、unreviewed advisories 和 malware advisories；reviewed advisories 会映射到 GitHub 支持生态中的包，unreviewed 来自 NVD feed 且不触发 Dependabot alerts，malware advisories 目前与 npm 生态相关。每条公告可含 GHSA ID、CVE、生态、包、受影响版本、修复版本、严重性、CVSS、EPSS 等。来源：`https://docs.github.com/en/code-security/concepts/vulnerability-reporting-and-management/github-advisory-database`。

deps.dev 的 GetVersion 返回直接影响该版本的安全公告，但文档明确该字段不包括影响依赖项的公告；GetAdvisory 返回 OSV 托管公告。这适合做“包版本本身是否被公告影响”的判断，不适合直接替代项目级传递风险扫描。来源：`https://docs.deps.dev/api/v3/`。

Go 生态中，Go Vulnerability Database 是官方漏洞数据源之一，govulncheck 通过静态分析源代码或二进制符号表，尽量缩小到实际可达漏洞；pkg.go.dev 与 Go vulnerability management 集成。这个模型比“包版本有公告”更精细，因为它关心 vulnerable symbol 是否被调用。来源：`https://go.dev/doc/security/vuln/`，`https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck`。

许可证字段同样有事实源差异。包元数据中的 license、源码仓库中的 LICENSE 文件、GitHub 许可证识别、deps.dev SPDX-like 归一化、Libraries.io normalized_licenses、Maven POM 中 licenses、NuGet package license expression 并不总是一致。Libraries.io 已明确说明免费数据不被验证；GitHub Dependency Graph 只显示它能识别的许可证；pkg.go.dev 甚至会受 redistributable license 判断影响展示。跨生态模型中，许可证字段应至少保留 source 字段：registry_declared、repository_detected、aggregator_normalized、advisory_or_compliance_override。

## 6. 文档与 API 检视支持

文档和 API 检视是单生态表面最有价值的差异化能力。pkg.go.dev 从 Go 源码生成包文档、符号和模块信息，并提供官方 v1beta API，适合自动读取 Go 包 synopsis、symbols、versions 和漏洞。来源：`https://go.dev/blog/pkgsite-api`，`https://pkg.go.dev/api`。

FuGet 的核心价值是 .NET API browser：它结合 XML 文档和程序集元数据来浏览命名空间、类型、成员，并能比较版本差异。这是 NuGet.org 包元数据之外的程序集级检视能力。来源：`https://github.com/praeclarum/FuGetGallery`。

Swift Package Index 可生成并托管 DocC 文档，并用构建矩阵表达平台和 Swift 版本兼容性；这些字段来自源码克隆和构建流程，不是 SwiftPM 注册表字段。来源：`https://swiftpackageindex.com/faq`，`https://forums.swift.org/t/swift-package-index-auto-generating-auto-hosting-and-auto-updating-docc-documentation/57806`。

Maven Central / Central Search 与 MvnRepository 更偏坐标、POM、版本和依赖片段；Javadoc 是否可用通常依赖 artifact 是否发布了 javadoc/sources。Central Repository 发布指南说明发布到 Central 需要 artifact、POM、PGP 签名、sources/javadoc 等质量元数据。来源：`https://maven.apache.org/repository/guide-central-repository-upload.html`。

deps.dev 和 Libraries.io 不主要承担语言 API 文档角色。deps.dev 更适合依赖和安全图谱；Libraries.io 更适合包级搜索和维护信号。跨生态模型中，文档质量或 API 可检视性只能在有语言专门索引时作为生态内指标，不能简单跨语言比较。

## 7. 可用于统一模型的变量与指标

可跨生态比较的低风险变量包括：包标识符、生态/注册表、版本号、发布时间、最新版本、源码仓库 URL、声明许可证、归一化许可证、是否存在安全公告、公告严重性、直接依赖数量、解析依赖图是否可得、反向依赖或导入者是否可得、维护时间戳、仓库活动时间、是否有文档入口、是否有 README、是否有源码链接、是否有 SBOM 或可导出依赖清单。

需要谨慎比较的变量包括：下载量、stars、forks、SourceRank、dependent repositories、Used By、imported-by、MvnRepository Used By、Swift Package Index compatibility、pkg.go.dev redistributable、FuGet API surface。它们的采集范围和语义不同：下载量受代理缓存、CI、镜像和生态安装习惯影响；stars/forks 是仓库信号而非包质量；反向依赖可能只覆盖公开 GitHub 仓库或特定注册表；兼容性矩阵只对 Swift 生态自然成立；API surface 的大小不能直接代表质量。

统一模型应保留三层证据，而不是把所有字段合成为一个分数。第一层是注册表/官方事实：包名、版本、发布时间、发布者、artifact、POM/package manifest、官方 license expression、官方 deprecation/vulnerability 字段。第二层是派生图谱：解析依赖、传递依赖、反向依赖、项目到包映射、安全公告映射。第三层是发现和维护信号：README、文档、仓库活跃度、issue/PR、stars/forks、SourceRank、构建兼容性、API 可检视性。只有第一层适合高权重事实判断；第二层适合风险和影响面分析；第三层只适合排序提示或人工复核入口。

## 8. 事实源与可靠性问题

官方注册表数据也不等于绝对真相。包作者可填写不完整或错误的 license、repository、homepage、description；某些生态允许 yanking、unlisting、deprecation 或删除；同一个源码仓库可能发布多个包；同一个包名在不同注册表中可能无关；源码仓库链接可能迁移或被接管；manifest 中的依赖约束不等于一次具体安装的解析结果。

第三方聚合器的问题更明显。Libraries.io 明确说明免费数据未验证、未校正、未策展；因此它的许可证、维护和仓库关系字段只能作为观察线索。deps.dev 虽由 Google 托管并有清晰 API，但仍聚合多来源并生成派生图；其公告字段对直接版本有效，不覆盖依赖项公告。ecosyste.ms 以开放 API 和开放数据为优势，但其具体服务多、端点和字段需要按 API 分别验证。MvnRepository、FuGet、lib.rs 和 Swift Package Index 都是单生态增强面，不能替代对应官方发布事实。

反向依赖尤其容易被误读。GitHub Dependency Graph 的 dependents 仅对公共仓库报告，Used By 还有生态、公开仓库链接和数量门槛；Libraries.io dependent repositories 依赖其抓取覆盖；Maven used-by 可能只覆盖 Maven 索引关系；Go imported-by 是源码 import 关系；Swift dependency 是 Package.swift 声明和构建上下文。统一模型中应把这些字段命名为 `reverse_signal_type`，并保留 `coverage_scope`，例如 `github_public_repositories`、`registry_declared_dependents`、`source_importers`、`aggregator_observed_dependents`。

安全公告也存在来源层级。GitHub-reviewed advisory、NVD-imported unreviewed advisory、OSV advisory、Go vulnerability database、NuGet.org vulnerable 标记、MvnRepository 漏洞提示不是同一个审查级别。跨生态风险判断应记录 advisory source、review status、affected range、patched range、database update time 和 whether transitively applicable。不能把“某页面显示漏洞”直接等同于“项目一定可利用”。

许可证判断应优先使用官方声明与源码许可证检测的组合。若 registry_declared 与 repository_detected 冲突，应标记为冲突，而不是自动选择第三方归一化值。对合规用途，第三方字段只能作为待复核提示；对发现用途，第三方归一化字段可以提高召回。

## 9. 未解决不确定性与参考资料

未解决不确定性包括：第一，ecosyste.ms 在线 API 在本次环境中出现 HTTPS 连接失败，报告对其字段能力主要依据 GitHub README 和 OpenAPI 文件，不据此判断服务稳定性。第二，lib.rs 页面在命令行访问时触发 Cloudflare 挑战，因此本文只将其列为 Rust 发现增强面，不使用其页面字段作为自动化事实源。第三，MvnRepository 页面字段来自公开观察，缺少正式机器接口说明；严肃数据抽取应优先使用 Central Search、Maven Central index 或 Sonatype Central。第四，pkg.go.dev API 在 2026 年仍标记 `v1beta`，正式 v1 前字段可能变化。第五，跨生态“维护度”“健康度”“受欢迎度”没有统一物理量；任何统一评分都必须暴露来源权重和不可比字段。

参考资料：

- deps.dev / Open Source Insights 文档：`https://docs.deps.dev/`
- deps.dev API v3：`https://docs.deps.dev/api/v3/`
- Libraries.io API：`https://libraries.io/api`
- Libraries.io data and accuracy statement：`https://libraries.io/data`
- ecosyste.ms packages：`https://github.com/ecosyste-ms/packages`
- ecosyste.ms packages OpenAPI：`https://github.com/ecosyste-ms/packages/blob/main/openapi/api/v1/openapi.yaml`
- ecosyste.ms resolve：`https://github.com/ecosyste-ms/resolve`
- ecosyste.ms CLI API groups：`https://github.com/ecosyste-ms/ecosyste_ms_cli`
- crates.io data access policy：`https://crates.io/data-access`
- lib.rs observable crate page：`https://lib.rs/crates/serde`
- Swift Package Index FAQ：`https://swiftpackageindex.com/faq`
- Swift Package Index PackageList：`https://github.com/SwiftPackageIndex/PackageList`
- Swift Package Index server：`https://github.com/SwiftPackageIndex/SwiftPackageIndex-Server`
- Swift Package Index DocC announcement：`https://forums.swift.org/t/swift-package-index-auto-generating-auto-hosting-and-auto-updating-docc-documentation/57806`
- Maven Central documentation：`https://central.sonatype.org/`
- Central Search：`https://central.sonatype.org/search/`
- Central Search REST API：`https://central.sonatype.org/search/rest-api-guide/`
- Apache Maven Central index：`https://maven.apache.org/repository/central-index.html`
- Maven Central upload guide：`https://maven.apache.org/repository/guide-central-repository-upload.html`
- MvnRepository observation page：`https://mvnrepository.com/artifact/org.springframework/spring-core`
- FuGetGallery source and description：`https://github.com/praeclarum/FuGetGallery`
- NuGet API overview：`https://learn.microsoft.com/en-us/nuget/api/overview`
- NuGet catalog guide：`https://github.com/NuGet/docs.microsoft.com-nuget/blob/main/docs/guides/api/query-for-all-published-packages.md`
- pkg.go.dev API announcement：`https://go.dev/blog/pkgsite-api`
- pkg.go.dev API reference：`https://pkg.go.dev/api`
- Google Open Source Blog on pkg.go.dev API：`https://opensource.googleblog.com/2026/06/a-new-pkggodev-api-for-go.html`
- Go vulnerability management：`https://go.dev/doc/security/vuln/`
- govulncheck documentation：`https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck`
- GitHub Dependency Graph：`https://docs.github.com/en/code-security/concepts/supply-chain-security/dependency-graph`
- GitHub supported dependency graph ecosystems：`https://docs.github.com/en/code-security/reference/supply-chain-security/dependency-graph-supported-package-ecosystems`
- GitHub Advisory Database：`https://docs.github.com/en/code-security/concepts/vulnerability-reporting-and-management/github-advisory-database`
- GitHub SBOM REST API：`https://docs.github.com/en/rest/dependency-graph/sboms`
