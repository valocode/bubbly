// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/license"
	"github.com/valocode/bubbly/ent/licenseusage"
	"github.com/valocode/bubbly/ent/predicate"
)

// LicenseUpdate is the builder for updating License entities.
type LicenseUpdate struct {
	config
	hooks    []Hook
	mutation *LicenseMutation
}

// Where adds a new predicate for the LicenseUpdate builder.
func (lu *LicenseUpdate) Where(ps ...predicate.License) *LicenseUpdate {
	lu.mutation.predicates = append(lu.mutation.predicates, ps...)
	return lu
}

// SetSpdxID sets the "spdx_id" field.
func (lu *LicenseUpdate) SetSpdxID(s string) *LicenseUpdate {
	lu.mutation.SetSpdxID(s)
	return lu
}

// SetName sets the "name" field.
func (lu *LicenseUpdate) SetName(s string) *LicenseUpdate {
	lu.mutation.SetName(s)
	return lu
}

// SetReference sets the "reference" field.
func (lu *LicenseUpdate) SetReference(s string) *LicenseUpdate {
	lu.mutation.SetReference(s)
	return lu
}

// SetNillableReference sets the "reference" field if the given value is not nil.
func (lu *LicenseUpdate) SetNillableReference(s *string) *LicenseUpdate {
	if s != nil {
		lu.SetReference(*s)
	}
	return lu
}

// ClearReference clears the value of the "reference" field.
func (lu *LicenseUpdate) ClearReference() *LicenseUpdate {
	lu.mutation.ClearReference()
	return lu
}

// SetDetailsURL sets the "details_url" field.
func (lu *LicenseUpdate) SetDetailsURL(s string) *LicenseUpdate {
	lu.mutation.SetDetailsURL(s)
	return lu
}

// SetNillableDetailsURL sets the "details_url" field if the given value is not nil.
func (lu *LicenseUpdate) SetNillableDetailsURL(s *string) *LicenseUpdate {
	if s != nil {
		lu.SetDetailsURL(*s)
	}
	return lu
}

// ClearDetailsURL clears the value of the "details_url" field.
func (lu *LicenseUpdate) ClearDetailsURL() *LicenseUpdate {
	lu.mutation.ClearDetailsURL()
	return lu
}

// SetIsOsiApproved sets the "is_osi_approved" field.
func (lu *LicenseUpdate) SetIsOsiApproved(b bool) *LicenseUpdate {
	lu.mutation.SetIsOsiApproved(b)
	return lu
}

// SetNillableIsOsiApproved sets the "is_osi_approved" field if the given value is not nil.
func (lu *LicenseUpdate) SetNillableIsOsiApproved(b *bool) *LicenseUpdate {
	if b != nil {
		lu.SetIsOsiApproved(*b)
	}
	return lu
}

// AddComponentIDs adds the "components" edge to the Component entity by IDs.
func (lu *LicenseUpdate) AddComponentIDs(ids ...int) *LicenseUpdate {
	lu.mutation.AddComponentIDs(ids...)
	return lu
}

// AddComponents adds the "components" edges to the Component entity.
func (lu *LicenseUpdate) AddComponents(c ...*Component) *LicenseUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return lu.AddComponentIDs(ids...)
}

// AddUsageIDs adds the "usages" edge to the LicenseUsage entity by IDs.
func (lu *LicenseUpdate) AddUsageIDs(ids ...int) *LicenseUpdate {
	lu.mutation.AddUsageIDs(ids...)
	return lu
}

// AddUsages adds the "usages" edges to the LicenseUsage entity.
func (lu *LicenseUpdate) AddUsages(l ...*LicenseUsage) *LicenseUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lu.AddUsageIDs(ids...)
}

// Mutation returns the LicenseMutation object of the builder.
func (lu *LicenseUpdate) Mutation() *LicenseMutation {
	return lu.mutation
}

