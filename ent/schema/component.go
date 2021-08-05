package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Component struct {
	ent.Schema
}

func (Component) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "component"},
	}
}

func (Component) Fields() []ent.Field {
	return []ent.Field{
		field.Text("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.Text("vendor").NotEmpty().
			Annotations(
				entgql.OrderField("vendor"),
			),
		field.Text("version").NotEmpty().
			Annotations(
				entgql.OrderField("version"),
			),
		field.Text("description").NotEmpty(),
		field.Text("url").NotEmpty().
			Annotations(
				entgql.OrderField("url"),
			),
	}
}

func (Component) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vulnerabilities", Vulnerability.Type),
		edge.To("licenses", License.Type),
		edge.To("release", Release.Type),
	}
}

func (Component) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "vendor", "version").
			Unique(),
	}
}
