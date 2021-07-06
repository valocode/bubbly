package builtin

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

const (
	DBlockTableName  = "DBlock_Table"
	DBlockPolicyName = "DBlock_Policy"
	DBlockJoins      = "DBlock_Joins"
)

// ToDataBlocks takes an instance (or slice of instances) of one of the generated
// Go structs from the Bubbly Schema.
// The idea is that instead of creating core.DataBlocks directly, you can create
// an instance of a struct (typed) and generate the DataBlocks using this method.
// It is a fairly involved process, using lots of reflection, but it works nicely
func ToDataBlocks(value interface{}) core.DataBlocks {
	return toDataBlocks(value, false)
}
func toDataBlocks(value interface{}, ignoreNesting bool) core.DataBlocks {
	var dbs core.DataBlocks
	var blocks []interface{}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice:
		blocks = value.([]interface{})
	case reflect.Struct:
		blocks = []interface{}{value}
	default:
		panic("did not receive an []interface{} but " + reflect.TypeOf(value).String())
	}
	for _, block := range blocks {
		var (
			data core.Data
			val  = reflect.ValueOf(block)
			ty   = reflect.TypeOf(block)
		)
		for i := 0; i < ty.NumField(); i++ {
			var (
				fieldVal    = val.Field(i)
				structField = ty.Field(i)
			)
			// Handle the special cases
			switch structField.Name {
			case DBlockTableName:
				data.TableName = jsonTagName(structField)
				continue
			case DBlockPolicyName:
				data.Policy = core.DataBlockPolicy(fieldVal.String())
				continue
			case DBlockJoins:
				if structField.Type.Kind() != reflect.Slice && structField.Type.Elem().Kind() != reflect.String {
					panic(DBlockJoins + " type must be of type []string ... something strange has happened with the generation...")
				}
				for j := 0; j < fieldVal.Len(); j++ {
					data.Joins = append(data.Joins, fieldVal.Index(j).String())
				}
				continue
			}
			// If the struct field was not one of the "special fields", it is
			// an actual value or a relationship to other structs
			switch structField.Type.Kind() {
			case reflect.Struct:
				// Some structs we want to store as capsule types, like time.Time
				if structField.Type == reflect.TypeOf(time.Time{}) {
					dataFieldName := jsonTagName(structField)
					// We cannot Addr the struct field value, so create a new value
					// and copy the time.Time over into the new ptr
					timePtr := reflect.New(fieldVal.Type())
					timePtr.Elem().Set(fieldVal)
					dataFieldValue := cty.CapsuleVal(parser.TimeType, timePtr.Interface())
					data.Fields.Values[dataFieldName] = dataFieldValue
					continue
				}
				data.Data = append(data.Data, toDataBlocks(fieldVal.Interface(), false)...)
			case reflect.Slice:
				// If a slice, we want to iterate through the slices and recursively
				// generate the data blocks
				eType := structField.Type.Elem()
				if eType.Kind() == reflect.Struct {
					for j := 0; j < fieldVal.Len(); j++ {
						data.Data = append(data.Data, toDataBlocks(fieldVal.Index(j).Interface(), false)...)
					}
				} else {
					panic("slice of non-struct type not supported...")
				}
			case reflect.Ptr:
				// We don't care about Nil pointers
				if fieldVal.IsNil() {
					continue
				}
				// If a valid ptr, recursively generated data blocks for the pointer
				// BUT set the ignoreNesting flag to true because the data block
				// may not belong to the parent... It's all about the relationships.
				// If A --> B, and B has a pointer to C, then it doesnt mean
				// C belongs to A. By setting ignore nesting to true then we tell
				// Bubbly to ignore the fact that C gets nested under A
				dbs = append(dbs,
					toDataBlocks(fieldVal.Elem().Interface(), true)...,
				)
				data.Joins = append(data.Joins, jsonTagName(structField))
			default:
				// We also don't care about zero value (e.g. empty string, 0 int)
				if fieldVal.IsZero() {
					continue
				}
				// If no special cases then we have a data block field.
				// Make sure data fields are initialised
				if data.Fields == nil {
					data.Fields = &core.DataFields{
						Values: make(map[string]cty.Value),
					}
				}
				var (
					dataFieldName  = jsonTagName(structField)
					dataFieldValue cty.Value
				)
				ctyType, err := gocty.ImpliedType(fieldVal.Interface())
				if err != nil {
					panic(fmt.Sprintf("could not imply cty.Type from field %s: %s", structField.Name, err.Error()))
				}
				dataFieldValue, err = gocty.ToCtyValue(fieldVal.Interface(), ctyType)
				if err != nil {
					panic(fmt.Sprintf("could not convert field to cty.Value: %s: %s", structField.Name, err.Error()))
				}
				data.Fields.Values[dataFieldName] = dataFieldValue
			}
		}
		// Set the flat to ignore nesting
		data.IgnoreNesting = ignoreNesting
		dbs = append(dbs, data)
	}

	return dbs
}

func jsonTagName(field reflect.StructField) string {
	return strings.Split(field.Tag.Get("json"), ",")[0]
}