// ClearComponents clears all "components" edges to the Component entity.
func (lu *LicenseUpdate) ClearComponents() *LicenseUpdate {
	lu.mutation.ClearComponents()
	return lu
}

// RemoveComponentIDs removes the "components" edge to Component entities by IDs.
func (lu *LicenseUpdate) RemoveComponentIDs(ids ...int) *LicenseUpdate {
	lu.mutation.RemoveComponentIDs(ids...)
	return lu
}

// RemoveComponents removes "components" edges to Component entities.
func (lu *LicenseUpdate) RemoveComponents(c ...*Component) *LicenseUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return lu.RemoveComponentIDs(ids...)
}

// ClearUsages clears all "usages" edges to the LicenseUsage entity.
func (lu *LicenseUpdate) ClearUsages() *LicenseUpdate {
	lu.mutation.ClearUsages()
	return lu
}

// RemoveUsageIDs removes the "usages" edge to LicenseUsage entities by IDs.
func (lu *LicenseUpdate) RemoveUsageIDs(ids ...int) *LicenseUpdate {
	lu.mutation.RemoveUsageIDs(ids...)
	return lu
}

// RemoveUsages removes "usages" edges to LicenseUsage entities.
func (lu *LicenseUpdate) RemoveUsages(l ...*LicenseUsage) *LicenseUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lu.RemoveUsageIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (lu *LicenseUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(lu.hooks) == 0 {
		if err = lu.check(); err != nil {
			return 0, err
		}
		affected, err = lu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*LicenseMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = lu.check(); err != nil {
				return 0, err
			}
			lu.mutation = mutation
			affected, err = lu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(lu.hooks) - 1; i >= 0; i-- {
			mut = lu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, lu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (lu *LicenseUpdate) SaveX(ctx context.Context) int {
	affected, err := lu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lu *LicenseUpdate) Exec(ctx context.Context) error {
	_, err := lu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lu *LicenseUpdate) ExecX(ctx context.Context) {
	if err := lu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lu *LicenseUpdate) check() error {
	if v, ok := lu.mutation.SpdxID(); ok {
		if err := license.SpdxIDValidator(v); err != nil {
			return &ValidationError{Name: "spdx_id", err: fmt.Errorf("ent: validator failed for field \"spdx_id\": %w", err)}
		}
	}
	if v, ok := lu.mutation.Name(); ok {
		if err := license.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	return nil
}

func (lu *LicenseUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   license.Table,
			Columns: license.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: license.FieldID,
			},
		},
	}
	if ps := lu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := lu.mutation.SpdxID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldSpdxID,
		})
	}
	if value, ok := lu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldName,
		})
	}
	if value, ok := lu.mutation.Reference(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldReference,
		})
	}
	if lu.mutation.ReferenceCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: license.FieldReference,
		})
	}
	if value, ok := lu.mutation.DetailsURL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldDetailsURL,
		})
	}
	if lu.mutation.DetailsURLCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: license.FieldDetailsURL,
		})
	}
	if value, ok := lu.mutation.IsOsiApproved(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: license.FieldIsOsiApproved,
		})
	}
	if lu.mutation.ComponentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   license.ComponentsTable,
			Columns: license.ComponentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: component.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.RemovedComponentsIDs(); len(nodes) > 0 && !lu.mutation.ComponentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   license.ComponentsTable,
			Columns: license.ComponentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: component.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.ComponentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   license.ComponentsTable,
			Columns: license.ComponentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: component.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lu.mutation.UsagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   license.UsagesTable,
			Columns: []string{license.UsagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: licenseusage.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.RemovedUsagesIDs(); len(nodes) > 0 && !lu.mutation.UsagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   license.UsagesTable,
			Columns: []string{license.UsagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: licenseusage.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.UsagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   license.UsagesTable,
			Columns: []string{license.UsagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: licenseusage.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, lu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{license.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// LicenseUpdateOne is the builder for updating a single License entity.
type LicenseUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *LicenseMutation
}

// SetSpdxID sets the "spdx_id" field.
func (luo *LicenseUpdateOne) SetSpdxID(s string) *LicenseUpdateOne {
	luo.mutation.SetSpdxID(s)
	return luo
}

// SetName sets the "name" field.
func (luo *LicenseUpdateOne) SetName(s string) *LicenseUpdateOne {
	luo.mutation.SetName(s)
	return luo
}

// SetReference sets the "reference" field.
func (luo *LicenseUpdateOne) SetReference(s string) *LicenseUpdateOne {
	luo.mutation.SetReference(s)
	return luo
}

// SetNillableReference sets the "reference" field if the given value is not nil.
func (luo *LicenseUpdateOne) SetNillableReference(s *string) *LicenseUpdateOne {
	if s != nil {
		luo.SetReference(*s)
	}
	return luo
}

// ClearReference clears the value of the "reference" field.
func (luo *LicenseUpdateOne) ClearReference() *LicenseUpdateOne {
	luo.mutation.ClearReference()
	return luo
}

// SetDetailsURL sets the "details_url" field.
func (luo *LicenseUpdateOne) SetDetailsURL(s string) *LicenseUpdateOne {
	luo.mutation.SetDetailsURL(s)
	return luo
}

// SetNillableDetailsURL sets the "details_url" field if the given value is not nil.
func (luo *LicenseUpdateOne) SetNillableDetailsURL(s *string) *LicenseUpdateOne {
	if s != nil {
		luo.SetDetailsURL(*s)
	}
	return luo
}

// ClearDetailsURL clears the value of the "details_url" field.
func (luo *LicenseUpdateOne) ClearDetailsURL() *LicenseUpdateOne {
	luo.mutation.ClearDetailsURL()
	return luo
}

// SetIsOsiApproved sets the "is_osi_approved" field.
func (luo *LicenseUpdateOne) SetIsOsiApproved(b bool) *LicenseUpdateOne {
	luo.mutation.SetIsOsiApproved(b)
	return luo
}

// SetNillableIsOsiApproved sets the "is_osi_approved" field if the given value is not nil.
func (luo *LicenseUpdateOne) SetNillableIsOsiApproved(b *bool) *LicenseUpdateOne {
	if b != nil {
		luo.SetIsOsiApproved(*b)
	}
	return luo
}

// AddComponentIDs adds the "components" edge to the Component entity by IDs.
func (luo *LicenseUpdateOne) AddComponentIDs(ids ...int) *LicenseUpdateOne {
	luo.mutation.AddComponentIDs(ids...)
	return luo
}

// AddComponents adds the "components" edges to the Component entity.
func (luo *LicenseUpdateOne) AddComponents(c ...*Component) *LicenseUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return luo.AddComponentIDs(ids...)
}

// AddUsageIDs adds the "usages" edge to the LicenseUsage entity by IDs.
func (luo *LicenseUpdateOne) AddUsageIDs(ids ...int) *LicenseUpdateOne {
	luo.mutation.AddUsageIDs(ids...)
	return luo
}

// AddUsages adds the "usages" edges to the LicenseUsage entity.
func (luo *LicenseUpdateOne) AddUsages(l ...*LicenseUsage) *LicenseUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return luo.AddUsageIDs(ids...)
}

// Mutation returns the LicenseMutation object of the builder.
func (luo *LicenseUpdateOne) Mutation() *LicenseMutation {
	return luo.mutation
}

// ClearComponents clears all "components" edges to the Component entity.
func (luo *LicenseUpdateOne) ClearComponents() *LicenseUpdateOne {
	luo.mutation.ClearComponents()
	return luo
}

// RemoveComponentIDs removes the "components" edge to Component entities by IDs.
func (luo *LicenseUpdateOne) RemoveComponentIDs(ids ...int) *LicenseUpdateOne {
	luo.mutation.RemoveComponentIDs(ids...)
	return luo
}

// RemoveComponents removes "components" edges to Component entities.
func (luo *LicenseUpdateOne) RemoveComponents(c ...*Component) *LicenseUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return luo.RemoveComponentIDs(ids...)
}

