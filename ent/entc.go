//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"

	"github.com/valocode/bubbly/ent/extensions/entfilter"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
	"github.com/valocode/bubbly/ent/extensions/entts"
)

func main() {

	ex, err := entgql.NewExtension(
		entgql.WithWhereFilters(true),
		entgql.WithMapScalarFunc(func(f *gen.Field, o gen.Op) string {
			if f.StructField() == "Metadata" {
				return "Map"
			}
			return ""
		}),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}

	if err := entc.Generate("./schema",
		&gen.Config{
			Features: []gen.Feature{
				// gen.FeaturePrivacy,
				// gen.FeatureUpsert,
			},
			// Templates: templates,
		},
		entc.Extensions(
			ex,
			entfilter.NewExtension(),
			entmodel.NewExtension(),
			entts.NewExtension(),
		),
	); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
