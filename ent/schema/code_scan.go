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

type CodeScan struct {
	ent.Schema
}

func (CodeScan) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "code_scan"},
		entmodel.Annotation{},
	}
}

func (CodeScan) Fields() []ent.Field {
	return []ent.Field{
		field.Text("tool").Immutable().NotEmpty().
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

func (CodeScan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("release", Release.Type).Unique().Required(),
		edge.From("entry", ReleaseEntry.Type).Ref("code_scan").Unique(),
		// Different results from the scan
		edge.From("issues", CodeIssue.Type).Ref("scan"),
		edge.From("vulnerabilities", ReleaseVulnerability.Type).Ref("scans"),
		// edge.From("licenses", LicenseUse.Type).Ref("scans"),
		edge.From("components", ReleaseComponent.Type).Ref("scans"),
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
