#
# Example for testing the GraphQL `limit` argument.
#
table "events" {

	join "location" {}

	field "severity" {
		type = string
	}

	# UNIX epoch: https://www.epochconverter.com/
	field "timestamp" {
		type = number
	}
}

table "location" {	
	field "friendly_name" {
		type = string
		unique = true
	}

	# UNIX epoch: https://www.epochconverter.com/
	field "created_at" {
		type = string
	}
}
