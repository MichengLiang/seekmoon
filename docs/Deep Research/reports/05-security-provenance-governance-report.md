# Security, Provenance, and Governance Across Package Ecosystems 研究报告

## 1. Scope and object boundary

本文讨论的对象是跨包生态的安全、来源证明、发布身份、漏洞状态与注册表治理信号。对象边界不是“某个包是否安全”这一单一结论，而是包在发现、安装、解析、发布、撤回、维护和治理过程中可以被观察、验证或记录的不同事实层。必须区分五类对象：漏洞状态、发布者身份、构建来源证明、制品完整性、项目治理。

漏洞状态描述某个包名、版本、提交或版本范围是否被已知公告覆盖。OSV 将漏洞记录表达为可机器读取的格式，目标是使漏洞能精确映射到开源包版本或提交哈希，并通过 `affected`、`ranges`、`versions`、`severity`、`withdrawn` 等字段表达影响范围、严重度和撤回状态（https://ossf.github.io/osv-schema/）。GitHub Advisory Database 也区分 GitHub-reviewed advisories、unreviewed advisories 和 malware advisories，并说明恶意软件与普通漏洞的处理语义不同（https://docs.github.com/en/code-security/concepts/vulnerability-reporting-and-management/github-advisory-database）。

发布者身份描述谁或哪个自动化环境被注册表接受为发布主体。npm trusted publishing 使用 OIDC 在 npm 与 CI/CD provider 之间建立信任关系，允许特定 workflow 直接发布包，避免长期 npm token（https://docs.npmjs.com/trusted-publishers/）。PyPI Trusted Publishing 同样以 OIDC 在受信第三方服务和 PyPI 之间交换短生命周期身份 token，并由 PyPI 铸造短生命周期、项目作用域 API token（https://docs.pypi.org/trusted-publishers/）。

构建来源证明描述制品从哪个源码、哪个构建过程、哪个构建平台、哪些输入产生。SLSA v1.2 将供应链安全组织为不同 track 和 level，其中 Build track 关注构建过程的来源证明和防篡改保证（https://slsa.dev/spec/v1.2/）。SLSA provenance predicate 用 `subject`、`buildDefinition`、`externalParameters`、`resolvedDependencies`、`runDetails.builder.id` 等字段描述软件制品在哪里、何时、如何产生（https://slsa.dev/spec/v1.2/build-provenance）。

制品完整性描述下载到的包文件是否对应预期摘要、签名或证明材料。npm provenance 文档说明，当包带有 provenance 时，npm 会通过 Sigstore 公共服务签名并记录到透明日志；该证明并不保证无恶意代码，而是提供可验证的源码和构建指令链接（https://docs.npmjs.com/generating-provenance-statements/）。PyPI digital attestations 将每个 release distribution 绑定到内容摘要，使 PyPI 和下游用户能够验证某个包文件由特定身份完成 attestation（https://docs.pypi.org/attestations/）。

项目治理描述命名空间、所有权、维护状态、权限分配、2FA、token、撤回、yank、deprecation、name transfer 等注册表制度。治理信号不是密码学证明，但会影响包发现和风险解释。npm scopes 允许用户或组织拥有自己的 scope，只有该用户或组织可以在该 scope 下添加包（https://docs.npmjs.com/about-scopes/）。PyPI name retention 明确包索引服务社区，保留已发布内容的历史价值，但在废弃项目、名称争议、无效项目、恶意软件或侵权场景中允许转移或移除（https://docs.pypi.org/project-management/name-retention/）。

## 2. Security/provenance lifecycle placement

安全和来源证明信号应按生命周期归位。发现阶段可以观察包名、命名空间、主页、仓库链接、维护活动、license、最新版本、deprecation/yank 状态、是否存在公开 provenance、是否存在已知 advisory、是否有 OpenSSF Scorecard 结果。这个阶段的判断还没有执行安装，也没有解析完整依赖图，因此只能处理包级元数据、项目级元数据和公开公告。

