BINARY := bin/unifi-mcp
CMD     := ./cmd/unifi-mcp

# Pinned tool versions — update these intentionally when upgrading
GOFUMPT_VERSION       := v0.7.0
GOSEC_VERSION         := v2.22.8
GOVULNCHECK_VERSION   := v1.1.4
GOLANGCILINT_VERSION  := v2.10.1

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
	go build -o $(BINARY) $(CMD)

## check – run all quality gates in order (pre-commit gate)
check: fix fmt vet lint sec vulncheck test build

## clean – remove build artifacts
clean:
	rm -f $(BINARY)
