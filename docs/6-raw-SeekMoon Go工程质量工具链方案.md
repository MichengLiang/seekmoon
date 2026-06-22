# SeekMoon Go 工程质量工具链确认

我调查后给出明确结论：Go 生态里不需要照搬 Rust 的 `cargo-nextest` 形态。Go 的测试执行核心就是官方 `go test`，它同时承担编译、单元测试、包级测试、并发测试入口、覆盖率、fuzz 入口和缓存机制。Go 生态里真正值得补的一层是 `gotestsum`：它不是替代 `go test` 的新执行模型，而是对 `go test -json` 的人类可读输出、CI 汇总和 JUnit 报告层。类比 Rust 的 nextest 时，SeekMoon 的答案应写成：执行核心使用 `go test`；本地与 CI 的测试呈现使用 `gotestsum`。

下面是我为 SeekMoon 固定的 Go 工程质量工具链。版本以本次调查和 Go module 元数据为准。

## 1. 基础 Go 工具链

当前本机 Go：

```text
go version go1.26.3 linux/amd64
```

SeekMoon 的工程工具链基线采用 Go 1.26.x。原因是现有环境已经是 `go1.26.3`，并且前面选择的部分工具链版本已经声明 Go 1.25 或 1.26 级别的 go.mod 版本要求。为了避免本地、CI、release 三处版本漂移，书中应明确：SeekMoon Go implementation 使用 Go 1.26 toolchain。

官方必选命令：

```bash
go test ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -fuzz=Fuzz -run=^$ ./...
go vet ./...
go mod tidy
go mod verify
govulncheck ./...
```

这些是 Go 工程质量的底座，不是第三方偏好。

`go test ./...` 是功能正确性主入口。SeekMoon 的 model、source、service、store、output、contract 都通过它验证。

`go test -race ./...` 是并发风险入口。SeekMoon 会有 HTTP client、source cache、store、probe runner、output writer 等并发或共享资源边界，race detector 必须进入标准检查。

`go test -coverprofile=coverage.out ./...` 和 `go tool cover -func=coverage.out` 是覆盖率入口。覆盖率不作为唯一质量指标，但它能暴露 source reader、error surface、schema 和 jq projection 是否缺测试。

`go test -fuzz=Fuzz -run=^$ ./...` 是 fuzz 入口。SeekMoon 适合 fuzz 的对象是 module coordinate parser、package path to relpath、local index JSONL parser、source URL/path parser。gojq parser 本身由 gojq 负责，SeekMoon 只测表达式错误映射。

`go vet ./...` 是官方静态检查入口，覆盖一批 Go 语言层面的常见错误。

`go mod tidy` 和 `go mod verify` 负责依赖图卫生。前者保持 `go.mod`/`go.sum` 与源码一致；后者验证 module cache 中的依赖内容与 `go.sum` 匹配。

`govulncheck ./...` 是官方漏洞可达性扫描，不只是查依赖库名字，而是结合调用可达性判断风险。SeekMoon 读网络、处理压缩包、执行本地命令、访问 GitHub API，必须把它列入质量门。

官方资料入口：

- https://go.dev/doc/tutorial/add-a-test
- https://go.dev/doc/tutorial/fuzz
- https://go.dev/doc/tutorial/govulncheck
- https://go.dev/doc/articles/race_detector
- https://pkg.go.dev/cmd/go

## 2. 测试运行体验：`gotestsum@v1.13.0`

固定选择：

```text
gotest.tools/gotestsum@v1.13.0
```

模块元数据：tag `refs/tags/v1.13.0`，时间 `2025-09-11T03:13:16Z`。

`gotestsum` 的职责是测试输出层和 CI 报告层。Go 的测试核心仍然是 `go test`，但 `go test` 默认输出对大型项目和 CI 不够好读。`gotestsum` 消费 `go test -json`，提供稳定格式、失败摘要、包级状态和 JUnit XML。

SeekMoon 的标准测试命令应是：

```bash
gotestsum --format=pkgname -- -count=1 ./...
gotestsum --format=pkgname --junitfile=.artifacts/test/unit.xml -- -count=1 ./...
```

CI 里使用 JUnit 输出：

```bash
gotestsum --format=github-actions --junitfile=.artifacts/test/unit.xml -- -count=1 ./...
```

这里的判断是：Go 没有必要找一个“替代 go test 的 nextest”。`gotestsum` 是 Go 生态里成熟的测试呈现工具，正好补 Go 官方测试输出的可读性和 CI 汇总。

来源：

- https://github.com/gotestyourself/gotestsum
- https://pkg.go.dev/gotest.tools/gotestsum

