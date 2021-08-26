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

type GitCommit struct {
	ent.Schema
}

func (GitCommit) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "commit"},
		entmodel.Annotation{},
	}
}

func (GitCommit) Fields() []ent.Field {
	return []ent.Field{
		field.String("hash").Immutable().NotEmpty().
			Annotations(
				entgql.OrderField("hash"),
			),
		field.String("branch").Immutable().NotEmpty().
			Annotations(
				entgql.OrderField("branch"),
			),
		field.String("tag").Immutable().Optional().
			Annotations(
				entgql.OrderField("tag"),
			),
		field.Time("time").Immutable().
			Annotations(
				entgql.OrderField("time"),
			),
	}
}

func (GitCommit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("repo", Repo.Type).
			Unique().Required(),
		edge.To("release", Release.Type).
			Unique(),
	}
}

func (GitCommit) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("hash").
			Edges("repo").
			Unique(),
	}
}
