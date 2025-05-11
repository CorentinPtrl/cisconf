package cisconf

type Config struct {
	Hostname     string           `reg:"hostname ([[:print:]]+)" cmd:"hostname %s"`
	Vlans        []Vlan           `preg:"^vlan (\\d+)$"`
	Interfaces   []CiscoInterface `preg:"(?m)^\\s*interface ([\\w\\/\\.\\-\\:]+)"`
	OSPFProcess  []Ospf           `preg:"(?m)^\\s*router ospf (\\d+)( vrf ([[:print:]]+))?"`
	EIGRPProcess []Eigrp          `preg:"(?m)^\\s*router eigrp (\\d+)"`
	EigrpNamed   []EigrpNamed     `preg:"(?m)^\\s*router eigrp ([a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*)"`
}