// ClearUsages clears all "usages" edges to the LicenseUsage entity.
func (luo *LicenseUpdateOne) ClearUsages() *LicenseUpdateOne {
	luo.mutation.ClearUsages()
	return luo
}

// RemoveUsageIDs removes the "usages" edge to LicenseUsage entities by IDs.
func (luo *LicenseUpdateOne) RemoveUsageIDs(ids ...int) *LicenseUpdateOne {
	luo.mutation.RemoveUsageIDs(ids...)
	return luo
}

// RemoveUsages removes "usages" edges to LicenseUsage entities.
func (luo *LicenseUpdateOne) RemoveUsages(l ...*LicenseUsage) *LicenseUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return luo.RemoveUsageIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (luo *LicenseUpdateOne) Select(field string, fields ...string) *LicenseUpdateOne {
	luo.fields = append([]string{field}, fields...)
	return luo
}

// Save executes the query and returns the updated License entity.
func (luo *LicenseUpdateOne) Save(ctx context.Context) (*License, error) {
	var (
		err  error
		node *License
	)
	if len(luo.hooks) == 0 {
		if err = luo.check(); err != nil {
			return nil, err
		}
		node, err = luo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*LicenseMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = luo.check(); err != nil {
				return nil, err
			}
			luo.mutation = mutation
			node, err = luo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(luo.hooks) - 1; i >= 0; i-- {
			mut = luo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, luo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (luo *LicenseUpdateOne) SaveX(ctx context.Context) *License {
	node, err := luo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (luo *LicenseUpdateOne) Exec(ctx context.Context) error {
	_, err := luo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (luo *LicenseUpdateOne) ExecX(ctx context.Context) {
	if err := luo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (luo *LicenseUpdateOne) check() error {
	if v, ok := luo.mutation.SpdxID(); ok {
		if err := license.SpdxIDValidator(v); err != nil {
			return &ValidationError{Name: "spdx_id", err: fmt.Errorf("ent: validator failed for field \"spdx_id\": %w", err)}
		}
	}
	if v, ok := luo.mutation.Name(); ok {
		if err := license.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	return nil
}

func (luo *LicenseUpdateOne) sqlSave(ctx context.Context) (_node *License, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   license.Table,
			Columns: license.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: license.FieldID,
			},
		},
	}
	id, ok := luo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing License.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := luo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, license.FieldID)
		for _, f := range fields {
			if !license.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != license.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := luo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := luo.mutation.SpdxID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldSpdxID,
		})
	}
	if value, ok := luo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldName,
		})
	}
	if value, ok := luo.mutation.Reference(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldReference,
		})
	}
	if luo.mutation.ReferenceCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: license.FieldReference,
		})
	}
	if value, ok := luo.mutation.DetailsURL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldDetailsURL,
		})
	}
	if luo.mutation.DetailsURLCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: license.FieldDetailsURL,
		})
	}
	if value, ok := luo.mutation.IsOsiApproved(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: license.FieldIsOsiApproved,
		})
	}
	if luo.mutation.ComponentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   license.ComponentsTable,
			Columns: license.ComponentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: component.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.RemovedComponentsIDs(); len(nodes) > 0 && !luo.mutation.ComponentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   license.ComponentsTable,
			Columns: license.ComponentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: component.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.ComponentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   license.ComponentsTable,
			Columns: license.ComponentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: component.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if luo.mutation.UsagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   license.UsagesTable,
			Columns: []string{license.UsagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: licenseusage.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.RemovedUsagesIDs(); len(nodes) > 0 && !luo.mutation.UsagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   license.UsagesTable,
			Columns: []string{license.UsagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: licenseusage.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.UsagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   license.UsagesTable,
			Columns: []string{license.UsagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: licenseusage.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &License{config: luo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, luo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{license.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
