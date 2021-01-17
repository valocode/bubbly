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
		desc string
		in   core.DataBlocks
		refs dataRefs
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
			refs: dataRefs{},
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
			refs: dataRefs{
				"a": map[string]interface{}{
					"b": nil,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			var sc = newSaveContext()

			data := flatten(c.in, "")
			orderDataRefs(sc, data)

			assert.Equalf(t, sc.DataRefs, c.refs, "Refs not equal")
		})
	}
}
