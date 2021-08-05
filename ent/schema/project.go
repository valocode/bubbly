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

type Project struct {
	ent.Schema
}

func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "project"},
	}
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.Text("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repos", Repo.Type).Ref("project"),
		edge.From("releases", Release.Type).Ref("project"),
		edge.From("cve_rules", CVERule.Type).Ref("project"),
	}
}

func (Project) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}
