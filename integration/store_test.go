// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	testData "github.com/verifa/bubbly/integration/testdata"
	"github.com/verifa/bubbly/server"
)

func TestStore(t *testing.T) {

	s := server.GetStore()

	data, err := testData.TestAutomationData(".")
	assert.NoError(t, err)

	// Create baseline data.
	err = s.Save(data)

	assert.NoError(t, err)

	// err = s.Save(core.DataBlocks{
	// 	{
	// 		TableName: "test_run",
	// 		Fields: []core.DataField{
	// 			{
	// 				Name:  "name",
	// 				Value: cty.StringVal("run 1"),
	// 			},
	// 		},
	// 		Data: core.DataBlocks{
	// 			{
	// 				TableName: "test_set",
	// 				Fields: []core.DataField{
	// 					{
	// 						Name:  "name",
	// 						Value: cty.StringVal("set 1"),
	// 					},
	// 				},
	// 				Data: core.DataBlocks{
	// 					{
	// 						TableName: "test_case",
	// 						Fields: []core.DataField{
	// 							{
	// 								Name:  "name",
	// 								Value: cty.StringVal("case 1.1"),
	// 							},
	// 							{
	// 								Name:  "status",
	// 								Value: cty.StringVal("PASS"),
	// 							},
	// 						},
	// 					},
	// 					{
	// 						TableName: "test_case",
	// 						Fields: []core.DataField{
	// 							{
	// 								Name:  "name",
	// 								Value: cty.StringVal("case 1.2"),
	// 							},
	// 							{
	// 								Name:  "status",
	// 								Value: cty.StringVal("PASS"),
	// 							},
	// 						},
	// 					},
	// 					{
	// 						TableName: "test_case",
	// 						Fields: []core.DataField{
	// 							{
	// 								Name:  "name",
	// 								Value: cty.StringVal("case 1.3"),
	// 							},
	// 							{
	// 								Name:  "status",
	// 								Value: cty.StringVal("FAIL"),
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				TableName: "test_set",
	// 				Fields: []core.DataField{
	// 					{
	// 						Name:  "name",
	// 						Value: cty.StringVal("set 2"),
	// 					},
	// 				},
	// 				Data: core.DataBlocks{
	// 					{
	// 						TableName: "test_case",
	// 						Fields: []core.DataField{
	// 							{
	// 								Name:  "name",
	// 								Value: cty.StringVal("case 2.1"),
	// 							},
	// 							{
	// 								Name:  "status",
	// 								Value: cty.StringVal("FAIL"),
	// 							},
	// 						},
	// 					},
	// 					{
	// 						TableName: "test_case",
	// 						Fields: []core.DataField{
	// 							{
	// 								Name:  "name",
	// 								Value: cty.StringVal("case 2.2"),
	// 							},
	// 							{
	// 								Name:  "status",
	// 								Value: cty.StringVal("FAIL"),
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// })

	// assert.NoError(t, err)

	n, err := s.Query(`{
	   	test_set(name: "set 1") {
	   			name
				test_case(status: "PASS") {
					name
	   				status
	   			}
	   		}
		}`)

	assert.NoError(t, err)

	b, err := json.MarshalIndent(n, "", "\t")
	assert.NoError(t, err)

	fmt.Println(string(b))
}
