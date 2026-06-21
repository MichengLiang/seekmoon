# C/C++ 与 Swift 包复用机制研究报告

## 1. 范围与对象边界

本文只研究 C/C++ 与 Swift 生态中“包如何被声明、发现、解析、锁定并实现为可消费构件”的机制，不研究具体项目迁移策略，也不提供代码或实施计划。C/C++ 侧以 Conan 2、ConanCenter、vcpkg manifest mode、vcpkg registries、Conan lockfiles、Conan profiles、Conan settings/options、vcpkg triplets 为对象。Swift 侧以 Swift Package Manager、`Package.swift`、`PackageDescription`、Swift Package Registry Service 与 Swift Package Index 为对象。

对象边界必须先区分两类表面。第一类是包管理器的执行机制：Conan 与 vcpkg 负责根据配方、清单、配置、版本约束、变体变量和远端/注册表状态生成依赖图并取得源码或二进制构件；SwiftPM 负责读取 `Package.swift`、解析依赖版本、生成 `Package.resolved` 并构建 targets/products。第二类是发现与评估表面：ConanCenter、vcpkg.io 包页面与 Swift Package Index 帮助人寻找包、查看元数据、版本、兼容性或维护信号，但它们本身不等同于包管理器的解析算法。

本文的核心判断是：C/C++ 包身份不能只用“名称 + 版本 + 源码位置”描述。Conan 官方二进制模型明确把 `package_id` 视为二进制包标识，并由 settings、options 与依赖版本等配置参与计算；其文档说明 `package_id` 是通过对配置进行哈希来计算的包二进制标识，配置包括 settings、options 和 dependency versions（https://docs.conan.io/2/reference/binary_model.html）。vcpkg 的 triplet 概念也把目标环境集中为 CPU、OS、compiler、runtime 等变量的组合（https://learn.microsoft.com/en-us/vcpkg/concepts/triplets）。因此，在 C/C++ 中，包复用评估必须把 ABI、平台、编译器、运行时、构建类型、链接方式和选项作为一等变量，而不是把它们当作安装后的附属信息。

## 2. C/C++ 生命周期映射

Conan 的生命周期从 recipe 开始。`conanfile.py` 中的 `name`、`version`、可选的 `user/channel` 形成包引用；官方文档说明 ConanCenter recipes 不使用 `user/channel`，通常为 `pkg/version` 形式（https://docs.conan.io/2/reference/conanfile/attributes.html）。依赖可以通过 `requires`、`tool_requires`、`test_requires` 等属性或方法声明，其中 `requires` 表示 host context 中的普通库依赖，也支持版本范围（https://docs.conan.io/2/reference/conanfile/attributes.html）。同一个 recipe 在不同 settings/options 下会产生不同二进制包，最终由 `package_id` 区分。

Conan 的变体变量分为两层。`settings.yml` 提供全局设置维度，例如 `os`、`arch`、`compiler`、`compiler.version`、`compiler.libcxx`、`compiler.cppstd`、`compiler.runtime`、`compiler.runtime_type`、`build_type` 等；这些变量覆盖操作系统、体系结构、编译器、C/C++ 标准库、C/C++ 标准、运行时链接和 Debug/Release 等 ABI 相关条件（https://docs.conan.io/2/reference/config_files/settings.html）。recipe 的 `options` 则表达当前包自己的可配置特征，例如 `shared`、`fPIC` 或项目自定义选项；Conan 文档明确说明 option 值变化默认会改变 `package_id`（https://docs.conan.io/2/reference/conanfile/attributes.html）。profiles 把 settings、options、环境变量、tool requirements、configuration variables 等聚合为可复用配置文件（https://docs.conan.io/2/reference/config_files/profiles.html）。

vcpkg manifest mode 的生命周期从 `vcpkg.json` 开始。Microsoft 文档说明 manifest mode 使用声明式 JSON 文件描述项目或包元数据，文件名必须为 `vcpkg.json`；项目 manifest 的主要用途是声明依赖，并可以使用版本约束和 overrides 锁定特定版本（https://learn.microsoft.com/en-us/vcpkg/concepts/manifest-mode）。`vcpkg.json` 的顶层字段包括 `dependencies`、`features`、`default-features`、`builtin-baseline`、`overrides`、`supports`、`version` 与 `port-version` 等；其中 dependencies 可为字符串或对象，依赖对象可声明 `name`、`features`、`default-features`、`host`、`platform`、`version>=`（https://learn.microsoft.com/en-us/vcpkg/reference/vcpkg-json）。

