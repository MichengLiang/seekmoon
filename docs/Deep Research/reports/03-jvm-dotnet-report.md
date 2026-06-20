# JVM/Kotlin 与 .NET 包生态研究报告

## 1. Scope and object boundary

本文研究对象是“软件包复用”在 JVM/Kotlin 与 .NET/C# 两个生态中的身份、元数据、发现、解析、兼容性、可复现性与安全暴露面。这里的“包”不是单一安装器命令，也不是源码仓库本身，而是被发布到包仓库、由构建工具或 IDE 消费、并通过元数据参与依赖图求解的可复用软件单元。

JVM/Kotlin 侧的核心对象包括 Maven 坐标、POM、Maven Central、Gradle 依赖管理、Gradle Module Metadata、Gradle variant model 与 Kotlin Multiplatform 发布格式。Maven 官方命名文档把 Maven coordinates 定义为 `groupId`、`artifactId`、`version` 三个标识符的组合，并说明 artifact 需要由这组三元组唯一识别；Sonatype Central 文档也把 GAV 坐标视为 Maven Central 发布位置的决定因素（https://maven.apache.org/guides/mini/guide-naming-conventions.html，https://central.sonatype.org/publish/requirements/）。.NET 侧的核心对象包括 NuGet 包、`.nuspec`、`PackageReference`、`packages.lock.json`、target frameworks 与 NuGet audit。NuGet 官方版本文档明确说，具体包由 package identifier 与精确版本号共同指称（https://learn.microsoft.com/en-us/nuget/concepts/package-versioning）。

本文不把 Maven、Gradle、NuGet 视为可互换安装器。Maven 主要以 POM 与 Maven repository layout 组织 artifact；Gradle 在可消费 Maven/Ivy 元数据之外，引入 variant-aware resolution 与 Gradle Module Metadata；NuGet 则围绕 `.nupkg`、`.nuspec`、TFM 与 MSBuild 项目文件中的 `PackageReference` 工作。这些对象都支持包复用，但它们暴露的身份粒度、兼容性判断变量和冲突调解规则不同。

## 2. JVM/Kotlin lifecycle mapping

JVM 包生命周期可以分为声明、发布、发现、解析、选择 artifact、构建使用与安全审查几个阶段。Maven 的声明面是 `pom.xml`。Apache Maven 的 POM 介绍把 Project Object Model 称为 Maven 的基本工作单元，POM 是包含项目信息与构建配置的 XML 文件（https://maven.apache.org/guides/introduction/introduction-to-the-pom.html）。POM 参考还说明 Maven repository artifact 在仓库中按 `$BASE_REPO/groupId/artifactId/version/artifactId-version.pom` 布局，因此坐标不仅是逻辑身份，也是仓库路径结构的一部分（https://maven.apache.org/pom.html）。

发布阶段由仓库规则约束。Maven Central 要求发布组件包含 POM 元数据；Central 文档列出正确坐标、项目名称、描述、URL、license、developer、SCM 等最低元数据要求，并强调依赖声明用于消费者正确解析传递依赖（https://central.sonatype.org/publish/requirements/）。Apache Maven 的 Central 上传指南还要求部署 artifact、POM 及其 PGP signatures，并说明 Central release 不会在发布后被修改或替换，且通常需要 javadoc 与 sources 支持 IDE 查阅（https://maven.apache.org/repository/guide-central-repository-upload.html）。

Gradle 的生命周期不是简单复制 Maven。Gradle 可以消费 Maven POM，但其依赖模型以 configuration、attribute、capability、variant 与 resolution result 为核心。Gradle 依赖管理文档把外部模块身份表述为 group、name、version，并把 `implementation`、`api` 等 configuration 作为依赖应用范围（https://docs.gradle.org/current/userguide/getting_started_dep_man.html）。Gradle Module Metadata 则用于序列化 Gradle component model；官方文档说明它类似 POM 或 Ivy descriptor，但目标是提供 multi-platform 与 variant-aware 的发布模型（https://docs.gradle.org/current/userguide/publishing_gradle_module_metadata.html）。

