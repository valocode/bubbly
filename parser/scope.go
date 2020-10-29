package parser

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

// Scope is one of the main types in parsing HCL as it contains the SymbolTable.
// The Scope is the main type responsible for decoding an hcl.Body into the
// given value by constructing the SymbolTable and then using it during the
// decoding process
type Scope struct {
	*SymbolTable

	Parent *Scope
	// TODO: if traversals have a cyclic dependency then this logic will recurse
	// forever... this needs to be solved by detecting cyclic dependencies.
	// A task for later...
	// TransitiveTraversal []hcl.Traversal
}

// NewScope creates a new Scope
func NewScope() *Scope {
	return &Scope{
		SymbolTable: NewSymbolTable(),
		Parent:      nil,
	}
}

// NestedScope creates a new Scope which includes the parent Scope but can be
// used to decode a specific scope within the HCL (e.g. the spec{} within a
// resource which can reference it "self").
func (s *Scope) NestedScope(inputs cty.Value) *Scope {
	return &Scope{
		SymbolTable: s.NestedSymbolTable(inputs),
		Parent:      s,
	}
}

// InsertValue takes a cty.Value and a path, denoted as a slice of strings,
// and inserts the given value at the location given by the path (in the form
// of a traversal), e.g. cty.StringVal("hello"), ["self", "hello"]
// will make the StringVal "hello" available at the traversal "self.hello"
func (s *Scope) InsertValue(val cty.Value, path []string) {
	traversal := hcl.Traversal{
		hcl.TraverseRoot{Name: path[0]},
	}
	for _, step := range path[1:] {
		traversal = append(traversal, hcl.TraverseAttr{Name: step})
	}
	s.insert(val, traversal)
}

// resolveVariables takes a list of traversals and rsolves them
func (s *Scope) resolveVariables(decodeCtx *decodeContext, traversals []hcl.Traversal) error {

	for _, traversal := range traversals {
		// let's see if the traversal already exists... in case of inputs they
		// can exist before we do any resolving
		_, exists := s.lookup(traversal)
		if exists {
			// continue to the next traversal
			continue
		}
		// create a new resolveContext for resolving this traversal
		resolveCtx := newResolveContext(decodeCtx, traversal)
		if err := s.resolveVariable(resolveCtx); err != nil {
			return err
		}
	}
	return nil
}

// resolveVariable takes a specific traversal, resolves it and inserts it into
// the SymbolTable.
// If the traversal has transitive traversals (traverersals that need to be
// resolved before it can be resolved) will be recursively resolved.
func (s *Scope) resolveVariable(ctx *resolveContext) error {

	// recursively traverse as many blocks as possible
	attrs, err := s.traverseBodyToBlock(ctx, nil)
	if err != nil {
		return fmt.Errorf(`Failed to resolve variable "%s": %s`, traversalString(ctx.OrigTraversal), err.Error())
	}

	attr, exists := (*attrs)[ctx.Attribute]
	if exists {
		return s.resolveExpr(ctx, attr.Expr)
	}
	// if the attribute was not found, and there were no more blocks to traverse
	// into, then we are not able to resolve the variable
	return fmt.Errorf(`Could not find attribute "%s" for traversal "%s"`, ctx.Attribute, traversalString(ctx.OrigTraversal))
}