安装阶段可以观察解析后的具体版本、传递依赖、lockfile、registry 返回的文件摘要、签名或 registry integrity 信息、安装脚本、平台条件、peer/optional dependency 选择、审计报告。npm audit 明确说明，它会提交当前包依赖描述到默认 registry，请求已知漏洞报告，并在 `npm install` 时自动运行；它检查 direct、dev、bundled、optional dependencies，但不检查 peerDependencies（https://docs.npmjs.com/auditing-package-dependencies-for-security-vulnerabilities/）。因此，漏洞判断在安装阶段比发现阶段更具体，因为解析结果已经确定。

管理阶段包括发布、撤回、yank、deprecate、transfer、token rotation、2FA enforcement、组织权限调整、CI trusted publisher 配置、advisory 发布与修订。很多治理事件只有项目所有者、注册表或包管理器在后续操作中才能确认。例如 PyPI yanking 是非破坏性删除替代；被 yanked 的 release 通常会被 installer 忽略，除非它是唯一匹配显式版本约束的 release（https://docs.pypi.org/project-management/yanking/）。这类状态会改变下游解析行为，不能只作为说明性文本处理。

## 3. Pre-install discovery signals

安装前可观察信号分为可直接展示、可外部聚合、可推断三类。可直接展示的信号包括包名、scope 或 namespace、版本发布时间、维护者或组织、仓库 URL、license、README、deprecation message、yanked reason、provenance badge 或 attestation availability。npm provenance 文档指出 provenance 可以在下载前帮助开发者验证包在哪里以及如何构建；同时也明确 provenance 不保证包没有恶意代码（https://docs.npmjs.com/generating-provenance-statements/）。因此，发现页可以把 provenance 作为“可审计性”信号，而不能把它显示为“安全保证”。

可外部聚合的信号包括 OSV/GHSA 漏洞公告、OpenSSF Scorecard、依赖更新工具存在性、维护活跃度、CI 与 review 机制、SBOM 可用性、security policy。OSV.dev 是多个采用 OSV 格式的漏洞数据库聚合器，并提供 API 给工具查询版本或提交是否受已知漏洞影响（https://google.github.io/osv.dev/）。GitHub Advisory Database 从 GitHub security advisories、NVD、npm Security advisories、Python Packaging Advisory Database、RustSec 等来源导入数据，并以 OSV JSON 发布 advisory（https://docs.github.com/en/code-security/concepts/vulnerability-reporting-and-management/github-advisory-database）。

可推断信号包括命名空间可信度、所有者稳定性、发布频率异常、版本年龄、包名相似度、维护者数量、组织治理结构、token 发布风险迹象。它们不能替代官方状态。比如 scope 可以降低公共命名空间被抢注的风险，但并不证明 scope 内每个包安全。npm 官方文档只说明 scope 允许创建与其他用户或组织同名但不冲突的包，并能作为组织官方包的一种信号（https://docs.npmjs.com/about-scopes/）。该事实只能转化为命名治理变量，不能转化为代码安全变量。

发现排序中应显式分离状态：存在已知未修复漏洞、存在恶意软件 advisory、包被 deprecate、release 被 yanked、存在 provenance、存在 trusted publisher、Scorecard 某项低分、维护活动停止。这些状态的含义不同。普通漏洞意味着某些版本受影响且可能有 fix；恶意软件意味着包本身被判定为有害行为，GitHub 文档明确建议移除 malware 并寻找更安全替代（https://docs.github.com/en/code-security/concepts/vulnerability-reporting-and-management/github-advisory-database）。二者不能用同一扣分项表示。

## 4. Install-time and management-time checks

安装时检查处理的是实际依赖解。发现阶段看到的是候选包；安装阶段知道的是具体版本集合、依赖边、lockfile 和实际下载文件。此时可以执行漏洞查询、文件完整性校验、签名验证、provenance 验证、yank/deprecation 决策、安装脚本策略、平台兼容性判断。

