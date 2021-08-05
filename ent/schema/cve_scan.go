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

type CVEScan struct {
	ent.Schema
}

func (CVEScan) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "cve_scan"},
	}
}

func (CVEScan) Fields() []ent.Field {
	return []ent.Field{
		field.Text("tool").Immutable().NotEmpty().
			Annotations(
				entgql.OrderField("tool"),
			),
	}
}

func (CVEScan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("release", Release.Type).Unique().Required(),
		edge.From("entry", ReleaseEntry.Type).Ref("cve_scan").Unique(),
		edge.From("vulnerabilities", Vulnerability.Type).Ref("scan"),
	}
}

func (CVEScan) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.CVEScanFunc(func(ctx context.Context, m *gen.CVEScanMutation) (ent.Value, error) {
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