Kotlin Multiplatform 把这个生命周期进一步扩展。Kotlin 官方发布文档说明，当使用 `maven-publish` 时，Kotlin 插件会为当前 host 可构建的每个 target 自动创建 publication，Android target 需要额外配置（https://kotlinlang.org/docs/multiplatform/multiplatform-publish-lib-setup.html）。这意味着 KMP 包不是“一个 jar 对所有平台”的普通 Maven artifact。它通常有根 publication、metadata publication 与各平台 target publication；消费者需要根据 JVM、JS、Native、Android 等目标及源集关系选择正确变体。因而 KMP 的可复用对象边界必须包括目标平台、源集、编译产物类型与 Gradle metadata。

## 3. .NET lifecycle mapping

.NET 包生命周期以 NuGet 包创建、发布、发现、项目引用、restore、兼容性检查、锁定和 audit 为主。NuGet 官方介绍将 `.nupkg` 定义为带 `.nupkg` 扩展名的 ZIP 文件，包含编译代码、相关文件与描述性 manifest；NuGet Gallery 是包作者和消费者使用的中心包仓库（https://learn.microsoft.com/en-us/nuget/what-is-nuget，https://www.nuget.org/）。

NuGet 的包元数据面是 `.nuspec`。Microsoft 的 `.nuspec` 参考明确说明 `.nuspec` 是包含包元数据的 XML manifest，用于构建包并向消费者提供信息，而且 manifest 总是包含在包中（https://learn.microsoft.com/en-us/nuget/reference/nuspec）。现代 SDK-style 项目可以把通常出现在 `.nuspec` 的属性放入项目文件，由 `dotnet pack` 或 MSBuild pack target 生成包；非 SDK-style 或特殊打包场景仍可直接使用 `.nuspec`。

消费阶段以 `PackageReference` 为主。Microsoft 文档定义 `PackageReference` 为 MSBuild items，直接在项目文件中指定 NuGet 包依赖，而不是放在独立 `packages.config` 文件中；NuGet.Config 中的 package sources 等设置仍然参与 restore（https://learn.microsoft.com/en-us/nuget/consume-packages/package-references-in-project-files）。Restore 阶段不只是下载包，还会依据包依赖、项目目标框架、源配置、锁文件与审计配置求解 dependency graph。

## 4. Metadata and identity surfaces

Maven identity 是分层坐标。`groupId` 通常承担组织或命名空间职责，Maven getting started 文档说明它通常基于组织的 fully qualified domain name；`artifactId` 是同一 group 下的组件名；`version` 标识发布版本（https://maven.apache.org/guides/getting-started/，https://maven.apache.org/guides/mini/guide-naming-conventions.html）。这与扁平包名不同：`jackson-core` 这种 artifact 名不能单独稳定标识一个包，必须与 `com.fasterxml.jackson.core` 和具体版本组合。Maven dependency mechanism 文档还指出在 `dependencyManagement` 匹配中，完整匹配键实际包括 `groupId`、`artifactId`、`type`、`classifier`，只是在普通 jar 且无 classifier 时可简写为 `groupId`、`artifactId`（https://maven.apache.org/guides/introduction/introduction-to-dependency-mechanism.html）。

POM 元数据覆盖项目身份、父 POM、packaging、name、description、url、licenses、developers、SCM、properties、dependencies、dependencyManagement、repositories、build plugins 等对象。POM 既是构建描述，也是消费者解析传递依赖的元数据源。Maven Central 的发布要求把正确 dependencies 视为推荐元数据，因为构建工具依赖这些信息解析传递依赖（https://central.sonatype.org/publish/requirements/）。

Gradle Module Metadata 的身份仍可落在 Maven 仓库的 group/name/version 上，但它增加了 POM 难以表达的变体结构。Gradle 官方文档说明 GMM 包含 variants 细节，并用于改进 multi-platform、variant-aware dependency resolution（https://docs.gradle.org/current/userguide/variant_attributes.html，https://docs.gradle.org/current/userguide/publishing_gradle_module_metadata.html）。一个 component 可以有 `apiElements`、`runtimeElements`、平台 target、usage、category、library elements、JVM version、Kotlin platform type 等属性组合；这些不是普通 POM dependency scope 的同义词，而是消费者选择 artifact 的输入变量。

