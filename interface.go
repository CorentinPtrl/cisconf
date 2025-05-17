package cisconf

type CiscoInterfaceParent struct {
	Identifier   string `reg:"interface ([\\w\\/\\.\\-\\:]+)" cmd:" %s"`
	SubInterface int    `reg:"interface [\\w\\/\\.\\-\\:]+\\.(\\d+)?" cmd:" %d"`
}

type CiscoInterface struct {
	Parent                CiscoInterfaceParent `reg:"interface.*" cmd:"interface" parent:"true"`
	AccessVlan            int                  `reg:"switchport access vlan ([0-9]+)" cmd:"switchport access vlan %d"`
	Access                bool                 `reg:"switchport mode access" cmd:"switchport mode access"`
	VoiceVlan             int                  `reg:"switchport voice vlan ([0-9]+)" cmd:"switchport voice vlan %d"`
	PortSecurityMaximum   int                  `reg:"switchport port-security maximum ([0-9]+)" cmd:"switchport port-security maximum %d"`
	PortSecurityViolation string               `reg:"switchport port-security violation (protect|restrict|shutdown)" cmd:"switchport port-security violation %s"`
	PortSecurityAgingTime int                  `reg:"switchport port-security aging time ([0-9]+)" cmd:"switchport port-security aging time %d"`
	PortSecurityAgingType string               `reg:"switchport port-security aging type (absolute|inactivity)" cmd:"switchport port-security aging type %s"`
	PortSecurity          bool                 `reg:"switchport port-security" cmd:"switchport port-security"`
	Description           string               `reg:"description ([[:print:]]+)" cmd:"description %s"`
	NativeVlan            int                  `reg:"switchport trunk native vlan ([0-9]+)" cmd:"switchport trunk native vlan %d"`
	Encapsulation         string               `reg:"switchport trunk encapsulation ([[:print:]]+)" cmd:"switchport trunk encapsulation %s"`
	Trunk                 bool                 `reg:"switchport mode trunk" cmd:"switchport mode trunk"`
	TrunkAllowedVlan      []int                `reg:"switchport trunk allowed vlan( add)? ([\\d,-]+)" cmd:"switchport trunk allowed vlan %s"`
	Shutdown              bool                 `reg:"shutdown" cmd:"shutdown" default:"false"`
	SCBroadcastLevel      float64              `reg:"storm-control broadcast level ([0-9\\.]+)" cmd:"storm-control broadcast level %.2f"`
	STPPortFast           string               `reg:"spanning-tree portfast (disable|edge|network)" cmd:"spanning-tree portfast %s"`
	STPBpduGuard          string               `reg:"spanning-tree bpduguard (disable|enable)" cmd:"spanning-tree bpduguard %s"`
	ServicePolicyInput    string               `reg:"service-policy input ([[:print:]]+)" cmd:"service-policy input %s"`
	ServicePolicyOutput   string               `reg:"service-policy output ([[:print:]]+)" cmd:"service-policy output %s"`
	Switchport            bool                 `cmd:"switchport" reg:"switchport" default:"true"`
	DhcpSnoopingThrust    bool                 `reg:"ip dhcp snooping trust" cmd:"ip dhcp snooping trust"`
	Ips                   []Ip                 `reg:"ip address.*" cmd:"ip address"`
	IPHelperAddresses     []string             `reg:"ip helper-address (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})" cmd:"ip helper-address %s"`
	Vrf                   string               `reg:"ip vrf forwarding ([[:print:]]+)" cmd:"ip vrf forwarding %s"`
	OspfNetwork           string               `reg:"ip ospf network (broadcast|non-broadcast|point-to-multipoint|point-to-point)" cmd:"ip ospf network %s"`
}

type Ip struct {
	Ip        string `reg:"ip address (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})" cmd:" %s"`
	Subnet    string `reg:"ip address (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})" cmd:" %s"`
	Secondary bool   `reg:"ip address (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})( secondary)(?: vrf ([\\w\\-]+))?" cmd:" secondary"`
	VRF       string `reg:"ip address (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})(?: secondary)?( vrf ([\\w\\-]+))" cmd:" vrf %s"`
	DHCP      bool   `reg:"ip address dhcp" cmd:" dhcp"`
}
