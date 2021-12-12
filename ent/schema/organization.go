package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
)

type Organization struct {
	ent.Schema
}

func (Organization) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "organization"},
		// entgql.Skip(),
		entmodel.Annotation{},
	}
}

func (Organization) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().
			Unique(),
	}
}

func (Organization) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("projects", Project.Type).Ref("owner"),
		edge.From("repositories", Repository.Type).Ref("owner"),
	}
}
