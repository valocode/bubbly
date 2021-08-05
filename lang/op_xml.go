package lang

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// decodeXML reads in, decodes, and validates the format of data
func decodeXML(r io.Reader, ty cty.Type) (cty.Value, error) {

	data, err := mxj.NewMapXmlReader(r, true)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to decode XML: %w", err)
	}

	if err := fixListsInXML(&data, ty); err != nil {
		return cty.NilVal, err
	}

	val, err := gocty.ToCtyValue(data, ty)
	if err != nil {
		return cty.NilVal, err
	}

	return val, nil
}

// TODO: fixListsInXML could do with extensive unit testing of edge cases and better documentation

// fixListsInXML updates those elements in XML tree who should have been
// the only member of a list of length one. In the absence of a schema or
// a document type definition, the XML parser cannot tell which elements
// must be placed into a list of length one. This is because XML does not
// have syntax for lists, unlike JSON. But the format of data is known to _us_
// via the resource definition. As there is no easy way to communicate this
// information to the XML parser, the second best approach is to add a post-
// processing step. It traverses the data format definition recursively, identifying
// the lists in it, and validates selected branches of the XML tree, updating
// those elements which should have been places in a list.
func fixListsInXML(data *mxj.Map, ty cty.Type) error {

	// Forward declaration: recursive inner function over the format spec
	var f func(*mxj.Map, cty.Type, []string, int) error

	// Implementation: recursive inner function over the format spec
	f = func(data *mxj.Map, ty cty.Type, path []string, idx int) error {

		// The full path to the current node
		// in a format that XML manipulation functions understand.
		pathStr := strings.Join(path, ".")

		if idx > 0 {
			pathStr += fmt.Sprint("[", idx, "]")
		}

		// Traverse all fields of the objects recursively,
		// as they may contain lists as well.
		if ty.IsObjectType() {
			for x := range ty.AttributeTypes() {
				path = append(path, x)
				pathIdx := len(path) - 1

				f(data, ty.AttributeType(x), path, 0)
				path = path[0:pathIdx]
			}
		}

		// For any list declarated in the format spec,
		// the XML tree has to be checked and updated,
		// if necessary...
		if ty.IsListType() {

			// Get the XML tree elements at that level
			vs, err := data.ValuesForPath(pathStr)
			if err != nil {
				return fmt.Errorf("wrong path (%s) in XML syntax tree: %w", pathStr, err)
			}

			n := len(vs)

			// If there is only one conforming element present in the XML tree, it means
			// the parser would not know that it had to make a list for it. To fix that,
			// create a list of length one for that element and update the XML tree.
			switch n {
			case 0:
				return fmt.Errorf("xml data structure inconsistent state, ValuesForPath are zero at %s", pathStr)
			case 1:
				v := vs[0]

				if reflect.TypeOf(v).Kind() == reflect.Map {
					vv := make([]interface{}, 0)
					vv = append(vv, v)
					if err := data.SetValueForPath(vv, pathStr); err != nil {
						return fmt.Errorf("cannot convert at path %s, error %w", pathStr, err)
					}
				}
				fallthrough
			default:
				for i := range vs {
					return f(data, ty.ElementType(), path, i)
				}
			}
		}

		return nil
	}

	path := make([]string, 0)
	return f(data, ty, path, 0)
}
