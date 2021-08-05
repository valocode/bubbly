// +build ignore

package main

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"

	"github.com/valocode/bubbly/ent/extensions"
)

func main() {
	var templates []*gen.Template
	templates = append(
		templates,
		gen.MustParse(gen.NewTemplate("static").
			Funcs(template.FuncMap{
				"filterNodeIndexFields": filterNodeIndexFields,
				"filterNodeIndexEdges":  filterNodeIndexEdges,
				"entToCtyFunc":          entToCtyFunc,
				"stringsJoin":           strings.Join,
			}).ParseDir("templates")),
	)

	ex, err := entgql.NewExtension(
		entgql.WithSchemaPath("../gql/ent.graphql"),
		entgql.WithConfigPath("../gql/gqlgen.yaml"),
		entgql.WithCustomRelaySpec(true, func(name string) string {
			return name + "_connection"
		}),
		entgql.WithNaming("snake"),
		entgql.WithWhereFilters(true),
		entgql.WithOrderBy(true),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}

	tsex, err := extensions.NewExtension()
	if err != nil {
		log.Fatalf("creating tsmodel extension: %v", err)
	}

	if err := entc.Generate("./schema",
		&gen.Config{
			Templates: templates,
		},
		entc.Extensions(ex, tsex),
	); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
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
