// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/testcase"
	"github.com/valocode/bubbly/ent/testrun"
)

// TestCaseUpdate is the builder for updating TestCase entities.
type TestCaseUpdate struct {
	config
	hooks    []Hook
	mutation *TestCaseMutation
}

// Where appends a list predicates to the TestCaseUpdate builder.
func (tcu *TestCaseUpdate) Where(ps ...predicate.TestCase) *TestCaseUpdate {
	tcu.mutation.Where(ps...)
	return tcu
}

// SetResult sets the "result" field.
func (tcu *TestCaseUpdate) SetResult(b bool) *TestCaseUpdate {
	tcu.mutation.SetResult(b)
	return tcu
}

// SetMessage sets the "message" field.
func (tcu *TestCaseUpdate) SetMessage(s string) *TestCaseUpdate {
	tcu.mutation.SetMessage(s)
	return tcu
}

// SetElapsed sets the "elapsed" field.
func (tcu *TestCaseUpdate) SetElapsed(f float64) *TestCaseUpdate {
	tcu.mutation.ResetElapsed()
	tcu.mutation.SetElapsed(f)
	return tcu
}

// SetNillableElapsed sets the "elapsed" field if the given value is not nil.
func (tcu *TestCaseUpdate) SetNillableElapsed(f *float64) *TestCaseUpdate {
	if f != nil {
		tcu.SetElapsed(*f)
	}
	return tcu
}

// AddElapsed adds f to the "elapsed" field.
func (tcu *TestCaseUpdate) AddElapsed(f float64) *TestCaseUpdate {
	tcu.mutation.AddElapsed(f)
	return tcu
}

// SetRunID sets the "run" edge to the TestRun entity by ID.
func (tcu *TestCaseUpdate) SetRunID(id int) *TestCaseUpdate {
	tcu.mutation.SetRunID(id)
	return tcu
}

// SetRun sets the "run" edge to the TestRun entity.
func (tcu *TestCaseUpdate) SetRun(t *TestRun) *TestCaseUpdate {
	return tcu.SetRunID(t.ID)
}

// Mutation returns the TestCaseMutation object of the builder.
func (tcu *TestCaseUpdate) Mutation() *TestCaseMutation {
	return tcu.mutation
}