## 3. Lint 聚合：`golangci-lint@v2.12.2`

固定选择：

```text
github.com/golangci/golangci-lint/v2@v2.12.2
```

模块元数据：tag `refs/tags/v2.12.2`，时间 `2026-05-06T11:01:25Z`。

Go 官方有 `go vet`，但 SeekMoon 需要更完整的工程质量门：未处理错误、低质量结构、重复代码、复杂度、无效赋值、静态安全、拼写、导入顺序、注释规则、unused、ineffassign、staticcheck 等。`golangci-lint` 是 Go 生态事实标准的 lint 聚合器，适合作为 CI 质量门。

SeekMoon 的 lint 命令：

```bash
golangci-lint run ./...
```

配置文件建议采用：

```text
.golangci.yml
```

启用方向应包括：

```yaml
version: "2"
linters:
  enable:
    - govet
    - staticcheck
    - errcheck
    - ineffassign
    - unused
    - revive
    - gosec
    - misspell
    - gocritic
    - unparam
    - nakedret
    - prealloc
```

`gosec` 不作为单独命令进入主工具链，而是由 `golangci-lint` 管理。这样安全静态检查仍存在，但不会多出一个重复质量入口。真正的漏洞可达性扫描由官方 `govulncheck` 承担。

来源：

- https://golangci-lint.run/
- https://github.com/golangci/golangci-lint

## 4. 格式化：`gofmt` + `gofumpt@v0.10.0`

固定选择：

```text
mvdan.cc/gofumpt@v0.10.0
```

模块元数据：tag `refs/tags/v0.10.0`，时间 `2026-05-04T15:05:39Z`。

Go 官方格式化是 `gofmt`，这仍是底线。SeekMoon 可以使用更严格的 `gofumpt` 作为主格式化工具，因为它是 gofmt 的严格超集，能减少无意义风格分歧。

标准命令：

```bash
gofumpt -w .
gofumpt -l .
```

CI 检查：

```bash
test -z "$(gofumpt -l .)"
```

导入整理使用官方工具：

```bash
goimports -w .
```

如果要减少工具数量，`goimports` 可以只作为编辑器或 CI 辅助；主格式化质量门用 `gofumpt`。SeekMoon 文档里应明确：格式化结果由 `gofumpt` 固定，导入分组由 `goimports` 或 golangci-lint 对应 linter 约束。

来源：

- https://pkg.go.dev/cmd/gofmt
- https://github.com/mvdan/gofumpt

## 5. 漏洞扫描：`govulncheck@v1.4.0`

固定选择：

```text
golang.org/x/vuln@v1.4.0
```

模块元数据：tag `refs/tags/v1.4.0`，时间 `2026-06-17T18:58:06Z`。

`govulncheck` 是 Go 官方漏洞扫描工具，读取 Go vulnerability database，并结合源码调用路径判断 vulnerability 是否可达。SeekMoon 处理网络响应、zip、GitHub API、local command 和 JSON，因此依赖安全必须是常规检查。

标准命令：

```bash
govulncheck ./...
```

CI 中它应作为独立 job。它不是 lint，也不是 test；它是 dependency vulnerability exposure check。

来源：

- https://go.dev/doc/tutorial/govulncheck
- https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck

## 6. 发布工具：`goreleaser@v2.16.0`

固定选择：

```text
github.com/goreleaser/goreleaser/v2@v2.16.0
```

模块元数据：tag `refs/tags/v2.16.0`，时间 `2026-05-24T14:47:07Z`。

SeekMoon 是 CLI 工具，最终需要多平台二进制、checksums、archives、GitHub Release artifact。GoReleaser 是 Go CLI 发布生态里的成熟选择。它承担 release artifact 生成，不承担测试和 lint。

标准命令：

```bash
goreleaser check
goreleaser release --snapshot --clean
goreleaser release --clean
```

本地验证 release 配置用 `check`，本地预演用 snapshot，正式 tag release 用 release。

来源：

- https://goreleaser.com/
- https://github.com/goreleaser/goreleaser

## 7. 任务编排：`just`

固定选择：

```text
just 1.50.0
```

本机已经安装：`just 1.50.0`。

Go 生态有 `mage@v1.17.2`，它是 Go-native task runner。但 SeekMoon 位于一个多语言 workbench 里，已有书架是 Node/pnpm 构建，未来 Go CLI、AsciiDoc、MoonBit probe、GitHub release 都会共存。这个场景更适合语言无关的 `just`，而不是引入 Go-only 的 Mage。Rust 项目里常用 Just；这里同样采用 Just 作为 repo-level orchestration。

