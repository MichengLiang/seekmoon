# SeekMoon Help 体系完整设计稿

以下是我准备写入 SeekMoon 代码之前的完整设计说明。本文分四层：第一层是我的注意思想；第二层是实现结构；第三层是每个命令的中文审查稿和英文公共稿；第四层是为什么这样安排。你确认之后，我再把它写进代码。

## 一、注意思想

### 1. 当前对象

当前对象不是“国际化系统”，不是“语言偏好配置”，不是“用户可切换语言的 CLI 功能”。当前对象是 SeekMoon 的 help 表面。这个表面服务第一次使用命令的 Contract reader，使使用者在执行命令前知道当前命令的动作、输入、证据边界、输出模式和下一步连接。

因此，代码里需要两份文案：中文审查稿和英文公共稿。中文审查稿服务你对语义、力度和层位的审查；英文公共稿服务最终 CLI help 输出。中文稿进入代码，不是为了运行时切换语言，而是为了让公共英文 help 有一个可审查的中文对象结构。这个边界必须保持清楚。

### 2. 正面定义优先

Help 正文优先写当前命令是什么、读取什么、写入什么、产出什么、如何进入下一步。它不以“不要做什么”作为核心句。相邻对象只在会造成高成本误读时进入，并且尽量写成对象归属句、状态归属句或证据归属句。

例如，`probe` 的核心句写“记录本地验证证据”，不写成“不要把 probe 当成质量证明”。`source` 的核心句写“定位 registry 发布版本对应的源码材料”，不围绕 repository main branch 打转。`search --target` 的核心句写“加入 target 上下文，target 证据来自已读取来源或后续命令”，不把 unsupported 这个错误对象推到首位。

### 3. Help 是行为场景

根 help、子命令 help、示例顺序、flag 文案和命令排列共同形成 CLI 行为场景。它们会诱导使用者走某条路径。当前根 help 只有默认 Cobra 列表，所以它把 SeekMoon 呈现成普通命令集合。新的 help 要把 SeekMoon 呈现成 discovery workbench：先检查环境和数据口径，再生成候选，再下钻证据，再本地验证，再记录和报告。

这不是增加教程噪音，而是让命令表面和文档中的工作台对象同形。

### 4. Help 不承担普通输出职责

普通 pretty text 继续保持低噪声；JSON、jq、shape、schema 继续承担机器处理和契约学习。Help 只在用户主动查看 `--help` 时出现，不进入正常命令结果，不污染 pipeline，不改变 JSON schema 和字段名。

### 5. 编号记忆必须进入 help

当前代码已经有 `.seekmoon/sessions/default.json` 的默认 session candidate map。`search` 和 `skill search` 产生编号候选，`view 1`、`api 1`、`source 1`、`compare 1 2`、`probe 1`、`record 1` 可以消费编号。这个能力已经在代码和测试中成立，但 help 没写。新的 help 必须把编号输入作为工作台 affordance 明确表达。

表达方式要诚实：编号来自当前项目的默认 session；失效时重新运行 search 或传完整坐标。这里不引入 session id 命令，不承诺多会话管理。

### 6. 中文审查稿与英文公共稿共享对象结构

中文稿和英文稿不逐词对应，但同一个命令必须有同一个对象结构。中文稿用于审查主语、动作、证据边界和语义强度；英文稿用于实际 Cobra `Long`、`Example`、flag usage。英文稿不能因为我对英文词感受不敏感就写得强硬、焦虑或冒犯。

### 7. 每个 help 单元的必要性

每个命令 help 只保留五类内容：命令动作、输入解释、读写边界、输出模式、下一步示例。删掉以后不影响用户正确行动的句子不进入 help。实现历史、设计者心理、理论名词、未来功能、相邻功能想象不进入 help。

## 二、实现结构设计

### 1. 新增包

我会新增：

`internal/helpdoc/doc.go`

职责：定义 help 文档数据结构和注入函数。

拟定结构：

```go
type CommandDoc struct {
    Key     string
    ShortEN string
    LongEN  string
    ExampleEN string
    ReviewZH string
    Flags map[string]FlagDoc
}

type FlagDoc struct {
    UsageEN string
    ReviewZH string
}
```

`ShortEN`、`LongEN`、`ExampleEN` 注入 Cobra。`ReviewZH` 保留中文审查稿。`Flags` 让公共 flag usage 也能被审查。

`internal/helpdoc/en.go`

职责：英文公共稿。每个命令的 Short、Long、Example、flag usage 都集中写在这里。

