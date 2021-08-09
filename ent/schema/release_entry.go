package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ReleaseEntry struct {
	ent.Schema
}

func (ReleaseEntry) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "release_entry"},
	}
}

func (ReleaseEntry) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").
			Values("artifact", "deploy", "code_scan", "test_run").
			Immutable().
			Annotations(
				entgql.OrderField("type"),
			),
		field.Time("time").
			Default(func() time.Time { return time.Now() }).
			Annotations(
				entgql.OrderField("time"),
			),
	}
}

func (ReleaseEntry) Edges() []ent.Edge {
	return []ent.Edge{
		// Edges to different "events" are option, but it must have at least one
		// of these
		edge.To("artifact", Artifact.Type).Unique(),
		edge.To("code_scan", CodeScan.Type).Unique(),
		edge.To("test_run", TestRun.Type).Unique(),
		// edge.To("cve_scan", CVEScan.Type).Unique(),
		// edge.To("license_scan", LicenseScan.Type).Unique(),
		// Edge to a release is required
		edge.To("release", Release.Type).Unique().Required(),
	}
}

func (ReleaseEntry) Hooks() []ent.Hook {
	return []ent.Hook{
		// func(next ent.Mutator) ent.Mutator {
		// 	return hook.ReleaseEntryFunc(func(ctx context.Context, m *gen.ReleaseEntryMutation) (ent.Value, error) {
		// 		// TODO: how to validate *exactly one* of these edges exists?
		// 		// _, hasArtifact := m.ArtifactID()
		// 		// _, hasCodeScan := m.CodeScanID()
		// 		// _, hasCVEScan := m.CveScanID()
		// 		// _, hasLicenseScan := m.LicenseScanID()
		// 		return next.Mutate(ctx, m)
		// 	})
		// },
	}
}
