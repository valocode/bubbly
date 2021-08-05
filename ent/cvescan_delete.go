// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/cvescan"
	"github.com/valocode/bubbly/ent/predicate"
)

// CVEScanDelete is the builder for deleting a CVEScan entity.
type CVEScanDelete struct {
	config
	hooks    []Hook
	mutation *CVEScanMutation
}

// Where adds a new predicate to the CVEScanDelete builder.
func (csd *CVEScanDelete) Where(ps ...predicate.CVEScan) *CVEScanDelete {
	csd.mutation.predicates = append(csd.mutation.predicates, ps...)
	return csd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (csd *CVEScanDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(csd.hooks) == 0 {
		affected, err = csd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CVEScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			csd.mutation = mutation
			affected, err = csd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(csd.hooks) - 1; i >= 0; i-- {
			mut = csd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, csd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (csd *CVEScanDelete) ExecX(ctx context.Context) int {
	n, err := csd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (csd *CVEScanDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: cvescan.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: cvescan.FieldID,
			},
		},
	}
	if ps := csd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, csd.driver, _spec)
}

// CVEScanDeleteOne is the builder for deleting a single CVEScan entity.
type CVEScanDeleteOne struct {
	csd *CVEScanDelete
}

// Exec executes the deletion query.
func (csdo *CVEScanDeleteOne) Exec(ctx context.Context) error {
	n, err := csdo.csd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{cvescan.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (csdo *CVEScanDeleteOne) ExecX(ctx context.Context) {
	csdo.csd.ExecX(ctx)
}
