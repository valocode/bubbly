package policy

entries := release_entries()

gosec_scan {
	entries[_].type == "code_scan"
	entries[_].edges.code_scan.tool == "gosec"
}

require[violation] {
	not gosec_scan
	violation = {
		"message": "missing gosec scan",
		"severity": "warning",
	}
}
