package store

import (
	"github.com/graphql-go/graphql"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

// Provider providea an interface for persisting readiness data.
type provider interface {
	Create([]core.Table) error
	Save(core.DataBlocks) ([]core.Table, error)
	ResolveScalar(graphql.ResolveParams) (interface{}, error)
	ResolveList(graphql.ResolveParams) (interface{}, error)
	LastValue(tableName, field string) (cty.Value, error)
}

// ProviderType is a store provider.
type ProviderType string

const (
	// Postgres is a Postgres provider.
	Postgres ProviderType = "postgres"
)
