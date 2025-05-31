package cisconf

type Eigrp struct {
	Asn     int            `reg:"router eigrp (\\d+)" cmd:"router eigrp %d" parent:"true"`
	Network []EigrpNetwork `reg:"network.*" cmd:"network"`
}

type EigrpNetwork struct {
	NetworkNumber string `reg:"network (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})(?: (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}))?" cmd:" %s"`
	WildCard      string `reg:"network (?:\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) (\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})" cmd:" %s" default:"0.0.0.255"`
}
