package store

import (
	"github.com/graphql-go/graphql"
)

// Provider provides an interface for persisting readiness data.
type provider interface {
	Apply(*bubblySchema) error
	Save(*bubblySchema, dataTree) error
	ResolveQuery(graph *schemaGraph, params graphql.ResolveParams) (interface{}, error)
}
