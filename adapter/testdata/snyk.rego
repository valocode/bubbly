package adapter

code_scan[scan] {
	scan := {"tool": "snyk"}
}

component := [comp |
	vuln := input[_].vulnerabilities[_]
	comp := {
		"name": vuln.name,
		"version": vuln.version,
		"vulnerabilities": [{
			"vid": vuln.id,
			"severity_score": vuln.cvssScore,
			# if you want to add a patch which becomes a vulnerability review
			# "patch": {"message": "test patch"},
		}],
	}
]
