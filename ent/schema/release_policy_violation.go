package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/valocode/bubbly/ent/extensions/entmodel"
)

type ReleasePolicyViolation struct {
	ent.Schema
}

func (ReleasePolicyViolation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "release_policy_violation"},
		entmodel.Annotation{},
	}
}

func (ReleasePolicyViolation) Fields() []ent.Field {
	return []ent.Field{
		field.String("message").NotEmpty(),
		field.Enum("severity").
			Values("suggestion", "warning", "error", "blocking"),
	}
}

func (ReleasePolicyViolation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("policy", ReleasePolicy.Type).Unique().Required(),
		edge.To("release", Release.Type).Unique().Required(),
	}
}
