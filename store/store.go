package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cornelk/hashmap"
	"github.com/graphql-go/graphql"
	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

const DefaultTenantName = "default"
//
// The Bubbly Store is an abstraction for structured data stored in Bubbly,
// as well as the metadata describing it.
//
// The Bubbly Store brings together the Bubbly Schema, the Bubbly Schema Graph,
// the GraphQL Schema, and the underlying RDBMS. The latter is known as a "provider".
//

// New creates a new Store for the given config.
func New(bCtx *env.BubblyContext) (*Store, error) {
	var (
		s = &Store{
			bCtx:     bCtx,
			graphs:   &hashmap.HashMap{},
			schemas:  &hashmap.HashMap{},
			triggers: &hashmap.HashMap{},
		}
		err error
	)

	// Connect to the provider's database RetryAttempts times, with a RetrySleep
	for attempt := 1; attempt <= bCtx.StoreConfig.RetryAttempts; attempt++ {
		switch bCtx.StoreConfig.Provider {
		case config.PostgresStore:
			s.p, err = newPostgres(bCtx)
		case config.CockroachDBStore:
			s.p, err = newCockroachdb(bCtx)
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

	if err := s.initStoreSchemas(); err != nil {
		return nil, fmt.Errorf("failed to initialize the store schemas: %w", err)
	}

	return s, nil
}

// Store provides access to persisted readiness data.
type Store struct {
	bCtx *env.BubblyContext
	p    provider

	graphs   *hashmap.HashMap
	schemas  *hashmap.HashMap
	triggers *hashmap.HashMap
}

// CreateTenant creates a tenant schema in the provider
func (s *Store) CreateTenant(tenant string) error {
	if err := s.p.CreateTenant(tenant); err != nil {
		return fmt.Errorf("error creating tenant %s: %w", tenant, err)
	}
	// We should check that a schema already exists, and if not, we should
	// initialize one
	ok, err := s.p.HasTable(tenant, schemaTable)
	if err != nil {
		return fmt.Errorf("error checking if provider has schema for tenant %s: %w", tenant, err)
	}
	if !ok {
		// Initialize the internal schema for the tenant
		if err := s.p.Apply(tenant, newBubblySchema()); err != nil {
			return fmt.Errorf("failed to initialize schema with internal tables for tenant %s: %w", tenant, err)
		}
	}
	return nil
}

// Query queries the store.
func (s *Store) Query(tenant string, query string) (*graphql.Result, error) {
	schema, ok := s.schemas.GetStringKey(tenant)
	if !ok {
		return nil, fmt.Errorf("no schema exists for tenant %s", tenant)
	}
	return graphql.Do(graphql.Params{
		Schema:        schema.(graphql.Schema),
		RequestString: query,
	}), nil
}

// Apply applies a schema corresponding to a set of tables.
func (s *Store) Apply(tenant string, tables core.Tables) error {
	currentSchema, err := s.currentBubblySchema(tenant)
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
	addImplicitJoins(newSchema, tables, nil)

	// Calculate the schema diff
	cl, err := compareSchema(*currentSchema, *newSchema)
	if err != nil {
		return fmt.Errorf("failed to compare schemas: %w", err)
	}
	newSchema.changelog = cl

	// Perform the migration based on the schemaUpdates
	if err := s.p.Migrate(tenant, newSchema, cl); err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}

	// Update the store cache
	if err := s.updateSchema(tenant, newSchema); err != nil {
		return fmt.Errorf("failed to sync schema: %w", err)
	}

	return nil
}

// Save saves data into the store.
func (s *Store) Save(tenant string, data core.DataBlocks) error {
	var (
		graph    *schemaGraph
		triggers []*trigger
	)

	dataTree, err := createDataTree(data)
	if err != nil {
		return fmt.Errorf("failed to create tree of data blocks for storing: %w", err)
	}
	graphVal, ok := s.graphs.GetStringKey(tenant)
	if !ok {
		return fmt.Errorf("no schema exists for tenant %s", tenant)
	}
	graph = graphVal.(*schemaGraph)

	triggersVal, ok := s.triggers.GetStringKey(tenant)
	if !ok {
		return fmt.Errorf("no triggers exist for tenant %s", tenant)
	}
	triggers = triggersVal.([]*trigger)

	if err := s.p.Save(s.bCtx, tenant, graph, dataTree); err != nil {
		return fmt.Errorf("falied to save data in provider: %w", err)
	}

	triggersTree, err := HandleTriggers(s.bCtx, dataTree, triggers, Active)
	if err != nil {
		return fmt.Errorf("data triggers failed: %w", err)
	}

	if err := s.p.Save(s.bCtx, tenant, graph, triggersTree); err != nil {
		return fmt.Errorf("falied to save data in provider: %w", err)
	}

	_, err = HandleTriggers(s.bCtx, dataTree, triggers, Passive)
	if err != nil {
		return fmt.Errorf("passive triggers failed: %w", err)
	}

	return nil
}

// Close closes the connection to the store's own database and the provider
func (s *Store) Close() {
	// Close the provider's connection
	s.p.Close()
}

func (s *Store) initStoreSchemas() error {
	var tenants []string
	// If multitenancy is enabled, fetch the tenants from the store
	if s.bCtx.AuthConfig.MultiTenancy {
		var err error
		tenants, err = s.p.Tenants()
		if err != nil {
			return fmt.Errorf("failed to get tenants from provider: %w", err)
		}
	} else {
		tenants = []string{DefaultTenantName}
		// Make sure we have created the default tenant, as this will not get
		// called otherwise as there is no idea of creating a tenant in single
		// tenancy mode
		if err := s.CreateTenant(DefaultTenantName); err != nil {
			return fmt.Errorf("failed creating default tenant %s: %w", DefaultTenantName, err)
		}
	}
	// Iterate through the tenants and initialize the cache of schemas
	for _, tenant := range tenants {
		// Check if the provider database has the schema table.
		// If not, it indicates a fresh database that should be initialized
		schemaExists, err := s.p.HasTable(tenant, schemaTable)
		if err != nil {
			return fmt.Errorf("failed to check existing schema table: %w", err)
		}
		if !schemaExists {
			// If the schema table does not exist yet we should create it
			if err := s.p.Apply(tenant, newBubblySchema()); err != nil {
				return fmt.Errorf("failed to initialize the provider database with internal tables")
			}
		}

		if err := s.syncSchema(tenant); err != nil {
			return fmt.Errorf("failed to sync schema for tenant %s: %w", tenant, err)
		}
	}

	return nil
}

// syncSchema is used by a store instance to sync it's internally stored schema
// for the specified tenant in the provider's databases
func (s *Store) syncSchema(tenant string) error {
	bubblySchema, err := s.currentBubblySchema(tenant)
	if err != nil {
		return fmt.Errorf("failed to get current schema: %w", err)
	}
	s.updateSchema(tenant, bubblySchema)
	return nil
}

func (s *Store) currentBubblySchema(tenant string) (*bubblySchema, error) {
	var (
		schema graphql.Schema
		err    error
	)
	schemaVal, ok := s.schemas.GetStringKey(tenant)
	if !ok {
		// If there was no schema for this tenant, it might be because we are
		// initialising the store. In order to do this, we need at least some
		// minimum viable graphq schema to query the provider for the existing
		// schema
		graph := internalSchemaGraph()
		schemaVal, err = newGraphQLSchema(graph, func(p graphql.ResolveParams) (interface{}, error) {
			return s.p.ResolveQuery(tenant, graph, p)
		})
		if err != nil {
			return nil, fmt.Errorf("failed creating GraphQL schema of internal tables: %w", err)
		}
	}
	schema = schemaVal.(graphql.Schema)

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: schemaQuery,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get current schema: %w", err)
	}
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

	var bSchema bubblySchema
	b, err := json.Marshal(val[len(val)-1])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal graphql response: %w", err)
	}
	err = json.Unmarshal(b, &bSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal graphql response: %w", err)
	}
	return &bSchema, nil
}

// updateSchema creates a new GraphQL schema from a provided Bubbly Schema,
// and binds that GraphQL schema to the Bubbly Store instance.
func (s *Store) updateSchema(tenant string, bubblySchema *bubblySchema) error {
	graph, err := newSchemaGraphFromMap(bubblySchema.Tables)
	if err != nil {
		return fmt.Errorf("failed to build schema graph: %w", err)
	}
	schema, err := newGraphQLSchema(graph, func(p graphql.ResolveParams) (interface{}, error) {
		return s.p.ResolveQuery(tenant, graph, p)
	})
	if err != nil {
		return fmt.Errorf("failed to create GraphQL schema from graph: %w", err)
	}

	s.graphs.Set(tenant, graph)
	s.schemas.Set(tenant, schema)
	s.triggers.Set(tenant, internalTriggers)
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

// internalTables is the minimum set of tables in a valid Bubbly Schema.
// This set of tables is created when a new Bubbly Schema is created.
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
