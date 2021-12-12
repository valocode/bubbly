package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
	types "github.com/valocode/bubbly/ent/schema/types"
)

type Repository struct {
	ent.Schema
}

func (Repository) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entmodel.Annotation{},
	}
}

func (Repository) Fields() []ent.Field {
	return []ent.Field{
		// github.com/valocode/bubbly
		field.String("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.String("default_branch").
			NotEmpty().
			Default("main"),
		field.JSON("labels", types.Labels{}).Optional(),
	}
}

func (Repository) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owner", Organization.Type).Unique().Required(),
		edge.To("project", Project.Type).Unique().Required(),
		edge.To("head", Release.Type).Unique(),
		edge.From("commits", GitCommit.Type).Ref("repository"),
	}
}

func (Repository) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Edges("owner").
			Unique(),
	}
}
