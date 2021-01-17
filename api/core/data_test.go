package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

func TestJSONData(t *testing.T) {
	// TODO: make some proper test cases
	// tcs := []struct {
	// 	desc     string
	// 	dblocks  DataBlocks
	// 	expected DataBlocks
	// }{}
	dBlocks := DataBlocks{
		Data{
			TableName: "TestTable",
			Fields: DataFields{
				DataField{
					Name: "TestField",
					Value: cty.ObjectVal(
						map[string]cty.Value{
							"attribute": cty.ObjectVal(
								map[string]cty.Value{
									"nested": cty.StringVal("some value"),
								},
							),
						},
					),
				},
				DataField{
					Name: "DataRef",
					Value: cty.CapsuleVal(parser.DataRefType, &parser.DataRef{
						TableName: "my_table", Field: "my_field",
					}),
				},
			},
			Joins: DataJoins{
				DataJoin{
					Table: "TestJoin",
					Value: cty.ObjectVal(
						map[string]cty.Value{
							"attribute": cty.ObjectVal(
								map[string]cty.Value{
									"nested": cty.StringVal("some value"),
								},
							),
						},
					),
				},
				DataJoin{
					Table: "DataRef",
					Value: cty.CapsuleVal(parser.DataRefType, &parser.DataRef{
						TableName: "my_table", Field: "my_field",
					}),
				},
			},
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
	assert.NoErrorf(t, err, "failed to unmarshal json data blocks")
	assert.Equalf(t, dBlocks, testBlocks, "JSON returned from transform does not equal unmarshalled dataBlocks")
}
