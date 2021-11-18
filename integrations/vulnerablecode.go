package integrations

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasecomponent"
	"github.com/valocode/bubbly/store"
)

type (
	VulnerableCodeMonitor struct {
		ctx   context.Context
		store *store.Store
	}
)

func (m *VulnerableCodeMonitor) Do() error {

	releaseIDs, err := m.store.Client().Release.Query().Where(
		func(s *sql.Selector) {
			s.Where(
				sqljson.ValueContains(release.FieldLabels, "true", sqljson.Path("bubbly/release/active")),
			)
		},
	).IDs(m.ctx)

	dbComponents, err := m.store.Client().Component.Query().
		Where(
			component.HasUsesWith(
				releasecomponent.HasReleaseWith(
					release.IDIn(releaseIDs...),
				),
			),
		).Unique(true).
		All(m.ctx)

	// TODO: For each dbComponent... check vulnerable code
}
