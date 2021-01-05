package store

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/parser"
)

// tableRefs stores instances of a tableRef using the tableRef.Name as the key.
// This enables retrieval of values such as the ID for a given table, such that
// data refs can be resolved.
type tableRefs map[string]tableRef

// tableRef stores the ID of a record, the table name, and the value of the
// instance of the struct that go-pg uses when inserting data to postgres
type tableRef struct {
	ID    int64
	Name  string
	Value interface{}
}

// Takes the data blocks as they come to the store (with mixed cty.Values and
// smuggled DataRef CapsuleVals). This sets all the normal blocks of data, that
// do not contain DataRefs, to altData and all the blocks of data containing
// DataRefs to ref
func prepareDataRefs(data core.DataBlocks, altData *core.DataBlocks, refs *core.DataBlocks) error {
	dataMap := make(map[string]core.Data)
	if err := flatten(data, altData, refs, dataMap, ""); err != nil {
		return err
	}
	if err := validateDataRefs(dataMap, *refs); err != nil {
		return err
	}

	return nil
}

func flatten(data core.DataBlocks, altData *core.DataBlocks, refs *core.DataBlocks, dataMap map[string]core.Data, parentTable string) error {
	for _, d := range data {

		hasDataRef := false
		// Check if this data block has any data refs in its fields
		for _, f := range d.Fields {
			// We are interested in DataRefs which are capsule types, so ignore
			// other types
			if f.Value.Type() != parser.DataRefType {
				continue
			}
			hasDataRef = true
		}

		// Assign the ParentTable before we make a flat list of the DataBlocks
		// after which the hierarchy is lost
		d.ParentTable = parentTable

		// If the data block has a data ref, then assign it to the dataRef slice
		// Else add it to the normal data blocks.
		if hasDataRef {
			*refs = append(*refs, d)
		} else {
			*altData = append(*altData, d)
		}

		// Add each table to dataMap so that we can validate that the data refs
		// refer to a table in this bunch of DataBlocks
		dataMap[d.TableName] = d

		if err := flatten(d.Data, altData, refs, dataMap, d.TableName); err != nil {
			return err
		}
	}

	return nil
}

func validateDataRefs(dataMap map[string]core.Data, refs core.DataBlocks) error {
	for _, d := range refs {
		for _, f := range d.Fields {
			// We are interested in DataRefs which are capsule types, so ignore
			// other types
			if f.Value.Type() != parser.DataRefType {
				continue
			}
			// We have already evaluated that any capsule types are DataRefs
			ref := f.Value.EncapsulatedValue().(*parser.DataRef)
			// Validate that the table which the data ref references exists in
			// the original set of DataBlocks we received. Otherwise we should
			// error
			_, ok := dataMap[ref.TableName]
			if !ok {
				return fmt.Errorf("reference to unspecified data block: %s.%s", ref.TableName, ref.Field)
			}
		}
	}
	return nil
}

func keyName(table, field string) string {
	return table + "." + field
}

func fieldName(table, field string) string {
	return table + "_" + field
}
