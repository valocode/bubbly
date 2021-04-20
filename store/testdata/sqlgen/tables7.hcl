#
# Multi-row example for testing `order_by` 
# and `distinct_on` arguments working together.
#
table "test_run" {

	join "location" {}
	join "configuration" {}
	join "version" {}

	field "ok" {
		type = bool
	}

	# UNIX epoch: https://www.epochconverter.com/
	# This indicates when the test run finished.
	# For the test_run records tied at (location, configuration, version),
	# the most recent one is the one with the greatest finish_epoch value.
	field "finish_epoch" {
		type = string
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
