# UniFi MCP Server — Project Plan

## Summary

MCP server in Go exposing UniFi home network operations as tools.
Stack: `modelcontextprotocol/go-sdk`, custom UniFi HTTP client, API key auth, `golangci-lint` / `gosec` / `govulncheck`.

---

## ✅ Completed — v0.1.0 + v0.2.0

### Infrastructure (v0.1.0)
- Go module `github.com/gordcurrie/unifi-mcp`, `modelcontextprotocol/go-sdk v1.4.0`
- `Makefile`: `build`, `fix`, `fmt`, `vet`, `lint`, `sec`, `vulncheck`, `test`, `check`
- `.golangci.yml` with `gosec`, `govet`, `staticcheck`, `errcheck`, `bodyclose`,
  `noctx`, `gofumpt`, `revive`, `gocritic`, `unparam`, `unconvert`
- `internal/unifi/client.go` — `get`, `post`, `postWithBody`, `put`; `X-API-Key` auth,
  TLS-skip opt-in, 10 MiB response cap, 30 s timeout
- `tools/helpers.go` — `jsonResult`, `textResult`
- `cmd/unifi-mcp/main.go` — env var bootstrap, `--transport stdio|http`, `--addr`,
  `SIGINT`/`SIGTERM` graceful shutdown
- `httptest`-based unit tests for client auth, error propagation, JSON decoding

### v1 API Rewrite (v0.2.0) — path prefix `/proxy/network/integration/v1`
- Dropped all legacy `/api/s/{site}/...` usage entirely
- `UNIFI_SITE_ID` is now a UUID (not slug); `ErrSiteNotFound` sentinel added
- **13 tools shipped and live-tested against UCG-Max running Network 10.1.85:**

| File              | Tool                    | Read-only |
|-------------------|-------------------------|-----------|
| `sites.go`        | `get_application_info`  | ✅        |
| `sites.go`        | `list_sites`            | ✅        |
| `sites.go`        | `get_site`              | ✅        |
| `devices.go`      | `list_devices`          | ✅        |
| `devices.go`      | `get_device`            | ✅        |
| `devices.go`      | `get_device_stats`      | ✅        |
| `devices.go`      | `restart_device`        |           |
| `clients.go`      | `list_clients`          | ✅        |
| `network.go`      | `list_wifi_broadcasts`  | ✅        |
| `network.go`      | `list_networks`         | ✅        |
| `network.go`      | `list_firewall_policies`| ✅        |
| `network.go`      | `list_firewall_zones`   | ✅        |
| `network.go`      | `list_acl_rules`        | ✅        |

- `.env.example`, `README.md` tool tables, `PLAN.md` all kept in sync per PR convention
- `make check` clean; tagged `v0.2.0`; GitHub release published

---

## ✅ Phase 3 — MCP Correctness & Security Hardening

Findings from MCP spec review (2025-06-18). Must pass `make check` before committing.
Target tag: `v0.3.0`.

### 3a — isError semantics (spec §6) ✅

MCP spec distinguishes protocol errors (unknown tool, bad schema) from tool execution
errors. API/business failures must be returned as `isError: true` results, not as Go
errors that bubble up as protocol-level failures.

- Add `errorResult(err error)` helper to `tools/helpers.go`
- Change every `return nil, nil, fmt.Errorf("tool: %w", err)` in all tool files to
  `return errorResult(fmt.Errorf("tool: %w", err))`
- Files: `tools/devices.go`, `tools/sites.go`, `tools/clients.go`, `tools/network.go`

### 3b — Required input validation (spec §7) ✅

Spec §7: servers MUST validate all tool inputs. Empty `device_id` produces malformed
URLs (`.../devices//statistics/latest`) and confusing downstream errors.

- In `get_device`, `restart_device`, `get_device_stats` handlers, guard:
  ```go
  if input.DeviceID == "" {
      return errorResult(fmt.Errorf("get_device: device_id is required"))
  }
  ```

### 3c — `DestructiveHint` on `restart_device` ✅

