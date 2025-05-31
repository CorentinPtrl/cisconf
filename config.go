package cisconf

type Config struct {
	Hostname     string           `reg:"hostname ([[:print:]]+)" cmd:"hostname %s"`
	Vlans        []Vlan           `preg:"^vlan (\\d+)\\s*$"`
	Interfaces   []CiscoInterface `preg:"(?m)^\\s*interface ([\\w\\/\\.\\-\\:]+)"`
	OSPFProcess  []Ospf           `preg:"(?m)^\\s*router ospf (\\d+)( vrf ([[:print:]]+))?"`
	EIGRPProcess []Eigrp          `preg:"(?m)^\\s*router eigrp (\\d+)"`
	EigrpNamed   []EigrpNamed     `preg:"(?m)^\\s*router eigrp ([a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*)"`
	RoutesType
}

type RoutesType struct {
	Routes []Route `reg:"(?m)^\\s*ip route.*" cmd:"ip route"`
}

type Route struct {
	Prefix    string `reg:"ip route (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})(?: \\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})(?: (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}))?" cmd:" %s"`
	Mask      string `reg:"ip route (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})(?: \\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})" cmd:" %s"`
	IpAddress string `reg:"ip route (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})(?: \\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})" cmd:" %s"`
}
