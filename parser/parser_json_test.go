package parser

import (
	"testing"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func TestJSON(t *testing.T) {

	p, err := NewParserFromFilename("testdata")
	if err != nil {
		t.Errorf("Failed to create parser: %s", err.Error())
		t.FailNow()
	}
	if err := p.Parse(); err != nil {
		t.Errorf("Failed to decode parser: %s", err.Error())
		t.FailNow()
	}

	// create a new parser to load the JSON resources into
	p2 := newParser(nil, nil)
	for _, resMap := range p.Resources {
		for _, resource := range resMap {
			t.Logf("Converting resource %s to JSON", resource.String())
			bJSON, err := resource.JSON(p.Context(cty.NilVal))
			if err != nil {
				t.Errorf("Failed to convert to json for resource %s: %s", resource.String(), err.Error())
			}

			println(string(bJSON))

			_, err = p2.JSONToResource(bJSON)
			if err != nil {
				t.Errorf("Failed to convert json to resource %s: %s", resource.String(), err.Error())
			}
		}
	}

	res, err := p2.GetResource(core.TranslatorResourceKind, "junit")
	if err != nil {
		t.Errorf("Couldnt get resource: %s", "junit")
		t.FailNow()
	}
	inputs := cty.ObjectVal(map[string]cty.Value{
		"input": cty.ObjectVal(
			map[string]cty.Value{
				"data": cty.ListVal([]cty.Value{cty.StringVal("WALALALALA")}),
			},
		),
	})
	if out := res.Apply(p2.Context(inputs)); out.Error != nil {
		t.Errorf("Failed to decode translator at the end: %s", out.Error.Error())
	}
}
