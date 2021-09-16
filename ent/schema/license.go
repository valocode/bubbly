package schema

import (
	"regexp"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
)

type License struct {
	ent.Schema
}

func (License) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "license"},
		entmodel.Annotation{},
	}
}

func (License) Fields() []ent.Field {
	return []ent.Field{
		// License ID is, for example, the SPDX ID: https://spdx.dev/ids/
		field.String("license_id").NotEmpty().Unique().
			Match(regexp.MustCompile(`^[A-Za-z0-9-_\.]+$`)).
			Annotations(
				entgql.OrderField("license_id"),
			),
		field.String("name").Optional().
			Annotations(
				entgql.OrderField("name"),
			),
		// reference points to a url where more information is available
		field.String("reference").Optional(),
		field.String("details_url").Optional(),
		field.Bool("is_osi_approved").Default(false),
	}
}

func (License) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owner", Organization.Type).Unique().Required(),
		edge.From("components", Component.Type).Ref("licenses"),
		edge.From("instances", ReleaseLicense.Type).Ref("license"),
	}
}
