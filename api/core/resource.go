package core

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
)

const ResourceTableName = "_resource"
const SchemaTableName = "_schema"
const EventTableName = "_event"

// ResourceBlocks is a wrapper for a slice of type ResourceBlock
type ResourceBlocks []*ResourceBlock

// ResourceBlockHCLWrapper is a simple wrapper for a ResourceBlock which is
// used when decoding a single resource block in HCL/JSON.
type ResourceBlockHCLWrapper struct {
	ResourceBlock ResourceBlock `hcl:"resource,block"`
}

// ResourceBlock represents the resource{} block in HCL.
type ResourceBlock struct {
	ResourceKind       string            `hcl:",label" json:"kind" mapstructure:"kind"`
	ResourceName       string            `hcl:",label" json:"name" mapstructure:"name"`
	ResourceAPIVersion APIVersion        `hcl:"api_version,optional" json:"api_version" mapstructure:"api_version"`
	Metadata           *Metadata         `hcl:"metadata,block" json:"metadata" mapstructure:"metadata"`
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

func (r ResourceBlock) ID() string {
	return fmt.Sprintf("%s/%s", r.Kind(), r.Name())
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
		"%s/%s",
		r.ResourceKind, r.ResourceName,
	)
}

// Data returns a Data block representation of the resource which can be
// sent to the bubbly store.
func (r ResourceBlock) Data() (Data, error) {
	resJSON, err := r.ResourceBlockJSON()
	if err != nil {
		return Data{}, fmt.Errorf("failed to create ResourceBlockJSON for resource: %s: %w", r.String(), err)
	}
	return resJSON.Data(), nil
}

// ResourceBlockJSON returns a ResourceBlockJSON representation of this
// ResourceBlock
func (r ResourceBlock) ResourceBlockJSON() (*ResourceBlockJSON, error) {
	specBytes, err := r.specBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get raw spec for resource: %s: %w", r.String(), err)
	}
	return NewResourceBlockJSON(r, string(specBytes)), nil
}

// MarshalJSON is customized to marshal a ResourceBlock, and thereby a resource
func (r ResourceBlock) MarshalJSON() ([]byte, error) {
	resJSON, err := r.ResourceBlockJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to create ResourceBlockJSON for resource: %s: %w", r.String(), err)
	}
	b, err := json.Marshal(resJSON)
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

func (r ResourceBlock) specBytes() ([]byte, error) {
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
	fileBytes, err := os.ReadFile(srcRange.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read resource file: %w", err)
	}
	if !srcRange.CanSliceBytes(fileBytes) {
		return nil, fmt.Errorf("cannot slice bytes for resource %s in filename %s", r.String(), srcRange.Filename)
	}
	specBytes := srcRange.SliceBytes(fileBytes)
	// specBytes contains the block paranthesis "{" and "}". Remove them
	return specBytes[1 : len(specBytes)-1], nil

}

func NewResourceBlockJSON(resBlock ResourceBlock, specRaw string) *ResourceBlockJSON {
	return &ResourceBlockJSON{
		ResourceBlockAlias: ResourceBlockAlias(resBlock),
		SpecRaw:            specRaw,
	}
}

type ResourceBlockAlias ResourceBlock
type ResourceBlockJSON struct {
	ResourceBlockAlias `mapstructure:",squash"`
	SpecRaw            string `json:"spec" mapstructure:"spec"`
}

func (r *ResourceBlockJSON) ResourceBlock() (ResourceBlock, error) {
	resBlock := ResourceBlock(r.ResourceBlockAlias)
	err := parser.ParseResource(env.NewBubblyContext(), []byte(r.SpecRaw), &resBlock.SpecHCL)
	return resBlock, err
}

// Data produces a core.Data type of this resource.
// The Data type is produced so that it can be sent to the store as any other
// piece of data, and therefore the store does not need to implement anything
// specific for a resource.
func (r *ResourceBlockJSON) Data() Data {

	var metaMap = make(map[string]cty.Value, 2)

	if r.Metadata != nil {
		var metaLabels = make(map[string]cty.Value, len(r.Metadata.Labels))
		for k, v := range r.Metadata.Labels {
			metaLabels[k] = cty.StringVal(v)
		}
		metaMap["labels"] = cty.ObjectVal(metaLabels)
	}
	d := Data{
		TableName: ResourceTableName,
		Fields: map[string]cty.Value{
			"id":          cty.StringVal(r.String()),
			"name":        cty.StringVal(r.ResourceName),
			"kind":        cty.StringVal(r.ResourceKind),
			"api_version": cty.StringVal(string(r.ResourceAPIVersion)),
			"metadata":    cty.ObjectVal(metaMap),
			"spec":        cty.StringVal(r.SpecRaw),
		},
	}
	return d
}

// String returns a human-friendly string ID for the resource
func (r *ResourceBlockJSON) String() string {
	return fmt.Sprintf(
		"%s/%s",
		r.ResourceKind, r.ResourceName,
	)
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
	// RunResourceKind represents the resource kind run
	RunResourceKind = "run"
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
		RunResourceKind,
	}
}

func ResourceRunKinds() []ResourceKind {
	return []ResourceKind{
		RunResourceKind,
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
	Labels map[string]string `json:"labels,omitempty" hcl:"labels,attr"`
}

// ResourceBlockSpec represents the spec{} block within a resource
type ResourceBlockSpec struct {
	Body hcl.Body `hcl:",remain"`
}
