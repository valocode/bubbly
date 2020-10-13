package core

import (
	"github.com/hashicorp/hcl/v2"
)

// ResourceBlocks is a wrapper for a slice of type ResourceBlock
type ResourceBlocks []*ResourceBlock

// ResourceBlock represents the resource{} block in HCL.
type ResourceBlock struct {
	Kind       string            `hcl:",label"`
	Name       string            `hcl:",label"`
	APIVersion APIVersion        `hcl:"apiVersion,attr"`
	Metadata   *Metadata         `hcl:"metadata,block"`
	Spec       ResourceBlockSpec `hcl:"spec,block"`
}

// ResourceKind represents the different kinds of resources
type ResourceKind string

const (
	// ImporterResourceKind represents the kind importer
	ImporterResourceKind   ResourceKind = "importer"
	TranslatorResourceKind              = "translator"
)

// APIVersion represents the apiVersion of different resources
type APIVersion string

// Metadata represents the metadata{} block in a resource... This could
// probably also be versioned?
type Metadata struct {
}

// ResourceBlockSpec represents the spec{} block within a resource
type ResourceBlockSpec struct {
	Body hcl.Body `hcl:",remain"`
}
