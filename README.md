# unifi-mcp

An [MCP](https://modelcontextprotocol.io) server that exposes [UniFi Network](https://ui.com/introduction/unifi) operations as tools, built in Go using the official [go-sdk](https://github.com/modelcontextprotocol/go-sdk).

## Tools

> All `list_*` tools accept optional `offset` and `limit` parameters for pagination and return a `Page[T]` object with `data`, `totalCount`, `offset`, `limit`, and `count` fields. Most tools also accept an optional `site_id`; omit it to use the default configured via `UNIFI_SITE_ID`.

### Sites

| Tool | Description | Parameters |
|---|---|---|
| `get_application_info` | UniFi controller version and type | — |
| `list_sites` | List sites on the controller | `offset`, `limit` (optional) |
| `get_site` | Details for a specific site | `site_id` (optional) |

### Devices

| Tool | Description | Parameters |
|---|---|---|
| `list_devices` | Adopted devices (APs, switches, gateways) | `offset`, `limit` (optional) |
| `get_device` | Details for a specific device | `device_id` |
| `get_device_stats` | Latest CPU, memory, and uptime stats | `device_id` |
| `list_pending_devices` | Devices visible on the network but not yet adopted | `offset`, `limit` (optional) |
| `restart_device` | Restart a device | `device_id`, `confirmed` (must be `true`) |
| `power_cycle_port` | Power-cycle a PoE port on a switch | `device_id`, `port_idx`, `confirmed` (must be `true`) |

### Clients

| Tool | Description | Parameters |
|---|---|---|
| `list_clients` | Currently connected clients | `offset`, `limit` (optional) |
| `get_client` | Details for a specific client | `client_id` |
| `authorize_guest_client` | Authorize a client for guest network access | `client_id`, `confirmed` (must be `true`), `time_limit_minutes` (optional), `data_limit_mb` (optional), `download_bandwidth_kbps` (optional), `upload_bandwidth_kbps` (optional) |

### Network

| Tool | Description | Parameters |
|---|---|---|
| `list_wifi_broadcasts` | WiFi broadcast (SSID) configurations | `offset`, `limit` (optional) |
| `get_wifi_broadcast` | Details for a specific WiFi broadcast | `broadcast_id` |
| `set_wifi_broadcast_enabled` | Enable or disable a WiFi broadcast | `broadcast_id`, `enabled`, `confirmed` (must be `true`) |
| `list_networks` | LAN/VLAN network configurations | `offset`, `limit` (optional) |
| `list_firewall_policies` | Firewall policies | `offset`, `limit` (optional) |
| `get_firewall_policy` | Details for a specific firewall policy | `policy_id` |
| `set_firewall_policy_enabled` | Enable or disable a firewall policy | `policy_id`, `enabled`, `confirmed` (must be `true`) |
| `list_firewall_zones` | Firewall zones | `offset`, `limit` (optional) |
| `get_firewall_zone` | Details for a specific firewall zone | `zone_id` |
| `create_firewall_zone` | Create a new firewall zone | `name`, `network_ids` (optional, comma-separated) |
| `update_firewall_zone` | Update an existing firewall zone | `zone_id`, `name`, `network_ids` (optional, comma-separated) |
| `list_acl_rules` | ACL rules | `offset`, `limit` (optional) |
| `get_acl_rule` | Details for a specific ACL rule | `rule_id` |
| `get_acl_rule_ordering` | Current ACL rule evaluation order | — |
| `list_traffic_matching_lists` | Traffic matching lists (IP/port sets used by firewall policies) | `offset`, `limit` (optional) |
| `get_traffic_matching_list` | Details for a specific traffic matching list | `list_id` |
| `list_wans` | WAN interface definitions | `offset`, `limit` (optional) |
| `list_vpn_tunnels` | Site-to-site VPN tunnels | `offset`, `limit` (optional) |
| `list_vpn_servers` | VPN server configurations | `offset`, `limit` (optional) |
| `list_dns_policies` | Local DNS A-record policies | `offset`, `limit` (optional) |
| `get_dns_policy` | Details for a specific DNS policy | `policy_id` |
| `create_dns_policy` | Create a new local DNS A-record policy | `type`, `domain`, `ipv4_address` (optional), `ttl_seconds`, `enabled` |
| `update_dns_policy` | Update an existing DNS policy | `policy_id`, `type`, `domain`, `ipv4_address` (optional), `ttl_seconds`, `enabled` |
| `list_vouchers` | Hotspot vouchers | `offset`, `limit` (optional) |
| `get_voucher` | Details for a specific hotspot voucher | `voucher_id` |
| `create_vouchers` | Generate one or more hotspot vouchers | `count`, `name` (optional), `time_limit_minutes` (optional), `data_limit_mb` (optional), `confirmed` (must be `true`) |
| `list_device_tags` | Device tags for the site | `offset`, `limit` (optional) |
| `list_dpi_categories` | DPI application categories used in firewall matching | `offset`, `limit` (optional) |
| `list_dpi_applications` | DPI applications used in firewall matching | `offset`, `limit` (optional) |
| `list_radius_profiles` | RADIUS profiles for the site | `offset`, `limit` (optional) |

### Destructive (opt-in)

These tools are **not registered by default**. Set `UNIFI_ALLOW_DESTRUCTIVE=true` to enable them.

| Tool | Description | Parameters |
|---|---|---|
| `delete_dns_policy` | Permanently delete a DNS policy | `policy_id`, `confirmed` (must be `true`) |
| `delete_firewall_policy` | Permanently delete a firewall policy | `policy_id`, `confirmed` (must be `true`) |
| `delete_firewall_zone` | Permanently delete a firewall zone | `zone_id`, `confirmed` (must be `true`) |
| `create_acl_rule` | Create a new ACL rule | `type` (`IPV4`\|`MAC`), `name`, `action` (`ALLOW`\|`BLOCK`), `enabled`, `confirmed` (must be `true`) |
| `update_acl_rule` | Update an existing ACL rule | `rule_id`, `type` (`IPV4`\|`MAC`), `name`, `action` (`ALLOW`\|`BLOCK`), `enabled`, `confirmed` (must be `true`) |
| `set_acl_rule_enabled` | Enable or disable an ACL rule | `rule_id`, `enabled`, `confirmed` (must be `true`) |
| `reorder_acl_rules` | Set the ACL rule evaluation order | `rule_ids` (comma-separated, in desired order), `confirmed` (must be `true`) |
| `delete_acl_rule` | Permanently delete an ACL rule | `rule_id`, `confirmed` (must be `true`) |
| `delete_voucher` | Permanently revoke a hotspot voucher | `voucher_id`, `confirmed` (must be `true`) |

> **Why are all ACL writes destructive-gated?** Any ACL mutation directly controls which traffic is allowed or blocked. A misplaced `BLOCK` rule — or a reorder that promotes one — can cause a complete network outage. `UNIFI_ALLOW_DESTRUCTIVE=true` is the primary guard; `confirmed: true` is the per-call secondary guard.

## Installation

### Download a pre-built binary

Download the latest release for your platform from the [Releases](https://github.com/gordcurrie/unifi-mcp/releases) page.

| Platform | Binary |
|---|---|
| Linux (amd64) | `unifi-mcp_linux_amd64` |
| Linux (arm64) | `unifi-mcp_linux_arm64` |
| macOS (amd64) | `unifi-mcp_darwin_amd64` |
| macOS (arm64) | `unifi-mcp_darwin_arm64` |
| Windows (amd64) | `unifi-mcp_windows_amd64.exe` |

Make it executable and place it on your `PATH` (substitute the filename for your platform):

```bash
chmod +x <binary-name>
mv <binary-name> /usr/local/bin/unifi-mcp
```

> Windows users: rename `unifi-mcp_windows_amd64.exe` to `unifi-mcp.exe` and place it in a directory on your `%PATH%`.

### Build from source

Requires Go 1.26+. You will also need a UniFi OS console (UCG-Max, UDM-Pro, etc.) with an API key generated under *UniFi OS → Settings → API*.

```bash
git clone https://github.com/gordcurrie/unifi-mcp
cd unifi-mcp
cp .env.example .env   # copy the example env file
$EDITOR .env           # set UNIFI_* values (see table below)
make build             # binary lands in bin/unifi-mcp
```

## Configuration

All configuration is via environment variables:

| Variable | Required | Description |
|---|---|---|
| `UNIFI_BASE_URL` | yes | e.g. `https://192.168.1.1/proxy/network` |
| `UNIFI_API_KEY` | yes | API key from *UniFi OS → Settings → API* |
| `UNIFI_SITE_ID` | yes | Default site UUID — find it with `list_sites` |
| `UNIFI_INSECURE` | no | `true` to skip TLS verification (self-signed certs) |
| `UNIFI_ALLOW_DESTRUCTIVE` | no | `true` to register ACL write, delete, and revoke tools (default: disabled) |

Source your `.env` file before running:

```bash
set -a && source .env && set +a
```

## Running

### stdio (default — for local MCP clients)

```bash
unifi-mcp
# or, if running directly from a source build:
./bin/unifi-mcp
```

### HTTP (streamable — for remote/shared deployments)

> The HTTP transport has no built-in authentication. Use a loopback address or place it behind a reverse proxy before exposing it on a shared network.

```bash
unifi-mcp --transport http --addr 127.0.0.1:8080
```

## VS Code Copilot configuration

Create `.vscode/mcp.json` in your workspace (already gitignored):

```json
{
  "servers": {
    "unifi-mcp": {
      "type": "stdio",
      "command": "/path/to/unifi-mcp/bin/unifi-mcp",
      "env": {
        "UNIFI_BASE_URL": "https://192.168.1.1/proxy/network",
        "UNIFI_API_KEY": "your-api-key",
        "UNIFI_SITE_ID": "your-site-uuid"
      }
    }
  }
}
```

Then open the Copilot chat panel, switch to **Agent** mode, and the `unifi-mcp` server will appear in the available tools.

## Claude Desktop configuration

Add the server to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "unifi-mcp": {
      "command": "/path/to/unifi-mcp/bin/unifi-mcp",
      "env": {
        "UNIFI_BASE_URL": "https://192.168.1.1/proxy/network",
        "UNIFI_API_KEY": "your-api-key",
        "UNIFI_SITE_ID": "your-site-uuid"
      }
    }
  }
}
```

Restart Claude Desktop after saving the config — the UniFi tools will appear in the tool selector.

## OpenCode configuration

Add the server to `opencode.json` in your project root (or `~/.config/opencode/opencode.json` for global config):

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "unifi-mcp": {
      "type": "local",
      "command": ["/path/to/unifi-mcp/bin/unifi-mcp"],
      "enabled": true,
      "environment": {
        "UNIFI_BASE_URL": "https://192.168.1.1/proxy/network",
        "UNIFI_API_KEY": "your-api-key",
        "UNIFI_SITE_ID": "your-site-uuid"
      }
    }
  }
}
```

## Development

```bash
make install-tools   # install golangci-lint, gosec, govulncheck, gofumpt
make check           # full quality gate: fix, fmt, vet, lint, sec, vulncheck, test, build
make test            # tests only (with race detector)
make build           # build only → bin/unifi-mcp
make clean           # remove bin/unifi-mcp
```