`internal/helpdoc/zh.go`

职责：中文审查稿。这里写每个命令的中文对象说明、边界说明和示例说明。它不参与运行时语言切换。

`internal/helpdoc/apply.go`

职责：遍历 Cobra 命令树，按 command key 注入英文 help。没有文档的公共命令会在测试中失败。

### 2. CLI 接入

在 `internal/cli/root.go` 中，创建完整命令树后调用：

```go
helpdoc.Apply(cmd)
```

它按命令路径映射文案，例如：

- `seekmoon`
- `seekmoon doctor`
- `seekmoon sync`
- `seekmoon search`
- `seekmoon view`
- `seekmoon api`
- `seekmoon source`
- `seekmoon skill`
- `seekmoon skill search`
- `seekmoon skill view`
- `seekmoon compare`
- `seekmoon probe`
- `seekmoon record`
- `seekmoon report`
- `seekmoon raw`

Cobra 默认的 `completion` 和 `help` 可以保留默认，不纳入 SeekMoon 对象文案。它们是 Cobra infrastructure，不是 SeekMoon discovery workbench 命令。

### 3. 不新增 language 命令

不会新增 `seekmoon language`。不会新增语言配置文件。不会让用户在运行时切换 help 语言。

中文稿作为代码内审查材料存在。它可以被测试读取，确认每个公共命令都有中文审查稿，但不进入默认 CLI 输出。

### 4. 测试安排

新增测试：

- 根 help 包含 workbench、first-use entry、numbered candidates、output modes。
- 每个 SeekMoon 公共命令都有 Long help 和 Example。
- 每个 SeekMoon 公共命令都有中文审查稿。
- `search --help` 包含 numbered candidates/session 的英文表达。
- `probe --help` 包含 local validation evidence 的英文表达。
- `skill --help` 或 `skill search --help` 包含 executable skill entry 的英文表达。
- `--shape`、`--schema` 仍然不执行普通命令路径。
- JSON schema id、字段名、枚举不随 help 文案变化。

## 三、通用输出模式文案

所有输出型命令共用四个 flag。当前 flag usage 太短，我会改成下面这种英文公共稿，并在中文审查稿里保存对应说明。

### 中文审查稿

`--json`：输出当前命令结果的 JSON 投影。JSON 服务脚本、CI 和报告生成，不服务终端排版。

`--jq <expr>`：先生成当前命令的 JSON 投影，再用内置 jq 求值表达式。

`--shape`：输出当前命令 JSON 投影的字段树。该模式用于写 jq、理解字段路径和学习输出契约。

`--schema`：输出当前命令 JSON 投影的 JSON Schema。该模式用于严格校验和代码生成。

### 英文公共稿

`--json`: `render the command JSON projection for scripts and automation`

`--jq`: `evaluate a jq expression against the command JSON projection`

`--shape`: `show the command JSON field tree without running the data action`

`--schema`: `show the command JSON Schema without running the data action`

## 四、根命令 help

### 中文审查稿

SeekMoon 是 MoonBit 包发现工作台。它帮助依赖消费者在引入包之前完成候选发现、证据下钻、本地验证、采纳记录和报告输出。

首次使用 SeekMoon 时，从顶层 help 进入。首次使用某个命令时，从该命令 help 进入；命令 help 说明动作、输入、证据边界、输出模式和编号候选输入。

常用工作流：先运行 `doctor` 检查本地环境，再运行 `sync` 固定数据口径，然后用 `search` 生成候选。候选可以继续进入 `view`、`api`、`source`、`compare` 和 `probe`。判断结果用 `record` 保存，用 `report` 输出。

`search` 和 `skill search` 会写入编号候选。后续命令可以用 `1`、`2` 这样的编号引用当前项目默认 session 中的候选；编号不可用时，重新运行 search 或传完整坐标。

默认输出服务终端阅读。自动化使用 `--json` 或 `--jq`。字段学习使用 `--shape`。严格校验使用 `--schema`。

### 英文公共稿

Short:

`MoonBit package discovery workbench`

Long:

```text
SeekMoon is a MoonBit package discovery workbench. It helps dependency consumers discover candidate packages, inspect evidence, run local validation, record adoption judgments, and render investigation reports before adding a dependency.

Start here when using SeekMoon for the first time. Start with a command's help before using that command for the first time; command help explains the action, inputs, evidence boundary, output modes, and numbered candidate inputs.

Common path:
  doctor -> sync -> search -> view/api/source/compare -> probe -> record -> report

search and skill search write numbered candidates into the current project's default session. Later commands can use numbers such as 1 or 2. When a number is unavailable, run search again or pass a full coordinate.

Pretty output is for terminal reading. Use --json or --jq for automation, --shape for field paths, and --schema for strict validation.
```

