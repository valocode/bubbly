package core

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

// ResourceBlocks is a wrapper for a slice of type ResourceBlock
type ResourceBlocks []*ResourceBlock

type ResourceBlockHCLWrapper struct {
	ResourceBlock ResourceBlock `hcl:"resource,block"`
}

// ResourceBlock represents the resource{} block in HCL.
type ResourceBlock struct {
	ResourceKind       string            `hcl:",label" json:"kind"`
	ResourceName       string            `hcl:",label" json:"name"`
	ResourceAPIVersion APIVersion        `hcl:"api_version,attr" json:"api_version"`
	Metadata           *Metadata         `hcl:"metadata,block" json:"metadata"`
	SpecHCL            ResourceBlockSpec `hcl:"spec,block" json:"-"`
}

func (r *ResourceBlock) Kind() ResourceKind {
	return ResourceKind(r.ResourceKind)
}

func (r *ResourceBlock) Name() string {
	return r.ResourceName
}

func (r *ResourceBlock) APIVersion() APIVersion {
	return r.ResourceAPIVersion
}

func (r *ResourceBlock) SpecHCLBody() hcl.Body {
	return r.SpecHCL.Body
}

func (r *ResourceBlock) String() string {
	return fmt.Sprintf(
		"%s.%s.%s",
		r.APIVersion(), r.ResourceKind, r.ResourceName,
	)
}

// ResourceKind represents the different kinds of resources
type ResourceKind string

const (
	// ImporterResourceKind represents the kind importer
	ImporterResourceKind    ResourceKind = "importer"
	TranslatorResourceKind               = "translator"
	PublishResourceKind                  = "publish"
	PipelineResourceKind                 = "pipeline"
	PipelineRunResourceKind              = "pipeline_run"
)

func ResourceKindPriority() []ResourceKind {
	return []ResourceKind{
		ImporterResourceKind, TranslatorResourceKind, PublishResourceKind,
		PipelineResourceKind, PipelineRunResourceKind,
	}
}

// APIVersion represents the api_version of different resources
type APIVersion string

func (a *APIVersion) String() string {
	return string(*a)
}

// Metadata represents the metadata{} block in a resource... This could
// probably also be versioned?
type Metadata struct {
}

// ResourceBlockSpec represents the spec{} block within a resource
type ResourceBlockSpec struct {
	Body hcl.Body `hcl:",remain"`
}

// func (rb *ResourceBlockSpec) MarshalJSON() ([]byte, error) {
// 	return []byte(`{"value": "yoyo"}`), nil
// }