漏洞查询需要用解析后的版本范围。OSV schema 的核心价值在于精确表达受影响版本、修复版本、提交区间和别名关系（https://ossf.github.io/osv-schema/）。如果 discovery 只看包名级别 advisory，容易误判：同一个包可能只有旧版本受影响，也可能某个 advisory 已 withdrawn。OSV 的 `withdrawn` 字段表示记录被撤回；该状态应从“存在漏洞”变量中剥离，作为公告质量和公告状态变量处理。

完整性和来源证明检查处理的是下载文件与证明材料之间的关系。PyPI Integrity API 提供 `GET /integrity/<project>/<version>/<filename>/provenance`，返回某个文件的 provenance object，并说明 attestation object 针对单个文件，provenance object 将一个文件的多个 attestations 与产生它们的 identity 绑定（https://docs.pypi.org/api/integrity/）。这意味着 PyPI 证明的粒度是 release file，而不是仅项目名。安装器或管理工具如果只按项目显示“有证明”，会丢失文件级差异。

管理时检查处理发布与治理风险。npm trusted publishing 文档建议在配置 trusted publisher 后限制传统 token-based publishing access，并提供“Require two-factor authentication and disallow tokens”的发布访问设置（https://docs.npmjs.com/trusted-publishers/）。这说明 token 风险不是抽象风险，而是 registry 允许的认证路径和包设置共同产生的风险。PyPI Trusted Publishing 则把 OIDC token 交换成 15 分钟有效的项目作用域 API token，避免长期发布 token 被泄露后持续可用（https://docs.pypi.org/trusted-publishers/）。

## 5. OpenSSF Scorecard variables

OpenSSF Scorecard 是项目级自动化安全健康检查，不是包文件级证明。其 checks 文档说明每项检查包含 scoring criteria、remediation steps 和低分风险解释，并强调 checks 会持续变化（https://github.com/ossf/scorecard/blob/main/docs/checks.md）。因此 Scorecard 适合作为项目治理与开发流程变量来源，不适合作为单个包版本的完整性结论。

可转为评价变量的 Scorecard 检查包括：

- Binary-Artifacts：源仓库是否包含生成的可执行二进制，低分风险是不可审查代码。该变量对应“源码可审查性”。
- Branch-Protection：默认和 release 分支是否受保护，是否要求 review、status checks、防止 force push。该变量对应“主线变更控制”。
- CI-Tests：合并前是否运行测试。该变量对应“自动验证存在性”，但 Scorecard 明确低分不必然表示项目有风险，因为自动检测可能漏判。
- Code-Review：近期变更是否有人类 review。该变量对应“恶意或误提交进入主线的阻力”。
- Dangerous-Workflow：GitHub Actions 是否存在 `pull_request_target` 等危险模式或脚本注入风险。该变量对应“CI 被贡献者输入劫持的风险”。
- Dependency-Update-Tool：是否使用 Dependabot 或 Renovate。该变量对应“已知漏洞更新机制存在性”，但不保证 PR 被合并。
- Fuzzing、SAST：是否存在动态输入探索和静态安全分析。它们对应“缺陷发现机制”。
- Maintained：项目是否活跃维护。该变量对应“漏洞修复可能性”，但低维护不必然表示所有小型稳定包都高风险。
- Pinned-Dependencies：构建和发布过程依赖是否 pin 到 hash 或明确版本。该变量对应“构建可复现性与依赖替换风险”。
- SBOM：是否存在软件物料清单。该变量对应“依赖透明度”。
- Security-Policy：是否有漏洞报告渠道和披露说明。该变量对应“协调披露能力”。
- Signed-Releases：release artifacts 是否有签名或 SLSA provenance 文件；Scorecard 文档明确该检查不验证签名。该变量对应“证明材料存在性”，不是“证明材料有效性”。
- Token-Permissions：自动化 workflow token 是否最小权限。该变量对应“CI token 横向影响面”。
- Vulnerabilities：是否存在 open、unfixed vulnerabilities；Scorecard 使用 OSV 服务。该变量对应“已知未修复漏洞暴露”。
- Webhooks：webhook 是否配置 token 认证。该变量对应“外部回调入口认证”。

