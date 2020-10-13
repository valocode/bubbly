package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

// decodeContext provides a simple struct type to encapsulate data needed to
// decode a given Body into a Value
type decodeContext struct {
	Body  hcl.Body
	Value interface{}
}

// newDecodeContext produces a new decodeContext
func newDecodeContext(body hcl.Body, val interface{}) *decodeContext {
	return &decodeContext{
		Body:  body,
		Value: val,
	}
}

// Type returns the reflect.Type of the Value
func (d *decodeContext) Type() reflect.Type {
	return reflect.TypeOf(d.Value)
}

// resolveContext is a wrapper around decodeContext and encapsulates the data
// needed to resolve variables in HCL
type resolveContext struct {
	*decodeContext

	OrigTraversal hcl.Traversal

	StepBody      hcl.Body
	StepTraversal hcl.Traversal
	StepType      reflect.Type
}

// newResolveContext produces a new resolveContext
func newResolveContext(decodeCtx *decodeContext, traversal hcl.Traversal) *resolveContext {
	return &resolveContext{
		decodeContext: decodeCtx,
		OrigTraversal: traversal,

		StepBody:      decodeCtx.Body,
		StepTraversal: traversal,
		StepType:      decodeCtx.Type(),
	}
}

// Scope is one of the main types in parsing HCL. Every Module will have a Scope
// and a Scope contains the SymbolTable. The Scope is the main type responsible
// for decoding an hcl.Body into the given value by constructing the
// SymbolTable and then using it during the decoding process
type Scope struct {
	*SymbolTable

	// TODO: if traversals have a cyclic dependency then this logic will recurse
	// forever... this needs to be solved by detecting cyclic dependencies.
	// A task for later...
	// TransitiveTraversal []hcl.Traversal
}

// NewScope creates a new Scope
func NewScope() *Scope {
	return NewScopeWithInput(
		map[string]cty.Value{},
	)
}

// NewScopeWithInput creates a new Scope with some input values; those needed to resolve
// inputs to a module
func NewScopeWithInput(inputs map[string]cty.Value) *Scope {
	return &Scope{
		SymbolTable: NewSymbolTableWithInputs(inputs),
	}
}

// Decode is the main entrypoint method to decode the scope body into the
// provided val. It also takes care of considering locals and inputs, and
// returns any outputs that have been defined
func (s *Scope) Decode(body hcl.Body, val interface{}) hcl.Diagnostics {
	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("Received non pointer type %s", rv.String()))
	}

	var diags hcl.Diagnostics
	// create a decode context which will be passed around
	decodeCtx := newDecodeContext(body, val)

	// first resolve only references that are needed to expand any dynamic
	// blocks
	traversals := walkExpandVariables(decodeCtx.Body, decodeCtx.Type())
	diags = append(diags, s.resolveVariables(decodeCtx, traversals)...)
	if diags.HasErrors() {
		return diags
	}

	// we have to expand before we resolve variables, otherwise the variables
	// will not exist
	decodeCtx.Body = s.expandBody(body)

	// second resolve all references that are needed to decode the entire body
	traversals = walkVariables(decodeCtx.Body, decodeCtx.Type())
	fmt.Sprintf("TRAVERSALS SCOPE: %s", traversalsString(traversals))
	// resolve the remaining variables on the expanded body
	diags = append(diags, s.resolveVariables(decodeCtx, traversals)...)

	if diags.HasErrors() {
		return diags
	}
	// thirdly decode the wanted value
	diags = append(diags, gohcl.DecodeBody(decodeCtx.Body, s.EvalContext(), val)...)

	return diags
}

// resolveVariables takes a list of traversals and rsolves them
func (s *Scope) resolveVariables(decodeCtx *decodeContext, traversals []hcl.Traversal) hcl.Diagnostics {

	var diags hcl.Diagnostics
	for _, traversal := range traversals {
		// create a new resolveContext for resolving this traversal
		resolveCtx := newResolveContext(decodeCtx, traversal)
		// let's see if the traversal already exists... in case of inputs they
		// can exist before we do any resolving
		_, exists := s.lookup(traversal)
		if exists {
			// continue to the next traversal
			continue
		}
		varDiags := s.resolveVariable(resolveCtx)
		diags = append(diags, varDiags...)
	}

	return diags
}

