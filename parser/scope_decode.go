package parser

import (
	"fmt"
	"reflect"

	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/gocty"
)

// DecodeExpandBody is a wrapper around decodeBody which also expands dynamic
// blocks
func (s *Scope) DecodeExpandBody(bCtx *env.BubblyContext, body hcl.Body, val interface{}) error {

	// create a decode context which will be passed around
	decodeCtx := newDecodeContext(body, val)

	// first resolve only references that are needed to expand any dynamic
	// blocks
	traversals := walkExpandVariables(decodeCtx.Body, decodeCtx.Type())
	if err := s.resolveVariables(bCtx, decodeCtx, traversals); err != nil {
		return fmt.Errorf(`failed to decode body using type "%s": %w`, decodeCtx.Type().String(), err)
	}

	// we have to expand before we resolve variables, otherwise the variables
	// will not exist
	body = dynblock.Expand(body, s.EvalContext)
	if diags := s.decodeBody(bCtx, body, val); diags.HasErrors() {
		return fmt.Errorf(`failed to decode body using type "%s": %s`, decodeCtx.Type().String(), diags.Error())
	}

	return nil
}

// decodeBody extracts the configuration within the given body into the given
// value. This value must be a non-nil pointer to either a struct or
// a map, where in the former case the configuration will be decoded using
// struct tags and in the latter case only attributes are allowed and their
// values are decoded into the map.
//
// The given EvalContext is used to resolve any variables or functions in
// expressions encountered while decoding. This may be nil to require only
// constant values, for simple applications that do not support variables or
// functions.
//
// The returned diagnostics should be inspected with its HasErrors method to
// determine if the populated value is valid and complete. If error diagnostics
// are returned then the given value may have been partially-populated but
// may still be accessed by a careful caller for static analysis and editor
// integration use-cases.
func (s *Scope) decodeBody(bCtx *env.BubblyContext, body hcl.Body, val interface{}) hcl.Diagnostics {
	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("target value must be a pointer, not %s", rv.Type().String()))
	}

	et := rv.Type().Elem()
	if et.Kind() != reflect.Struct {
		panic(fmt.Sprintf("target value must be pointer to struct, not %s", et.String()))
	}

	// create a decode context which will be passed around
	decodeCtx := newDecodeContext(body, val)

	return s.decodeBodyToStruct(bCtx, decodeCtx, body, rv.Elem())
}

// decodeBodyToStruct decodes the given hcl.Body into the given val, which must
// be a struct
func (s *Scope) decodeBodyToStruct(bCtx *env.BubblyContext, ctx *decodeContext, body hcl.Body, val reflect.Value) hcl.Diagnostics {
	schema, partial := gohcl.ImpliedBodySchema(val.Interface())

	var content *hcl.BodyContent
	var leftovers hcl.Body
	var diags hcl.Diagnostics
	if partial {
		content, leftovers, diags = body.PartialContent(schema)
	} else {
		content, diags = body.Content(schema)
	}
	if content == nil {
		return diags
	}

	tags := getFieldTags(val.Type())

	if tags.Remain != nil {
		fieldIdx := *tags.Remain
		field := val.Type().Field(fieldIdx)
		fieldV := val.Field(fieldIdx)
		switch {
		case bodyType.AssignableTo(field.Type):
			fieldV.Set(reflect.ValueOf(leftovers))
		case attrsType.AssignableTo(field.Type):
			attrs, attrsDiags := leftovers.JustAttributes()
			if len(attrsDiags) > 0 {
				diags = append(diags, attrsDiags...)
			}
			fieldV.Set(reflect.ValueOf(attrs))
		default:
			diags = append(diags, s.decodeBodyToStruct(bCtx, ctx, leftovers, fieldV)...)
		}
	}

	for name, fieldIdx := range tags.Attributes {
		attr := content.Attributes[name]
		field := val.Type().Field(fieldIdx)
		fieldV := val.Field(fieldIdx)

		if attr == nil {
			if !exprType.AssignableTo(field.Type) {
				continue
			}

			// As a special case, if the target is of type hcl.Expression then
			// we'll assign an actual expression that evalues to a cty null,
			// so the caller can deal with it within the cty realm rather
			// than within the Go realm.
			synthExpr := hcl.StaticExpr(cty.NullVal(cty.DynamicPseudoType), body.MissingItemRange())
			fieldV.Set(reflect.ValueOf(synthExpr))
			continue
		}

		switch {
		case attrType.AssignableTo(field.Type):
			fieldV.Set(reflect.ValueOf(attr))
		case exprType.AssignableTo(field.Type):
			fieldV.Set(reflect.ValueOf(attr.Expr))
		default:
			diags = append(diags, s.decodeExpression(
				bCtx, ctx, attr.Expr, fieldV.Addr().Interface(),
			)...)
		}
	}

	blocksByType := content.Blocks.ByType()

	for typeName, fieldIdx := range tags.Blocks {
		blocks := blocksByType[typeName]
		field := val.Type().Field(fieldIdx)

		ty := field.Type
		isSlice := false
		isPtr := false
		if ty.Kind() == reflect.Slice {
			isSlice = true
			ty = ty.Elem()
		}
		if ty.Kind() == reflect.Ptr {
			isPtr = true
			ty = ty.Elem()
		}

		if len(blocks) > 1 && !isSlice {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Duplicate %s block", typeName),
				Detail: fmt.Sprintf(
					"Only one %s block is allowed. Another was defined at %s.",
					typeName, blocks[0].DefRange.String(),
				),
				Subject: &blocks[1].DefRange,
			})
			continue
		}

		if len(blocks) == 0 {
			if isSlice || isPtr {
				if val.Field(fieldIdx).IsNil() {
					val.Field(fieldIdx).Set(reflect.Zero(field.Type))
				}
			} else {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("Missing %s block", typeName),
					Detail:   fmt.Sprintf("A %s block is required.", typeName),
					Subject:  body.MissingItemRange().Ptr(),
				})
			}
			continue
		}

		switch {

		case isSlice:
			elemType := ty
			if isPtr {
				elemType = reflect.PtrTo(ty)
			}
			sli := val.Field(fieldIdx)
			if sli.IsNil() {
				sli = reflect.MakeSlice(reflect.SliceOf(elemType), len(blocks), len(blocks))
			}

			for i, block := range blocks {
				if isPtr {
					if i >= sli.Len() {
						sli = reflect.Append(sli, reflect.New(ty))
					}
					v := sli.Index(i)
					if v.IsNil() {
						v = reflect.New(ty)
					}
					diags = append(diags, s.decodeBlockToValue(bCtx, ctx, block, v.Elem())...)
					sli.Index(i).Set(v)
				} else {
					diags = append(diags, s.decodeBlockToValue(bCtx, ctx, block, sli.Index(i))...)
				}
			}

			if sli.Len() > len(blocks) {
				sli.SetLen(len(blocks))
			}

			val.Field(fieldIdx).Set(sli)

		default:
			block := blocks[0]
			if isPtr {
				v := val.Field(fieldIdx)
				if v.IsNil() {
					v = reflect.New(ty)
				}
				diags = append(diags, s.decodeBlockToValue(bCtx, ctx, block, v.Elem())...)
				val.Field(fieldIdx).Set(v)
			} else {
				diags = append(diags, s.decodeBlockToValue(bCtx, ctx, block, val.Field(fieldIdx))...)
			}

		}

	}

	return diags
}

