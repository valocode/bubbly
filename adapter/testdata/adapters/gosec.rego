package adapter

somevar := "asd"

code_scan[scan] {
	scan := {"tool": "gosec"}
}

code_issue[issue] {
	results = json("./testdata/adapters/gosec.json")
	some i
	iss := results.Issues[i]
	issue := {
		# providing the i is necessary so that we get all unique code_issues
		"i": i,
		"rule_id": iss.rule_id,
		"message": iss.details,
		"severity": lower(iss.severity),
		"type": "security",
	}
}
