package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/event"
	"github.com/valocode/bubbly/ent/spdxlicense"
	"github.com/valocode/bubbly/store"
)

const (
	spdxLatestReleaseURL = "https://api.github.com/repos/spdx/license-list-data/releases/latest"
	spdxLicenseListURL   = "https://raw.githubusercontent.com/spdx/license-list-data/%s/json/licenses.json"
)

type (
	SPDXMonitor struct {
		ctx   context.Context
		store *store.Store
	}

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

func NewSPDXMonitor(opts ...func(m *SPDXMonitor)) (*SPDXMonitor, error) {
	monitor := SPDXMonitor{}
	for _, opt := range opts {
		opt(&monitor)
	}

	if monitor.store == nil {
		return nil, fmt.Errorf("spdx monitor requires a store")
	}
	// Set a default context if none was given. In production we'd expect
	// to receive a context
	if monitor.ctx == nil {
		monitor.ctx = context.Background()
	}

	return &monitor, nil
}

func WithContext(ctx context.Context) func(m *SPDXMonitor) {
	return func(m *SPDXMonitor) {
		m.ctx = ctx
	}
}

func WithStore(bstore *store.Store) func(m *SPDXMonitor) {
	return func(m *SPDXMonitor) {
		m.store = bstore
	}
}

func (m *SPDXMonitor) Do() error {
	client := m.store.Client()
	licenseList, err := fetchSPDXLicenses()
	if err != nil {
		// We cannot return an error, but save an event with stauts error
		_, eventErr := client.Event.Create().
			SetMessage(fmt.Sprintf("error fetching SPDX license list: %s", err.Error())).
			SetType(event.TypeMonitor).
			SetStatus(event.StatusError).
			Save(m.ctx)
		if eventErr != nil {
			// Not much we can do... Just log the error
			log.Errorf("saving event: %w", eventErr)
		}
		return err
	}

	// TODO: perhaps this logic should be in the store, rather than here and
	// "using" features of the store?
	txErr := store.WithTx(m.ctx, m.store.Client(), func(tx *ent.Tx) error {
		for _, lic := range licenseList.Licenses {
			spdxLic, err := tx.SPDXLicense.Query().Where(
				spdxlicense.LicenseID(lic.LicenseID),
			).Only(m.ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return store.HandleEntError(err, "query for SPDX license")
				}
				// If not found, create the SPDX license
				_, err = tx.SPDXLicense.Create().
					SetLicenseID(lic.LicenseID).
					SetName(lic.Name).
					SetReference(lic.Reference).
					SetDetailsURL(lic.DetailsURL).
					SetIsOsiApproved(lic.IsOSIApproved).
					Save(m.ctx)
				if err != nil {
					fmt.Printf("err: %#v\n", err)
					return store.HandleEntError(err, "create SPDX license")
				}
				// Continue to the next license
				continue
			}
			// If it was found, update it
			if _, err := tx.SPDXLicense.UpdateOne(spdxLic).
				SetName(lic.Name).
				SetReference(lic.Reference).
				SetDetailsURL(lic.DetailsURL).
				SetIsOsiApproved(lic.IsOSIApproved).
				Save(m.ctx); err != nil {
				return store.HandleEntError(err, "update SPDX license")
			}
		}
		return nil
	})
	if txErr != nil {
		// We cannot return an error, but save an event with stauts error
		_, eventErr := client.Event.Create().
			SetMessage(fmt.Sprintf("error updating SPDX license list: %s", txErr.Error())).
			SetType(event.TypeMonitor).
			SetStatus(event.StatusError).
			Save(m.ctx)
		if eventErr != nil {
			// Not much we can do... Just log the error
			log.Errorf("saving event: %w", eventErr)
		}
		return txErr
	}
	// If all worked successfully, create an event
	_, eventErr := client.Event.Create().
		SetMessage("SPDX license list successfully updated").
		SetType(event.TypeMonitor).
		SetStatus(event.StatusOk).
		Save(m.ctx)
	if eventErr != nil {
		// Not much we can do... Just log the error
		log.Errorf("saving event: %w", eventErr)
	}
	return nil
}

func fetchSPDXLicenses() (*SpdxLicenseList, error) {
	timeout := time.Duration(time.Second * 20)

	var (
		body    io.Reader
		release SpdxRelease
		list    SpdxLicenseList
	)
	{
		// Make the request to get the latest SPDX release
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
		// Make the HTTP request to get the SPDX license list
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
