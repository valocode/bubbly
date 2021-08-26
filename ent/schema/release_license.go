package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
)

type ReleaseLicense struct {
	ent.Schema
}

func (ReleaseLicense) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "release_license"},
	}
}

func (ReleaseLicense) Fields() []ent.Field {
	return []ent.Field{}
}

func (ReleaseLicense) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge to the actual vulnerability for this
		edge.To("license", License.Type).Unique().Required(),
		edge.To("component", ReleaseComponent.Type).Unique(),
		edge.To("release", Release.Type).Unique().Required(),
		// TODO: edge.To("reviews", LicenseReview.Type),
		// Scan is optional and only used if the license was identified by
		// a tool integrated with an adapter
		edge.To("scans", CodeScan.Type),
	}
}
