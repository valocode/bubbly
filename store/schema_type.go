package store

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

const (
	idFieldName         = "ID"
	schemaTypeFieldName = "SchemaType"
	titledIDFieldName   = "Id" // if "id" is titled it becomes "Id"
)

var (
	reflectBool   = reflect.TypeOf(false)
	reflectInt    = reflect.TypeOf(0)
	reflectString = reflect.TypeOf("")
)

// newSchemaType creates a map of names->reflect.Type
// based on the given Tables. This function processes
// all subtables recursively
func newSchemaTypes(tables core.Tables) map[string]schemaType {
	types := make(map[string]schemaType)
	for _, t := range tables {
		addSchemaType(t, types)
	}
	return types
}

func addSchemaType(t core.Table, types map[string]schemaType) {
	// Gather all fields from the table and map them to their
	// respective types from the reflect package.
	fields := make([]reflect.StructField, 0, len(t.Fields))
	for _, f := range t.Fields {
		fields = append(fields, reflect.StructField{
			// Make sure that the attribute is exported on the struct.
			Name: strings.Title(f.Name),
			Type: reflectFieldType(f),
		})
	}

	// Store the name of the type on the type itself.
	fields = append(fields, reflect.StructField{
		Name: schemaTypeFieldName,
		Type: reflectString,
	})

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

func reflectFieldType(f core.TableField) reflect.Type {
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

// Empty creates a zero'd instanced of the schemaType.
func (t schemaType) Empty() interface{} {
	return reflect.New(t.rt).Interface()
}

// EmptySlice creates a slice of zero'd instanced of the schemaType.
func (t schemaType) EmptySlice() interface{} {
	return reflect.New(reflect.SliceOf(reflect.PtrTo(t.rt))).Interface()
}

// New creates a new instance of schemaType with the fields
// set based on the values in d.
func (t schemaType) New(d core.Data, tableRefs tableRefs) (interface{}, error) {
	var (
		val  = reflect.New(t.rt)
		elem = val.Elem()
	)

	// We now have a value of our dynamic struct.
	// We can set the fields based on our data now.
	for _, f := range d.Fields {
		// TODO(andrewhare): Revist this to be more flexible, again
		// we have to title-case the attr to match what we did when
		// we created the type itself.
		var (
			name = strings.Title(f.Name)
			fval = elem.FieldByName(name)
		)
		switch f.Value.Type() {
		case cty.Bool:
			var n bool
			if err := gocty.FromCtyValue(f.Value, &n); err != nil {
				return nil, fmt.Errorf("falied to extract bool value for %s.%s: %w", d.TableName, name, err)
			}
			fval.SetBool(n)
		case cty.Number:
			// TODO: Support different numeric types.
			var n int64
			if err := gocty.FromCtyValue(f.Value, &n); err != nil {
				return nil, fmt.Errorf("falied to extract numeric value for %s.%s: %w", d.TableName, name, err)
			}
			fval.SetInt(n)
		case cty.String:
			var n string
			if err := gocty.FromCtyValue(f.Value, &n); err != nil {
				return nil, fmt.Errorf("falied to extract string value for %s.%s: %w", d.TableName, name, err)
			}
			fval.SetString(n)
		case parser.DataRefType:
			ref := f.Value.EncapsulatedValue().(*parser.DataRef)
			val, ok := tableRefs[ref.TableName]
			if !ok {
				return nil, fmt.Errorf("could not find referenced data table: %s", ref.TableName)
			}
			// get the name of the field according to the schema type struct
			fieldName := strings.Title(ref.Field)
			// if the "id" field is referenced, the titled form is "Id" not "ID"
			if fieldName == titledIDFieldName {
				fieldName = idFieldName
			}
			rval := reflect.ValueOf(val.Value).Elem().FieldByName(fieldName)
			if !rval.IsValid() {
				return nil, fmt.Errorf(`field "%s" for type "%s" does not exist`, fieldName, reflect.TypeOf(val.Value).String())
			}
			fval.Set(rval)
		default:
			return nil, fmt.Errorf("unsupported cty type: %s", f.Value.Type().FriendlyName())
		}
	}

	elem.FieldByName(schemaTypeFieldName).SetString(d.TableName)

	if d.ParentTable != "" {
		parent, ok := tableRefs[d.ParentTable]
		if !ok {
			return nil, fmt.Errorf("could not find parent table: %s", d.ParentTable)
		}
		elem.FieldByName(strings.Title(d.ParentTable + "_id")).SetInt(parent.ID)
	}

	return val.Interface(), nil
}

// schemaTypeID returns the ID for the given value.
func schemaTypeID(n interface{}) int64 {
	return reflect.ValueOf(n).Elem().FieldByName(idFieldName).Int()
}

// schemaTypeName returns the name of the schema type for the given value.
func schemaTypeName(n interface{}) string {
	return reflect.ValueOf(n).Elem().FieldByName(schemaTypeFieldName).String()
}

func schemaTypeVal(n interface{}, field string) (cty.Value, error) {
	f := reflect.ValueOf(n).Elem().FieldByName(field).Interface()
	ty, err := gocty.ImpliedType(f)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to get implied cty type: %w", err)
	}
	val, err := gocty.ToCtyValue(f, ty)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to create cty value: %w", err)
	}
	return val, nil
}