Example:

```text
seekmoon doctor
seekmoon sync
seekmoon search markdown --target js
seekmoon view 1
seekmoon api 1 --package moonbitlang/core/argparse
seekmoon probe 1 --target js
seekmoon record 1 --conclusion continue-verification
seekmoon report --format markdown
```

## 五、`doctor` help

### 中文审查稿

`doctor` 检查 SeekMoon 当前运行所需的本地环境。它读取 MoonBit 工具链、registry 路径、网络可达性和当前项目上下文。

`doctor` 不创建 snapshot，不更新 registry，不写 record。它的输出用于确认后续命令是否具备本地运行条件。

首次在一个项目中使用 SeekMoon 时，先运行 `doctor`。环境状态异常时，先根据错误表面处理本地工具链、路径或网络问题。

### 英文公共稿

Short:

`Check the local environment`

Long:

```text
doctor checks the local environment that SeekMoon commands depend on. It reads MoonBit toolchain availability, registry paths, network reachability, and the current project context.

doctor does not create a snapshot, update the registry, or write adoption records. Its output tells later commands which local capabilities are available.

Run doctor first in a project when using SeekMoon there for the first time.
```

Example:

```text
seekmoon doctor
seekmoon doctor --json
seekmoon doctor --shape
```

## 六、`sync` help

### 中文审查稿

`sync` 创建带时间戳的数据 snapshot。Snapshot 固定当前 Mooncakes API、统计信息、本地 index 和工具链信息的数据口径。

后续 `search`、`record` 和 `report` 可以引用 snapshot，使调查结果回到同一读取时间和来源状态。

`sync` 可以执行本地 registry 更新，并把每个来源的状态写入 snapshot。部分来源失败时，失败状态保留在 snapshot 或错误表面中。

### 英文公共稿

Short:

`Create a data snapshot`

Long:

```text
sync creates a dated evidence snapshot. The snapshot fixes the current Mooncakes API data, statistics, local registry index summary, and toolchain context used by later investigation steps.

search, record, and report can refer back to the snapshot so the investigation keeps a stable data basis.

sync records source states. When a source fails, the failure belongs to that source action and is reported through the snapshot or error surface.
```

Example:

```text
seekmoon sync
seekmoon sync --json
seekmoon sync --schema
```

## 七、`search` help

### 中文审查稿

`search` 从 query 生成 library module 候选。它读取 snapshot 中的 Modules API 数据，并按 module、description、keywords 和 repository declaration 进行本地匹配。没有可用 snapshot 时，命令可以读取当前 API 数据形成临时口径。

`search` 会把可见候选写入当前项目默认 session。输出表中的编号可以被 `view`、`api`、`source`、`compare`、`probe` 和 `record` 继续使用。

`--target` 将 target 上下文加入候选列表。target 证据只来自已读取来源或后续命令的派生结果。

下一步通常是 `view 1` 查看单候选证据，`compare 1 2` 比较多个候选，或 `probe 1 --target js` 做本地验证。

### 英文公共稿

Short:

`Search library module candidates`

Long:

```text
search turns a query into library module candidates. It reads module summaries from the current snapshot and matches module name, description, keywords, and repository declaration. If no usable snapshot exists, search can read current API data for a transient search basis.

search writes visible candidates into the current project's default session. The numbers in the output table can be used by view, api, source, compare, probe, and record.

--target adds target context to the candidate list. Target evidence appears only when it is loaded from a source or derived by a later command.

After search, inspect a candidate with view, compare multiple candidates, or run a local probe.
```

Example:

```text
seekmoon search markdown
seekmoon search markdown --target js
seekmoon search markdown --json
seekmoon search markdown --jq '.results[].module'
seekmoon view 1
seekmoon compare 1 2
```

Flag:

`--target`: `add target context to the candidate list`

## 八、`view` help

### 中文审查稿

`view` 展示单个 library module 的证据画像。输入可以是完整 module coordinate，也可以是 `search` 产生的编号。

`view` 读取 Manifest API，并按 manifest version 读取 module index asset。输出包含 description、license、repository declaration、downloads、build status、docs URL、package index 状态和 package 摘要。

`view` 不展开完整 API symbol。API 详情由 `api` 命令读取。需要源码材料时进入 `source`。

