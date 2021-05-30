package server

import (
	"github.com/appleboy/gofight/v2"
	"github.com/valocode/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func DataStructFailing() map[string]interface{} {
	dataMap := gofight.D{
		"badName": "test",
		"nofields": core.DataFields{Values: map[string]cty.Value{
			"help": cty.ObjectVal(map[string]cty.Value{
				"value": cty.StringVal("TestValue1"),
			})},
		},
	}
	return dataMap
}
