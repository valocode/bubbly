package v1

import "github.com/zclconf/go-cty/cty"

// ExpectedType returns the data structure describing the type which
// the data extracted by JSON Extract in this test case will conform to
func ExpectedType() cty.Type {

	ctyType := cty.Object(map[string]cty.Type{
		"issues": cty.List(cty.Object(map[string]cty.Type{
			"engineId": cty.String,
			"ruleId":   cty.String,
			"severity": cty.String,
			"type":     cty.String,
			"primaryLocation": cty.Object(map[string]cty.Type{
				"message":  cty.String,
				"filePath": cty.String,
				"textRange": cty.Object(map[string]cty.Type{
					"startLine":   cty.Number,
					"endLine":     cty.Number,
					"startColumn": cty.Number,
					"endColumn":   cty.Number,
				}),
			}),
		})),
	})

	return ctyType
}

// ExpectedValue returns the data structure as it should be after
// the value returned from the JSON unmarhalling library is converted
// to cty.Value using the type information provided by the ExpectedType()
func ExpectedValue() cty.Value {

	expected := cty.ObjectVal(map[string]cty.Value{
		"issues": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"engineId": cty.StringVal("test"),
				"primaryLocation": cty.ObjectVal(map[string]cty.Value{
					"filePath": cty.StringVal("sources/A.java"),
					"message":  cty.StringVal("fully-fleshed issue"),
					"textRange": cty.ObjectVal(
						map[string]cty.Value{
							"endColumn":   cty.NumberIntVal(14),
							"endLine":     cty.NumberIntVal(30),
							"startColumn": cty.NumberIntVal(9),
							"startLine":   cty.NumberIntVal(30),
						}),
				}),
				"ruleId":   cty.StringVal("rule1"),
				"severity": cty.StringVal("BLOCKER"),
				"type":     cty.StringVal("CODE_SMELL"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"engineId": cty.StringVal("test"),
				"primaryLocation": cty.ObjectVal(map[string]cty.Value{
					"filePath": cty.StringVal("sources/Measure.java"),
					"message":  cty.StringVal("minimal issue raised at file level"),
					"textRange": cty.NullVal(
						cty.Object(
							map[string]cty.Type{
								"endColumn":   cty.Number,
								"endLine":     cty.Number,
								"startColumn": cty.Number,
								"startLine":   cty.Number,
							}),
					)},
				),
				"ruleId":   cty.StringVal("rule2"),
				"severity": cty.StringVal("INFO"),
				"type":     cty.StringVal("BUG"),
			}),
		}),
	})

	return expected
}
