package store

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/zclconf/go-cty/cty"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
)

// New creates a new Store for the given config.
func New(bCtx *env.BubblyContext) (*Store, error) {
	var (
		s   = &Store{}
		p   provider
		err error
	)

	for attempt := 1; attempt <= bCtx.StoreConfig.RetryAttempts; attempt++ {
		switch bCtx.StoreConfig.Provider {
		case config.PostgresStore:
			p, err = newPostgres(bCtx)
		case config.CockroachDBStore:
			p, err = newCockroachdb(bCtx)
		default:
			return nil, fmt.Errorf("invalid provider: %s", bCtx.StoreConfig.Provider)
		}
		// If the connection succeeded then break out of the attempt loop
		if err == nil {
			break
		}
		bCtx.Logger.Warn().Err(err).Msgf("Store connection attempt %d failed. %d attempts left", attempt, bCtx.StoreConfig.RetryAttempts-attempt)

		// Sleep for the specified amount of time
		time.Sleep(time.Second * time.Duration(bCtx.StoreConfig.RetrySleep))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to provider: %s: %w", bCtx.StoreConfig.Provider, err)
	}
	s.p = p

	// Check if the provider database has the schema table.
	// If not, it indicates a fresh database that should be initialized
	schemaExists, err := s.p.HasTable(schemaTable)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing schema table: %w", err)
	}

	if !schemaExists {
		// If the schema table does not exist yet we should create it
		if err := s.p.Apply(newBubblySchema()); err != nil {
			return nil, fmt.Errorf("failed to initialize the provider database with internal tables")
		}
	}

	if err := s.syncSchema(); err != nil {
		return nil, fmt.Errorf("failed to sync schema: %w", err)
	}
	return s, nil
}

// Store provides access to persisted readiness data.
type Store struct {
	p provider

	mu           sync.RWMutex
	graph        *schemaGraph
	bubblySchema *bubblySchema
	schema       *graphql.Schema
}

// Schema gets the graphql schema for the store.
func (s *Store) Schema() *graphql.Schema {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.schema
}

// Query queries the store.
func (s *Store) Query(query string) *graphql.Result {

	s.mu.RLock()
	defer s.mu.RUnlock()

	return graphql.Do(graphql.Params{
		Schema:        *s.schema,
		RequestString: query,
	})
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
	if len(cl) == 0 {
		return nil
	}

	if err := s.p.Apply(currentSchema); err != nil {
		return fmt.Errorf("failed to apply schema in provider: %w", err)
	}

	if err := s.updateSchema(currentSchema); err != nil {
		return fmt.Errorf("failed to sync schema: %w", err)
	}

	return nil
}

// Save saves data into the store.
func (s *Store) Save(data core.DataBlocks) error {

	dataTree, err := createDataTree(data)
	if err != nil {
		return fmt.Errorf("failed to create tree of data blocks for storing: %w", err)
	}
	if err := s.p.Save(s.bubblySchema, dataTree); err != nil {
		return fmt.Errorf("falied to save data in provider: %w", err)
	}

	return nil
}

func (s *Store) syncSchema() error {

	bubblySchema, err := s.currentBubblySchema()
	if err != nil {
		return fmt.Errorf("failed to get current schema: %w", err)
	}
	s.updateSchema(bubblySchema)

	return nil
}

func (s *Store) updateSchema(bubblySchema *bubblySchema) error {
	graph, err := newSchemaGraphFromMap(bubblySchema.Tables)
	if err != nil {
		return fmt.Errorf("failed to build schema graph: %w", err)
	}

	gqlSchema, err := newGraphQLSchema(graph, s)
	if err != nil {
		return fmt.Errorf("falied to build GraphQL schema: %w", err)
	}

	s.mu.Lock()
	s.graph = graph
	s.bubblySchema = bubblySchema
	s.schema = &gqlSchema
	s.mu.Unlock()

	return nil
}

func (s *Store) resolveQuery(params graphql.ResolveParams) (interface{}, error) {
	return s.p.ResolveQuery(s.graph, params)
}

func (s *Store) currentBubblySchema() (*bubblySchema, error) {
	// Check if the schema has been initialized
	if s.schema == nil {
		// This is fine, as we might be creating the schema in this process,
		// so initialize the schema and set the graphql schema
		if err := s.updateSchema(newBubblySchema()); err != nil {
			return nil, fmt.Errorf("failed to initialize graphql schema with internal tables: %w", err)
		}
	}

	result := s.Query(schemaQuery)
	if result.HasErrors() {
		return nil, fmt.Errorf("failed to get current tables: %v", result.Errors)
	}
	if !isValidValue(result.Data) {
		return newBubblySchema(), nil
	}
	val := result.Data.(map[string]interface{})[core.SchemaTableName].([]interface{})
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

const tableIDField = "_id"
const tableJoinSuffix = "_id"

var schemaQuery = fmt.Sprintf(`
{
	%s {
		tables
	}
}
`, core.SchemaTableName)

var internalTables = core.Tables{resourceTable, schemaTable, eventTable}
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
