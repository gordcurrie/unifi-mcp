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
| `UNIFI_ALLOW_DESTRUCTIVE` | No | Set `true` to enable permanent-delete tools (e.g. `delete_dns_policy`) |

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

### Pagination

All `list_*` tools support optional pagination via `offset` and `limit` fields:

| Field | Type | Description |
|---|---|---|
| `offset` | `int` | Number of items to skip (0 = use API default, i.e. return from the start) |
| `limit` | `int` | Maximum number of items to return (0 = use API default page size) |

All list tools return a `Page[T]` object containing:

```json
{
  "data": [...],
  "totalCount": 42,
  "offset": 0,
  "limit": 25,
  "count": 25
}
```

Omit both fields (or pass `0`) to fetch the API's default first page.

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
| `list_pending_devices` | Devices visible on the network not yet adopted |
| `power_cycle_port` | Power-cycle a single PoE port on a switch — requires `confirmed: true` |

### Clients

| Tool | Description |
|---|---|
| `list_clients` | All currently connected clients |
| `get_client` | Details for a single connected client by ID |
| `authorize_guest_client` | Authorize a connected client for guest network access — requires `confirmed: true` |

### Network

| Tool | Description |
|---|---|
| `list_wifi_broadcasts` | All WiFi broadcast (SSID) configurations |
| `get_wifi_broadcast` | Details for a single WiFi broadcast (SSID) by ID |
| `set_wifi_broadcast_enabled` | Enable or disable a WiFi broadcast (SSID) — requires `confirmed: true` |
| `list_networks` | LAN/VLAN network configurations |
| `list_firewall_policies` | All firewall policies |
| `get_firewall_policy` | Details for a single firewall policy by ID |
| `set_firewall_policy_enabled` | Enable or disable a firewall policy — requires `confirmed: true` |
| `list_firewall_zones` | All firewall zones |
| `get_firewall_zone` | Details for a single firewall zone by ID |
| `create_firewall_zone` | Create a new firewall zone |
| `update_firewall_zone` | Update an existing firewall zone by ID |
| `list_acl_rules` | All ACL rules |
| `get_acl_rule` | Details for a single ACL rule by ID |
| `get_acl_rule_ordering` | Current ACL rule evaluation order |
| `list_traffic_matching_lists` | All traffic matching lists (IP/port sets used by firewall policies) |
| `get_traffic_matching_list` | Details for a single traffic matching list by ID |
| `list_wans` | All WAN interface definitions |
| `list_vpn_tunnels` | All site-to-site VPN tunnels |
| `list_vpn_servers` | All VPN server configurations |
| `list_dns_policies` | All local DNS A-record policies |
| `get_dns_policy` | Details for a single DNS policy by ID |
| `create_dns_policy` | Create a new local DNS A-record policy |
| `update_dns_policy` | Update an existing DNS policy by ID |
| `list_vouchers` | All hotspot vouchers |
| `get_voucher` | Details for a single hotspot voucher by ID |
| `create_vouchers` | Generate one or more hotspot vouchers — requires `confirmed: true` |
| `list_device_tags` | All device tags defined for the site |
| `list_dpi_categories` | All DPI application categories (used in firewall matching rules) |
| `list_dpi_applications` | All DPI applications (used in firewall matching rules) |
| `list_radius_profiles` | All RADIUS profiles configured for the site |

### Destructive (requires `UNIFI_ALLOW_DESTRUCTIVE=true`)

| Tool | Description |
|---|---|
| `delete_dns_policy` | Permanently delete a DNS policy — requires `confirmed: true` |
| `delete_firewall_policy` | Permanently delete a firewall policy — requires `confirmed: true` |
| `delete_firewall_zone` | Permanently delete a firewall zone — requires `confirmed: true` |
| `create_acl_rule` | Create a new ACL rule — requires `confirmed: true` |
| `update_acl_rule` | Update an existing ACL rule by ID — requires `confirmed: true` |
| `set_acl_rule_enabled` | Enable or disable an ACL rule — requires `confirmed: true` |
| `reorder_acl_rules` | Set the ACL rule evaluation order — requires `confirmed: true` |
| `delete_acl_rule` | Permanently delete an ACL rule — requires `confirmed: true` |
| `delete_voucher` | Permanently revoke a hotspot voucher — requires `confirmed: true` |

> **Why are all ACL writes destructive-gated?** Unlike firewall zones (organisational
> containers), any ACL mutation directly controls which traffic is allowed or blocked.
> A misplaced `BLOCK` rule — or a reorder that promotes one — can cause a complete
> network outage. Gating every write on `UNIFI_ALLOW_DESTRUCTIVE=true` ensures that
> an AI session without the flag cannot issue any ACL write at all. The `confirmed: true`
> field is a secondary per-call guard; the env flag is the primary one.

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
