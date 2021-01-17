package core

import "github.com/verifa/bubbly/env"

// Resource is the interface for any resources, such as Extract, Transform,
// etc.
type Resource interface {
	SubResource

	Name() string
	// Kind returns the ResourceKind
	Kind() ResourceKind
	APIVersion() APIVersion
	Namespace() string
	// Return a string representation of the resource, mainly for diagnostics
	String() string
	// Data returns a Data block representation of the resource which can be
	// sent to the bubbly store
	Data() (Data, error)
}

type SubResource interface {
	Apply(*env.BubblyContext, *ResourceContext) ResourceOutput
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
	SubResource
}

// TaskRun interface is for any resources of type TaskRun
type TaskRun interface {
	Resource
}

// Query interface is for any resources of type Query
type Query interface {
	Resource
}

// Criteria interface is for any resources of type Criteria
type Criteria interface {
	SubResource
}

// Condition interface represents a condition inside a Criteria
type Condition interface {
	SubResource
}

// Operation interface represents the operation inside a Criteria
type Operation interface {
	SubResource
}
