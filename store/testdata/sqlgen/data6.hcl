#
# Data for tables6.hcl
#

### Hideaway 1

data "hideaways" {	
	fields = {
		"location":         "Deep Dark Wood",
		"sophistication":   "simple",
		"distance_from_x":  1500,
		"ready":            true,
	}
}

# The Totoros are inserted into the DB out-of-order.
# The ASC sorting order would be: Chibi-Totoro, Chuu-Totoro, Oh-Totoro
# The DESC sorting order would be: Oh-Totoro, Chuu-Totoro, Chibi-Totoro
# We insert: Chibi-Totoro, Chuu-Totoro, Oh-Totoro
# In the absence of fuzzy testing, this should provide better assurance
# that our order_by is working correctly.

data "characters" {
	fields = {
		"name": "Oh-Totoro",
	}
}

data "crew" {
	joins = [
		"hideaways",
		"characters",
		]
	fields = {
		"count": 1,
	}
}

data "characters" {
	fields = {
		"name": "Chibi-Totoro",
	}
}

data "crew" {
	joins = [
		"hideaways",
		"characters",
		]
	fields = {
		"count": 1,
	}
}

data "characters" {
	fields = {
		"name": "Chuu-Totoro",
	}
}

data "crew" {
	joins = [
		"hideaways",
		"characters",
		]
	fields = {
		"count": 1,
	}
}

### Hideaway 2

data "hideaways" {	
	fields = {
		"location":         "Secret Underground Facility on the Moon",
		"sophistication":   "incredible",
		"distance_from_x":  384400,
		"ready":            true,
	}
}

# The Green Men are inserted into the DB in ASC order: 
#  - "Little Green Man"
#  - "Little Green Man Leader"
# When sorted by (crew.count ASC, characters.name DESC),
# the "Little Green Man Leader" would come before 
# the "Little Green Man". This would be a check for ordering
# not just on string values, but also on numbers.

data "characters" {
	fields = {
		"name": "Little Green Man",
	}
}

data "crew" {
	joins = [
		"hideaways",
		"characters",
		]
	fields = {
		"count": 42,
	}
}

data "characters" {
	fields = {
		"name": "Little Green Man Leader",
	}
}

data "crew" {
	joins = [
		"hideaways",
		"characters",
		]
	fields = {
		"count": 1,
	}
}

### Hideaway 3

data "hideaways" {	
	fields = {
		"location":         "Gold Coast Villa",
		"sophistication":   "average",
		"distance_from_x":  12500,
		"ready":            false,
	}
}

data "characters" {
	fields = {
		"name": "Wombat",
	}
}

data "crew" {
	joins = [
		"hideaways",
		"characters",
		]
	fields = {
		"count": 1,
	}
}

data "characters" {
	fields = {
		"name": "Smol European Mouse",
	}
}

data "crew" {
	joins = [
		"hideaways",
		"characters",
		]
	fields = {
		"count": 2,
	}
}

data "characters" {
	fields = {
		"name": "Old Horse with Shady Past",
	}
}

data "crew" {
	joins = [
		"hideaways",
		"characters",
		]
	fields = {
		"count": 1,
	}
}