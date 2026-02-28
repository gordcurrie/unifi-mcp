package unifi

// ApplicationInfo is returned by GET /v1/info.
type ApplicationInfo struct {
	ApplicationVersion string `json:"applicationVersion"`
	ControllerType     string `json:"controllerType,omitempty"`
}

// Site is returned by GET /v1/sites and GET /v1/sites/{siteID}.
type Site struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Device is returned by GET /v1/sites/{siteID}/devices.
type Device struct {
	ID              string  `json:"id"`
	MAC             string  `json:"mac"`
	Name            string  `json:"name,omitempty"`
	Model           string  `json:"model,omitempty"`
	Type            string  `json:"type,omitempty"`
	State           int     `json:"state"`
	IP              string  `json:"ip,omitempty"`
	Version         string  `json:"version,omitempty"`
	Uptime          int64   `json:"uptime,omitempty"`
	Adopted         bool    `json:"adopted"`
	Disabled        bool    `json:"disabled"`
	UpdateAvailable bool    `json:"update_available,omitempty"`
	CPUUsage        float64 `json:"cpu_usage,omitempty"`
	MemUsage        float64 `json:"mem_usage,omitempty"`
	NumSta          int     `json:"num_sta,omitempty"`
	Experience      int     `json:"experience,omitempty"`
}

// ActiveClient is returned by GET /v1/sites/{siteID}/clients/active.
type ActiveClient struct {
	ID        string `json:"id"`
	MAC       string `json:"mac"`
	IP        string `json:"ip,omitempty"`
	Hostname  string `json:"hostname,omitempty"`
	Name      string `json:"name,omitempty"`
	Network   string `json:"network,omitempty"`
	SSID      string `json:"ssid,omitempty"`
	APName    string `json:"ap_name,omitempty"`
	Signal    int    `json:"signal,omitempty"`
	RxBytes   int64  `json:"rx_bytes"`
	TxBytes   int64  `json:"tx_bytes"`
	Uptime    int64  `json:"uptime"`
	IsWired   bool   `json:"is_wired"`
	Blocked   bool   `json:"blocked"`
	FirstSeen int64  `json:"first_seen,omitempty"`
	LastSeen  int64  `json:"last_seen,omitempty"`
}

// KnownClient is returned by GET /v1/sites/{siteID}/clients/history.
type KnownClient struct {
	ID        string `json:"id"`
	MAC       string `json:"mac"`
	Name      string `json:"name,omitempty"`
	Hostname  string `json:"hostname,omitempty"`
	IP        string `json:"ip,omitempty"`
	Network   string `json:"network,omitempty"`
	FirstSeen int64  `json:"first_seen,omitempty"`
	LastSeen  int64  `json:"last_seen,omitempty"`
	Blocked   bool   `json:"blocked"`
	IsWired   bool   `json:"is_wired"`
}

// SiteStats is returned by GET /v1/sites/{siteID}/statistics/site.
type SiteStats struct {
	NumAdopted      int     `json:"num_adopted"`
	NumDisconnected int     `json:"num_disconnected"`
	NumNew          int     `json:"num_new"`
	NumSta          int     `json:"num_sta"`
	NumUser         int     `json:"num_user"`
	NumGuest        int     `json:"num_guest"`
	WLANNumUser     int     `json:"wlan-num_user,omitempty"`
	LANNumUser      int     `json:"lan-num_user,omitempty"`
	WANTxBytes      int64   `json:"wan-tx_bytes,omitempty"`
	WANRxBytes      int64   `json:"wan-rx_bytes,omitempty"`
	LANTxBytes      int64   `json:"lan-tx_bytes,omitempty"`
	LANRxBytes      int64   `json:"lan-rx_bytes,omitempty"`
	WLANTxBytes     int64   `json:"wlan-tx_bytes,omitempty"`
	WLANRxBytes     int64   `json:"wlan-rx_bytes,omitempty"`
	Latency         int     `json:"latency,omitempty"`
	Uptime          int64   `json:"uptime,omitempty"`
	Drops           int     `json:"drops,omitempty"`
	TxBytes         int64   `json:"tx_bytes,omitempty"`
	RxBytes         int64   `json:"rx_bytes,omitempty"`
	XputUp          float64 `json:"xput_up,omitempty"`
	XputDown        float64 `json:"xput_down,omitempty"`
	SpeedtestPing   float64 `json:"speedtest_ping,omitempty"`
}

