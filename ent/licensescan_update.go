// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/licensescan"
	"github.com/valocode/bubbly/ent/licenseusage"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releaseentry"
)

// LicenseScanUpdate is the builder for updating LicenseScan entities.
type LicenseScanUpdate struct {
	config
	hooks    []Hook
	mutation *LicenseScanMutation
}

// Where adds a new predicate for the LicenseScanUpdate builder.
func (lsu *LicenseScanUpdate) Where(ps ...predicate.LicenseScan) *LicenseScanUpdate {
	lsu.mutation.predicates = append(lsu.mutation.predicates, ps...)
	return lsu
}

// SetReleaseID sets the "release" edge to the Release entity by ID.
func (lsu *LicenseScanUpdate) SetReleaseID(id int) *LicenseScanUpdate {
	lsu.mutation.SetReleaseID(id)
	return lsu
}

// SetRelease sets the "release" edge to the Release entity.
func (lsu *LicenseScanUpdate) SetRelease(r *Release) *LicenseScanUpdate {
	return lsu.SetReleaseID(r.ID)
}

// SetEntryID sets the "entry" edge to the ReleaseEntry entity by ID.
func (lsu *LicenseScanUpdate) SetEntryID(id int) *LicenseScanUpdate {
	lsu.mutation.SetEntryID(id)
	return lsu
}

// SetNillableEntryID sets the "entry" edge to the ReleaseEntry entity by ID if the given value is not nil.
func (lsu *LicenseScanUpdate) SetNillableEntryID(id *int) *LicenseScanUpdate {
	if id != nil {
		lsu = lsu.SetEntryID(*id)
	}
	return lsu
}

// SetEntry sets the "entry" edge to the ReleaseEntry entity.
func (lsu *LicenseScanUpdate) SetEntry(r *ReleaseEntry) *LicenseScanUpdate {
	return lsu.SetEntryID(r.ID)
}

// AddLicenseIDs adds the "licenses" edge to the LicenseUsage entity by IDs.
func (lsu *LicenseScanUpdate) AddLicenseIDs(ids ...int) *LicenseScanUpdate {
	lsu.mutation.AddLicenseIDs(ids...)
	return lsu
}

// AddLicenses adds the "licenses" edges to the LicenseUsage entity.
func (lsu *LicenseScanUpdate) AddLicenses(l ...*LicenseUsage) *LicenseScanUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lsu.AddLicenseIDs(ids...)
}

// Mutation returns the LicenseScanMutation object of the builder.
func (lsu *LicenseScanUpdate) Mutation() *LicenseScanMutation {
	return lsu.mutation
}

// ClearRelease clears the "release" edge to the Release entity.
func (lsu *LicenseScanUpdate) ClearRelease() *LicenseScanUpdate {
	lsu.mutation.ClearRelease()
	return lsu
}

// ClearEntry clears the "entry" edge to the ReleaseEntry entity.
func (lsu *LicenseScanUpdate) ClearEntry() *LicenseScanUpdate {
	lsu.mutation.ClearEntry()
	return lsu
}

// ClearLicenses clears all "licenses" edges to the LicenseUsage entity.
func (lsu *LicenseScanUpdate) ClearLicenses() *LicenseScanUpdate {
	lsu.mutation.ClearLicenses()
	return lsu
}

// RemoveLicenseIDs removes the "licenses" edge to LicenseUsage entities by IDs.
func (lsu *LicenseScanUpdate) RemoveLicenseIDs(ids ...int) *LicenseScanUpdate {
	lsu.mutation.RemoveLicenseIDs(ids...)
	return lsu
}

