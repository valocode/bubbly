package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"
)

func CreateRelease(bCtx *env.BubblyConfig, req *api.ReleaseCreateRequest) error {
	return handlePostRequest(bCtx, req, "releases")
}

func GetRelease(bCtx *env.BubblyConfig, req *api.ReleaseGetRequest) (*api.ReleaseGetResponse, error) {
	var r api.ReleaseGetResponse
	params, err := structToStringMap(req)
	if err != nil {
		return nil, err
	}
	if err := handleGetRequest(bCtx, &r, "releases", params); err != nil {
		return nil, err
	}
	return &r, nil
}

func SaveCodeScan(bCtx *env.BubblyConfig, req *api.CodeScanRequest) error {
	return handlePostRequest(bCtx, req, "codescans")
}

func SaveTestRun(bCtx *env.BubblyConfig, req *api.TestRunRequest) error {
	return handlePostRequest(bCtx, req, "testruns")
}

func GetAdapter(bCtx *env.BubblyConfig, req *api.AdapterGetRequest) (*api.AdapterGetResponse, error) {
	var a api.AdapterGetResponse
	params, err := structToStringMap(req)
	if err != nil {
		return nil, err
	}

	if err := handleGetRequest(bCtx, &a, "adapters/"+*req.Name, params); err != nil {
		return nil, err
	}
	return &a, nil
}

func SaveAdapter(bCtx *env.BubblyConfig, req *api.AdapterSaveRequest) error {
	return handlePostRequest(bCtx, req, "adapters")
}

func GetPolicy(bCtx *env.BubblyConfig, req *api.ReleasePolicyGetRequest) (*api.ReleasePolicyGetResponse, error) {
	var r api.ReleasePolicyGetResponse
	if err := handleGetRequest(bCtx, &r, "policies/"+*req.Name, nil); err != nil {
		return nil, err
	}
	return &r, nil
}

func SavePolicy(bCtx *env.BubblyConfig, req *api.ReleasePolicySaveRequest) error {
	return handlePostRequest(bCtx, req, "policies")
}

func SetPolicy(bCtx *env.BubblyConfig, req *api.ReleasePolicySetRequest) error {
	return handlePutRequest(bCtx, req, "policies")
}

func handleGetRequest(bCtx *env.BubblyConfig, resp interface{}, urlsuffix string, params map[string]string) error {
	url := bCtx.ClientConfig.V1() + "/" + urlsuffix
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	q := httpReq.URL.Query()
	for name, value := range params {
		q.Add(name, value)
	}
	httpReq.URL.RawQuery = q.Encode()

	bCtx.Logger.Debug().
		Str("query", q.Encode()).
		Str("method", http.MethodGet).
		Msg("Client making request")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	if err := handleResponseError(httpResp); err != nil {
		return err
	}
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return fmt.Errorf("decoding adapter response: %w", err)
	}
	if err := validator.New().Struct(resp); err != nil {
		return err
	}
	return nil
}

func handlePostRequest(bCtx *env.BubblyConfig, req interface{}, urlsuffix string) error {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(req); err != nil {
		return err
	}

	bCtx.Logger.Debug().
		Bytes("data", b.Bytes()).
		Str("endpoint", urlsuffix).
		Str("method", http.MethodPost).
		Msg("Client making request")

	url := bCtx.ClientConfig.V1() + "/" + urlsuffix
	httpReq, err := http.NewRequest(http.MethodPost, url, b)
	if err != nil {
		return err
	}
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	if err := handleResponseError(resp); err != nil {
		return err
	}
	return nil
}

func handlePutRequest(bCtx *env.BubblyConfig, req interface{}, urlsuffix string) error {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(req); err != nil {
		return err
	}

	bCtx.Logger.Debug().
		Bytes("data", b.Bytes()).
		Str("endpoint", urlsuffix).
		Str("method", http.MethodPut).
		Msg("Client making request")

	url := bCtx.ClientConfig.V1() + "/" + urlsuffix
	httpReq, err := http.NewRequest(http.MethodPut, url, b)
	if err != nil {
		return err
	}
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp, err := http.DefaultClient.Do(httpReq)
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

// TODO: this won't work because we cannot distinguish between path params and query params...
// Solution is to use mapstructure and add the skip tag for those fields which should
// be path params and not query params. See GetEvents example.
func structToStringMap(req interface{}) (map[string]string, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encoding request to JSON: %w", err)
	}
	var params map[string]string
	if err := json.Unmarshal(b, &params); err != nil {
		return nil, fmt.Errorf("decoding request to string map: %w", err)
	}
	return params, nil
}
