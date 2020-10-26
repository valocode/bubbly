package client

import (
	"github.com/verifa/bubbly/config"
)

var publishDataCases = []struct {
	desc         string
	sc           config.ServerConfig
	inputFile    string
	route        string
	expected     bool
	responseCode int
	response     map[string]interface{}
}{
	{
		desc: "basic importer resource describe",
		sc: config.ServerConfig{
			Protocol: "http",
			Host:     "localhost",
			Auth:     false,
			Port:     "8080",
		},
		inputFile:    "./testdata/publish/publish_output.json",
		route:        "/alpha1/upload",
		expected:     true,
		responseCode: 200,
		response: map[string]interface{}{
			"status": "uploaded",
		},
	},
}
