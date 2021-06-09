package parser

import (
	"fmt"
	"reflect"
	"time"

	"github.com/zclconf/go-cty/cty"
)

// NOTE: these types should probably not be here, but time has it that they
// will live here until we find a better place

var DataRefType = cty.CapsuleWithOps(
	"DataRef", reflect.TypeOf(DataRef{}),
	&cty.CapsuleOps{
		GoString: func(val interface{}) string { return fmt.Sprintf("%#v", val) },
		RawEquals: func(a, b interface{}) bool {
			d1 := a.(*DataRef)
			d2 := a.(*DataRef)
			return d1.TableName == d2.TableName && d1.Field == d2.Field
		},
	},
)

// DataRef is a data block that does not contain a
// static value but references a value from another
// data block.
type DataRef struct {
	TableName string `json:"table"`
	Field     string `json:"field"`
}

var TimeType = cty.CapsuleWithOps(
	"Time", reflect.TypeOf(time.Time{}),
	&cty.CapsuleOps{
		GoString: func(val interface{}) string { return fmt.Sprintf("%#v", val) },
		RawEquals: func(a, b interface{}) bool {
			t1 := a.(*time.Time)
			t2 := b.(*time.Time)
			return t1.Equal(*t2)
		},
	},
)
