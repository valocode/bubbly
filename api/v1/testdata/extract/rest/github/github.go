package v1

import "github.com/zclconf/go-cty/cty"

// ExpectedType returns the data structure
// describing the dynamic type (cty.Type) of
// the value (cty.Value) for this test fixture.
func ExpectedType() cty.Type {
	return cty.List(cty.Object(map[string]cty.Type{
		"name": cty.String,
		"commit": cty.Object(map[string]cty.Type{
			"sha": cty.String,
			"url": cty.String,
		}),
		"protected": cty.Bool,
	}))
}

// ExpectedValue returns the data structure describing
// the value (cty.Value) for this test fixture
func ExpectedValue() cty.Value {
	return cty.ListVal([]cty.Value{

		cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal("master"),
			"commit": cty.ObjectVal(map[string]cty.Value{
				"sha": cty.StringVal("7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"),
				"url": cty.StringVal("https://api.github.com/repos/octocat/Hello-World/commits/7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"),
			}),
			"protected": cty.BoolVal(false),
		}),

		cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal("octocat-patch-1"),
			"commit": cty.ObjectVal(map[string]cty.Value{
				"sha": cty.StringVal("b1b3f9723831141a31a1a7252a213e216ea76e56"),
				"url": cty.StringVal("https://api.github.com/repos/octocat/Hello-World/commits/b1b3f9723831141a31a1a7252a213e216ea76e56"),
			}),
			"protected": cty.BoolVal(false),
		}),

		cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal("test"),
			"commit": cty.ObjectVal(map[string]cty.Value{
				"sha": cty.StringVal("b3cbd5bbd7e81436d2eee04537ea2b4c0cad4cdf"),
				"url": cty.StringVal("https://api.github.com/repos/octocat/Hello-World/commits/b3cbd5bbd7e81436d2eee04537ea2b4c0cad4cdf"),
			}),
			"protected": cty.BoolVal(false),
		}),
	})
}