NuGet identity 更接近“扁平 ID + 版本”。官方版本文档说明具体包由 package identifier 与 exact version number 指称（https://learn.microsoft.com/en-us/nuget/concepts/package-versioning）。`.nuspec` 元数据包括 `id`、`version`、`authors`、`description`、license、projectUrl、icon、readme、repository、tags、dependencies、framework assemblies、content files、package types 等字段；其中 dependencies 可按 target framework 分组，以表达不同 TFM 下的依赖闭包（https://learn.microsoft.com/en-us/nuget/reference/nuspec）。NuGet 包名没有 Maven 那样内置的 `groupId` 层级；组织归属、可信度和命名空间通常通过 ID 前缀、作者、verified owner、repository metadata 与 nuget.org 页面来辅助判断。

## 5. Discovery surfaces

Maven Central 是 JVM 生态的一手发现与分发面。Central Portal 首页将自己标识为 official search，并提供 Browse、Search、API Doc、Namespace、Publishing 等入口（https://central.sonatype.com/）。`search.maven.org` 也显示为 Maven Central Repository 维护者的 official search（https://search.maven.org/）。这些界面以坐标、版本、artifact 文件、依赖、发布者命名空间、校验与安全信息为发现变量。

MvnRepository 是第三方发现面，不应等同于 Maven Central 本身。它提供分类浏览、artifact 页面、版本列表、依赖片段和安全提示等聚合信息，但其数据解释和页面排序不是 Maven 官方契约。本文仅将它作为第三方发现界面示例（https://mvnrepository.com/）。

Gradle Plugin Portal 是 Gradle 插件发现面。Portal 首页提供插件搜索与发布文档入口（https://plugins.gradle.org/），Gradle 官方发布文档说明，把插件发布到 Gradle Plugin Portal 是使插件在全球范围内可发现和可使用的主要方式（https://docs.gradle.org/current/userguide/publishing_gradle_plugins.html）。它发现的是 Gradle plugin id 与插件版本，不等同于普通 Maven library dependency 的发现。

NuGet.org 是 .NET 生态中心发现面。NuGet.org 首页说明 NuGet Gallery 是包作者和消费者使用的 central package repository（https://www.nuget.org/）。包页面通常展示 ID、版本、owners、下载量、README、frameworks、dependencies、vulnerabilities、deprecation、license、repository、install 命令等信息。Visual Studio 的 NuGet Package Manager 也是 IDE 发现面；Microsoft 文档说明 Browse tab 会按当前 package source 显示包，Updates tab 显示可更新包，Solution 管理界面还包含 Consolidate tab 用于发现同一 solution 中不同项目使用的不同包版本（https://learn.microsoft.com/en-us/nuget/consume-packages/install-use-packages-visual-studio）。

JVM 侧 IDE 发现也不是独立包仓库。JetBrains 文档说明 IntelliJ IDEA 支持在 POM 中添加 Maven dependency、导入依赖、查看依赖图并分析 unresolved、conflict 与 transitive dependencies（https://www.jetbrains.com/help/idea/work-with-maven-dependencies.html）；对 Gradle 项目，IDEA 通过 Gradle 模型显示依赖图、版本、group 和 artifact ID，并把依赖管理交回 build.gradle（https://www.jetbrains.com/help/idea/work-with-gradle-dependency-diagram.html）。因此 IDE 包管理器更准确地说是发现、编辑、分析和可视化界面，不是依赖解析规则的来源。

## 6. Dependency resolution, variants, and compatibility

Maven 的依赖解析以传递依赖、scope、dependencyManagement、exclusion、optional 与 mediation 组成。官方依赖机制文档说明 Maven 会读取依赖项目的 project files，并自动包含传递依赖；当多个版本出现时，Maven 使用 “nearest definition”，即选择依赖树中离当前项目最近的版本，同深度时 first declaration wins（https://maven.apache.org/guides/introduction/introduction-to-dependency-mechanism.html）。`dependencyManagement` 可以在不直接添加依赖的情况下管理版本，并且对传递依赖的版本选择具有优先作用。这个模型的兼容性判断主要围绕同一 GAV artifact、scope 与类路径形成，并不内建“一个组件多个消费变体”的完整语义。

