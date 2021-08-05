// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/predicate"
)

// CodeIssueDelete is the builder for deleting a CodeIssue entity.
type CodeIssueDelete struct {
	config
	hooks    []Hook
	mutation *CodeIssueMutation
}

// Where adds a new predicate to the CodeIssueDelete builder.
func (cid *CodeIssueDelete) Where(ps ...predicate.CodeIssue) *CodeIssueDelete {
	cid.mutation.predicates = append(cid.mutation.predicates, ps...)
	return cid
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cid *CodeIssueDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(cid.hooks) == 0 {
		affected, err = cid.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CodeIssueMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			cid.mutation = mutation
			affected, err = cid.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(cid.hooks) - 1; i >= 0; i-- {
			mut = cid.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cid.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (cid *CodeIssueDelete) ExecX(ctx context.Context) int {
	n, err := cid.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cid *CodeIssueDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: codeissue.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: codeissue.FieldID,
			},
		},
	}
	if ps := cid.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, cid.driver, _spec)
}

// CodeIssueDeleteOne is the builder for deleting a single CodeIssue entity.
type CodeIssueDeleteOne struct {
	cid *CodeIssueDelete
}

// Exec executes the deletion query.
func (cido *CodeIssueDeleteOne) Exec(ctx context.Context) error {
	n, err := cido.cid.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{codeissue.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cido *CodeIssueDeleteOne) ExecX(ctx context.Context) {
	cido.cid.ExecX(ctx)
}
