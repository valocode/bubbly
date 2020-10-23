package interim

import (
	"fmt"
	"reflect"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

var (
	reflectBool   = reflect.TypeOf(false)
	reflectInt    = reflect.TypeOf(0)
	reflectString = reflect.TypeOf("")
)

// newSchemaType creates a map of names->reflect.Type
// based on the given Tables. This function processes
// all subtables recursively
func newSchemaTypes(tables []Table) map[string]schemaType {
	types := make(map[string]schemaType)
	for _, t := range tables {
		addSchemaType(t, types)
	}
	return types
}

func addSchemaType(t Table, types map[string]schemaType) {
	// Initalize the list with the "id" and "fk" index fields
	// that are required for memdb.
	fields := []reflect.StructField{
		{
			Name: idFieldName,
			Type: reflectString,
		},
		{
			Name: fkFieldName,
			Type: reflectString,
		},
	}

	// Gather all fields from the table and map them to their
	// respective types from the reflect package.
	for _, f := range t.Fields {
		fields = append(fields, reflect.StructField{
			Name: f.Name,
			Type: reflectFieldType(f),
		})
	}

	// Set the type to the name of the table. Note that
	// this architecture assumes no two tables at any
	// point in the hierarchy will have the same names.
	// We should validate this or add internal namespaces
	// but keep in mind that namespacing will complicate querying.
	types[t.Name] = schemaType{
		rt: reflect.StructOf(fields),
	}

	// Recursively process all sub tables.
	for _, sub := range t.Tables {
		addSchemaType(sub, types)
	}
}

func reflectFieldType(f Field) reflect.Type {
	switch f.Type {
	case cty.Bool:
		return reflectBool
	case cty.Number:
		return reflectInt
	case cty.String:
		return reflectString
	}

	return nil
}

// schemaType wraps a reflect.Type with convenience
// methods for creating instances with values.
type schemaType struct {
	rt reflect.Type
}

// New creates a new instance of schemaType with the fields
// set based on the values in d.
func (t schemaType) New(d core.Data, id, fk string) (interface{}, error) {
	val := reflect.New(t.rt).Elem()

	// We now have a value of our dynamic struct.
	// We can set the fields based on our data now.
	for _, f := range d.Fields {
		fval := val.FieldByName(f.Name)
		switch f.Value.Type() {
		case cty.Bool:
			var n bool
			if err := gocty.FromCtyValue(f.Value, &n); err != nil {
				return nil, fmt.Errorf("falied to extract bool value for %s[%s]: %w", d.Name, f.Name, err)
			}
			fval.SetBool(n)
		case cty.Number:
			// TODO: Support different numeric types.
			var n int64
			if err := gocty.FromCtyValue(f.Value, &n); err != nil {
				return nil, fmt.Errorf("falied to extract numeric value for %s[%s]: %w", d.Name, f.Name, err)
			}
			fval.SetInt(n)
		case cty.String:
			var n string
			if err := gocty.FromCtyValue(f.Value, &n); err != nil {
				return nil, fmt.Errorf("falied to extract string value for %s[%s]: %w", d.Name, f.Name, err)
			}
			fval.SetString(n)
		}
	}

	val.FieldByName(idFieldName).SetString(id)
	val.FieldByName(fkFieldName).SetString(fk)

	return val.Interface(), nil
}
