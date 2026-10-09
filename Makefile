BINARY := bin/unifi-mcp
CMD     := ./cmd/unifi-mcp

# Pinned tool versions — update these intentionally when upgrading
GOFUMPT_VERSION       := v0.12.0
# gosec: pinned to master (7b1b5ce) for Go 1.27.2 support; v2.29.0 cannot read
# Go 1.27.2 export data (securego/gosec#1771). Switch to v2.29.1+ once tagged.
GOSEC_VERSION         := v2.29.1-0.20261009120814-7b1b5cebe007
GOVULNCHECK_VERSION   := v1.8.0
GOLANGCILINT_VERSION  := v2.14.0

.PHONY: all install-tools fix fmt vet lint sec vulncheck test build check clean

all: check

## install-tools – installs all required dev tools
install-tools:
	go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)
	go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

## fix – apply go fix to update deprecated API usage
fix:
	go fix ./...

## fmt – format all Go source files with gofumpt
fmt:
	gofumpt -w .

## vet – run go vet
vet:
	go vet ./...

## lint – run golangci-lint
lint:
	golangci-lint run ./...

## sec – run gosec standalone security scanner
sec:
	gosec ./...

## vulncheck – check for known vulnerabilities in dependencies
vulncheck:
	govulncheck ./...

## test – run tests with race detector
test:
	go test -race -count=1 ./...

## build – compile the binary
build:
	@mkdir -p bin
	go build -ldflags "-X main.version=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)" -o $(BINARY) $(CMD)

## check – run all quality gates in order (pre-commit gate)
check: fix fmt vet lint sec vulncheck test build

## clean – remove build artifacts
clean:
	rm -f $(BINARY)
