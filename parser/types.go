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
	},
)
