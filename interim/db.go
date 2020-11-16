package interim

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/verifa/bubbly/api/core"

	"github.com/graphql-go/graphql"
	"github.com/hashicorp/go-memdb"
)

// DB provides access to data by means of a query.
type DB struct {
	memDB       *memdb.MemDB
	graphQL     graphql.Schema
	schemaTypes map[string]schemaType
}

// Query queries the DB.
func (db *DB) Query(query string) (interface{}, error) {
	res := graphql.Do(graphql.Params{
		Schema:        db.graphQL,
		RequestString: query,
	})

	if res.HasErrors() {
		return nil, fmt.Errorf("failed to execute query: %v", res.Errors)
	}

	return res.Data, nil
}

// Import imports transform data into the DB.
func (db *DB) Import(data core.DataBlocks) error {
	txn := db.memDB.Txn(true)

	if err := db.insert(data, txn); err != nil {
		return fmt.Errorf("failed to import: %w", err)
	}

	txn.Commit()

	return nil
}

func (db *DB) insert(data core.DataBlocks, txn *memdb.Txn) error {
	for _, d := range data {
		// Retrieve the schema type for the data
		// we are trying to insert.
		st, ok := db.schemaTypes[d.TableName]
		if !ok {
			return fmt.Errorf("unkown type: %s", d.TableName)
		}

		// Create a new instance of a predefined struct
		// that corresponds to both the data and the schema.
		n, err := st.New(d, uuid.New().String(), uuid.New().String())
		if err != nil {
			return fmt.Errorf("falied to create instace of %s: %w", d.TableName, err)
		}

		// Insert the data into membd
		if err := txn.Insert(d.TableName, n); err != nil {
			return fmt.Errorf("falied to insert %s: %w", d.TableName, err)
		}

		// Recursively insert all sub-data.
		return db.insert(d.Data, txn)
	}

	return nil
}

// NewDB creates a new DB for the given tables.
func NewDB(tables []core.Table) (*DB, error) {
	memDB, err := newMemDB(tables)
	if err != nil {
		return nil, fmt.Errorf("failed to create memDB: %w", err)
	}
	graphQL, err := newGraphQL(tables, memDB)
	if err != nil {
		return nil, fmt.Errorf("failed to create graphQL: %w", err)
	}

	return &DB{
		memDB:       memDB,
		graphQL:     graphQL,
		schemaTypes: newSchemaTypes(tables),
	}, nil
}

// // core.Table is a schema table. It may
// // contain fields, tables, or any
// // combination of the two.
// type core.Table struct {
// 	Name   string
// 	Fields []core.TableField
// 	Tables []core.Table
// }

// // core.TableField is a schema field.
// type core.TableField struct {
// 	Name   string
// 	Unique bool
// 	Type   cty.Type
// }
