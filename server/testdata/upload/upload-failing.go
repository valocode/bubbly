package server

import (
	"github.com/appleboy/gofight/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func DataStructFailing() map[string]interface{} {
	dataMap := gofight.D{
		"badName": "test",
		"nofields": []core.DataField{
			core.DataField{
				Name: "help",
				Value: cty.ObjectVal(map[string]cty.Value{
					"value": cty.StringVal("TestValue1"),
				}),
			},
		},
		"table": []core.Data{},
	}
	return dataMap
}
