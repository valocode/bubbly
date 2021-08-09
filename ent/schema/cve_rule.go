package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type CVERule struct {
	ent.Schema
}

func (CVERule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "cve_rule"},
	}
}

func (CVERule) Fields() []ent.Field {
	return []ent.Field{
		field.Text("name").Optional().
			Annotations(
				entgql.OrderField("name"),
			),
	}
}

func (CVERule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cve", CVE.Type).Unique().Required(),
		edge.To("project", Project.Type),
		edge.To("repo", Repo.Type),
	}
}

// func (CVERule) Indexes() []ent.Index {
// 	return []ent.Index{
// 		index.Fields("name").
// 			Edges("cve", "project", "repo").
// 			Unique(),
// 	}
// }

func (CVERule) Hooks() []ent.Hook {
	return []ent.Hook{
		// hook.On(
		// 	func(next ent.Mutator) ent.Mutator {
		// 		return hook.CVERuleFunc(func(ctx context.Context, m *gen.CVERuleMutation) (ent.Value, error) {
		// 			if n, ok := m.Name(); ok {
		// 				if n == "" {
		// 					if cveID, ok := m.CveID(); ok {
		// 						cve, err := m.Client().CVE.Get(ctx, cveID)
		// 						if err == nil {
		// 							m.SetName("Rule for " + cve.CveID)
		// 						}
		// 					}
		// 				}
		// 			}
		// 			return next.Mutate(ctx, m)
		// 		})
		// 	},
		// 	ent.OpCreate,
		// ),
	}
}