vcpkg 的变体核心是 triplet。官方 triplet 文档把 triplet 定义为跨编译语境中完整捕捉 target environment 的标准术语，内容包括 CPU、OS、compiler、runtime 等；vcpkg 中一次构建最多消耗 target triplet 与 host triplet 两个 triplet（https://learn.microsoft.com/en-us/vcpkg/concepts/triplets）。这说明 vcpkg 的包身份不是纯源码身份，而是“端口 + 版本 + 特性 + 目标/宿主 triplet + 工具链与链接配置 + 注册表版本状态”的组合。vcpkg 的 feature 是构建行为和依赖的布尔扩展项，官方文档称 features 会添加额外行为和依赖（https://learn.microsoft.com/en-us/vcpkg/reference/vcpkg-json）。

## 3. Swift 生命周期映射

SwiftPM 的生命周期从 `Package.swift` manifest 开始。Swift 官方 `PackageDescription` 文档将 `Package` 定义为 Swift package 的配置对象，并说明初始化参数提供 package name、targets、products、dependencies 和其他配置（https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html）。`Package.swift` 必须以 `// swift-tools-version:` 开头；该 tools version 声明 `PackageDescription` 库版本、处理 manifest 所需的最低 Swift tools 版本和语言兼容版本（https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html）。

Swift 包结构的关键对象是 package、product、target、dependency 与 supported platform。Product 是对外可见的 build artifact，由一个或多个 target 的构建产物组装而成；Swift 官方文档把 product 分为 library 与 executable，并说明 library product 可静态或动态链接，未显式指定时由 SwiftPM 按消费者偏好选择（https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html）。Target 是 Swift 包的基本构建单元，每个 target 包含一组源码，被编译为 module 或 test suite；target 可以依赖同包 target，也可以依赖外部 package dependency 提供的 product（https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html）。

SwiftPM 的 dependency 在常规 Git 模式下由源码仓库 URL 与版本 requirement 组成。官方文档说明 SwiftPM 会执行 dependency resolution 以确定可用的精确依赖版本，解析结果记录在顶层 `Package.resolved` 文件中（https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html）。版本规则以 SemVer 为基础，`Version` 被定义为符合 semantic versioning 的版本；常见 requirement 包括从某版本到下一 major、到下一 minor、精确版本、范围、branch 与 revision，但官方文档明确提示 branch/commit dependency 不适合作为发布版本依赖的公共表面（https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html）。

Swift 的 platform 约束在 `platforms` API 中声明。官方文档说明 SwiftPM 会为支持平台分配默认最低部署版本；如果依赖的 deployment target 与顶层 package 的 deployment version 不兼容，SwiftPM 会报错，且依赖的部署目标必须小于或等于顶层包在该平台上的部署目标（https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html）。这使 SwiftPM 的兼容性评估主要围绕 Swift tools version、语言版本、平台最低部署版本、target/product graph 与 SemVer 版本关系展开，而不是围绕 C++ ABI 矩阵展开。

## 4. 元数据、变体与兼容性表面

C/C++ 的元数据层与变体层必须分离。名称、版本、license、homepage、description、topics 等帮助识别上游项目和发现包；Conan 文档说明 `topics` 可用于 ConanCenter 搜索过滤，license 推荐使用 SPDX 标识（https://docs.conan.io/2/reference/conanfile/attributes.html）。但这些字段不能决定二进制可用性。二进制可用性由 settings/options、依赖版本、构建配置与 ABI 规则决定。Conan 的 `settings.yml` 显式列出编译器家族、版本、标准库、runtime、runtime_type、cppstd/cstd 等变量；这就是 ABI/platform/compiler 变量在 C/C++ 中成为一等变量的直接证据（https://docs.conan.io/2/reference/config_files/settings.html）。

