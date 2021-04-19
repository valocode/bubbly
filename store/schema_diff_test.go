package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/api/core"
)

var schemaDiffTests = []struct {
	name    string
	s1      core.Tables
	s2      core.Tables
	want    schemaUpdates
	wantErr bool
}{
	{
		name:    "Test 2 equal schemas",
		s1:      schema1,
		s2:      schema1,
		want:    nil,
		wantErr: false,
	},
	{
		name: "Test create table",
		s1:   core.Tables{},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: create, TableInfo: tableInfo{TableName: "a", ElementName: "a", ElementType: "table"}, From: nil, To: core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		},
		wantErr: false,
	},
	{
		name: "Test delete table",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		s2:   core.Tables{},
		want: schemaUpdates{
			changeEntry{Action: remove, TableInfo: tableInfo{TableName: "a", ElementName: "a", ElementType: "table"}, From: core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}, To: nil},
		},
		wantErr: false,
	},
	{
		name: "Test add field",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}, {Name: "b", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: create, TableInfo: tableInfo{TableName: "a", ElementName: "b", ElementType: "field"}, From: nil, To: core.TableField{Name: "b", Type: cty.String}},
		},
		wantErr: false,
	},
	{
		name: "Test update field type",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.Number}}}},
		want: schemaUpdates{
			changeEntry{Action: update, TableInfo: tableInfo{TableName: "a", ElementName: "a", ElementType: "field"}, From: cty.String, To: cty.Number},
		},
		wantErr: false,
	},
	{
		name: "Test update field set unique",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String, Unique: true}}}},
		want: schemaUpdates{
			changeEntry{Action: update, TableInfo: tableInfo{TableName: "a", ElementName: "a", ElementType: "unique"}, From: false, To: true},
		},
		wantErr: false,
	},
	{
		name: "Test update field set not unique",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String, Unique: true}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: update, TableInfo: tableInfo{TableName: "a", ElementName: "a", ElementType: "unique"}, From: true, To: false},
		},
		wantErr: false,
	},
	{
		name: "Test delete field",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}, {Name: "b", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: remove, TableInfo: tableInfo{TableName: "a", ElementName: "b", ElementType: "field"}, From: core.TableField{Name: "b", Type: cty.String}, To: nil},
		},
		wantErr: false,
	},
	{
		name: "Test create table with implicit join",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}, Tables: []core.Table{{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}}}},
		want: schemaUpdates{
			changeEntry{Action: create, TableInfo: tableInfo{TableName: "b", ElementName: "b", ElementType: "table"}, From: nil, To: core.Table{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}, Joins: []core.TableJoin{{Table: "a"}}}},
		},
		wantErr: false,
	},
	{
		name: "Test create implicit join",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}, core.Table{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}, Tables: []core.Table{{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}}}},
		want: schemaUpdates{
			changeEntry{Action: create, TableInfo: tableInfo{TableName: "b", ElementName: "a", ElementType: "join"}, From: nil, To: core.TableJoin{Table: "a"}},
		},
		wantErr: false,
	},
	{
		name: "Test create implicit single join",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}, core.Table{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}, Tables: []core.Table{{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}, Unique: true}}}},
		want: schemaUpdates{
			changeEntry{Action: create, TableInfo: tableInfo{TableName: "b", ElementName: "a", ElementType: "join"}, From: nil, To: core.TableJoin{Table: "a", Single: true}},
		},
		wantErr: false,
	},
	{
		name: "Test remove join",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}, Tables: []core.Table{{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}, core.Table{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: remove, TableInfo: tableInfo{TableName: "b", ElementName: "a", ElementType: "join"}, From: core.TableJoin{Table: "a"}, To: nil},
		},
		wantErr: false,
	},
	{
		name: "Test remove table with join",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}, Tables: []core.Table{{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: remove, TableInfo: tableInfo{TableName: "b", ElementName: "b", ElementType: "table"}, From: core.Table{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}, Joins: []core.TableJoin{{Table: "a"}}}, To: nil},
		},
		wantErr: false,
	},
}

func TestCompareSchema(t *testing.T) {
	for _, tt := range schemaDiffTests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := newBubblySchemaFromTables(tt.s1)
			s2 := newBubblySchemaFromTables(tt.s2)
			changes, err := compareSchema(s1, s2)
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, changes)
		})
	}
}

var schema1 = core.Tables{
	core.Table{
		Name: "table1",
		Fields: []core.TableField{
			{
				Name: "field1",
				Type: cty.String,
			},
			{
				Name: "field_delete",
				Type: cty.Number,
			},
		},
		Joins: []core.TableJoin{
			{
				Table:  "table2",
				Single: true,
			},
		},
		Tables: []core.Table{
			{
				Name:   "tables_1",
				Unique: false,
				Fields: []core.TableField{
					{
						Name:   "field1",
						Type:   cty.String,
						Unique: true,
					},
					{
						Name:   "field_delete",
						Type:   cty.String,
						Unique: false,
					},
				},
				Tables: []core.Table{
					{
						Name: "sub_table_1",
						Fields: []core.TableField{
							{
								Name:   "field1",
								Type:   cty.String,
								Unique: true,
							},
						},
						Tables: nil,
						Unique: false,
					},
				},
			},
			{
				Name: "table_to_remove",
				Fields: []core.TableField{
					{
						Name:   "field1",
						Type:   cty.String,
						Unique: true,
					},
				},
				Tables: nil,
				Unique: false,
			},
		},
		Unique: true,
	},
	core.Table{
		Name: "table2",
		Fields: []core.TableField{
			{
				Name: "field2",
				Type: cty.String,
			},
		},
		Unique: false,
	},
}
