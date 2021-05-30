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

type DataBlocks []Data

// Data will reference a Table name, and assign the Field values into the
// corresponding Field values in the Table
type Data struct {
	TableName     string          `hcl:",label" json:"data"`
	Fields        *DataFields     `hcl:"fields,block" json:"fields,omitempty"`
	Joins         []string        `hcl:"joins,optional" json:"joins,omitempty"`
	Policy        DataBlockPolicy `hcl:"policy,optional" json:"policy,omitempty"`
	IgnoreNesting bool            `hcl:"ignore_nesting,optional" json:"ignore_nesting,omitempty"`
	Data          DataBlocks      `hcl:"data,block" json:"nested_data,omitempty"`
}

// DataBlockPolicy defines the policy for how the data block shall be handled.
// When the bubbly store goes to save a data block, it should consider whether
// it should create and/or update the data block (default behaviour), only
// create the data block (fail on conflict), or only reference an existing data
// block with matching field values so that another data block can join to it
type DataBlockPolicy string

const (
	EmptyPolicy DataBlockPolicy = ""
	// DefaultPolicy is to create or update
	DefaultPolicy DataBlockPolicy = CreateUpdatePolicy
	// CreateUpdatePolicy means either create or update an existing data block
	// based on the unique constraints applied to the schema table that this data
	// block refers to
	CreateUpdatePolicy DataBlockPolicy = "create_update"
	// CreatePolicy means only create. If a conflict occurs on unique constraints
	// on the corresponding schema table, then error
	CreatePolicy DataBlockPolicy = "create"
	// ReferencePolicy means do not create or update, but only retrieve a reference
	// to an already saved data block, with the matching field values
	ReferencePolicy DataBlockPolicy = "reference"
	// ReferenceIfExistsPolicy is the same as ReferencePolicy but it does not
	// error in case a reference does not exist
	ReferenceIfExistsPolicy DataBlockPolicy = "reference_if_exists"
)

// DataFields is a slice of DataField
// type DataFields map[string]cty.Value
type DataFields struct {
	Values map[string]cty.Value `hcl:",remain"`
}

type dataBlockWrapper struct {
	Data []struct {
		Name   string `hcl:",label"`
		Type   string `hcl:"type,attr"`
		Fields struct {
			Values map[string]cty.Value `hcl:",remain"`
			// Values hcl.Attributes `hcl:",remain"`
		} `hcl:"fields,block"`
	} `hcl:"data,block"`
}

func TestDecode(t *testing.T) {
	// bCtx := env.NewBubblyContext()
	file, diags := hclparse.NewParser().ParseHCL([]byte("value = self.data.my_table.my_field"), "testing")
	assert.Equalf(t, diags.HasErrors(), false, diags.Error())
	var val testHCLValue
	err := DecodeExpandBody(file.Body, &val, cty.EmptyObjectVal)
	assert.NoErrorf(t, err, "failed to decode body")
}

func TestDecodeDynamicRemain(t *testing.T) {
	hcl := `
dynamic "data" {
	for_each = [1,2,3]
	iterator = it
	labels = ["whatever"]
	content {
		type = "${it.value}"
		fields {
			val = it.value
		}
	}
}
	`
	file, diags := hclparse.NewParser().ParseHCL([]byte(hcl), "testing")
	assert.Equalf(t, diags.HasErrors(), false, diags.Error())
	var val dataBlockWrapper
	err := DecodeExpandBody(file.Body, &val, cty.EmptyObjectVal)
	assert.NoErrorf(t, err, "failed to decode body")

	// for _, d := range val.Data {
	// 	for a, b := range d.Fields.Values {
	// 		t.Logf("%s : %#v", a, b.Expr.Value())
	// 	}
	// }
}
