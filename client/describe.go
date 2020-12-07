package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/events"
)

type DescribeResourceReturn struct {
	Exists bool           `json:"exists"`
	Status string         `json:"status"`
	Events []events.Event `json:"events"`
}

func (c *Client) DescribeResource(bCtx *env.BubblyContext, rType, rName, rVersion string) (DescribeResourceReturn, error) {
	route := "describe"
	u, err := url.Parse(c.HostURL)

	if err != nil {
		return DescribeResourceReturn{}, fmt.Errorf("failed to parse bubbly server URL: %w", err)
	}

	u.Path = path.Join(u.Path, route, rType, rVersion, rName)
	requestRoute := u.String()

	bCtx.Logger.Debug().Str("route", requestRoute).Msg("attempting to describe resource via GET request")
	// describe resource
	req, err := http.NewRequest(http.MethodGet, requestRoute, nil)

	if err != nil {
		return DescribeResourceReturn{}, fmt.Errorf("failed to form GET request: %w", err)
	}

	rc, err := c.do(req)

	if err != nil {
		return DescribeResourceReturn{}, fmt.Errorf("failed to make GET request to describe resource: %w", err)
	}

	defer rc.Close()
	body, err := ioutil.ReadAll(rc)

	if err != nil {
		return DescribeResourceReturn{}, fmt.Errorf("failed to read GET response: %w", err)
	}

	// parse response body
	drrs := DescribeResourceReturn{}
	err = json.Unmarshal(body, &drrs)
	if err != nil {
		return DescribeResourceReturn{}, fmt.Errorf("failed to unmarshal describe GET response: %w", err)
	}

	return drrs, nil
}

func (c *Client) DescribeResourceGroup(bCtx *env.BubblyContext, rType, rVersion string) (map[string]DescribeResourceReturn, error) {
	route := "describe"
	u, err := url.Parse(c.HostURL)

	if err != nil {
		return nil, fmt.Errorf("failed to parse bubbly server URL: %w", err)
	}

	u.Path = path.Join(u.Path, route, rType, rVersion)
	requestRoute := u.String()

	bCtx.Logger.Debug().Str("route", requestRoute).Msg("attempting to describe resource group via GET request")

	// describe resource
	req, err := http.NewRequest(http.MethodGet, requestRoute, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to describe resource group: %w", err)
	}

	rc, err := c.do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to describe resource group: %w", err)
	}

	defer rc.Close()
	body, err := ioutil.ReadAll(rc)

	if err != nil {
		return nil, fmt.Errorf("failed to read GET response: %w", err)
	}

	// parse response body
	drrs := map[string]DescribeResourceReturn{}
	err = json.Unmarshal(body, &drrs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal describe GET response: %w", err)
	}

	return drrs, nil
}
