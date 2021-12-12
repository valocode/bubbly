package client

import (
	"fmt"
	"net/http"

	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/store/api"
)

func CreateRelease(bCtx *config.BubblyConfig, req *api.ReleaseCreateRequest) error {
	return handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodPost), WithPayload(req), WithRequestURL("releases"),
	)
}

func GetReleases(bCtx *config.BubblyConfig, req *api.ReleaseGetRequest) (*api.ReleaseGetResponse, error) {
	var r api.ReleaseGetResponse
	if err := handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodGet),
		WithResponse(&r), WithRequestURL("releases"), WithQueryParamsStruct(req),
	); err != nil {
		return nil, err
	}
	return &r, nil
}

func SaveCodeScan(bCtx *config.BubblyConfig, req *api.CodeScanRequest) error {
	return handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("codescans"),
	)
}

func SaveTestRun(bCtx *config.BubblyConfig, req *api.TestRunRequest) error {
	return handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("testruns"),
	)
}

func GetAdapters(bCtx *config.BubblyConfig, req *api.AdapterGetRequest) (*api.AdapterGetResponse, error) {
	var a api.AdapterGetResponse
	if err := handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodGet),
		WithRequestURL("adapters"), WithQueryParamsStruct(req), WithResponse(&a),
	); err != nil {
		return nil, err
	}
	return &a, nil
}

func SaveAdapter(bCtx *config.BubblyConfig, req *api.AdapterSaveRequest) error {
	return handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("adapters"),
	)
}

func GetPolicies(bCtx *config.BubblyConfig, req *api.ReleasePolicyGetRequest) (*api.ReleasePolicyGetResponse, error) {
	var r api.ReleasePolicyGetResponse
	if err := handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodGet),
		WithRequestURL("policies"), WithQueryParamsStruct(req), WithResponse(&r),
	); err != nil {
		return nil, err
	}
	return &r, nil
}

func SavePolicy(bCtx *config.BubblyConfig, req *api.ReleasePolicyCreateRequest) error {
	return handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("policies"),
	)
}

func SetPolicy(bCtx *config.BubblyConfig, req *api.ReleasePolicyUpdateRequest) error {
	return handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodPut),
		WithPayload(req), WithRequestURL(fmt.Sprintf("policies/%d", *req.ID)),
	)
}

func GetVulnerabilityReviews(bCtx *config.BubblyConfig, req *api.VulnerabilityReviewGetRequest) (*api.VulnerabilityReviewGetResponse, error) {
	var r api.VulnerabilityReviewGetResponse
	if err := handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodGet),
		WithRequestURL("vulnerabilityreviews"), WithQueryParamsStruct(req), WithResponse(&r),
	); err != nil {
		return nil, err
	}
	return &r, nil
}

func SaveVulnerabilityReview(bCtx *config.BubblyConfig, req *api.VulnerabilityReviewSaveRequest) error {
	return handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("vulnerabilityreviews"),
	)
}

func UpdateVulnerabilityReview(bCtx *config.BubblyConfig, req *api.VulnerabilityReviewUpdateRequest) error {
	return handleRequest(
		WithBubblyConfig(bCtx), WithAPIV1(true), WithMethod(http.MethodPut),
		WithPayload(req), WithRequestURL(fmt.Sprintf("vulnerabilityreviews/%d", *req.ID)),
	)
}
