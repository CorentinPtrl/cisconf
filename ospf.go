package cisconf

type OspfParent struct {
	Asn        int    `reg:"router ospf (\\d+)" cmd:" %d"`
	ProcessVRF string `reg:"router ospf \\d+ vrf ([[:print:]]+)" cmd:" vrf %s"`
}
type Ospf struct {
	Parent                  OspfParent `reg:"router ospf.*" cmd:"router ospf" parent:"true"`
	LogAdjacencyChange      bool       `reg:"log-adjacency-changes detail" cmd:"log-adjacency-changes detail"`
	PassiveInterfaceDefault bool       `reg:"passive-interface default" cmd:"passive-interface default"`
	PassiveInterface        []int
	Network                 []OspfNetwork `reg:"network.*" cmd:"network"`
}

type OspfNetwork struct {
	NetworkNumber string `reg:"network (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) area (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}|\\d+)" cmd:" %s"`
	WildCard      string `reg:"network (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) area (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}|\\d+)" cmd:" %s"`
	Area          string `reg:"network (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) area (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}|\\d+)" cmd:" area %s"`
}
