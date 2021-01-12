package component

type Publications []Publication

// Publications are used by bubbly components to publish to NATS channels
type Publication struct {
	Subject Subject
	Value   interface{}
}
