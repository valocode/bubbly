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

type License struct {
	ent.Schema
}

func (License) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "license"},
	}
}

func (License) Fields() []ent.Field {
	return []ent.Field{
		// License ID is follows the SPDX license IDs: https://spdx.dev/ids/
		field.Text("spdx_id").NotEmpty().
			Annotations(
				entgql.OrderField("spdx_id"),
			),
		field.Text("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		// reference points to a url where more information is available
		field.Text("reference").Optional(),
		field.Text("details_url").Optional(),
		field.Bool("is_osi_approved").Default(false),
	}
}

func (License) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("components", Component.Type).Ref("licenses"),
		edge.From("usages", LicenseUsage.Type).Ref("license"),
	}
}

func (License) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("spdx_id").Unique(),
	}
}
