package store

import (
	"errors"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

type testProvider struct {
	data map[string]map[string]cty.Value
}

func (p *testProvider) Create(_ core.Tables) error {
	panic("not implemented")
}

func (p *testProvider) Save(data core.DataBlocks) (core.Tables, error) {
	panic("not implemented")
}

func (p *testProvider) ResolveScalar(graphql.ResolveParams) (interface{}, error) {
	panic("not implemented")
}

func (p *testProvider) ResolveList(graphql.ResolveParams) (interface{}, error) {
	panic("not implemented")
}

func (p *testProvider) LastValue(tableName, field string) (cty.Value, error) {
	t, ok := p.data[tableName]
	if !ok {
		return cty.NilVal, errors.New("no such table")
	}

	f, ok := t[field]
	if !ok {
		return cty.NilVal, errors.New("no such field")
	}

	return f, nil
}

func TestResolver(t *testing.T) {
	cases := []struct {
		desc        string
		p           provider
		in          core.DataBlocks
		out         core.DataBlocks
		errContains string
	}{
		{
			desc: "no data refs",
			in: core.DataBlocks{
				{
					TableName: "a",
					Fields: core.DataFields{
						{
							Name:  "b",
							Value: cty.NumberIntVal(111),
						},
					},
				},
			},
			out: core.DataBlocks{
				{
					TableName: "a",
					Fields: core.DataFields{
						{
							Name:  "b",
							Value: cty.NumberIntVal(111),
						},
					},
				},
			},
		},
		{
			desc: "in memory ref",
			in: core.DataBlocks{
				{
					TableName: "a",
					Fields: core.DataFields{
						{
							Name:  "b",
							Value: cty.NumberIntVal(111),
						},
					},
				},
				{
					TableName: "c",
					DataRefs: core.DataRefs{
						{
							TableName: "a",
							Field:     "b",
						},
					},
				},
			},
			out: core.DataBlocks{
				{
					TableName: "a",
					Fields: core.DataFields{
						{
							Name:  "b",
							Value: cty.NumberIntVal(111),
						},
					},
				},
				{
					TableName: "c",
					Fields: core.DataFields{
						{
							Name:  "a_b",
							Value: cty.NumberIntVal(111),
						},
					},
				},
			},
		},
		{
			desc: "provider ref",
			p: &testProvider{
				data: map[string]map[string]cty.Value{
					"a": map[string]cty.Value{
						"b": cty.NumberIntVal(111),
					},
				},
			},
			in: core.DataBlocks{
				{
					TableName: "c",
					DataRefs: core.DataRefs{
						{
							TableName: "a",
							Field:     "b",
						},
					},
				},
			},
			out: core.DataBlocks{
				{
					TableName: "c",
					Fields: core.DataFields{
						{
							Name:  "a_b",
							Value: cty.NumberIntVal(111),
						},
					},
				},
			},
		},
		{
			desc: "missing ref",
			p: &testProvider{
				data: map[string]map[string]cty.Value{
					"a": map[string]cty.Value{
						"___": cty.NumberIntVal(111),
					},
				},
			},
			in: core.DataBlocks{
				{
					TableName: "c",
					DataRefs: core.DataRefs{
						{
							TableName: "a",
							Field:     "b",
						},
					},
				},
			},
			errContains: "no such field",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			var (
				r        = &resolver{p: c.p}
				out, err = r.Resolve(c.in)
			)

			if c.errContains != "" {
				assert.Contains(t, err.Error(), c.errContains)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, out, c.out)
		})
	}
}
