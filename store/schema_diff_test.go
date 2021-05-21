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
			changeEntry{Action: create, TableInfo: tableInfo{TableName: "a", ElementName: "a", ElementType: tableElement}, From: nil, To: core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		},
		wantErr: false,
	},
	{
		name: "Test delete table",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		s2:   core.Tables{},
		want: schemaUpdates{
			changeEntry{Action: remove, TableInfo: tableInfo{TableName: "a", ElementName: "a", ElementType: tableElement}, From: core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}, To: nil},
		},
		wantErr: false,
	},
	{
		name: "Test add field",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}, {Name: "b", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: create, TableInfo: tableInfo{TableName: "a", ElementName: "b", ElementType: fieldElement}, From: nil, To: core.TableField{Name: "b", Type: cty.String}},
		},
		wantErr: false,
	},
	{
		name: "Test update field type",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.Number}}}},
		want: schemaUpdates{
			changeEntry{Action: update, TableInfo: tableInfo{TableName: "a", ElementName: "a", ElementType: fieldType}, From: cty.String, To: cty.Number},
		},
		wantErr: false,
	},
	{
		name: "Test update field set unique",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String, Unique: true}}}},
		want: schemaUpdates{
			changeEntry{Action: update, TableInfo: tableInfo{TableName: "a", ElementName: "a", ElementType: fieldUniqueAttr}, From: false, To: true},
		},
		wantErr: false,
	},
	{
		name: "Test update field set not unique",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String, Unique: true}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: update, TableInfo: tableInfo{TableName: "a", ElementName: "a", ElementType: fieldUniqueAttr}, From: true, To: false},
		},
		wantErr: false,
	},
	{
		name: "Test delete field",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}, {Name: "b", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: remove, TableInfo: tableInfo{TableName: "a", ElementName: "b", ElementType: fieldElement}, From: core.TableField{Name: "b", Type: cty.String}, To: nil},
		},
		wantErr: false,
	},
	{
		name: "Test create table with implicit join",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}, Tables: []core.Table{{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}}}},
		want: schemaUpdates{
			changeEntry{Action: create, TableInfo: tableInfo{TableName: "b", ElementName: "b", ElementType: tableElement}, From: nil, To: core.Table{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}, Joins: []core.TableJoin{{Table: "a"}}}},
		},
		wantErr: false,
	},
	{
		name: "Test create implicit join",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}, core.Table{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}, Tables: []core.Table{{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}}}},
		want: schemaUpdates{
			changeEntry{Action: create, TableInfo: tableInfo{TableName: "b", ElementName: "a", ElementType: joinElement}, From: nil, To: core.TableJoin{Table: "a"}},
		},
		wantErr: false,
	},
	{
		name: "Test create implicit single join",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}, core.Table{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}, Tables: []core.Table{{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}, Single: true}}}},
		want: schemaUpdates{
			changeEntry{Action: create, TableInfo: tableInfo{TableName: "b", ElementName: "a", ElementType: joinElement}, From: nil, To: core.TableJoin{Table: "a", Single: true}},
		},
		wantErr: false,
	},
	{
		name: "Test remove join",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}, Tables: []core.Table{{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}, core.Table{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: remove, TableInfo: tableInfo{TableName: "b", ElementName: "a", ElementType: joinElement}, From: core.TableJoin{Table: "a"}, To: nil},
		},
		wantErr: false,
	},
	{
		name: "Test remove table with join",
		s1:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}, Tables: []core.Table{{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}}}}},
		s2:   core.Tables{core.Table{Name: "a", Fields: []core.TableField{{Name: "a", Type: cty.String}}}},
		want: schemaUpdates{
			changeEntry{Action: remove, TableInfo: tableInfo{TableName: "b", ElementName: "b", ElementType: tableElement}, From: core.Table{Name: "b", Fields: []core.TableField{{Name: "b", Type: cty.String}}, Joins: []core.TableJoin{{Table: "a"}}}, To: nil},
		},
		wantErr: false,
	},
	{
		name: "Add unique constraint on join",
		s1:   core.Tables{core.Table{Name: "a", Tables: []core.Table{{Name: "b"}}}},
		s2:   core.Tables{core.Table{Name: "a", Tables: []core.Table{{Name: "b", Unique: true}}}},
		want: schemaUpdates{
			changeEntry{Action: update, TableInfo: tableInfo{TableName: "b", ElementName: "a", ElementType: joinUniqueAttr}, From: false, To: true},
		},
		wantErr: false,
	},
	{
		name: "Remove unique constraint on join",
		s1:   core.Tables{core.Table{Name: "a", Tables: []core.Table{{Name: "b", Unique: true}}}},
		s2:   core.Tables{core.Table{Name: "a", Tables: []core.Table{{Name: "b"}}}},
		want: schemaUpdates{
			changeEntry{Action: update, TableInfo: tableInfo{TableName: "b", ElementName: "a", ElementType: joinUniqueAttr}, From: true, To: false},
		},
		wantErr: false,
	},
	{
		name: "Add single attribute on join",
		s1:   core.Tables{core.Table{Name: "a", Tables: []core.Table{{Name: "b"}}}},
		s2:   core.Tables{core.Table{Name: "a", Tables: []core.Table{{Name: "b", Single: true}}}},
		want: schemaUpdates{
			changeEntry{Action: update, TableInfo: tableInfo{TableName: "b", ElementName: "a", ElementType: joinSingleAttr}, From: false, To: true},
		},
		wantErr: false,
	},
	{
		name: "Remove single attribute on join",
		s1:   core.Tables{core.Table{Name: "a", Tables: []core.Table{{Name: "b", Single: true}}}},
		s2:   core.Tables{core.Table{Name: "a", Tables: []core.Table{{Name: "b"}}}},
		want: schemaUpdates{
			changeEntry{Action: update, TableInfo: tableInfo{TableName: "b", ElementName: "a", ElementType: joinSingleAttr}, From: true, To: false},
		},
		wantErr: false,
	},
}

func TestCompareSchema(t *testing.T) {
	for _, tt := range schemaDiffTests {
		t.Run(tt.name, func(t *testing.T) {
			s1, err := newBubblySchemaFromTables(tt.s1, true)
			require.NoError(t, err)
			s2, err := newBubblySchemaFromTables(tt.s2, true)
			require.NoError(t, err)
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
				Single: false,
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
						Single: false,
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
				Single: false,
			},
		},
		Single: true,
	},
	core.Table{
		Name: "table2",
		Fields: []core.TableField{
			{
				Name: "field2",
				Type: cty.String,
			},
		},
		Single: false,
	},
}
