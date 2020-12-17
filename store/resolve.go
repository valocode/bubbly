package store

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

type resolver struct {
	p provider
}

func (r *resolver) Resolve(data core.DataBlocks) (core.DataBlocks, error) {
	m := make(map[string]cty.Value)
	r.flatten(m, data)
	return r.resolveRefs(m, data)
}

func (r *resolver) resolveRefs(m map[string]cty.Value, data core.DataBlocks) (core.DataBlocks, error) {
	altData := make(core.DataBlocks, 0, len(data))

	for _, d := range data {
		for _, ref := range d.DataRefs {
			key := r.keyName(ref.TableName, ref.Field)

			val, ok := m[key]
			if !ok {
				var err error
				val, err = r.p.LastValue(ref.TableName, ref.Field)
				if err != nil {
					return nil, fmt.Errorf("failed to get last value from provider: %w", err)
				}
			}
			d.Fields = append(d.Fields, core.DataField{
				Name:  r.fieldName(ref.TableName, ref.Field),
				Value: val,
			})
		}
		d.DataRefs = nil
		altData = append(altData, d)
	}

	return altData, nil
}

func (r *resolver) flatten(m map[string]cty.Value, data core.DataBlocks) {
	for _, d := range data {
		for _, f := range d.Fields {
			m[r.keyName(d.TableName, f.Name)] = f.Value
		}
		r.flatten(m, d.Data)
	}
}

func (r *resolver) keyName(table, field string) string {
	return table + "." + field
}

func (r *resolver) fieldName(table, field string) string {
	return table + "_" + field
}
