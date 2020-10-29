package server

import (
	"github.com/appleboy/gofight/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func DataStruct() map[string]interface{} {
	dataMap := gofight.D{"name": "test1",
		"table": []core.Data{
			{
				Name: "sub",
				Fields: []core.Field{
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