// SpeedTestStatus is returned by the devmgr speedtest-status command.
type SpeedTestStatus struct {
	Status     string  `json:"status"`
	LatencyAvg float64 `json:"latency_avg,omitempty"`
	XputDown   float64 `json:"xput_download,omitempty"`
	XputUp     float64 `json:"xput_upload,omitempty"`
	RunDate    int64   `json:"rundate,omitempty"`
}

// Event is returned by GET /api/s/{site}/stat/event.
type Event struct {
	ID        string `json:"_id"`
	Key       string `json:"key"`
	Msg       string `json:"msg"`
	Time      int64  `json:"time"`
	SiteID    string `json:"site_id,omitempty"`
	MAC       string `json:"mac,omitempty"`
	APName    string `json:"ap,omitempty"`
	Subsystem string `json:"subsystem,omitempty"`
	IsAdmin   bool   `json:"is_admin,omitempty"`
}

// Alarm is returned by GET /api/s/{site}/stat/alarm.
type Alarm struct {
	ID        string `json:"_id"`
	Key       string `json:"key"`
	Msg       string `json:"msg"`
	Time      int64  `json:"time"`
	Archived  bool   `json:"archived"`
	SiteID    string `json:"site_id,omitempty"`
	MAC       string `json:"mac,omitempty"`
	Subsystem string `json:"subsystem,omitempty"`
	Handled   bool   `json:"handled,omitempty"`
}

// WLAN is returned by GET /api/s/{site}/rest/wlanconf.
type WLAN struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Security    string `json:"security,omitempty"`
	WPAMode     string `json:"wpa_mode,omitempty"`
	VLANID      int    `json:"vlan_id,omitempty"`
	NetworkID   string `json:"networkconf_id,omitempty"`
	WLANGroupID string `json:"wlangroup_id,omitempty"`
	IsGuest     bool   `json:"is_guest,omitempty"`
	HideSSID    bool   `json:"hide_ssid,omitempty"`
}

// NetworkConf is returned by GET /api/s/{site}/rest/networkconf.
type NetworkConf struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Purpose     string `json:"purpose"`
	VLANID      int    `json:"vlan_id,omitempty"`
	IPSubnet    string `json:"ip_subnet,omitempty"`
	DhcpEnabled bool   `json:"dhcpd_enabled,omitempty"`
	DhcpStart   string `json:"dhcpd_start,omitempty"`
	DhcpStop    string `json:"dhcpd_stop,omitempty"`
	Enabled     bool   `json:"enabled"`
}

// FirewallRule is returned by GET /api/s/{site}/rest/firewallrule.
type FirewallRule struct {
	ID                  string   `json:"_id"`
	Name                string   `json:"name"`
	Enabled             bool     `json:"enabled"`
	Action              string   `json:"action"`
	Ruleset             string   `json:"ruleset"`
	RuleIndex           int      `json:"rule_index,omitempty"`
	Protocol            string   `json:"protocol,omitempty"`
	SrcFirewallGroupIDs []string `json:"src_firewallgroup_ids,omitempty"`
	DstFirewallGroupIDs []string `json:"dst_firewallgroup_ids,omitempty"`
	SrcAddress          string   `json:"src_address,omitempty"`
	DstAddress          string   `json:"dst_address,omitempty"`
	SrcPort             string   `json:"src_port,omitempty"`
	DstPort             string   `json:"dst_port,omitempty"`
	Logging             bool     `json:"logging,omitempty"`
	StateNew            bool     `json:"state_new,omitempty"`
	StateEstablished    bool     `json:"state_established,omitempty"`
	StateInvalid        bool     `json:"state_invalid,omitempty"`
	StateRelated        bool     `json:"state_related,omitempty"`
}

// PortForward is returned by GET /api/s/{site}/rest/portforward.
type PortForward struct {
	ID      string `json:"_id"`
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled"`
	Src     string `json:"src,omitempty"`
	DstPort string `json:"dst_port,omitempty"`
	Fwd     string `json:"fwd,omitempty"`
	FwdPort string `json:"fwd_port,omitempty"`
	Proto   string `json:"proto,omitempty"`
	Log     bool   `json:"log,omitempty"`
}

// deviceCmdRequest is the body sent to POST /api/s/{site}/cmd/devmgr.
type deviceCmdRequest struct {
	Cmd string `json:"cmd"`
	MAC string `json:"mac,omitempty"`
}

// clientCmdRequest is the body sent to POST /api/s/{site}/cmd/stamgr.
type clientCmdRequest struct {
	Cmd string `json:"cmd"`
	MAC string `json:"mac"`
}

// wlanEnabledRequest is the body sent to PUT /api/s/{site}/rest/wlanconf/{id}.
type wlanEnabledRequest struct {
	Enabled bool `json:"enabled"`
}