Scorecard 变量进入排序时应保持分项可解释。一个项目可能 CI 配置良好但存在未修复漏洞；也可能 provenance 完整但维护者很少。把这些变量压成单一“安全分”会隐藏风险来源。更合理的表达是分维度展示：开发流程、构建发布、依赖治理、漏洞响应、维护活跃度。

## 6. SLSA/Sigstore/provenance variables

SLSA 贡献的是供应链完整性语言和等级化要求。SLSA v1.2 明确它是描述并逐步改善供应链安全的 specification，由多个 levels 和 tracks 组成，level 越高表示安全保证越强（https://slsa.dev/spec/v1.2/）。Build L1 要求存在 provenance，显示包如何构建；L2 要求 hosted build platform 生成并签名 provenance，以防构建后篡改；L3 要求 hardened build platform，以防构建期间被篡改（https://slsa.dev/spec/v1.2/build-track-basics）。这些等级应进入“构建来源可信度”变量，而不是“漏洞状态”变量。

SLSA provenance 的可比变量包括：`subject` 是否绑定具体 artifact digest；`buildType` 是否可解释；`externalParameters` 是否完整且可验证；`resolvedDependencies` 是否记录构建时依赖；`builder.id` 是否标识受信构建平台；`metadata.startedOn/finishedOn` 是否提供构建时间；signer-builder pair 是否符合消费者预期。SLSA 文档要求消费者只接受特定 signer-builder pairs，并说明 builder 与 signer 是不同对象（https://slsa.dev/spec/v1.2/build-provenance）。该点对于跨生态评分很重要：签名者身份不等于构建平台安全边界。

Sigstore 贡献的是身份绑定签名、短生命周期证书和透明日志。Sigstore keyless signing 将身份而非长期密钥绑定到 artifact signature；Fulcio 基于 OIDC 身份签发短生命周期证书；Rekor 记录签名事件并提供可审计时间记录（https://docs.sigstore.dev/cosign/signing/overview/）。签名流程中会生成临时密钥对、获取 identity token、由 CA 验证身份并签发证书，私钥随后销毁，消费者通过透明日志条目而不是长期私钥管理来验证（https://docs.sigstore.dev/cosign/signing/overview/）。因此 Sigstore 变量应表达为“身份可追溯性”“签名事件透明度”“证明材料可验证性”，不能表达为“代码无害”。

npm 与 PyPI 的 provenance 机制在对象粒度上不同但可比较。npm provenance 包含 provenance attestation 和 publish attestation；前者提供源码和构建指令链接，后者由 registry 在授权用户发布时生成；包带 provenance 时会由 Sigstore 公共服务签名并记录到公共透明账本（https://docs.npmjs.com/generating-provenance-statements/）。PyPI attestations 基于 PEP 740 和 in-toto framework，支持 SLSA Provenance 与 PyPI Publish attestation，并把每个文件与内容 digest 和身份绑定（https://docs.pypi.org/attestations/）。可比变量包括：是否有 attestation、attestation predicate 类型、subject digest、publisher identity、CI identity、透明日志材料、registry-hosted provenance endpoint、公开仓库限制。

## 7. Registry governance variables

注册表治理变量应描述谁能声明名字、谁能发布、谁能转移、谁能撤回，以及下游如何看见这些状态。

命名空间变量包括 npm scope、组织 scope、PyPI project name、名称保留政策。npm scope 属于用户或组织，能避免其他人抢占同一 scope 下的包名，并能表示组织官方包集合（https://docs.npmjs.com/about-scopes/）。PyPI 没有 npm 同形态 scope；其治理重点是 project name retention。PyPI name retention 规定废弃项目需满足 owner 不可联系、过去十二个月无 release、项目主页无 owner 活动等条件才被视为 abandoned；即使废弃，也不会仅因废弃从索引移除，因为 artifacts 有历史价值（https://docs.pypi.org/project-management/name-retention/）。

