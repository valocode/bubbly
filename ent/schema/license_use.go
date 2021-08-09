package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
)

type LicenseUse struct {
	ent.Schema
}

func (LicenseUse) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "license_use"},
	}
}

func (LicenseUse) Fields() []ent.Field {
	return []ent.Field{}
}

func (LicenseUse) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("license", License.Type).Unique().Required(),
		// edge.To("scans", CodeScan.Type).Required(),
		// edge.To("release", Release.Type).Unique().Required(),

		// edge.To("component", ComponentUse.Type).Unique().Required(),
	}
}

func (LicenseUse) Indexes() []ent.Index {
	return []ent.Index{}
}
