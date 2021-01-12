package component

type Subscriptions []Subscription

// A Subscription is a Go-native representation of a queue-group NATS
// Subscription.
// Bubbly components use Subscriptions to subscribe to NATS channels and
// therefore for cross-component communication
type Subscription struct {
	Subject Subject
	Queue   Queue
}

type Subjects []Subject

// Subject represents a string matching a NATS Subject.
// Bubbly components use Subjects within their Publish and
// Subscribe method signatures in order to communicate with one another
// Therefore, any cross-component communication requires a Subject
type Subject string

// Any Subjects that components use to communicate with one another should be
// defined centrally here
const (
	WorkerPipelineRunIntervalSubject Subject = "pipeline_run_interval"
)

type Queues []Queue

// Queue represents a string matching a NATS Queue group,
// which NATS uses to load balance messages across groups of subscribers
// Bubbly components use Queues when subscribing on a given subject.
type Queue string

// Any Queue that components use as a part of their Subjects should be
// defined centrally here
const (
	QueueWorker Queue = "worker"
)
