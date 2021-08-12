package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"
)

func CreateRelease(bCtx *env.BubblyContext, req *api.ReleaseCreateRequest) error {
	return handlePushRequest(bCtx, req, "releases")
}

func SaveCodeScan(bCtx *env.BubblyContext, req *api.CodeScanRequest) error {
	return handlePushRequest(bCtx, req, "codescans")
}

func GetAdapter(bCtx *env.BubblyContext, req *api.AdapterGetRequest) (*adapter.Adapter, error) {
	url := "http://localhost:8111/api/v1/adapters/" + *req.Name
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	q := httpReq.URL.Query()
	if req.Tag != nil {
		q.Add("tag", *req.Tag)
	}
	if req.Type != nil {
		q.Add("type", *req.Type)
	}
	httpReq.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if err := handleResponseError(resp); err != nil {
		return nil, err
	}
	var a api.AdapterGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, fmt.Errorf("decoding adapter response: %w", err)
	}
	return adapter.FromModel(a.AdapterModel)
}

func SaveAdapter(bCtx *env.BubblyContext, req *api.AdapterSaveRequest) error {
	return handlePushRequest(bCtx, req, "adapters")
}

func handlePushRequest(bCtx *env.BubblyContext, req interface{}, urlsuffix string) error {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(req); err != nil {
		return err
	}
	url := "http://localhost:8111/api/v1/" + urlsuffix
	resp, err := http.Post(url, echo.MIMEApplicationJSON, b)
	if err != nil {
		return err
	}
	if err := handleResponseError(resp); err != nil {
		return err
	}
	return nil
}

func handleResponseError(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	var httpErr echo.HTTPError
	if err := json.NewDecoder(resp.Body).Decode(&httpErr); err != nil {
		return fmt.Errorf("decoding HTTP error from the server: %w", err)
	}
	return fmt.Errorf("HTTP Status: %d, Message: %s", resp.StatusCode, httpErr.Message)
}