// traverseBodyToBlock traverses a body of HCL as far as possible, using the
// implied structure from a type (ty), for the given traversal, and returns the
// deepest block found and the BodyContent of that block.
//
// If the traversal does not match any blocks then the returned block and
// body content are nil
func (s *Scope) traverseBodyToBlock(ctx *resolveContext, parent *hcl.Block) (*hcl.Attributes, error) {
	nestedType := nestedElem(ctx.StepType)
	// fmt.Printf("traverseBodyToBlock: %s\t%s == %d\n", nestedType.String(), traversalString(ctx.StepTraversal), len(ctx.StepTraversal))
	zeroVal := reflect.Zero(nestedType)
	schema, _ := gohcl.ImpliedBodySchema(zeroVal.Interface())
	content, _, diags := ctx.StepBody.PartialContent(schema)

	// not much point to continue if there are errors
	if diags.HasErrors() {
		return nil, fmt.Errorf(`Could not traverse body using type "%s": %s`, nestedType.String(), diags.Error())
	}

	// if the traversal is length one then the reamining traversal is an
	// attribute, so finish finding the block
	if len(ctx.StepTraversal) == 0 {
		// then we have finished
		return &content.Attributes, nil
	}

	blockName := traverserName(ctx.StepTraversal[0])
	tags := getFieldTags(nestedType)
	for _, block := range content.Blocks.OfType(blockName) {
		// fmt.Printf("Checking Block: %s\n", block.Type)
		// if a block then the variables should be resolved using the
		// block labels
		// fmt.Printf("stepTraversal: %s == %d\n", traversalString(ctx.StepTraversal[1:]), len(ctx.StepTraversal[1:]))
		if remTraversal := traverseLabels(ctx.StepTraversal[1:], block.Labels); remTraversal != nil {
			// fmt.Printf("remTraversal: %s == %d\n", traversalString(remTraversal), len(remTraversal))
			fieldIdx, exists := tags.Blocks[blockName]
			if !exists {
				panic(fmt.Sprintf(`Could not find HCL block type "%s" in type "%s"`, blockName, nestedType.String()))
			}
			ctx.StepType = nestedType.Field(fieldIdx).Type
			ctx.StepTraversal = remTraversal
			ctx.StepBody = block.Body
			return s.traverseBodyToBlock(ctx, block)
		}
	}
	// if there are no more matching blocks and we did not exhaust all the step
	// traversals then there is an error
	return nil, fmt.Errorf(
		`Could not find a matching block for traversal "%s", with leftover traversal "%s"`,
		traversalString(ctx.OrigTraversal),
		traversalString(ctx.StepTraversal),
	)
}

// resolveExpr ...
func (s *Scope) resolveExpr(ctx *resolveContext, expr hcl.Expression) error {
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
	if err := s.resolveVariables(ctx.decodeContext, unresolvedTraversals); err != nil {
		return fmt.Errorf(
			`Failed to resolve transitive traversals "%s": %s`,
			traversalsString(unresolvedTraversals),
			err.Error(),
		)
	}

	// if all transitive traversals were resolvable
	val, diags := expr.Value(s.EvalContext)
	if diags.HasErrors() {
		return fmt.Errorf(`Failed to get value of hcl experssion: %s`, diags.Error())
	}

	s.insert(val, ctx.AliasTraversal)
	return nil
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

	OrigTraversal  hcl.Traversal
	AliasTraversal hcl.Traversal

	StepBody      hcl.Body
	StepTraversal hcl.Traversal
	StepType      reflect.Type

	Attribute string
}

// newResolveContext produces a new resolveContext
func newResolveContext(decodeCtx *decodeContext, traversal hcl.Traversal) *resolveContext {
	return newResolveContextWithAlias(decodeCtx, traversal, traversal)
}

// newResolveContextWithAlias produces a new resolveContext
func newResolveContextWithAlias(decodeCtx *decodeContext, traversal hcl.Traversal, alias hcl.Traversal) *resolveContext {
	var attrName string
	switch traversal.RootName() {
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
	case "self":
		attrName = traverserName(traversal[len(traversal)-1])
	default:
		log.Fatalf(`Unknown root traversal "%s" for reference "%s"`, traversal.RootName(), traversalString(traversal))
	}

	return &resolveContext{
		decodeContext:  decodeCtx,
		OrigTraversal:  traversal,
		AliasTraversal: alias,

		StepBody:      decodeCtx.Body,
		StepTraversal: traversal,
		StepType:      decodeCtx.Type(),

		Attribute: attrName,
	}
}
