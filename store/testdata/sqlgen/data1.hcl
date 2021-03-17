#
# Data for tables1.hcl
#
data "A" {
	fields = {
		"whaat": "va1"
	}
}

data "B" {
	joins = ["A"]
	fields = {
		"whbbt": "vb1"
	}
}
