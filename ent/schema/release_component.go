package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type ReleaseComponent struct {
	ent.Schema
}

func (ReleaseComponent) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "release_component"},
	}
}

func (ReleaseComponent) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").
			Comment(`
The type indicates how the component is used in the project,
e.g. whether it is embedded into the build (static link) or just
distributed (dynamic link) or just a development dependency`,
			).
			Values("embedded", "distributed", "development").
			Default("embedded"),
	}
}

func (ReleaseComponent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("release", Release.Type).Required().Unique(),
		edge.To("scans", CodeScan.Type).Required(),
		edge.To("component", Component.Type).Required().Unique(),
		edge.From("vulnerabilities", ReleaseVulnerability.Type).Ref("component"),
		edge.From("licenses", ReleaseLicense.Type).Ref("component"),
	}
}

func (ReleaseComponent) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Edges("release", "component").
			Unique(),
	}
}

// func (ReleaseComponent) Hooks() []ent.Hook {
// 	return []ent.Hook{
// 		hook.On(
// 			inheritVulnerabiliesAndReviews,
// 			ent.OpCreate,
// 		),
// 		// hook.On(
// 		// 	inheritVulnerabiliesAndReviews,
// 		// 	ent.OpCreate,
// 		// ),
// 	}
// }

// func inheritVulnerabiliesAndReviews(next ent.Mutator) ent.Mutator {
// 	return hook.ReleaseComponentFunc(func(ctx context.Context, m *gen.ReleaseComponentMutation) (ent.Value, error) {
// 		client := clientOrTxClient(m)
// 		rID, hasRelease := m.ReleaseID()
// 		cpID, hasComponent := m.ComponentID()
// 		if !hasComponent || !hasRelease {
// 			return nil, errors.New("release component must have a component and a release")
// 		}
// 		// Get the release with it's project and repo
// 		rQuery := client.Release.Query().
// 			Where(release.IDEQ(rID)).QueryCommit().QueryRepository()
// 		repoID, err := rQuery.OnlyID(ctx)
// 		if err != nil {
// 			return nil, fmt.Errorf("error getting repo from release component: %w", err)
// 		}
// 		// TODO: is it one project, or projects
// 		projectIDs, err := rQuery.QueryProject().IDs(ctx)
// 		if err != nil {
// 			return nil, fmt.Errorf("error getting projects from release component: %w", err)
// 		}

// 		// Get the list of vulnerability reviews associated with the repo
// 		// and project(s) that this release component is associated with
// 		vulns, err := client.Component.Query().
// 			Where(component.IDEQ(cpID)).
// 			QueryVulnerabilities().
// 			WithReviews(func(vrq *gen.VulnerabilityReviewQuery) {
// 				vrq.Where(vulnerabilityreview.Or(
// 					vulnerabilityreview.HasProjectsWith(project.IDIn(projectIDs...)),
// 					vulnerabilityreview.HasReposWith(repo.IDEQ(repoID)),
// 				))
// 			}).
// 			All(ctx)
// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, vuln := range vulns {
// 			relVuln, err := client.ReleaseVulnerability.Create().
// 				SetReleaseID(rID).
// 				SetVulnerability(vuln).
// 				// Add also reviews that came with it
// 				AddReviews(vuln.Edges.Reviews...).
// 				Save(ctx)
// 			if err != nil {
// 				return nil, fmt.Errorf("error creating implicit release vulnerability from component vulnerability: %w", err)
// 			}
// 			// Add the release vulnerability to this release component
// 			m.AddVulnerabilityIDs(relVuln.ID)
// 		}
// 		return next.Mutate(ctx, m)
// 	})
// }
