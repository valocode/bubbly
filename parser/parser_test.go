package parser

import (
	"testing"

	"github.com/verifa/bubbly/env"
)

func TestParser(t *testing.T) {
	bCtx := env.NewBubblyContext()
	p, err := NewParserFromFilename("testdata")
	if err != nil {
		t.Errorf("Could not initialise parser: %s", err.Error())
		t.FailNow()
	}

	err = p.Parse(bCtx)
	if err != nil {
		t.Errorf("Failed to decode parser body: %s", err.Error())
		t.FailNow()
	}
}
