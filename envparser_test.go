package envparser

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

type config struct {
	AStruct struct {
		ASubStruct struct {
			AString string            `yaml:"astring" json:"astring"`
			AMap    map[string]string `yaml:"amap" json:"amap"`
		} `yaml:"asubstruct" json:"asubstruct"`
	} `yaml:"astruct" json:"astruct"`
	BStruct struct {
		BInt   int      `yaml:"bint" json:"bint"`
		BBool  bool     `yaml:"bbool" json:"bbool"`
		BSlice []string `yaml:"bslice" json:"bslice"`
	} `yaml:"bstruct" json:"bstruct"`
	CString string `yaml:"cstring" json:"cstring"`
}

var yamlC = `
astruct: 
  asubstruct: 
    amap: 
      oriKey1: oriValue1
      oriKey2: oriValue2
    astring: originstring
bstruct: 
  bbool: true
  bint: 666
  bslice: 
    - oriA
    - oriB
    - oriC
cstring: oriCstring
`

var jsonC = `
{
	"astruct": {
		"asubstruct": {
			"amap": {
				"oriKey1": "oriValue1",
				"oriKey2": "oriValue2"
			},
			"astring": "originstring"
		}
	},
	"bstruct": {
		"bbool": true,
		"bint": 666,
		"bslice": [
			"oriA",
			"oriB",
			"oriC"
		]
	},
	"cstring": "oriCString"
}
`

func TestYamlConfigOverWrite(t *testing.T) {
	yamlConf := &config{}
	err := yaml.Unmarshal([]byte(yamlC), yamlConf)
	if err != nil {
		t.Fatalf("unmarshal yamC failed: %s", err)
	}

	os.Setenv("YAMLCONFIG_ASTRUCT_ASUBSTRUCT_ASTRING", "overwrite string")
	os.Setenv("YAMLCONFIG_ASTRUCT_ASUBSTRUCT_AMAP", "{owkey1: owvalue1, owkey2: owvalue2,owkey3: owvalue3 }")
	os.Setenv("YAMLCONFIG_BSTRUCT_BINT", "888")
	os.Setenv("YAMLCONFIG_BSTRUCT_BBOOL", "false")
	os.Setenv("YAMLCONFIG_BSTRUCT_BSLICE", "['a','b','c']")

	ep := New("yamlconfig", yaml.Unmarshal)
	ep.Parse(yamlConf)
	//test string overwrite in the nested struct
	expect(t, "overwrite string", yamlConf.AStruct.ASubStruct.AString)

	//test map
	amap := yamlConf.AStruct.ASubStruct.AMap
	expect(t, 3, len(amap))
	expect(t, "owvalue1", amap["owkey1"])
	expect(t, "owvalue2", amap["owkey2"])
	expect(t, "owvalue3", amap["owkey3"])

	//test int overwrite
	expect(t, 888, yamlConf.BStruct.BInt)
	//test boolean overwrite
	expect(t, false, yamlConf.BStruct.BBool)
	//test slice overwrite
	aslice := yamlConf.BStruct.BSlice
	expect(t, 3, len(aslice))
	expect(t, "a", aslice[0])
	//test orignal string
	expect(t, "oriCstring", yamlConf.CString)
}

func TestJsonConfigOverWrite(t *testing.T) {
	jsonConf := &config{}
	err := json.Unmarshal([]byte(jsonC), jsonConf)
	if err != nil {
		t.Fatalf("unmarshal jsonC failed: %s", err)
	}

	os.Setenv("JSONCONFIG_ASTRUCT_ASUBSTRUCT_ASTRING", "\"overwrite string\"")
	os.Setenv("JSONCONFIG_ASTRUCT_ASUBSTRUCT_AMAP", "{\"owkey1\": \"owvalue1\", \"owkey2\": \"owvalue2\",\"owkey3\": \"owvalue3\"}")
	os.Setenv("JSONCONFIG_BSTRUCT_BINT", "888")
	os.Setenv("JSONCONFIG_BSTRUCT_BBOOL", "false")
	os.Setenv("JSONCONFIG_BSTRUCT_BSLICE", "[\"a\",\"b\",\"c\"]")

	ep := New("jsonconfig", json.Unmarshal)
	ep.Parse(jsonConf)
	//test string overwrite in the nested struct
	expect(t, "overwrite string", jsonConf.AStruct.ASubStruct.AString)

	//test map
	amap := jsonConf.AStruct.ASubStruct.AMap
	expect(t, 3, len(amap))
	expect(t, "owvalue1", amap["owkey1"])
	expect(t, "owvalue2", amap["owkey2"])
	expect(t, "owvalue3", amap["owkey3"])

	//test int overwrite
	expect(t, 888, jsonConf.BStruct.BInt)
	//test boolean overwrite
	expect(t, false, jsonConf.BStruct.BBool)
	//test slice overwrite
	aslice := jsonConf.BStruct.BSlice
	expect(t, 3, len(aslice))
	expect(t, "a", aslice[0])
	//test orignal string
	expect(t, "oriCString", jsonConf.CString)
}

func expect(t *testing.T, expected interface{}, actual interface{}) {
	if actual != expected {
		t.Errorf("Expected ||%#v|| (type %v) - Got ||%#v|| (type %v)", expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
	}
}
