package store

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

// New creates a new Store for the given config.
func New(bCtx *env.BubblyContext) (*Store, error) {
	var (
		s   = &Store{bCtx: bCtx}
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

	s.triggers = internalTriggers
	return s, nil
}

// Store provides access to persisted readiness data.
type Store struct {
	bCtx *env.BubblyContext
	p    provider

	// bubblySchema *bubblySchema
	graph     *schemaGraph
	mu        sync.RWMutex
	gqlSchema *graphql.Schema
	triggers  []*trigger
}

// Schema gets the graphql schema for the store.
func (s *Store) Schema() *graphql.Schema {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.gqlSchema
}

// Query queries the store.
func (s *Store) Query(query string) *graphql.Result {

	s.mu.RLock()
	defer s.mu.RUnlock()

	return graphql.Do(graphql.Params{
		Schema:        *s.gqlSchema,
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
	addImplicitJoins(newSchema, tables, nil)

	// Calculate the schema diff
	cl, err := compareSchema(*currentSchema, *newSchema)
	if err != nil {
		return fmt.Errorf("failed to compare schemas: %w", err)
	}
	newSchema.Changelog = cl

	// perform the schema migration if there is anything in the changelog
	if len(cl) > 0 {
		m, err := s.p.GenerateMigration(s.bCtx, cl)
		if err != nil {
			return err
		}
		err = s.p.Migrate(m)
		if err != nil {
			return err
		}

		if err := s.updateSchema(currentSchema); err != nil {
			return fmt.Errorf("failed to sync schema: %w", err)
		}
	} else {
		// since there is no schema, apply
		if err := s.p.Apply(newSchema); err != nil {
			return fmt.Errorf("failed to apply schema in provider: %w", err)
		}
		if err := s.updateSchema(currentSchema); err != nil {
			return fmt.Errorf("failed to sync schema: %w", err)
		}
	}

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

	if err := s.p.Save(s.bCtx, s.graph, dataTree); err != nil {
		return fmt.Errorf("falied to save data in provider: %w", err)
	}

	triggersTree, err := HandleTriggers(s.bCtx, dataTree, s.triggers, Active)
	if err != nil {
		return fmt.Errorf("data triggers failed: %w", err)
	}

	if err := s.p.Save(s.bCtx, s.graph, triggersTree); err != nil {
		return fmt.Errorf("falied to save data in provider: %w", err)
	}

	_, err = HandleTriggers(s.bCtx, dataTree, s.triggers, Passive)
	if err != nil {
		return fmt.Errorf("passive triggers failed: %w", err)
	}

	return nil
}

// resolveQuery is the handler for resolving graphql queries
func (s *Store) resolveQuery(params graphql.ResolveParams) (interface{}, error) {
	return s.p.ResolveQuery(s.graph, params)
}

// syncSchema is used by a store instance to sync it's internally stored schema
// with the current schema in the provider
func (s *Store) syncSchema() error {
	bubblySchema, err := s.currentBubblySchema()
	if err != nil {
		return fmt.Errorf("failed to get current schema: %w", err)
	}
	s.updateSchema(bubblySchema)

	return nil
}

func (s *Store) currentBubblySchema() (*bubblySchema, error) {
	// Check if the schema has been initialized
	if s.gqlSchema == nil {
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
	b, err := json.Marshal(val[len(val)-1])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal graphql response: %w", err)
	}
	err = json.Unmarshal(b, &schema)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal graphql response: %w", err)
	}
	return &schema, nil
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
	s.gqlSchema = &gqlSchema
	s.mu.Unlock()

	return nil
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
					Single: parent.Unique,
				})
			}
		}

		addImplicitJoins(schema, t.Tables, &t)
		// Clear the child tables
		t.Tables = nil
		schema.Tables[t.Name] = t
	}
}

const (
	tableIDField    = "_id"
	tableJoinSuffix = "_id"
)

// TODO: add "limit" arg to this query
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
			Name: "error",
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
