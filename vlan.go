package cisconf

type Vlan struct {
	Id   int    `reg:"vlan (\\d+)" cmd:"vlan %d" parent:"true"`
	Name string `reg:"name ([[:print:]]+)" cmd:"name %s"`
}