vcpkg 同样把发现元数据与构建变体分开。`vcpkg.json` 中的 description、homepage、license、maintainers 支持理解与发现；dependencies、features、supports、host、platform、version>=、builtin-baseline、overrides 进入解析和选择过程（https://learn.microsoft.com/en-us/vcpkg/reference/vcpkg-json）。Triplet 则负责表达目标环境集合，且允许自定义 triplet 覆盖或扩展默认配置（https://learn.microsoft.com/en-us/vcpkg/concepts/triplets）。因此，vcpkg 中“同名同版本端口”在不同 triplet、features、host/target 组合下不应被统一视为同一个可复用 artifact。

Swift 的元数据与兼容性表面更接近源码包模型，但仍存在构建条件。`Package.swift` 中 `name`、products、targets、dependencies、platforms、swiftLanguageVersions、cLanguageStandard、cxxLanguageStandard 等构成包声明（https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html）。Target 可声明 C/C++/Swift/linker settings；binary target 可引用预构建二进制 artifact，并要求 URL 与 checksum，且官方文档说明 binary targets 只支持 Apple platforms（https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html）。这说明 SwiftPM 并非完全没有二进制构件，但 Swift 源码包的主身份仍由 manifest、URL 或 registry identifier、SemVer release、products/targets 与 platform constraints 主导。

## 5. 发现表面与数据来源

ConanCenter 是 Conan 生态的中心发现表面。ConanCenter 页面自称为 Conan libraries and tools central repository，并展示 recipes、references、热门 recipes、最新 revisions 与新 upstream releases（https://conan.io/center）。ConanCenter 的作用是让使用者发现 recipe、版本和相关元数据；实际解析、下载、构建和二进制选择仍由 Conan client、remotes、recipe、profile、settings/options 与 lockfile 共同完成。

vcpkg 的发现表面包括 vcpkg.io 包浏览页、vcpkg curated registry 和 registry 文件结构。Microsoft 文档说明 vcpkg 在 GitHub 上维护一组 ports，称为 curated registry，但 vcpkg 不限于 curated registry，用户可创建 custom registries；registry 是按特定结构组织的 ports 与 helper files 集合（https://learn.microsoft.com/en-us/vcpkg/concepts/registries）。vcpkg.io 的包浏览页提供官方包搜索入口（https://vcpkg.io/en/packages）。这里也要区分：vcpkg.io 是发现页面，registry 是解析时可被配置和读取的数据源，`vcpkg-configuration.json` 或 manifest 内嵌 configuration 才决定当前项目使用哪些 registry。

Swift Package Index 是第三方发现与评估索引，不是 SwiftPM 的官方注册表解析机制。其 FAQ 明确称自己是支持 Swift Package Manager 的 packages 搜索引擎，并强调选择依赖不只是找到代码，还包括维护、开发历史、测试等评估问题（https://swiftpackageindex.com/faq）。SPI 的数据来自一个 repository list；当仓库被加入列表时，站点会完整 clone 源码，并从源码、Git 历史和 GitHub 托管仓库提取元数据，之后每几小时轮询变化（https://swiftpackageindex.com/faq）。SPI 还运行真实构建，用多种 Swift compiler 和平台评估 Swift version compatibility 与 platform compatibility；其 FAQ 明确说这些兼容性数据不能从 manifest 单独推导（https://swiftpackageindex.com/faq）。

## 6. 依赖解析、锁定与 artifact 实现

Conan lockfile 的对象是依赖解析结果的可复现性。官方文档说明 lockfiles 用于在依赖新版本或新 revisions 被创建后仍实现 reproducible dependencies；示例中 `conan.lock` 捕获已解析的 `matrix/1.0` 及其 revision，使后续即使范围内出现 `matrix/1.1` 也继续使用锁定版本（https://docs.conan.io/2/tutorial/versioning/lockfiles.html）。这说明 Conan 的锁定不只是包名版本，也包括 recipe revision，并且锁定内容受配置影响；同一依赖图在不同架构条件下可能出现条件依赖，锁文件必须覆盖相应配置才完整。

