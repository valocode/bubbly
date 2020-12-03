package core

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

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
func (r *ResourceBlock) Kind() ResourceKind {
	return ResourceKind(r.ResourceKind)
}

// Name returns the name of the resource
func (r *ResourceBlock) Name() string {
	return r.ResourceName
}

// Namespace returns the namespace of the resource
func (r *ResourceBlock) Namespace() string {
	if md := r.Metadata; md != nil {
		if md.Namespace != "" {
			return md.Namespace
		}
	}
	return "default"
}

// Labels returns the labels of the resource
func (r *ResourceBlock) Labels() map[string]string {
	if md := r.Metadata; md != nil {
		return md.Labels
	}
	return nil
}

// APIVersion returns the APIVersion of the resource
func (r *ResourceBlock) APIVersion() APIVersion {
	return r.ResourceAPIVersion
}

// String returns a human-friendly string ID for the resource
func (r *ResourceBlock) String() string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		r.APIVersion(), r.Namespace(), r.ResourceKind, r.ResourceName,
	)
}

// MarshalJSON makes sure no-body marshals this type the default way
func (r *ResourceBlock) MarshalJSON() ([]byte, error) {
	panic("MarshalJSON not supported. Use the JSON() method.")
}

// UnmarshalJSON makes sure no-body unmarshals this type the default way
func (r *ResourceBlock) UnmarshalJSON([]byte) error {
	panic("UnmarshalJSON not supported. Use the JSON() method.")
}

// JSON returns a JSON representation of this resource block using the given
// ResourceContext.
func (r *ResourceBlock) JSON(ctx *ResourceContext) ([]byte, error) {
	// get the resource spec{} block as JSON
	sBody := r.SpecHCL.Body.(*hclsyntax.Body)
	bodyJSON, err := ctx.BodyToJSON(sBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal resource %s to json: %s", r.String(), err.Error())
	}

	// create the resource{} block as JSON
	resObj := gabs.New()
	resObj.Set(r.APIVersion(), "api_version")
	resObj.Set(bodyJSON, "spec")

	// create the top level JSON object that contains the resource
	jsonObj := gabs.New()
	jsonObj.Set(resObj.Data(), "resource", string(r.Kind()), r.Name())
	return jsonObj.Bytes(), nil
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
)

// ResourceKindPriority returns a list of the resource kinds by their priority
func ResourceKindPriority() []ResourceKind {
	return []ResourceKind{
		ExtractResourceKind, TransformResourceKind, LoadResourceKind,
		PipelineResourceKind, PipelineRunResourceKind, TaskRunResourceKind,
	}
}

// ResourceID is a type representation the ID of a resource
type ResourceID struct {
	Kind      ResourceKind
	Name      string
	Namespace string
}

// NewResourceIDFromString takes a string representation of an ID, e.g. from
// an HCL file, and returns a ResourceID to be used to identify the resource
// programmatically.
func NewResourceIDFromString(resStr string) *ResourceID {
	rID := &ResourceID{}
	subStr := strings.Split(resStr, "/")
	switch l := len(subStr); l {
	case 2:
		// We currently do not know the namespace of the resource if
		// not explicitely provided by the user. Fortunately, it is not needed
		// to successfully Get the resource during a task.Apply().
		rID.Kind = ResourceKind(subStr[0])
		rID.Name = subStr[1]
	case 3:
		rID.Namespace = subStr[0]
		rID.Kind = ResourceKind(subStr[1])
		rID.Name = subStr[2]
	default:
		panic(fmt.Sprintf("Unsupported resource string representation %s", resStr))
	}
	return rID
}

// String returns a string representation of a ResourceID
func (r *ResourceID) String() string {
	return strings.Join([]string{r.Namespace, string(r.Kind), r.Name}, "/")
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
