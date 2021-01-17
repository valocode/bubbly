package store

import (
	"github.com/graphql-go/graphql"
)

// Provider provides an interface for persisting readiness data.
type provider interface {
	Apply(*bubblySchema) error
	Save(*saveContext) error
	ResolveScalar(graphql.ResolveParams) (interface{}, error)
	ResolveList(graphql.ResolveParams) (interface{}, error)

	// GetResource(id string) (io.Reader, error)
	// PutResource(id, val string) error
}
