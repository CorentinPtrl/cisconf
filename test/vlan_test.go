package test

import (
	"github.com/CorentinPtrl/cisconf"
	"testing"
)

func TestParseVlan(t *testing.T) {
	test := "vlan 300\n name office\n!"
	var vlan cisconf.Vlan
	err := cisconf.Unmarshal(test, &vlan)
	if err != nil {
		t.Error(err)
	}

	if vlan.Id != 300 {
		t.Errorf("Wrong VLan Id wants: `300` got: `%v`", vlan.Id)
	}

	if vlan.Name != "office" {
		t.Errorf("Wrong VLan name wants: `office` got: `%v`", vlan.Id)
	}
}
