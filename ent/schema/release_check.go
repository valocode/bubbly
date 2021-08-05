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

type ReleaseCheck struct {
	ent.Schema
}

func (ReleaseCheck) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "release_check"},
	}
}

func (ReleaseCheck) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").
			Values("artifact", "unit_test", "security").
			Annotations(
				entgql.OrderField("type"),
			),
	}
}

func (ReleaseCheck) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("release", Release.Type).Unique().Required(),
	}
}

func (ReleaseCheck) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type").
			Edges("release").
			Unique(),
	}
}
