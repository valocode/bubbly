package entmodel

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/graphql-go/graphql"
)

const (
	schemaFile = "../gql/types.graphql"
)

func genGraphQLTypes(g *gen.Graph) error {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# #######################################\n")
	fmt.Fprintf(&b, "# Code is generated using a custom ent extension.\n")
	fmt.Fprintf(&b, "# DO NOT MODIFY.\n")
	fmt.Fprintf(&b, "# #######################################\n")

	fmt.Fprintf(&b, `
directive @goField(forceResolver: Boolean, name: String) on INPUT_FIELD_DEFINITION
	| FIELD_DEFINITION

scalar Map

"""
Maps a Time GraphQL scalar to a Go time.Time struct.
Generated by ent.
"""
scalar Time

"""
Define a Relay Cursor type:
https://relay.dev/graphql/connections.htm#sec-Cursor
Generated by ent.
"""
scalar Cursor

interface Node {
  id: ID!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: Cursor
  endCursor: Cursor
}

enum OrderDirection {
	ASC
	DESC
  }

`)

	nodes, err := filterNodes(g.Nodes)
	if err != nil {
		return err
	}
	for _, node := range nodes {

		fmt.Fprintf(&b, "\"\"\"\n")
		fmt.Fprintf(&b, "%s\n", node.Name)
		fmt.Fprintf(&b, "\"\"\"\n")

		//
		// Create the base GraphQL type
		//
		// TODO: once we remove the Base suffix, it should also implement the Noder
		// interface defined in ent
		// fmt.Fprintf(&b, "type %s implements Node {\n", node.Name)
		fmt.Fprintf(&b, "type %sBase {\n", node.Name)
		fmt.Fprintf(&b, "\t%s: ID!\n", node.ID.Name)
		for _, field := range node.Fields {
			fmt.Fprintf(&b, "\t%s: %s\n", field.Name, entFieldToGraphQL(node, field))
		}
		fmt.Fprintf(&b, "}\n\n")

		//
		// Create the ordering enum definition (if ordering fields exist)
		//
		orderFields := orderByFields(node)
		if len(orderFields) > 0 {
			fmt.Fprintf(&b, "enum %sOrderField {\n", node.Name)
			for _, ev := range orderFields {
				fmt.Fprintf(&b, "\t%s\n", ev)
			}
			fmt.Fprintf(&b, "}\n\n")

			fmt.Fprintf(&b, "input %sOrder {\n", node.Name)
			fmt.Fprintf(&b, "\tdirection: OrderDirection!\n")
			fmt.Fprintf(&b, "\tfield: %sOrderField!\n", node.Name)
			fmt.Fprintf(&b, "}\n\n")
		}

		//
		// Create any enum definitions
		//
		for _, enum := range node.EnumFields() {
			fmt.Fprintf(&b, "enum %s {\n", node.Name+strings.Title(enum.Name))
			for _, ev := range enum.EnumValues() {
				fmt.Fprintf(&b, "\t%s\n", ev)
			}
			fmt.Fprintf(&b, "}\n\n")
		}
	}

	if err := os.WriteFile(schemaFile, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing to file %s: %w", schemaFile, err)
	}

	return nil
}

func orderByFields(n *gen.Type) []string {
	var orderFields []string
	for _, f := range n.Fields {
		// Check if graphql ordering has been applied to this field through an annotation.
		// If it has, add it to the order by enum for this node.
		var ant entgql.Annotation
		if i, ok := f.Annotations[ant.Name()]; ok && ant.Decode(i) == nil && ant.OrderField != "" {
			orderFields = append(orderFields, f.Name)
		}
	}
	return orderFields
}

func entFieldToGraphQL(n *gen.Type, f *gen.Field) string {
	switch t := f.Type.Type; {
	case t.Integer():
		return graphql.Int.Name()
	case t.Float():
		return graphql.Float.Name()
	case t == field.TypeString:
		return graphql.String.Name()
	case f.IsBool():
		return "Boolean"
	case f.IsEnum():
		return n.Name + strings.Title(f.Name)
	case f.StructField() == "Metadata":
		return "Map"
	case f.IsTime():
		return "Time"
	}
	return "UNKNOWN_TYPE"
}
