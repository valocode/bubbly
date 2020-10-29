package core

// Resource is the interface for any resources, such as Importer, Translator,
// etc.
type Resource interface {
	Apply(*ResourceContext) ResourceOutput

	Name() string
	// Kind returns the ResourceKind
	Kind() ResourceKind
	APIVersion() APIVersion
	// Return a string representation of the resource, mainly for diagnostics
	String() string

	// TODO
	JSON(*ResourceContext) ([]byte, error)
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

// Pipeline interface is for any resources of type Pipeline
type Pipeline interface {
	Resource
}

// PipelineRun interface is for any resources of type PipelineRun
type PipelineRun interface {
	Resource
}

// Task interface represents a task inside a pipeline
type Task interface {
	Apply(*ResourceContext) error
}
