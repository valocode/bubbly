package store

import (
	"io"

	"github.com/graphql-go/graphql"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

// Provider provides an interface for persisting readiness data.
type provider interface {
	Create(core.Tables) error
	Save(core.DataBlocks, core.DataBlocks) (core.Tables, error)
	ResolveScalar(graphql.ResolveParams) (interface{}, error)
	ResolveList(graphql.ResolveParams) (interface{}, error)
	LastValue(tableName, field string) (cty.Value, error)

	GetResource(id string) (io.Reader, error)
	PutResource(id string, val string) error
}
