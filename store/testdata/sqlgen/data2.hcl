#
# Data for tables2.hcl
#
data "A" {
	fields {
		whaat = "va1"
	}
}

data "B" {
	joins = ["A"]
	fields {
		whbbt = "vb1"
	}
}

data "C" {
	joins = ["A"]
	fields {
		whcct = "vc1"
	}
}
