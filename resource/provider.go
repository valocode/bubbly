package resource

import "github.com/verifa/bubbly/env"

type provider interface {
	Query(*env.BubblyContext, string) (string, error)
	Save(*env.BubblyContext, string, string) error
}
