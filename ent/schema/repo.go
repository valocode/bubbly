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

type Repo struct {
	ent.Schema
}

func (Repo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "repo"},
		entmodel.Annotation{},
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
		edge.To("projects", Project.Type),
		edge.From("commits", GitCommit.Type).Ref("repo"),
		edge.From("vulnerability_reviews", VulnerabilityReview.Type).Ref("repos"),
	}
}

func (Repo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}
