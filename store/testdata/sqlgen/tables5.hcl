#
# Multi-row example with sibling fields.
#
table "test_run" {

	join "location" {}
	join "configuration" {}
	join "version" {}

	field "ok" {
		type = bool
	}
}

table "location" {	
	field "name" {
		type = string
		unique = true
	}
}

table "configuration" {
	field "name" {
		type = string
		unique = true
	}
}

table "version" {
	field "name" {
		type = string
		unique = true
	}
}
