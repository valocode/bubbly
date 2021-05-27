package core

import "github.com/valocode/bubbly/env"

// Resource is the interface for any resources, such as Extract, Transform,
// etc.
type Resource interface {
	SubResource

	Name() string
	// Kind returns the ResourceKind
	Kind() ResourceKind
	APIVersion() APIVersion
	ID() string
	// Return a string representation of the resource, mainly for diagnostics
	String() string
	// Data returns a Data block representation of the resource which can be
	// sent to the bubbly store
	Data() (Data, error)
}

type SubResource interface {
	// Run is the method called when a Resource (or SubResource) is run.
	// This can be considered the "main" function or entrypoint for a resource
	Run(*env.BubblyContext, *ResourceContext) ResourceOutput
}

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

// Run interface is for any resources of type Run
type Run interface {
	Resource
}

// Task interface represents a task inside a pipeline
type Task interface {
	SubResource
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
