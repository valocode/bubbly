package v1

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestImporterJSON(t *testing.T) {
	ctyType := cty.Object(map[string]cty.Type{
		"issues": cty.List(cty.Object(map[string]cty.Type{
			"engineId": cty.String,
			"ruleId":   cty.String,
			"severity": cty.String,
			"type":     cty.String,
			"primaryLocation": cty.Object(map[string]cty.Type{
				"message":  cty.String,
				"filePath": cty.String,
				"textRange": cty.Object(map[string]cty.Type{
					"startLine":   cty.Number,
					"endLine":     cty.Number,
					"startColumn": cty.Number,
					"endColumn":   cty.Number,
				}),
			}),
		})),
	})

	importer := Importer{
		Spec: ImporterSpec{
			Type:   jsonImporterType,
			Format: ctyType,
			Source: &JSONSource{
				File: "testdata/importer/json/sonarqube-example.json",
			},
		},
	}

	val := importer.Output()
	if val.Error != nil {
		t.Errorf("Failed to Resolve() JSON importer: %s", val.Error.Error())
		t.Fail()
	}
	if val.Value.IsNull() {
		t.Errorf("Received Null type value")
	}
	t.Logf("JSON Importer returned value: %s", val.Value.GoString())
}

func TestImporterXML(t *testing.T) {
	ctyType := cty.Object(map[string]cty.Type{
		"testsuites": cty.Object(map[string]cty.Type{
			"duration": cty.Number,
			"testsuite": cty.List(cty.Object(map[string]cty.Type{
				"failures": cty.Number,
				"name":     cty.String,
				"package":  cty.String,
				"tests":    cty.Number,
				"time":     cty.Number,
				"testcase": cty.List(cty.Object(map[string]cty.Type{
					"classname": cty.String,
					"name":      cty.String,
					"time":      cty.Number,
				})),
			})),
		}),
	})
	importer := Importer{
		Spec: ImporterSpec{
			Type:   xmlImporterType,
			Format: ctyType,
			Source: &XMLSource{
				File: "testdata/importer/json/sonarqube-example.json",
			},
		},
	}

	val := importer.Output()
	if val.Error != nil {
		t.Errorf("Failed to Resolve() XML importer: %s", val.Error.Error())
		t.Fail()
	}
	if val.Value.IsNull() {
		t.Errorf("Received Null type value")
	}
	t.Logf("XML Importer returned value: %s", val.Value.GoString())
}