Gradle 的解析分为 graph resolution 与 artifact resolution。官方文档说明 graph resolution 根据声明依赖和 metadata 构造依赖图，并使用被解析 configuration 的 request attributes；artifact resolution 再为 resolved components 选择具体 variants 与 artifacts（https://docs.gradle.org/current/userguide/dependency_resolution.html）。Gradle 冲突处理不同于 Maven：官方文档说明 Gradle 处理 version conflicts 与 capability conflicts；默认情况下，版本冲突会考虑图中所有请求版本并选择最高版本，而 capability conflict 通常会使构建失败，除非消费者显式选择（https://docs.gradle.org/current/userguide/dependency_constraints_conflicts.html）。

Gradle variant model 改变了“包兼容”的含义。官方 variant-aware resolution 文档说明 Gradle 通过匹配 consumer attributes 与 producer variants 的 attributes 来选择变体，variant 名称主要用于调试和错误信息，不参与匹配，真正参与匹配的是 attributes（https://docs.gradle.org/current/userguide/variant_aware_resolution.html）。因此 `com.example:lib:1.0` 对消费者是否兼容，不只取决于版本号，还取决于是否存在满足 consumer usage、category、platform、JVM version、library elements、Kotlin platform 等属性的 producer variant。

Kotlin Multiplatform 要求 variant-aware metadata 的原因正在这里。一个 KMP library 可能同时发布 JVM jar、JS artifact、Native klib、metadata artifact 与 Android artifact。消费者在 `commonMain`、`jvmMain`、`iosArm64Main` 等源集中的依赖请求不相同；如果只用 Maven POM 的单一依赖列表表达，无法充分描述平台 target 与变体选择。Kotlin 发布文档说明插件自动为 target 创建 publications（https://kotlinlang.org/docs/multiplatform/multiplatform-publish-lib-setup.html），Gradle Module Metadata 文档又说明 GMM 支持 multi-platform 与 variant-aware 发布（https://docs.gradle.org/current/userguide/publishing_gradle_module_metadata.html）。二者合起来说明 KMP 的包兼容性必须被建模为“组件版本 + 平台变体 + 属性匹配”，而不是单纯“坐标存在即可使用”。

NuGet 的兼容性中心变量是 target framework。NuGet TFM 文档说明 NuGet 使用 target framework references 在多个位置识别并隔离包中依赖 framework 的组件（https://learn.microsoft.com/en-us/nuget/reference/target-frameworks）。NuGet “what is NuGet” 文档将 compatible package 定义为包含至少一个与消费项目 target framework 兼容的程序集（https://learn.microsoft.com/en-us/nuget/what-is-nuget）。多目标包文档说明 NuGet 支持在单一包中放入同一 library 面向多个 .NET Framework 版本的不同程序集，并通过 `lib/<TFM>` 与 dependency groups 表达（https://learn.microsoft.com/en-us/nuget/create-packages/supporting-multiple-target-frameworks）。因此 NuGet 包评估变量必须包括项目 TFM、包内 lib/ref/runtimes 资产、dependency group 与 fallback/兼容规则。

NuGet 冲突调解也有独立规则。官方 dependency resolution 文档列出 direct dependency wins、cousin dependencies、lowest applicable version、floating versions 等规则，并指出 direct dependency wins 可能导致 downgrade，NuGet 会发出警告（https://learn.microsoft.com/en-us/nuget/concepts/dependency-resolution）。这与 Maven nearest definition 和 Gradle highest-version 默认策略不同。

## 7. Reproducibility, lock files, and audits

Maven 可复现性主要依赖显式版本、dependencyManagement/BOM、仓库不可变 release、插件版本固定和 reproducible builds 配置，而不是官方内建 lockfile。Maven reproducible builds 文档说明可通过升级支持可复现构建的插件，并在 POM 中设置 `project.build.outputTimestamp` 启用插件层面的可复现输出；该文档也提供 `mvn artifact:compare` 检查两次构建输出是否一致（https://maven.apache.org/guides/mini/guide-reproducible-builds.html）。Maven Central 上传指南说明 release 上传后不会改变或替换，这提供了仓库层面的不可变性前提（https://maven.apache.org/repository/guide-central-repository-upload.html）。但这不等于存在标准锁文件；如果使用 version ranges、SNAPSHOT 或不固定插件，解析结果和产物仍可能变化。

