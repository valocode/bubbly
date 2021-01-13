package component

type Publications []Publication

// Publications are used by bubbly components to publish to NATS channels,
// either when reaching out to subscriptions OR when replying directly to a
// NATS request
type Publication struct {
	Subject Subject
	Data    []byte
	Encoder string
	Error   error
}
