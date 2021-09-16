// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/license"
	"github.com/valocode/bubbly/ent/organization"
	"github.com/valocode/bubbly/ent/releaselicense"
)

// LicenseCreate is the builder for creating a License entity.
type LicenseCreate struct {
	config
	mutation *LicenseMutation
	hooks    []Hook
}

// SetLicenseID sets the "license_id" field.
func (lc *LicenseCreate) SetLicenseID(s string) *LicenseCreate {
	lc.mutation.SetLicenseID(s)
	return lc
}

// SetName sets the "name" field.
func (lc *LicenseCreate) SetName(s string) *LicenseCreate {
	lc.mutation.SetName(s)
	return lc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (lc *LicenseCreate) SetNillableName(s *string) *LicenseCreate {
	if s != nil {
		lc.SetName(*s)
	}
	return lc
}

// SetReference sets the "reference" field.
func (lc *LicenseCreate) SetReference(s string) *LicenseCreate {
	lc.mutation.SetReference(s)
	return lc
}

// SetNillableReference sets the "reference" field if the given value is not nil.
func (lc *LicenseCreate) SetNillableReference(s *string) *LicenseCreate {
	if s != nil {
		lc.SetReference(*s)
	}
	return lc
}

// SetDetailsURL sets the "details_url" field.
func (lc *LicenseCreate) SetDetailsURL(s string) *LicenseCreate {
	lc.mutation.SetDetailsURL(s)
	return lc
}

// SetNillableDetailsURL sets the "details_url" field if the given value is not nil.
func (lc *LicenseCreate) SetNillableDetailsURL(s *string) *LicenseCreate {
	if s != nil {
		lc.SetDetailsURL(*s)
	}
	return lc
}

// SetIsOsiApproved sets the "is_osi_approved" field.
func (lc *LicenseCreate) SetIsOsiApproved(b bool) *LicenseCreate {
	lc.mutation.SetIsOsiApproved(b)
	return lc
}

// SetNillableIsOsiApproved sets the "is_osi_approved" field if the given value is not nil.
func (lc *LicenseCreate) SetNillableIsOsiApproved(b *bool) *LicenseCreate {
	if b != nil {
		lc.SetIsOsiApproved(*b)
	}
	return lc
}

// SetOwnerID sets the "owner" edge to the Organization entity by ID.
func (lc *LicenseCreate) SetOwnerID(id int) *LicenseCreate {
	lc.mutation.SetOwnerID(id)
	return lc
}

// SetOwner sets the "owner" edge to the Organization entity.
func (lc *LicenseCreate) SetOwner(o *Organization) *LicenseCreate {
	return lc.SetOwnerID(o.ID)
}

// AddComponentIDs adds the "components" edge to the Component entity by IDs.
func (lc *LicenseCreate) AddComponentIDs(ids ...int) *LicenseCreate {
	lc.mutation.AddComponentIDs(ids...)
	return lc
}

// AddComponents adds the "components" edges to the Component entity.
func (lc *LicenseCreate) AddComponents(c ...*Component) *LicenseCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return lc.AddComponentIDs(ids...)
}

// AddInstanceIDs adds the "instances" edge to the ReleaseLicense entity by IDs.
func (lc *LicenseCreate) AddInstanceIDs(ids ...int) *LicenseCreate {
	lc.mutation.AddInstanceIDs(ids...)
	return lc
}

// AddInstances adds the "instances" edges to the ReleaseLicense entity.
func (lc *LicenseCreate) AddInstances(r ...*ReleaseLicense) *LicenseCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return lc.AddInstanceIDs(ids...)
}

// Mutation returns the LicenseMutation object of the builder.
func (lc *LicenseCreate) Mutation() *LicenseMutation {
	return lc.mutation
}