Gradle 有明确 dependency locking。官方文档说明 dependency locking 会保存 resolved versions 到 lock file，确保后续构建使用相同依赖版本，并防止 dependency graph 意外变化（https://docs.gradle.org/current/userguide/dependency_locking.html）。Gradle 还提供 dependency verification：官方文档说明创建 `gradle/verification-metadata.xml` 后，Gradle 可验证 checksums 与 signatures；该机制覆盖 JAR、ZIP、POM、Ivy descriptors、Gradle Module Metadata、plugins 和高级解析 API 获得的 artifacts（https://docs.gradle.org/current/userguide/dependency_verification.html）。因此 Gradle 侧需要区分“版本锁定”和“来源/内容验证”两个不同机制。

NuGet 的锁文件是 `packages.lock.json`。`PackageReference` 文档说明如果项目存在 lock file，NuGet restore 会使用它；没有定义依赖变化时不会重新评估 package dependencies；在 locked mode 下，restore 要么恢复锁文件列出的 exact packages，要么失败（https://learn.microsoft.com/en-us/nuget/consume-packages/package-references-in-project-files）。NuGet 官方博客还说明 lock file 会持久化所有 packages 的 content hash，restore 时会校验实际恢复包的 hash，防止同 ID 同版本但内容不同的 package content mismatch（https://devblogs.microsoft.com/dotnet/enable-repeatable-package-restores-using-a-lock-file/）。

安全审计方面，Maven Central 和 NuGet.org 暴露面不同。Maven Central 的发布侧要求 PGP signatures、POM、sources/javadoc 等规则（https://maven.apache.org/repository/guide-central-repository-upload.html），Central Portal 自身也提供 search 与安全相关聚合面（https://central.sonatype.com/）。Maven/Gradle 项目常借助第三方 SCA、IDE inspection 或 Gradle dependency verification 完成更严格审查。JetBrains IDEA 文档说明 IDE 可以高亮 `pom.xml` 或 `build.gradle` 中被认为 vulnerable 的 packages（https://www.jetbrains.com/help/idea/package-analysis.html）。

NuGet audit 是 NuGet 内建 restore 审计面。Microsoft 文档定义包管理器安全审计为分析项目中包及其依赖的安全性、识别漏洞和风险并提供缓解建议的过程（https://learn.microsoft.com/en-us/nuget/concepts/auditing-packages）。NuGet vulnerability info API 文档说明 NuGet client 可下载已知包漏洞信息，用于 restore 时检查包；package metadata resource 也包含已知漏洞信息（https://learn.microsoft.com/en-us/nuget/api/vulnerability-info）。.NET SDK 文档说明 `dotnet restore` 在 .NET 8 SDK 及以后默认产生安全漏洞警告（https://learn.microsoft.com/en-us/dotnet/core/compatibility/sdk/8.0/dotnet-restore-audit）。因此 .NET 侧的统一模型应把 audit source、advisory URL、severity、direct/transitive mode、suppressions 与 warning/error 策略纳入变量。

## 8. Variables and indicators useful for the unified model

统一模型应至少包含以下变量：

1. 身份变量：Maven/Gradle 的 `groupId`/`artifactId` 或 group/name、version、type、classifier；NuGet 的 package ID、version、package type。Maven 的 identity 是分层命名空间加 artifact 名，NuGet 的 identity 是扁平 package ID 加版本。

2. 发布位置变量：repository URL、repository type、namespace ownership、package source、source mapping、仓库不可变性、是否中心仓库、是否第三方镜像或私有 feed。

3. 元数据变量：POM、Gradle Module Metadata、`.nuspec`、MSBuild project package properties、README/license/repository/SCM/developers/tags、dependency groups、dependency scopes、optional/exclusion、capabilities、attributes。

4. 依赖边变量：direct/transitive、scope/configuration、dependencyManagement/BOM/platform/constraints、version range/floating version、optional/private assets、excluded dependencies、development dependency、runtime assets。

5. 解析规则变量：Maven nearest definition 与同深度 first declaration wins；Gradle 默认 highest requested version、rich versions、constraints、capability conflict；NuGet direct dependency wins、lowest applicable version、cousin dependencies、floating versions。

