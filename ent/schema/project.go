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

type Project struct {
	ent.Schema
}

func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "project"},
		entmodel.Annotation{},
	}
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().
			Unique().
			Annotations(
				entgql.OrderField("name"),
			),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repos", Repo.Type).Ref("projects"),
		edge.From("vulnerability_reviews", VulnerabilityReview.Type).Ref("projects"),
	}
}

func (Project) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}
