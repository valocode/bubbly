package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"
)

// ComponentInclude - component_include
// ComponentDependency - component_dep
// ComponentUse - component_use
type ComponentUse struct {
	ent.Schema
}

func (ComponentUse) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "component_use"},
	}
}

func (ComponentUse) Fields() []ent.Field {
	return []ent.Field{}
}

func (ComponentUse) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("release", Release.Type).Required().Unique(),
		edge.To("scans", CodeScan.Type).Required(),
		edge.To("component", Component.Type).Required().Unique(),

		// edge.From("vulnerabilities", Vulnerability.Type).Ref("component"),
		// edge.From("licenses", LicenseUse.Type).Ref("component"),
		// edge.To("licenses", License.Type),
		// edge.To("vulnerabilities", Vulnerability.Type),
		// edge.To("release", Release.Type),
	}
}

func (ComponentUse) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Edges("release", "component").
			Unique(),
	}
}
