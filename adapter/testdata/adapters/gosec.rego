package adapter

code_scan[scan] {
	scan := {
		"tool": "gosec",
		"metadata": {"env": {"some_var": "some_value"}},
	}
}

code_issue[issue] {
	some i
	iss := input[_].Issues[i]
	issue := {
		# providing the i is necessary so that we get all unique code_issues
		"i": i,
		"rule_id": iss.rule_id,
		"message": iss.details,
		"severity": lower(iss.severity),
		"type": "security",
	}
}
