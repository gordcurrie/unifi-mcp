# unifi-mcp

MCP server written in Go that exposes UniFi home network operations as tools,
built on `modelcontextprotocol/go-sdk`.

## Architecture

```
cmd/unifi-mcp/main.go          # entrypoint — flags, env vars, server bootstrap
internal/unifi/                # UniFi HTTP client (v1 API only)
tools/                         # MCP tool registration
```

Only the **v1 API** (`/proxy/network/integration/v1/sites/{siteID}/...`) is used.

## Environment variables

| Variable                  | Notes                            |
|---------------------------|----------------------------------|
| `UNIFI_BASE_URL`          | required — e.g. `https://192.168.1.1/proxy/network` |
| `UNIFI_API_KEY`           | required                         |
| `UNIFI_SITE_ID`           | required — default site          |
| `UNIFI_INSECURE`          | `true` to skip TLS verification  |
| `UNIFI_ALLOW_DESTRUCTIVE` | `true` to enable mutating tools  |

## Build and quality gate

```bash
make check   # fix → fmt → vet → lint → sec → vulncheck → test → build
```

## Network security audit

Use the `audit-network-security` skill. Ask to audit, review, or assess the
network security posture and the skill will be loaded automatically.
