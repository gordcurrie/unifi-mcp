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
| `--addr`      | `:8080`          | Listen address when `--transport=http`          |

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

## Tool Inventory (27 always-on + 2 opt-in destructive)

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
| `devices.go`    | `upgrade_device`            |           |
| `devices.go`    | `run_speed_test`            |           |
| `devices.go`    | `get_speed_test_status`     | ✅        |
| `clients.go`    | `list_active_clients`       | ✅        |
| `clients.go`    | `list_known_clients`        | ✅        |
| `clients.go`    | `block_client`              |           |
| `clients.go`    | `unblock_client`            |           |
| `clients.go`    | `kick_client`               |           |
| `statistics.go` | `get_site_stats`            | ✅        |
| `statistics.go` | `get_device_stats`          | ✅        |
| `statistics.go` | `get_client_stats`          | ✅        |
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
- `RunSpeedTest(site string)` → `POST /api/s/{site}/cmd/devmgr/speedtest`
- `GetSpeedTestStatus(site string)` → `GET /api/s/{site}/stat/speedtest-status`

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
- `GetDeviceStats(siteID, deviceID string)` → `GET /v1/sites/{siteID}/statistics/devices/{deviceID}`
- `GetClientStats(siteID, clientID string)` → `GET /v1/sites/{siteID}/statistics/clients/{clientID}`

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

---

## Phase 2 — v1-Only Rewrite

### Decision: Drop legacy API support (Option A)

Live testing against UCG-Max running Network 10.1.85 revealed that the v1 API path
prefix is `/integration/v1/` (not `/v1/`), and the v1 API surface is comprehensive
enough for all read and most action use cases. Rather than maintain dual-mode routing
(UUID site IDs for v1, slug site IDs for legacy), we drop all legacy API usage entirely.

**Tools dropped** (no v1 equivalent):
- `locate_device`, `unlocate_device` — devmgr `set-locate`
- `run_speed_test`, `get_speed_test_status` — devmgr `speedtest`
- `block_client`, `unblock_client`, `kick_client` — stamgr (v1 client actions limited to guest auth)
- `forget_client`, `force_reprovision_device` — devmgr / stamgr destructive
- `list_known_clients` — `stat/alluser`
- `list_events`, `list_alarms` — `stat/event`, `stat/alarm`
- `get_site_stats`, `get_client_stats` — no v1 stats equivalent
- `set_wlan_enabled` — v1 WiFi Broadcasts update is complex; deferred to Phase 3

**Tools renamed/replaced:**
- `list_wlans` → `list_wifi_broadcasts` (path: `/v1/sites/{id}/wifi/broadcasts`)
- `list_firewall_rules` → `list_firewall_policies` + `list_firewall_zones` (path: `/v1/sites/{id}/firewall/...`)
- `list_port_forwards` — removed (no v1 equivalent; deferred)
- `get_device_stats` — re-pointed at `/v1/sites/{id}/devices/{deviceId}/statistics/latest`

**New tool added:**
- `list_acl_rules` → `GET /v1/sites/{id}/acl-rules`

### Discoveries

| Area | Assumed | Actual |
|---|---|---|
| v1 path prefix | `/proxy/network/v1/...` | `/proxy/network/integration/v1/...` |
| Site ID format | slug `"default"` | UUID only |
| `GET /v1/info` response | `{"data": {...}}` | raw `{"applicationVersion": "..."}` |
| `GET /v1/sites/{id}` | exists | **404** — list+filter only |
| Device fields | snake_case | camelCase (`macAddress`, `ipAddress`, `firmwareVersion`) |
| Device `state` | `int` | `string` (`"ONLINE"`, `"OFFLINE"`) |
| Active clients | `/clients/active` | `/clients` (all connected) |

All phases must pass `make check` before committing.

---

### Phase 2a — client.go: path prefix + site ID simplification ✅

