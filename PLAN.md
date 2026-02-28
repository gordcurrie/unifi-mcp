# UniFi MCP Server — Project Plan

## Summary

Build an MCP server in Go that exposes UniFi home network operations as MCP tools.
Uses the official `modelcontextprotocol/go-sdk`, a custom UniFi HTTP client (no third-party
UniFi library), API key auth, and strict linting enforced via `golangci-lint`, `gosec`,
`govulncheck`, and `go fix`.

---

## Project Structure

```
unifi_mcp/
├── cmd/unifi-mcp/main.go          # entrypoint, CLI flags, signal handling
├── internal/unifi/
│   ├── client.go                  # HTTP client: API key header, TLS skip, base URL
│   ├── client_test.go             # httptest unit tests
│   ├── types.go                   # all request/response structs
│   ├── sites.go                   # /v1/info, /v1/sites
│   ├── devices.go                 # /v1/sites/{id}/devices + devmgr commands
│   ├── clients.go                 # /v1/sites/{id}/clients/active|history + stamgr
│   ├── statistics.go              # /v1/sites/{id}/statistics/* + events/alarms
│   └── network.go                 # legacy /api/s/{site}/rest/wlanconf, firewallrule, etc.
└── tools/
    ├── register.go                # RegisterAll + Config{AllowDestructive bool}
    ├── sites.go                   # site info tools
    ├── devices.go                 # device list/status + restart/locate/upgrade
    ├── clients.go                 # active/known clients + block/unblock/kick
    ├── statistics.go              # stats, events, alarms, speed test, DPI
    ├── network.go                 # WLANs, firewall rules, port forwards (read + limited edit)
    └── destructive.go             # forget_client, force_reprovision (opt-in)
```

---

## Environment Variables

| Variable                 | Required | Description                                                                 |
|--------------------------|----------|-----------------------------------------------------------------------------|
| `UNIFI_BASE_URL`         | yes      | e.g. `https://192.168.1.1/proxy/network`                                   |
| `UNIFI_API_KEY`          | yes      | API key from UniFi Site Manager Settings → API Keys                         |
| `UNIFI_SITE_ID`          | yes      | Site ID UUID from the Site Manager                                          |
| `UNIFI_INSECURE`         | no       | `true` to skip TLS verification (typical for local UCG self-signed cert)    |
| `UNIFI_ALLOW_DESTRUCTIVE`| no       | `true` to enable `forget_client` and `force_reprovision_device`             |

## CLI Flags

| Flag          | Default          | Description                                     |
|---------------|------------------|-------------------------------------------------|
| `--transport` | `stdio`          | `stdio` or `http`                               |
| `--addr`      | `localhost:8080` | Listen address when `--transport=http`          |

---

## API Coverage Notes

The UCG-Max uses the UDM-style local API with `https://<console-ip>/proxy/network` as
the base path. Two API path styles are used under the same client and API key:

- **New v1 API** (`/v1/sites/{siteId}/...`) — used for sites, devices, clients, and
  statistics. Official and documented at developer.ui.com.
- **Legacy local API** (`/api/s/{site}/...`) — used for management commands
  (`devmgr`, `stamgr`) and network config endpoints (WLANs, firewall rules, port
  forwards) that the v1 API does not yet expose.

Both path styles authenticate with the same `X-API-Key` bearer header on the UCG-Max.
`UNIFI_SITE_ID` sets a default site so tools don't require it as a parameter for
single-site home lab use, but it is still accepted as an optional input parameter.

---

## Tool Inventory (22 always-on + 2 opt-in destructive)

