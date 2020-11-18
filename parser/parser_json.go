package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

// ResourceToJSON takes a resource and produces a JSON representation of it
// func (p *Parser) ResourceToJSON(resource core.Resource) ([]byte, error) {

// 	// get the resource spec{} block as JSON
// 	sBody := resource.SpecHCLBody().(*hclsyntax.Body)
// 	bodyJSON := p.bodyToJSON(sBody)

// 	// create the resource{} block as JSON
// 	resObj := gabs.New()
// 	resObj.Set(resource.APIVersion(), "api_version")
// 	resObj.Set(bodyJSON, "spec")

// 	// create the top level JSON object that contains the resource
// 	jsonObj := gabs.New()
// 	jsonObj.Set(resObj.Data(), "resource", string(resource.Kind()), resource.Name())
// 	return jsonObj.Bytes(), nil
// }

// JSONToResource takes a JSON representation of HCL as input and returns a
// core.Resource.
func (p *Parser) JSONToResource(json []byte) (*core.ResourceBlock, error) {
	// parse the json and pass in a unique filename for the parser, because
	// otherwise it returns the same json file again...
	file, diags := p.HCLParser.ParseJSON(json, fmt.Sprintf("json-%d", len(p.HCLParser.Files())))
	if diags.HasErrors() {
		return nil, fmt.Errorf("Failed to parse JSON: %s", diags.Error())
	}

	resWrap := &core.ResourceBlockHCLWrapper{}
	if diags := p.Scope.decodeBody(file.Body, resWrap); diags.HasErrors() {
		return nil, fmt.Errorf("Failed to decode JSON: %s", diags.Error())
	}

	// create a resource and store it in the parser
	p.Resources.NewResource(&resWrap.ResourceBlock)

	return &resWrap.ResourceBlock, nil
}

// BodyToJSON converts an hcl body written in hcl, to a JSON equivalent
func (p *Parser) BodyToJSON(body *hclsyntax.Body) (interface{}, error) {

	bodyJSON := gabs.New()

	for _, block := range body.Blocks {
		childJSON, err := p.BodyToJSON(block.Body)
		if err != nil {
			return nil, fmt.Errorf(`Failed to convert block %s to json: %s`, block.Type, err.Error())
		}

		path := append([]string{block.Type}, block.Labels...)
		bodyJSON.Set(childJSON, path...)
	}

	for _, attr := range body.Attributes {
		obj, err := p.exprToContainer(attr.Expr)
		if err != nil {
			return nil, fmt.Errorf(`Failed to convert expression for attribute %s to json: %s`, attr.Name, err.Error())
		}
		bodyJSON.Set(obj, attr.Name)
	}

	return bodyJSON.Data(), nil
}

func (p *Parser) exprToContainer(expr hcl.Expression) (interface{}, error) {
	// first try to evaluate the expression... if that works, great,
	// otherwise we store the raw value
	val, diags := expr.Value(p.Scope.EvalContext)
	if !diags.HasErrors() {
		return ctyValueToJSON(val)
	}
	switch ty := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		// LiteralValueExpr are just literal expressions, so there already is
		// a cty.Value
		return ctyValueToJSON(ty.Val)
	case *hclsyntax.FunctionCallExpr:
		// FunctionCallExpr is a function call
		ctyType, diags := typeexpr.TypeConstraint(expr)
		if !diags.HasErrors() {
			return typeexpr.TypeString(ctyType), nil
			// return gabs.Consume(typeexpr.TypeString(ctyType))
		}
		// TODO... what about ordinary function calls?
	case *hclsyntax.ScopeTraversalExpr:
		// ScopeTraversalExpr is a simple traversal, like my.variable.reference
		if len(ty.Traversal) == 1 {
			return traversalString(ty.Traversal), nil
		}
		return fmt.Sprintf("${%s}", traversalString(ty.Traversal)), nil
	case *hclsyntax.TemplateExpr:
		// TemplateExpr is an expression like: "raw_string_with_${a.var}"
		exprs := []string{}
		for _, exprEl := range ty.Parts {
			ret, err := p.exprToContainer(exprEl)
			if err != nil {
				return nil, err
			}
			// TODO: this might not work so well... let's see
			if retStr, ok := ret.(string); ok {
				exprs = append(exprs, retStr)
			} else {
				return nil, fmt.Errorf("Failed to convert return value to string: %v", ret)
			}
		}
		return strings.Join(exprs, ""), nil
	case *hclsyntax.TupleConsExpr:
		// TupleConsExpr is a tuple of expressions... so store those in an
		// interface{} list for gabs to deal with
		var list []interface{}
		for _, exprEl := range ty.ExprList() {
			ret, err := p.exprToContainer(exprEl)
			if err != nil {
				return nil, err
			}
			if retStr, ok := ret.(string); ok {
				list = append(list, strings.Trim(retStr, "\""))
			} else {
				return nil, fmt.Errorf("Failed to convert return value to string: %v", ret)
			}
		}
		return list, nil
	default:
	}
	panic(fmt.Sprintf(`exprToContainer: unsupported expression type "%s"`, reflect.TypeOf(expr).String()))
}

// ctyValueToJSON takes a cty value and returns a go interface that will be
// consumed by gabs to create a JSON representation
func ctyValueToJSON(val cty.Value) (interface{}, error) {
	if !val.IsKnown() {
		return "", fmt.Errorf("trying to convert not known cty.Value to a go value")
	}
	if val.IsNull() {
		return "nil", nil
	}

	t := val.Type()

	switch {
	case t.IsPrimitiveType():
		switch {
		case t == cty.Bool:
			if val.True() {
				return true, nil
			}
			return false, nil
		case t == cty.String:
			return val.AsString(), nil
		case t == cty.Number:
			return val.AsBigFloat(), nil
		default:
			panic(fmt.Sprintf("Unknown primitive cty type %s", t.GoString()))
		}
	case t.IsListType(), t.IsSetType(), t.IsTupleType():
		var list []interface{}
		it := val.ElementIterator()
		for it.Next() {
			_, eVal := it.Element()
			e, err := ctyValueToJSON(eVal)
			if err != nil {
				return "", fmt.Errorf(`Failed to get cty value of element %s: %s`, eVal.GoString(), err.Error())
			}
			list = append(list, e)
		}
		return list, nil
	case t.IsMapType(), t.IsObjectType():
		// TODO
		// var obj map[string]interface{}
	}
	panic(fmt.Sprintf(`Unsupported cty.Value type %s`, t.GoString()))
}