### 英文公共稿

Short:

`View a library module profile`

Long:

```text
view shows the evidence profile for one library module. The input can be a full module coordinate or a numbered candidate produced by search.

view reads the Manifest API and the module index asset for the manifest version. The profile includes description, license, repository declaration, downloads, build status, docs URL, package index state, and package summary.

view does not expand full API symbols. Use api for package API details and source for published source material.
```

Example:

```text
seekmoon view mizchi/markdown
seekmoon view 1
seekmoon view 1 --json
seekmoon api 1 --package mizchi/markdown/src/api
seekmoon source 1
```

## 九、`api` help

### 中文审查稿

`api` 展示某个 package 的 API profile。输入是 module coordinate 或候选编号，并且必须提供 `--package <path>`。

`api` 先读取 module index，确认 package path 并派生 package relpath，然后读取同版本的 `package_data.json`。输出包含类型、函数、trait、docstring、signature 和 source location。

package path 不存在时，错误表面使用 module index 中的已知 package paths 帮助恢复。

### 英文公共稿

Short:

`View a package API profile`

Long:

```text
api shows the API profile for one package. The input is a module coordinate or numbered candidate, plus --package <path>.

api reads the module index, resolves the package path to a package relpath, and then reads package_data.json for the same module version. The output contains types, values, traits, docstrings, signatures, and source locations.

When the package path is absent, the error surface uses known package paths from the module index to support recovery.
```

Example:

```text
seekmoon api mizchi/markdown --package mizchi/markdown/src/api
seekmoon api 1 --package mizchi/markdown/src/api
seekmoon api 1 --package mizchi/markdown/src/api --jq '.values[].name'
seekmoon api --shape
```

Flag:

`--package`: `package path inside the selected module`

## 十、`source` help

### 中文审查稿

`source` 定位 registry 发布版本对应的源码材料。输入可以是 module coordinate、带版本的 coordinate，或 `search` 产生的编号。

`source` 记录每次来源解析尝试，包括 `moon fetch`、source zip、本地 cache、core 本地源码和 repository signal。选中的源码来源来自成功的解析尝试。

源码定位会产生本地文件系统结果。默认结果属于 SeekMoon 控制的 source 或 cache 边界。

### 英文公共稿

Short:

`Locate published source`

Long:

```text
source locates source material for the registry-published module version. The input can be a module coordinate, a versioned coordinate, or a numbered candidate from search.

source records each resolution attempt, including moon fetch, source zip, local cache, core local source, and repository signal. The selected source is derived from successful attempts.

Source resolution can produce local filesystem material under SeekMoon-controlled source or cache paths.
```

Example:

```text
seekmoon source mizchi/markdown@0.6.2
seekmoon source 1
seekmoon source 1 --json
seekmoon source 1 --jq '.selected_source.path'
```

## 十一、`skill` help

### 中文审查稿

`skill` 处理 Mooncakes Skills API 中的 executable skill entry。Skill entry 是可执行 Wasm 或 runwasm 对象。

`skill search` 生成 skill 候选，并写入当前项目默认 session。`skill view` 展示 skill profile、asset 状态和 pinned runwasm coordinate。

Skill 使用独立的 search/view 路径。记录 skill 判断时使用 `record --kind skill`。

### 英文公共稿

Short:

`Search or view executable skill entries`

Long:

```text
skill works with executable skill entries from the Mooncakes Skills API. A skill entry is an executable Wasm or runwasm object.

skill search creates numbered skill candidates in the current project's default session. skill view shows the skill profile, asset states, and pinned runwasm coordinate.

Use record --kind skill when recording a skill investigation result.
```

Example:

```text
seekmoon skill search cowsay
seekmoon skill view 1
seekmoon record 1 --kind skill --conclusion continue-verification
```

## 十二、`skill search` help

### 中文审查稿

`skill search` 从 query 生成 executable skill 候选。它读取 Skills API，并按 skill name、module、package 和 metadata description 匹配。

输出编号写入当前项目默认 session。后续可以用 `skill view 1` 读取 skill profile，也可以用 `record 1 --kind skill` 保存判断。

### 英文公共稿

Short:

`Search skill entries`

Long:

```text
skill search turns a query into executable skill candidates. It reads the Skills API and matches skill name, module, package, and metadata description.

The output numbers are written into the current project's default session. Use skill view to inspect a candidate and record --kind skill to save an investigation result.
```

Example:

```text
seekmoon skill search cowsay
seekmoon skill search cowsay --json
seekmoon skill view 1
```

