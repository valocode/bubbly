package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/parser"
)

const DefaultNamespace = "default"

// ResourceBlocks is a wrapper for a slice of type ResourceBlock
type ResourceBlocks []*ResourceBlock

// ResourceBlockHCLWrapper is a simple wrapper for a ResourceBlock which is
// used when decoding a single resource block in HCL/JSON.
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

// Kind returns the resource kind
func (r ResourceBlock) Kind() ResourceKind {
	return ResourceKind(r.ResourceKind)
}

// Name returns the name of the resource
func (r ResourceBlock) Name() string {
	return r.ResourceName
}

// Namespace returns the namespace of the resource
func (r ResourceBlock) Namespace() string {
	if md := r.Metadata; md != nil {
		if md.Namespace != "" {
			return md.Namespace
		}
	}
	return "default"
}

// Labels returns the labels of the resource
func (r ResourceBlock) Labels() map[string]string {
	if md := r.Metadata; md != nil {
		return md.Labels
	}
	return nil
}

// APIVersion returns the APIVersion of the resource
func (r ResourceBlock) APIVersion() APIVersion {
	return r.ResourceAPIVersion
}

// String returns a human-friendly string ID for the resource
func (r ResourceBlock) String() string {
	return fmt.Sprintf(
		"%s/%s/%s",
		r.Namespace(), r.ResourceKind, r.ResourceName,
	)
}

// MarshalJSON is customized to marshal a ResourceBlock, and thereby a resource
func (r ResourceBlock) MarshalJSON() ([]byte, error) {
	// get the source range of the hcl spec{} block, so that we can extract it
	// as raw text
	var srcRange hcl.Range
	switch body := r.SpecHCL.Body.(type) {
	case *hclsyntax.Body:
		srcRange = body.SrcRange
	default:
		return nil, fmt.Errorf("cannot get src range for unknown hcl.Body type %s", reflect.TypeOf(body).String())
	}
	// read the bubbly file containing the HCL
	fileBytes, err := ioutil.ReadFile(srcRange.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read resource file: %w", err)
	}
	if !srcRange.CanSliceBytes(fileBytes) {
		return nil, fmt.Errorf("cannot slice bytes for resource %s in filename %s", r.String(), srcRange.Filename)
	}
	specBytes := srcRange.SliceBytes(fileBytes)
	// specBytes contains the block paranthesis "{" and "}". Remove them
	specBytes = specBytes[1 : len(specBytes)-1]
	b, err := json.Marshal(NewResourceBlockJSON(r, string(specBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ResourceBLockJSON: %w", err)
	}
	return b, nil
}

// UnmarshalJSON is customized to unmarshal a resource
func (r *ResourceBlock) UnmarshalJSON(data []byte) error {
	var resJSON ResourceBlockJSON
	if err := json.Unmarshal(data, &resJSON); err != nil {
		return fmt.Errorf("failed to unmarshal ResourceBlockJSON: %w", err)
	}
	var err error
	*r, err = resJSON.ResourceBlock()
	return err
}

func NewResourceBlockJSON(resBlock ResourceBlock, specRaw string) *ResourceBlockJSON {
	return &ResourceBlockJSON{
		ResourceBlockAlias: ResourceBlockAlias(resBlock),
		SpecRaw:            specRaw,
	}
}

type ResourceBlockAlias ResourceBlock
type ResourceBlockJSON struct {
	ResourceBlockAlias
	SpecRaw string `json:"spec"`
}

func (r *ResourceBlockJSON) ResourceBlock() (ResourceBlock, error) {
	resBlock := ResourceBlock(r.ResourceBlockAlias)
	err := parser.ParseResource(env.NewBubblyContext(), []byte(r.SpecRaw), &resBlock.SpecHCL)
	return resBlock, err
}

func (r *ResourceBlockJSON) Validate() error {
	// TODO: check the fields which are needed
	return nil
}

// ResourceKind represents the different kinds of resources
type ResourceKind string

const (
	// ExtractResourceKind represents the resource kind extract
	ExtractResourceKind ResourceKind = "extract"
	// TransformResourceKind represents the resource kind transform
	TransformResourceKind = "transform"
	// LoadResourceKind represents the resource kind load
	LoadResourceKind = "load"
	// PipelineResourceKind represents the resource kind pipeline
	PipelineResourceKind = "pipeline"
	// PipelineRunResourceKind represents the resource kind pipeline_run
	PipelineRunResourceKind = "pipeline_run"
	// TaskRunResourceKind represents the resource kind task_run
	TaskRunResourceKind = "task_run"
	// QueryResourceKind represents the resource kind query
	QueryResourceKind = "query"
	// CriteriaResourceKind represents the resource kind criteria
	CriteriaResourceKind = "criteria"
)

// ResourceKindPriority returns a list of the resource kinds by their priority
func ResourceKindPriority() []ResourceKind {
	return []ResourceKind{
		ExtractResourceKind,
		TransformResourceKind,
		LoadResourceKind,
		QueryResourceKind,
		// Pipeline and Criteria both reference other resources
		PipelineResourceKind,
		CriteriaResourceKind,
		// last in the priority come the "Run" kinds
		PipelineRunResourceKind,
		TaskRunResourceKind,
	}
}

func ResourceRunKinds() []ResourceKind {
	return []ResourceKind{
		TaskRunResourceKind,
		PipelineRunResourceKind,
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
	Labels    map[string]string `json:"labels,omitempty" hcl:"labels,attr"`
	Namespace string            `json:"namespace,omitempty" hcl:"namespace,attr"`
}

// ResourceBlockSpec represents the spec{} block within a resource
type ResourceBlockSpec struct {
	Body hcl.Body `hcl:",remain"`
}
