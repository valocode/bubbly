package store

import (
	"fmt"

	"github.com/valocode/bubbly/test"
)

func (s *Store) PopulateStoreWithDummyData() error {
	data := test.CreateDummyData()
	for _, repo := range data {
		for _, release := range repo.Releases {
			if _, err := s.CreateRelease(release.Release); err != nil {
				return err
			}
			for _, art := range release.Artifacts {
				if _, err := s.LogArtifact(art); err != nil {
					return err
				}
			}
			for _, scan := range release.CodeScans {
				if _, err := s.SaveCodeScan(scan); err != nil {
					return err
				}
			}
			for _, run := range release.TestRuns {
				if _, err := s.SaveTestRun(run); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Store) PopulateStoreWithPolicies(basedir string) error {
	reqs, err := test.ParsePolicies(basedir)
	if err != nil {
		return err
	}
	for _, req := range reqs {
		if _, err := s.SaveReleasePolicy(req); err != nil {
			return err
		}
		fmt.Println("saved policy: ", *req.Policy.Name)
	}
	return nil
}
