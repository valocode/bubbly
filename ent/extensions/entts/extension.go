package entts

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

const (
	tsSchemaFile = "../typescript/schema_gen.ts"
)

type (
	Extension struct {
		entc.DefaultExtension
	}
)

func NewExtension() *Extension {
	return &Extension{}
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

// filterNodes - copied from ent/contrib/entgql...
func filterNodes(nodes []*gen.Type) ([]*gen.Type, error) {
	var filteredNodes []*gen.Type
	for _, n := range nodes {
		ant := &entgql.Annotation{}
		if n.Annotations != nil && n.Annotations[ant.Name()] != nil {
			if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
				return nil, err
			}
			if ant.Skip {
				continue
			}
		}
		filteredNodes = append(filteredNodes, n)
	}
	return filteredNodes, nil
}

func genTypescriptInterfaces(graph *gen.Graph) error {

	var b bytes.Buffer
	fmt.Fprintf(&b, "// #######################################\n")
	fmt.Fprintf(&b, "// Code is generated using a custom ent extension.\n")
	fmt.Fprintf(&b, "// DO NOT MODIFY.\n")
	fmt.Fprintf(&b, "// Currently it is manually copied over from the bubbly repository where it is generated.\n")
	fmt.Fprintf(&b, "// #######################################\n\n")
	nodes, err := filterNodes(graph.Nodes)
	if err != nil {
		return err
	}
	for _, node := range nodes {

		fmt.Fprintf(&b, "// #######################################\n")
		fmt.Fprintf(&b, "// %s\n", node.Name)
		fmt.Fprintf(&b, "// #######################################\n")

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

	if err := os.WriteFile(tsSchemaFile, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing to file %s: %w", tsSchemaFile, err)
	}

	return nil
}
