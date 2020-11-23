package client

import (
	"net/http"

	"github.com/verifa/bubbly/config"
)

var loadDataCases = []struct {
	desc         string
	sc           config.ServerConfig
	inputFile    string
	route        string
	expected     bool
	responseCode int
	response     map[string]interface{}
}{
	{
		desc: "basic extract resource describe",
		sc: config.ServerConfig{
			Protocol: "http",
			Host:     "localhost",
			Auth:     false,
			Port:     "8080",
		},
		inputFile:    "./testdata/load/load_output.json",
		route:        "/alpha1/upload",
		expected:     true,
		responseCode: http.StatusOK,
		response: map[string]interface{}{
			"status": "uploaded",
		},
	},
}
