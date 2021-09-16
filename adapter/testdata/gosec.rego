package adapter

code_scan[scan] {
	scan := {
		"tool": "gosec",
		"metadata": {"env": {"some_var": "some_value"}},
	}
}

# Use array comprehension to create a list of code_issue
# https://www.openpolicyagent.org/docs/latest/policy-language/#array-comprehensions
code_issue := [issue |
	iss := input[_].Issues[_]
	issue := {
		"rule_id": iss.rule_id,
		"message": iss.details,
		"severity": lower(iss.severity),
		"type": "security",
	}
]