| File            | Tool                        | Read-only |
|-----------------|-----------------------------|-----------|
| `sites.go`      | `get_application_info`      | ✅        |
| `sites.go`      | `list_sites`                | ✅        |
| `sites.go`      | `get_site`                  | ✅        |
| `devices.go`    | `list_devices`              | ✅        |
| `devices.go`    | `get_device`                | ✅        |
| `devices.go`    | `restart_device`            |           |
| `devices.go`    | `locate_device`             |           |
| `devices.go`    | `unlocate_device`           |           |
| `devices.go`    | `upgrade_device_firmware`   |           |
| `devices.go`    | `run_speed_test`            |           |
| `devices.go`    | `get_speed_test_status`     | ✅        |
| `clients.go`    | `list_active_clients`       | ✅        |
| `clients.go`    | `list_known_clients`        | ✅        |
| `clients.go`    | `block_client`              |           |
| `clients.go`    | `unblock_client`            |           |
| `clients.go`    | `kick_client`               |           |
| `statistics.go` | `get_site_statistics`       | ✅        |
| `statistics.go` | `get_device_statistics`     | ✅        |
| `statistics.go` | `get_client_statistics`     | ✅        |
| `statistics.go` | `list_events`               | ✅        |
| `statistics.go` | `list_alarms`               | ✅        |
| `network.go`    | `list_wlans`                | ✅        |
| `network.go`    | `list_networks`             | ✅        |
| `network.go`    | `list_firewall_rules`       | ✅        |
| `network.go`    | `list_port_forwards`        | ✅        |
| `network.go`    | `enable_wlan`               |           |
| `network.go`    | `disable_wlan`              |           |
| destructive     | `forget_client`             | opt-in    |
| destructive     | `force_reprovision_device`  | opt-in    |

---

## Implementation Steps

### Step 1 — Bootstrap

- `go mod init github.com/gordcurrie/unifi-mcp`
- `go get github.com/modelcontextprotocol/go-sdk@v1.4.0`
- Create `Makefile` with `build`, `lint`, `test`, and `check` targets
- Create `.golangci.yml` mirroring the proxmox-mcp linter config

### Step 2 — `internal/unifi/client.go`

HTTP client struct with `baseURL`, `apiKey`, `siteID`, and `httpClient` (TLS skip
configurable). Authentication via `X-API-Key: <key>` header on every request.
Implement `get()`, `post()`, `postWithBody()`, `put()` helpers returning `[]byte` + error.
Convenience helpers `jsonResult()`, `textResult()` in `tools/helpers.go`.

### Step 3 — `internal/unifi/client_test.go`

`httptest`-based unit tests covering auth header injection, error propagation, and
JSON decoding.

### Step 4 — `internal/unifi/types.go`

Structs for all API responses:

- `Site`, `SiteHealth`
- `Device`, `DevicePort`, `DeviceRadio`
- `Client`, `ActiveClient`
- `SiteStats`, `DeviceStats`, `ClientStats`
- `Event`, `Alarm`
- `WLAN`, `NetworkConf`, `FirewallRule`, `PortForward`
- Command request types: `DeviceCommandRequest`, `ClientCommandRequest`

### Step 5 — `internal/unifi/sites.go`

- `GetInfo()` → `GET /v1/info`
- `ListSites()` → `GET /v1/sites`
- `GetSite(siteID string)` → `GET /v1/sites/{siteID}`

### Step 6 — `internal/unifi/devices.go`

v1 endpoints:
- `ListDevices(siteID string)` → `GET /v1/sites/{siteID}/devices`
- `GetDevice(siteID, deviceID string)` → `GET /v1/sites/{siteID}/devices/{deviceID}`

Legacy devmgr commands (`POST /api/s/{site}/cmd/devmgr`):
- `RestartDevice(site, mac string)`
- `LocateDevice(site, mac string)`
- `UnlocateDevice(site, mac string)`
- `UpgradeDevice(site, mac string)`
- `ForceReprovisionDevice(site, mac string)`

Legacy stat endpoints:
- `RunSpeedTest(site string)` → `POST /api/s/{site}/cmd/devmgr` `{"cmd":"speedtest"}`
- `GetSpeedTestStatus(site string)` → `POST /api/s/{site}/cmd/devmgr` `{"cmd":"speedtest-status"}`

### Step 7 — `internal/unifi/clients.go`

v1 endpoints:
- `ListActiveClients(siteID string)` → `GET /v1/sites/{siteID}/clients/active`
- `ListKnownClients(siteID string)` → `GET /v1/sites/{siteID}/clients/history`

Legacy stamgr commands (`POST /api/s/{site}/cmd/stamgr`):
- `BlockClient(site, mac string)`
- `UnblockClient(site, mac string)`
- `KickClient(site, mac string)` (disconnect, not ban)
- `ForgetClient(site, mac string)`

