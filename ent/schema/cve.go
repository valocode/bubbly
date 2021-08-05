package schema

import (
	"context"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	gen "github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/cve"
	"github.com/valocode/bubbly/ent/hook"
)

type CVE struct {
	ent.Schema
}

func (CVE) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "cve"},
	}
}

func (CVE) Fields() []ent.Field {
	return []ent.Field{
		field.Text("cve_id").NotEmpty().
			Annotations(
				entgql.OrderField("cve_id"),
			),
		field.Text("description").Optional().
			Annotations(
				entgql.OrderField("description"),
			),
		field.Float("severity_score").
			Default(0.0).
			Annotations(
				entgql.OrderField("severity_score"),
			),
		field.Enum("severity").
			Default("None").
			Values("None", "Low", "Medium", "High", "Critical").
			Annotations(
				entgql.OrderField("severity"),
			),
		field.Time("published_data").Optional().
			Annotations(
				entgql.OrderField("published_data"),
			),
		field.Time("modified_data").Optional().
			Annotations(
				entgql.OrderField("modified_data"),
			),
	}
}

func (CVE) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("found", Vulnerability.Type).Ref("cve"),
		edge.From("rules", CVERule.Type).Ref("cve"),
	}
}

func (CVE) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("cve_id").Unique(),
	}
}

func (CVE) Hooks() []ent.Hook {
	return []ent.Hook{
		// Automatically set the CVE severity if it has not already been set
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.CVEFunc(func(ctx context.Context, m *gen.CVEMutation) (ent.Value, error) {
					if score, exists := m.SeverityScore(); exists {
						// Ranges: https://nvd.nist.gov/vuln-metrics/cvss
						var severity cve.Severity
						switch {
						case score == 0.0:
							severity = cve.SeverityNone
						case score >= 0.1 && score < 4.0:
							severity = cve.SeverityLow
						case score >= 4.0 && score < 7.0:
							severity = cve.SeverityMedium
						case score >= 7.0 && score < 9.0:
							severity = cve.SeverityHigh
						case score >= 9.0:
							severity = cve.SeverityCritical
						default:
							severity = cve.SeverityNone
						}
						m.SetSeverity(severity)
					}
					return next.Mutate(ctx, m)
				})
			},
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),
	}
}
