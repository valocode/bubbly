package entmodel

import (
	"encoding/json"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
)

const (
	antName = "EntModel"
)

var (
	_ schema.Annotation = (*Annotation)(nil)
)

type (
	Extensions struct {
		entc.DefaultExtension
	}

	Annotation struct {
		IsRoot bool
	}
)

func NewExtension() *Extensions {
	return &Extensions{}
}

// Templates of the extension.
func (e *Extensions) Templates() []*gen.Template {
	var templates []*gen.Template
	templates = append(
		templates,
		gen.MustParse(gen.NewTemplate("model").
			Funcs(template.FuncMap{
				"filterModelNodes": filterModelNodes,
				"fieldsAndID":      fieldsAndID,
			}).ParseDir(
			"extensions/entmodel/templates",
		)))

	return templates
}

func (Annotation) Name() string {
	return antName
}

func Root() Annotation {
	return Annotation{IsRoot: true}
}

// Decode unmarshal annotation
func DecodeAnnotation(annotation interface{}) (Annotation, error) {
	var ant Annotation
	buf, err := json.Marshal(annotation)
	if err != nil {
		return Annotation{}, err
	}
	if err := json.Unmarshal(buf, &ant); err != nil {
		return Annotation{}, err
	}
	return ant, nil
}

func filterModelNodes(nodes []*gen.Type) ([]*gen.Type, error) {
	var aptNodes []*gen.Type
	for _, n := range nodes {
		if n.Annotations != nil && n.Annotations[antName] != nil {
			aptNodes = append(aptNodes, n)
			// for _, e := range filterModelNodeEdges(n.Edges) {
			// 	aptNodes = append(aptNodes, e.Type)

			// 	// aptNodes = append(aptNodes, []*gen)
			// }
		}
	}
	return aptNodes, nil
}

func fieldsAndID(node *gen.Type) []*gen.Field {
	return append(node.Fields, node.ID)
}
