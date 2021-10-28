package schema

import (
	"context"
	"fmt"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	gen "github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/hook"
	types "github.com/valocode/bubbly/ent/schema/types"
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
		field.String("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.String("version").NotEmpty().
			Annotations(
				entgql.OrderField("version"),
			),

		field.JSON("labels", types.Labels{}).Optional().
			Annotations(
				entgql.Skip(),
			),
	}
}

func (Release) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("dependencies", Release.Type).From("subreleases"),
		edge.From("commit", GitCommit.Type).Ref("release").Unique().Required(),
		edge.From("head_of", Repo.Type).Ref("head").Unique(),
		edge.From("log", ReleaseEntry.Type).Ref("release"),
		edge.From("violations", ReleasePolicyViolation.Type).Ref("release"),
		edge.From("artifacts", Artifact.Type).Ref("release"),
		edge.From("components", ReleaseComponent.Type).Ref("release"),
		edge.From("vulnerabilities", ReleaseVulnerability.Type).Ref("release"),
		edge.From("licenses", ReleaseLicense.Type).Ref("release"),
		edge.From("code_scans", CodeScan.Type).Ref("release"),
		edge.From("test_runs", TestRun.Type).Ref("release"),

		edge.From("vulnerability_reviews", VulnerabilityReview.Type).Ref("releases"),
	}
}

func (Release) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			updateReleaseHead,
			// Limit the hook only for these operations.
			ent.OpCreate,
		),
	}
}

func (Release) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("commit").
			Unique(),
	}
}

func updateReleaseHead(next ent.Mutator) ent.Mutator {
	return hook.ReleaseFunc(func(ctx context.Context, m *gen.ReleaseMutation) (ent.Value, error) {
		cID, ok := m.CommitID()
		if !ok {
			return nil, fmt.Errorf("release does not have a commit")
		}
		client := m.Client()
		dbCommit, err := client.GitCommit.Query().
			Where(gitcommit.ID(cID)).
			WithRepo(func(rq *gen.RepoQuery) {
				rq.WithHead(func(rq *gen.ReleaseQuery) {
					rq.WithCommit()
				})
			}).Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("error quering commit and repo data for release: %w", err)
		}
		dbRepo := dbCommit.Edges.Repo
		if dbRepo == nil {
			return nil, fmt.Errorf("release commit does not have a repo")
		}
		dbHead := dbRepo.Edges.Head
		if dbHead == nil {
			// If not found, then the repo does not have a head yet, so assign this
			m.SetHeadOfID(dbRepo.ID)
			if dbCommit.Branch == "master" {
				_, err := client.Repo.UpdateOne(dbRepo).
					SetDefaultBranch("master").
					Save(ctx)
				if err != nil {
					return nil, fmt.Errorf("error setting default branch for repo: %w", err)
				}
			}
			return next.Mutate(ctx, m)
		}
		// If there is a head already, check if this release should become
		// the new head
		if dbCommit.Branch == dbRepo.DefaultBranch {
			if dbCommit.Time.After(dbHead.Edges.Commit.Time) {
				m.SetHeadOfID(dbRepo.ID)
				// Clear the existing head
				_, err := client.Release.UpdateOne(dbHead).ClearHeadOf().Save(ctx)
				if err != nil {
					return nil, err
				}
			}
		}

		return next.Mutate(ctx, m)
	})
}
