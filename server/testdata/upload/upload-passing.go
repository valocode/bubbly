package server

import (
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func DataStruct() core.DataBlocks {
	dataMap := core.DataBlocks{
		{
			TableName: "product",
			Fields: []core.DataField{
				{
					Name:  "Name",
					Value: cty.StringVal("1234"),
				},
				{
					Name: "testField1",
					Value: cty.ObjectVal(map[string]cty.Value{
						"value": cty.StringVal("TestValue1"),
					}),
				},
				{
					Name: "testField2",
					Value: cty.ObjectVal(map[string]cty.Value{
						"value": cty.StringVal("TestValue2"),
					}),
				},
			},
		},
	}
	return dataMap
}
