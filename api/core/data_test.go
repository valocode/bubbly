package core

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

func TestJSONData(t *testing.T) {
	nowTime := time.Now()
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
				"time": cty.CapsuleVal(parser.TimeType, &nowTime),
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
	require.Equal(t, len(dBlocks), len(testBlocks))
	for i, b1 := range dBlocks {
		b2 := testBlocks[i]
		for name, f1 := range b1.Fields.Values {
			f2 := b2.Fields.Values[name]
			assert.Truef(t, f1.Equals(f2).True(), "fields not equal for %s.%s", b1.TableName, name)
		}
	}
}
