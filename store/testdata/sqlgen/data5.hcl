#
# Data for tables5.hcl
#

### Entry 1

data "location" {	
	fields = {
		"name": "Deep Dark Wood"
	}
}

data "configuration" {
	fields = {
		"name": "Primitive"
	}
}

data "version" {
	fields = {
		"name": "v1.0.1"
	}
}

data "testrun" {
	joins = [
		"location",
		"configuration",
		"version",
		]
	fields = {
		"ok": true
	}
}

### Entry 2

data "location" {
	fields = {
		"name": "Gold Coast City Skyline"
	}
}

data "configuration" {
	fields = {
		"name": "Approved by Wombats"
	}
}

data "version" {
	fields = {
		"name": "v2.5"
	}
}


data "testrun" {
	joins = [
		"version",
		"configuration",
		"location",
	]
	fields = {
		"ok": false
	}
}

### Entry 3

data "location" {
	fields = {
		"name": "Secret Underground Facility on the Moon"
	}
}

data "configuration" {
	fields = {
		"name": "Magic"
	}
}

data "version" {
	fields = {
		"name": "X.Y.Z-1"
	}
}

data "testrun" {
	joins = [
		"location",
		"version",
		"configuration",
	]
	fields = {
		"ok": true
	}
}
