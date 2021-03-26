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

// data "version" {
// 	fields = {
// 		"name": "v1.0.1"
// 	}
// }

data "test_run" {
	joins = [
		"location",
		"configuration",
//		"version",
		]
	fields = {
		"ok": true
	}
}

### Entry 2

// data "location" {
// 	fields = {
// 		"name": "Gold Coast City Skyline"
// 	}
// }

// data "test_run" {
// 	joins = [
// 		"location"
// 	]
// 	fields = {
// 		"ok": false
// 	}
// }

### Entry 3

// data "location" {
// 	fields = {
// 		"name": "Secret Underground Facility on the Moon"
// 	}
// }

// data "test_run" {
// 	joins = [
// 		"location"
// 	]
// 	fields = {
// 		"ok": true
// 	}
// }