所有权和转移变量包括 owner 数量、maintainer 数量、组织团队权限、名称争议处理、废弃项目转移条件。PyPI name transfer 需要候选者证明联系原 owner 失败、已有 fork 改进、不能用不同名称作为合理替代，且索引维护者没有额外保留意见；可联系 owner 明确反对时不得强制转移（https://docs.pypi.org/project-management/name-retention/）。这类变量反映名称稳定性和接管风险，不直接反映代码质量。

撤回、删除、yank、deprecation 变量必须按生态语义分开。PyPI yanking 是非破坏性替代删除，通常被 installer 忽略，但显式 pin 到唯一匹配版本时仍可能安装；yank reason 会显示在 release page 和 index APIs 中（https://docs.pypi.org/project-management/yanking/）。npm deprecation 是维护者对包或版本标记“不再建议使用”的机制，并允许附带消息；npm unpublish policy 则限制何时可以从 registry 删除包或版本，删除后历史可用性和名称复用风险都会改变（https://docs.npmjs.com/deprecating-and-undeprecating-packages-or-package-versions/；https://docs.npmjs.com/policies/unpublish/）。跨生态比较时，deprecation 是维护者对使用者的语义警告，yank 是解析行为信号，unpublish/delete 是可用性和历史可追溯性风险。

认证治理变量包括 2FA、token 类型、token 作用域、token 是否允许绕过 2FA、trusted publishing 是否启用、传统 token publishing 是否禁用。npm 文档说明所有包发布和设置修改要求 2FA 或启用 bypass 2FA 的 granular access token，并在 trusted publisher 场景中建议 disallow tokens（https://docs.npmjs.com/requiring-2fa-for-package-publishing-and-settings-modification/；https://docs.npmjs.com/trusted-publishers/）。PyPI 自 2024 年 1 月 1 日起要求所有用户启用 2FA，这是账号级治理基线；Trusted Publishing 则进一步减少长期发布 token 风险（https://blog.pypi.org/posts/2023-12-13-2fa-enforcement/；https://docs.pypi.org/trusted-publishers/）。

## 8. How security variables enter package-level and ecosystem-level scoring

包级评分应按对象层拆分，而不是把所有信号压成一个“安全”。可以使用以下组件语义：

- 漏洞暴露：当前候选版本是否被 OSV/GHSA 覆盖，是否有 fixed version，是否 withdrawn，是否是 malware advisory。
- 发布身份：发布是否来自 trusted publisher，publisher identity 是否可见，传统 token 发布是否仍允许。
- 构建来源：是否存在 SLSA provenance、buildType、builder.id、resolvedDependencies、subject digest，是否达到 SLSA Build L1/L2/L3 语义。
- 制品完整性：下载文件是否有 registry digest、signature、attestation、透明日志材料，证明粒度是否到文件级。
- 项目治理：Scorecard 分项、security policy、maintenance、code review、branch protection、token permissions、namespace/ownership、yank/deprecation 状态。

这些组件之间不能互相抵消。一个包存在 verified provenance 但被公告为 malware，发现排序应强烈降权或阻断，因为 provenance 只能说明恶意包来自哪里，不能使恶意行为变安全。一个包没有 provenance 但无已知漏洞，也不能被称为安全，只能称为缺少可验证构建来源。一个包 Scorecard 高分但当前版本被 yanked，解析排序仍应避免默认选择该 release，因为 yanking 是注册表对 release 级可用性的明确信号。

生态级评分应衡量生态能力，而不是单包状态。可比生态变量包括：注册表是否支持 trusted publishing；是否支持 hosted attestations 或 provenance API；是否公开 yanked/deprecated 状态；是否有恶意包报告通道；是否有 advisory database 与 OSV export；是否要求或鼓励 2FA；是否支持 scoped tokens；是否支持组织和团队权限；是否有名称争议和废弃项目转移政策；是否允许删除历史 artifacts；是否公开包级所有权变化。生态级评分的对象是“注册表和工具链能否提供可观察信号”，不是“生态中所有包都安全”。

