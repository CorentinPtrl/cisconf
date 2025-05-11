package cisconf

type EigrpNamed struct {
	ProcessName  string          `reg:"router eigrp ([[:print:]]+)" cmd:"router eigrp %s" parent:"true"`
	AdressFamily []AddressFamily `preg:"(?m)^\\s*address-family\\s+ipv4\\s+unicast\\s+autonomous-system\\s+(\\d+)"`
}

type AddressFamily struct {
	Asn        int                      `reg:"address-family ipv4 unicast autonomous-system (\\d+)" cmd:"address-family ipv4 unicast autonomous-system %d" exit:"exit-address-family" parent:"true"`
	Interfaces []AddressFamilyInterface `preg:"(?m)^\\s*af-interface ([\\w\\/\\.\\-\\:]+)"`
	Network    []EigrpNetwork           `reg:"network.*" cmd:"network"`
}

type AddressFamilyInterface struct {
	InterfaceName string `reg:"af-interface ([\\w\\/\\.\\-\\:]+)" cmd:"af-interface %s" exit:"exit-af-interface" parent:"true"`
	Passive       bool   `reg:"passive-interface" cmd:"passive-interface" default:"false"`
}