### Step 8 — `internal/unifi/statistics.go`

v1 endpoints:
- `GetSiteStats(siteID string)` → `GET /v1/sites/{siteID}/statistics/site`
- `GetDeviceStats(siteID string)` → `GET /v1/sites/{siteID}/statistics/devices`
- `GetClientStats(siteID string, macs []string)` → `GET /v1/sites/{siteID}/statistics/clients`

Legacy stat endpoints:
- `ListEvents(site string)` → `GET /api/s/{site}/stat/event`
- `ListAlarms(site string)` → `GET /api/s/{site}/stat/alarm`

### Step 9 — `internal/unifi/network.go`

Legacy REST endpoints (read + limited mutation):
- `ListWLANs(site string)` → `GET /api/s/{site}/rest/wlanconf`
- `SetWLANEnabled(site, wlanID string, enabled bool)` → `PUT /api/s/{site}/rest/wlanconf/{wlanID}`
- `ListNetworks(site string)` → `GET /api/s/{site}/rest/networkconf`
- `ListFirewallRules(site string)` → `GET /api/s/{site}/rest/firewallrule`
- `ListPortForwards(site string)` → `GET /api/s/{site}/rest/portforward`

### Step 10 — `tools/register.go`

```go
type Config struct {
    AllowDestructive bool
}

func RegisterAll(s *mcp.Server, client *unifi.Client, cfg Config) {
    registerSiteTools(s, client)
    registerDeviceTools(s, client)
    registerClientTools(s, client)
    registerStatisticsTools(s, client)
    registerNetworkTools(s, client)
    if cfg.AllowDestructive {
        registerDestructiveTools(s, client)
    }
}
```

### Step 11 — `tools/sites.go`

3 read-only tools: `get_application_info`, `list_sites`, `get_site`.
All annotated with `ReadOnlyHint: true`.

### Step 12 — `tools/devices.go`

- Read-only: `list_devices`, `get_device`, `get_speed_test_status`
- Mutating: `restart_device`, `locate_device`, `unlocate_device`,
  `upgrade_device_firmware`, `run_speed_test`

### Step 13 — `tools/clients.go`

- Read-only: `list_active_clients`, `list_known_clients`
- Mutating: `block_client`, `unblock_client`, `kick_client`

### Step 14 — `tools/statistics.go`

5 read-only tools: `get_site_statistics`, `get_device_statistics`,
`get_client_statistics`, `list_events`, `list_alarms`.

### Step 15 — `tools/network.go`

- Read-only: `list_wlans`, `list_networks`, `list_firewall_rules`,`list_port_forwards`
- Mutating: `enable_wlan`, `disable_wlan`

### Step 16 — `tools/destructive.go`

`forget_client` and `force_reprovision_device` — only registered when
`UNIFI_ALLOW_DESTRUCTIVE=true`. Both require `confirmed: true` field to proceed.
Both annotated with `DestructiveHint: true`.

### Step 17 — `cmd/unifi-mcp/main.go`

- Read env vars (`UNIFI_BASE_URL`, `UNIFI_API_KEY`, `UNIFI_SITE_ID`,
  `UNIFI_INSECURE`, `UNIFI_ALLOW_DESTRUCTIVE`)
- Parse `--transport` and `--addr` flags
- Create `unifi.Client`
- Create `mcp.Server` with name `unifi-mcp`
- Call `tools.RegisterAll`
- Run with `StdioTransport` or `StreamableHTTPHandler`
- Graceful shutdown on `SIGINT`/`SIGTERM`

### Step 18 — `README.md`

- Tool list
- Prerequisites (Go 1.21+, API key setup)
- Build instructions (`make build`)
- Environment variable table
- VS Code Copilot `.vscode/mcp.json` example
- Claude Desktop config example
- Development (`make check`)

---

## Verification

- `make build` → binary at `bin/unifi-mcp`
- `make check` → `golangci-lint`, `gosec`, `govulncheck` all pass
- `make test` → `httptest`-based client unit tests pass
- Manual: configure `.vscode/mcp.json` with real env vars, test `list_sites` and
  `list_active_clients` in Copilot Agent mode against the live UCG-Max
