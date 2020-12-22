package parser

// TODO: cannot test using api.ResourceParserType because of a cyclic dependency...
// This needs to be re-done when the parser is reworked. There is some work here...

// func TestParser(t *testing.T) {
// 	bCtx := env.NewBubblyContext()
// 	value := api.ResourcesParserType{}
// 	p, err := NewParserFromFilename("testdata")
// 	if err != nil {
// 		t.Errorf("Could not initialise parser: %s", err.Error())
// 		t.FailNow()
// 	}

// 	err = p.Parse(bCtx, &value)
// 	if err != nil {
// 		t.Errorf("Failed to decode parser body: %s", err.Error())
// 		t.FailNow()
// 	}
// }
