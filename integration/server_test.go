// +build integration

package integration

import (
	"testing"

	"github.com/verifa/bubbly/bubbly"
	"github.com/verifa/bubbly/config"
)

func TestServer(t *testing.T) {

	// data, err := testData.TestAutomationData(".")
	// assert.NoError(t, err)

	// url := "http://" + hostURL
	// client, err := client.NewUnauthClient(&url)
	// assert.NoError(t, err)

	// err = client.Load(data)
	// assert.NoError(t, err)

	serverConfig := config.ServerConfig{
		Protocol: "http",
		Port:     "8112",
		Host:     "localhost",
		Auth:     false,
	}
	bubbly.Apply("./testdata/testautomation/golang/pipeline.bubbly", serverConfig)

	// TODO this query should be sent by the client over HTTP, not directly
	// to the store
	// s := server.GetStore()
	// n, err := s.Query(`{
	//    	test_set(name: "set 1") {
	//    			name
	// 			test_case(status: "PASS") {
	// 				name
	//    				status
	//    			}
	//    		}
	// 	}`)

	// assert.NoError(t, err)

	// b, err := json.MarshalIndent(n, "", "\t")
	// assert.NoError(t, err)

	// fmt.Println(string(b))
}
