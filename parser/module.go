package parser

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/imdario/mergo"
	"github.com/rs/zerolog/log"
	"github.com/verifa/bubbly/api"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

// Modules is a wrapper for a slice of type Module
type Modules []*Module

// Module is the main type for decoding HCL.
// When invoking the parser, a "root" Module should be created which is
// responsible for fetching the HCL Body of code that the module references.
// E.g. if it is a "local" module type then it can point to a directory which
// contains some HCL code
//
// The structure is such that a Module is initialised from a ModuleBlock, and
// then contains that ModuleBlock.
type Module struct {
	//  ModuleBlock is the result of decoding a HCL module {} block, and is
	// the basis for creating a new Module
	ModuleBlock *core.ModuleBlock
	// ModuleBody stores the body of HCL that is referenced by the ModuleBlock
	ModuleBody hcl.Body
	// Name comes from the ModuleBlock and is the name of the module
	Name string
	// Scope is used to decode HCL and store the SymbolTable which keeps track
	// of the variables/traversals which are resolved. A Scope is not concerned
	// with Modules, however, it is only capable of resolving variables within
	// it's given Scope (HCL body)
	Scope *Scope
	// SubModules stores the slice of modules which are referenced within this
	// Module
	SubModules map[string]*Module
	// Value stores the resulting Value after parsing the body of HCL code that
	// this Module references, i.e. the Value is what we want in the end from
	// a Module
	Value *core.HCLMainType
}

// NewRootModule creates a new root module and is the entrypoint for the parser
func NewRootModule(baseDir string) *Module {
	mb := core.NewLocalModuleBlock(baseDir)
	return newModule(mb)
}

// newModule creates a new Module and initialises some values
func newModule(mb *core.ModuleBlock) *Module {
	return &Module{
		ModuleBlock: mb,
		Name:        mb.Name,
		Scope:       NewScope(),
		SubModules:  map[string]*Module{},
	}
}

func (m *Module) NewSubModuleFromBlock(block *hcl.Block) (*Module, hcl.Diagnostics) {
	mb := &core.ModuleBlock{}
	diags := m.Scope.Decode(block.Body, mb)
	mb.Name = block.Labels[0]
	subModule := newModule(mb)
	m.SubModules[subModule.Name] = subModule

	return subModule, diags
}

// Resolve is the main function for resolving/processing a module with the end
// goal of decoding the body of HCL that it refers to into a value.
//
// The logic can be summarised as follows:
// 1. The ModuleBlock references some HCL. Fetch that HCL into a hcl.Body
// 2. Resolve any variables/traversals needed to decode that body and then
//    decode that body
// 3. Create the list of SubModules from the module blocks defined within the
//    body that was decoded
// 4. For each sub-module, resolve it and then aggregate the return value
//    together
func (m *Module) Resolve() error {
	// if the Module is already resolved... then skip
	if m.IsResolved() {
		return nil
	}
	fmt.Printf("\n\nRESOLVING MODULE: %s\n\n", m.Name)
	// initialise the Value
	m.Value = &core.HCLMainType{}

	// get the HCL body that this module references
	body, err := m.ModuleBlock.Body()
	if err != nil {
		return fmt.Errorf(`Failed to resolve module "%s". Could not get hcl.Body: %s`, m.Name, err.Error())
	}
	m.ModuleBody = body

	body, err = m.expandBody(body, reflect.TypeOf(m.Value))
	if err != nil {
		return fmt.Errorf(`Failed to expand body for module "%s": %s`, m.Name, err.Error())
	}
	// assign the expanded body to the ModuleBody
	m.ModuleBody = body

	// resolve the expanded body variables
	traversals := walkVariables(body, reflect.TypeOf(m.Value))
	if err = m.resolveVariables(traversals); err != nil {
		return fmt.Errorf(`Failed to resolve variables for module "%s": %s`, m.Name, err.Error())
	}

	// after solving the variables, decode the body
	diags := m.Scope.Decode(body, m.Value)
	if diags.HasErrors() {
		return fmt.Errorf("Failed to decode module body: %s", diags.Error())
	}

	// create sub modules from the module {} blocks defined in this module
	for _, mb := range m.Value.ModuleBlocks {
		m.SubModules[mb.Name] = newModule(mb)
	}

	// Go through the SubModules and resolve them.
	// The slightly awkward logic is that the variables/traversals defined in a
	// ModuleBlock need to be resolved inside the parent module (which is where
	// they are defined). So before the SubModule can resolve itself, we need
	// to resolve the SubModule variables/traversals
	for _, subModule := range m.SubModules {
		if err := m.resolveSubModule(body, subModule); err != nil {
			// fmt.Printf("ASDASDASD: %s", err.Error())
			return fmt.Errorf("Failed to resolve sub module %s. %v", subModule.Name, err)
		}
		// merge values from submodule to parent module Value
		if err := mergo.Merge(m.Value, subModule.Value, mergo.WithAppendSlice); err != nil {
			return fmt.Errorf(
				`Failed to merge value from module submodule "%s" into parent module "%s": %s`,
				subModule.Name, m.Name, err.Error(),
			)
		}
	}

	// post process modules... this probably belongs somewhere else!
	for _, resBlock := range m.Value.ResourceBlocks {
		resource := api.NewResource(resBlock)
		if err := resource.Decode(m.decodeBody); err != nil {
			return fmt.Errorf(`Failed to decode resource "%s": %s`, resource.String(), err.Error())
		}
		switch resType := resource.(type) {
		case core.Importer:
			value, err := resType.Resolve()
			if err != nil {
				return fmt.Errorf(`Failed to resolve importer "%s": %s`, resType.String(), err.Error())
			}
			fmt.Printf("Importer %s Value: %s\n", resType.String(), value.GoString())
		case core.Translator:
			json, err := resType.JSON()
			if err != nil {
				return fmt.Errorf(`Failed to resolve translator "%s": %s`, resType.String(), err.Error())
			}
			fmt.Printf("Translator %s Value: %s\n", resType.String(), json)
		default:
			log.Warn().Msgf(`Resource Type "%s" not implemented yet...`, resType.String())
		}
	}
	return nil
}

