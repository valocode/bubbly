package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// Block is the interface for any type that will include HCL blocks
type Block interface {
	// Decode decodes the given HCL body into the given interface value.
	// It is implemented by types to enable them to decode any HCL bodies,
	// such as the spec{} of a resource, and any nested bodies within spec
	Decode(DecodeResourceFn) error
}

// Resource is the interface for any resources, such as Importer, Translator,
// etc.
type Resource interface {
	// All resources must implement the Block interface
	Block
	Name() string
	// Kind returns the ResourceKind
	Kind() ResourceKind
	APIVersion() APIVersion
	// Return a string representation of the resource, mainly for diagnostics
	String() string
	SpecHCLBody() hcl.Body
	SpecValue() ResourceSpec
	// ResourceBlock() *ResourceBlock
	// ResourceSpec() ResourceSpec
	// Every resource will provide an Output, specifying simple things like
	// if the Resource was applied successfully, and also any resulting value
	// for the Resource.
	Output() ResourceOutput
	// Spec() *ResourceBlockSpec
}

// ResourceSpec is the spec{} block inside a ResourceBlock
type ResourceSpec interface{}

// Importer interface is for any resources of type Importer
type Importer interface {
	Resource
}

// Translator interface is for any resources of type Translator
type Translator interface {
	Resource
}

// Publish interface is for any resources of type Publish
type Publish interface {
	Resource
}

type Pipeline interface {
	Resource
	Task(name string) Task
}

type PipelineRun interface {
	Resource
	// returns the ID of the pipeline to run
	Pipeline() string
	Inputs() cty.Value
}

type Task interface {
	ResourceKind() ResourceKind
	ResourceName() string
	TaskName() string
	Output() cty.Value
}
