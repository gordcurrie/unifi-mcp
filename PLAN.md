# UniFi MCP Server — Project Plan

## Summary

MCP server in Go exposing UniFi home network operations as tools.
Stack: `modelcontextprotocol/go-sdk`, custom UniFi HTTP client, API key auth, `golangci-lint` / `gosec` / `govulncheck`.

Only the v1 API (`/proxy/network/integration/v1/...`) is used. The legacy
`/api/s/{site}/...` path is not supported and will not be added — it is not
versioned, undocumented, and will be removed in a future Network release.

---

## Current tool inventory

| File          | Tool                          | Read-only |
|---------------|-------------------------------|-----------|
| `sites.go`    | `get_application_info`        | ✅        |
| `sites.go`    | `list_sites`                  | ✅        |
| `sites.go`    | `get_site`                    | ✅        |
| `devices.go`  | `list_devices`                | ✅        |
| `devices.go`  | `get_device`                  | ✅        |
| `devices.go`  | `get_device_stats`            | ✅        |
| `devices.go`  | `list_pending_devices`        | ✅        |
| `devices.go`  | `restart_device`              |           |
| `devices.go`  | `power_cycle_port`            |           |
| `clients.go`  | `list_clients`                | ✅        |
| `clients.go`  | `get_client`                  | ✅        |
| `clients.go`  | `authorize_guest_client`      |           |
| `network.go`  | `list_wifi_broadcasts`        | ✅        |
| `network.go`  | `get_wifi_broadcast`          | ✅        |
| `network.go`  | `set_wifi_broadcast_enabled`  |           |
| `network.go`  | `list_networks`               | ✅        |
| `network.go`  | `list_firewall_policies`      | ✅        |
| `network.go`  | `get_firewall_policy`         | ✅        |
| `network.go`  | `set_firewall_policy_enabled` |           |
| `network.go`  | `list_firewall_zones`         | ✅        |
| `network.go`  | `get_firewall_zone`           | ✅        |
| `network.go`  | `create_firewall_zone`        |           |
| `network.go`  | `update_firewall_zone`        |           |
| `network.go`  | `list_acl_rules`              | ✅        |
| `network.go`  | `get_acl_rule`                | ✅        |
| `network.go`  | `get_acl_rule_ordering`       | ✅        |
| `network.go`  | `list_traffic_matching_lists` | ✅        |
| `network.go`  | `get_traffic_matching_list`   | ✅        |
| `network.go`  | `list_dns_policies`           | ✅        |
| `network.go`  | `get_dns_policy`              | ✅        |
| `network.go`  | `create_dns_policy`           |           |
| `network.go`  | `update_dns_policy`           |           |
| `network.go`  | `list_wans`                   | ✅        |
| `network.go`  | `list_vpn_tunnels`            | ✅        |
| `network.go`  | `list_vpn_servers`            | ✅        |
| `network.go`  | `list_vouchers`               | ✅        |
| `network.go`  | `get_voucher`                 | ✅        |
| `network.go`  | `delete_voucher`              |           |
| `network.go`  | `create_vouchers`             |           |
| `network.go`  | `list_radius_profiles`        | ✅        |
| `network.go`  | `list_device_tags`            | ✅        |
| `network.go`  | `list_dpi_categories`         | ✅        |
| `network.go`  | `list_dpi_applications`       | ✅        |

Destructive tools (require `UNIFI_ALLOW_DESTRUCTIVE=true` + `confirmed: true`):
`restart_device`, `power_cycle_port`, `set_wifi_broadcast_enabled`,
`set_firewall_policy_enabled`, `create_firewall_zone`, `update_firewall_zone`,
`create_dns_policy`, `update_dns_policy`, `create_vouchers`, `delete_voucher`,
`authorize_guest_client`.

---

## Deferred (intentionally not implemented)

- `POST /v1/sites/{id}/firewall/policies` — create firewall policy (schema too complex; manual creation in UI is safer)
- `DELETE` on any resource except vouchers — too high blast radius for a home lab MCP; use the UI
- ACL rule write operations (`create_acl_rule`, `update_acl_rule`, `delete_acl_rule`, `reorder_acl_rules`, `set_acl_rule_enabled`) — any mutation directly controls traffic; deferred until there is a clear use case