`ReadOnlyHint: false` is not the same as `DestructiveHint: true`. The spec uses
`DestructiveHint` as the mechanism for MCP clients to trigger a confirmation step
("human in the loop") before invocation.

- Change `restart_device` annotation to `Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue}`

### 3d — HTTP transport auth warning ✅

When `--transport http` is used there is no authentication layer. Anyone reaching the
port can invoke `restart_device`. Add a visible runtime warning.

- In `cmd/unifi-mcp/main.go` `http` case, before `ListenAndServe`:
  ```go
  slog.Warn("HTTP transport has no authentication — restrict network access to trusted hosts only")
  ```

---

## Phase 4 — Expand Functionality

Candidates in rough priority order. Each item requires `make check` + README/PLAN.md
updates before merging.

> **API path note:** All `/v1/...` paths below are relative to the client's `baseURL`,
> which is `https://<console>/proxy/network/integration`. The full resolved path is
> `https://<console>/proxy/network/integration/v1/...` — identical to the existing tools.

### ✅ 4a — WiFi broadcast mutation

Enable/disable a WiFi broadcast (WLAN on/off for a given SSID).
- `GET /v1/sites/{id}/wifi/broadcasts/{broadcastId}` — get single ✅ (shipped in Phase 4 PR 1)
- `PUT /v1/sites/{id}/wifi/broadcasts/{broadcastId}` — full update; set `enabled` field
- New tools: `set_wifi_broadcast_enabled` (mutating, `confirmed bool`)

### 4b — Client lookup by ID ✅

Direct single-client lookup is confirmed in the API (`GET /v1/sites/{id}/clients/{clientId}`).
The v1 filtering syntax also supports querying by property (e.g. `macAddress.eq(...)`).
- New tool: `get_client` (read-only, `client_id` required)

### ✅ 4c — Port power cycle (PoE reboot)

Power-cycle a single PoE port on a switch without restarting the whole device.
Very useful for rebooting cameras, APs, or other PoE devices.
- `POST /v1/sites/{id}/devices/{deviceId}/interfaces/ports/{portIdx}/actions` body `{"action":"POWER_CYCLE"}`
- New tool: `power_cycle_port` (mutating, `confirmed bool`, `device_id`, `port_idx`)

### 4d — Pending device discovery ✅

List devices visible on the network but not yet adopted.
- `GET /v1/pending-devices`
- New tool: `list_pending_devices` (read-only)

### ✅ 4e — DNS policies

Local DNS A-record management — map hostnames to IPs on the site (e.g. `nas.home → 192.168.1.100`).
High value for home lab.
- `GET /v1/sites/{id}/dns/policies` — list
- `GET /v1/sites/{id}/dns/policies/{id}` — get single
- `POST /v1/sites/{id}/dns/policies` — create (`type`, `domain`, `ipv4Address`, `ttlSeconds`, `enabled`)
- `PUT /v1/sites/{id}/dns/policies/{id}` — update
- `DELETE /v1/sites/{id}/dns/policies/{id}` — delete (destructive, `confirmed bool`)
- New tools: `list_dns_policies`, `get_dns_policy`, `create_dns_policy`, `update_dns_policy`, `delete_dns_policy`

### ✅ 4f — Firewall policy & zone CRUD

Extend beyond read-only to full create/enable-disable/delete.
The API also provides `PATCH` for just `loggingEnabled` and a dedicated ordering endpoint.
- `PUT /v1/sites/{id}/firewall/policies/{id}` — full update (includes `enabled`) ✅
- `DELETE /v1/sites/{id}/firewall/policies/{id}` — delete (destructive, `confirmed bool`, requires `UNIFI_ALLOW_DESTRUCTIVE`) ✅
- `POST /v1/sites/{id}/firewall/zones` — create custom zone ✅
- `PUT /v1/sites/{id}/firewall/zones/{id}` — update zone (`name`, `networkIds`) ✅
- `DELETE /v1/sites/{id}/firewall/zones/{id}` — delete custom zone (destructive, `confirmed bool`, requires `UNIFI_ALLOW_DESTRUCTIVE`) ✅
- Deferred: `POST /v1/sites/{id}/firewall/policies` — create (complex schema)
- Deferred: `PATCH /v1/sites/{id}/firewall/policies/{id}` — partial update
- Deferred: `GET/PUT /v1/sites/{id}/firewall/policies/ordering` — reorder policies
- New tools: `get_firewall_policy`, `set_firewall_policy_enabled`, `delete_firewall_policy`, `get_firewall_zone`, `create_firewall_zone`, `update_firewall_zone`, `delete_firewall_zone`

