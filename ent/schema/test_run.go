package schema

import (
	"context"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	gen "github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
	"github.com/valocode/bubbly/ent/hook"
)

type TestRun struct {
	ent.Schema
}

func (TestRun) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "test_run"},
		entmodel.Annotation{},
	}
}

func (TestRun) Fields() []ent.Field {
	return []ent.Field{
		field.String("tool").Immutable().NotEmpty().
			Annotations(
				entgql.OrderField("tool"),
			),
		field.Time("time").
			Immutable().
			Default(func() time.Time { return time.Now() }).
			Annotations(
				entgql.OrderField("time"),
			),
	}
}

func (TestRun) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("release", Release.Type).Unique().Required(),
		edge.From("entry", ReleaseEntry.Type).Ref("test_run").Unique(),
		edge.From("tests", TestCase.Type).Ref("run"),
	}
}

func (TestRun) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.TestRunFunc(func(ctx context.Context, m *gen.TestRunMutation) (ent.Value, error) {
					if err := createReleaseEntry(ctx, m); err != nil {
						return nil, err
					}
					return next.Mutate(ctx, m)
				})
			},
			// Limit the hook only for these operations.
			ent.OpCreate,
		),
	}
}