6. 兼容性变量：JVM bytecode level、Java/Kotlin platform、Gradle attributes、KMP target/source set、Android variant、NuGet TFM、包内 asset folder、runtime identifier、dependency group target framework。

7. 可复现变量：显式版本、动态版本使用、Maven outputTimestamp、Gradle lockfile、Gradle verification metadata、NuGet `packages.lock.json`、locked mode、content hash、仓库内容是否可替换。

8. 安全变量：签名、checksum、verified namespace/owner、vulnerability advisory、severity、affected range、fixed version、direct/transitive exposure、audit source、suppression reason、IDE/SCA warning surface。

9. 发现变量：搜索关键词、分类、下载量、版本新旧、owner/publisher、verified prefix、frameworks tab、dependency graph view、IDE dependency analyzer、插件 portal plugin id。

这些变量不构成实现计划，而是跨生态描述包复用状态所需的对象字段。其核心用途是避免把“包名相同”“仓库可搜索到”“构建能下载”误认为“当前项目可兼容、可复现、可安全使用”。

## 9. Cross-ecosystem comparison: JVM/Kotlin vs .NET

JVM/Kotlin 与 .NET 的第一差异是身份结构。Maven/Gradle 坐标天然分层，`groupId`/group 承担命名空间，`artifactId`/name 承担组件名，version 承担发布版本；NuGet 以 package ID 和 version 指称具体包，组织信息更多通过 ID 前缀、owner、verified status、repository 等辅助字段表达。由此，JVM 统一模型必须保留 group 维度，不能把 artifactId 当作全局包名；.NET 模型必须把 package ID 视为全局搜索与引用入口。

第二差异是兼容性语义。传统 Maven artifact 的兼容性主要在类路径、scope 和 Java 运行/编译要求上体现；Gradle/KMP 则把兼容性前移到 variant attributes 匹配。NuGet 则以 TFM 为核心，对包内不同 framework assets 和 dependency groups 进行选择。换言之，Gradle/KMP 的“兼容”是 attribute-based variant selection，NuGet 的“兼容”是 target-framework-based asset selection，Maven POM 的“兼容”更多是坐标解析与 classpath 形成后的构建/运行结果。

第三差异是冲突调解。Maven 默认最近定义，Gradle 默认最高版本，NuGet 有 direct dependency wins 和 lowest applicable version 等规则。相同依赖图在三个工具中可能产生不同版本选择，因此统一模型不能只存 dependency edges，还必须存 resolver family 和 mediation rule。

第四差异是可复现手段。Maven 依赖显式版本和构建输出可复现配置，没有同 NuGet/Gradle 一样的一等官方 lockfile；Gradle 有 dependency locking 与 verification metadata；NuGet 有 `packages.lock.json`、locked mode 和 contentHash 校验。统一模型应把“版本锁定”“内容校验”“构建输出可复现”拆开，而不是用一个 reproducible 布尔值覆盖。

第五差异是安全暴露面。Maven Central 发布侧强调 namespace、POM、签名、sources/javadoc 和不可变 release；Gradle 侧在消费时可通过 dependency verification 强化 checksum/signature；NuGet 侧把 vulnerability audit 深度集成到 restore、Visual Studio 与 API metadata 中。统一模型应同时记录发布端信任信号和消费端审计信号。

## 10. Open uncertainties and references

开放不确定性如下。

第一，Maven Central 页面和第三方索引对漏洞的展示方式会随平台产品演进变化；本文只把 Maven Central/IDE/SCA 视为安全信息暴露面，不声称 Maven 本体具有与 NuGet Audit 等价的内建漏洞审计能力。

第二，Kotlin Multiplatform 的实际 publication layout 会随 Kotlin Gradle Plugin、Android Gradle Plugin 和目标平台组合变化；本文确认的稳定结论是 KMP 需要 target-aware/variant-aware metadata，具体 artifact 形态仍应以发布项目生成的 `.module`、POM 和仓库文件为准。

第三，NuGet target framework 兼容性规则涉及 TFM 等价、fallback、asset selection、RID-specific assets 等更多细节；本文只把 TFM 作为统一模型的必要变量，不展开全部 NuGet resolver 算法。