func (s *Scope) decodeBlockToValue(bCtx *env.BubblyContext, ctx *decodeContext, block *hcl.Block, v reflect.Value) hcl.Diagnostics {
	var diags hcl.Diagnostics

	ty := v.Type()

	switch {
	case blockType.AssignableTo(ty):
		v.Elem().Set(reflect.ValueOf(block))
	case bodyType.AssignableTo(ty):
		v.Elem().Set(reflect.ValueOf(block.Body))
	case attrsType.AssignableTo(ty):
		attrs, attrsDiags := block.Body.JustAttributes()
		if len(attrsDiags) > 0 {
			diags = append(diags, attrsDiags...)
		}
		v.Elem().Set(reflect.ValueOf(attrs))
	default:
		diags = append(diags, s.decodeBodyToStruct(bCtx, ctx, block.Body, v)...)

		if len(block.Labels) > 0 {
			blockTags := getFieldTags(ty)
			for li, lv := range block.Labels {
				lfieldIdx := blockTags.Labels[li].FieldIndex
				v.Field(lfieldIdx).Set(reflect.ValueOf(lv))
			}
		}

	}

	return diags
}

// decodeExpression extracts the value of the given expression into the given
// value. This method has been adapted from the gohcl package to also support
// typeexpr and cty.Type as a target value.
//
// Additionally, before trying to evaluate the expression it checks for any
// variable references (traversals) and attempts to resolve those.
func (s *Scope) decodeExpression(bCtx *env.BubblyContext, ctx *decodeContext, expr hcl.Expression, val interface{}) hcl.Diagnostics {
	if _, ok := val.(*cty.Type); ok {
		ctyType, diags := typeexpr.TypeConstraint(expr)
		if diags.HasErrors() {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsuitable type expr",
				Detail:   fmt.Sprintf("Unsuitable type expr: %s", diags.Error()),
				Subject:  expr.StartRange().Ptr(),
				Context:  expr.Range().Ptr(),
			})
			return diags
		}

		// assign the ctyType to the target field
		target := reflect.ValueOf(val).Elem()
		target.Set(reflect.ValueOf(ctyType))
		return diags
	}

	// switch ty := expr.(type) {
	// case *hclsyntax.ScopeTraversalExpr:
	// 	if ty.Traversal.RootName() == "self" {
	// 		// TODO this should not just set zero...
	// 		target := reflect.ValueOf(val).Elem()
	// 		target.Set(reflect.Zero(target.Type()))
	// 		return hcl.Diagnostics{}
	// 	}
	// }

	// TODO get traversals and resolve!
	traversals := expr.Variables()
	if err := s.resolveVariables(bCtx, ctx, traversals); err != nil {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unable to resolve traversals/variables",
			Detail:   fmt.Sprintf("Unable to resolve traversals/variables: %s", err.Error()),
			Subject:  expr.StartRange().Ptr(),
			Context:  expr.Range().Ptr(),
		}}
	}

	srcVal, diags := expr.Value(s.EvalContext)

	convTy, err := gocty.ImpliedType(val)
	if err != nil {
		panic(fmt.Sprintf("unsuitable DecodeExpression target: %s", err))
	}

	if srcVal.IsNull() {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Null pointer dereference",
			Detail:   fmt.Sprintf("Null pointer dereference: trying to assign from cty NilVal"),
			Subject:  expr.StartRange().Ptr(),
			Context:  expr.Range().Ptr(),
		})
		return diags
	}

	srcVal, err = convert.Convert(srcVal, convTy)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsuitable value type",
			Detail:   fmt.Sprintf("Unsuitable value: %s", err.Error()),
			Subject:  expr.StartRange().Ptr(),
			Context:  expr.Range().Ptr(),
		})
		return diags
	}

	err = gocty.FromCtyValue(srcVal, val)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsuitable value type",
			Detail:   fmt.Sprintf("Unsuitable value: %s", err.Error()),
			Subject:  expr.StartRange().Ptr(),
			Context:  expr.Range().Ptr(),
		})
	}

	return diags
}
