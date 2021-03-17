#
# Multiple one-to-many relationships as siblings:
#   The table `A` is the parent table, and 
#     the table `B` is the related table.
#     the table `C` is another related table.
#
table "A" {
	field "whaat" {
		type = string
	}

	table "B" {
		field "whbbt" {
			type = string
		}
	}

	table "C" {
		field "whcct" {
			type = string
		}
	}
}
