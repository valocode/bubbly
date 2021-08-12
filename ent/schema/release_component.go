package schema

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"
	gen "github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/hook"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/repo"
	"github.com/valocode/bubbly/ent/vulnerabilityreview"
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
	return []ent.Field{}
}

func (ReleaseComponent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("release", Release.Type).Required().Unique(),
		edge.To("scans", CodeScan.Type).Required(),
		edge.To("component", Component.Type).Required().Unique(),
		edge.From("vulnerabilities", ReleaseVulnerability.Type).Ref("components"),
	}
}

func (ReleaseComponent) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Edges("release", "component").
			Unique(),
	}
}

func (ReleaseComponent) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.ReleaseComponentFunc(func(ctx context.Context, m *gen.ReleaseComponentMutation) (ent.Value, error) {
					client := clientOrTxClient(m)
					rID, hasRelease := m.ReleaseID()
					cpID, hasComponent := m.ComponentID()
					if !hasComponent || !hasRelease {
						return nil, errors.New("release component must have a component and a release")
					}
					// Get the release with it's project and repo
					rQuery := client.Release.Query().
						Where(release.IDEQ(rID)).QueryCommit().QueryRepo()
					repoID, err := rQuery.OnlyID(ctx)
					if err != nil {
						return nil, fmt.Errorf("error getting repo from release component: %w", err)
					}
					projectIDs, err := rQuery.QueryProjects().IDs(ctx)
					if err != nil {
						return nil, fmt.Errorf("error getting projects from release component: %w", err)
					}

					vulns, err := client.Component.Query().
						Where(component.IDEQ(cpID)).
						QueryVulnerabilities().
						WithReviews(func(vrq *gen.VulnerabilityReviewQuery) {
							vrq.Where(vulnerabilityreview.Or(
								vulnerabilityreview.HasProjectsWith(project.IDIn(projectIDs...)),
								vulnerabilityreview.HasReposWith(repo.IDEQ(repoID)),
							))
						}).
						All(ctx)
					if err != nil {
						return nil, err
					}

					for _, vuln := range vulns {
						relVuln, err := client.ReleaseVulnerability.Create().
							SetReleaseID(rID).
							SetVulnerability(vuln).
							// Add also reviews that came with it
							AddReviews(vuln.Edges.Reviews...).
							Save(ctx)
						if err != nil {
							return nil, fmt.Errorf("error creating implicit release vulnerability from component vulnerability: %w", err)
						}
						// Add the release vulnerability to this release component
						m.AddVulnerabilityIDs(relVuln.ID)
					}
					return next.Mutate(ctx, m)
				})
			},
			// Limit the hook only for these operations.
			ent.OpCreate,
		),
	}
}
