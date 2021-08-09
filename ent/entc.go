// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"

	"github.com/valocode/bubbly/ent/extensions"
	"github.com/valocode/bubbly/ent/extensions/entapt"
)

func main() {
	var templates []*gen.Template
	templates = append(
		templates,
		gen.MustParse(gen.NewTemplate("static").
			ParseDir("templates")),
	)

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

	tsex, err := extensions.NewTSExtension()
	if err != nil {
		log.Fatalf("creating tsmodel extension: %v", err)
	}
	adapterExt, err := entapt.NewExtension()
	if err != nil {
		log.Fatalf("creating adapter extension: %v", err)
	}

	if err := entc.Generate("./schema",
		&gen.Config{
			Features: []gen.Feature{
				// gen.FeaturePrivacy,
				// gen.FeatureUpsert,
			},
			Templates: templates,
		},
		entc.Extensions(ex, tsex, adapterExt),
	); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
