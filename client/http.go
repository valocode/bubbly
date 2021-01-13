package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HTTP) handleResponse(resp *http.Response, err error) (*http.Response,
	error) {
	if err != nil {
		return resp, err
	}

	// check the status code
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, fmt.Errorf(`failed to read body of respose with status "%s": %w`, resp.Status, err)
		}

		var httpError echo.HTTPError
		json.Unmarshal(body, &httpError)
		return nil, fmt.Errorf(`%s: %s`, resp.Status, httpError.Message)
	}

	return resp, nil
}
