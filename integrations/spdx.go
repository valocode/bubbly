package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	spdxLatestReleaseURL = "https://api.github.com/repos/spdx/license-list-data/releases/latest"
	spdxLicenseListURL   = "https://raw.githubusercontent.com/spdx/license-list-data/%s/json/licenses.json"
)

type (
	SpdxRelease struct {
		TagName string `json:"tag_name"`
	}

	SpdxLicenseList struct {
		ListVersion string        `json:"licenseListVersion"`
		Licenses    []SpdxLicense `json:"licenses"`
		ReleaseDate string        `json:"releaseDate"`
	}

	SpdxLicense struct {
		LicenseID     string `json:"licenseId"`
		Name          string `json:"name"`
		Reference     string `json:"reference"`
		DetailsURL    string `json:"detailsUrl"`
		IsOSIApproved bool   `json:"isOsiApproved"`
	}
)

func FetchSPDXLicenses() (*SpdxLicenseList, error) {
	timeout := time.Duration(time.Minute * 5)

	var body io.Reader
	var release SpdxRelease
	var list SpdxLicenseList
	{
		req, err := http.NewRequest(http.MethodGet, spdxLatestReleaseURL, body)
		if err != nil {
			return nil, err
		}
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
		if err := dec.Decode(&release); err != nil {
			return nil, err
		}
	}

	{
		spdxListURL := fmt.Sprintf(spdxLicenseListURL, release.TagName)
		req, err := http.NewRequest(http.MethodGet, spdxListURL, body)
		if err != nil {
			return nil, err
		}
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
		if err := dec.Decode(&list); err != nil {
			return nil, err
		}

	}

	return &list, nil
}
