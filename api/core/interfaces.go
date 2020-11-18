package core

import "github.com/verifa/bubbly/env"

// Resource is the interface for any resources, such as Extract, Transform,
// etc.
type Resource interface {
	Apply(*env.BubblyContext, *ResourceContext) ResourceOutput

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

// Extract interface is for any resources of type Extract
type Extract interface {
	Resource
}

// Transform interface is for any resources of type Transform
type Transform interface {
	Resource
}

// Load interface is for any resources of type Load
type Load interface {
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
	Apply(*env.BubblyContext, *ResourceContext) error
}