// RemoveLicenses removes "licenses" edges to LicenseUsage entities.
func (lsu *LicenseScanUpdate) RemoveLicenses(l ...*LicenseUsage) *LicenseScanUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lsu.RemoveLicenseIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (lsu *LicenseScanUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(lsu.hooks) == 0 {
		if err = lsu.check(); err != nil {
			return 0, err
		}
		affected, err = lsu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*LicenseScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = lsu.check(); err != nil {
				return 0, err
			}
			lsu.mutation = mutation
			affected, err = lsu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(lsu.hooks) - 1; i >= 0; i-- {
			mut = lsu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, lsu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (lsu *LicenseScanUpdate) SaveX(ctx context.Context) int {
	affected, err := lsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lsu *LicenseScanUpdate) Exec(ctx context.Context) error {
	_, err := lsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lsu *LicenseScanUpdate) ExecX(ctx context.Context) {
	if err := lsu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lsu *LicenseScanUpdate) check() error {
	if _, ok := lsu.mutation.ReleaseID(); lsu.mutation.ReleaseCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"release\"")
	}
	return nil
}

func (lsu *LicenseScanUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   licensescan.Table,
			Columns: licensescan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: licensescan.FieldID,
			},
		},
	}
	if ps := lsu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if lsu.mutation.ReleaseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   licensescan.ReleaseTable,
			Columns: []string{licensescan.ReleaseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: release.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lsu.mutation.ReleaseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   licensescan.ReleaseTable,
			Columns: []string{licensescan.ReleaseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: release.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lsu.mutation.EntryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   licensescan.EntryTable,
			Columns: []string{licensescan.EntryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releaseentry.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lsu.mutation.EntryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   licensescan.EntryTable,
			Columns: []string{licensescan.EntryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releaseentry.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lsu.mutation.LicensesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   licensescan.LicensesTable,
			Columns: []string{licensescan.LicensesColumn},
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
	if nodes := lsu.mutation.RemovedLicensesIDs(); len(nodes) > 0 && !lsu.mutation.LicensesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   licensescan.LicensesTable,
			Columns: []string{licensescan.LicensesColumn},
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
	if nodes := lsu.mutation.LicensesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   licensescan.LicensesTable,
			Columns: []string{licensescan.LicensesColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, lsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{licensescan.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// LicenseScanUpdateOne is the builder for updating a single LicenseScan entity.
type LicenseScanUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *LicenseScanMutation
}

// SetReleaseID sets the "release" edge to the Release entity by ID.
func (lsuo *LicenseScanUpdateOne) SetReleaseID(id int) *LicenseScanUpdateOne {
	lsuo.mutation.SetReleaseID(id)
	return lsuo
}

// SetRelease sets the "release" edge to the Release entity.
func (lsuo *LicenseScanUpdateOne) SetRelease(r *Release) *LicenseScanUpdateOne {
	return lsuo.SetReleaseID(r.ID)
}

// SetEntryID sets the "entry" edge to the ReleaseEntry entity by ID.
func (lsuo *LicenseScanUpdateOne) SetEntryID(id int) *LicenseScanUpdateOne {
	lsuo.mutation.SetEntryID(id)
	return lsuo
}

// SetNillableEntryID sets the "entry" edge to the ReleaseEntry entity by ID if the given value is not nil.
func (lsuo *LicenseScanUpdateOne) SetNillableEntryID(id *int) *LicenseScanUpdateOne {
	if id != nil {
		lsuo = lsuo.SetEntryID(*id)
	}
	return lsuo
}

// SetEntry sets the "entry" edge to the ReleaseEntry entity.
func (lsuo *LicenseScanUpdateOne) SetEntry(r *ReleaseEntry) *LicenseScanUpdateOne {
	return lsuo.SetEntryID(r.ID)
}

// AddLicenseIDs adds the "licenses" edge to the LicenseUsage entity by IDs.
func (lsuo *LicenseScanUpdateOne) AddLicenseIDs(ids ...int) *LicenseScanUpdateOne {
	lsuo.mutation.AddLicenseIDs(ids...)
	return lsuo
}

// AddLicenses adds the "licenses" edges to the LicenseUsage entity.
func (lsuo *LicenseScanUpdateOne) AddLicenses(l ...*LicenseUsage) *LicenseScanUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lsuo.AddLicenseIDs(ids...)
}

// Mutation returns the LicenseScanMutation object of the builder.
func (lsuo *LicenseScanUpdateOne) Mutation() *LicenseScanMutation {
	return lsuo.mutation
}

// ClearRelease clears the "release" edge to the Release entity.
func (lsuo *LicenseScanUpdateOne) ClearRelease() *LicenseScanUpdateOne {
	lsuo.mutation.ClearRelease()
	return lsuo
}

// ClearEntry clears the "entry" edge to the ReleaseEntry entity.
func (lsuo *LicenseScanUpdateOne) ClearEntry() *LicenseScanUpdateOne {
	lsuo.mutation.ClearEntry()
	return lsuo
}

// ClearLicenses clears all "licenses" edges to the LicenseUsage entity.
func (lsuo *LicenseScanUpdateOne) ClearLicenses() *LicenseScanUpdateOne {
	lsuo.mutation.ClearLicenses()
	return lsuo
}

// RemoveLicenseIDs removes the "licenses" edge to LicenseUsage entities by IDs.
func (lsuo *LicenseScanUpdateOne) RemoveLicenseIDs(ids ...int) *LicenseScanUpdateOne {
	lsuo.mutation.RemoveLicenseIDs(ids...)
	return lsuo
}

// RemoveLicenses removes "licenses" edges to LicenseUsage entities.
func (lsuo *LicenseScanUpdateOne) RemoveLicenses(l ...*LicenseUsage) *LicenseScanUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lsuo.RemoveLicenseIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (lsuo *LicenseScanUpdateOne) Select(field string, fields ...string) *LicenseScanUpdateOne {
	lsuo.fields = append([]string{field}, fields...)
	return lsuo
}

// Save executes the query and returns the updated LicenseScan entity.
func (lsuo *LicenseScanUpdateOne) Save(ctx context.Context) (*LicenseScan, error) {
	var (
		err  error
		node *LicenseScan
	)
	if len(lsuo.hooks) == 0 {
		if err = lsuo.check(); err != nil {
			return nil, err
		}
		node, err = lsuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*LicenseScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = lsuo.check(); err != nil {
				return nil, err
			}
			lsuo.mutation = mutation
			node, err = lsuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(lsuo.hooks) - 1; i >= 0; i-- {
			mut = lsuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, lsuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (lsuo *LicenseScanUpdateOne) SaveX(ctx context.Context) *LicenseScan {
	node, err := lsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (lsuo *LicenseScanUpdateOne) Exec(ctx context.Context) error {
	_, err := lsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lsuo *LicenseScanUpdateOne) ExecX(ctx context.Context) {
	if err := lsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lsuo *LicenseScanUpdateOne) check() error {
	if _, ok := lsuo.mutation.ReleaseID(); lsuo.mutation.ReleaseCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"release\"")
	}
	return nil
}

func (lsuo *LicenseScanUpdateOne) sqlSave(ctx context.Context) (_node *LicenseScan, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   licensescan.Table,
			Columns: licensescan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: licensescan.FieldID,
			},
		},
	}
	id, ok := lsuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing LicenseScan.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := lsuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, licensescan.FieldID)
		for _, f := range fields {
			if !licensescan.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != licensescan.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := lsuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if lsuo.mutation.ReleaseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   licensescan.ReleaseTable,
			Columns: []string{licensescan.ReleaseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: release.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lsuo.mutation.ReleaseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   licensescan.ReleaseTable,
			Columns: []string{licensescan.ReleaseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: release.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lsuo.mutation.EntryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   licensescan.EntryTable,
			Columns: []string{licensescan.EntryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releaseentry.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lsuo.mutation.EntryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   licensescan.EntryTable,
			Columns: []string{licensescan.EntryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releaseentry.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lsuo.mutation.LicensesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   licensescan.LicensesTable,
			Columns: []string{licensescan.LicensesColumn},
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
	if nodes := lsuo.mutation.RemovedLicensesIDs(); len(nodes) > 0 && !lsuo.mutation.LicensesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   licensescan.LicensesTable,
			Columns: []string{licensescan.LicensesColumn},
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
	if nodes := lsuo.mutation.LicensesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   licensescan.LicensesTable,
			Columns: []string{licensescan.LicensesColumn},
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
	_node = &LicenseScan{config: lsuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, lsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{licensescan.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
