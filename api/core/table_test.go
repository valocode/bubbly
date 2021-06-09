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
	Tables []Table `hcl:"table,block"`
}

func TestTable(t *testing.T) {
	file, diags := hclparse.NewParser().ParseHCL([]byte(testSchema), "testSchema")
	require.Empty(t, diags.Errs())

	var val tableWrapper
	diags = gohcl.DecodeBody(file.Body, nil, &val)
	assert.Empty(t, diags)

	for _, tbl := range val.Tables {
		err := tbl.Resolve()
		assert.NoError(t, err)
	}

	t.Run("json", func(t *testing.T) {
		b, err := json.Marshal(val)
		require.NoError(t, err)

		var val2 tableWrapper
		err = json.Unmarshal(b, &val2)
		require.NoError(t, err)
		require.Equal(t, len(val.Tables), len(val2.Tables))
		for idx, t1 := range val.Tables {
			t2 := val2.Tables[idx]
			require.Equal(t, len(t1.Fields), len(t2.Fields))
			for j, f1 := range t1.Fields {
				f2 := t2.Fields[j]
				assert.Equal(t, f1.Name, f2.Name)
				assert.True(t, f1.Type.Equals(f2.Type))
			}
		}
	})
}
