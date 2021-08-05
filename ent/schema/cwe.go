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

type CWE struct {
	ent.Schema
}

func (CWE) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "cwe"},
	}
}

func (CWE) Fields() []ent.Field {
	return []ent.Field{
		field.Text("cwe_id").NotEmpty().
			Annotations(
				entgql.OrderField("cwe_id"),
			),
		field.Text("description").Optional().
			Annotations(
				entgql.OrderField("description"),
			),
		field.Float("url").Optional(),
	}
}

func (CWE) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("issues", CodeIssue.Type).Ref("cwe"),
	}
}

func (CWE) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("cwe_id").Unique(),
	}
}
