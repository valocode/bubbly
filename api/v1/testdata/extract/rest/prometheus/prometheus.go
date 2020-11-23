package v1

import "github.com/zclconf/go-cty/cty"

// ExpectedType returns the data structure
// describing the dynamic type (cty.Type) of
// the value (cty.Value) for this test fixture
func ExpectedType() cty.Type {
	return cty.Object(map[string]cty.Type{
		"status": cty.String,
		"data": cty.Object(map[string]cty.Type{
			"startTime":           cty.String,
			"CWD":                 cty.String,
			"reloadConfigSuccess": cty.Bool,
			"lastConfigTime":      cty.String,
			"timeSeriesCount":     cty.Number,
			"corruptionCount":     cty.Number,
			"goroutineCount":      cty.Number,
			"GOMAXPROCS":          cty.Number,
			"GOGC":                cty.String,
			"GODEBUG":             cty.String,
			"storageRetention":    cty.String,
		}),
	})
}

// ExpectedValue returns the data structure describing
// the value (cty.Value) for this test fixture
func ExpectedValue() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"status": cty.StringVal("success"),
		"data": cty.ObjectVal(map[string]cty.Value{
			"startTime":           cty.StringVal("2019-11-02T17:23:59.301361365+01:00"),
			"CWD":                 cty.StringVal("/"),
			"reloadConfigSuccess": cty.BoolVal(true),
			"lastConfigTime":      cty.StringVal("2019-11-02T17:23:59+01:00"),
			"timeSeriesCount":     cty.NumberIntVal(873),
			"corruptionCount":     cty.NumberIntVal(0),
			"goroutineCount":      cty.NumberIntVal(48),
			"GOMAXPROCS":          cty.NumberIntVal(4),
			"GOGC":                cty.StringVal(""),
			"GODEBUG":             cty.StringVal(""),
			"storageRetention":    cty.StringVal("15d"),
		}),
	})
}
