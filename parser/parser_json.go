package parser

import (
	"fmt"
	"log"
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
func (p *Parser) ResourceToJSON(resource core.Resource) ([]byte, error) {

	// get the resource spec{} block as JSON
	sBody := resource.SpecHCLBody().(*hclsyntax.Body)
	bodyJSON := p.bodyToJSON(sBody)

	// create the resource{} block as JSON
	resObj := gabs.New()
	resObj.Set(resource.APIVersion(), "api_version")
	resObj.Set(bodyJSON, "spec")

	// create the top level JSON object that contains the resource
	jsonObj := gabs.New()
	jsonObj.Set(resObj.Data(), "resource", string(resource.Kind()), resource.Name())
	return jsonObj.Bytes(), nil
}

// JSONToResource takes a JSON representation of HCL as input and returns a
// core.Resource, whilst also adding it to the Parser store of Resources.
func (p *Parser) JSONToResource(json []byte) (core.Resource, error) {
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

	resource := p.Resources.NewResource(&resWrap.ResourceBlock)

	return resource, nil
}

// bodyToJSON converts an hcl body written in hcl, to a JSON equivalent
func (p *Parser) bodyToJSON(body *hclsyntax.Body) interface{} {

	bodyJSON := gabs.New()

	for _, block := range body.Blocks {
		childJSON := p.bodyToJSON(block.Body)

		path := append([]string{block.Type}, block.Labels...)
		bodyJSON.Set(childJSON, path...)
	}

	for _, attr := range body.Attributes {
		obj, err := p.exprToContainer(attr.Expr)
		if err != nil {
			log.Fatalf(`Failed to convert expression to string: %s`, err.Error())
		}
		bodyJSON.Set(obj.Data(), attr.Name)
	}

	return bodyJSON.Data()
}

func (p *Parser) exprToContainer(expr hcl.Expression) (*gabs.Container, error) {
	// first try to evaluate the expression... if that works, great,
	// otherwise we store the raw value
	val, diags := expr.Value(p.Scope.EvalContext)
	if !diags.HasErrors() {
		return ctyValueToGabs(val)
	}
	switch ty := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		// LiteralValueExpr ar ejust literal expressions, so there already is
		// a cty.Value
		return ctyValueToGabs(ty.Val)
	case *hclsyntax.FunctionCallExpr:
		// FunctionCallExpr is a function call
		ctyType, diags := typeexpr.TypeConstraint(expr)
		if !diags.HasErrors() {
			return gabs.Consume(typeexpr.TypeString(ctyType))
		}
		// TODO... what about ordinary function calls?
	case *hclsyntax.ScopeTraversalExpr:
		// ScopeTraversalExpr is a simple traversal, like my.variable.reference
		if len(ty.Traversal) == 1 {
			return gabs.Consume(traversalString(ty.Traversal))
		}
		return gabs.Consume(fmt.Sprintf("${%s}", traversalString(ty.Traversal)))
	case *hclsyntax.TemplateExpr:
		// TemplateExpr is an expression like: "raw_string_with_${a.var}"
		exprs := []string{}
		for _, exprEl := range ty.Parts {
			retStr, err := p.exprToContainer(exprEl)
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, strings.Trim(retStr.String(), "\""))
		}
		obj, err := gabs.Consume(strings.Join(exprs, ""))
		return obj, err
	case *hclsyntax.TupleConsExpr:
		// TupleConsExpr is a tuple of expressions... so store those in an
		// interface{} list for gabs to deal with
		var list []interface{}
		for _, exprEl := range ty.ExprList() {
			retStr, err := p.exprToContainer(exprEl)
			if err != nil {
				return nil, err
			}
			list = append(list, retStr.Data())
		}
		return gabs.Consume(list)
	default:
	}
	panic(fmt.Sprintf(`exprToContainer: unsupported expression type "%s"`, reflect.TypeOf(expr).String()))
}

// TODO: this has not been tested at all
func (p *Parser) callExprToString(expr *hclsyntax.FunctionCallExpr) string {
	fmt.Printf("CALL EXPR!")
	static := expr.ExprCall()
	var retStr strings.Builder
	retStr.WriteString(static.Name)
	// for _, arg := range static.Arguments {
	// 	retStr.WriteString(" .. ")
	// 	// retStr.WriteString(p.exprToContainer(arg))
	// }

	return retStr.String()
}

func ctyValueToGabs(val cty.Value) (*gabs.Container, error) {
	if !val.IsKnown() {
		return nil, fmt.Errorf("trying to convern not known cty.Value to a go value")
	}
	if val.IsNull() {
		return gabs.New(), nil
	}

	t := val.Type()

	switch {
	case t == cty.Bool:
		return gabs.Consume(val.True())
	case t == cty.String:
		return gabs.Consume(val.AsString())
	case t == cty.Number:
		return gabs.Consume(val.AsBigFloat())
	case t.IsListType() || t.IsSetType() || t.IsTupleType():
		// TODO
	case t.IsMapType() || t.IsObjectType():
		// TODO
	}
	panic(fmt.Sprintf(`Unsupported cty.Value type %s`, t.GoString()))
}
