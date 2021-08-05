// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/codescan"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releaseentry"
)

// CodeScanUpdate is the builder for updating CodeScan entities.
type CodeScanUpdate struct {
	config
	hooks    []Hook
	mutation *CodeScanMutation
}

// Where adds a new predicate for the CodeScanUpdate builder.
func (csu *CodeScanUpdate) Where(ps ...predicate.CodeScan) *CodeScanUpdate {
	csu.mutation.predicates = append(csu.mutation.predicates, ps...)
	return csu
}

// SetReleaseID sets the "release" edge to the Release entity by ID.
func (csu *CodeScanUpdate) SetReleaseID(id int) *CodeScanUpdate {
	csu.mutation.SetReleaseID(id)
	return csu
}

// SetRelease sets the "release" edge to the Release entity.
func (csu *CodeScanUpdate) SetRelease(r *Release) *CodeScanUpdate {
	return csu.SetReleaseID(r.ID)
}

// AddIssueIDs adds the "issues" edge to the CodeIssue entity by IDs.
func (csu *CodeScanUpdate) AddIssueIDs(ids ...int) *CodeScanUpdate {
	csu.mutation.AddIssueIDs(ids...)
	return csu
}

// AddIssues adds the "issues" edges to the CodeIssue entity.
func (csu *CodeScanUpdate) AddIssues(c ...*CodeIssue) *CodeScanUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return csu.AddIssueIDs(ids...)
}

// SetEntryID sets the "entry" edge to the ReleaseEntry entity by ID.
func (csu *CodeScanUpdate) SetEntryID(id int) *CodeScanUpdate {
	csu.mutation.SetEntryID(id)
	return csu
}

// SetNillableEntryID sets the "entry" edge to the ReleaseEntry entity by ID if the given value is not nil.
func (csu *CodeScanUpdate) SetNillableEntryID(id *int) *CodeScanUpdate {
	if id != nil {
		csu = csu.SetEntryID(*id)
	}
	return csu
}

// SetEntry sets the "entry" edge to the ReleaseEntry entity.
func (csu *CodeScanUpdate) SetEntry(r *ReleaseEntry) *CodeScanUpdate {
	return csu.SetEntryID(r.ID)
}

// Mutation returns the CodeScanMutation object of the builder.
func (csu *CodeScanUpdate) Mutation() *CodeScanMutation {
	return csu.mutation
}

// ClearRelease clears the "release" edge to the Release entity.
func (csu *CodeScanUpdate) ClearRelease() *CodeScanUpdate {
	csu.mutation.ClearRelease()
	return csu
}

// ClearIssues clears all "issues" edges to the CodeIssue entity.
func (csu *CodeScanUpdate) ClearIssues() *CodeScanUpdate {
	csu.mutation.ClearIssues()
	return csu
}

// RemoveIssueIDs removes the "issues" edge to CodeIssue entities by IDs.
func (csu *CodeScanUpdate) RemoveIssueIDs(ids ...int) *CodeScanUpdate {
	csu.mutation.RemoveIssueIDs(ids...)
	return csu
}

// RemoveIssues removes "issues" edges to CodeIssue entities.
func (csu *CodeScanUpdate) RemoveIssues(c ...*CodeIssue) *CodeScanUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return csu.RemoveIssueIDs(ids...)
}