// ClearRun clears the "run" edge to the TestRun entity.
func (tcu *TestCaseUpdate) ClearRun() *TestCaseUpdate {
	tcu.mutation.ClearRun()
	return tcu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tcu *TestCaseUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(tcu.hooks) == 0 {
		if err = tcu.check(); err != nil {
			return 0, err
		}
		affected, err = tcu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TestCaseMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = tcu.check(); err != nil {
				return 0, err
			}
			tcu.mutation = mutation
			affected, err = tcu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(tcu.hooks) - 1; i >= 0; i-- {
			if tcu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = tcu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tcu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (tcu *TestCaseUpdate) SaveX(ctx context.Context) int {
	affected, err := tcu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tcu *TestCaseUpdate) Exec(ctx context.Context) error {
	_, err := tcu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcu *TestCaseUpdate) ExecX(ctx context.Context) {
	if err := tcu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tcu *TestCaseUpdate) check() error {
	if v, ok := tcu.mutation.Message(); ok {
		if err := testcase.MessageValidator(v); err != nil {
			return &ValidationError{Name: "message", err: fmt.Errorf("ent: validator failed for field \"message\": %w", err)}
		}
	}
	if v, ok := tcu.mutation.Elapsed(); ok {
		if err := testcase.ElapsedValidator(v); err != nil {
			return &ValidationError{Name: "elapsed", err: fmt.Errorf("ent: validator failed for field \"elapsed\": %w", err)}
		}
	}
	if _, ok := tcu.mutation.RunID(); tcu.mutation.RunCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"run\"")
	}
	return nil
}

func (tcu *TestCaseUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   testcase.Table,
			Columns: testcase.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: testcase.FieldID,
			},
		},
	}
	if ps := tcu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tcu.mutation.Result(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: testcase.FieldResult,
		})
	}
	if value, ok := tcu.mutation.Message(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: testcase.FieldMessage,
		})
	}
	if value, ok := tcu.mutation.Elapsed(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: testcase.FieldElapsed,
		})
	}
	if value, ok := tcu.mutation.AddedElapsed(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: testcase.FieldElapsed,
		})
	}
	if tcu.mutation.RunCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   testcase.RunTable,
			Columns: []string{testcase.RunColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: testrun.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tcu.mutation.RunIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   testcase.RunTable,
			Columns: []string{testcase.RunColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: testrun.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tcu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{testcase.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// TestCaseUpdateOne is the builder for updating a single TestCase entity.
type TestCaseUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TestCaseMutation
}

// SetResult sets the "result" field.
func (tcuo *TestCaseUpdateOne) SetResult(b bool) *TestCaseUpdateOne {
	tcuo.mutation.SetResult(b)
	return tcuo
}

// SetMessage sets the "message" field.
func (tcuo *TestCaseUpdateOne) SetMessage(s string) *TestCaseUpdateOne {
	tcuo.mutation.SetMessage(s)
	return tcuo
}

// SetElapsed sets the "elapsed" field.
func (tcuo *TestCaseUpdateOne) SetElapsed(f float64) *TestCaseUpdateOne {
	tcuo.mutation.ResetElapsed()
	tcuo.mutation.SetElapsed(f)
	return tcuo
}

// SetNillableElapsed sets the "elapsed" field if the given value is not nil.
func (tcuo *TestCaseUpdateOne) SetNillableElapsed(f *float64) *TestCaseUpdateOne {
	if f != nil {
		tcuo.SetElapsed(*f)
	}
	return tcuo
}

// AddElapsed adds f to the "elapsed" field.
func (tcuo *TestCaseUpdateOne) AddElapsed(f float64) *TestCaseUpdateOne {
	tcuo.mutation.AddElapsed(f)
	return tcuo
}

// SetRunID sets the "run" edge to the TestRun entity by ID.
func (tcuo *TestCaseUpdateOne) SetRunID(id int) *TestCaseUpdateOne {
	tcuo.mutation.SetRunID(id)
	return tcuo
}

// SetRun sets the "run" edge to the TestRun entity.
func (tcuo *TestCaseUpdateOne) SetRun(t *TestRun) *TestCaseUpdateOne {
	return tcuo.SetRunID(t.ID)
}

// Mutation returns the TestCaseMutation object of the builder.
func (tcuo *TestCaseUpdateOne) Mutation() *TestCaseMutation {
	return tcuo.mutation
}

// ClearRun clears the "run" edge to the TestRun entity.
func (tcuo *TestCaseUpdateOne) ClearRun() *TestCaseUpdateOne {
	tcuo.mutation.ClearRun()
	return tcuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tcuo *TestCaseUpdateOne) Select(field string, fields ...string) *TestCaseUpdateOne {
	tcuo.fields = append([]string{field}, fields...)
	return tcuo
}

// Save executes the query and returns the updated TestCase entity.
func (tcuo *TestCaseUpdateOne) Save(ctx context.Context) (*TestCase, error) {
	var (
		err  error
		node *TestCase
	)
	if len(tcuo.hooks) == 0 {
		if err = tcuo.check(); err != nil {
			return nil, err
		}
		node, err = tcuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TestCaseMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = tcuo.check(); err != nil {
				return nil, err
			}
			tcuo.mutation = mutation
			node, err = tcuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(tcuo.hooks) - 1; i >= 0; i-- {
			if tcuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = tcuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tcuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (tcuo *TestCaseUpdateOne) SaveX(ctx context.Context) *TestCase {
	node, err := tcuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tcuo *TestCaseUpdateOne) Exec(ctx context.Context) error {
	_, err := tcuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcuo *TestCaseUpdateOne) ExecX(ctx context.Context) {
	if err := tcuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tcuo *TestCaseUpdateOne) check() error {
	if v, ok := tcuo.mutation.Message(); ok {
		if err := testcase.MessageValidator(v); err != nil {
			return &ValidationError{Name: "message", err: fmt.Errorf("ent: validator failed for field \"message\": %w", err)}
		}
	}
	if v, ok := tcuo.mutation.Elapsed(); ok {
		if err := testcase.ElapsedValidator(v); err != nil {
			return &ValidationError{Name: "elapsed", err: fmt.Errorf("ent: validator failed for field \"elapsed\": %w", err)}
		}
	}
	if _, ok := tcuo.mutation.RunID(); tcuo.mutation.RunCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"run\"")
	}
	return nil
}

func (tcuo *TestCaseUpdateOne) sqlSave(ctx context.Context) (_node *TestCase, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   testcase.Table,
			Columns: testcase.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: testcase.FieldID,
			},
		},
	}
	id, ok := tcuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing TestCase.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := tcuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, testcase.FieldID)
		for _, f := range fields {
			if !testcase.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != testcase.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tcuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tcuo.mutation.Result(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: testcase.FieldResult,
		})
	}
	if value, ok := tcuo.mutation.Message(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: testcase.FieldMessage,
		})
	}
	if value, ok := tcuo.mutation.Elapsed(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: testcase.FieldElapsed,
		})
	}
	if value, ok := tcuo.mutation.AddedElapsed(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: testcase.FieldElapsed,
		})
	}
	if tcuo.mutation.RunCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   testcase.RunTable,
			Columns: []string{testcase.RunColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: testrun.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tcuo.mutation.RunIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   testcase.RunTable,
			Columns: []string{testcase.RunColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: testrun.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &TestCase{config: tcuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tcuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{testcase.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}