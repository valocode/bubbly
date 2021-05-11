package builtin

import "github.com/valocode/bubbly/api/core"

type SchemaWrapper struct {
	Tables core.Tables `hcl:"table,block"`
}
