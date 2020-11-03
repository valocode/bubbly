package bubbly

import (
	"testing"
)

func TestApply(t *testing.T) {
	// TODO: should provide an apply that has a pipeline_run
	if err := Apply("../parser/testdata/git"); err != nil {
		t.Errorf("Failed to apply resources: %s", err.Error())
	}
}
