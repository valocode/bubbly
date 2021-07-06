package core

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testSchema = `
table "test" {
	field "string" {
		type = string
	}
	field "map" {
		type = map(string)
	}
	field "object" {
		type = object({val: string})
	}
	field "time" {
		type = time
	}
	
}
`

type tableWrapper struct {
	Tables []TableHCL `hcl:"table,block"`
}

func TestTable(t *testing.T) {
	file, diags := hclparse.NewParser().ParseHCL([]byte(testSchema), "testSchema")
	require.Empty(t, diags.Errs())

	var val tableWrapper
	diags = gohcl.DecodeBody(file.Body, nil, &val)
	assert.Empty(t, diags)

	tables1, err := TablesFromHCL(val.Tables)
	require.NoError(t, err)

	t.Run("json", func(t *testing.T) {
		b, err := json.Marshal(tables1)
		require.NoError(t, err)

		var tables2 []Table = make([]Table, 0)
		err = json.Unmarshal(b, &tables2)
		require.NoError(t, err)
		require.Equal(t, len(val.Tables), len(tables2))
		for idx, t1 := range tables1 {
			t2 := tables2[idx]
			require.Equal(t, len(t1.Fields), len(t2.Fields))
			for j, f1 := range t1.Fields {
				f2 := t2.Fields[j]
				assert.Equal(t, f1.Name, f2.Name)
				assert.True(t, f1.Type.Equals(f2.Type))
			}
		}
	})
}
