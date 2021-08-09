package lang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type jsonOp struct {
	File          string         `hcl:"file,optional"`
	Preprocess    hcl.Expression `hcl:"preprocess,optional"`
	Structure     cty.Type
	StructureExpr hcl.Expression `hcl:"structure,attr"`
}

// Run parses the JSON file/contents and converts the value into a cty value
// based on the give format
func (o *jsonOp) Run(adapter *Adapter, opts AdapterOptions) (cty.Value, error) {
	var (
		file string
		r    io.Reader
		err  error
	)
	// Set the default file name
	file = adapter.Name + ".bubbly.json"
	// If the adapter sets a default filename, that will override it
	if o.File != "" {
		file = o.File
	}
	// If a filename was provided during the adapter run, that will override all
	if opts.Filename != "" {
		file = opts.Filename
	}

	r, err = os.Open(file)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error opening file %s: %w", o.File, err)
	}

	// fmt.Printf("Preproc: %#v\n", o.Preprocess)
	if !isNilExpression(o.Preprocess) {
		bytes, err := io.ReadAll(r)
		if err != nil {
			return cty.NilVal, fmt.Errorf("error reading json: %w", err)
		}
		variables := map[string]cty.Value{
			"raw_data": cty.StringVal(string(bytes)),
		}
		eCtx := NewEvalContext(variables)
		var jsonString string
		diags := gohcl.DecodeExpression(o.Preprocess, eCtx, &jsonString)
		if diags.HasErrors() {
			return cty.NilVal, fmt.Errorf("error during preprocessing: %w", diags)
		}
		r = strings.NewReader(jsonString)
	}

	return decodeJSON(r, o.Structure)
}

func (o *jsonOp) Decode(body hcl.Body) error {
	diags := gohcl.DecodeBody(body, nil, o)
	if diags.HasErrors() {
		return diags
	}
	if !isNilExpression(o.StructureExpr) {
		o.Structure, diags = typeexpr.TypeConstraint(o.StructureExpr)
		if diags.HasErrors() {
			return fmt.Errorf("invalid structure: %w", diags)
		}
	}

	return nil
}

func decodeJSON(r io.Reader, ty cty.Type) (cty.Value, error) {
	if ty == cty.NilType {
		b, err := io.ReadAll(r)
		if err != nil {
			return cty.NilVal, fmt.Errorf("error reading bytes whilst implying the JSON structure: %w", err)
		}
		// TODO: annoyingly the cty json package takes a []byte to convert into
		// an io.Reader... rather than just taking an io.Reader
		ty, err = ctyjson.ImpliedType(b)
		if err != nil {
			return cty.NilVal, fmt.Errorf("error getting implied json type: %w", err)
		}
		r = bytes.NewReader(b)
	}
	var data interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return cty.NilVal, fmt.Errorf("failed to decode JSON: %w", err)
	}
	val, err := gocty.ToCtyValue(data, ty)
	if err != nil {
		return cty.NilVal, err
	}

	return val, nil
}

// isNilExpression returns true if the given expression was not provided in the
// decoded hcl config. Unfortunately the gohcl package turns expressions into
// staticExpr so we cannot simply check if they are nil
func isNilExpression(expr hcl.Expression) bool {
	if expr == nil {
		return true
	}
	// If there are no variables, and the expr returns a null value, then
	// consider it to be a nil expression
	if expr.Variables() == nil {
		if val, diags := expr.Value(nil); !diags.HasErrors() && val.IsNull() {
			return true
		}
	}

	return false
}
