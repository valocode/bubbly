package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
)

type CodeIssue struct {
	ent.Schema
}

func (CodeIssue) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "code_issue"},
		entmodel.Annotation{},
	}
}

func (CodeIssue) Fields() []ent.Field {
	return []ent.Field{
		field.String("rule_id").Immutable().NotEmpty().
			Annotations(
				entgql.OrderField("rule_id"),
			),
		field.String("message").Immutable().NotEmpty(),
		field.Enum("severity").Immutable().
			Values("low", "medium", "high").
			Annotations(
				entgql.OrderField("severity"),
			),
		field.Enum("type").Immutable().
			Values("style", "security", "bug").
			Annotations(
				entgql.OrderField("type"),
			),
	}
}

func (CodeIssue) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("scan", CodeScan.Type).Unique().Required(),
	}
}
