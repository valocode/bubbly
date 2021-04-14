package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/api/core"
)

func TestCompareSchema(t *testing.T) {
	type args struct {
		s1 bubblySchema
		s2 bubblySchema
	}
	tests := []struct {
		name    string
		args    args
		want    schemaUpdates
		wantErr bool
	}{
		{
			name: "Test 2 equal schemas",
			args: args{
				s1: schema1,
				s2: schema1,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test 2 unequal schemas",
			args: args{
				s1: schema1,
				s2: schema2,
			},
			want:    expectedChanges,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := compareSchema(tt.args.s1, tt.args.s2)
			require.NoError(t, gotErr)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

var schema1 = bubblySchema{Tables: map[string]core.Table{
	"table1": {
		Name: "table1",
		Fields: []core.TableField{
			{
				Name: "field1",
				Type: cty.String,
			},
			{
				Name: "fieldToBeDeleted",
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
						Name:   "fieldToBeDeleted",
						Type:   cty.String,
						Unique: false,
					},
				},
				Joins: []core.TableJoin{
					{
						Table:  "table1",
						Single: true,
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
				Name: "tableToBeRemoved",
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
	"table2": {
		Name: "table2",
		Fields: []core.TableField{
			{
				Name: "field2",
				Type: cty.String,
			},
		},
		Unique: false,
	},
}}
var schema2 = bubblySchema{Tables: map[string]core.Table{
	"table1": {
		Name: "table1",
		Fields: []core.TableField{
			{
				Name: "field1",
				Type: cty.Number,
			},
			{
				Name:   "field11",
				Type:   cty.Number,
				Unique: true,
			},
		},
		Joins: []core.TableJoin{
			{
				Table:  "table2",
				Single: true,
			},
		},
		Unique: true,
		Tables: []core.Table{
			{
				Name:   "tables_1",
				Unique: true,
				Joins: []core.TableJoin{
					{
						Table:  "tableJoin",
						Single: false,
					},
				},
				Fields: []core.TableField{
					{
						Name:   "field1",
						Type:   cty.Number,
						Unique: false,
					},
				},
				Tables: []core.Table{
					{
						Name:   "sub_table_1",
						Fields: nil,
						Tables: nil,
						Unique: true,
					},
				},
			},
			{
				Name: "tables_new",
				Tables: []core.Table{
					{
						Name: "tables_tables_1",
					},
				},
			},
		},
	},
	"table2": {
		Name: "table2",
		Fields: []core.TableField{
			{
				Name: "field22",
				Type: cty.String,
			},
		},
		Unique: true,
	},
	"table3": {
		Name: "table3",
		Fields: []core.TableField{
			{
				Name: "field22",
				Type: cty.String,
			},
		},
		Unique: false,
	},
}}

var expectedChanges = []changeEntry{
	{Action: "update", TableInfo: tableInfo{TableName: "table1", ElementName: "field1", ElementType: "field"}, From: cty.String, To: cty.Number},
	{Action: "delete", TableInfo: tableInfo{TableName: "table1", ElementName: "table1", ElementType: "field"}, From: core.TableField{Name: "fieldToBeDeleted", Unique: false, Type: cty.Number}, To: nil},
	{Action: "create", TableInfo: tableInfo{TableName: "table1", ElementName: "field11", ElementType: "field"}, From: nil, To: core.TableField{Name: "field11", Unique: true, Type: cty.Number}},
	{Action: "update", TableInfo: tableInfo{TableName: "tables_1", ElementName: "field1", ElementType: "field"}, From: cty.String, To: cty.Number},
	{Action: "update", TableInfo: tableInfo{TableName: "tables_1", ElementName: "field1", ElementType: "unique"}, From: true, To: false},
	{Action: "delete", TableInfo: tableInfo{TableName: "tables_1", ElementName: "tables_1", ElementType: "field"}, From: core.TableField{Name: "fieldToBeDeleted", Unique: false, Type: cty.String}, To: nil},
	{Action: "delete", TableInfo: tableInfo{TableName: "tables_1", ElementName: "table1", ElementType: "join"}, From: core.TableJoin{Table: "table1", Single: true}, To: nil},
	{Action: "create", TableInfo: tableInfo{TableName: "tables_1", ElementName: "tableJoin", ElementType: "join"}, From: nil, To: core.TableJoin{Table: "tableJoin", Single: false}},
	{Action: "delete", TableInfo: tableInfo{TableName: "sub_table_1", ElementName: "sub_table_1", ElementType: "field"}, From: core.TableField{Name: "field1", Unique: true, Type: cty.String}, To: nil},
	{Action: "update", TableInfo: tableInfo{TableName: "sub_table_1", ElementName: "name", ElementType: "unique"}, From: false, To: true},
	{Action: "update", TableInfo: tableInfo{TableName: "tables_1", ElementName: "name", ElementType: "unique"}, From: false, To: true},
	{Action: "delete", TableInfo: tableInfo{TableName: "tableToBeRemoved", ElementName: "tableToBeRemoved", ElementType: "table"}, From: core.Table{Name: "tableToBeRemoved", Fields: []core.TableField{{Name: "field1", Unique: true, Type: cty.String}}, Joins: nil, Unique: false, Tables: nil}, To: nil},
	{Action: "create", TableInfo: tableInfo{TableName: "tables_new", ElementName: "tables_new", ElementType: "table"}, From: nil, To: core.Table{Name: "tables_new", Fields: nil, Unique: false, Tables: []core.Table{{Name: "tables_tables_1", Fields: nil, Unique: false, Tables: nil}}}},
	{Action: "delete", TableInfo: tableInfo{TableName: "table2", ElementName: "table2", ElementType: "field"}, From: core.TableField{Name: "field2", Unique: false, Type: cty.String}, To: nil},
	{Action: "create", TableInfo: tableInfo{TableName: "table2", ElementName: "field22", ElementType: "field"}, From: nil, To: core.TableField{Name: "field22", Unique: false, Type: cty.String}},
	{Action: "update", TableInfo: tableInfo{TableName: "table2", ElementName: "name", ElementType: "unique"}, From: false, To: true},
	{Action: "create", TableInfo: tableInfo{TableName: "table3", ElementName: "table3", ElementType: "table"}, From: nil, To: core.Table{Name: "table3", Fields: []core.TableField{{Name: "field22", Unique: false, Type: cty.String}}, Joins: nil, Unique: false, Tables: nil}},
}
