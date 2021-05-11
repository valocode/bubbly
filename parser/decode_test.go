package parser

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

type testHCLValue struct {
	Value cty.Value `hcl:"value,attr"`
}

func TestDecode(t *testing.T) {
	// bCtx := env.NewBubblyContext()
	file, diags := hclparse.NewParser().ParseHCL([]byte("value = self.data.my_table.my_field"), "testing")
	assert.Equalf(t, diags.HasErrors(), false, diags.Error())
	var val testHCLValue
	err := DecodeExpandBody(file.Body, &val, cty.EmptyObjectVal)
	assert.NoErrorf(t, err, "failed to decode body")
}
