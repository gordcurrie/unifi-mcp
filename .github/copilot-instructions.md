# GitHub Copilot Instructions — unifi-mcp

## Project Overview

This is an MCP server written in Go that exposes UniFi home network operations as tools,
built on the official `modelcontextprotocol/go-sdk`.

## Architecture

```
cmd/unifi-mcp/main.go          # entrypoint — flags, env vars, server bootstrap
internal/unifi/                # custom UniFi HTTP client package (no third-party UniFi libs)
tools/                         # MCP tool registration
```

## UniFi API

Two API path styles are used under the same `X-API-Key` auth header on UCG-Max:

- **v1 API** (`/v1/sites/{siteID}/...`) — sites, devices, clients, statistics
- **Legacy API** (`/api/s/{site}/...`) — management commands (devmgr, stamgr) and
  network config (wlanconf, firewallrule, portforward) not yet in v1

The base URL is always `https://<console-ip>/proxy/network`. Both path styles are
appended relative to this base.

The default site is configured via the `UNIFI_SITE_ID` env var on the `unifi.Client`
struct so tools don't need to require it as a parameter for single-site home lab use.
Where a tool accepts an optional `site_id`, it falls back to the client default.

## Code Conventions

### Tool Registration Pattern

Each tool group lives in `tools/<group>.go` and exports a single unexported
`register<Group>Tools(s *mcp.Server, client *unifi.Client)` function. All functions
are wired in `tools/register.go`'s `RegisterAll`.

```go
func registerSiteTools(s *mcp.Server, client *unifi.Client) {
    type getSiteInput struct {
        SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
    }
    mcp.AddTool(s, &mcp.Tool{
        Name:        "get_site",
        Description: "...",
        Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
    }, func(ctx context.Context, _ *mcp.CallToolRequest, input getSiteInput) (*mcp.CallToolResult, any, error) {
        site, err := client.GetSite(ctx, input.SiteID)
        if err != nil {
            return nil, nil, fmt.Errorf("get_site: %w", err)
        }
        return jsonResult(site)
    })
}
```

### Input Structs

- Define input structs **inline** inside the `register*Tools` function, scoped as tightly
  as possible (reuse across tools in the same function when they share the same shape).
- All fields use `json` + `jsonschema` struct tags.
- Optional fields use `omitempty`.
- Use `*bool` for optional booleans to distinguish "not provided" from `false`.

### Result Helpers (`tools/helpers.go`)

```go
// jsonResult marshals v to a JSON TextContent result.
func jsonResult(v any) (*mcp.CallToolResult, any, error)

// textResult wraps a plain string in a TextContent result.
func textResult(s string) (*mcp.CallToolResult, any, error)
```

Always return errors wrapped with `fmt.Errorf("tool_name: %w", err)`.

### Destructive Tools

- Only registered when `UNIFI_ALLOW_DESTRUCTIVE=true`.
- Always require `confirmed bool` field — return error if `false`.
- Annotated with `DestructiveHint: &destructiveTrue` where `destructiveTrue = true`.

### HTTP Client (`internal/unifi/client.go`)

```go
type Client struct {
    baseURL    string       // e.g. "https://192.168.1.1/proxy/network"
    apiKey     string
    siteID     string       // default site for single-site use
    httpClient *http.Client // TLS skip configurable
}
```

- Private methods: `get(ctx, path)`, `post(ctx, path)`,
  `postWithBody(ctx, path, body)`, `put(ctx, path, body)` — all return `([]byte, error)`.
- Unmarshal via a typed `apiResponse[T]` generic wrapper.
- v1 responses: `{"data": T, ...}` — legacy responses: `{"data": [T,...], "meta": {...}}`.

### Error Wrapping in Client

```go
return nil, fmt.Errorf("ListDevices %s: %w", siteID, err)
```

### Linting

- `golangci-lint`, `gosec`, `govulncheck` must all pass before a PR.
- Run `make check` to execute all gates.
- No use of `interface{}` — use `any`.
- Exported types/functions must have doc comments.

## Env Vars

| Variable                  | Default | Notes                            |
|---------------------------|---------|----------------------------------|
| `UNIFI_BASE_URL`          | —       | required                         |
| `UNIFI_API_KEY`           | —       | required                         |
| `UNIFI_SITE_ID`           | —       | required; default site           |
| `UNIFI_INSECURE`          | `false` | skip TLS verification            |
| `UNIFI_ALLOW_DESTRUCTIVE` | `false` | enables forget/reprovision tools |

## Code Style — Idiomatic Go

- Wrap errors: `fmt.Errorf("doing X: %w", err)` — never discard errors
- Sentinel errors: `var ErrNotFound = errors.New("resource not found")`
- `context.Context` is always the first parameter on any function that does I/O
- Pointer receivers on `Client` and mutable types; value receivers on small read-only structs
- No `init()` functions anywhere — explicit initialization only
- No global mutable state — inject dependencies
- All exported types and functions must have doc comments
- `json` tags on all API types; `jsonschema` tags on MCP tool input structs
- Table-driven tests using `t.Run` subtests
- Use `gofumpt` formatting (stricter than `gofmt`)

## Documentation — Required for Every PR

Every PR that adds or changes tools must update both files before committing:

- **README.md** — add each new tool to the appropriate table (Sites, Devices, Clients,
  Statistics, Network, or Destructive). Include the tool name and a one-line description.
  Destructive tools go in the Destructive table.

- **PLAN.md** — mark the relevant section as completed (add ✅ if not already there) and
  update the running tool count at the bottom of the phase section.

Do not skip these updates under any circumstances — documentation is part of the definition
of done for every PR, the same as passing `make check`.

## Git Commits

Always write multi-line commit messages via a temp file to avoid shell quoting issues:

```bash
python3 -c "open('/tmp/msg.txt','w').write('''subject line\n\nbody line 1\nbody line 2\n''')"
git add . && git commit -F /tmp/msg.txt
```

Never pass multi-line messages with `-m` — the shell mangles them.

## PR Size Guidelines

- **Target < 500 lines changed per PR** (diffs shown by `git diff --stat`, excluding generated files).
- Split large features into logical sub-PRs — e.g. read-only tools first, then mutations.
- Each PR must be independently reviewable, `make check` clean, and merge-safe on its own.
- A good split heuristic: one new `tools/<group>.go` file + its matching `internal/unifi/` client methods per PR.

## Git Push Rules

- Use plain `git push` for normal commits (new commits on top of the branch).
- Use `git push --force-with-lease` only when history has been rewritten (rebase, amend, etc.).
- Never use `git push --force` under any circumstances without explicit user permission.

## Quality Gates — `make check` Must Pass Before Every Commit

```bash
make check
```

Runs: `fix → fmt → vet → lint → sec → vulncheck → test → build`.
Config is in `.golangci.yml`. Key linters: `gosec`, `govet`, `staticcheck`, `errcheck`,
`bodyclose`, `noctx`, `gofumpt`, `revive`, `gocritic`, `unparam`, `unconvert`.
