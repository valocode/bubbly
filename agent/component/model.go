package component

import (
	"time"
)

type Request struct {
	Subject Subject       `json:"subject"`
	Data    []byte        `json:"data"`
	Timeout time.Duration `json:"timeout"`
	Reply   *Reply        `json:"reply"`
}

type Reply struct {
	Data  []byte `json:"data"`
	Error string `json:"error"`
}
