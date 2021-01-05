package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

func TestDataRefs(t *testing.T) {
	cases := []struct {
		desc        string
		in          core.DataBlocks
		altData     core.DataBlocks
		refs        core.DataBlocks
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
			altData: core.DataBlocks{
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
			refs: core.DataBlocks{},
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
					Fields: core.DataFields{
						{
							Name: "d",
							Value: cty.CapsuleVal(
								parser.DataRefType,
								&parser.DataRef{
									TableName: "a",
									Field:     "b",
								},
							),
						},
					},
				},
			},
			altData: core.DataBlocks{
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
			refs: core.DataBlocks{
				{
					TableName: "c",
					Fields: core.DataFields{
						{
							Name: "d",
							Value: cty.CapsuleVal(
								parser.DataRefType,
								&parser.DataRef{
									TableName: "a",
									Field:     "b",
								},
							),
						},
					},
				},
			},
		},
		{
			// this test case checks if a data ref refers to a table that was
			// not provided as a DataBlock
			desc: "unspecified data block",
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
					Fields: core.DataFields{
						{
							Name: "d",
							Value: cty.CapsuleVal(
								parser.DataRefType,
								&parser.DataRef{
									TableName: "e",
									Field:     "f",
								},
							),
						},
					},
				},
			},
			errContains: "reference to unspecified data block",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {

			altData := make(core.DataBlocks, 0)
			refs := make(core.DataBlocks, 0)
			err := prepareDataRefs(c.in, &altData, &refs)

			if c.errContains != "" {
				assert.Contains(t, err.Error(), c.errContains)
				return
			}

			assert.NoError(t, err)
			assert.Equalf(t, altData, c.altData, "altData not equal")
			assert.Equalf(t, refs, c.refs, "refs not equal")
		})
	}
}