vcpkg 的复现机制主要通过 manifest versioning、registry baseline 与 overrides 表达。Microsoft versioning 文档说明 versioning 允许在 manifest 文件内确定性控制项目使用的依赖 precise revisions，并且只适用于 manifest mode（https://learn.microsoft.com/en-us/vcpkg/users/versioning）。`builtin-baseline` 是默认 registry 版本解析的 baseline shortcut，代表提供全局 minimum version information 的 vcpkg commit；`overrides` 是顶层项目用于强制特定依赖版本的精确 pin（https://learn.microsoft.com/en-us/vcpkg/reference/vcpkg-json）。registry 的 versions database 包含 baseline file 与 version files，version files 列出可用版本并指向检索位置；Git registry 的 version entry 包含 `git-tree`，即 port 目录内容 hash，可从 Git 历史中取回旧版本 port files（https://learn.microsoft.com/en-us/vcpkg/concepts/registries）。

vcpkg 的 artifact realization 还包括 binary caching。Microsoft 文档说明 binary caching 通过 `VCPKG_BINARY_SOURCES` 或 `--binarysource` 配置，可使用 files、NuGet、HTTP、Azure Blob、GCS、AWS S3 等来源，并且 read/write/readwrite 控制下载与上传行为（https://learn.microsoft.com/en-us/vcpkg/reference/binarycaching）。该机制不是包身份本身，但它会影响同一解析结果是否从缓存取得二进制、是否需要本地构建，以及组织内部如何复用构建产物。

SwiftPM 的锁定表面是 `Package.resolved`。Swift 官方文档说明 `Package.resolved` 记录 dependency resolution 的结果（https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html）。Swift Package Registry Service 的 SE-0292 进一步扩展了 registry 下载路径：registry 模式下 SwiftPM 可通过 scoped identifier `scope.package-name` 解析依赖，先请求 `/{scope}/{name}` 获取 releases，再取对应版本 manifest，解析完成后下载 `/{scope}/{name}/{version}.zip` 源码归档（https://github.com/swiftlang/swift-evolution/blob/main/proposals/0292-package-registry-service.md）。SE-0292 还要求 registry release archive 使用 checksum；若服务器提供 checksum 与 `Package.resolved` 中已有 checksum 或下载产物 checksum 不匹配，SwiftPM 会拒绝下载依赖（https://github.com/swiftlang/swift-evolution/blob/main/proposals/0292-package-registry-service.md）。这提供了 Swift registry 模式下的完整性约束。

## 7. 统一模型中有用的变量与指标

统一模型不应只保存“生态、名称、版本、URL”。对 C/C++，模型至少需要以下变量：包管理器；包引用或端口名；上游版本；包装修订或 recipe revision/port-version；registry/remote 身份；registry baseline 或 remote revision 状态；依赖声明；版本范围与 override；host dependency/tool dependency；target dependency；feature/options；settings；profile；triplet；target OS；target arch；host OS/arch；compiler family；compiler version；C/C++ standard；standard library；runtime linkage；runtime type；build type；shared/static link mode；fPIC 或等价 ABI 影响项；dependency graph hash 或 lock entry；binary cache key；source checksum 或 archive hash；license 与 homepage 等发现元数据。

对 Swift，模型至少需要：package identity 类型，即 URL-derived identity 或 registry scoped identifier；repository URL 或 registry scope/name；SemVer release；tools version；PackageDescription API version；products；product type；target list；target dependency graph；target type，包括 regular、executable、test、system library、binary；platform minimum deployment versions；Swift language version；C/C++ language standard；binary target checksum；dependency requirement 类型，包括 range、exact、branch、revision；`Package.resolved` pin；registry URL 与 scope mapping；source archive checksum；SPI 兼容性指标，包括 Swift compiler matrix、platform build results、license、last activity、stars、product filter、keyword 等发现信号。

统一模型还应保留“声明兼容性”和“实测兼容性”的区别。vcpkg 的 `supports` 与 SwiftPM 的 `platforms` 是声明或约束；Swift Package Index 的 platform compatibility 来自真实构建结果（https://swiftpackageindex.com/faq）。Conan 和 vcpkg 的二进制可用性则取决于特定 settings/options/triplet 下是否已有 binary package 或 binary cache 命中，以及是否能从源码构建。

## 8. 跨生态比较：C/C++ vs Swift

