#
# Multi-level nested and sibling one-to-many relationships:
#   The table `A` is the parent table to `B` and `D`,
#     the table `B` is the related table to `A`.
#     the table `D` is the related table to `A`.
#
#   The table `B` is a parent table to `C`,
#     the table `C` is the related table to `B.
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

	table "D" {
		field "whddt" {
			type = string
		}
	}
}
