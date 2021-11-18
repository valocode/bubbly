package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
	types "github.com/valocode/bubbly/ent/schema/types"
)

type Component struct {
	ent.Schema
}

func (Component) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "component"},
		entmodel.Annotation{},
	}
}

func (Component) Fields() []ent.Field {
	// Follow the Package URL (purl) schema for components:
	// https://github.com/package-url/purl-spec
	return []ent.Field{
		field.String("scheme").NotEmpty().
			Annotations(
				entgql.OrderField("scheme"),
			),
		field.String("namespace").Default("").
			Annotations(
				entgql.OrderField("namespace"),
			),
		field.String("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.String("version").NotEmpty().
			Annotations(
				entgql.OrderField("version"),
			),
		field.String("description").Optional(),
		field.String("url").Optional(),
		field.JSON("metadata", types.Metadata{}).Optional(),
		field.JSON("labels", types.Labels{}).Optional().
			Annotations(
				entgql.Skip(),
			),
	}
}

func (Component) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owner", Organization.Type).Unique().Required(),
		edge.To("vulnerabilities", Vulnerability.Type),
		edge.To("licenses", License.Type),
		edge.From("uses", ReleaseComponent.Type).Ref("component"),
	}
}

func (Component) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("scheme", "namespace", "name", "version").
			Unique(),
	}
}
