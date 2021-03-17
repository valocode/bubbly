package server

import (
	"github.com/valocode/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func DataStruct() core.DataBlocks {
	dataMap := core.DataBlocks{
		{
			TableName: "product",
			Fields: core.DataFields{
				"Name": cty.StringVal("1234"),
				"testField1": cty.ObjectVal(map[string]cty.Value{
					"value": cty.StringVal("TestValue1"),
				}),
				"testField2": cty.ObjectVal(map[string]cty.Value{
					"value": cty.StringVal("TestValue2"),
				}),
			},
		},
	}
	return dataMap
}
