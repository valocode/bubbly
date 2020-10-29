package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	p, err := NewParserFromFilename("testdata/example")
	if err != nil {
		t.Errorf("Could not initialise parser: %s", err.Error())
		t.FailNow()
	}

	err = p.Parse()
	if err != nil {
		t.Errorf("Failed to decode parser body: %s", err.Error())
		t.FailNow()
	}
}
