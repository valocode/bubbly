package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/graphql-go/graphql"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

// New creates a new Store for the given config.
func New(bCtx *env.BubblyContext) (*Store, error) {
	var (
		s   = &Store{}
		p   provider
		err error
	)

	switch bCtx.StoreConfig.Provider {
	case config.PostgresStore:
		p, err = newPostgres(bCtx)
	case config.CockroachDBStore:
		p, err = newCockroachdb(bCtx)
	default:
		return nil, fmt.Errorf(`invalid provider: "%s"`, bCtx.StoreConfig.Provider)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	s.p = p
	schema, err := s.currentBubblySchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get current schema: %w", err)
	}

	if err := s.setGraphQLSchema(schema); err != nil {
		return nil, fmt.Errorf("failed to set graphql schema: %w", err)
	}

	return s, nil
}

// Store provides access to persisted readiness data.
type Store struct {
	p provider

	mu     sync.RWMutex
	schema *graphql.Schema
}

// Schema gets the graphql schema for the store.
func (s *Store) Schema() *graphql.Schema {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.schema
}

// Query queries the store.
func (s *Store) Query(query string) (interface{}, error) {
	if s.schema == nil {
		return nil, errors.New("cannot perform query without a schema. Apply a schema first")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := graphql.Do(graphql.Params{
		Schema:        *s.schema,
		RequestString: query,
	})

	if res.HasErrors() {
		return nil, fmt.Errorf("failed to execute query: %v", res.Errors)
	}

	return res.Data, nil
}

// Apply applies a schema corresponding to a set of tables.
func (s *Store) Apply(tables core.Tables) error {
	currentSchema, err := s.currentBubblySchema()
	if err != nil {
		return fmt.Errorf("failed to get current schema: %w", err)
	}

	// Append the internal tables containing definition of the schema and
	// resource tables.
	tables = append(tables, internalTables...)
	newSchemaTables := make(map[string]core.Table)
	for _, table := range tables {
		newSchemaTables[table.Name] = table
	}
	newSchema := &bubblySchema{
		Tables: newSchemaTables,
	}
	addImplicitJoins(currentSchema, tables, nil)
	// addImplicitJoins(newSchema, tables, nil)

	// Calculate the schema diff
	cl, err := compareSchema(*currentSchema, *newSchema)
	if err != nil {
		return fmt.Errorf("failed to compare schemas: %w", err)
	}
	newSchema.Changelog = cl
	// TODO call schema migration

	if err := s.p.Apply(currentSchema); err != nil {
		return fmt.Errorf("failed to apply schema in provider: %w", err)
	}

	if err := s.setGraphQLSchema(currentSchema); err != nil {
		return fmt.Errorf("failed to set graphql schema: %w", err)
	}

	return nil
}

// Save saves data into the store.
func (s *Store) Save(data core.DataBlocks) error {

	dataTree, err := createDataTree(data)
	if err != nil {
		return fmt.Errorf("failed to create tree of data blocks for storing: %w", err)
	}
	bubblySchema, err := s.currentBubblySchema()
	if err != nil {
		return fmt.Errorf("failed to get current schema: %w", err)
	}

	if err := s.p.Save(bubblySchema, dataTree); err != nil {
		return fmt.Errorf("falied to save data in provider: %w", err)
	}

	gqlSchema, err := newGraphQLSchema(bubblySchema, s.p)
	if err != nil {
		return fmt.Errorf("falied to build GraphQL schema: %w", err)
	}

	s.mu.Lock()
	s.schema = &gqlSchema
	s.mu.Unlock()

	return nil
}

func (s *Store) setGraphQLSchema(schema *bubblySchema) error {
	// Cannot create a graphql schema from an empty bubbly schema
	if len(schema.Tables) == 0 {
		return nil
	}

	gqlSchema, err := newGraphQLSchema(schema, s.p)
	if err != nil {
		return fmt.Errorf("falied to build GraphQL schema: %w", err)
	}

	s.mu.Lock()
	s.schema = &gqlSchema
	s.mu.Unlock()

	return nil
}

func (s *Store) currentBubblySchema() (*bubblySchema, error) {
	// Check if the schema has been initialized
	if s.schema == nil {
		// This is fine, as we might be creating the schema in this process,
		// so initialize the schema
		return newBubblySchema(), nil
	}

	schemaQuery := fmt.Sprintf(`
		{
			%s (first: 1){
				tables
			}
		}
	`, core.SchemaTableName)

	result, err := s.Query(schemaQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get current tables: %w", err)
	}
	val := result.(map[string]interface{})[core.SchemaTableName].([]interface{})
	if len(val) == 0 {
		return newBubblySchema(), nil
	}

	var schema bubblySchema
	b, err := json.Marshal(val[0])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal graphql response: %w", err)
	}
	err = json.Unmarshal(b, &schema)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal graphql response: %w", err)
	}
	return &schema, nil
}

func addImplicitJoins(schema *bubblySchema, tables core.Tables, parent *core.Table) {
	for _, t := range tables {
		if parent != nil {
			var hasParentID bool
			// Check if the parent was already added to the schema
			for _, f := range t.Joins {
				if f.Table == parent.Name {
					hasParentID = true
				}
			}
			if !hasParentID {
				t.Joins = append(t.Joins, core.TableJoin{
					Table:  parent.Name,
					Unique: parent.Unique,
				})
			}
		}

		addImplicitJoins(schema, t.Tables, &t)
		// Clear the child tables
		t.Tables = nil
		schema.Tables[t.Name] = t
	}
}

var internalTables = core.Tables{resourceTable, eventTable, schemaTable}
var resourceTable = core.Table{
	Name: core.ResourceTableName,
	Fields: []core.TableField{
		{
			Name:   "id",
			Type:   cty.String,
			Unique: true,
		},
		{
			Name: "name",
			Type: cty.String,
		},
		{
			Name: "kind",
			Type: cty.String,
		},
		{
			Name: "metadata",
			Type: cty.Object(map[string]cty.Type{}),
		},
		{
			Name: "api_version",
			Type: cty.String,
		},
		{
			Name: "spec",
			Type: cty.String,
		},
	},
}

var eventTable = core.Table{
	Name: core.EventTableName,
	Fields: []core.TableField{
		{
			Name: "status",
			Type: cty.String,
		},
		{
			Name: "time",
			Type: cty.String,
		},
	},
	Joins: []core.TableJoin{
		{
			Table: core.ResourceTableName,
		},
	},
}

var schemaTable = core.Table{
	Name: core.SchemaTableName,
	Fields: []core.TableField{
		{
			Name:   "tables",
			Type:   cty.Object(map[string]cty.Type{}),
			Unique: false,
		},
	},
}