## 十三、`skill view` help

### 中文审查稿

`skill view` 展示一个 executable skill 的证据画像。输入可以是 skill entry coordinate 或 `skill search` 产生的编号。

`skill view` 读取 skill detail、SKILL.md、Wasm asset、checksum asset，并派生 pinned runwasm coordinate。

### 英文公共稿

Short:

`View a skill profile`

Long:

```text
skill view shows the evidence profile for one executable skill. The input can be a skill entry coordinate or a numbered candidate produced by skill search.

skill view reads skill detail, SKILL.md, Wasm asset state, checksum asset state, and derives the pinned runwasm coordinate.
```

Example:

```text
seekmoon skill view Yoorkin/cowsay
seekmoon skill view 1
seekmoon skill view 1 --jq '.runwasm_coordinate'
```

## 十四、`compare` help

### 中文审查稿

`compare` 把多个候选放在同一个证据表面中。输入可以是多个编号，也可以是多个 module coordinate。

`compare` 对齐 manifest、package index、source、probe 和已加载 repository signal 等证据字段。它展示证据差异，供消费者继续下钻、验证或记录判断。

下一步通常是对差异明显的候选运行 `view`、`api`、`source` 或 `probe`。

### 英文公共稿

Short:

`Compare candidate evidence`

Long:

```text
compare places multiple candidates on one evidence surface. Inputs can be numbered candidates or module coordinates.

compare aligns evidence fields such as manifest data, package index state, source state, probe state, and loaded repository signals. The output shows evidence differences for further inspection, validation, or recording.

After compare, inspect a specific difference with view, api, source, or probe.
```

Example:

```text
seekmoon compare 1 2
seekmoon compare 1 2 3 --json
seekmoon compare mizchi/markdown moonbit-community/cmark
seekmoon probe 1 --target js
```

## 十五、`probe` help

### 中文审查稿

`probe` 记录一个候选在当前工具链、版本、target 与命令序列下的本地验证证据。输入可以是候选编号或 module coordinate。

默认 probe 在隔离目录中创建验证项目，执行 `moon add`、`moon check`、`moon test` 和 target check/build。每个步骤记录 command、cwd、exit code、状态和 log path。

probe 结果属于 local derived evidence。上游来源字段保持各自的来源证据身份。

### 英文公共稿

Short:

`Run local validation for a candidate`

Long:

```text
probe records local validation evidence for one candidate in the current toolchain, version, target, and command sequence. The input can be a numbered candidate or module coordinate.

The default probe creates an isolated validation project and runs moon add, moon check, moon test, and target check/build steps. Each step records command, cwd, exit code, status, and log path.

A probe result is local derived evidence. Upstream source fields keep their own source-evidence identity.
```

Example:

```text
seekmoon probe 1
seekmoon probe 1 --target js
seekmoon probe mizchi/markdown@0.6.2 --target wasm-gc
seekmoon probe 1 --json
```

Flag:

`--target`: `target backend for local validation steps`

## 十六、`record` help

### 中文审查稿

`record` 保存一次采纳判断。输入是候选编号或 coordinate，并且必须提供 `--conclusion <value>`。

record 写入候选、版本、项目上下文、snapshot、证据引用、结论和 note。结论使用稳定枚举，JSON 输出保持英文枚举值。

library 候选默认使用 `--kind library`。skill 判断使用 `--kind skill`。

### 英文公共稿

Short:

`Save an adoption judgment`

Long:

```text
record saves one adoption judgment. The input is a numbered candidate or coordinate, plus --conclusion <value>.

record writes candidate, version, project context, snapshot, evidence references, conclusion, and note. Conclusions use stable enum values, and JSON output keeps those enum values in English.

Library candidates use --kind library by default. Use --kind skill when recording a skill investigation.
```

Example:

```text
seekmoon record 1 --conclusion continue-verification
seekmoon record 1 --conclusion reject-for-now --note "API coverage does not match this project."
seekmoon record 1 --kind skill --conclusion continue-verification
seekmoon record 1 --json
```

Flags:

`--conclusion`: `adoption conclusion enum value`

`--note`: `human note stored with the record`

`--kind`: `candidate kind recorded by this judgment`

## 十七、`report` help

### 中文审查稿

`report` 从已有 records、snapshot、项目上下文和证据引用生成调查报告。它输出调查轨迹的文档投影。

Report 只列已经被记录或引用的来源。没有执行的验证动作不会出现在报告中。

