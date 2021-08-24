// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/cwe"
	"github.com/valocode/bubbly/ent/predicate"
)

// CWEUpdate is the builder for updating CWE entities.
type CWEUpdate struct {
	config
	hooks    []Hook
	mutation *CWEMutation
}

// Where appends a list predicates to the CWEUpdate builder.
func (cu *CWEUpdate) Where(ps ...predicate.CWE) *CWEUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetCweID sets the "cwe_id" field.
func (cu *CWEUpdate) SetCweID(s string) *CWEUpdate {
	cu.mutation.SetCweID(s)
	return cu
}

// SetDescription sets the "description" field.
func (cu *CWEUpdate) SetDescription(s string) *CWEUpdate {
	cu.mutation.SetDescription(s)
	return cu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cu *CWEUpdate) SetNillableDescription(s *string) *CWEUpdate {
	if s != nil {
		cu.SetDescription(*s)
	}
	return cu
}

// ClearDescription clears the value of the "description" field.
func (cu *CWEUpdate) ClearDescription() *CWEUpdate {
	cu.mutation.ClearDescription()
	return cu
}

// SetURL sets the "url" field.
func (cu *CWEUpdate) SetURL(f float64) *CWEUpdate {
	cu.mutation.ResetURL()
	cu.mutation.SetURL(f)
	return cu
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (cu *CWEUpdate) SetNillableURL(f *float64) *CWEUpdate {
	if f != nil {
		cu.SetURL(*f)
	}
	return cu
}

// AddURL adds f to the "url" field.
func (cu *CWEUpdate) AddURL(f float64) *CWEUpdate {
	cu.mutation.AddURL(f)
	return cu
}

// ClearURL clears the value of the "url" field.
func (cu *CWEUpdate) ClearURL() *CWEUpdate {
	cu.mutation.ClearURL()
	return cu
}

// AddIssueIDs adds the "issues" edge to the CodeIssue entity by IDs.
func (cu *CWEUpdate) AddIssueIDs(ids ...int) *CWEUpdate {
	cu.mutation.AddIssueIDs(ids...)
	return cu
}

// AddIssues adds the "issues" edges to the CodeIssue entity.
func (cu *CWEUpdate) AddIssues(c ...*CodeIssue) *CWEUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cu.AddIssueIDs(ids...)
}

// Mutation returns the CWEMutation object of the builder.
func (cu *CWEUpdate) Mutation() *CWEMutation {
	return cu.mutation
}

// ClearIssues clears all "issues" edges to the CodeIssue entity.
func (cu *CWEUpdate) ClearIssues() *CWEUpdate {
	cu.mutation.ClearIssues()
	return cu
}

// RemoveIssueIDs removes the "issues" edge to CodeIssue entities by IDs.
func (cu *CWEUpdate) RemoveIssueIDs(ids ...int) *CWEUpdate {
	cu.mutation.RemoveIssueIDs(ids...)
	return cu
}

// RemoveIssues removes "issues" edges to CodeIssue entities.
func (cu *CWEUpdate) RemoveIssues(c ...*CodeIssue) *CWEUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cu.RemoveIssueIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *CWEUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(cu.hooks) == 0 {
		if err = cu.check(); err != nil {
			return 0, err
		}
		affected, err = cu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CWEMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = cu.check(); err != nil {
				return 0, err
			}
			cu.mutation = mutation
			affected, err = cu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(cu.hooks) - 1; i >= 0; i-- {
			if cu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = cu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CWEUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CWEUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CWEUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cu *CWEUpdate) check() error {
	if v, ok := cu.mutation.CweID(); ok {
		if err := cwe.CweIDValidator(v); err != nil {
			return &ValidationError{Name: "cwe_id", err: fmt.Errorf("ent: validator failed for field \"cwe_id\": %w", err)}
		}
	}
	return nil
}

func (cu *CWEUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cwe.Table,
			Columns: cwe.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: cwe.FieldID,
			},
		},
	}
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.CweID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cwe.FieldCweID,
		})
	}
	if value, ok := cu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cwe.FieldDescription,
		})
	}
	if cu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: cwe.FieldDescription,
		})
	}
	if value, ok := cu.mutation.URL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: cwe.FieldURL,
		})
	}
	if value, ok := cu.mutation.AddedURL(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: cwe.FieldURL,
		})
	}
	if cu.mutation.URLCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: cwe.FieldURL,
		})
	}
	if cu.mutation.IssuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cwe.IssuesTable,
			Columns: cwe.IssuesPrimaryKey,
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
	if nodes := cu.mutation.RemovedIssuesIDs(); len(nodes) > 0 && !cu.mutation.IssuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cwe.IssuesTable,
			Columns: cwe.IssuesPrimaryKey,
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
	if nodes := cu.mutation.IssuesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cwe.IssuesTable,
			Columns: cwe.IssuesPrimaryKey,
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
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cwe.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// CWEUpdateOne is the builder for updating a single CWE entity.
type CWEUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *CWEMutation
}

// SetCweID sets the "cwe_id" field.
func (cuo *CWEUpdateOne) SetCweID(s string) *CWEUpdateOne {
	cuo.mutation.SetCweID(s)
	return cuo
}

