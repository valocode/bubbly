#
# Multi-row, multi-table example for testing 
# the `order_by` GraphQL argument.
#

# This example is about a crew of characters,
# each crew housed in a hideaway.

# NB: This was **not** meant as an example of excellent database design. 
# It was built to provide the data necessary for testing, with a simple story.
# To find out how we can help you manage your comprehensive hideways list please contact Sales.

table "hideaways" {
	field "location" {
		type = string
		unique = true
	}

	field "sophistication" {
		type = string
	}

	field "distance_from_x" {
		type = number
	}

	field "ready" {
		type = bool
	}
}

table "characters" {
	field "name" {
		type = string
		unique = true
	}
}

table "crew" {
	join "characters" {}
	join "hideaways" {}

	field "count" {
		type = number
	}
}
