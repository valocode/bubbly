package entapt

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

const (
	antName = "EntAdapter"
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

func NewExtension() (*Extensions, error) {
	ex := Extensions{}
	return &ex, nil
}

// Templates of the extension.
func (e *Extensions) Templates() []*gen.Template {
	var templates []*gen.Template
	templates = append(
		templates,
		gen.MustParse(gen.NewTemplate("adapter").
			Funcs(template.FuncMap{
				"filterAdapterNodes":     filterAdapterNodes,
				"filterAdapterRootNodes": filterAdapterRootNodes,
				"filterAdapterEdges":     filterAdapterEdges,
				"filterNodeIndexFields":  filterNodeIndexFields,
				"filterNodeIndexEdges":   filterNodeIndexEdges,
				"entToCtyFunc":           entToCtyFunc,
				"stringsJoin":            strings.Join,
			}).ParseDir(
			"extensions/entapt/templates",
			// "extensions/templates/graph.tmpl",
			// "extensions/templates/node.tmpl",
			// "extensions/templates/types.tmpl",
		)))

	return templates
}

// Hooks of the extension.
// func (e *Extensions) Hooks() []gen.Hook {
// 	return []gen.Hook{
// 		func(next gen.Generator) gen.Generator {
// 			return gen.GenerateFunc(func(g *gen.Graph) error {
// 				if err := next.Generate(g); err != nil {
// 					return err
// 				}

// 				return genTypescriptInterfaces(g)
// 			})
// 		},
// 	}
// }

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

func filterAdapterNodes(nodes []*gen.Type) ([]*gen.Type, error) {
	var aptNodes []*gen.Type
	for _, n := range nodes {
		if n.Annotations != nil && n.Annotations[antName] != nil {
			aptNodes = append(aptNodes, n)
		}
	}
	return aptNodes, nil
}

func filterAdapterRootNodes(nodes []*gen.Type) ([]*gen.Type, error) {
	var aptNodes []*gen.Type
	for _, n := range nodes {
		if n.Annotations != nil && n.Annotations[antName] != nil {
			ant, err := DecodeAnnotation(n.Annotations[antName])
			if err != nil {
				return nil, err
			}
			if ant.IsRoot {
				aptNodes = append(aptNodes, n)
			}
		}
	}
	return aptNodes, nil
}

func filterAdapterEdges(edges []*gen.Edge) ([]*gen.Edge, error) {
	var aptEdges []*gen.Edge
	for _, e := range edges {
		n := e.Type
		if n.Annotations != nil && n.Annotations[antName] != nil {
			ant, err := DecodeAnnotation(n.Annotations[antName])
			if err != nil {
				return nil, err
			}
			// Skip edges to root nodes
			if ant.IsRoot {
				continue
			}
			aptEdges = append(aptEdges, e)
		}
	}
	return aptEdges, nil
}

// filterNodeIndexFields takes a type and returns the fields which are part of the
// unique index
func filterNodeIndexFields(t *gen.Type) []*gen.Field {
	var (
		fields []*gen.Field
		index  *gen.Index
	)
	for _, idx := range t.Indexes {
		if idx.Unique {
			if index != nil {
				log.Fatalf("cannot have more than one unique index for type %s", t.Name)
			}
			index = idx
		}
	}
	if index == nil {
		return []*gen.Field{
			t.ID,
		}
	}
	for _, col := range index.Columns {
		for _, field := range t.Fields {
			if col == field.Name {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// filterNodeIndexEdges takes a type and returns the edges which are part of the
// unique index
func filterNodeIndexEdges(t *gen.Type) []*gen.Edge {
	var (
		edges []*gen.Edge
		index *gen.Index
	)
	for _, ind := range t.Indexes {
		if !ind.Unique {
			continue
		}
		if index != nil {
			log.Fatal("cannot have more than one unique index for type: ", t.Name)
		}
		index = ind
	}
	if index != nil {
		for _, col := range index.Columns {
			for _, edge := range t.Edges {
				if col == edge.Rel.Column() {
					edges = append(edges, edge)
				}
			}
		}
		// Finish here
		return edges
	}
	// Check unique fields because this way we can see if a value already
	// exists
	for _, edge := range t.Edges {
		if edge.Unique {
			edges = append(edges, edge)
		}
	}
	return edges
}

func entToCtyFunc(ty *field.Type, v string) string {
	switch *ty {
	case field.TypeInt8, field.TypeInt16, field.TypeInt32, field.TypeInt, field.TypeInt64:
		return fmt.Sprintf("cty.NumberIntVal(int64(%s))", v)
	case field.TypeUint8, field.TypeUint16, field.TypeUint32, field.TypeUint, field.TypeUint64:
		return fmt.Sprintf("cty.NumberUIntVal(uint64(%s))", v)
	case field.TypeFloat32, field.TypeFloat64:
		return fmt.Sprintf("cty.NumberFloatVal(%s)", v)
	case field.TypeBool:
		return fmt.Sprintf("cty.BoolVal(%s)", v)
	case field.TypeString:
		return fmt.Sprintf("cty.StringVal(%s)", v)
	case field.TypeEnum:
		return fmt.Sprintf("cty.StringVal(%s.String())", v)
	case field.TypeTime:
		return fmt.Sprintf("cty.StringVal(%s.Format(time.RFC3339))", v)
	case field.TypeJSON:
		return "UNSUPPORTED_TYPE_" + ty.String()
	case field.TypeUUID:
		return "UNSUPPORTED_TYPE_" + ty.String()
	case field.TypeBytes:
		return "UNSUPPORTED_TYPE_" + ty.String()
	case field.TypeOther:
		return "UNSUPPORTED_TYPE_" + ty.String()
	default:
		return "UNSUPPORTED_TYPE_" + ty.String()
	}
}
