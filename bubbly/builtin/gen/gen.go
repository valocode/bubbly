// This package generates go structs for the builtin bubbly schema
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/parser"
	"github.com/valocode/bubbly/store"
	"github.com/zclconf/go-cty/cty"
)

const (
	goSchemaFile    = "schema_gen.go"
	goTablesFile    = "tables_gen.go"
	tsSchemaFile    = "schema_gen.ts"
	packageName     = "builtin"
	importStatement = `
import (
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)	
`
)

func main() {
	tables, err := builtin.BuiltinSchema()
	if err != nil {
		color.Red(fmt.Errorf("error getting schema tables: %w", err).Error())
		os.Exit(1)
	}
	// Flatten the tables and add implicit joins
	tables = store.FlattenTables(tables, nil)
	graph, err := store.NewSchemaGraph(tables)
	if err != nil {
		color.Red(fmt.Errorf("error creating schema graph").Error())
		os.Exit(1)
	}

	if err := genTablesFromSchema(graph); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	color.Green("Go tables successfully written to %s", goTablesFile)

	if err := genStructsFromSchema(graph); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	color.Green("Go structs successfully written to %s", goSchemaFile)

	if err := genTSInterfaceFromSchema(graph); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	color.Green("Typescript interfaces successfully written to %s", tsSchemaFile)
}

func genTablesFromSchema(graph *store.SchemaGraph) error {

	var b bytes.Buffer
	// Use this comment when testing to ignore compile errors...
	// fmt.Fprintf(&b, "//+build ignore\n\n")
	fmt.Fprintf(&b, "package %s\n\n", packageName)
	fmt.Fprint(&b, importStatement)
	fmt.Fprintf(&b, "var BuiltinTables = []core.Table{\n\n")
	graph.Traverse(func(node *store.SchemaNode) error {
		var (
			table = node.Table
			// tableName = camelToPascal(table.Name)
		)
		fmt.Fprintf(&b, "// #######################################\n")
		fmt.Fprintf(&b, "// %s\n", strings.ToUpper(table.Name))
		fmt.Fprintf(&b, "// #######################################\n")
		fmt.Fprintf(&b, "\ttable(\"%s\",\n", table.Name)
		fmt.Fprintf(&b, "\t\tfields(\n")
		for _, field := range table.Fields {
			fmt.Fprintf(&b, "\t\t\tfield(\"%s\", %s, %t),\n", field.Name, ctyTypeToGoCtyString(field.Type), field.Unique)
		}
		fmt.Fprintf(&b, "\t\t),\n")
		fmt.Fprintf(&b, "\t\tjoins(\n")
		for _, join := range table.Joins {
			fmt.Fprintf(&b, "\t\t\tjoin(\"%s\", %t, %t),\n", join.Table, join.Single, join.Unique)
		}
		fmt.Fprintf(&b, "\t\t),\n")
		fmt.Fprintf(&b, "\t),\n")
		return nil
	})
	// Close the BuiltinTables var
	fmt.Fprintf(&b, "}\n\n")

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		return fmt.Errorf("error formatting generated code: %w", err)
	}
	err = os.WriteFile(goTablesFile, formatted, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", goTablesFile, err)
	}

	return nil
}