// resolveSubModule is responsible for resolving a SubModule in the parent
// Module context (using the parent Module Scope).
// Traversals and variables that are referenced in the module {} block cannot be
// resolved from the sub-module itself, but need to be resolved in the parent
// Module.
func (m *Module) resolveSubModule(body hcl.Body, subModule *Module) error {

	// resolve the sub module parameters using the parent module eval context
	// to return the inputs for the sub module
	inputs := m.Inputs(subModule)
	subModule.Scope.SetInputs(inputs)

	err := subModule.ModuleBlock.Decode(m.decodeBody)
	if err != nil {
		return fmt.Errorf(`Failed to decode module block "%s": %s`, m.ModuleBlock.Name, err.Error())
	}

	if err := subModule.Resolve(); err != nil {
		return fmt.Errorf(`Failed to resolve sub module "%s": "%s"`, subModule.Name, err.Error())
	}

	outputs := subModule.Outputs()
	// set the outputs from the sub module in the parent module
	m.Scope.SetOutputs(subModule.Name, outputs)

	return nil
}

// decodeBody takes a body and an interface and decodes the body into the
// interface val, making sure to resolve all variables referenced in body,
// including references to other modules.
//
// This method is passed as a parameter to the different APIs which enable them
// to decode themselves without requiring a dependency on the parser
func (m *Module) decodeBody(body hcl.Body, val interface{}) error {
	ty := reflect.TypeOf(val)
	body, err := m.expandBody(body, ty)
	if err != nil {
		return fmt.Errorf(`Failed to expand body using type "%s": %s`, ty.String(), err.Error())
	}

	traversals := walkVariables(body, ty)
	if err := m.resolveVariables(traversals); err != nil {
		return fmt.Errorf(`Failed to resolve variables of body using type "%s": %s`, ty.String(), err.Error())
	}

	if diags := m.Scope.Decode(body, val); diags.HasErrors() {
		return fmt.Errorf(`Failed to decode body using type "%s": %s`, ty.String(), diags.Error())
	}
	return nil
}

func (m *Module) expandBody(body hcl.Body, ty reflect.Type) (hcl.Body, error) {
	traversals := walkExpandVariables(body, ty)
	if err := m.resolveVariables(traversals); err != nil {
		return nil, fmt.Errorf(`Failed to resolve variables of body using type "%s": %s`, ty.String(), err.Error())
	}

	return m.Scope.expandBody(body), nil
}

// resolveVariables is responsible for handling the logic to resolve variables
// (traversals). If the variable is an ordinary variable that does not require
// resolving a Module to obtain the value, then the work is delegated to the
// Scope to resolve the variable. If however, the variable references an output
// from a Module, then that Module is first resolved.
func (m *Module) resolveVariables(traversals []hcl.Traversal) error {

	for _, traversal := range traversals {
		// check first of all if the variable/traversal has already been
		// resolved
		_, exists := m.Scope.lookup(traversal)
		if exists {
			// the traversal has already been resolved, continue to next
			continue
		}
		// decide how to resolve based on the traversal root (e.g. "module")
		switch traversal.RootName() {
		// then we need to resolve a sub module
		case "module":
			// get the name of the module
			moduleName := traverserName(traversal[1])
			if _, exists := m.SubModules[moduleName]; exists {
				log.Warn().Msgf("Trying to resolve already resolved Sub Module %s", moduleName)
			}
			// traversal for retrieving the module block
			moduleTraversal := traversal[:2]
			ty := reflect.TypeOf(m.Value)
			// get the block for the module traversal
			block, _, diags := m.Scope.traverseBodyToBlock(m.ModuleBody, &ty, &moduleTraversal, nil)
			if diags.HasErrors() {
				return fmt.Errorf(`Failed to traverse body for module "%s": %s`, moduleName, diags.Error())
			}

			// create a new sub module from the retrieved block
			subModule, diags := m.NewSubModuleFromBlock(block)
			if diags.HasErrors() {
				return fmt.Errorf(`Failed to create new sub module module "%s": %s`, moduleName, diags.Error())
			}

			// resolve the sub module to get the necessary variable
			if err := m.resolveSubModule(m.ModuleBody, subModule); err != nil {
				return fmt.Errorf(`Failed to resolve module sub "%s": %s`, subModule.Name, err.Error())
			}
		default:
			// if default then delete the resolving to the scope
			diags := m.Scope.ResolveVariable(m.ModuleBody, m.Value, traversal)
			if diags.HasErrors() {
				return fmt.Errorf(`Failed to resolve variable "%s": %s`, traversalString(traversal), diags.Error())
			}
		}
	}
	return nil
}

// Inputs returns the input cty.Value for a module, based on the parameters
// specified in a module block
func (m *Module) Inputs(subModule *Module) cty.Value {
	inputValues := map[string]cty.Value{}
	for _, input := range subModule.ModuleBlock.Inputs {
		inputValues[input.Name] = input.Value
	}
	return cty.ObjectVal(inputValues)
}

// Outputs returns the output values for a module
func (m *Module) Outputs() cty.Value {
	outputValues := map[string]cty.Value{}

	for _, output := range m.Value.Outputs {
		outputValues[output.Name] = output.Value
	}

	return cty.ObjectVal(outputValues)
}

// IsResolved returns true if the module is already resolved, else false
func (m *Module) IsResolved() bool {
	return m.Value != nil
}