Remove `decodeV1Single` envelope wrapper (GetInfo returns raw object now).
Remove `decodeLegacy`, `checkLegacyRC`, `legacyMeta`, `legacyResponse` — no legacy calls remain.
Rename internal v1 route prefix to `/integration/v1`.
`UNIFI_SITE_ID` must now be the UUID from the console. Document in `.env.example`.

---

### Phase 2b — sites.go + types ✅

- `GetInfo` — decode raw JSON (not `{"data":…}` envelope)
- `ListSites` — path `/integration/v1/sites`; `Site` struct gains `InternalReference`
- `GetSite` — implement as `ListSites` + filter by ID (no single-get endpoint exists)

---

### Phase 2c — devices.go + types ✅

- `ListDevices` / `GetDevice` — path `/integration/v1/sites/{id}/...`
- `GetDeviceStats` — new method: `GET /integration/v1/sites/{id}/devices/{deviceId}/statistics/latest`
- `RestartDevice` — re-point to `POST /integration/v1/sites/{id}/devices/{deviceId}/actions` body `{"action":"RESTART"}`
- Remove: `LocateDevice`, `UnlocateDevice`, `UpgradeDevice`, `ForceReprovisionDevice`, `RunSpeedTest`, `GetSpeedTestStatus`
- `Device` struct: camelCase JSON tags, `State string`, `FirmwareUpdatable bool`
- New `DeviceStats` struct matching `/statistics/latest` response

---

### Phase 2d — clients.go + types ✅

- `ListClients` — path `/integration/v1/sites/{id}/clients` (rename from `ListActiveClients`)
- Remove `ListKnownClients`
- `Client` struct (rename from `ActiveClient`): camelCase fields — `macAddress`, `ipAddress`, `connectedAt`, `type`, `uplinkDeviceId`

---

### Phase 2e — network.go + types ✅

Remove `SetWLANEnabled`, `ListFirewallRules`, `ListPortForwards`.
- `ListWiFiBroadcasts` — `GET /integration/v1/sites/{id}/wifi/broadcasts`
- `ListNetworks` — `GET /integration/v1/sites/{id}/networks` (camelCase: `vlanId`, `enabled`)
- `ListFirewallPolicies` — `GET /integration/v1/sites/{id}/firewall/policies`
- `ListFirewallZones` — `GET /integration/v1/sites/{id}/firewall/zones`
- `ListACLRules` — `GET /integration/v1/sites/{id}/acl-rules`
- Update `WiFiBroadcast`, `NetworkConf`, new `FirewallPolicy`, `FirewallZone`, `ACLRule` structs

---

### Phase 2f — statistics.go ✅

Remove `GetSiteStats`, `GetClientStats`, `ListEvents`, `ListAlarms` (all legacy).
Keep only `GetDeviceStats` (v1, moved to devices.go).
Delete `statistics.go` entirely (or repurpose in a later phase).

---

### Phase 2g — tools/ layer ✅

Remove tool registrations for all dropped methods.
Rename/update tool registrations to match new method names.
Add `list_wifi_broadcasts`, `list_firewall_policies`, `list_firewall_zones`, `list_acl_rules`.
Remove `destructive.go` tools that have no v1 equivalent.
Update `client_iface.go`.

---

### Phase 2h — tests + make check + commit ✅

Update all test path assertions and fixture JSON to match new paths and camelCase fields.
`make check` must be clean.
Rebuild binary, restart VS Code MCP server, smoke-test each tool live.
Commit + tag `v0.2.0`.

---

## Phase 3 — Legacy Action Tools (deferred)

Re-add legacy-only tools behind a compile tag or `UNIFI_ALLOW_LEGACY=true`:
- `locate_device`, `unlocate_device`
- `run_speed_test`, `get_speed_test_status`
- `block_client`, `unblock_client`, `kick_client`, `forget_client`
- `list_events`, `list_alarms`
- `set_wlan_enabled`

Requires dual site ID routing (UUID for v1, slug for legacy) implemented cleanly.
