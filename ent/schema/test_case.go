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
	types "github.com/valocode/bubbly/ent/schema/types"
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
		field.String("name").Immutable().NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.Bool("result"),
		field.String("message").NotEmpty(),
		field.Float("elapsed").Validate(func(f float64) error {
			if f < 0 {
				return errors.New("value cannot be negative")
			}
			return nil
		}).Default(0),
		field.JSON("metadata", types.Metadata{}).Optional(),
	}
}

func (TestCase) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("run", TestRun.Type).Unique().Required(),
	}
}
