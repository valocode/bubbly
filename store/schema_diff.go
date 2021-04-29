package store

import (
	"fmt"

	"github.com/valocode/bubbly/api/core"
)

// compareSchema will take 2 schemas as arguments and return a schemaUpdates
// elements are matched on Name, and if a Name is not matched in the 2nd schema,
// this will be treated as a deletion.
func compareSchema(s1 *bubblySchema, s2 *bubblySchema) (schemaUpdates, error) {
	var changelog schemaUpdates
	for tableName, table1 := range s1.Tables {
		if tableName != table1.Name {
			return nil, fmt.Errorf("map key and table name do not match for table %s", table1.Name)
		}
		_, ok := s2.Tables[table1.Name]
		// The key exists in both tables, and will be checked for updates
		if ok {
			// ALTER
			calculateDiff(table1, s2.Tables[table1.Name], &changelog)
		} else {
			// DELETE
			// The key does not exist in the second table and will be removed
			changelog = append(changelog, changeEntry{
				Action: remove,
				TableInfo: tableInfo{
					TableName:   table1.Name,
					ElementName: table1.Name,
					ElementType: tableElement,
				},
				From: table1,
				To:   nil,
			})
		}
	}
	for _, table2 := range s2.Tables {
		_, ok := s1.Tables[table2.Name]
		if !ok {
			// CREATE
			// The key exist in only the second table and will be created
			changelog = append(changelog, changeEntry{
				Action: create,
				TableInfo: tableInfo{
					TableName:   table2.Name,
					ElementName: table2.Name,
					ElementType: tableElement,
				},
				From: nil,
				To:   table2,
			})
		}
	}
	return changelog, nil
}

type DiffAction string
type Element string

var (
	update DiffAction = "update"
	remove DiffAction = "delete"
	create DiffAction = "create"

	tableElement    Element = "table"
	fieldElement    Element = "field"
	fieldType       Element = "fieldType"
	fieldUniqueAttr Element = "fieldUnique"
	joinElement     Element = "join"
	joinSingleAttr  Element = "joinSingle"
	joinUniqueAttr  Element = "joinUnique"
)

// schemaUpdates is a list of expectedChanges that will be applied by the migration
type schemaUpdates []changeEntry

type tableInfo struct {
	TableName   string
	ElementName string
	ElementType Element
}

type changeEntry struct {
	Action    DiffAction
	TableInfo tableInfo
	From      interface{}
	To        interface{}
}

// calculateDiff will calculate the difference between two schemas.
// In this case, all elements will be matched on id, if 2 ids are different, they will
// be treated as separate elements. For example:
// field1: "hello world" -> field1: "lizards"
// will be views as an update on field1, but if field1 has its name changed:
// field1: "hello world" -> field2: "hello world"
// These will be treated as 2 separate entities, field1 will be seen as deleted, and field2 will be added
func calculateDiff(t1 core.Table, t2 core.Table, cl *schemaUpdates) {
	compareFields(t1, t2, cl)
	compareJoins(t1, t2, cl)
}

func compareFields(t1, t2 core.Table, cl *schemaUpdates) {
	for _, field1 := range t1.Fields {
		found := false
		for _, field2 := range t2.Fields {
			if field1.Name != field2.Name {
				continue
			}
			found = true
			// Check if same field but different type
			if !field1.Type.Equals(field2.Type) {
				*cl = append(*cl, changeEntry{
					Action: update,
					TableInfo: tableInfo{
						TableName:   t2.Name,
						ElementName: field2.Name,
						ElementType: fieldType,
					},
					From: field1.Type,
					To:   field2.Type,
				})
			}
			// Check if same field but "Unique" has changed
			if field1.Unique != field2.Unique {
				*cl = append(*cl, changeEntry{
					Action: update,
					TableInfo: tableInfo{
						TableName:   t2.Name,
						ElementName: field2.Name,
						ElementType: fieldUniqueAttr,
					},
					From: field1.Unique,
					To:   field2.Unique,
				})
			}
		}
		if !found {
			*cl = append(*cl, changeEntry{
				Action: remove,
				TableInfo: tableInfo{
					TableName:   t1.Name,
					ElementName: field1.Name,
					ElementType: fieldElement,
				},
				From: field1,
				To:   nil,
			})
		}
	}
	for _, schema2Field := range t2.Fields {
		found := false
		for _, schema1Field := range t1.Fields {
			if schema2Field.Name == schema1Field.Name {
				found = true
				break
			}
		}
		if !found {
			*cl = append(*cl, changeEntry{
				Action: create,
				TableInfo: tableInfo{
					TableName:   t2.Name,
					ElementName: schema2Field.Name,
					ElementType: fieldElement,
				},
				From: nil,
				To:   schema2Field,
			})
		}
	}
}

// compareJoins takes two tables and adds any differences in the joins to schemaUpdates
func compareJoins(t1, t2 core.Table, cl *schemaUpdates) {
	for _, join1 := range t1.Joins {
		found := false
		for _, join2 := range t2.Joins {
			// Check whether the join's match by name. If not, continue to the
			// next join
			if join1.Table != join2.Table {
				continue
			}
			found = true
			if join1.Single != join2.Single {
				*cl = append(*cl, changeEntry{
					Action: update,
					TableInfo: tableInfo{
						TableName:   t2.Name,
						ElementName: join2.Table,
						ElementType: joinSingleAttr,
					},
					From: join1.Single,
					To:   join2.Single,
				})
			}
			if join1.Unique != join2.Unique {
				*cl = append(*cl, changeEntry{
					Action: update,
					TableInfo: tableInfo{
						TableName:   t2.Name,
						ElementName: join2.Table,
						ElementType: joinUniqueAttr,
					},
					From: join1.Unique,
					To:   join2.Unique,
				})
			}
		}
		// If not found, the join in t1 has been REMOVED
		if !found {
			*cl = append(*cl, changeEntry{
				Action: remove,
				TableInfo: tableInfo{
					TableName:   t1.Name,
					ElementName: join1.Table,
					ElementType: joinElement,
				},
				From: join1,
				To:   nil,
			})
		}
	}

	// Find the joins from t2 that have been CREATED
	for _, join2 := range t2.Joins {
		found := false
		for _, join1 := range t1.Joins {
			if join2.Table == join1.Table {
				found = true
				break
			}
		}
		if !found {
			*cl = append(*cl, changeEntry{
				Action: create,
				TableInfo: tableInfo{
					TableName:   t2.Name,
					ElementName: join2.Table,
					ElementType: joinElement,
				},
				From: nil,
				To:   join2,
			})
		}
	}
}
