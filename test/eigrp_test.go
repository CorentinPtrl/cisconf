package test

import (
	"github.com/CorentinPtrl/cisconf"
	"io/ioutil"
	"testing"
)

func TestEigrp_Parse(t *testing.T) {
	var ospfFiles []testOspf
	process := testOspf{
		FileName: "process.txt",
	}
	ospfFiles = append(ospfFiles, process)
	for _, ospfFile := range ospfFiles {
		content, err := ioutil.ReadFile("res/eigrp/" + ospfFile.FileName)
		if err != nil {
			t.Error(err)
		}

		var eigrp cisconf.Eigrp
		err = cisconf.Unmarshal(string(content), &eigrp)
		if err != nil {
			t.Error(err)
		}

		generated, err := cisconf.Marshal(eigrp)
		if err != nil {
			t.Error(err)
		}
		if generated != string(content) {
			t.Error("Config wrong parsed or generated \n File: \n", string(content), "\n Generated: \n", generated)
		}
	}
}
