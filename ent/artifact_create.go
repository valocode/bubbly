// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/artifact"

	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releaseentry"
	schema "github.com/valocode/bubbly/ent/schema/types"
)

// ArtifactCreate is the builder for creating a Artifact entity.
type ArtifactCreate struct {
	config
	mutation *ArtifactMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (ac *ArtifactCreate) SetName(s string) *ArtifactCreate {
	ac.mutation.SetName(s)
	return ac
}

// SetSha256 sets the "sha256" field.
func (ac *ArtifactCreate) SetSha256(s string) *ArtifactCreate {
	ac.mutation.SetSha256(s)
	return ac
}

// SetType sets the "type" field.
func (ac *ArtifactCreate) SetType(a artifact.Type) *ArtifactCreate {
	ac.mutation.SetType(a)
	return ac
}

// SetTime sets the "time" field.
func (ac *ArtifactCreate) SetTime(t time.Time) *ArtifactCreate {
	ac.mutation.SetTime(t)
	return ac
}

// SetNillableTime sets the "time" field if the given value is not nil.
func (ac *ArtifactCreate) SetNillableTime(t *time.Time) *ArtifactCreate {
	if t != nil {
		ac.SetTime(*t)
	}
	return ac
}

// SetMetadata sets the "metadata" field.
func (ac *ArtifactCreate) SetMetadata(s schema.Metadata) *ArtifactCreate {
	ac.mutation.SetMetadata(s)
	return ac
}

// SetReleaseID sets the "release" edge to the Release entity by ID.
func (ac *ArtifactCreate) SetReleaseID(id int) *ArtifactCreate {
	ac.mutation.SetReleaseID(id)
	return ac
}

// SetNillableReleaseID sets the "release" edge to the Release entity by ID if the given value is not nil.
func (ac *ArtifactCreate) SetNillableReleaseID(id *int) *ArtifactCreate {
	if id != nil {
		ac = ac.SetReleaseID(*id)
	}
	return ac
}

// SetRelease sets the "release" edge to the Release entity.
func (ac *ArtifactCreate) SetRelease(r *Release) *ArtifactCreate {
	return ac.SetReleaseID(r.ID)
}

// SetEntryID sets the "entry" edge to the ReleaseEntry entity by ID.
func (ac *ArtifactCreate) SetEntryID(id int) *ArtifactCreate {
	ac.mutation.SetEntryID(id)
	return ac
}

// SetNillableEntryID sets the "entry" edge to the ReleaseEntry entity by ID if the given value is not nil.
func (ac *ArtifactCreate) SetNillableEntryID(id *int) *ArtifactCreate {
	if id != nil {
		ac = ac.SetEntryID(*id)
	}
	return ac
}

// SetEntry sets the "entry" edge to the ReleaseEntry entity.
func (ac *ArtifactCreate) SetEntry(r *ReleaseEntry) *ArtifactCreate {
	return ac.SetEntryID(r.ID)
}

// Mutation returns the ArtifactMutation object of the builder.
func (ac *ArtifactCreate) Mutation() *ArtifactMutation {
	return ac.mutation
}

// Save creates the Artifact in the database.
func (ac *ArtifactCreate) Save(ctx context.Context) (*Artifact, error) {
	var (
		err  error
		node *Artifact
	)
	if err := ac.defaults(); err != nil {
		return nil, err
	}
	if len(ac.hooks) == 0 {
		if err = ac.check(); err != nil {
			return nil, err
		}
		node, err = ac.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ArtifactMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ac.check(); err != nil {
				return nil, err
			}
			ac.mutation = mutation
			if node, err = ac.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(ac.hooks) - 1; i >= 0; i-- {
			if ac.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ac.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ac.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ac *ArtifactCreate) SaveX(ctx context.Context) *Artifact {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ac *ArtifactCreate) Exec(ctx context.Context) error {
	_, err := ac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ac *ArtifactCreate) ExecX(ctx context.Context) {
	if err := ac.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ac *ArtifactCreate) defaults() error {
	if _, ok := ac.mutation.Time(); !ok {
		if artifact.DefaultTime == nil {
			return fmt.Errorf("ent: uninitialized artifact.DefaultTime (forgotten import ent/runtime?)")
		}
		v := artifact.DefaultTime()
		ac.mutation.SetTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (ac *ArtifactCreate) check() error {
	if _, ok := ac.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "name"`)}
	}
	if v, ok := ac.mutation.Name(); ok {
		if err := artifact.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "name": %w`, err)}
		}
	}
	if _, ok := ac.mutation.Sha256(); !ok {
		return &ValidationError{Name: "sha256", err: errors.New(`ent: missing required field "sha256"`)}
	}
	if v, ok := ac.mutation.Sha256(); ok {
		if err := artifact.Sha256Validator(v); err != nil {
			return &ValidationError{Name: "sha256", err: fmt.Errorf(`ent: validator failed for field "sha256": %w`, err)}
		}
	}
	if _, ok := ac.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "type"`)}
	}
	if v, ok := ac.mutation.GetType(); ok {
		if err := artifact.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "type": %w`, err)}
		}
	}
	if _, ok := ac.mutation.Time(); !ok {
		return &ValidationError{Name: "time", err: errors.New(`ent: missing required field "time"`)}
	}
	return nil
}

func (ac *ArtifactCreate) sqlSave(ctx context.Context) (*Artifact, error) {
	_node, _spec := ac.createSpec()
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (ac *ArtifactCreate) createSpec() (*Artifact, *sqlgraph.CreateSpec) {
	var (
		_node = &Artifact{config: ac.config}
		_spec = &sqlgraph.CreateSpec{
			Table: artifact.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: artifact.FieldID,
			},
		}
	)
	if value, ok := ac.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: artifact.FieldName,
		})
		_node.Name = value
	}
	if value, ok := ac.mutation.Sha256(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: artifact.FieldSha256,
		})
		_node.Sha256 = value
	}
	if value, ok := ac.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: artifact.FieldType,
		})
		_node.Type = value
	}
	if value, ok := ac.mutation.Time(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: artifact.FieldTime,
		})
		_node.Time = value
	}
	if value, ok := ac.mutation.Metadata(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: artifact.FieldMetadata,
		})
		_node.Metadata = value
	}
	if nodes := ac.mutation.ReleaseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   artifact.ReleaseTable,
			Columns: []string{artifact.ReleaseColumn},
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
		_node.artifact_release = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ac.mutation.EntryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   artifact.EntryTable,
			Columns: []string{artifact.EntryColumn},
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
		_node.release_entry_artifact = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ArtifactCreateBulk is the builder for creating many Artifact entities in bulk.
type ArtifactCreateBulk struct {
	config
	builders []*ArtifactCreate
}

// Save creates the Artifact entities in the database.
func (acb *ArtifactCreateBulk) Save(ctx context.Context) ([]*Artifact, error) {
	specs := make([]*sqlgraph.CreateSpec, len(acb.builders))
	nodes := make([]*Artifact, len(acb.builders))
	mutators := make([]Mutator, len(acb.builders))
	for i := range acb.builders {
		func(i int, root context.Context) {
			builder := acb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ArtifactMutation)
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
					_, err = mutators[i+1].Mutate(root, acb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, acb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, acb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (acb *ArtifactCreateBulk) SaveX(ctx context.Context) []*Artifact {
	v, err := acb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (acb *ArtifactCreateBulk) Exec(ctx context.Context) error {
	_, err := acb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acb *ArtifactCreateBulk) ExecX(ctx context.Context) {
	if err := acb.Exec(ctx); err != nil {
		panic(err)
	}
}
