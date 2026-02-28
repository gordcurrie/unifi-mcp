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
| `UNIFI_ALLOW_DESTRUCTIVE` | No | Set `true` to enable `forget_client` and `force_reprovision_device` |

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
| `restart_device` | Restart a device |
| `locate_device` | Flash device LEDs to locate it |
| `unlocate_device` | Stop LED location flash |
| `upgrade_device` | Trigger a firmware upgrade |
| `run_speed_test` | Start an internet speed test |
| `get_speed_test_status` | Poll speed test progress and results |

### Clients

| Tool | Description |
|---|---|
| `list_active_clients` | Currently connected clients |
| `list_known_clients` | All ever-seen clients |
| `block_client` | Block a client by MAC address |
| `unblock_client` | Unblock a client by MAC address |
| `kick_client` | Disconnect a client (they can reconnect) |

### Statistics

| Tool | Description |
|---|---|
| `get_site_stats` | Aggregate traffic and client counts for the site |
| `get_device_stats` | Per-device traffic statistics |
| `get_client_stats` | Per-client traffic statistics |
| `list_events` | Recent UniFi events |
| `list_alarms` | Active (or archived) alarms |

### Network

| Tool | Description |
|---|---|
| `list_wlans` | All WLAN configurations |
| `enable_wlan` | Enable a WLAN by ID |
| `disable_wlan` | Disable a WLAN by ID |
| `list_networks` | LAN/VLAN network configurations |
| `list_firewall_rules` | All firewall rules |
| `list_port_forwards` | All port-forward rules |

### Destructive (requires `UNIFI_ALLOW_DESTRUCTIVE=true`)

All destructive tools require `"confirmed": true` in the call.

| Tool | Description |
|---|---|
| `forget_client` | Permanently remove a client from history |
| `force_reprovision_device` | Force-reprovision a device (disrupts traffic) |

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
  sites.go / devices.go / clients.go / statistics.go / network.go
tools/                     — MCP tool registration
  register.go              — RegisterAll wires all tool groups
  sites.go / devices.go / clients.go / statistics.go / network.go / destructive.go
  helpers.go               — jsonResult / textResult helpers
```

Both v1 (`/v1/sites/{id}/...`) and legacy (`/api/s/{site}/...`) UniFi API paths are used,
both authenticated with the `X-API-Key` header.
