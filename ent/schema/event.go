package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
)

type Event struct {
	ent.Schema
}

func (Event) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "event"},
		entmodel.Annotation{},
	}
}

func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.String("message").Default("").Immutable(),
		field.Enum("status").
			Values("ok", "error").
			Default("ok").
			Immutable(),
		field.Enum("type").
			Values("evaluate_release", "monitor").
			Immutable().
			Annotations(
				entgql.OrderField("type"),
			),
		field.Time("time").
			Immutable().
			Default(func() time.Time { return time.Now() }).
			Annotations(
				entgql.OrderField("time"),
			),
	}
}

func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		// Edges to a release, repo or project (i.e. the "scope" of the event) are optional
		edge.To("release", Release.Type).Unique(),
		edge.To("repo", Repo.Type).Unique(),
		edge.To("project", Project.Type).Unique(),
	}
}
