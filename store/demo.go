package store

import (
	"fmt"

	"github.com/valocode/bubbly/store/api"
	"github.com/valocode/bubbly/test"
)

func (h *Handler) PopulateStoreWithDummyData() error {
	data := test.CreateDummyData()
	for _, repo := range data {
		for _, release := range repo.Releases {
			if _, err := h.CreateRelease(release.Release); err != nil {
				return err
			}
			for _, art := range release.Artifacts {
				if _, err := h.LogArtifact(art); err != nil {
					return err
				}
			}
			for _, scan := range release.CodeScans {
				if _, err := h.SaveCodeScan(scan); err != nil {
					return err
				}
			}
			for _, run := range release.TestRuns {
				if _, err := h.SaveTestRun(run); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (h *Handler) PopulateStoreWithPolicies(basedir string) error {
	reqs, err := test.ParsePolicies(basedir)
	if err != nil {
		return err
	}
	for _, req := range reqs {
		req.Affects = &api.ReleasePolicyAffects{
			Projects: []string{"demo"},
		}
		if _, err := h.SaveReleasePolicy(req); err != nil {
			return err
		}
		fmt.Println("saved policy: ", *req.Policy.Name)
	}
	return nil
}
