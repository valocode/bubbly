package builtin

import "github.com/valocode/bubbly/api/core"

type SchemaWrapper struct {
	Tables []core.TableHCL `hcl:"table,block"`
}
