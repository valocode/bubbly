package bubbly

violation[violation] {
	issues := code_issues()
	high_issues := [issue | issues[i].severity == "high"; issue := issues[i]]
	count(high_issues) > 0
	violation = {
		"msg": sprintf("%d high issue(s)", [count(high_issues)]),
		"severity": "error",
	}
}