`--format` 指定报告格式。当前公共路径以 markdown 和 json 为主要目标。

### 英文公共稿

Short:

`Render an investigation report`

Long:

```text
report renders an investigation report from existing records, snapshots, project context, and evidence references. It is the document projection of the recorded investigation path.

A report lists sources that were recorded or referenced. Validation steps that were not run do not appear as validation results.

--format selects the report format.
```

Example:

```text
seekmoon report --format markdown
seekmoon report --format json
seekmoon report --format markdown --json
```

Flag:

`--format`: `report format to render`

## 十八、`raw` help

### 中文审查稿

`raw` 读取指定来源的原始 payload。它保留上游字段名和原始 shape，并附带 source status 和 metadata。

`raw` 服务来源审计、字段复查和失败复现。普通 discovery 路径优先使用 `search`、`view`、`api`、`source` 和 `skill`。

### 英文公共稿

Short:

`Read a raw source payload`

Long:

```text
raw reads the requested source payload without normalizing its upstream field names or shape. The output includes source status and metadata.

raw serves source audit, field inspection, and failure reproduction. Use search, view, api, source, and skill for the ordinary discovery path.
```

Example:

```text
seekmoon raw modules
seekmoon raw manifest mizchi/markdown
seekmoon raw module-index mizchi/markdown@0.6.2
seekmoon raw package-data mizchi/markdown@0.6.2 mizchi/markdown/src/api
seekmoon raw skills
```

## 十九、为什么这样安排

### 1. 根 help 建立行为场景

根 help 必须说明 SeekMoon 是 workbench，并给出工作流。没有这层，使用者只看到命令列表，无法形成“从候选到证据到验证到记录”的行动路径。根 help 不展开每个命令细节，因为那会把入口文本变成手册。它只负责建立进入路径。

### 2. 子命令 help 只服务当前动作

每个子命令只写当前动作所需信息。`api` 不讲 `probe` 的完整意义，`probe` 不讲 Mooncakes API 全部来源，`record` 不讲 discovery 理论。这样遵守量的准则和关系准则，避免把使用者带离当前动作。

### 3. 编号引用写在相关命令里

编号引用改变用户实际输入方式，是公共契约的一部分。`search`、`skill search` 写入编号；`view`、`api`、`source`、`compare`、`probe`、`record` 消费编号。相关 help 都要写。没有编号能力的 `doctor`、`sync`、`report`、`raw` 不写编号。

### 4. 边界句放在对象归属位置

`probe` 的边界写成“probe result is local derived evidence”；`source` 的边界写成“source material for registry-published module version”；`report` 的边界写成“lists sources that were recorded or referenced”。这些句子不是焦虑式提醒，而是对象归属规则。

### 5. 示例给出行动邀请

示例不是越多越好。每个命令只保留能改变下一步行动的示例：普通用法、编号用法、JSON/jq/shape 用法、下一步连接。示例顺序从普通终端路径到自动化路径，不把高级路径放在第一位。

### 6. 中文审查稿留在代码里

中文审查稿让每个英文 help 背后的对象结构可审查。它不是运行时语言功能。这样既满足你逐字审查的需要，又不引入未授权的 CLI 公共承诺。

## 二十、我会避免的写法

我不会在 help 正文写设计过程，例如“这个命令是为了避免用户误解而设计”。

我不会在 help 正文写理论名词，例如“根据媒介环境学，help 是行为场景”。理论只约束写作者，不进入使用者 help。

我不会把否定性警告放在核心位置，例如“不把 probe 当质量证明”。需要边界时写成“probe result is local derived evidence”。

我不会把未来可能功能写进 help，例如 session id、语言切换、repository enrichment 默认加载策略。

我不会翻译 JSON 字段、schema id、枚举值。机器契约保持英文稳定。

## 二十一、确认后执行步骤

你确认文案后，我会按以下步骤写代码：

1. 新增 `internal/helpdoc` 包，写入中英文案和注入逻辑。
2. 在 `internal/cli/root.go` 调用 help 注入。
3. 调整通用 output flags 的英文 usage。
4. 为命令专属 flags 调整 usage。
5. 增加 help 文案测试和中文审查稿存在性测试。
6. 运行 `go test ./...`。
7. 实际运行 `go run ./cmd/seekmoon --help` 和几个关键子命令 help，检查输出表面。

这套设计的核心是：help 成为 SeekMoon 工作台的入口表面；中文成为代码内审查材料；英文成为公共 CLI help；不新增语言切换命令；不把我的反思、理论和焦虑写给使用者。