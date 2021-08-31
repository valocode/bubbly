package schema

import (
	"context"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	gen "github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
	"github.com/valocode/bubbly/ent/hook"
	types "github.com/valocode/bubbly/ent/schema/types"
)

type Artifact struct {
	ent.Schema
}

func (Artifact) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "artifact"},
		entmodel.Annotation{},
	}
}

func (Artifact) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Immutable().NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.String("sha256").Immutable().NotEmpty().
			Annotations(
				entgql.OrderField("sha256"),
			),
		field.Enum("type").Immutable().
			Values("docker", "file").
			Annotations(
				entgql.OrderField("type"),
			),
		field.Time("time").
			Immutable().
			Default(func() time.Time { return time.Now() }).
			Annotations(
				entgql.OrderField("time"),
			),
		field.JSON("metadata", types.Metadata{}).Optional(),
	}
}

func (Artifact) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("release", Release.Type).Unique(),
		edge.From("entry", ReleaseEntry.Type).Ref("artifact").Unique(),
	}
}

func (Artifact) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sha256").
			Unique(),
	}
}

func (Artifact) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.ArtifactFunc(func(ctx context.Context, m *gen.ArtifactMutation) (ent.Value, error) {
					if err := createReleaseEntry(ctx, m); err != nil {
						return nil, err
					}
					return next.Mutate(ctx, m)
				})
			},
			// Limit the hook only for these operations.
			ent.OpCreate,
		),
	}
}
