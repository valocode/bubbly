package schema

import (
	"regexp"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/valocode/bubbly/ent/extensions/entmodel"
)

type ReleasePolicy struct {
	ent.Schema
}

func (ReleasePolicy) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "release_policy"},
		entmodel.Annotation{},
	}
}

func (ReleasePolicy) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique().
			Match(regexp.MustCompile("^[a-z0-9_]+$")).
			Annotations(
				entgql.OrderField("name"),
			),
		field.String("module").NotEmpty().
			Comment("module stores the rego module defining the violation rules"),
		// TODO: compiling the rego modules become more complicted with
		// all the custom functions
		// Validate(func(s string) error {
		// 	// Validate that the provided module can be compiled.
		// 	// There might be a performance hit with doing this, but should
		// 	// ensure that all modules going into bubbly do compile.
		// 	// Creating modules is also not something time-critical in terms
		// 	// of CRUD operations
		// 	_, err := ast.CompileModules(map[string]string{
		// 		"module": s,
		// 	})
		// 	if err != nil {
		// 		return fmt.Errorf("error parsing rego module: %w", err)
		// 	}
		// 	return nil
		// }),
	}
}

func (ReleasePolicy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owner", Organization.Type).Unique().Required(),

		edge.From("violations", ReleasePolicyViolation.Type).Ref("policy"),
	}
}