func (s *Scope) ResolveVariable(body hcl.Body, val interface{}, traversal hcl.Traversal) hcl.Diagnostics {
	dCtx := newDecodeContext(body, val)
	return s.resolveVariables(dCtx, []hcl.Traversal{traversal})
}

// resolveVariable takes a specific traversal, resolves it and inserts it into
// the SymbolTable.
// If the traversal has transitive traversals (traverersals that need to be
// resolved before it can be resolved) will be recursively resolved.
func (s *Scope) resolveVariable(ctx *resolveContext) hcl.Diagnostics {

	// recursively traverse as many blocks as possible
	block, content, diags := s.traverseBodyToBlock(ctx.Body, &ctx.StepType, &ctx.StepTraversal, nil)
	if diags.HasErrors() {
		return diags
	}

	if block == nil {
		panic(fmt.Sprintf(`Something went wrong... No diagnostics and an empty block whilst looking for "%s"`, traversalString(ctx.OrigTraversal)))
	}

	if len(ctx.StepTraversal) > 1 {
		return diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  `Unable to resolve variable`,
			Detail:   fmt.Sprintf(`Unable to resolve varible "%s" with leftover traversal: "%s"`, traversalString(ctx.OrigTraversal), traversalString(ctx.StepTraversal)),
			// TODO this Range is all wrong...
			// Subject: &content.MissingItemRange,
		})
	}

	// if there is 0 traversal steps left then we are looking for a
	// "special" attribute (e.g. local, input). If there is 1 traversal step
	// left then it *could* be a block, but a block without labels which is
	// uncommon
	attrName := s.nextAttr(ctx.StepTraversal, ctx.OrigTraversal)
	attr, exists := content.Attributes[attrName]
	if exists {
		diags = append(diags, s.resolveExpr(ctx, attr.Expr)...)
		return diags
	}
	// if the attribute was not found, and there were no more blocks to traverse
	// into, then we are not able to resolve the variable
	return diags.Append(&hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  `Could not find attribute`,
		Detail:   fmt.Sprintf(`Could not find attribute "%s" for traversal "%s"`, attrName, traversalString(ctx.OrigTraversal)),
		// TODO this Range is all wrong...
		// Subject:  &content.MissingItemRange,
	})
}

// traverseBodyToBlock traverses a body of HCL as far as possible, using the
// implied structure from a type (ty), for the given traversal, and returns the
// deepest block found and the BodyContent of that block.
//
// If the traversal does not match any blocks then the returned block and
// body content are nil
func (s *Scope) traverseBodyToBlock(body hcl.Body, ty *reflect.Type, traversal *hcl.Traversal, parent *hcl.Block) (*hcl.Block, *hcl.BodyContent, hcl.Diagnostics) {
	nestedType := nestedElem(*ty)
	zeroVal := reflect.Zero(nestedType)
	schema, _ := gohcl.ImpliedBodySchema(zeroVal.Interface())
	content, _, diags := body.PartialContent(schema)

	// not much point to continue if there are errors
	if diags.HasErrors() {
		return nil, nil, diags
	}

	if len(*traversal) > 0 {
		blockName := traverserName((*traversal)[0])
		tags := getFieldTags(nestedType)
		for _, block := range content.Blocks.OfType(blockName) {
			// fmt.Printf("Checking Block: %s\n", block.Type)
			// if a block then the variables should be resolved using the
			// block labels
			if remTraversal := traverseLabels((*traversal)[1:], block.Labels); remTraversal != nil {
				fieldIdx, exists := tags.Blocks[blockName]
				if !exists {
					panic(fmt.Sprintf(`Could not find HCL block type "%s" in type "%s"`, blockName, nestedType.String()))
				}
				*ty = nestedType.Field(fieldIdx).Type
				*traversal = remTraversal
				return s.traverseBodyToBlock(block.Body, ty, traversal, block)
			}
		}
	}
	// if there are no more matching blocks then there is an attribute that
	// needs to be resolved in traversal so return the matched block (if any)
	return parent, content, diags
}

