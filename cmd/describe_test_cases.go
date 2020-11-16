package cmd

import (
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/events"
)

var describeValidResourceCases = []struct {
	desc             string
	rName            string
	rType            string
	rVersion         string
	expectedContains string
}{
	{
		desc:             "basic valid resource describe",
		rType:            "extract",
		rVersion:         "v1",
		rName:            "example_extract",
		expectedContains: "failed to describe resource",
	},
}

var describeInvalidResourceCases = []struct {
	desc     string
	rName    string
	rType    string
	rVersion string
	expected string
}{
	{
		desc:     "basic invalid resource describe",
		rType:    "destroyer",
		rVersion: "v1",
		rName:    "example_destroyer",
		expected: "Error: Invalid resource type destroyer. Use 'bubbly api-resources' for a complete list of supported resources.\n",
	},
}

var describeResourceReturnCases = []struct {
	desc         string
	sc           config.ServerConfig
	route        string
	responseCode int
	response     map[string]interface{}
	token        string
	rName        string
	rType        string
	rVersion     string
	expected     string
}{
	{
		desc: "basic extract resource describe",
		sc: config.ServerConfig{
			Host: "http://localhost",
			Auth: false,
			Port: "8080",
		},
		rType:        "extract",
		rVersion:     "v1",
		rName:        "example_extract",
		route:        "/describe/extract/v1/example_extract",
		token:        "token",
		responseCode: 200,
		response: map[string]interface{}{
			"exists": true,
			"status": "creating",
			"events": []events.Event{
				{
					Status:  events.CreatingResource,
					Age:     "24h",
					Message: "Creating resource 'extract/v1/example_extract'",
				},
				{
					Status:  events.KilledResource,
					Age:     "6h",
					Message: "Killed resource 'extract/v1/example_extract'",
				},
			},
		},
		expected: "EXISTS: true, STATUS: creating, EVENTS:\n	status: Creating, age: 24h, message Creating resource 'extract/v1/example_extract'\n	status: Killed, age: 6h, message Killed resource 'extract/v1/example_extract'\n",
	},
}

var describeResourceGroupReturnCases = []struct {
	desc             string
	sc               config.ServerConfig
	route            string
	responseCode     int
	response         map[string]interface{}
	token            string
	rName            string
	rType            string
	rVersion         string
	expected         string
	expectedContains []string
}{
	{
		desc: "basic extract resource group describe",
		sc: config.ServerConfig{
			Host: "http://localhost",
			Auth: false,
			Port: "8080",
		},
		rType:        "extract",
		rVersion:     "v1",
		route:        "/describe/extract/v1",
		token:        "token",
		responseCode: 200,
		response: map[string]interface{}{
			"example_extract": map[string]interface{}{
				"exists": true,
				"status": "creating",
				"events": []events.Event{
					{
						Status:  events.CreatingResource,
						Age:     "24h",
						Message: "Creating resource 'extract/v1/example_extract'",
					},
					{
						Status:  events.KilledResource,
						Age:     "6h",
						Message: "Killed resource 'extract/v1/example_extract'",
					},
				},
			},
			"example_extract_2": map[string]interface{}{
				"exists": true,
				"status": "creating",
				"events": []events.Event{
					{
						Status:  events.CreatingResource,
						Age:     "24h",
						Message: "Creating resource 'extract/v1/example_extract_2'",
					},
					{
						Status:  events.KilledResource,
						Age:     "2h",
						Message: "Killed resource 'extract/v1/example_extract_2'",
					},
				},
			},
		},
		expectedContains: []string{"RESOURCE: example_extract", "RESOURCE: example_extract_2"},
	},
}

var describeWithVersionArgCases = []struct {
	desc     string
	rName    string
	rType    string
	rVersion string
	expected string
	version  string
}{
	{
		desc:     "basic resource describe with port argument",
		rType:    "destroyer",
		rVersion: "v1",
		rName:    "example_destroyer",
		version:  "994icidj",
		expected: "994icidj",
	},
}

var describeWithServerConfigsSetupCases = []struct {
	desc     string
	rName    string
	rType    string
	rVersion string
	expected config.Config
	port     string
	flags    map[string]string
}{
	{
		desc:     "basic resource describe with server configs pre-configured",
		rType:    "destroyer",
		rVersion: "v1",
		rName:    "example_destroyer",
		port:     "6040",
		flags: map[string]string{
			"port":  "5050",
			"host":  "localhost",
			"auth":  "false",
			"token": "",
		},
		expected: config.Config{
			ServerConfig: &config.ServerConfig{
				Protocol: "http",
				Port:     "5050",
				Host:     "localhost",
				Auth:     false,
			},
		},
	},
	{
		desc:     "basic resource describe with server configs pre-configured, auth true",
		rType:    "destroyer",
		rVersion: "v1",
		rName:    "example_destroyer",
		flags: map[string]string{
			"port":  "5050",
			"host":  "localhost",
			"auth":  "true",
			"token": "example_token",
		},
		expected: config.Config{
			ServerConfig: &config.ServerConfig{
				Protocol: "http",
				Port:     "5050",
				Host:     "localhost",
				Auth:     true,
				Token:    "example_token",
			},
		},
	},
}