第四，IDE 包管理器是发现和分析界面，不是包身份或解析规则的唯一权威来源；同一 IDE 在不同版本和插件组合下展示信息可能不同。

参考 URL：

- Apache Maven POM Reference: https://maven.apache.org/pom.html
- Apache Maven Introduction to the POM: https://maven.apache.org/guides/introduction/introduction-to-the-pom.html
- Apache Maven Naming Convention of Maven Coordinates: https://maven.apache.org/guides/mini/guide-naming-conventions.html
- Apache Maven Dependency Mechanism: https://maven.apache.org/guides/introduction/introduction-to-dependency-mechanism.html
- Apache Maven Central Repository Upload Guide: https://maven.apache.org/repository/guide-central-repository-upload.html
- Apache Maven Reproducible Builds Guide: https://maven.apache.org/guides/mini/guide-reproducible-builds.html
- Sonatype Central Publishing Requirements: https://central.sonatype.org/publish/requirements/
- Sonatype Central Portal: https://central.sonatype.com/
- Maven Central Search: https://search.maven.org/
- MvnRepository third-party discovery surface: https://mvnrepository.com/
- Gradle Dependency Management: https://docs.gradle.org/current/userguide/getting_started_dep_man.html
- Gradle Dependency Resolution: https://docs.gradle.org/current/userguide/dependency_resolution.html
- Gradle Dependency Constraints and Conflict Resolution: https://docs.gradle.org/current/userguide/dependency_constraints_conflicts.html
- Gradle Module Metadata: https://docs.gradle.org/current/userguide/publishing_gradle_module_metadata.html
- Gradle Variant-Aware Resolution: https://docs.gradle.org/current/userguide/variant_aware_resolution.html
- Gradle Variants and Attributes: https://docs.gradle.org/current/userguide/variant_attributes.html
- Gradle Dependency Locking: https://docs.gradle.org/current/userguide/dependency_locking.html
- Gradle Dependency Verification: https://docs.gradle.org/current/userguide/dependency_verification.html
- Gradle Plugin Portal: https://plugins.gradle.org/
- Gradle Plugin Publishing: https://docs.gradle.org/current/userguide/publishing_gradle_plugins.html
- Kotlin Multiplatform Publishing: https://kotlinlang.org/docs/multiplatform/multiplatform-publish-lib-setup.html
- NuGet overview: https://learn.microsoft.com/en-us/nuget/what-is-nuget
- NuGet.org: https://www.nuget.org/
- NuGet Package Versioning: https://learn.microsoft.com/en-us/nuget/concepts/package-versioning
- NuGet `.nuspec` reference: https://learn.microsoft.com/en-us/nuget/reference/nuspec
- NuGet PackageReference: https://learn.microsoft.com/en-us/nuget/consume-packages/package-references-in-project-files
- NuGet Dependency Resolution: https://learn.microsoft.com/en-us/nuget/concepts/dependency-resolution
- NuGet Target Frameworks: https://learn.microsoft.com/en-us/nuget/reference/target-frameworks
- NuGet multi-target packages: https://learn.microsoft.com/en-us/nuget/create-packages/supporting-multiple-target-frameworks
- NuGet lock file announcement: https://devblogs.microsoft.com/dotnet/enable-repeatable-package-restores-using-a-lock-file/
- NuGet Auditing Packages: https://learn.microsoft.com/en-us/nuget/concepts/auditing-packages
- NuGet Vulnerability Info API: https://learn.microsoft.com/en-us/nuget/api/vulnerability-info
- .NET restore audit behavior: https://learn.microsoft.com/en-us/dotnet/core/compatibility/sdk/8.0/dotnet-restore-audit
- Visual Studio NuGet Package Manager: https://learn.microsoft.com/en-us/nuget/consume-packages/install-use-packages-visual-studio
- IntelliJ Maven dependencies: https://www.jetbrains.com/help/idea/work-with-maven-dependencies.html
- IntelliJ Gradle dependencies: https://www.jetbrains.com/help/idea/work-with-gradle-dependency-diagram.html
- IntelliJ package vulnerability analysis: https://www.jetbrains.com/help/idea/package-analysis.html