func (s *Scope) expandBody(body hcl.Body) hcl.Body {
	return dynblock.Expand(body, s.EvalContext())
}

// resolveExpr resolves an expression and adds the value to the symbol table.
// Returns nil if it resolved correctly, otherwise returns the list of
// unresolved traversals
func (s *Scope) resolveExpr(ctx *resolveContext, expr hcl.Expression) hcl.Diagnostics {
	var diags hcl.Diagnostics
	varTraversals := expr.Variables()
	unresolvedTraversals := []hcl.Traversal{}
	for _, tr := range varTraversals {
		_, exists := s.lookup(tr)
		if !exists {
			unresolvedTraversals = append(unresolvedTraversals, tr)
		}
	}
	// Let's resolve any transitive unresolved traversals.
	// Doing this recursively means any recurise transitive dependencies
	// will also be solved
	diags = append(diags, s.resolveVariables(ctx.decodeContext, unresolvedTraversals)...)

	// if all transitive traversals were resolvable
	val, diags := expr.Value(s.EvalContext())
	if diags.HasErrors() {
		return diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  `Failed to get value of hcl expression`,
			Detail:   fmt.Sprintf(`Failed to get value of hcl expression at "%s" which requires variables "%s"`, expr.Range().String(), traversalsString(varTraversals)),
			// TODO this Range is all wrong...
			// Subject:  &content.MissingItemRange,
		})
	}

	s.insert(val, ctx.OrigTraversal)
	return diags
}

// nextAttr takes the current traversal step and returns the logical next
// attribute name that should be used to continue traversing.
// It should be noted that this attribute does not necessarily have to exist
// in the case that len(stepTrav) is greater than 0, in which case it *could*
// be a block type
func (s *Scope) nextAttr(stepTrav hcl.Traversal, origTrav hcl.Traversal) string {
	var attrName string

	switch len(stepTrav) {
	case 0:
		// if we have ran out of traversal steps, then it means we are
		// trying to resolve from a "special" value (e.g. local or input)
		// or something has just gone wrong!
		switch origTrav.RootName() {
		case "local":
			// we are trying to access a local value, so return the value
			// attribute
			attrName = "value"
		case "input":
			// we are trying to access an input value.
			// inputs are resolved on entry as they should be provided to
			// the scope at the beginning. Thus, if we are here, let's use
			// the default value
			attrName = "default"
		default:
			// what are the leftover cases? Looking for an id?
			// TODO diagnostic
			panic(fmt.Sprintf(`Ran out of traversals whilst resolving variable "%s"`, traversalString(origTrav)))
		}
	default:
		// if there are traversal steps left, use the next step as the attribute
		// name in this hcl body content
		attrName = traverserName(stepTrav[0])
	}
	return attrName
}

// traverseLabels returns the remaining hcl traversal if it is a match, nil
// otherwise
func traverseLabels(traversal hcl.Traversal, labels []string) hcl.Traversal {
	// if there are more labels than traversals then it cannot be a match
	if len(labels) > len(traversal) {
		return nil
	}
	for idx, label := range labels {
		name := traverserName(traversal[idx])
		if name != label {
			return nil
		}
	}
	// if we have iterated over the labels, and they all matched, then we
	// have  match
	return traversal[len(labels):]
}

// traversalsString is a helper to pring a string representation of traversals
func traversalsString(traversals []hcl.Traversal) string {
	strTraversals := []string{}
	for _, traversal := range traversals {
		strTraversals = append(strTraversals, traversalString(traversal))
	}
	return strings.Join(strTraversals, ", ")
}

// traversalString is a helper to pring a string representation of a traversal
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
	default:
		panic("Unknown type of traverser")
	}
}
