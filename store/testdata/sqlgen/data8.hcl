#
# Data for tables8.hcl
#

# The Deep Dark Wood had been created in the beginning of time.
# It has multiple events associated with it.
data "location" {	
	fields {
		friendly_name = "Deep Dark Wood"
		created_at = 0
	}
}

data "events" {
	joins = [
		"location"
	]
	fields {
		severity = "INFO"
		timestamp = 10
	}
}

data "events" {
	joins = [
		"location",
	]
	fields {
		severity = "DEBUG"
		timestamp = 20
	}
}

data "events" {
	joins = [
		"location",
	]
	fields {
		severity = "CRITICAL"
		timestamp = 30
	}
}

# Another location created much later...
# It has only a single event associated with it.
data "location" {
	fields {
		friendly_name = "Secret Underground Facility on the Moon"
		created_at = 1000
	}
}

data "events" {
	joins = [
		"location",
	]
	fields {
		severity = "DEBUG"
		timestamp = 1010
	}
}


# Yet another location created last.
# It has no events associated with it.
data "location" {
	fields {
		friendly_name = "Gold Coast Villa"
		created_at = 2000
	}
}
