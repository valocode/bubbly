package v1

import "github.com/zclconf/go-cty/cty"

// ExpectedType returns the data structure describing the type which
// the data imported by XML Importer in this test case will conform to.
func ExpectedType() cty.Type {
	return cty.Object(map[string]cty.Type{
		"testsuites": cty.Object(map[string]cty.Type{
			"duration": cty.Number,
			"testsuite": cty.List(cty.Object(map[string]cty.Type{
				"failures": cty.Number,
				"name":     cty.String,
				"package":  cty.String,
				"tests":    cty.Number,
				"time":     cty.Number,
				"testcase": cty.List(cty.Object(map[string]cty.Type{
					"classname": cty.String,
					"name":      cty.String,
					"time":      cty.Number,
				})),
			})),
		}),
	})
}

// func ExpectedValueX() cty.Value:
//   returns the data structure as it should be after
//   the value returned from the XML parser is converted to cty.Value,
//   using the type information returned by the ExpectedType() function.

// ExpectedValue0 for junit0.xml
func ExpectedValue0() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{

		"testsuites": cty.ObjectVal(map[string]cty.Value{

			"duration": cty.NumberFloatVal(50.5),
			"testsuite": cty.ListVal([]cty.Value{

				cty.ObjectVal(map[string]cty.Value{
					"failures": cty.NumberIntVal(0),
					"name":     cty.StringVal("Untitled suite in /Users/niko/Sites/casperjs/tests/suites/casper/agent.js"),
					"package":  cty.StringVal("tests/suites/casper/agent"),
					"tests":    cty.NumberIntVal(3),
					"time":     cty.NumberFloatVal(0.256),

					"testcase": cty.ListVal([]cty.Value{

						cty.ObjectVal(map[string]cty.Value{
							"classname": cty.StringVal("tests/suites/casper/agent"),
							"name":      cty.StringVal("Default user agent matches /CasperJS/"),
							"time":      cty.NumberFloatVal(0.103),
						}),

						cty.ObjectVal(map[string]cty.Value{
							"classname": cty.StringVal("tests/suites/casper/agent"),
							"name":      cty.StringVal("Default user agent matches /plop/"),
							"time":      cty.NumberFloatVal(0.146),
						}),

						cty.ObjectVal(map[string]cty.Value{
							"classname": cty.StringVal("tests/suites/casper/agent"),
							"name":      cty.StringVal("Default user agent matches /plop/"),
							"time":      cty.NumberFloatVal(0.007),
						}),
					}),
				}),

				cty.ObjectVal(map[string]cty.Value{
					"failures": cty.NumberIntVal(0),
					"name":     cty.StringVal("Untitled suite in /Users/niko/Sites/casperjs/tests/suites/casper/auth.js"),
					"package":  cty.StringVal("tests/suites/casper/auth"),
					"tests":    cty.NumberIntVal(2),
					"time":     cty.NumberFloatVal(0.101),

					"testcase": cty.ListVal([]cty.Value{

						cty.ObjectVal(map[string]cty.Value{
							"classname": cty.StringVal("tests/suites/casper/auth"),
							"name":      cty.StringVal("Subject equals the expected value"),
							"time":      cty.NumberFloatVal(0.1),
						}),

						cty.ObjectVal(map[string]cty.Value{
							"classname": cty.StringVal("tests/suites/casper/auth"),
							"name":      cty.StringVal("Subject equals the expected value"),
							"time":      cty.NumberFloatVal(0),
						}),
					}),
				}),
			}),
		}),
	})
}

// ExpectedValue1 for junit1.xml
func ExpectedValue1() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{

		"testsuites": cty.ObjectVal(map[string]cty.Value{

			"duration": cty.NumberFloatVal(50.5),
			"testsuite": cty.ListVal([]cty.Value{

				cty.ObjectVal(map[string]cty.Value{
					"failures": cty.NumberIntVal(0),
					"name":     cty.StringVal("Untitled suite in /Users/niko/Sites/casperjs/tests/suites/casper/agent.js"),
					"package":  cty.StringVal("tests/suites/casper/agent"),
					"tests":    cty.NumberIntVal(3),
					"time":     cty.NumberFloatVal(0.256),

					"testcase": cty.ListVal([]cty.Value{

						cty.ObjectVal(map[string]cty.Value{
							"classname": cty.StringVal("tests/suites/casper/agent"),
							"name":      cty.StringVal("Default user agent matches /CasperJS/"),
							"time":      cty.NumberFloatVal(0.103),
						}),

						cty.ObjectVal(map[string]cty.Value{
							"classname": cty.StringVal("tests/suites/casper/agent"),
							"name":      cty.StringVal("Default user agent matches /plop/"),
							"time":      cty.NumberFloatVal(0.146),
						}),

						cty.ObjectVal(map[string]cty.Value{
							"classname": cty.StringVal("tests/suites/casper/agent"),
							"name":      cty.StringVal("Default user agent matches /plop/"),
							"time":      cty.NumberFloatVal(0.007),
						}),
					}),
				}),
			}),
		}),
	})
}

// ExpectedValue2 for junit2.xml
func ExpectedValue2() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{

		"testsuites": cty.ObjectVal(map[string]cty.Value{

			"duration": cty.NumberFloatVal(50.5),
			"testsuite": cty.ListVal([]cty.Value{

				cty.ObjectVal(map[string]cty.Value{
					"failures": cty.NumberIntVal(0),
					"name":     cty.StringVal("Untitled suite in /Users/niko/Sites/casperjs/tests/suites/casper/agent.js"),
					"package":  cty.StringVal("tests/suites/casper/agent"),
					"tests":    cty.NumberIntVal(3),
					"time":     cty.NumberFloatVal(0.256),

					"testcase": cty.ListVal([]cty.Value{

						cty.ObjectVal(map[string]cty.Value{
							"classname": cty.StringVal("tests/suites/casper/agent"),
							"name":      cty.StringVal("Default user agent matches /CasperJS/"),
							"time":      cty.NumberFloatVal(0.103),
						}),
					}),
				}),
			}),
		}),
	})
}
