package bubbly

violation[{"msg": msg, "severity": severity}] {
	tests := test_cases()
	failed_tests := [test | not tests[i].result; test := tests[i]]
	count(failed_tests) > 0
	msg := sprintf("%d failing tests(s)", [count(failed_tests)])
	severity := "error"
}
