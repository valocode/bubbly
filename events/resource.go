package events

// Event represents a lifecycle change to a Resource.
type Event struct {
	Error  interface{} // The error responsible for a failure status
	Status string
	Time   string // The time at which the Event occurred
}

type Status int

const (
	ResourceCreatedUpdated Status = iota // store V1 does not distinguish between these two lifecycle states
	ResourceDestroyed
	ResourceApplySuccess
	ResourceApplyFailure
	ResourceRunSuccess
	ResourceRunFailure
)

func (s Status) String() string {
	return [...]string{
		"Created/Updated",
		"Destroyed",
		"ApplySuccess",
		"ApplyFailure",
		"RunSuccess",
		"RunFailure"}[s]
}

// ResourceOutputStatus represents the output statuses for a resource
type ResourceOutputStatus string

// String gets a string value of a ResourceOutputStatus
func (r *ResourceOutputStatus) String() string {
	return string(*r)
}
