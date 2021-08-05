package extensions

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

const (
	tsSchemaFile = "../typescript/schema_gen.ts"
)

type (
	Extension struct {
		entc.DefaultExtension
	}
)

func NewExtension() (*Extension, error) {
	ex := Extension{}

	return &ex, nil
}

// Hooks of the extension.
func (e *Extension) Hooks() []gen.Hook {
	return []gen.Hook{
		func(next gen.Generator) gen.Generator {
			return gen.GenerateFunc(func(g *gen.Graph) error {
				if err := next.Generate(g); err != nil {
					return err
				}

				return genTypescriptInterfaces(g)
			})
		},
	}
}

func genTypescriptInterfaces(graph *gen.Graph) error {

	var b bytes.Buffer
	fmt.Fprintf(&b, "// #######################################\n")
	fmt.Fprintf(&b, "// Code is generated using a custom ent extension.\n")
	fmt.Fprintf(&b, "// DO NOT MODIFY.\n")
	fmt.Fprintf(&b, "// Currently it is manually copied over from the bubbly repository where it is generated.\n")
	fmt.Fprintf(&b, "// #######################################\n\n")
	for _, node := range graph.Nodes {

		fmt.Fprintf(&b, "// #######################################\n")
		fmt.Fprintf(&b, "// %s\n", node.Name)
		fmt.Fprintf(&b, "// #######################################\n")

		//
		// Create a JSON wrapper for the interface
		//
		fmt.Fprintf(&b, "export interface %s {\n", node.Name+"_Json")
		fmt.Fprintf(&b, "\t%s?: %s;\n", node.Table(), node.Name+"[]")
		fmt.Fprintf(&b, "}\n\n")

		//
		// Create the interface definition
		//
		fmt.Fprintf(&b, "export interface %s {\n", node.Name)
		fmt.Fprintf(&b, "\tid?: %s;\n", entFieldToTSType(node, node.ID))
		for _, field := range node.Fields {
			fmt.Fprintf(&b, "\t%s?: %s;\n", field.Name, entFieldToTSType(node, field))
		}
		for _, edge := range node.Edges {
			fmt.Fprintf(&b, "\t%s?: %s;\n", edge.Name, entEdgeToTSType(edge))
		}
		fmt.Fprintf(&b, "}\n\n")

		//
		// Create a relay spec interface
		//
		fmt.Fprintf(&b, "export interface %s {\n", node.Name+"_Relay")
		fmt.Fprintf(&b, "\t%s_connection?: %s_Conn;\n", node.Table(), node.Name)
		fmt.Fprintf(&b, "}\n\n")

		fmt.Fprintf(&b, "export interface %s {\n", node.Name+"_Conn")
		fmt.Fprintf(&b, "\ttotalCount?: number;\n")
		fmt.Fprintf(&b, "\tpageInfo?: pageInfo;\n")
		fmt.Fprintf(&b, "\tedges?: %s[];\n", node.Name+"_Edge")
		fmt.Fprintf(&b, "}\n\n")

		//
		// Create a relay spec edge
		//
		fmt.Fprintf(&b, "export interface %s {\n", node.Name+"_Edge")
		fmt.Fprintf(&b, "\tnode?: %s;\n", node.Name)
		fmt.Fprintf(&b, "}\n\n")

		//
		// Create the enum definitions
		//
		for _, enum := range node.EnumFields() {
			fmt.Fprintf(&b, "export enum %s {\n", node.Name+strings.Title(enum.Name))
			for _, ev := range enum.EnumValues() {
				fmt.Fprintf(&b, "\t%s = \"%s\",\n", ev, ev)
			}
			fmt.Fprintf(&b, "}\n\n")
		}
	}

	//
	// Create the relay spec page_info
	//
	fmt.Fprintf(&b, "export interface pageInfo {\n")
	fmt.Fprintf(&b, "\thasNextPage?: boolean;\n")
	fmt.Fprintf(&b, "\thasPreviousPage?: boolean;\n")
	fmt.Fprintf(&b, "\tstartCursor?: string;\n")
	fmt.Fprintf(&b, "\tendCursor?: string;\n")
	fmt.Fprintf(&b, "}\n\n")

	err := os.WriteFile(tsSchemaFile, b.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", tsSchemaFile, err)
	}

	return nil
}

func entFieldToTSType(n *gen.Type, f *gen.Field) string {
	// Handle all numerics nice and easy
	if f.Type.Numeric() {
		return "number"
	}
	switch f.Type.Type {
	case field.TypeBool:
		return "boolean"
	case field.TypeString:
		return "string"
	case field.TypeEnum:
		return n.Name + strings.Title(f.Name)
	case field.TypeTime:
		return "Date"
	}
	return "UNKNOWN_TYPE_" + f.Type.String()
}

func entEdgeToTSType(e *gen.Edge) string {
	var tsStr string

	tsStr += e.Type.Name
	if !e.Unique {
		tsStr += "[]"
	}
	return tsStr
}
