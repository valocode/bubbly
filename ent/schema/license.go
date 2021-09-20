package schema

import (
	"context"
	"fmt"
	"regexp"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	gen "github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
	"github.com/valocode/bubbly/ent/hook"
	"github.com/valocode/bubbly/ent/spdxlicense"
)

type License struct {
	ent.Schema
}

func (License) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "license"},
		entmodel.Annotation{},
	}
}

func (License) Fields() []ent.Field {
	return []ent.Field{
		field.String("license_id").NotEmpty().Unique().
			Match(regexp.MustCompile(`^[A-Za-z0-9-_\.\+]+$`)).
			Annotations(
				entgql.OrderField("license_id"),
			),
		field.String("name").Optional().
			Annotations(
				entgql.OrderField("name"),
			),
	}
}

func (License) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owner", Organization.Type).Unique().Required(),
		edge.To("spdx", SPDXLicense.Type).Unique(),
		edge.From("components", Component.Type).Ref("licenses"),
		edge.From("instances", ReleaseLicense.Type).Ref("license"),
	}
}

func (License) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			matchLicenseIDToSPDX,
			// Limit the hook only for these operations.
			ent.OpCreate|ent.OpUpdate,
		),
	}
}

// matchLicenseIDToSPDX checks if the licenseID matches an SPDX license ID.
// If it does match, then it creates an edge between the licenses
func matchLicenseIDToSPDX(next ent.Mutator) ent.Mutator {
	return hook.LicenseFunc(func(ctx context.Context, m *gen.LicenseMutation) (ent.Value, error) {

		licenseID, ok := m.LicenseID()
		if !ok {
			return nil, fmt.Errorf("license must have a license id")
		}
		// If a current SPDX mapping exists, clear it.
		// If the licenseID did not change it does not matter as it will get reset anyway.
		// This ensures we do not leave dangling mappings if the licenseID were to change.
		if _, ok := m.SpdxID(); ok {
			m.ClearSpdx()
		}

		client := m.Client()
		spdxID, err := client.SPDXLicense.Query().
			Where(spdxlicense.LicenseID(licenseID)).
			OnlyID(ctx)
		if err != nil {
			// If not found, just continue
			if gen.IsNotFound(err) {
				return next.Mutate(ctx, m)
			}
			return nil, fmt.Errorf("checking for matching SPDX license: %w", err)
		}

		m.SetSpdxID(spdxID)

		return next.Mutate(ctx, m)
	})
}