// SetDescription sets the "description" field.
func (cuo *CWEUpdateOne) SetDescription(s string) *CWEUpdateOne {
	cuo.mutation.SetDescription(s)
	return cuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cuo *CWEUpdateOne) SetNillableDescription(s *string) *CWEUpdateOne {
	if s != nil {
		cuo.SetDescription(*s)
	}
	return cuo
}

// ClearDescription clears the value of the "description" field.
func (cuo *CWEUpdateOne) ClearDescription() *CWEUpdateOne {
	cuo.mutation.ClearDescription()
	return cuo
}

// SetURL sets the "url" field.
func (cuo *CWEUpdateOne) SetURL(f float64) *CWEUpdateOne {
	cuo.mutation.ResetURL()
	cuo.mutation.SetURL(f)
	return cuo
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (cuo *CWEUpdateOne) SetNillableURL(f *float64) *CWEUpdateOne {
	if f != nil {
		cuo.SetURL(*f)
	}
	return cuo
}

// AddURL adds f to the "url" field.
func (cuo *CWEUpdateOne) AddURL(f float64) *CWEUpdateOne {
	cuo.mutation.AddURL(f)
	return cuo
}

// ClearURL clears the value of the "url" field.
func (cuo *CWEUpdateOne) ClearURL() *CWEUpdateOne {
	cuo.mutation.ClearURL()
	return cuo
}

// AddIssueIDs adds the "issues" edge to the CodeIssue entity by IDs.
func (cuo *CWEUpdateOne) AddIssueIDs(ids ...int) *CWEUpdateOne {
	cuo.mutation.AddIssueIDs(ids...)
	return cuo
}

// AddIssues adds the "issues" edges to the CodeIssue entity.
func (cuo *CWEUpdateOne) AddIssues(c ...*CodeIssue) *CWEUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cuo.AddIssueIDs(ids...)
}

// Mutation returns the CWEMutation object of the builder.
func (cuo *CWEUpdateOne) Mutation() *CWEMutation {
	return cuo.mutation
}

// ClearIssues clears all "issues" edges to the CodeIssue entity.
func (cuo *CWEUpdateOne) ClearIssues() *CWEUpdateOne {
	cuo.mutation.ClearIssues()
	return cuo
}

// RemoveIssueIDs removes the "issues" edge to CodeIssue entities by IDs.
func (cuo *CWEUpdateOne) RemoveIssueIDs(ids ...int) *CWEUpdateOne {
	cuo.mutation.RemoveIssueIDs(ids...)
	return cuo
}

// RemoveIssues removes "issues" edges to CodeIssue entities.
func (cuo *CWEUpdateOne) RemoveIssues(c ...*CodeIssue) *CWEUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cuo.RemoveIssueIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *CWEUpdateOne) Select(field string, fields ...string) *CWEUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated CWE entity.
func (cuo *CWEUpdateOne) Save(ctx context.Context) (*CWE, error) {
	var (
		err  error
		node *CWE
	)
	if len(cuo.hooks) == 0 {
		if err = cuo.check(); err != nil {
			return nil, err
		}
		node, err = cuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CWEMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = cuo.check(); err != nil {
				return nil, err
			}
			cuo.mutation = mutation
			node, err = cuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(cuo.hooks) - 1; i >= 0; i-- {
			if cuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = cuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CWEUpdateOne) SaveX(ctx context.Context) *CWE {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *CWEUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CWEUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cuo *CWEUpdateOne) check() error {
	if v, ok := cuo.mutation.CweID(); ok {
		if err := cwe.CweIDValidator(v); err != nil {
			return &ValidationError{Name: "cwe_id", err: fmt.Errorf("ent: validator failed for field \"cwe_id\": %w", err)}
		}
	}
	return nil
}

func (cuo *CWEUpdateOne) sqlSave(ctx context.Context) (_node *CWE, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cwe.Table,
			Columns: cwe.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: cwe.FieldID,
			},
		},
	}
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing CWE.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cwe.FieldID)
		for _, f := range fields {
			if !cwe.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != cwe.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.CweID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cwe.FieldCweID,
		})
	}
	if value, ok := cuo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cwe.FieldDescription,
		})
	}
	if cuo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: cwe.FieldDescription,
		})
	}
	if value, ok := cuo.mutation.URL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: cwe.FieldURL,
		})
	}
	if value, ok := cuo.mutation.AddedURL(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: cwe.FieldURL,
		})
	}
	if cuo.mutation.URLCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: cwe.FieldURL,
		})
	}
	if cuo.mutation.IssuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cwe.IssuesTable,
			Columns: cwe.IssuesPrimaryKey,
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
	if nodes := cuo.mutation.RemovedIssuesIDs(); len(nodes) > 0 && !cuo.mutation.IssuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cwe.IssuesTable,
			Columns: cwe.IssuesPrimaryKey,
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
	if nodes := cuo.mutation.IssuesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cwe.IssuesTable,
			Columns: cwe.IssuesPrimaryKey,
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
	_node = &CWE{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cwe.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}