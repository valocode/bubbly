#
# Data for tables4.hcl
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

data "C" {
	joins = ["B"]
	fields = {
		"whcct": "vc1"
	}
}

data "D" {
	joins = ["A"]
	fields = {
		"whddt": "vd1"
	}
}
