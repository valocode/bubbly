package policy

violation[item] {
	tests := test_cases()

	failed_tests = [test | test := tests[_]; not test.result]

	count(failed_tests) > 0
	item = {
		"message": sprintf("%d failing tests(s)", [count(failed_tests)]),
		"severity": "error",
	}
}
