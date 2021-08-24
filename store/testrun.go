package store

import (
	"errors"
	"fmt"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

func (s *Store) SaveTestRun(req *api.TestRunRequest) (*ent.TestRun, error) {
	release, err := s.releaseFromCommit(req.Commit)
	if err != nil {
		return nil, err
	}
	return s.saveTestRun(release, req.TestRun)
}

func (s *Store) saveTestRun(release *ent.Release, run *api.TestRun) (*ent.TestRun, error) {
	var testRun *ent.TestRun
	txErr := s.WithTx(func(tx *ent.Tx) error {
		if run.Tool == nil {
			return errors.New("tool is required")
		}
		dbRun, err := tx.TestRun.Create().
			SetModelCreate(&run.TestRunModelCreate).
			SetRelease(release).
			Save(s.ctx)
		if err != nil {
			return fmt.Errorf("error creating test run: %w", err)
		}
		for _, tc := range run.TestCases {
			_, err := tx.TestCase.Create().
				SetRun(dbRun).
				SetModelCreate(&tc.TestCaseModelCreate).
				Save(s.ctx)
			if err != nil {
				return fmt.Errorf("error creating test case: %w", err)
			}
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return testRun, nil
}