C/C++ 与 Swift 的根本差异在于 artifact identity 的层位。Swift 源码包通常先以 source package identity 进入解析：URL 或 registry identifier、SemVer 版本、manifest 中的 products/targets/platforms 决定依赖图；构建产物由当前工具链生成。C/C++ 包复用则经常直接面对二进制可复用问题：同名同版本源码在不同 compiler、libc++/libstdc++、MSVC runtime、Debug/Release、shared/static、架构、操作系统和依赖 ABI 下可能不兼容。Conan 将这些变量纳入 `package_id`；vcpkg 将 target environment 聚合到 triplet，并通过 features、host/target triplet、versioning 和 binary cache 组织构建结果。

Conan 与 vcpkg 的差异也很明显。Conan 的 recipe 模型直接暴露 settings/options/profiles，并以 `package_id` 管理二进制包；lockfile 明确捕获依赖版本与 revision。vcpkg 的 port/manifest/registry 模型更强调 curated registry、baseline、version database、features 与 triplets；版本复现主要通过 manifest 中的 baseline、version constraints 与 overrides 控制，二进制复用通过 binary caching 配置。二者都不能被简化为 npm/PyPI 风格的“从注册表下载一个源包版本”。

SwiftPM 与 Swift Package Index 的关系也必须分开。SwiftPM 是官方包管理器和解析器，`Package.swift` 与 `Package.resolved` 是执行机制中的核心文件。Swift Package Index 是发现、搜索和兼容性评估服务，提供搜索过滤、维护信号、构建矩阵与文档托管等信息，但不替代 SwiftPM 的依赖解析。Swift Package Registry Service 则是 SwiftPM 可集成的 registry protocol，其对象是从 registry 下载 releases、manifest 与 source archives，而不是搜索排名或项目质量评估。

## 9. 开放不确定性与参考资料

第一，vcpkg 的公开文档强调 manifest versioning、baseline、overrides、registry version database 与 binary caching；本文不把 vcpkg 内部安装树中的临时状态文件作为稳定公共 lockfile 契约处理，因为任务要求优先使用官方/primary sources，而 Microsoft Learn 的稳定复现表面主要是 versioning 与 registry baseline。第二，Swift Package Registry Service 已在 SE-0292 标为 implemented in Swift 5.7，但不同 registry 实现、私有 registry 的可用性、生态采用程度不在本文证据范围内。第三，Swift Package Index 的 build compatibility 是重要发现信号，但其结果依赖其构建矩阵和调度策略；本文只使用其 FAQ 中明确披露的行为，不推断生态流行度。

参考资料：

- Conan 2 Binary model: https://docs.conan.io/2/reference/binary_model.html
- Conan 2 conanfile attributes: https://docs.conan.io/2/reference/conanfile/attributes.html
- Conan 2 settings.yml: https://docs.conan.io/2/reference/config_files/settings.html
- Conan 2 profiles: https://docs.conan.io/2/reference/config_files/profiles.html
- Conan 2 lockfiles: https://docs.conan.io/2/tutorial/versioning/lockfiles.html
- ConanCenter: https://conan.io/center
- vcpkg manifest mode: https://learn.microsoft.com/en-us/vcpkg/concepts/manifest-mode
- vcpkg.json reference: https://learn.microsoft.com/en-us/vcpkg/reference/vcpkg-json
- vcpkg registries concepts: https://learn.microsoft.com/en-us/vcpkg/concepts/registries
- vcpkg triplets: https://learn.microsoft.com/en-us/vcpkg/concepts/triplets
- vcpkg versioning reference: https://learn.microsoft.com/en-us/vcpkg/users/versioning
- vcpkg binary caching: https://learn.microsoft.com/en-us/vcpkg/reference/binarycaching
- vcpkg package browser: https://vcpkg.io/en/packages
- Swift PackageDescription: https://docs.swift.org/package-manager/PackageDescription/PackageDescription.html
- Swift Evolution SE-0292 Package Registry Service: https://github.com/swiftlang/swift-evolution/blob/main/proposals/0292-package-registry-service.md
- Swift Package Index FAQ: https://swiftpackageindex.com/faq
- Swift Package Index Add a Package: https://swiftpackageindex.com/add-a-package
