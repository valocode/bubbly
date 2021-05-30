package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

func TestJSONData(t *testing.T) {
	dBlocks := DataBlocks{
		Data{
			TableName: "TestTable",
			Fields: &DataFields{Values: map[string]cty.Value{
				"TestField": cty.ObjectVal(
					map[string]cty.Value{
						"attribute": cty.ObjectVal(
							map[string]cty.Value{
								"nested": cty.StringVal("some value"),
							},
						),
					},
				),
				"string_field": cty.StringVal("mystring"),
				"DataRef": cty.CapsuleVal(parser.DataRefType, &parser.DataRef{
					TableName: "my_table", Field: "my_field",
				}),
			}},
			Joins: []string{"TestJoin", "DataRef"},
		},
	}

	jBytes, err := json.Marshal(dBlocks)
	if err != nil {
		t.Errorf("Failed to marshal JSON: %s", err.Error())
		t.FailNow()
	}

	t.Logf("Transform JSON() returned: %s", string(jBytes))

	testBlocks := DataBlocks{}
	err = json.Unmarshal(jBytes, &testBlocks)
	require.NoErrorf(t, err, "failed to unmarshal json data blocks")
	require.Equalf(t, dBlocks, testBlocks, "JSON returned from transform does not equal unmarshalled dataBlocks")
}
