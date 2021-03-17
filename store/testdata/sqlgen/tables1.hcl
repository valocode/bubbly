#
# One-to-many relationship:
# 	The table `A` is the parent table, and 
# 	the table `B` is the related table.
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
}
