package schema

import (
	"context"
	"fmt"
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	gen "github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
	"github.com/valocode/bubbly/ent/hook"
	"github.com/valocode/bubbly/ent/license"
)

type SPDXLicense struct {
	ent.Schema
}

func (SPDXLicense) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "spdx_license"},
		entmodel.Annotation{},
	}
}

func (SPDXLicense) Fields() []ent.Field {
	return []ent.Field{
		// License ID is the SPDX ID: https://spdx.dev/ids/
		field.String("license_id").NotEmpty().Unique().
			Match(regexp.MustCompile(`^[A-Za-z0-9-_\.\+]+$`)),
		field.String("name").NotEmpty(),
		// reference points to a url where more information is available
		field.String("reference").Optional(),
		field.String("details_url").Optional(),
		field.Bool("is_osi_approved").Default(false),
	}
}

func (SPDXLicense) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			matchSPDXToLicenseID,
			// Limit the hook only for these operations.
			ent.OpCreate|ent.OpUpdate,
		),
	}
}

// matchSPDXToLicenseID checks if the spdx license ID matches a license ID.
// This is so that if new SPDX licenses are added *after* a license ID is added
// it will make the connection
func matchSPDXToLicenseID(next ent.Mutator) ent.Mutator {
	return hook.SPDXLicenseFunc(func(ctx context.Context, m *gen.SPDXLicenseMutation) (ent.Value, error) {

		spdxID, ok := m.LicenseID()
		if !ok {
			return nil, fmt.Errorf("spdx license must have a license id")
		}

		client := m.Client()
		// Get the license that matches this SPDX ID
		dbLicense, err := client.License.Query().
			Where(license.LicenseID(spdxID)).
			WithSpdx().
			Only(ctx)
		if err != nil {
			// If not found
			if gen.IsNotFound(err) {
				return next.Mutate(ctx, m)
			}
			return nil, fmt.Errorf("checking for matching license from SPDX license ID: %w", err)
		}

		if dbLicense.Edges.Spdx != nil {
			if dbLicense.Edges.Spdx.LicenseID == spdxID {
				// If it's already matching, then do nothing
				return next.Mutate(ctx, m)
			}
			// If it did not match, then something is a bit strange. Either the
			// licenseID changed and did not remove the edge, or the SPDX ID
			// has changed. Either way, this SPDX ID not matches the license ID
			// and we will update it below
		}
		// Continue the operation so that we ensure we have a DB ID
		val, err := next.Mutate(ctx, m)
		if err != nil {
			return nil, fmt.Errorf("saving SPDX license: %w", err)
		}

		// Not possible for there not to be an ID after the operation has completed
		dbID, _ := m.ID()
		if _, err := client.License.UpdateOne(dbLicense).
			SetSpdxID(dbID).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("clearing SPDX from license: %w", err)
		}
		return val, nil
	})
}
