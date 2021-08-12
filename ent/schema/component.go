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
	return []ent.Field{
		field.Text("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.Text("vendor").Default("").
			Annotations(
				entgql.OrderField("vendor"),
			),
		field.Text("version").NotEmpty().
			Annotations(
				entgql.OrderField("version"),
			),
		field.Text("description").Optional(),
		field.Text("url").Optional(),
	}
}

func (Component) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vulnerabilities", Vulnerability.Type),
		edge.To("licenses", License.Type),
		edge.From("uses", ReleaseComponent.Type).Ref("component"),
	}
}

func (Component) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "vendor", "version").
			Unique(),
	}
}
