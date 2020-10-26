package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/verifa/bubbly/events"
)

type DescribeResourceReturn struct {
	Exists bool           `json:"exists"`
	Status string         `json:"status"`
	Events []events.Event `json:"events"`
}

func (c *Client) DescribeResource(rType, rName, rVersion string) (DescribeResourceReturn, error) {
	fmt.Printf("Making request to Bubbly server from client.DescribeResource to describe resource %s of type %s and API version %s\n", rType, rName, rVersion)
	// describe resource
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/describe/%s/%s/%s", c.HostURL, rType, rVersion, rName), nil)

	if err != nil {
		return DescribeResourceReturn{}, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return DescribeResourceReturn{}, err
	}

	// parse response body
	drrs := DescribeResourceReturn{}
	err = json.Unmarshal(body, &drrs)
	if err != nil {
		return DescribeResourceReturn{}, err
	}

	return drrs, nil
}

func (c *Client) DescribeResourceGroup(rType, rVersion string) (map[string]DescribeResourceReturn, error) {
	fmt.Printf("Making request to Bubbly server from client.DescribeResourceGroup to describe all resources of type %s and API version %s\n", rType, rVersion)

	// describe resource
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/describe/%s/%s", c.HostURL, rType, rVersion), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return nil, err
	}

	// parse response body
	drrs := map[string]DescribeResourceReturn{}
	err = json.Unmarshal(body, &drrs)
	if err != nil {
		return nil, err
	}

	return drrs, nil
}
