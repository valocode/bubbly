package schema

import (
	"regexp"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/valocode/bubbly/ent/extensions/entmodel"
)

type Adapter struct {
	ent.Schema
}

func (Adapter) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "adapter"},
		entmodel.Annotation{},
	}
}

func (Adapter) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().
			Match(regexp.MustCompile("^[a-z0-9_]+$")).
			Annotations(
				entgql.OrderField("name"),
			),
		field.String("tag").NotEmpty().
			Annotations(
				entgql.OrderField("tag"),
			),
		field.String("module").NotEmpty(),
	}
}

func (Adapter) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owner", Organization.Type).Unique().Required(),
	}
}

func (Adapter) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "tag").
			Unique(),
	}
}
