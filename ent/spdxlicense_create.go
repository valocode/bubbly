// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/spdxlicense"
)

// SPDXLicenseCreate is the builder for creating a SPDXLicense entity.
type SPDXLicenseCreate struct {
	config
	mutation *SPDXLicenseMutation
	hooks    []Hook
}

// SetLicenseID sets the "license_id" field.
func (slc *SPDXLicenseCreate) SetLicenseID(s string) *SPDXLicenseCreate {
	slc.mutation.SetLicenseID(s)
	return slc
}

// SetName sets the "name" field.
func (slc *SPDXLicenseCreate) SetName(s string) *SPDXLicenseCreate {
	slc.mutation.SetName(s)
	return slc
}

// SetReference sets the "reference" field.
func (slc *SPDXLicenseCreate) SetReference(s string) *SPDXLicenseCreate {
	slc.mutation.SetReference(s)
	return slc
}

// SetNillableReference sets the "reference" field if the given value is not nil.
func (slc *SPDXLicenseCreate) SetNillableReference(s *string) *SPDXLicenseCreate {
	if s != nil {
		slc.SetReference(*s)
	}
	return slc
}

// SetDetailsURL sets the "details_url" field.
func (slc *SPDXLicenseCreate) SetDetailsURL(s string) *SPDXLicenseCreate {
	slc.mutation.SetDetailsURL(s)
	return slc
}

// SetNillableDetailsURL sets the "details_url" field if the given value is not nil.
func (slc *SPDXLicenseCreate) SetNillableDetailsURL(s *string) *SPDXLicenseCreate {
	if s != nil {
		slc.SetDetailsURL(*s)
	}
	return slc
}

// SetIsOsiApproved sets the "is_osi_approved" field.
func (slc *SPDXLicenseCreate) SetIsOsiApproved(b bool) *SPDXLicenseCreate {
	slc.mutation.SetIsOsiApproved(b)
	return slc
}

// SetNillableIsOsiApproved sets the "is_osi_approved" field if the given value is not nil.
func (slc *SPDXLicenseCreate) SetNillableIsOsiApproved(b *bool) *SPDXLicenseCreate {
	if b != nil {
		slc.SetIsOsiApproved(*b)
	}
	return slc
}

// Mutation returns the SPDXLicenseMutation object of the builder.
func (slc *SPDXLicenseCreate) Mutation() *SPDXLicenseMutation {
	return slc.mutation
}

// Save creates the SPDXLicense in the database.
func (slc *SPDXLicenseCreate) Save(ctx context.Context) (*SPDXLicense, error) {
	var (
		err  error
		node *SPDXLicense
	)
	if err := slc.defaults(); err != nil {
		return nil, err
	}
	if len(slc.hooks) == 0 {
		if err = slc.check(); err != nil {
			return nil, err
		}
		node, err = slc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SPDXLicenseMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = slc.check(); err != nil {
				return nil, err
			}
			slc.mutation = mutation
			if node, err = slc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(slc.hooks) - 1; i >= 0; i-- {
			if slc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = slc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, slc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (slc *SPDXLicenseCreate) SaveX(ctx context.Context) *SPDXLicense {
	v, err := slc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (slc *SPDXLicenseCreate) Exec(ctx context.Context) error {
	_, err := slc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (slc *SPDXLicenseCreate) ExecX(ctx context.Context) {
	if err := slc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (slc *SPDXLicenseCreate) defaults() error {
	if _, ok := slc.mutation.IsOsiApproved(); !ok {
		v := spdxlicense.DefaultIsOsiApproved
		slc.mutation.SetIsOsiApproved(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (slc *SPDXLicenseCreate) check() error {
	if _, ok := slc.mutation.LicenseID(); !ok {
		return &ValidationError{Name: "license_id", err: errors.New(`ent: missing required field "license_id"`)}
	}
	if v, ok := slc.mutation.LicenseID(); ok {
		if err := spdxlicense.LicenseIDValidator(v); err != nil {
			return &ValidationError{Name: "license_id", err: fmt.Errorf(`ent: validator failed for field "license_id": %w`, err)}
		}
	}
	if _, ok := slc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "name"`)}
	}
	if v, ok := slc.mutation.Name(); ok {
		if err := spdxlicense.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "name": %w`, err)}
		}
	}
	if _, ok := slc.mutation.IsOsiApproved(); !ok {
		return &ValidationError{Name: "is_osi_approved", err: errors.New(`ent: missing required field "is_osi_approved"`)}
	}
	return nil
}

func (slc *SPDXLicenseCreate) sqlSave(ctx context.Context) (*SPDXLicense, error) {
	_node, _spec := slc.createSpec()
	if err := sqlgraph.CreateNode(ctx, slc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (slc *SPDXLicenseCreate) createSpec() (*SPDXLicense, *sqlgraph.CreateSpec) {
	var (
		_node = &SPDXLicense{config: slc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: spdxlicense.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: spdxlicense.FieldID,
			},
		}
	)
	if value, ok := slc.mutation.LicenseID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: spdxlicense.FieldLicenseID,
		})
		_node.LicenseID = value
	}
	if value, ok := slc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: spdxlicense.FieldName,
		})
		_node.Name = value
	}
	if value, ok := slc.mutation.Reference(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: spdxlicense.FieldReference,
		})
		_node.Reference = value
	}
	if value, ok := slc.mutation.DetailsURL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: spdxlicense.FieldDetailsURL,
		})
		_node.DetailsURL = value
	}
	if value, ok := slc.mutation.IsOsiApproved(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: spdxlicense.FieldIsOsiApproved,
		})
		_node.IsOsiApproved = value
	}
	return _node, _spec
}

// SPDXLicenseCreateBulk is the builder for creating many SPDXLicense entities in bulk.
type SPDXLicenseCreateBulk struct {
	config
	builders []*SPDXLicenseCreate
}

// Save creates the SPDXLicense entities in the database.
func (slcb *SPDXLicenseCreateBulk) Save(ctx context.Context) ([]*SPDXLicense, error) {
	specs := make([]*sqlgraph.CreateSpec, len(slcb.builders))
	nodes := make([]*SPDXLicense, len(slcb.builders))
	mutators := make([]Mutator, len(slcb.builders))
	for i := range slcb.builders {
		func(i int, root context.Context) {
			builder := slcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SPDXLicenseMutation)
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
					_, err = mutators[i+1].Mutate(root, slcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, slcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, slcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (slcb *SPDXLicenseCreateBulk) SaveX(ctx context.Context) []*SPDXLicense {
	v, err := slcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (slcb *SPDXLicenseCreateBulk) Exec(ctx context.Context) error {
	_, err := slcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (slcb *SPDXLicenseCreateBulk) ExecX(ctx context.Context) {
	if err := slcb.Exec(ctx); err != nil {
		panic(err)
	}
}
