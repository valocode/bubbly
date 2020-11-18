package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestJSONData(t *testing.T) {
	dBlocks := DataBlocks{
		Data{
			TableName: "TestTable",
			// RowName:   "TestRow",
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
	if err != nil {
		t.Errorf("Failed to Unmarshal JSON dataBlocks: %s", err.Error())
	}
	assert.Equalf(t, dBlocks, testBlocks, "JSON returned from transform equals unmarshalled dataBlocks")
}
