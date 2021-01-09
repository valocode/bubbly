package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

func DecodeBody(bCtx *env.BubblyContext, body hcl.Body, val interface{}, inputs cty.Value) error {
	if diags := gohcl.DecodeBody(body, newEvalContext(inputs), val); diags.HasErrors() {
		return fmt.Errorf(`failed to decode body using type "%s": %s`, reflect.TypeOf(val).String(), diags.Error())
	}
	return nil
}

func DecodeExpandBody(bCtx *env.BubblyContext, body hcl.Body, val interface{}, inputs cty.Value) error {

	// expand the body so that dynamic blocks are processed
	node := dynblock.WalkVariables(body)

	// get the list of variables/traversals that exist
	traversals := walkVariables(bCtx, node, reflect.TypeOf(val))

	inputs, err := processVariables(bCtx, inputs, traversals)
	if err != nil {
		return fmt.Errorf("failed to process variables: %w", err)
	}

	eCtx := newEvalContext(inputs)
	expBody := dynblock.Expand(body, eCtx)
	if diags := gohcl.DecodeBody(expBody, eCtx, val); diags.HasErrors() {
		return fmt.Errorf(`failed to decode body using type "%s": %s`, reflect.TypeOf(val).String(), diags.Error())
	}

	return nil
}

func processVariables(bCtx *env.BubblyContext, inputs cty.Value, traversals []hcl.Traversal) (cty.Value, error) {
	dataRefs := cty.EmptyObjectVal
	for _, tr := range traversals {
		switch tr.RootName() {
		case "self":
			if len(tr) == 1 {
				// TODO: do we want to be strict here and just error if we have
				// unknown traversals?
				return cty.NilVal, fmt.Errorf("unknown variable %s", traversalString(bCtx, tr))
			}
			switch traverserName(bCtx, tr[1]) {
			case "data":
				dataRef, err := newDataRef(bCtx, tr)
				if err != nil {
					return cty.NilVal, fmt.Errorf("could not create DataRef from variable %s", traversalString(bCtx, tr))
				}
				fmt.Println(dataRef.GoString())
				dataRefs = appendDataRef(dataRefs, dataRef, traverserName(bCtx, tr[2]), traverserName(bCtx, tr[3]))
			}

		default:
			return cty.NilVal, fmt.Errorf(`unknown variable reference "%s". Only references to self are supported`, traversalString(bCtx, tr))
		}
	}

	var inputsMap objMapVal
	if inputs.IsNull() {
		inputsMap = make(objMapVal)
	} else {
		inputsMap = inputs.AsValueMap()
		if inputsMap == nil {
			inputsMap = make(objMapVal)
		}
	}

	inputsMap["data"] = dataRefs
	return cty.ObjectVal(inputsMap), nil
}

func newEvalContext(inputs cty.Value) *hcl.EvalContext {
	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"self": inputs,
		},
		Functions: stdfunctions(),
	}
}

func newDataRef(bCtx *env.BubblyContext, traversal hcl.Traversal) (cty.Value, error) {
	if len(traversal) != 4 {
		return cty.NilVal, fmt.Errorf("data reference must consist of four parts, e.g. self.data.table_name.field")
	}
	return cty.CapsuleVal(DataRefType, &DataRef{
		TableName: traverserName(bCtx, traversal[2]),
		Field:     traverserName(bCtx, traversal[3]),
	}), nil
}

func appendDataRef(dataRefs cty.Value, dataRef cty.Value, table string, field string) cty.Value {
	dataMap := dataRefs.AsValueMap()
	if dataMap == nil {
		dataMap = make(objMapVal)
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
		inputMap = make(objMapVal)
	}
	inputMap[attr] = val
	return cty.ObjectVal(inputMap)
}

type objMapVal map[string]cty.Value

var DataRefType = cty.CapsuleWithOps(
	"DataRef", reflect.TypeOf(DataRef{}),
	&cty.CapsuleOps{
		GoString: func(val interface{}) string { return fmt.Sprintf("%#v", val) },
	},
)

// DataRef is a data block that does not contain a
// static value but references a value from another
// data block.
type DataRef struct {
	TableName string `json:"table"`
	Field     string `json:"field"`
}

// traversalsString is a helper to return a string representation of traversals
func traversalsString(bCtx *env.BubblyContext, traversals []hcl.Traversal) string {
	strTraversals := []string{}
	for _, traversal := range traversals {
		strTraversals = append(strTraversals, traversalString(bCtx, traversal))
	}
	return strings.Join(strTraversals, ", ")
}

// traversalString is a helper to return a string representation of a traversal
func traversalString(bCtx *env.BubblyContext, traversal hcl.Traversal) string {
	if len(traversal) == 0 {
		return ""
	}
	retStr := traverserName(bCtx, traversal[0])
	for _, tr := range traversal[1:] {
		retStr = fmt.Sprintf("%s.%s", retStr, traverserName(bCtx, tr))
	}
	return retStr
}

// traverserName gets the Name or the given traverser
func traverserName(bCtx *env.BubblyContext, tr hcl.Traverser) string {
	switch tt := tr.(type) {
	case hcl.TraverseRoot:
		return tt.Name
	case hcl.TraverseAttr:
		return tt.Name
	default:
		bCtx.Logger.Panic().Msgf("Unknown type of traverser: %s", reflect.TypeOf(tt).String())
	}
	// redundant but stops go linter from complaining
	return ""
}