### ✅ 4g — ACL rule CRUD

Extend beyond read-only to full create/enable-disable/delete/reorder.
- `POST /v1/sites/{id}/acl-rules` — create
- `PUT /v1/sites/{id}/acl-rules/{id}` — full update (includes `enabled`)
- `DELETE /v1/sites/{id}/acl-rules/{id}` — delete (destructive, `confirmed bool`, requires `UNIFI_ALLOW_DESTRUCTIVE`)
- `GET/PUT /v1/sites/{id}/acl-rules/ordering` — reorder
- New tools: `get_acl_rule`, `create_acl_rule`, `update_acl_rule`, `delete_acl_rule`, `reorder_acl_rules`, `set_acl_rule_enabled`, `get_acl_rule_ordering`

> **Safety note:** All ACL write tools (`create`, `update`, `set_enabled`, `reorder`, `delete`)
> are gated on `UNIFI_ALLOW_DESTRUCTIVE=true`, not just delete. Unlike firewall zones
> (organisational containers), any ACL mutation directly controls traffic flow; a misplaced
> BLOCK rule or a bad reorder can cause a full network outage. The flag is the primary guard;
> `confirmed: true` is secondary.

### ✅ 4h — Traffic matching lists (read-only)

Port/IP address lists referenced by firewall policies. Read-only surface is enough to
let the AI understand existing policy configurations.
- `GET /v1/sites/{id}/traffic-matching-lists` — list
- `GET /v1/sites/{id}/traffic-matching-lists/{id}` — get single
- New tools: `list_traffic_matching_lists`, `get_traffic_matching_list`

### ✅ 4i — WAN interfaces & VPN (read-only)

Useful for context when reasoning about firewall rules and routing.
- `GET /v1/sites/{id}/wans` — list WAN interface definitions
- `GET /v1/sites/{id}/vpn/site-to-site-tunnels` — list S2S tunnels
- `GET /v1/sites/{id}/vpn/servers` — list VPN servers
- New tools: `list_wans`, `list_vpn_tunnels`, `list_vpn_servers`

### ✅ 4j — Hotspot vouchers

Full CRUD for guest Hotspot vouchers (time/data-limited access codes).
Useful if running a guest or hotspot network.
- `GET /v1/sites/{id}/hotspot/vouchers` — list
- `GET /v1/sites/{id}/hotspot/vouchers/{id}` — get single
- `POST /v1/sites/{id}/hotspot/vouchers` — generate (`count`, `name`, `timeLimitMinutes`, optional limits)
- `DELETE /v1/sites/{id}/hotspot/vouchers/{id}` — revoke single (destructive, `confirmed bool`, requires `UNIFI_ALLOW_DESTRUCTIVE`)
- New tools: `list_vouchers`, `get_voucher`, `create_vouchers`, `delete_voucher` (destructive, `confirmed bool`)

### 4k — Guest client authorization

Authorize or revoke guest network access for a connected client.
- `POST /v1/sites/{id}/clients/{clientId}/actions` body `{"action":"AUTHORIZE_GUEST_ACCESS", ...}`
- New tool: `authorize_guest_client` (mutating, `confirmed bool`, optional time/data/rate limits)

### 4l — Supporting reference data (read-only)

Low-priority but useful for firewall policy creation context.
- `GET /v1/sites/{id}/device-tags` — device tags (used in WiFi broadcast device filters)
- `GET /v1/dpi/categories` + `/v1/dpi/applications` — DPI app categories (firewall matching)
- `GET /v1/sites/{id}/radius/profiles` — RADIUS profiles
- New tools: `list_device_tags`, `list_dpi_categories`, `list_dpi_applications`


