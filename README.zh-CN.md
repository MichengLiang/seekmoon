# SeekMoon

[![CI](https://github.com/MichengLiang/seekmoon/actions/workflows/ci.yml/badge.svg)](https://github.com/MichengLiang/seekmoon/actions/workflows/ci.yml)
[![Pages](https://github.com/MichengLiang/seekmoon/actions/workflows/pages.yml/badge.svg)](https://github.com/MichengLiang/seekmoon/actions/workflows/pages.yml)

[English](README.md) | 简体中文

SeekMoon 是一个面向 MoonBit 依赖消费者的包发现工作台。它帮助使用者在引入依赖之前完成候选发现、证据下钻、本地验证、采纳判断记录和调查报告输出。

Go CLI 读取 Mooncakes API、Mooncakes assets、MoonBit 本地工具链状态、项目上下文、本地 registry/cache 状态、repository 信号和项目记录。上游事实与本地派生证据保持分离，因此 search、inspection、comparison、probe、record、report、JSON、jq、shape、schema 和 error output 可以共享同一套 evidence-state 词汇。

本仓库还包含 AsciiDoc bookshelf。Bookshelf 定义 SeekMoon 及包复用研究系列使用的证据模型、命令契约、输出契约、实现架构和验收旅程。

## 链接

| 项目 | 位置 |
| --- | --- |
| GitHub | <https://github.com/MichengLiang/seekmoon> |
| Go module | `github.com/MichengLiang/seekmoon` |
| 已发布 bookshelf | <https://michengliang.github.io/seekmoon/> |
| License | [Apache-2.0](LICENSE) |

## 仓库结构

| 路径 | 职责 |
| --- | --- |
| `cmd/seekmoon/` | CLI 进程入口。 |
| `internal/` | Go 实现包，包括 CLI、services、sources、stores、output、contracts 和 help text。 |
| `tests/` | Acceptance、black-box、integration 和 journey tests。 |
| `bookshelf/` | AsciiDoc bookshelf 源码和构建工作区。 |
| `docs/` | 研究笔记、验证报告和原始设计材料。 |
| `spike/` | MoonBit 与 CLI 探索性实验。 |
| `justfile` | 本地 Go 质量门入口。 |

## CLI 快速开始

从源码运行：

```bash
go run ./cmd/seekmoon --help
```

构建本地二进制：

```bash
go build -o seekmoon ./cmd/seekmoon
./seekmoon --help
```

首次使用某个命令时，先从该命令的 help 进入：

```bash
seekmoon search --help
seekmoon probe --help
seekmoon record --help
```

常用调查路径：

```text
doctor -> sync -> search -> view/api/source/compare -> probe -> record -> report
```

示例会话：

```bash
seekmoon doctor
seekmoon sync
seekmoon search markdown --target js
seekmoon view 1
seekmoon api 1 --package mizchi/markdown/src/api
seekmoon compare 1 2
seekmoon probe 1 --target js
seekmoon record 1 --conclusion continue-verification
seekmoon report --format markdown
```

`search` 和 `skill search` 会把编号候选写入当前项目的默认 session。后续命令可以用 `1` 或 `2` 这样的编号继续引用候选。编号不可用时，重新运行 search，或传入 `owner/module@version` 这样的完整坐标。

## 命令

| 命令 | 动作 |
| --- | --- |
| `doctor` | 检查本地 MoonBit、registry、network 和 project-context evidence。 |
| `sync` | 创建带时间戳的 evidence snapshot。 |
| `search` | 搜索 library module 候选。 |
| `view` | 查看一个 library module 的 evidence profile。 |
| `api` | 查看一个 package API profile。 |
| `source` | 定位 registry 发布版本对应的源码材料。 |
| `skill search` | 从 Skills API 搜索 executable skill entries。 |
| `skill view` | 查看一个 executable skill profile。 |
| `compare` | 在同一证据表面比较多个候选。 |
| `probe` | 对一个候选执行本地验证。 |
| `record` | 保存一次采纳判断。 |
| `report` | 从 records 和 evidence references 输出调查报告。 |
| `raw` | 读取未规范化的上游 source payload。 |

## 输出模式

每个公共输出命令都支持通用输出模式：

| 模式 | 职责 |
| --- | --- |
| 默认 pretty text | 终端阅读。它不是解析接口。 |
| `--json` | 面向脚本和自动化的命令 JSON 投影。 |
| `--jq <expr>` | 对命令 JSON 投影求值 jq 表达式。 |
| `--shape` | 不执行数据动作，显示命令 JSON 字段树。 |
| `--schema` | 不执行数据动作，显示命令 JSON Schema。 |

示例：

```bash
seekmoon search argparse --json
seekmoon search argparse --jq '.results[].module'
seekmoon search --shape
seekmoon search --schema
```

## 本地状态

SeekMoon 区分项目调查状态和可复用远程缓存。

项目状态位于当前项目的 `.seekmoon/` 目录：

```text
.seekmoon/
  snapshots/
  sessions/
  records/
  reports/
  probes/
  sources/
  logs/
```

可复用缓存位于用户 XDG cache 目录：

```text
$XDG_CACHE_HOME/seekmoon/
  mooncakes/
  assets/
  github/
```

## 开发

SeekMoon 使用 Go 实现 CLI，使用 pnpm 构建 bookshelf。

完整本地质量门需要以下工具：

| 工具 | 职责 |
| --- | --- |
| Go 1.26.x | 构建、测试、覆盖率、module 检查和 Go 工具链。 |
| `just` | 本地质量门编排。 |
| `gofumpt` | Go 格式化门。 |
| `golangci-lint` | 聚合 lint 门。 |
| `gotestsum` | 可读的 Go test 输出。 |
| `govulncheck` | Go 可达漏洞暴露检查。 |
| `goreleaser` | release 配置和 snapshot artifact 检查。 |

安装 Go module 依赖：

```bash
go mod download
```

运行主要检查：

```bash
just fmt-check
just lint
just test
just test-race
just cover
just vuln
just mod-check
just release-check
```

运行完整本地质量门：

```bash
PATH="$(go env GOPATH)/bin:$PATH" just ci
```

`just cover` 写入 `.artifacts/coverage.out`。该文件是生成物，并被 Git 忽略。

使用真实网络、GitHub、Moon CLI 命令、source downloads 或 probe mutation 的 integration tests 是显式启用项。默认测试使用 fixtures、fake source readers、fake filesystems 和 fake command runners。

## 持续集成

`CI` workflow 包含两个 job：

| Job | 检查 |
| --- | --- |
| `go` | `gofumpt`、`golangci-lint`、unit tests、race tests、coverage、`govulncheck`、module integrity 和 `goreleaser check`。 |
| `bookshelf` | pnpm install、结构检查和 bookshelf build。 |

`Pages` workflow 从 `main` 分支或手动触发构建 `bookshelf/build/html`，并部署到 GitHub Pages。

## Bookshelf

源码目录入口是 [bookshelf/catalog.adoc](bookshelf/catalog.adoc)。

本地构建 bookshelf：

```bash
cd bookshelf
pnpm install
pnpm run check
pnpm run build
```

生成的 HTML 位于 `bookshelf/build/html/`。

## Release

Release 配置位于 [.goreleaser.yaml](.goreleaser.yaml)。

创建本地 snapshot release：

```bash
just release-snapshot
```

Release build 目标覆盖 Linux、macOS 和 Windows 的 `amd64` 与 `arm64`。

## License

Apache-2.0。见 [LICENSE](LICENSE)。
