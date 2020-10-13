package core

import (
	"github.com/hashicorp/hcl/v2"
)

type HCLMainType struct {
	ResourceBlocks ResourceBlocks `hcl:"resource,block"`
	ModuleBlocks   ModuleBlocks   `hcl:"module,block"`
	Locals         Locals         `hcl:"local,block"`
	Inputs         Inputs         `hcl:"input,block"`
	Outputs        Outputs        `hcl:"output,block"`
}

type ResourceBlocks []*ResourceBlock

type ResourceBlock struct {
	Kind       string            `hcl:",label"`
	Name       string            `hcl:",label"`
	APIVersion APIVersion        `hcl:"apiVersion,attr"`
	Metadata   *Metadata         `hcl:"metadata,block"`
	Spec       ResourceBlockSpec `hcl:"spec,block"`
	// Resource   Resource
}

// func (r *ResourceBlock) Resource() Resource {
// 	switch ResourceKind(r.Kind) {
// 	case importerResource:
// 		return v1.NewImporter(r.Spec.Body)
// 	}

// 	return nil
// }

type ResourceKind string

const (
	ImporterResourceKind ResourceKind = "importer"
)

type APIVersion string

type Metadata struct {
}

type ResourceBlockSpec struct {
	Body hcl.Body `hcl:",remain"`
}