// ClearEntry clears the "entry" edge to the ReleaseEntry entity.
func (csu *CodeScanUpdate) ClearEntry() *CodeScanUpdate {
	csu.mutation.ClearEntry()
	return csu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (csu *CodeScanUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(csu.hooks) == 0 {
		if err = csu.check(); err != nil {
			return 0, err
		}
		affected, err = csu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CodeScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = csu.check(); err != nil {
				return 0, err
			}
			csu.mutation = mutation
			affected, err = csu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(csu.hooks) - 1; i >= 0; i-- {
			mut = csu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, csu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (csu *CodeScanUpdate) SaveX(ctx context.Context) int {
	affected, err := csu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (csu *CodeScanUpdate) Exec(ctx context.Context) error {
	_, err := csu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (csu *CodeScanUpdate) ExecX(ctx context.Context) {
	if err := csu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (csu *CodeScanUpdate) check() error {
	if _, ok := csu.mutation.ReleaseID(); csu.mutation.ReleaseCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"release\"")
	}
	return nil
}

func (csu *CodeScanUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   codescan.Table,
			Columns: codescan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: codescan.FieldID,
			},
		},
	}
	if ps := csu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if csu.mutation.ReleaseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   codescan.ReleaseTable,
			Columns: []string{codescan.ReleaseColumn},
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
	if nodes := csu.mutation.ReleaseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   codescan.ReleaseTable,
			Columns: []string{codescan.ReleaseColumn},
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
	if csu.mutation.IssuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   codescan.IssuesTable,
			Columns: []string{codescan.IssuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: codeissue.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := csu.mutation.RemovedIssuesIDs(); len(nodes) > 0 && !csu.mutation.IssuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   codescan.IssuesTable,
			Columns: []string{codescan.IssuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: codeissue.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := csu.mutation.IssuesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   codescan.IssuesTable,
			Columns: []string{codescan.IssuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: codeissue.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if csu.mutation.EntryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   codescan.EntryTable,
			Columns: []string{codescan.EntryColumn},
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
	if nodes := csu.mutation.EntryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   codescan.EntryTable,
			Columns: []string{codescan.EntryColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, csu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{codescan.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// CodeScanUpdateOne is the builder for updating a single CodeScan entity.
type CodeScanUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *CodeScanMutation
}

// SetReleaseID sets the "release" edge to the Release entity by ID.
func (csuo *CodeScanUpdateOne) SetReleaseID(id int) *CodeScanUpdateOne {
	csuo.mutation.SetReleaseID(id)
	return csuo
}

// SetRelease sets the "release" edge to the Release entity.
func (csuo *CodeScanUpdateOne) SetRelease(r *Release) *CodeScanUpdateOne {
	return csuo.SetReleaseID(r.ID)
}

// AddIssueIDs adds the "issues" edge to the CodeIssue entity by IDs.
func (csuo *CodeScanUpdateOne) AddIssueIDs(ids ...int) *CodeScanUpdateOne {
	csuo.mutation.AddIssueIDs(ids...)
	return csuo
}

// AddIssues adds the "issues" edges to the CodeIssue entity.
func (csuo *CodeScanUpdateOne) AddIssues(c ...*CodeIssue) *CodeScanUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return csuo.AddIssueIDs(ids...)
}

// SetEntryID sets the "entry" edge to the ReleaseEntry entity by ID.
func (csuo *CodeScanUpdateOne) SetEntryID(id int) *CodeScanUpdateOne {
	csuo.mutation.SetEntryID(id)
	return csuo
}

// SetNillableEntryID sets the "entry" edge to the ReleaseEntry entity by ID if the given value is not nil.
func (csuo *CodeScanUpdateOne) SetNillableEntryID(id *int) *CodeScanUpdateOne {
	if id != nil {
		csuo = csuo.SetEntryID(*id)
	}
	return csuo
}

// SetEntry sets the "entry" edge to the ReleaseEntry entity.
func (csuo *CodeScanUpdateOne) SetEntry(r *ReleaseEntry) *CodeScanUpdateOne {
	return csuo.SetEntryID(r.ID)
}

// Mutation returns the CodeScanMutation object of the builder.
func (csuo *CodeScanUpdateOne) Mutation() *CodeScanMutation {
	return csuo.mutation
}

// ClearRelease clears the "release" edge to the Release entity.
func (csuo *CodeScanUpdateOne) ClearRelease() *CodeScanUpdateOne {
	csuo.mutation.ClearRelease()
	return csuo
}

// ClearIssues clears all "issues" edges to the CodeIssue entity.
func (csuo *CodeScanUpdateOne) ClearIssues() *CodeScanUpdateOne {
	csuo.mutation.ClearIssues()
	return csuo
}

// RemoveIssueIDs removes the "issues" edge to CodeIssue entities by IDs.
func (csuo *CodeScanUpdateOne) RemoveIssueIDs(ids ...int) *CodeScanUpdateOne {
	csuo.mutation.RemoveIssueIDs(ids...)
	return csuo
}

// RemoveIssues removes "issues" edges to CodeIssue entities.
func (csuo *CodeScanUpdateOne) RemoveIssues(c ...*CodeIssue) *CodeScanUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return csuo.RemoveIssueIDs(ids...)
}

// ClearEntry clears the "entry" edge to the ReleaseEntry entity.
func (csuo *CodeScanUpdateOne) ClearEntry() *CodeScanUpdateOne {
	csuo.mutation.ClearEntry()
	return csuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (csuo *CodeScanUpdateOne) Select(field string, fields ...string) *CodeScanUpdateOne {
	csuo.fields = append([]string{field}, fields...)
	return csuo
}

// Save executes the query and returns the updated CodeScan entity.
func (csuo *CodeScanUpdateOne) Save(ctx context.Context) (*CodeScan, error) {
	var (
		err  error
		node *CodeScan
	)
	if len(csuo.hooks) == 0 {
		if err = csuo.check(); err != nil {
			return nil, err
		}
		node, err = csuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CodeScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = csuo.check(); err != nil {
				return nil, err
			}
			csuo.mutation = mutation
			node, err = csuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(csuo.hooks) - 1; i >= 0; i-- {
			mut = csuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, csuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (csuo *CodeScanUpdateOne) SaveX(ctx context.Context) *CodeScan {
	node, err := csuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (csuo *CodeScanUpdateOne) Exec(ctx context.Context) error {
	_, err := csuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (csuo *CodeScanUpdateOne) ExecX(ctx context.Context) {
	if err := csuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (csuo *CodeScanUpdateOne) check() error {
	if _, ok := csuo.mutation.ReleaseID(); csuo.mutation.ReleaseCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"release\"")
	}
	return nil
}

func (csuo *CodeScanUpdateOne) sqlSave(ctx context.Context) (_node *CodeScan, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   codescan.Table,
			Columns: codescan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: codescan.FieldID,
			},
		},
	}
	id, ok := csuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing CodeScan.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := csuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, codescan.FieldID)
		for _, f := range fields {
			if !codescan.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != codescan.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := csuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if csuo.mutation.ReleaseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   codescan.ReleaseTable,
			Columns: []string{codescan.ReleaseColumn},
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
	if nodes := csuo.mutation.ReleaseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   codescan.ReleaseTable,
			Columns: []string{codescan.ReleaseColumn},
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
	if csuo.mutation.IssuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   codescan.IssuesTable,
			Columns: []string{codescan.IssuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: codeissue.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := csuo.mutation.RemovedIssuesIDs(); len(nodes) > 0 && !csuo.mutation.IssuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   codescan.IssuesTable,
			Columns: []string{codescan.IssuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: codeissue.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := csuo.mutation.IssuesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   codescan.IssuesTable,
			Columns: []string{codescan.IssuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: codeissue.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if csuo.mutation.EntryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   codescan.EntryTable,
			Columns: []string{codescan.EntryColumn},
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
	if nodes := csuo.mutation.EntryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   codescan.EntryTable,
			Columns: []string{codescan.EntryColumn},
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
	_node = &CodeScan{config: csuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, csuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{codescan.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
