package client

import (
	"net/http"

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
		desc: "basic extract resource describe",
		sc: config.ServerConfig{
			Protocol: "http",
			Host:     "localhost",
			Auth:     false,
			Port:     "8080",
		},
		rType:        "extract",
		rVersion:     "v1",
		rName:        "example_extract",
		route:        "/describe/extract/v1/example_extract",
		token:        "token",
		responseCode: http.StatusOK,
		response: map[string]interface{}{
			"exists": true,
			"status": "creating",
			"events": []events.Event{
				{
					Status:  events.CreatingResource,
					Age:     "24h",
					Message: "Creating resource 'extract/v1/example_extract'",
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
					Message: "Creating resource 'extract/v1/example_extract'",
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
		desc: "basic extract resource group describe",
		sc: config.ServerConfig{
			Protocol: "http",
			Host:     "localhost",
			Auth:     false,
			Port:     "8080",
		},
		rType:        "extract",
		rVersion:     "v1",
		route:        "/describe/extract/v1",
		token:        "token",
		responseCode: http.StatusOK,
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
		expected: map[string]DescribeResourceReturn{
			"example_extract": {
				Exists: true,
				Status: "creating",
				Events: []events.Event{
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
			"example_extract_2": {
				Exists: true,
				Status: "creating",
				Events: []events.Event{
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
	},
}
