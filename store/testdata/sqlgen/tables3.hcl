#
# Multi-level nested one-to-many relationships:
# 	The table `A` is the parent table to `B`, and 
# 	the table `B` is the related table to `A`.
#
#   The table `B` is a parent table to `C`,
#   and the table `C` is the related table to `B.
#
table "A" {
	field "whaat" {
		type = string
	}

	table "B" {
		field "whbbt" {
			type = string
		}

		table "C" {
			field "whcct" {
				type = string
			}
		}
	}
}
