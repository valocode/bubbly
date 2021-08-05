package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type (
	CVEResponse struct {
		ResultsPerPage int `json:"resultsPerPage"`
		StartIndex     int `json:"startIndex"`
		TotalResults   int `json:"totalResults"`

		Result struct {
			CVEDataFormat string `json:"CVE_data_format"`
			CVEItems      []struct {
				CVEItem CVEItem   `json:"cve"`
				Impact  CVEImpact `json:"impact"`
			} `json:"CVE_Items"`
		} `json:"result"`
	}

	CVEItem struct {
		CVEDataMeta struct {
			ID string `json:"ID"`
		} `json:"CVE_data_meta"`
		Description struct {
			DescriptionData []struct {
				Lang  string `json:"lang"`
				Value string `json:"value"`
			} `json:"description_data"`
		} `json:"description"`
		PublishedDate    string `json:"publishedDate"`
		LastModifiedDate string `json:"lastModifiedDate"`
	}

	CVEImpact struct {
		BaseMetricV3 struct {
			CVSSV3 struct {
				BaseScore float32 `json:"baseScore"`
			} `json:"cvssV3"`
			ExploitabilityScore float32 `json:"exploitabilityScore"`
			ImpactScore         float32 `json:"impactScore"`
		} `json:"baseMetricV3"`
	}
)

type (
	CVEFetchOptions struct {
		StartIndex     int
		ResultsPerPage int
		ModStartDate   *time.Time
		ModEndDate     *time.Time
		Sleep          time.Duration
	}
)

func DefaultCVEFetchOptions() CVEFetchOptions {
	return CVEFetchOptions{
		ResultsPerPage: 1000,
	}
}

type CVEHandler func(*CVEResponse) error

func InitFetchAllCVEs(fn CVEHandler) error {
	var (
		complete   bool
		startIndex = 0
	)
	for !complete {
		opts := CVEFetchOptions{
			StartIndex:     startIndex,
			ResultsPerPage: 1000,
			Sleep:          time.Second * 2,
		}
		resp, err := FetchCVEs(opts)
		if err != nil {
			return err
		}
		if err := fn(resp); err != nil {
			return err
		}

		startIndex = resp.StartIndex + resp.ResultsPerPage
		if startIndex > resp.TotalResults {
			complete = true
		}
	}

	return nil
}

func FetchCVEs(opts CVEFetchOptions) (*CVEResponse, error) {
	url := "https://services.nvd.nist.gov/rest/json/cves/1.0"
	timeout := time.Duration(time.Minute * 5)

	var body io.Reader
	req, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("startIndex", strconv.Itoa(opts.StartIndex))
	q.Add("resultsPerPage", strconv.Itoa(opts.ResultsPerPage))
	req.URL.RawQuery = q.Encode()

	// Before performing the request, sleep the required time
	time.Sleep(opts.Sleep)
	c := http.Client{
		Timeout: timeout,
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error getting body of response: %w", err)
		}
		return nil, fmt.Errorf("HTTP response status code: %d: %s", resp.StatusCode, body)
	}

	dec := json.NewDecoder(resp.Body)
	var cveResp CVEResponse
	if err := dec.Decode(&cveResp); err != nil {
		return nil, err
	}

	return &cveResp, nil
}
