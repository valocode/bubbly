package bubbly

import (
	"testing"
)

func TestApply(t *testing.T) {
	if err := Apply("../parser/testdata/local-sq-json"); err != nil {
		t.Errorf("Failed to apply resources: %s", err.Error())
	}
}
