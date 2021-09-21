package entmodel

import (
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	"entgo.io/contrib/entgql"
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
		SkipCreate bool
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
				"filterModelNodes":        filterModelNodes,
				"fieldsModelCreate":       fieldsModelCreate,
				"fieldsModelRead":         fieldsModelRead,
				"fieldsRequiredNoDefault": fieldsRequiredNoDefault,
				"fieldCreateTag":          fieldCreateTag,
				"fieldUpdateTag":          fieldUpdateTag,
				"fieldsAndID":             fieldsAndID,
				"filterNodes":             filterNodes,
			}).ParseDir(
			"extensions/entmodel/templates",
		)))

	return templates
}

func (Annotation) Name() string {
	return antName
}

func SkipCreate() Annotation {
	return Annotation{SkipCreate: true}
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
		}
	}
	return aptNodes, nil
}

func fieldsModelCreate(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range node.Fields {
		if val, ok := field.Annotations[antName]; ok {
			ant, err := DecodeAnnotation(val)
			if err != nil {
				log.Fatal(err)
			}
			if ant.SkipCreate {
				continue
			}
		}
		fields = append(fields, field)
	}
	return fields
}

func fieldsModelRead(node *gen.Type) []*gen.Field {
	return append(node.Fields, node.ID)
}

func fieldsRequiredNoDefault(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range node.Fields {
		if isFieldRequiredNoDefault(field) {
			fields = append(fields, field)
		}
	}
	return fields
}

func fieldCreateTag(field *gen.Field) string {
	var (
		jsonTag     = fmt.Sprintf(`json:"%s,omitempty"`, field.Name)
		validateTag string
		mapTag      = fmt.Sprintf(`mapstructure:"%s"`, field.Name)
	)
	if isFieldRequiredNoDefault(field) {
		validateTag = `validate:"required"`
	}
	return jsonTag + " " + validateTag + " " + mapTag
}

// fieldUpdateTag is like fieldCreateTag but made for update operations where no
// fields are required - only the fields that should be updated are provided
func fieldUpdateTag(field *gen.Field) string {
	var (
		jsonTag     = fmt.Sprintf(`json:"%s,omitempty"`, field.Name)
		validateTag string
		mapTag      = fmt.Sprintf(`mapstructure:"%s"`, field.Name)
	)
	return jsonTag + " " + validateTag + " " + mapTag
}

func isFieldRequiredNoDefault(field *gen.Field) bool {
	if !field.Optional {
		if !field.Default {
			return true
		}
	}
	return false
}

func fieldsAndID(node *gen.Type) []*gen.Field {
	return append(node.Fields, node.ID)
}

// filterNodes - copied from ent/contrib/entgql...
func filterNodes(nodes []*gen.Type) ([]*gen.Type, error) {
	var filteredNodes []*gen.Type
	for _, n := range nodes {
		ant := &entgql.Annotation{}
		if n.Annotations != nil && n.Annotations[ant.Name()] != nil {
			if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
				return nil, err
			}
			if ant.Skip {
				continue
			}
		}
		filteredNodes = append(filteredNodes, n)
	}
	return filteredNodes, nil
}
