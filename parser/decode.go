package parser

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

func DecodeBody(body hcl.Body, val interface{}, inputs cty.Value) error {
	if diags := gohcl.DecodeBody(body, newEvalContext(inputs), val); diags.HasErrors() {
		return NewParserError(val, diags)
	}
	return nil
}

func DecodeExpandBody(body hcl.Body, val interface{}, inputs cty.Value) error {

	// expand the body so that dynamic blocks are processed
	node := dynblock.WalkVariables(body)

	// get the list of variables/traversals that exist
	traversals := walkVariables(node, reflect.TypeOf(val))

	inputs, diags := processVariables(inputs, traversals)
	if diags.HasErrors() {
		return NewParserError(val, diags)
	}

	eCtx := newEvalContext(inputs)
	expBody := dynblock.Expand(body, eCtx)
	if diags := gohcl.DecodeBody(expBody, eCtx, val); diags.HasErrors() {
		return NewParserError(val, diags)
	}

	return nil
}

func ExpressionValue(expr hcl.Expression, inputs cty.Value) (cty.Value, error) {
	eCtx := newEvalContext(inputs)
	value, diags := expr.Value(eCtx)
	if diags.HasErrors() {
		return cty.NilVal, NewParserError(nil, diags)
	}
	return value, nil
}

func processVariables(inputs cty.Value, traversals []hcl.Traversal) (cty.Value, hcl.Diagnostics) {
	var (
		diags    hcl.Diagnostics
		dataRefs = cty.EmptyObjectVal
	)
	for _, tr := range traversals {
		switch tr.RootName() {
		case "self":
			if len(tr) == 1 {
				// TODO: do we want to be strict here and just error if we have
				// unknown traversals?
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "unknown variable " + traversalString(tr),
					Detail:   "unknown variable found while traversing HCL: " + traversalString(tr),
					Subject:  tr.SourceRange().Ptr(),
				})
			}
			switch traverserName(tr[1]) {
			case "data":
				dataRef, err := newDataRef(tr)
				if err != nil {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "reference to data cannot be created for variable " + traversalString(tr),
						Detail:   "reference to data cannot be created for variable " + traversalString(tr),
						Subject:  tr.SourceRange().Ptr(),
					})
					break
				}
				dataRefs = appendDataRef(dataRefs, dataRef, traverserName(tr[2]), traverserName(tr[3]))
			}

		default:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "unknown variable reference " + traversalString(tr),
				Detail:   "only references to self are supported, unknown variable reference " + traversalString(tr),
				Subject:  tr.SourceRange().Ptr(),
			})
		}
	}

	var inputsMap map[string]cty.Value
	if inputs.IsNull() {
		inputsMap = make(map[string]cty.Value)
	} else {
		inputsMap = inputs.AsValueMap()
		if inputsMap == nil {
			inputsMap = make(map[string]cty.Value)
		}
	}

	inputsMap["data"] = dataRefs
	return cty.ObjectVal(inputsMap), diags
}

func newEvalContext(inputs cty.Value) *hcl.EvalContext {
	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"self": inputs,
		},
		Functions: stdfunctions(),
	}
}

func newDataRef(traversal hcl.Traversal) (cty.Value, error) {
	if len(traversal) != 4 {
		return cty.NilVal, fmt.Errorf("data reference must consist of four parts, e.g. self.data.table_name.field")
	}
	return cty.CapsuleVal(DataRefType, &DataRef{
		TableName: traverserName(traversal[2]),
		Field:     traverserName(traversal[3]),
	}), nil
}

func appendDataRef(dataRefs cty.Value, dataRef cty.Value, table string, field string) cty.Value {
	dataMap := dataRefs.AsValueMap()
	if dataMap == nil {
		dataMap = make(map[string]cty.Value)
	}

	fields, ok := dataMap[table]
	// if the table exists already, then append to the fields
	if ok {
		fields = addObjectAttr(fields, field, dataRef)
	} else {
		fields = addObjectAttr(cty.EmptyObjectVal, field, dataRef)
	}
	dataMap[table] = fields
	return cty.ObjectVal(dataMap)
}

func addObjectAttr(input cty.Value, attr string, val cty.Value) cty.Value {
	inputMap := input.AsValueMap()
	if inputMap == nil {
		inputMap = make(map[string]cty.Value)
	}
	inputMap[attr] = val
	return cty.ObjectVal(inputMap)
}

// traversalString is a helper to return a string representation of a traversal
func traversalString(traversal hcl.Traversal) string {
	if len(traversal) == 0 {
		return ""
	}
	retStr := traverserName(traversal[0])
	for _, tr := range traversal[1:] {
		retStr = fmt.Sprintf("%s.%s", retStr, traverserName(tr))
	}
	return retStr
}

// traverserName gets the Name or the given traverser
func traverserName(tr hcl.Traverser) string {
	switch tt := tr.(type) {
	case hcl.TraverseRoot:
		return tt.Name
	case hcl.TraverseAttr:
		return tt.Name
	case hcl.TraverseIndex:
		return tt.Key.AsString()
	default:
		panic("unknown type of traverser: " + reflect.TypeOf(tt).String())
	}
}
