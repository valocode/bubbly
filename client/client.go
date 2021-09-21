package client

import (
	"fmt"
	"net/http"

	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"
)

func CreateRelease(bCtx *env.BubblyContext, req *api.ReleaseCreateRequest) error {
	return handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodPost), WithPayload(req), WithRequestURL("releases"),
	)
}

func GetReleases(bCtx *env.BubblyContext, req *api.ReleaseGetRequest) (*api.ReleaseGetResponse, error) {
	var r api.ReleaseGetResponse
	if err := handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodGet),
		WithResponse(&r), WithRequestURL("releases"), WithQueryParamsStruct(req),
	); err != nil {
		return nil, err
	}
	return &r, nil
}

func SaveCodeScan(bCtx *env.BubblyContext, req *api.CodeScanRequest) error {
	return handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("codescans"),
	)
}

func SaveTestRun(bCtx *env.BubblyContext, req *api.TestRunRequest) error {
	return handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("testruns"),
	)
}

func GetAdapters(bCtx *env.BubblyContext, req *api.AdapterGetRequest) (*api.AdapterGetResponse, error) {
	var a api.AdapterGetResponse
	if err := handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodGet),
		WithRequestURL("adapters"), WithQueryParamsStruct(req), WithResponse(&a),
	); err != nil {
		return nil, err
	}
	return &a, nil
}

func SaveAdapter(bCtx *env.BubblyContext, req *api.AdapterSaveRequest) error {
	return handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("adapters"),
	)
}

func GetPolicies(bCtx *env.BubblyContext, req *api.ReleasePolicyGetRequest) (*api.ReleasePolicyGetResponse, error) {
	var r api.ReleasePolicyGetResponse
	if err := handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodGet),
		WithRequestURL("policies"), WithQueryParamsStruct(req), WithResponse(&r),
	); err != nil {
		return nil, err
	}
	return &r, nil
}

func SavePolicy(bCtx *env.BubblyContext, req *api.ReleasePolicySaveRequest) error {
	return handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("policies"),
	)
}

func SetPolicy(bCtx *env.BubblyContext, req *api.ReleasePolicyUpdateRequest) error {
	return handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodPut),
		WithPayload(req), WithRequestURL(fmt.Sprintf("policies/%d", *req.ID)),
	)
}

func GetVulnerabilityReviews(bCtx *env.BubblyContext, req *api.VulnerabilityReviewGetRequest) (*api.VulnerabilityReviewGetResponse, error) {
	var r api.VulnerabilityReviewGetResponse
	if err := handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodGet),
		WithRequestURL("vulnerabilityreviews"), WithQueryParamsStruct(req), WithResponse(&r),
	); err != nil {
		return nil, err
	}
	return &r, nil
}

func SaveVulnerabilityReview(bCtx *env.BubblyContext, req *api.VulnerabilityReviewSaveRequest) error {
	return handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("vulnerabilityreviews"),
	)
}

func UpdateVulnerabilityReview(bCtx *env.BubblyContext, req *api.VulnerabilityReviewUpdateRequest) error {
	return handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodPut),
		WithPayload(req), WithRequestURL(fmt.Sprintf("vulnerabilityreviews/%d", *req.ID)),
	)
}