// Save creates the License in the database.
func (lc *LicenseCreate) Save(ctx context.Context) (*License, error) {
	var (
		err  error
		node *License
	)
	lc.defaults()
	if len(lc.hooks) == 0 {
		if err = lc.check(); err != nil {
			return nil, err
		}
		node, err = lc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*LicenseMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = lc.check(); err != nil {
				return nil, err
			}
			lc.mutation = mutation
			if node, err = lc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(lc.hooks) - 1; i >= 0; i-- {
			if lc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = lc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, lc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (lc *LicenseCreate) SaveX(ctx context.Context) *License {
	v, err := lc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lc *LicenseCreate) Exec(ctx context.Context) error {
	_, err := lc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lc *LicenseCreate) ExecX(ctx context.Context) {
	if err := lc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (lc *LicenseCreate) defaults() {
	if _, ok := lc.mutation.IsOsiApproved(); !ok {
		v := license.DefaultIsOsiApproved
		lc.mutation.SetIsOsiApproved(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lc *LicenseCreate) check() error {
	if _, ok := lc.mutation.LicenseID(); !ok {
		return &ValidationError{Name: "license_id", err: errors.New(`ent: missing required field "license_id"`)}
	}
	if v, ok := lc.mutation.LicenseID(); ok {
		if err := license.LicenseIDValidator(v); err != nil {
			return &ValidationError{Name: "license_id", err: fmt.Errorf(`ent: validator failed for field "license_id": %w`, err)}
		}
	}
	if _, ok := lc.mutation.IsOsiApproved(); !ok {
		return &ValidationError{Name: "is_osi_approved", err: errors.New(`ent: missing required field "is_osi_approved"`)}
	}
	if _, ok := lc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner", err: errors.New("ent: missing required edge \"owner\"")}
	}
	return nil
}

func (lc *LicenseCreate) sqlSave(ctx context.Context) (*License, error) {
	_node, _spec := lc.createSpec()
	if err := sqlgraph.CreateNode(ctx, lc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (lc *LicenseCreate) createSpec() (*License, *sqlgraph.CreateSpec) {
	var (
		_node = &License{config: lc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: license.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: license.FieldID,
			},
		}
	)
	if value, ok := lc.mutation.LicenseID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldLicenseID,
		})
		_node.LicenseID = value
	}
	if value, ok := lc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldName,
		})
		_node.Name = value
	}
	if value, ok := lc.mutation.Reference(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldReference,
		})
		_node.Reference = value
	}
	if value, ok := lc.mutation.DetailsURL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: license.FieldDetailsURL,
		})
		_node.DetailsURL = value
	}
	if value, ok := lc.mutation.IsOsiApproved(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: license.FieldIsOsiApproved,
		})
		_node.IsOsiApproved = value
	}
	if nodes := lc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   license.OwnerTable,
			Columns: []string{license.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: organization.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.license_owner = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.ComponentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   license.ComponentsTable,
			Columns: license.ComponentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: component.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.InstancesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   license.InstancesTable,
			Columns: []string{license.InstancesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releaselicense.FieldID,
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

// LicenseCreateBulk is the builder for creating many License entities in bulk.
type LicenseCreateBulk struct {
	config
	builders []*LicenseCreate
}

// Save creates the License entities in the database.
func (lcb *LicenseCreateBulk) Save(ctx context.Context) ([]*License, error) {
	specs := make([]*sqlgraph.CreateSpec, len(lcb.builders))
	nodes := make([]*License, len(lcb.builders))
	mutators := make([]Mutator, len(lcb.builders))
	for i := range lcb.builders {
		func(i int, root context.Context) {
			builder := lcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*LicenseMutation)
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
					_, err = mutators[i+1].Mutate(root, lcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, lcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, lcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (lcb *LicenseCreateBulk) SaveX(ctx context.Context) []*License {
	v, err := lcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lcb *LicenseCreateBulk) Exec(ctx context.Context) error {
	_, err := lcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lcb *LicenseCreateBulk) ExecX(ctx context.Context) {
	if err := lcb.Exec(ctx); err != nil {
		panic(err)
	}
}
