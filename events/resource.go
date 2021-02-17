package events

type Event int

const (
	ResourceCreatedUpdated Event = iota // store V1 does not distinguish between these two lifecycle states
	ResourceDestroyed
	ResourceApplySuccess
	ResourceApplyFailure
	ResourceRunSuccess
	ResourceRunFailure
)

func (e Event) String() string {
	return [...]string{
		"Created/Updated",
		"Destroyed",
		"ApplySuccess",
		"ApplyFailure",
		"RunSuccess",
		"RunFailure"}[e]
}

// ResourceOutputStatus represents the output statuses for a resource
type ResourceOutputStatus string

// String gets a string value of a ResourceOutputStatus
func (r *ResourceOutputStatus) String() string {
	return string(*r)
}
