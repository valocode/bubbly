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

// type RequestBuilder interface {
// 	NewRequest(method, url string, body io.Reader) (*http.Request, error)
// }

// var (
// 	Request RequestBuilder
// )

// type HTTPRequestBuilder struct{}

// func (h HTTPRequestBuilder) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
// 	return http.NewRequest(method, url, body)
// }

// func init() {
// 	Request = &HTTPRequestBuilder{}
// }

func CreateRelease(bCtx *env.BubblyContext, req *api.ReleaseCreateRequest) error {
	return handlePushRequest(bCtx, req, "releases")
}

func SaveCodeScan(bCtx *env.BubblyContext, req *api.CodeScanRequest) error {
	return handlePushRequest(bCtx, req, "codescans")
}

func SaveTestRun(bCtx *env.BubblyContext, req *api.TestRunRequest) error {
	return handlePushRequest(bCtx, req, "testruns")
}

func GetAdapter(bCtx *env.BubblyContext, req *api.AdapterGetRequest) (*api.AdapterGetResponse, error) {
	url := bCtx.ClientConfig.V1() + "/adapters/" + *req.Name
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	q := httpReq.URL.Query()
	if req.Tag != nil {
		q.Add("tag", *req.Tag)
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
	if err := validator.New().Struct(a); err != nil {
		return nil, err
	}
	return &a, nil
}

func SaveAdapter(bCtx *env.BubblyContext, req *api.AdapterSaveRequest) error {
	return handlePushRequest(bCtx, req, "adapters")
}

func handlePushRequest(bCtx *env.BubblyContext, req interface{}, urlsuffix string) error {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(req); err != nil {
		return err
	}

	fmt.Printf("%s\n", b)
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
