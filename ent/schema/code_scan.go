package schema

import (
	"context"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	gen "github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/hook"
)

type CodeScan struct {
	ent.Schema
}

func (CodeScan) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "code_scan"},
	}
}

func (CodeScan) Fields() []ent.Field {
	return []ent.Field{
		field.Text("tool").Immutable().NotEmpty().
			Annotations(
				entgql.OrderField("tool"),
			),
	}
}

func (CodeScan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("release", Release.Type).Unique().Required(),
		edge.From("issues", CodeIssue.Type).Ref("scan"),
		edge.From("entry", ReleaseEntry.Type).Ref("code_scan").Unique(),
	}
}

func (CodeScan) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.CodeScanFunc(func(ctx context.Context, m *gen.CodeScanMutation) (ent.Value, error) {
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
