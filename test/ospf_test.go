package test

import (
	"github.com/CorentinPtrl/cisconf"
	"io/ioutil"
	"testing"
)

type testOspf struct {
	FileName string
}

func TestOspf_Parse(t *testing.T) {
	var ospfFiles []testOspf
	process := testOspf{
		FileName: "process.txt",
	}
	ospfFiles = append(ospfFiles, process)
	processVrf := testOspf{
		FileName: "processVrf.txt",
	}
	ospfFiles = append(ospfFiles, processVrf)

	for _, ospfFile := range ospfFiles {
		// Open File
		content, err := ioutil.ReadFile("res/ospf/" + ospfFile.FileName)
		if err != nil {
			t.Error(err)
		}

		var ospf cisconf.Ospf
		err = cisconf.Unmarshal(string(content), &ospf)
		if err != nil {
			t.Error(err)
		}

		generated, err := cisconf.Marshal(ospf)
		if err != nil {
			t.Error(err)
		}
		contentStr := string(content)
		if generated != contentStr {
			t.Error("Config wrong parsed or generated \n File: \n", string(content), "\n Generated: \n", generated)
		}
	}
}
