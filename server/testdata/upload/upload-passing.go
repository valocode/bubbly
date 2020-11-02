package server

import (
	"github.com/appleboy/gofight/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func DataStruct() map[string]interface{} {
	dataMap := gofight.D{
		"table": "test1",
		"data": []core.Data{
			{
				TableName: "sub",
				Fields: []core.DataField{
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
		},
	}
	return dataMap
}
