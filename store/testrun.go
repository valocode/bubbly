package store

import (
	"errors"
	"fmt"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

func (h *Handler) SaveTestRun(req *api.TestRunRequest) (*ent.TestRun, error) {
	release, err := h.releaseFromCommit(req.Commit)
	if err != nil {
		return nil, err
	}
	return h.saveTestRun(release, req.TestRun)
}

func (h *Handler) saveTestRun(release *ent.Release, run *api.TestRun) (*ent.TestRun, error) {
	var testRun *ent.TestRun
	txErr := WithTx(h.ctx, h.client, func(tx *ent.Tx) error {
		if run.Tool == nil {
			return errors.New("tool is required")
		}
		dbRun, err := tx.TestRun.Create().
			SetModelCreate(&run.TestRunModelCreate).
			SetRelease(release).
			Save(h.ctx)
		if err != nil {
			return fmt.Errorf("error creating test run: %w", err)
		}
		for _, tc := range run.TestCases {
			_, err := tx.TestCase.Create().
				SetRun(dbRun).
				SetModelCreate(&tc.TestCaseModelCreate).
				Save(h.ctx)
			if err != nil {
				return fmt.Errorf("error creating test case: %w", err)
			}
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	// Once transaction is complete, evaluate the release.
	_, err := h.EvaluateReleasePolicies(release.ID)
	if err != nil {
		return nil, NewServerError(err, "evaluating release policies")
	}

	return testRun, nil
}
