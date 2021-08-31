// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/releasepolicy"
	"github.com/valocode/bubbly/ent/releasepolicyviolation"
	"github.com/valocode/bubbly/ent/repo"
)

// ReleasePolicyCreate is the builder for creating a ReleasePolicy entity.
type ReleasePolicyCreate struct {
	config
	mutation *ReleasePolicyMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (rpc *ReleasePolicyCreate) SetName(s string) *ReleasePolicyCreate {
	rpc.mutation.SetName(s)
	return rpc
}

// SetModule sets the "module" field.
func (rpc *ReleasePolicyCreate) SetModule(s string) *ReleasePolicyCreate {
	rpc.mutation.SetModule(s)
	return rpc
}

// AddProjectIDs adds the "projects" edge to the Project entity by IDs.
func (rpc *ReleasePolicyCreate) AddProjectIDs(ids ...int) *ReleasePolicyCreate {
	rpc.mutation.AddProjectIDs(ids...)
	return rpc
}

// AddProjects adds the "projects" edges to the Project entity.
func (rpc *ReleasePolicyCreate) AddProjects(p ...*Project) *ReleasePolicyCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return rpc.AddProjectIDs(ids...)
}

// AddRepoIDs adds the "repos" edge to the Repo entity by IDs.
func (rpc *ReleasePolicyCreate) AddRepoIDs(ids ...int) *ReleasePolicyCreate {
	rpc.mutation.AddRepoIDs(ids...)
	return rpc
}

// AddRepos adds the "repos" edges to the Repo entity.
func (rpc *ReleasePolicyCreate) AddRepos(r ...*Repo) *ReleasePolicyCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rpc.AddRepoIDs(ids...)
}

// AddViolationIDs adds the "violations" edge to the ReleasePolicyViolation entity by IDs.
func (rpc *ReleasePolicyCreate) AddViolationIDs(ids ...int) *ReleasePolicyCreate {
	rpc.mutation.AddViolationIDs(ids...)
	return rpc
}

// AddViolations adds the "violations" edges to the ReleasePolicyViolation entity.
func (rpc *ReleasePolicyCreate) AddViolations(r ...*ReleasePolicyViolation) *ReleasePolicyCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rpc.AddViolationIDs(ids...)
}

// Mutation returns the ReleasePolicyMutation object of the builder.
func (rpc *ReleasePolicyCreate) Mutation() *ReleasePolicyMutation {
	return rpc.mutation
}

// Save creates the ReleasePolicy in the database.
func (rpc *ReleasePolicyCreate) Save(ctx context.Context) (*ReleasePolicy, error) {
	var (
		err  error
		node *ReleasePolicy
	)
	if len(rpc.hooks) == 0 {
		if err = rpc.check(); err != nil {
			return nil, err
		}
		node, err = rpc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ReleasePolicyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rpc.check(); err != nil {
				return nil, err
			}
			rpc.mutation = mutation
			if node, err = rpc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(rpc.hooks) - 1; i >= 0; i-- {
			if rpc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = rpc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rpc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rpc *ReleasePolicyCreate) SaveX(ctx context.Context) *ReleasePolicy {
	v, err := rpc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rpc *ReleasePolicyCreate) Exec(ctx context.Context) error {
	_, err := rpc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rpc *ReleasePolicyCreate) ExecX(ctx context.Context) {
	if err := rpc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rpc *ReleasePolicyCreate) check() error {
	if _, ok := rpc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "name"`)}
	}
	if v, ok := rpc.mutation.Name(); ok {
		if err := releasepolicy.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "name": %w`, err)}
		}
	}
	if _, ok := rpc.mutation.Module(); !ok {
		return &ValidationError{Name: "module", err: errors.New(`ent: missing required field "module"`)}
	}
	if v, ok := rpc.mutation.Module(); ok {
		if err := releasepolicy.ModuleValidator(v); err != nil {
			return &ValidationError{Name: "module", err: fmt.Errorf(`ent: validator failed for field "module": %w`, err)}
		}
	}
	return nil
}

func (rpc *ReleasePolicyCreate) sqlSave(ctx context.Context) (*ReleasePolicy, error) {
	_node, _spec := rpc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rpc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (rpc *ReleasePolicyCreate) createSpec() (*ReleasePolicy, *sqlgraph.CreateSpec) {
	var (
		_node = &ReleasePolicy{config: rpc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: releasepolicy.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: releasepolicy.FieldID,
			},
		}
	)
	if value, ok := rpc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: releasepolicy.FieldName,
		})
		_node.Name = value
	}
	if value, ok := rpc.mutation.Module(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: releasepolicy.FieldModule,
		})
		_node.Module = value
	}
	if nodes := rpc.mutation.ProjectsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   releasepolicy.ProjectsTable,
			Columns: releasepolicy.ProjectsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rpc.mutation.ReposIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   releasepolicy.ReposTable,
			Columns: releasepolicy.ReposPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repo.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rpc.mutation.ViolationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   releasepolicy.ViolationsTable,
			Columns: []string{releasepolicy.ViolationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releasepolicyviolation.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ReleasePolicyCreateBulk is the builder for creating many ReleasePolicy entities in bulk.
type ReleasePolicyCreateBulk struct {
	config
	builders []*ReleasePolicyCreate
}

// Save creates the ReleasePolicy entities in the database.
func (rpcb *ReleasePolicyCreateBulk) Save(ctx context.Context) ([]*ReleasePolicy, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rpcb.builders))
	nodes := make([]*ReleasePolicy, len(rpcb.builders))
	mutators := make([]Mutator, len(rpcb.builders))
	for i := range rpcb.builders {
		func(i int, root context.Context) {
			builder := rpcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ReleasePolicyMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, rpcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rpcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, rpcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rpcb *ReleasePolicyCreateBulk) SaveX(ctx context.Context) []*ReleasePolicy {
	v, err := rpcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rpcb *ReleasePolicyCreateBulk) Exec(ctx context.Context) error {
	_, err := rpcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rpcb *ReleasePolicyCreateBulk) ExecX(ctx context.Context) {
	if err := rpcb.Exec(ctx); err != nil {
		panic(err)
	}
}