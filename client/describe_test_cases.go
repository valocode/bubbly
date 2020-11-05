package client

import (
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/events"
)

var describeResourceCases = []struct {
	desc         string
	sc           config.ServerConfig
	route        string
	responseCode int
	response     map[string]interface{}
	token        string
	rName        string
	rType        string
	rVersion     string
	expected     DescribeResourceReturn
}{
	{
		desc: "basic importer resource describe",
		sc: config.ServerConfig{
			Protocol: "http",
			Host:     "localhost",
			Auth:     false,
			Port:     "8080",
		},
		rType:        "importer",
		rVersion:     "v1",
		rName:        "example_importer",
		route:        "/describe/importer/v1/example_importer",
		token:        "token",
		responseCode: 200,
		response: map[string]interface{}{
			"exists": true,
			"status": "creating",
			"events": []events.Event{
				{
					Status:  events.CreatingResource,
					Age:     "24h",
					Message: "Creating resource 'importer/v1/example_importer'",
				},
			},
		},
		expected: DescribeResourceReturn{
			Exists: true,
			Status: "creating",
			Events: []events.Event{
				{
					Status:  events.CreatingResource,
					Age:     "24h",
					Message: "Creating resource 'importer/v1/example_importer'",
				},
			},
		},
	},
}

var describeResourceGroupCases = []struct {
	desc         string
	sc           config.ServerConfig
	route        string
	responseCode int
	response     map[string]interface{}
	token        string
	rType        string
	rVersion     string
	expected     map[string]DescribeResourceReturn
}{
	{
		desc: "basic importer resource group describe",
		sc: config.ServerConfig{
			Protocol: "http",
			Host:     "localhost",
			Auth:     false,
			Port:     "8080",
		},
		rType:        "importer",
		rVersion:     "v1",
		route:        "/describe/importer/v1",
		token:        "token",
		responseCode: 200,
		response: map[string]interface{}{
			"example_importer": map[string]interface{}{
				"exists": true,
				"status": "creating",
				"events": []events.Event{
					{
						Status:  events.CreatingResource,
						Age:     "24h",
						Message: "Creating resource 'importer/v1/example_importer'",
					},
					{
						Status:  events.KilledResource,
						Age:     "6h",
						Message: "Killed resource 'importer/v1/example_importer'",
					},
				},
			},
			"example_importer_2": map[string]interface{}{
				"exists": true,
				"status": "creating",
				"events": []events.Event{
					{
						Status:  events.CreatingResource,
						Age:     "24h",
						Message: "Creating resource 'importer/v1/example_importer_2'",
					},
					{
						Status:  events.KilledResource,
						Age:     "2h",
						Message: "Killed resource 'importer/v1/example_importer_2'",
					},
				},
			},
		},
		expected: map[string]DescribeResourceReturn{
			"example_importer": {
				Exists: true,
				Status: "creating",
				Events: []events.Event{
					{
						Status:  events.CreatingResource,
						Age:     "24h",
						Message: "Creating resource 'importer/v1/example_importer'",
					},
					{
						Status:  events.KilledResource,
						Age:     "6h",
						Message: "Killed resource 'importer/v1/example_importer'",
					},
				},
			},
			"example_importer_2": {
				Exists: true,
				Status: "creating",
				Events: []events.Event{
					{
						Status:  events.CreatingResource,
						Age:     "24h",
						Message: "Creating resource 'importer/v1/example_importer_2'",
					},
					{
						Status:  events.KilledResource,
						Age:     "2h",
						Message: "Killed resource 'importer/v1/example_importer_2'",
					},
				},
			},
		},
	},
}