func genStructsFromSchema(graph *store.SchemaGraph) error {

	var b bytes.Buffer
	// Use this comment when testing to ignore compile errors...
	// fmt.Fprintf(&b, "//+build ignore\n\n")
	fmt.Fprintf(&b, "package %s\n\n", packageName)
	fmt.Fprintf(&b, "import (\"time\"\n\"github.com/valocode/bubbly/api/core\")\n\n")
	graph.Traverse(func(node *store.SchemaNode) error {
		var (
			table     = node.Table
			tableName = camelToPascal(table.Name)
		)
		fmt.Fprintf(&b, "// #######################################\n")
		fmt.Fprintf(&b, "// %s\n", strings.ToUpper(table.Name))
		fmt.Fprintf(&b, "// #######################################\n")
		fmt.Fprintf(&b, "type %s struct {\n", tableName)
		// Add the Data Block "metadata" fields, which indicate the table name
		// and policy
		fmt.Fprintf(&b, "\t%s\t%s\t`json:\"%s,omitempty\"`\n", builtin.DBlockTableName, "string", table.Name)
		fmt.Fprintf(&b, "\t%s\t%s\t`json:\"%s\"`\n", builtin.DBlockPolicyName, "core.DataBlockPolicy", "-")
		fmt.Fprintf(&b, "\t%s\t%s\t`json:\"%s\"`\n", builtin.DBlockJoins, "[]string", "-")
		for _, field := range table.Fields {
			fmt.Fprintf(&b, "\t%s\t%s\t`json:\"%s,omitempty\"`\n", camelToPascal(field.Name), ctyTypeToString(field.Type), field.Name)
		}
		for _, edge := range node.Edges {
			var (
				eTable = edge.Node.Table
				single = edge.Rel != store.OneToMany
			)
			fmt.Fprintf(&b, "\t%s\t%s\t`json:\"%s,omitempty\"`\n", camelToPascal(eTable.Name), joinToType(eTable.Name, single), eTable.Name)
		}

		fmt.Fprintf(&b, "}\n")

		// Create some wrappers for JSON
		fmt.Fprintf(&b, "type %s struct {\n", tableName+"_Wrap")
		fmt.Fprintf(&b, "\t%s\t%s\t`json:\"%s,omitempty\"`\n", tableName, "[]"+tableName, table.Name)
		fmt.Fprintf(&b, "}\n\n")
		return nil
	})

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		return fmt.Errorf("error formatting generated code: %w", err)
	}
	err = os.WriteFile(goSchemaFile, formatted, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", goSchemaFile, err)
	}

	return nil
}

func genTSInterfaceFromSchema(graph *store.SchemaGraph) error {

	var b bytes.Buffer
	graph.Traverse(func(node *store.SchemaNode) error {
		var (
			table = node.Table
		)
		fmt.Fprintf(&b, "// #######################################\n")
		fmt.Fprintf(&b, "// %s\n", strings.ToUpper(table.Name))
		fmt.Fprintf(&b, "// #######################################\n")
		fmt.Fprintf(&b, "export interface %s {\n", table.Name)
		for _, field := range table.Fields {
			fmt.Fprintf(&b, "\t%s?: %s;\n", field.Name, ctyTypeToTSString(field.Type))
		}
		for _, edge := range node.Edges {
			var (
				eTable = edge.Node.Table
				single = edge.Rel != store.OneToMany
			)
			fmt.Fprintf(&b, "\t%s?: %s;\n", eTable.Name, joinToTSType(eTable.Name, single))
		}

		fmt.Fprintf(&b, "}\n")

		// Create some wrappers for JSON
		fmt.Fprintf(&b, "export interface %s {\n", table.Name+"_wrap")
		fmt.Fprintf(&b, "\t%s?: %s;\n", table.Name, table.Name+"[]")
		fmt.Fprintf(&b, "}\n\n")
		return nil
	})

	err := os.WriteFile(tsSchemaFile, b.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", tsSchemaFile, err)
	}

	return nil
}

func ctyTypeToString(ty cty.Type) string {
	switch {
	case ty == cty.Bool:
		return "bool"
	case ty == cty.Number:
		return "int64"
	case ty == cty.String:
		return "string"
	case ty.IsObjectType():
		return "map[string]interface{}"
	case ty.IsMapType():
		return "map[string]interface{}"
	case ty.IsCapsuleType():
		if ty == parser.TimeType {
			return "time.Time"
		}
	}
	return "UNKNOWN_TYPE"
}

func ctyTypeToGoCtyString(ty cty.Type) string {
	if ty.IsCapsuleType() {
		if ty == parser.TimeType {
			return "parser.TimeType"
		}
		return cty.DynamicPseudoType.GoString()
	}
	return ty.GoString()
}

func ctyTypeToTSString(ty cty.Type) string {
	switch {
	case ty == cty.Bool:
		return "boolean"
	case ty == cty.Number:
		return "number"
	case ty == cty.String:
		return "string"
	case ty.IsObjectType():
		return "object"
	case ty.IsMapType():
		return "object"
	case ty.IsCapsuleType():
		if ty == parser.TimeType {
			return "Date"
		}
	}
	return "UNKNOWN_TYPE"
}

func joinToType(ty string, single bool) string {
	result := camelToPascal(ty)
	if single {
		// Make it a pointer, because then we can check for nil when joins do not exist
		return "*" + result
	}
	// If not single, make it a slice
	return "[]" + result
}

func joinToTSType(ty string, single bool) string {
	if single {
		return ty
	}
	// If not single, make it a slice
	return ty + "[]"
}

func camelToPascal(value string) string {
	var result string
	parts := strings.Split(value, "_")
	for _, part := range parts {
		result += strings.Title(part)
	}
	return result

}
