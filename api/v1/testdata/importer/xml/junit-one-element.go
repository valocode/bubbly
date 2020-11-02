package v1

import "github.com/zclconf/go-cty/cty"

// ExpectedValueOneElement returns the data structure as it should be after
// the value returned from the XML parser is converted to cty.Value,
// using the type information returned by the ExpectedType() from junit.go
func ExpectedValueOneElement() cty.Value {

	expected := cty.ObjectVal(map[string]cty.Value{

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

	return expected
}
