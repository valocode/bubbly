package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

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
	if !isValidValue(result) {
		return newBubblySchema(), nil
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

// TODO delete?
// func (s *Store) GetResource(id string) (io.Reader, error) {
// 	return s.p.GetResource(id)
// }
//
// // TODO delete?
// func (s *Store) PutResource(id string, val string, resBlock core.ResourceBlock) error {
// 	return s.p.PutResource(id, val, resBlock)
// }
//
// // TODO delete?
// func (s *Store) GetResourcesByKind(resourceKind core.ResourceKind) (io.Reader, error) {
// 	// TODO This currently mocks the store behavior until graphql is implemented
// 	resString := "{\"kind\":\"pipeline_run\",\"name\":\"sonarqube\",\"api_version\":\"v1\",\"metadata\":{\"labels\":{\"environment\":\"prod\"},\"namespace\":\"qa\"},\"spec\":\"\\n        // specify the name of the pipeline resource to execute\\n        interval = \\\"22s\\\"\\n        pipeline = \\\"pipeline/sonarqube\\\"\\n        // specify the pipeline input(s) required\\n        input \\\"file\\\" {\\n            value = \\\"./testdata/sonarqube/sonarqube-example.json\\\"\\n        }\\n        input \\\"repo\\\" {\\n            value = \\\"./testdata/git/repo1.git\\\"\\n        }\\n    \"}"
// 	resJSON := core.ResourceBlockJSON{}
// 	err := json.Unmarshal([]byte(resString), &resJSON)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	resourceJson, err := json.Marshal(resJSON)
// 	if err != nil {
// 		return nil, err
// 	}
// 	bytes.NewReader(resourceJson)
//
// 	return bytes.NewReader(resourceJson), nil
// }

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
