---
name: unifi-network
description: >
  Manage and inspect a UniFi home or small-office network using the unifi-mcp
  tools. Use this skill when asked about UniFi devices, access points, switches,
  gateways, connected clients, WiFi networks (SSIDs), VLANs, firewall policies
  or zones, ACL rules, local DNS records, VPNs, or hotspot vouchers — for
  example "what's connected to my network", "is my AP healthy", "turn off the
  guest WiFi", "add a DNS record for my NAS", or "restart the switch".
license: MIT
---

# Skill: UniFi Network (via unifi-mcp)

This skill tells you how to use the tools exposed by
[gordcurrie/unifi-mcp](https://github.com/gordcurrie/unifi-mcp) safely and
efficiently. For a full security review, use the `audit-network-security`
skill instead.

---

## Before you start

The tools are only available if the unifi-mcp server is configured in your MCP
client. If no `list_devices` / `list_clients` tools are available, stop and
point the user at the setup section of the
[README](https://github.com/gordcurrie/unifi-mcp#installation). Do not try to
reach the UniFi API directly with curl unless the user asks you to.

If a call fails with a TLS or certificate error, the console is most likely
using its default self-signed certificate: the user needs `UNIFI_INSECURE=true`.
If it fails with `401`/`403`, the API key is wrong or was revoked.

---

## Ground rules

1. **Read before you write.** Look up the current state with a `list_*` or
   `get_*` tool before changing anything, and show the user what will change.
2. **IDs, not names.** Tools that act on an existing resource take its UUID
   (`device_id`, `client_id`, `broadcast_id`, `policy_id`, …). Resolve the name
   to an ID from the matching `list_*` result first. Never guess or construct an
   ID. Create tools such as `create_dns_policy` and `create_firewall_zone` need
   no existing ID.
3. **Every change needs the user's approval.** Get explicit approval for each
   mutation (create, update, enable/disable, restart, delete) before calling
   the tool. A previous approval does not carry over to a different change.
   Many mutating tools also declare a `confirmed` parameter: set it to `true`
   only after that approval. Do not send `confirmed` to tools that do not
   declare it.
4. **Paginate fully.** `list_*` tools return `data`, `totalCount`, `offset` and
   `limit`. If `totalCount` is larger than the number of items you have, call
   again with `offset` increased by the number of items returned until you have
   them all. Do not answer "how many" or "is X present" questions from a
   partial page.
5. **Default site.** Omit `site_id` unless the user names a different site; the
   server falls back to `UNIFI_SITE_ID`. Use `list_sites` for multi-site setups.

---

## Common tasks

### What is connected to my network?

1. `list_clients` (paginate). Each client has `name`, `ipAddress`,
   `macAddress`, `type` (`WIRED` or `WIRELESS`) and `uplinkDeviceId`.
2. To say which AP or switch a client is on, match `uplinkDeviceId` against the
   `id` field from `list_devices`.
3. Present a table sorted by name; flag clients with no name as unidentified.

### Is a device healthy?

1. `list_devices` to find the device and check `state` and `firmwareUpdatable`.
2. `get_device_stats` for CPU %, memory %, load averages and uptime.
3. Flag: `state` other than `ONLINE`, CPU or memory above ~80 %, very short uptime
   (recent unexpected reboot), or available firmware updates.

### Restart a device or power-cycle a PoE port

`restart_device` and `power_cycle_port` disconnect everything behind that
device or port. Before calling, name the device (and port index), say what will
drop offline, and get explicit approval. Restarting the gateway drops the whole
network, including possibly the connection you are running over.

### Turn a WiFi network on or off

1. `list_wifi_broadcasts` to find the SSID and its `id`.
2. Confirm with the user, then `set_wifi_broadcast_enabled`. Warn if the user
   may be connected over that SSID.

### Local DNS records

- `list_dns_policies` to see existing A records; check for a duplicate domain
  before creating one.
- `create_dns_policy` / `update_dns_policy` to add or change a record.
- Deleting a record requires destructive tools (see below).

### Firewall and ACL

- `list_firewall_policies` returns user-defined policies by default; pass
  `user_only: false` to include system policies.
- `list_firewall_zones` and `list_networks` explain which networks a policy
  applies to.
- ACL rules are evaluated in order: always read `get_acl_rule_ordering`
  alongside `list_acl_rules`.
- A wrong firewall or ACL change can cut off part or all of the network. State
  the exact effect before asking for confirmation.

### Guest access and hotspot vouchers

- `authorize_guest_client` grants a client guest access, optionally with time,
  data and bandwidth limits.
- `create_vouchers` generates 1–100 vouchers; `list_vouchers` shows existing
  ones.

---

## Destructive tools

Delete operations and all ACL writes (`create_acl_rule`, `update_acl_rule`,
`set_acl_rule_enabled`, `reorder_acl_rules`, `delete_*`) are only registered
when the server runs with `UNIFI_ALLOW_DESTRUCTIVE=true`. If the user asks for
one of these and the tool is missing, explain that it is disabled by design and
how to enable it. Do not look for a workaround.

---

## Limits

unifi-mcp uses the UniFi Network **v1 Integration API** only. It cannot read
traffic statistics, DPI usage per client, event or alarm logs, port forwards,
or controller settings. If asked for these, say so and suggest the UniFi web UI.
