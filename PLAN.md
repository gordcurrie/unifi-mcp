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
      return errorResult(fmt.Errorf("device_id is required"))
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

### 4a — WiFi broadcast mutation

Enable/disable a WiFi broadcast (WLAN on/off for a given SSID).
- `PUT /integration/v1/sites/{id}/wifi/broadcasts/{broadcastId}` — partial update with `enabled` field
- New tool: `set_wifi_broadcast_enabled` (mutating, `confirmed bool`)

### 4b — Device statistics history

Time-range queries on device stats instead of just `statistics/latest`.
- `GET /integration/v1/sites/{id}/devices/{deviceId}/statistics` with start/end params
- New tool: `get_device_stats_history`

### 4c — Client lookup by MAC / IP

Allow looking up a specific connected client without listing all of them.
- Explore if `GET /integration/v1/sites/{id}/clients?mac=...` filtering is supported
- New tool: `get_client` (read-only)

### 4d — FirewallPolicy mutation

Enable/disable a firewall policy rule.
- `PUT /integration/v1/sites/{id}/firewall/policies/{id}` — partial update with `enabled`
- New tool: `set_firewall_policy_enabled` (mutating, `confirmed bool`)

### 4e — ACL rule mutation

Enable/disable an ACL rule.
- `PUT /integration/v1/sites/{id}/acl-rules/{id}` — partial update with `enabled`
- New tool: `set_acl_rule_enabled` (mutating, `confirmed bool`)


