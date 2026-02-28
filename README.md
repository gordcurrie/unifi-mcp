# unifi-mcp

An [MCP](https://modelcontextprotocol.io) server written in Go that exposes UniFi home-network operations as tools,
built on the official [`modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk).

---

## Prerequisites

- Go 1.26 or later
- UniFi OS console (UCG-Max, UDM-Pro, etc.) with the REST API enabled
- A UniFi **API key** generated under *UniFi OS → Settings → API*

---

## Build

```bash
make build          # produces bin/unifi-mcp
# or
go build -o bin/unifi-mcp ./cmd/unifi-mcp
```

---

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `UNIFI_BASE_URL` | Yes | Full proxy base URL, e.g. `https://192.168.1.1/proxy/network` |
| `UNIFI_API_KEY` | Yes | UniFi API key |
| `UNIFI_SITE_ID` | Yes | Default site UUID — find it with `list_sites` |
| `UNIFI_INSECURE` | No | Set `true` to skip TLS verification (self-signed certs) |

---

## Running

### stdio (default — for MCP clients)

```bash
export UNIFI_BASE_URL=https://192.168.1.1/proxy/network
export UNIFI_API_KEY=your-api-key
export UNIFI_SITE_ID=your-site-uuid

./bin/unifi-mcp
```

### HTTP streamable

> **Note:** The HTTP transport binds to the address as-is with no authentication.
> Use a loopback address (`:8080` binds all interfaces; prefer `127.0.0.1:8080`)
> or place it behind a reverse proxy with auth before exposing it on a shared network.

```bash
./bin/unifi-mcp --transport http --addr 127.0.0.1:8080
```

---

## MCP Client Configuration

### VS Code (`.vscode/mcp.json`)

```json
{
  "servers": {
    "unifi": {
      "type": "stdio",
      "command": "/path/to/bin/unifi-mcp",
      "env": {
        "UNIFI_BASE_URL": "https://192.168.1.1/proxy/network",
        "UNIFI_API_KEY": "your-api-key",
        "UNIFI_SITE_ID": "your-site-uuid",
        "UNIFI_INSECURE": "true"
      }
    }
  }
}
```

### Claude Desktop (`claude_desktop_config.json`)

```json
{
  "mcpServers": {
    "unifi": {
      "command": "/path/to/bin/unifi-mcp",
      "env": {
        "UNIFI_BASE_URL": "https://192.168.1.1/proxy/network",
        "UNIFI_API_KEY": "your-api-key",
        "UNIFI_SITE_ID": "your-site-uuid",
        "UNIFI_INSECURE": "true"
      }
    }
  }
}
```

---

## Tools Reference

### Sites

| Tool | Description |
|---|---|
| `get_application_info` | UniFi application version and build info |
| `list_sites` | All sites on the console |
| `get_site` | Details for a single site (defaults to `UNIFI_SITE_ID`) |

### Devices

| Tool | Description |
|---|---|
| `list_devices` | All managed devices (APs, switches, gateways) |
| `get_device` | Details for a single device by ID |
| `restart_device` | Restart a device by ID — requires `confirmed: true` |
| `get_device_stats` | Latest CPU, memory, and uptime stats for a device |

### Clients

| Tool | Description |
|---|---|
| `list_clients` | All currently connected clients |

### Network

| Tool | Description |
|---|---|
| `list_wifi_broadcasts` | All WiFi broadcast (SSID) configurations |
| `list_networks` | LAN/VLAN network configurations |
| `list_firewall_policies` | All firewall policies |
| `list_firewall_zones` | All firewall zones |
| `list_acl_rules` | All ACL rules |

---

## Development

```bash
make install-tools   # install golangci-lint, gosec, govulncheck
make check           # fix, fmt, vet, lint, sec, vulncheck, test, build
make clean           # remove bin/
```

---

## Architecture

```
cmd/unifi-mcp/main.go      — flags, env vars, server bootstrap
internal/unifi/            — UniFi HTTP client (no third-party UniFi libs)
  client.go                — core HTTP + generic response decoders
  types.go                 — all data types
  sites.go / devices.go / clients.go / network.go
tools/                     — MCP tool registration
  register.go              — RegisterAll wires all tool groups
  sites.go / devices.go / clients.go / network.go
  helpers.go               — jsonResult / textResult helpers
```

All endpoints use the integration v1 API (`/proxy/network/integration/v1/sites/{id}/...`),
authenticated with the `X-API-Key` header.
