package component

import "time"

type Requests []Request

// https://github.com/nats-io/nats.go/blob/v1.10.0/nats.go#L2793
// matches the nats.Conn.Request() signature, which returns a msg and err
type Request struct {
	Subject Subject
	Data    []byte
	Timeout time.Duration
}
