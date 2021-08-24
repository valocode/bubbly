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
	// var templates []*gen.Template
	// templates = append(
	// 	templates,
	// 	gen.MustParse(gen.NewTemplate("static").
	// 		ParseDir("templates")),
	// )

	ex, err := entgql.NewExtension(
		entgql.WithSchemaPath("../gql/ent.graphql"),
		entgql.WithConfigPath("../gql/gqlgen.yaml"),
		entgql.WithCustomRelaySpec(true, func(name string) string {
			return name + "_connection"
		}),
		entgql.WithNaming("snake"),
		entgql.WithWhereFilters(true),
		entgql.WithOrderBy(true),
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
