package cmd

import (
	"net/http"

	"github.com/verifa/bubbly/config"
)

var applyWithServerConfigsSetupCases = []struct {
	desc         string
	flags        map[string]string
	serverConfig *config.ServerConfig
	route        string
	address      string
	responseCode int
	response     map[string]string
	expected     bool
}{
	{
		desc: "basic resource apply with server configs pre-configured",
		serverConfig: &config.ServerConfig{
			Protocol: "http",
			Port:     "8070",
			Host:     "localhost",
			Auth:     false,
			Token:    "",
		},
		expected:     true,
		address:      "http://localhost:8070",
		route:        "/alpha1/upload",
		responseCode: http.StatusOK,
		response: map[string]string{
			"status": "uploaded",
		},
	},
	{
		desc: "basic resource apply with server configs pre-configured incorrectly (gock runs at different address)",
		serverConfig: &config.ServerConfig{
			Protocol: "http",
			Port:     "8060",
			Host:     "localhost",
			Auth:     false,
			Token:    "",
		},
		expected:     true,
		address:      "http://localhost:8070",
		route:        "/alpha1/upload",
		responseCode: http.StatusOK,
		response: map[string]string{
			"status": "uploaded",
		},
	},
}
