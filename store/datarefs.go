package store

import (
	"sort"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/parser"
)

// dataRefs is a map of maps to store resolved dataref values.
// The first map is the table name, and the second map is for the field name.
// The stored value is the one that is created by storing the data block that
// corresponds to that dataref.
type dataRefs map[string]map[string]interface{}

// RefFields takes a table name and returns the fields that are referenced
// (i.e. there exists a dataref for those fields) so that the provider can know
// which fields for a data ref table need to be resolved and values stored.
// Basically: which fields are referenced in a particular table.
func (d dataRefs) RefFields(tableName string) []string {
	var (
		fields     = d[tableName]
		fieldNames = make([]string, len(fields))
		i          = 0
	)
	// At the minimum return the id field to return
	if len(fieldNames) == 0 {
		return []string{tableIDField}
	}
	for k := range fields {
		fieldNames[i] = k
		i++
	}
	sort.Strings(fieldNames)
	return fieldNames
}

// flatten recursively goes through the data blocks and flattens them
func flatten(data core.DataBlocks, parentTable string) core.DataBlocks {
	altData := core.DataBlocks{}
	for _, d := range data {
		// Assign the ParentTable before we make a flat list of the DataBlocks
		// after which the hierarchy is lost
		d.ParentTable = parentTable

		// Add nested data blocks to data block
		altData = append(altData, flatten(d.Data, d.TableName)...)
		// Clear the nested data
		d.Data = nil
		altData = append(altData, d)
	}

	return data
}

// orderDataRefs takes the flattened list of data blocks and puts all the
// datarefs at the end of the list, so that they are resolved last.
// At this point there is no cleverness with the ordering of datarefs, so
// datarefs that depend on other datarefs may have an issue.
func orderDataRefs(sc *saveContext, data core.DataBlocks) {
	var dataRefs = make(core.DataBlocks, 0)
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
			// Get the underlying DataRef value
			ref := f.Value.EncapsulatedValue().(*parser.DataRef)
			dataRef, ok := sc.DataRefs[ref.TableName]
			if !ok {
				// If it does not yet exist, then initialize it with at least
				// length 1 (we are about to add something to it, right?).
				dataRef = make(map[string]interface{}, 1)
				sc.DataRefs[ref.TableName] = dataRef
			}
			// Set the value to nil, as this should get assigned properly when
			// we start resolving the data refs (otherwise nil indicates that
			// data refs have not been resolved properly).
			dataRef[ref.Field] = nil
		}

		// If the data block has a data ref, then assign it to the dataRef slice
		// Else add it to the normal data blocks.
		if hasDataRef {
			dataRefs = append(dataRefs, d)
		} else {
			sc.Data = append(sc.Data, d)
		}
	}
	// Append the DataRefs to the list of data blocks so that these get resolved
	// last
	sc.Data = append(sc.Data, dataRefs...)
	return
}
