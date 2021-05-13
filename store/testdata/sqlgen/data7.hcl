#
# Data for tables7.hcl
#

### Set 1: multiple entries for (location, configuration, version)
###        being ("Oliver's NEON workstation", "So-so", "v1.0.1"),
###        with the most recent entry being "ok", and a mixed bag
###        of results for the previous entries.

data "location" {	
	fields = {
		"name": "Oliver's NEON workstation"
	}
}

data "configuration" {
	fields = {
		"name": "So-so"
	}
}

data "version" {
	fields = {
		"name": "v1.0.1"
	}
}

# Latest `testrun` in this set
data "testrun" {
	joins = [
		"location",
		"configuration",
		"version",
		]
	fields = {
		"ok": true      
		"finish_epoch": "1611111111"
	}
}

data "testrun" {
	joins = [
		"location",
		"configuration",
		"version",
		]
	fields = {
		"ok": false
		"finish_epoch": "1311111111"
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
		"finish_epoch": "1411111111"
	}
}
