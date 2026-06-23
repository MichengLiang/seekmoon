set dotenv-load := false

fmt:
    command -v gofumpt >/dev/null
    gofumpt -w cmd internal tests

fmt-check:
    command -v gofumpt >/dev/null
    test -z "$(gofumpt -l cmd internal tests)"

lint:
    command -v golangci-lint >/dev/null
    golangci-lint run ./...

test:
    command -v gotestsum >/dev/null
    gotestsum --format=pkgname -- -count=1 ./...

test-race:
    command -v gotestsum >/dev/null
    gotestsum --format=pkgname -- -race -count=1 ./...

cover:
    mkdir -p .artifacts
    go test -coverprofile=.artifacts/coverage.out ./...
    go tool cover -func=.artifacts/coverage.out

fuzz:
    go test -fuzz=Fuzz -run=^$ ./internal/model ./internal/source

vuln:
    command -v govulncheck >/dev/null
    govulncheck ./...

mod-check:
    go mod tidy
    git diff --exit-code -- go.mod go.sum
    go mod verify

release-check:
    command -v goreleaser >/dev/null
    goreleaser check

release-snapshot:
    command -v goreleaser >/dev/null
    goreleaser release --snapshot --clean

ci: fmt-check lint test test-race cover vuln mod-check release-check
