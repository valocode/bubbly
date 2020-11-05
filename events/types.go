package events

// Resource Run event reason list
const (
	CreatedResourceRun = "CreatedRun"
	StartedResourceRun = "StartedRun"
	FailedResourceRun  = "FailedRun"
	KillingResourceRun = "KillingRun"
)

// Resource event reason list
const (
	CreatingResource = "Creating"
	KillingResource  = "Deleting"
	KilledResource   = "Killed"
)

// bubbly describe importer example_importer
// 'bubbly get all -o events' is useful to have Kind
type Event struct {
	Status  string
	Kind    string // Importer, System, ...
	Age     string // TODO: formalise using time pkg
	Message string
}
