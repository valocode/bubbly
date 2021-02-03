package store

import (
	"fmt"
	"testing"

	"github.com/r3labs/diff"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zclconf/go-cty/cty"

	"github.com/verifa/bubbly/api/core"
)

func Test_compareSchema(t *testing.T) {
	// Test cases
	// table added *
	// field added *
	// join added *
	// Tables added *
	// table removed *
	// field removed *
	// join removed * list len 10
	// Tables removed *
	// Table updated
	//	Table Unique Updated *
	// field updated
	//	field deleted *
	//	field Unique updated *
	//	field Type updated
	// Join updated
	// 	join unique updated

	// TODO move to file
	// might need to split it up since bubblySchema is not exported
	schema1 := bubblySchema{Tables: map[string]core.Table{
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
					Table:  "tableToBeRemoved",
					Unique: true,
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
							Name: "fieldToBeDeleted",
						},
					},
					Joins: []core.TableJoin{
						{
							Table:  "tableJoin",
							Unique: true,
						},
					},
					Tables: []core.Table{
						{
							Name:   "sub_table_1",
							Unique: false,
						},
					},
				},
				{
					Name:   "tableToBeRemoved",
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
		"tableToBeRemoved": {
			Name: "tableToBeRemoved",
		},
	}}
	schema2 := bubblySchema{Tables: map[string]core.Table{
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
					Unique: true,
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
							Unique: false,
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

	type args struct {
		s1 bubblySchema
		s2 bubblySchema
	}
	tests := []struct {
		name    string
		args    args
		want    Changelog
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
			want: []Entry{
				{Action: "delete", TableName: "tableToBeRemoved", ElementType: "table", From: core.Table{Name: "tableToBeRemoved", Fields: []core.TableField(nil), Joins: []core.TableJoin(nil), Unique: false, Tables: []core.Table(nil)}, To: interface{}(nil)},
				{Action: "update", TableName: "table1", ElementType: "field", From: cty.String, To: cty.Number},
				{Action: "delete", TableName: "table1", ElementType: "field", From: core.TableField{Name: "fieldToBeDeleted", Unique: false, Type: cty.Number}, To: interface{}(nil)},

				{Action: "create", TableName: "table1", ElementType: "field", From: interface{}(nil), To: core.TableField{Name: "field11", Unique: true, Type: cty.Number}},
				{Action: "delete", TableName: "table1", ElementType: "join", From: core.TableJoin{Table: "tableToBeRemoved", Unique: true}, To: interface{}(nil)},
				{Action: "create", TableName: "table1", ElementType: "join", From: interface{}(nil), To: core.TableJoin{Table: "table2", Unique: true}},
				{Action: "update", TableName: "tables_1", ElementType: "field", From: cty.String, To: cty.Number},

				{Action: "update", TableName: "tables_1", ElementType: "unique", From: true, To: false},
				{Action: "delete", TableName: "tables_1", ElementType: "field", From: core.TableField{Name: "fieldToBeDeleted", Unique: false, Type: cty.NilType}, To: interface{}(nil)},
				{Action: "update", TableName: "tables_1", ElementType: "join", From: true, To: false},
				{Action: "update", TableName: "sub_table_1", ElementType: "unique", From: false, To: true},
				{Action: "update", TableName: "tables_1", ElementType: "unique", From: false, To: true},
				{Action: "delete", TableName: "tableToBeRemoved", ElementType: "table", From: core.Table{Name: "tableToBeRemoved", Fields: []core.TableField(nil), Joins: []core.TableJoin(nil), Unique: false, Tables: []core.Table(nil)}, To: interface{}(nil)},
				{Action: "create", TableName: "tables_new", ElementType: "table", From: interface{}(nil), To: core.Table{Name: "tables_new", Fields: []core.TableField(nil), Joins: []core.TableJoin(nil), Unique: false, Tables: []core.Table{core.Table{Name: "tables_tables_1", Fields: []core.TableField(nil), Joins: []core.TableJoin(nil), Unique: false, Tables: []core.Table(nil)}}}},
				{Action: "delete", TableName: "table2", ElementType: "field", From: core.TableField{Name: "field2", Unique: false, Type: cty.String}, To: interface{}(nil)},
				{Action: "create", TableName: "table2", ElementType: "field", From: interface{}(nil), To: core.TableField{Name: "field22", Unique: false, Type: cty.String}},
				{Action: "update", TableName: "table2", ElementType: "unique", From: false, To: true},
				{Action: "create", TableName: "table3", ElementType: "table", From: interface{}(nil), To: core.Table{Name: "table3", Fields: []core.TableField{core.TableField{Name: "field22", Unique: false, Type: cty.String}}, Joins: []core.TableJoin(nil), Unique: false, Tables: []core.Table(nil)}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := compareSchema(tt.args.s1, tt.args.s2)
			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
			// fmt.Println(deep.Equal(tt.want, got))
			cl, _ := diff.Diff(tt.want, got)
			fmt.Println(cl)
		})
	}
}
