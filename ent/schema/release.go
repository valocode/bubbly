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

type Release struct {
	ent.Schema
}

func (Release) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "release"},
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
			Default("pending"),
	}
}

func (Release) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("dependencies", Release.Type).From("subreleases"),
		edge.To("project", Project.Type).Unique().Required(),
		edge.From("commit", GitCommit.Type).Ref("release").Unique().Required(),
		edge.From("artifacts", Artifact.Type).Ref("release"),
		edge.From("checks", ReleaseCheck.Type).Ref("release"),
		edge.From("log", ReleaseEntry.Type).Ref("release"),
		edge.From("code_scans", CodeScan.Type).Ref("release"),
		edge.From("cve_scans", CVEScan.Type).Ref("release"),
		edge.From("license_scans", LicenseScan.Type).Ref("release"),
		edge.From("test_runs", TestRun.Type).Ref("release"),
		edge.From("components", Component.Type).Ref("release"),
		// edge.From("licenses", LicenseUsage.Type).Ref("release"),
		// edge.From("vulnerabilities", Vulnerability.Type).Ref("release"),
	}
}

func (Release) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "version").
			Edges("project", "commit").
			Unique(),
	}
}
