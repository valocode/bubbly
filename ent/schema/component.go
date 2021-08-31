package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/valocode/bubbly/ent/extensions/entmodel"
	types "github.com/valocode/bubbly/ent/schema/types"
)

type Component struct {
	ent.Schema
}

func (Component) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "component"},
		entmodel.Annotation{},
	}
}

func (Component) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.String("vendor").Default("").
			Annotations(
				entgql.OrderField("vendor"),
			),
		field.String("version").NotEmpty().
			Annotations(
				entgql.OrderField("version"),
			),
		field.String("description").Optional(),
		field.String("url").Optional(),
		field.JSON("metadata", types.Metadata{}).Optional(),
	}
}

func (Component) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vulnerabilities", Vulnerability.Type),
		edge.To("licenses", License.Type),
		edge.From("uses", ReleaseComponent.Type).Ref("component"),
	}
}

func (Component) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "vendor", "version").
			Unique(),
	}
}
