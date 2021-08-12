package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
)

type Release struct {
	ent.Schema
}

func (Release) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "release"},
		entmodel.Annotation{},
	}
}

func (Release) Fields() []ent.Field {
	return []ent.Field{
		field.Text("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.Text("version").NotEmpty().
			Annotations(
				entgql.OrderField("version"),
			),
		field.Enum("status").
			Values("pending", "ready", "blocked").
			Default("pending").Annotations(
		// entmodel.Annotation{SkipCreate: true},
		),
	}
}

func (Release) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("dependencies", Release.Type).From("subreleases"),
		edge.From("commit", GitCommit.Type).Ref("release").Unique().Required(),
		edge.From("log", ReleaseEntry.Type).Ref("release"),
		edge.From("artifacts", Artifact.Type).Ref("release"),
		edge.From("components", ReleaseComponent.Type).Ref("release"),
		edge.From("vulnerabilities", ReleaseVulnerability.Type).Ref("release"),
		// edge.From("licenses", LicenseUse.Type).Ref("release"),
		edge.From("code_scans", CodeScan.Type).Ref("release"),
		edge.From("test_runs", TestRun.Type).Ref("release"),

		edge.From("vulnerability_reviews", VulnerabilityReview.Type).Ref("releases"),
	}
}

func (Release) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("commit").
			Unique(),
	}
}