建议 `justfile` 命令：

```make
fmt:
    gofumpt -w .

fmt-check:
    test -z "$(gofumpt -l .)"

lint:
    golangci-lint run ./...

test:
    gotestsum --format=pkgname -- -count=1 ./...

test-race:
    gotestsum --format=pkgname -- -race -count=1 ./...

cover:
    go test -coverprofile=.artifacts/coverage.out ./...
    go tool cover -func=.artifacts/coverage.out

fuzz:
    go test -fuzz=Fuzz -run=^$ ./...

vuln:
    govulncheck ./...

mod-check:
    go mod tidy
    git diff --exit-code -- go.mod go.sum
    go mod verify

release-check:
    goreleaser check

release-snapshot:
    goreleaser release --snapshot --clean

ci: fmt-check lint test test-race cover vuln mod-check release-check
```

这里 `just` 承担命令编排，不替代 Go 工具。每个 recipe 调用一个明确工具。

## 8. 覆盖率与测试套件分层

SeekMoon 测试套件应分层，不用一个 `go test ./...` 混掉所有语义。

标准测试层：

```bash
gotestsum --format=pkgname -- -count=1 ./...
```

竞态测试层：

```bash
gotestsum --format=pkgname -- -race -count=1 ./...
```

覆盖率层：

```bash
go test -coverprofile=.artifacts/coverage.out ./...
go tool cover -func=.artifacts/coverage.out
```

Fuzz 层：

```bash
go test -fuzz=Fuzz -run=^$ ./internal/model ./internal/source
```

集成测试层：

```bash
SEEKMOON_INTEGRATION=1 gotestsum --format=pkgname -- -count=1 ./...
```

集成测试默认关闭。原因是 SeekMoon 的真实网络、真实 Moon CLI、GitHub API、source zip 下载和 probe 项目都可能触碰外部状态。默认测试必须可重复、离线、无副作用；集成测试由环境变量显式开启。

## 9. CI 质量门

SeekMoon CI 应固定这些 job：

1. `format`：`gofumpt -l .` 必须为空。
2. `lint`：`golangci-lint run ./...`。
3. `test`：`gotestsum --junitfile=.artifacts/test/unit.xml -- -count=1 ./...`。
4. `race`：`gotestsum --junitfile=.artifacts/test/race.xml -- -race -count=1 ./...`。
5. `coverage`：生成 `coverage.out` 和 coverage summary。
6. `vuln`：`govulncheck ./...`。
7. `mod`：`go mod tidy` 后 `git diff --exit-code -- go.mod go.sum`，再 `go mod verify`。
8. `release-check`：`goreleaser check`。

这个矩阵覆盖：格式、静态质量、测试、并发、覆盖率、漏洞、依赖一致性和发布配置。

## 10. 工具安装方式

Go 工具建议用 `go install tool@version` 固定：

```bash
go install gotest.tools/gotestsum@v1.13.0
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
go install golang.org/x/vuln/cmd/govulncheck@v1.4.0
go install mvdan.cc/gofumpt@v0.10.0
go install github.com/goreleaser/goreleaser/v2@v2.16.0
```

`just` 不用 Go install，因为它已经是 repo-level task runner，当前环境已有 `just 1.50.0`。如果要在 CI 安装，可以用系统包、预构建二进制或 GitHub action；不把它作为 Go module 依赖。

## 最终工具链清单

明确采用：

```text
Go toolchain: go1.26.x，当前 go1.26.3
test runner display: gotest.tools/gotestsum@v1.13.0
lint: github.com/golangci/golangci-lint/v2@v2.12.2
format: mvdan.cc/gofumpt@v0.10.0 + gofmt baseline
vulnerability: golang.org/x/vuln/cmd/govulncheck@v1.4.0
release: github.com/goreleaser/goreleaser/v2@v2.16.0
task runner: just 1.50.0
standard checks: go test, go test -race, go test -coverprofile, go test -fuzz, go vet, go mod tidy, go mod verify
```

明确的工程质量入口：

```bash
just fmt-check
just lint
just test
just test-race
just cover
just vuln
just mod-check
just release-check
just ci
```

这套工具链的设计理由是：Go 官方工具承担执行核心；`gotestsum`改善测试输出和 CI 汇总；`golangci-lint`承担聚合静态质量门；`gofumpt`固定更严格格式；`govulncheck`承担官方漏洞可达性扫描；`goreleaser`承担 CLI 发布工件；`just`承担跨语言、跨工具的命令编排。它不是堆工具，而是每个工程质量问题只有一个主入口。
