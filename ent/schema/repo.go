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

type Repo struct {
	ent.Schema
}

func (Repo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "repo"},
	}
}

func (Repo) Fields() []ent.Field {
	return []ent.Field{
		field.Text("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
	}
}

func (Repo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).Unique(),
		edge.From("commits", GitCommit.Type).Ref("repo"),
		edge.From("cve_rules", CVERule.Type).Ref("repo"),
	}
}

func (Repo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}
