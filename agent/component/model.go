package component

import (
	"time"
)

type Request struct {
	Subject Subject       `json:"subject"`
	Data    MessageData   `json:"data"`
	Timeout time.Duration `json:"timeout"`
	Reply   *Reply        `json:"reply"`
}

// MessageData represents the data in a NATS message.
// For our use case, we need to send authentication data as well as the actual
// data. As most of the actual data comes straight from the REST API (and is
// already in []byte) it makes sense to just pass this through as []byte and
// handle any unmarshalling on the agent handler side
type MessageData struct {
	Auth *MessageAuth `json:"auth"`
	Data []byte       `json:"data"`
}

// MessageAuth contains information about the user making the request and the
// organization against which the request is being made.
// This needs to be passed to the agents so that they can serve the request for
// multiple tenants
type MessageAuth struct {
	Organization string `json:"organization"`
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
}

type Reply struct {
	Data  []byte `json:"data"`
	Error string `json:"error"`
}
