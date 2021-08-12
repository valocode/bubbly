package schema

import (
	"errors"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
)

type TestCase struct {
	ent.Schema
}

func (TestCase) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "test_case"},
		entmodel.Annotation{},
	}
}

func (TestCase) Fields() []ent.Field {
	return []ent.Field{
		field.Text("name").Immutable().NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.Bool("result"),
		field.Text("message").NotEmpty(),
		field.Float("elapsed").Validate(func(f float64) error {
			if f < 0 {
				return errors.New("value cannot be negative")
			}
			return nil
		}).Default(0),
	}
}

func (TestCase) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("run", TestRun.Type).Unique().Required(),
	}
}
