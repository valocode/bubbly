// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/cverule"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/repo"
)

// RepoCreate is the builder for creating a Repo entity.
type RepoCreate struct {
	config
	mutation *RepoMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (rc *RepoCreate) SetName(s string) *RepoCreate {
	rc.mutation.SetName(s)
	return rc
}

// SetProjectID sets the "project" edge to the Project entity by ID.
func (rc *RepoCreate) SetProjectID(id int) *RepoCreate {
	rc.mutation.SetProjectID(id)
	return rc
}

// SetNillableProjectID sets the "project" edge to the Project entity by ID if the given value is not nil.
func (rc *RepoCreate) SetNillableProjectID(id *int) *RepoCreate {
	if id != nil {
		rc = rc.SetProjectID(*id)
	}
	return rc
}

// SetProject sets the "project" edge to the Project entity.
func (rc *RepoCreate) SetProject(p *Project) *RepoCreate {
	return rc.SetProjectID(p.ID)
}

// AddCommitIDs adds the "commits" edge to the GitCommit entity by IDs.
func (rc *RepoCreate) AddCommitIDs(ids ...int) *RepoCreate {
	rc.mutation.AddCommitIDs(ids...)
	return rc
}

// AddCommits adds the "commits" edges to the GitCommit entity.
func (rc *RepoCreate) AddCommits(g ...*GitCommit) *RepoCreate {
	ids := make([]int, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return rc.AddCommitIDs(ids...)
}

// AddCveRuleIDs adds the "cve_rules" edge to the CVERule entity by IDs.
func (rc *RepoCreate) AddCveRuleIDs(ids ...int) *RepoCreate {
	rc.mutation.AddCveRuleIDs(ids...)
	return rc
}

// AddCveRules adds the "cve_rules" edges to the CVERule entity.
func (rc *RepoCreate) AddCveRules(c ...*CVERule) *RepoCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return rc.AddCveRuleIDs(ids...)
}

// Mutation returns the RepoMutation object of the builder.
func (rc *RepoCreate) Mutation() *RepoMutation {
	return rc.mutation
}

// Save creates the Repo in the database.
func (rc *RepoCreate) Save(ctx context.Context) (*Repo, error) {
	var (
		err  error
		node *Repo
	)
	if len(rc.hooks) == 0 {
		if err = rc.check(); err != nil {
			return nil, err
		}
		node, err = rc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RepoMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rc.check(); err != nil {
				return nil, err
			}
			rc.mutation = mutation
			if node, err = rc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(rc.hooks) - 1; i >= 0; i-- {
			mut = rc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rc *RepoCreate) SaveX(ctx context.Context) *Repo {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// check runs all checks and user-defined validators on the builder.
func (rc *RepoCreate) check() error {
	if _, ok := rc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New("ent: missing required field \"name\"")}
	}
	if v, ok := rc.mutation.Name(); ok {
		if err := repo.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	return nil
}

func (rc *RepoCreate) sqlSave(ctx context.Context) (*Repo, error) {
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (rc *RepoCreate) createSpec() (*Repo, *sqlgraph.CreateSpec) {
	var (
		_node = &Repo{config: rc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: repo.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: repo.FieldID,
			},
		}
	)
	if value, ok := rc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: repo.FieldName,
		})
		_node.Name = value
	}
	if nodes := rc.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   repo.ProjectTable,
			Columns: []string{repo.ProjectColumn},
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
		_node.repo_project = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.CommitsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   repo.CommitsTable,
			Columns: []string{repo.CommitsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: gitcommit.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.CveRulesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.CveRulesTable,
			Columns: repo.CveRulesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: cverule.FieldID,
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

// RepoCreateBulk is the builder for creating many Repo entities in bulk.
type RepoCreateBulk struct {
	config
	builders []*RepoCreate
}

// Save creates the Repo entities in the database.
func (rcb *RepoCreateBulk) Save(ctx context.Context) ([]*Repo, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Repo, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RepoMutation)
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
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
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
				id := specs[i].ID.Value.(int64)
				nodes[i].ID = int(id)
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *RepoCreateBulk) SaveX(ctx context.Context) []*Repo {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
