package unifi

// ApplicationInfo is returned by GET /integration/v1/info.
type ApplicationInfo struct {
	ApplicationVersion string `json:"applicationVersion"`
}

// Site is returned by GET /integration/v1/sites.
type Site struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	InternalReference string `json:"internalReference,omitempty"`
	Description       string `json:"description,omitempty"`
}

// Device is returned by GET /integration/v1/sites/{siteId}/devices.
type Device struct {
	ID                string `json:"id"`
	MAC               string `json:"macAddress"`
	IP                string `json:"ipAddress,omitempty"`
	Name              string `json:"name,omitempty"`
	Model             string `json:"model,omitempty"`
	State             string `json:"state"`
	FirmwareVersion   string `json:"firmwareVersion,omitempty"`
	FirmwareUpdatable bool   `json:"firmwareUpdatable"`
	AdoptedAt         string `json:"adoptedAt,omitempty"`
	ProvisionedAt     string `json:"provisionedAt,omitempty"`
}

// DeviceStats is returned by GET /integration/v1/sites/{siteId}/devices/{deviceId}/statistics/latest.
type DeviceStats struct {
	UptimeSec            int64   `json:"uptimeSec"`
	LastHeartbeatAt      string  `json:"lastHeartbeatAt,omitempty"`
	NextHeartbeatAt      string  `json:"nextHeartbeatAt,omitempty"`
	LoadAverage1Min      float64 `json:"loadAverage1Min"`
	LoadAverage5Min      float64 `json:"loadAverage5Min"`
	LoadAverage15Min     float64 `json:"loadAverage15Min"`
	CPUUtilizationPct    float64 `json:"cpuUtilizationPct"`
	MemoryUtilizationPct float64 `json:"memoryUtilizationPct"`
}

// deviceActionRequest is the body sent to POST /integration/v1/sites/{siteId}/devices/{deviceId}/actions.
type deviceActionRequest struct {
	Action string `json:"action"`
}

// NetworkClient is returned by GET /integration/v1/sites/{siteId}/clients.
type NetworkClient struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	Name           string `json:"name,omitempty"`
	ConnectedAt    string `json:"connectedAt,omitempty"`
	IP             string `json:"ipAddress,omitempty"`
	MAC            string `json:"macAddress"`
	UplinkDeviceID string `json:"uplinkDeviceId,omitempty"`
}

// PendingDevice is returned by GET /integration/v1/pending-devices.
type PendingDevice struct {
	ID              string `json:"id"`
	MAC             string `json:"macAddress"`
	IP              string `json:"ipAddress,omitempty"`
	Model           string `json:"model,omitempty"`
	FirmwareVersion string `json:"firmwareVersion,omitempty"`
	State           string `json:"state,omitempty"`
}

// WiFiBroadcast is returned by GET /integration/v1/sites/{siteId}/wifi/broadcasts.
type WiFiBroadcast struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// NetworkConf is returned by GET /integration/v1/sites/{siteId}/networks.
type NetworkConf struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Enabled    bool   `json:"enabled"`
	VLANID     int    `json:"vlanId,omitempty"`
	Management string `json:"management,omitempty"`
	Default    bool   `json:"default,omitempty"`
}

// FirewallPolicy is returned by GET /integration/v1/sites/{siteId}/firewall/policies.
type FirewallPolicy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description,omitempty"`
	Index       int    `json:"index"`
}

// FirewallZone is returned by GET /integration/v1/sites/{siteId}/firewall/zones.
type FirewallZone struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	NetworkIDs []string `json:"networkIds,omitempty"`
}

// ACLRule is returned by GET /integration/v1/sites/{siteId}/acl-rules.
type ACLRule struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Action  string `json:"action"`
	Index   int    `json:"index"`
}

// TrafficMatchingList is returned by GET /integration/v1/sites/{siteId}/traffic-matching-lists.
type TrafficMatchingList struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Type    string   `json:"type,omitempty"`
	Entries []string `json:"entries,omitempty"`
}

// WAN is returned by GET /integration/v1/sites/{siteId}/wans.
type WAN struct {
	ID        string   `json:"id"`
	Name      string   `json:"name,omitempty"`
	Type      string   `json:"type,omitempty"`
	Enabled   bool     `json:"enabled"`
	State     string   `json:"state,omitempty"`
	IPAddress string   `json:"ipAddress,omitempty"`
	Gateway   string   `json:"gateway,omitempty"`
	DNS       []string `json:"dns,omitempty"`
}

// VPNTunnel is returned by GET /integration/v1/sites/{siteId}/vpn/site-to-site-tunnels.
type VPNTunnel struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Enabled  bool   `json:"enabled"`
	State    string `json:"state,omitempty"`
	LocalIP  string `json:"localIp,omitempty"`
	RemoteIP string `json:"remoteIp,omitempty"`
}

// VPNServer is returned by GET /integration/v1/sites/{siteId}/vpn/servers.
type VPNServer struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Enabled bool   `json:"enabled"`
}

// DNSPolicy is returned by GET /integration/v1/sites/{siteId}/dns/policies.
type DNSPolicy struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Domain      string `json:"domain"`
	IPv4Address string `json:"ipv4Address,omitempty"`
	TTLSeconds  int    `json:"ttlSeconds,omitempty"`
	Enabled     bool   `json:"enabled"`
}

// DNSPolicyRequest is the body for POST and PUT to /integration/v1/sites/{siteId}/dns/policies.
type DNSPolicyRequest struct {
	Type        string `json:"type"`
	Domain      string `json:"domain"`
	IPv4Address string `json:"ipv4Address,omitempty"`
	TTLSeconds  int    `json:"ttlSeconds,omitempty"`
	Enabled     bool   `json:"enabled"`
}
