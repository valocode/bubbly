package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
)

type LicenseUsage struct {
	ent.Schema
}

func (LicenseUsage) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "license_usage"},
	}
}

func (LicenseUsage) Fields() []ent.Field {
	return []ent.Field{}
}

func (LicenseUsage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("license", License.Type).Unique().Required(),
		edge.To("scan", LicenseScan.Type).Unique().Required(),
	}
}

func (LicenseUsage) Indexes() []ent.Index {
	return []ent.Index{}
}
