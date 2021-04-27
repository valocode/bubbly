package store

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/test"
)

const (
	postgresDatabase = "bubbly"
	postgresUser     = "postgres"
	postgresPassword = "secret"
)

func TestApplyMigrationSchemaPostgres(t *testing.T) {
	// Set up complete. Now for the test:
	bCtx := env.NewBubblyContext()
	for _, tt := range schemaDiffTests {
		t.Run(tt.name, func(t *testing.T) {
			// Start postgres in docker
			resource := test.RunPostgresDocker(bCtx, t)
			bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
			// // Create the bubbly schemas
			s1 := newBubblySchemaFromTables(tt.s1)
			s2 := newBubblySchemaFromTables(tt.s2)

			// Get a handle to the bubbly store
			s, err := New(bCtx)
			assert.NoErrorf(t, err, "failed to initialize store")

			// Apply the initial schema
			err = s.Apply(DefaultTenantName, tt.s1)
			assert.NoError(t, err)
			curSchema, err := s.currentBubblySchema(DefaultTenantName)
			assert.NoError(t, err)
			// Verify that no changes are needed
			curChanges, err := compareSchema(curSchema, s1)
			assert.NoError(t, err)
			assert.Empty(t, curChanges)

			// Apply the second schema, which triggers a migration
			err = s.Apply(DefaultTenantName, tt.s2)
			assert.NoError(t, err)
			newSchema, err := s.currentBubblySchema(DefaultTenantName)
			assert.NoError(t, err)
			// Verify that no changes are needed
			newChanges, err := compareSchema(newSchema, s2)
			assert.NoError(t, err)
			assert.Empty(t, newChanges)

			sqlStr := "select table_name,column_name,data_type from information_schema.columns where table_schema = '" + psqlSchemaName(DefaultTenantName) + "';"
			postgres, ok := s.p.(*postgres)
			require.Truef(t, ok, "error casting store provider to postgres")
			rows, err := postgres.pool.Query(context.TODO(), sqlStr)
			require.NoError(t, err)
			defer rows.Close()

			type sqlColumn struct {
				table      string
				column     string
				columnType string
			}

			var sqlColumns = make([]sqlColumn, 0)
			for rows.Next() {
				var tableName, columnName, columnType string
				err := rows.Scan(&tableName, &columnName, &columnType)
				require.NoError(t, err)

				sqlColumns = append(sqlColumns, sqlColumn{table: tableName, column: columnName, columnType: columnType})
				table, ok := s2.Tables[tableName]
				assert.Truef(t, ok, "table %s exists postgres but not in the bubbly schema after migration", tableName)
				// Continue onto next in for loop if table not found
				if !ok {
					continue
				}
				var foundField bool
				for _, field := range table.Fields {
					if field.Name == columnName {
						foundField = true
						break
					}
				}
				// Check if _id field
				if !foundField && columnName == tableIDField {
					foundField = true
				}
				// Check if foreign key field because of join
				if !foundField && strings.HasSuffix(columnName, tableJoinSuffix) {
					joinTableName := columnName[:len(columnName)-len(tableJoinSuffix)]
					for _, join := range table.Joins {
						if join.Table == joinTableName {
							foundField = true
							break
						}
					}
				}
				assert.Truef(t, foundField, "table %s has column %s in postgres which does not exist in the bubbly schema after migration", tableName, columnName)
			}

			// Now iterate over the schema and check that everything exists in postgres
			for _, table := range s2.Tables {
				var foundTable bool
				for _, col := range sqlColumns {
					if col.table == table.Name {
						foundTable = true
						break
					}
				}
				assert.Truef(t, foundTable, "table %s exists in the bubbly schema but not in postgres after the migration", table.Name)
				if !foundTable {
					continue
				}

				// Check the fields
				for _, field := range table.Fields {
					var foundField bool
					for _, col := range sqlColumns {
						if col.table == table.Name && col.column == field.Name {
							foundField = true
							break
						}
					}
					assert.Truef(t, foundField, "field %s in table %s exists in the bubbly schema but not in postgres after the migration", field.Name, table.Name)
				}
				// Check the join
				for _, join := range table.Joins {
					var foundJoin bool
					for _, col := range sqlColumns {
						if col.table == table.Name && col.column == join.Table+tableJoinSuffix {
							foundJoin = true
							break
						}
					}
					assert.Truef(t, foundJoin, "join to table %s in table %s exists in the bubbly schema but not in postgres after the migration", join.Table, table.Name)
				}
			}
		})
	}
}
