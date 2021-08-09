// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/cve"
	"github.com/valocode/bubbly/ent/cverule"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/repo"
)

// CVERuleCreate is the builder for creating a CVERule entity.
type CVERuleCreate struct {
	config
	mutation *CVERuleMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (crc *CVERuleCreate) SetName(s string) *CVERuleCreate {
	crc.mutation.SetName(s)
	return crc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (crc *CVERuleCreate) SetNillableName(s *string) *CVERuleCreate {
	if s != nil {
		crc.SetName(*s)
	}
	return crc
}

// SetCveID sets the "cve" edge to the CVE entity by ID.
func (crc *CVERuleCreate) SetCveID(id int) *CVERuleCreate {
	crc.mutation.SetCveID(id)
	return crc
}

// SetCve sets the "cve" edge to the CVE entity.
func (crc *CVERuleCreate) SetCve(c *CVE) *CVERuleCreate {
	return crc.SetCveID(c.ID)
}

// AddProjectIDs adds the "project" edge to the Project entity by IDs.
func (crc *CVERuleCreate) AddProjectIDs(ids ...int) *CVERuleCreate {
	crc.mutation.AddProjectIDs(ids...)
	return crc
}

// AddProject adds the "project" edges to the Project entity.
func (crc *CVERuleCreate) AddProject(p ...*Project) *CVERuleCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return crc.AddProjectIDs(ids...)
}

// AddRepoIDs adds the "repo" edge to the Repo entity by IDs.
func (crc *CVERuleCreate) AddRepoIDs(ids ...int) *CVERuleCreate {
	crc.mutation.AddRepoIDs(ids...)
	return crc
}

// AddRepo adds the "repo" edges to the Repo entity.
func (crc *CVERuleCreate) AddRepo(r ...*Repo) *CVERuleCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return crc.AddRepoIDs(ids...)
}

// Mutation returns the CVERuleMutation object of the builder.
func (crc *CVERuleCreate) Mutation() *CVERuleMutation {
	return crc.mutation
}

// Save creates the CVERule in the database.
func (crc *CVERuleCreate) Save(ctx context.Context) (*CVERule, error) {
	var (
		err  error
		node *CVERule
	)
	if len(crc.hooks) == 0 {
		if err = crc.check(); err != nil {
			return nil, err
		}
		node, err = crc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CVERuleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = crc.check(); err != nil {
				return nil, err
			}
			crc.mutation = mutation
			if node, err = crc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(crc.hooks) - 1; i >= 0; i-- {
			if crc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = crc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, crc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (crc *CVERuleCreate) SaveX(ctx context.Context) *CVERule {
	v, err := crc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (crc *CVERuleCreate) Exec(ctx context.Context) error {
	_, err := crc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (crc *CVERuleCreate) ExecX(ctx context.Context) {
	if err := crc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (crc *CVERuleCreate) check() error {
	if _, ok := crc.mutation.CveID(); !ok {
		return &ValidationError{Name: "cve", err: errors.New("ent: missing required edge \"cve\"")}
	}
	return nil
}

func (crc *CVERuleCreate) sqlSave(ctx context.Context) (*CVERule, error) {
	_node, _spec := crc.createSpec()
	if err := sqlgraph.CreateNode(ctx, crc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (crc *CVERuleCreate) createSpec() (*CVERule, *sqlgraph.CreateSpec) {
	var (
		_node = &CVERule{config: crc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: cverule.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: cverule.FieldID,
			},
		}
	)
	if value, ok := crc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cverule.FieldName,
		})
		_node.Name = value
	}
	if nodes := crc.mutation.CveIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   cverule.CveTable,
			Columns: []string{cverule.CveColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: cve.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.cve_rule_cve = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := crc.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cverule.ProjectTable,
			Columns: cverule.ProjectPrimaryKey,
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
	if nodes := crc.mutation.RepoIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cverule.RepoTable,
			Columns: cverule.RepoPrimaryKey,
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
	return _node, _spec
}

// CVERuleCreateBulk is the builder for creating many CVERule entities in bulk.
type CVERuleCreateBulk struct {
	config
	builders []*CVERuleCreate
}

// Save creates the CVERule entities in the database.
func (crcb *CVERuleCreateBulk) Save(ctx context.Context) ([]*CVERule, error) {
	specs := make([]*sqlgraph.CreateSpec, len(crcb.builders))
	nodes := make([]*CVERule, len(crcb.builders))
	mutators := make([]Mutator, len(crcb.builders))
	for i := range crcb.builders {
		func(i int, root context.Context) {
			builder := crcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CVERuleMutation)
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
					_, err = mutators[i+1].Mutate(root, crcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, crcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, crcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (crcb *CVERuleCreateBulk) SaveX(ctx context.Context) []*CVERule {
	v, err := crcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (crcb *CVERuleCreateBulk) Exec(ctx context.Context) error {
	_, err := crcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (crcb *CVERuleCreateBulk) ExecX(ctx context.Context) {
	if err := crcb.Exec(ctx); err != nil {
		panic(err)
	}
}
