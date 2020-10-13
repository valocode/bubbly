package core

import "github.com/zclconf/go-cty/cty"

// Block is the interface for any type that will include HCL blocks
type Block interface {
	// Decode decodes the given HCL body into the given interface value.
	// It is implemented by types to enable them to decode any HCL bodies,
	// such as the spec{} of a resource, and any nested bodies within spec
	Decode(DecodeBodyFn) error
}

// Resource is the interface for any resources, such as Importer, Translator,
// etc.
type Resource interface {
	// All resources must implement the Block interface
	Block
	// Return a string representation of the resource, mainly for diagnostics
	String() string
}

// ResourceSpec is the spec{} block inside a ResourceBlock
type ResourceSpec interface{}

// Importer interface is for any resources of type Importer
type Importer interface {
	Resource

	Resolve() (cty.Value, error)
}

// Translator interface is for any resources of type Translator
type Translator interface {
	Resource

	JSON() ([]byte, error)
}

// Upload interface is for any resources of type Upload
type Upload interface {
	Resource

	JSON() ([]byte, error)
}
