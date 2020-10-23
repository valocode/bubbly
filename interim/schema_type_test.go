package interim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func TestSchemaType(t *testing.T) {
	type tortoise struct {
		ID    string
		FK    string
		Name  string
		Age   int
		Happy bool
	}

	cases := []struct {
		desc string
		id   string
		fk   string
		t    Table
		d    core.Data
		out  interface{}
	}{
		{
			desc: "basic",
			id:   "abc123",
			fk:   "def456",
			t: Table{
				Name: "tortoise",
				Fields: []Field{
					{
						Name: "Name",
						Type: cty.String,
					},
					{
						Name: "Age",
						Type: cty.Number,
					},
					{
						Name: "Happy",
						Type: cty.Bool,
					},
				},
			},
			d: core.Data{
				Name: "tortoise",
				Fields: []core.Field{
					{
						Name:  "Name",
						Value: cty.StringVal("Harold"),
					},
					{
						Name:  "Age",
						Value: cty.NumberIntVal(111),
					},
					{
						Name:  "Happy",
						Value: cty.BoolVal(true),
					},
				},
			},
			out: tortoise{
				ID:    "abc123",
				FK:    "def456",
				Name:  "Harold",
				Age:   111,
				Happy: true,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			st := newSchemaTypes([]Table{c.t})
			assert.Contains(t, st, c.t.Name)

			val, err := st["tortoise"].New(c.d, c.id, c.fk)
			assert.NoError(t, err)

			assert.EqualValues(t, c.out, val)
		})
	}
}
