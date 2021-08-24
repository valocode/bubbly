package schema

import (
	"encoding/json"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/valocode/bubbly/ent/extensions/entmodel"
)

type Adapter struct {
	ent.Schema
}

func (Adapter) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "adapter"},
		// Skip graphql generation for the Adapter because json.RawMessage and []byte types
		// cause problems...
		entgql.Skip(),
		entmodel.Annotation{},
	}
}

func (Adapter) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.String("tag").NotEmpty().
			Annotations(
				entgql.OrderField("tag"),
			),
		field.Enum("type").Immutable().
			Values("json", "csv", "xml", "yaml", "http").
			Annotations(
				entgql.OrderField("type"),
			),
		field.JSON("operation", json.RawMessage{}),
		field.Enum("results_type").Immutable().
			Values("code_scan", "test_run").
			Annotations(
				entgql.OrderField("results_type"),
			),
		field.Bytes("results"),
	}
}

func (Adapter) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "tag").
			Unique(),
	}
}