发现排名的规则应体现风险优先级。恶意包 advisory 是最高严重度负信号；已知可利用且无修复漏洞应显著降权；yanked release 应从默认候选中排除或强降权；deprecated package 应在包级发现中降权并显示替代建议；缺少 provenance 不是直接负罪状态，但在同类候选中可降低可审计性分；Scorecard 低分应作为需要审查的项目治理信号，而不是安装阻断条件。对小型稳定项目，Maintained 低分需要结合包用途解释，不能机械惩罚。

## 9. Open uncertainties and references

开放不确定性主要有六类。

第一，跨生态 provenance 的可比粒度尚不完全一致。npm 更强调包发布与 CI provenance 展示；PyPI 的 Integrity API 明确到 release file；其他生态可能只有源码 release 签名或 registry digest。评分时需要记录证明粒度，否则会把项目级、版本级、文件级证明混淆。

第二，Scorecard 自动检测存在平台和启发式边界。其多个检查明确说明低分不一定表示项目有风险，因为某些实践可能由工具未识别。Scorecard 应作为证据变量，不应作为保证或唯一排序依据。

第三，OSV/GHSA 覆盖取决于 advisory 来源、版本映射质量和撤回状态。没有公告不等于没有漏洞；撤回公告不等于包本身安全，只能说明该 advisory 记录不再成立。恶意包记录与普通漏洞记录应分开处理。

第四，trusted publishing 降低长期发布 token 风险，但它不自动保证构建脚本、依赖、源码 review 或发布意图安全。OIDC 证明的是某个授权 workflow 运行并获得发布能力，不是证明包内容良性。

第五，yank、deprecate、unpublish 的生态语义不同。PyPI yanking 会被 index APIs 和 installer 语义使用；npm deprecation 是使用者可见警告；unpublish/delete 影响历史可用性与名称风险。跨生态模型必须保留原生态语义，再映射到共同变量。

第六，治理变量中的“官方性”难以形式化。scope、组织名、仓库 URL、publisher identity、homepage 都能提供身份线索，但它们可能被迁移、转让、误配或冒充。发现排序可以使用这些信号降低不确定性，不能把它们视为真实性证明。

参考 URL：

- OpenSSF Scorecard checks: https://github.com/ossf/scorecard/blob/main/docs/checks.md
- SLSA v1.2 specification: https://slsa.dev/spec/v1.2/
- SLSA Build Track Basics: https://slsa.dev/spec/v1.2/build-track-basics
- SLSA Build Provenance: https://slsa.dev/spec/v1.2/build-provenance
- Sigstore keyless signing overview: https://docs.sigstore.dev/cosign/signing/overview/
- npm trusted publishing: https://docs.npmjs.com/trusted-publishers/
- npm provenance statements: https://docs.npmjs.com/generating-provenance-statements/
- npm audit: https://docs.npmjs.com/auditing-package-dependencies-for-security-vulnerabilities/
- npm scopes: https://docs.npmjs.com/about-scopes/
- npm deprecating and undeprecating: https://docs.npmjs.com/deprecating-and-undeprecating-packages-or-package-versions/
- npm unpublish policy: https://docs.npmjs.com/policies/unpublish/
- npm 2FA publishing settings: https://docs.npmjs.com/requiring-2fa-for-package-publishing-and-settings-modification/
- PyPI Trusted Publishers: https://docs.pypi.org/trusted-publishers/
- PyPI Digital Attestations: https://docs.pypi.org/attestations/
- PyPI Integrity API: https://docs.pypi.org/api/integrity/
- PyPI Yanking: https://docs.pypi.org/project-management/yanking/
- PyPI Name Retention: https://docs.pypi.org/project-management/name-retention/
- PyPI 2FA enforcement: https://blog.pypi.org/posts/2023-12-13-2fa-enforcement/
- OSV documentation: https://google.github.io/osv.dev/
- OSV schema: https://ossf.github.io/osv-schema/
- GitHub Advisory Database: https://docs.github.com/en/code-security/concepts/vulnerability-reporting-and-management/github-advisory-database
