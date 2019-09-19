package config_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/setlog/trivrost/pkg/launcher/config"
)

const testJson1 string = `{
	"foo": { "SHA256": "9E079B502D173FE926B04E87715F4534C34F23EDF8E91FBBB2510BC666FB6C76" },
	"bar": { "SHA256": "3CD33E6295AEEB0622990B149A78A648200DE73C5CC0BF573F57CFDC0A6F0074" },
	"bee": { "SHA256": "56175A1FF29A145F58FFEC4ACC361D69728FEAEB35447750B271A8B55F4FFF50" }
}`

const testJson2 string = `{
	"foo": { "SHA256": "84C95435C2EC37380F38E453B05782B9D177D745C4BB7B43005D0337025DA857" },
	"bar": { "SHA256": "3CD33E6295AEEB0622990B149A78A648200DE73C5CC0BF573F57CFDC0A6F0074" },
	"baz": { "SHA256": "CE2213DACD9C4A703AC3A960D006F06468AAE2EF9E2D09EE87872F55A0CBECE6" }
}`

func TestUnmarshal(t *testing.T) {
	fm := config.NewFileInfoMap()
	err := json.Unmarshal([]byte(testJson1), &fm)
	if err != nil {
		t.Errorf("Unmarshal failed: %v", err)
	}
	if !strings.EqualFold(fm["foo"].SHA256, "9E079B502D173FE926B04E87715F4534C34F23EDF8E91FBBB2510BC666FB6C76") {
		t.Error("Missing 9E079B502D173FE926B04E87715F4534C34F23EDF8E91FBBB2510BC666FB6C76")
	}
	if !strings.EqualFold(fm["bar"].SHA256, "3CD33E6295AEEB0622990B149A78A648200DE73C5CC0BF573F57CFDC0A6F0074") {
		t.Error("Missing 3CD33E6295AEEB0622990B149A78A648200DE73C5CC0BF573F57CFDC0A6F0074")
	}
	if !strings.EqualFold(fm["bee"].SHA256, "56175A1FF29A145F58FFEC4ACC361D69728FEAEB35447750B271A8B55F4FFF50") {
		t.Error("Missing 56175A1FF29A145F58FFEC4ACC361D69728FEAEB35447750B271A8B55F4FFF50")
	}
	if len(fm) != 3 {
		t.Error("Length not 3")
	}
}

func TestMakeUpdateMap(t *testing.T) {
	fm1 := make(config.FileInfoMap)
	fm2 := make(config.FileInfoMap)
	err := json.Unmarshal([]byte(testJson1), &fm1)
	if err != nil {
		t.Errorf("Could not unmarshal json: %v", err)
		t.FailNow()
	}
	err = json.Unmarshal([]byte(testJson2), &fm2)
	if err != nil {
		t.Errorf("Could not unmarshal json: %v", err)
		t.FailNow()
	}
	um := config.MakeDiffFileInfoMap(fm1, fm2)
	upFoo, upBee, upBaz := um["foo"], um["bee"], um["baz"]
	upBar, hasBar := um["bar"]
	if upFoo.SHA256 == "" {
		t.Error("Deleting outdated foo")
	}
	if upBee.SHA256 != "" {
		t.Error("Updating removed bee")
	}
	if upBaz.SHA256 == "" {
		t.Error("Deleting new baz")
	}
	if hasBar {
		if upBar.SHA256 == "" {
			t.Error("Incorrectly deleting unchanged bar")
		} else {
			t.Error("Redundantly updating unchanged bar")
		}
	}
}
