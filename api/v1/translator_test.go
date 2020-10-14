package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestTranslator(t *testing.T) {
	dBlocks := DataBlocks{
		&Data{
			TableName: "TestTable",
			RowName:   "TestRow",
			Fields: Fields{
				&Field{
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
	translator := Translator{
		Spec: TranslatorSpec{
			Data: dBlocks,
		},
	}

	value, err := translator.JSON()
	if err != nil {
		t.Errorf("Failed to Resolve() JSON importer: %s", err.Error())
		t.Fail()
	}

	t.Logf("Translator JSON() returned: %s", value)

	testBlocks := DataBlocks{}
	err = json.Unmarshal(value, &testBlocks)
	if err != nil {
		t.Errorf("Failed to Unmarshal JSON DataBlocks: %s", err.Error())
	}
	assert.Equalf(t, dBlocks, testBlocks, "JSON returned from translator equals unmarshalled DataBlocks")
}
