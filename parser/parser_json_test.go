package parser

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/verifa/bubbly/api/core"
	v1 "github.com/verifa/bubbly/api/v1"
	"github.com/zclconf/go-cty/cty"
)

func TestJSON(t *testing.T) {

	p, err := NewParser("testdata/example")
	if err != nil {
		t.Errorf("Failed to create parser: %s", err.Error())
		t.FailNow()
	}
	if err := p.Decode(); err != nil {
		t.Errorf("Failed to decode parser: %s", err.Error())
		t.FailNow()
	}

	// create a new parser to load the JSON resources into
	p2 := newParser(nil, nil)
	for _, resMap := range p.Resources {
		for _, resource := range resMap {
			t.Logf("Converting resource %s to JSON", resource.String())
			bJSON, err := p2.ResourceToJSON(resource)
			if err != nil {
				t.Errorf("Failed to convert to json for resource %s: %s", resource.String(), err.Error())
			}

			println(string(bJSON))

			_, err = p2.JSONToResource(bJSON)
			if err != nil {
				t.Errorf("Failed to convert json to resource %s: %s", resource.String(), err.Error())
			}

			// if err := retRes.Decode(p2.Scope.decodeResource); err != nil {
			// 	t.Errorf("Failed to decode resource %s: %s", retRes.String(), err.Error())
			// }
		}
	}

	res := p2.Resources.Resource(core.TranslatorResourceKind, "junit")
	p2.Scope = p2.Scope.NestedScope(cty.ObjectVal(map[string]cty.Value{
		"input": cty.ObjectVal(
			map[string]cty.Value{
				"data": cty.ListVal([]cty.Value{cty.StringVal("WALALALALA")}),
			},
		),
	}))
	if err := res.Decode(p2.Scope.decodeResource); err != nil {
		t.Errorf("Failed to decode translator at the end: %s", err.Error())
	}

	for _, d := range res.(*v1.Translator).Spec.Data {
		spew.Dump(d)
	}
}
