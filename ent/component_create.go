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
	"github.com/valocode/bubbly/ent/releasecomponent"
	schema "github.com/valocode/bubbly/ent/schema/types"
	"github.com/valocode/bubbly/ent/vulnerability"
)

// ComponentCreate is the builder for creating a Component entity.
type ComponentCreate struct {
	config
	mutation *ComponentMutation
	hooks    []Hook
}

// SetScheme sets the "scheme" field.
func (cc *ComponentCreate) SetScheme(s string) *ComponentCreate {
	cc.mutation.SetScheme(s)
	return cc
}

// SetNamespace sets the "namespace" field.
func (cc *ComponentCreate) SetNamespace(s string) *ComponentCreate {
	cc.mutation.SetNamespace(s)
	return cc
}

// SetNillableNamespace sets the "namespace" field if the given value is not nil.
func (cc *ComponentCreate) SetNillableNamespace(s *string) *ComponentCreate {
	if s != nil {
		cc.SetNamespace(*s)
	}
	return cc
}

// SetName sets the "name" field.
func (cc *ComponentCreate) SetName(s string) *ComponentCreate {
	cc.mutation.SetName(s)
	return cc
}

// SetVersion sets the "version" field.
func (cc *ComponentCreate) SetVersion(s string) *ComponentCreate {
	cc.mutation.SetVersion(s)
	return cc
}

// SetDescription sets the "description" field.
func (cc *ComponentCreate) SetDescription(s string) *ComponentCreate {
	cc.mutation.SetDescription(s)
	return cc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cc *ComponentCreate) SetNillableDescription(s *string) *ComponentCreate {
	if s != nil {
		cc.SetDescription(*s)
	}
	return cc
}

// SetURL sets the "url" field.
func (cc *ComponentCreate) SetURL(s string) *ComponentCreate {
	cc.mutation.SetURL(s)
	return cc
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (cc *ComponentCreate) SetNillableURL(s *string) *ComponentCreate {
	if s != nil {
		cc.SetURL(*s)
	}
	return cc
}

// SetMetadata sets the "metadata" field.
func (cc *ComponentCreate) SetMetadata(s schema.Metadata) *ComponentCreate {
	cc.mutation.SetMetadata(s)
	return cc
}

// SetLabels sets the "labels" field.
func (cc *ComponentCreate) SetLabels(s schema.Labels) *ComponentCreate {
	cc.mutation.SetLabels(s)
	return cc
}

// SetOwnerID sets the "owner" edge to the Organization entity by ID.
func (cc *ComponentCreate) SetOwnerID(id int) *ComponentCreate {
	cc.mutation.SetOwnerID(id)
	return cc
}

// SetOwner sets the "owner" edge to the Organization entity.
func (cc *ComponentCreate) SetOwner(o *Organization) *ComponentCreate {
	return cc.SetOwnerID(o.ID)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (cc *ComponentCreate) AddVulnerabilityIDs(ids ...int) *ComponentCreate {
	cc.mutation.AddVulnerabilityIDs(ids...)
	return cc
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (cc *ComponentCreate) AddVulnerabilities(v ...*Vulnerability) *ComponentCreate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return cc.AddVulnerabilityIDs(ids...)
}

// AddLicenseIDs adds the "licenses" edge to the License entity by IDs.
func (cc *ComponentCreate) AddLicenseIDs(ids ...int) *ComponentCreate {
	cc.mutation.AddLicenseIDs(ids...)
	return cc
}

// AddLicenses adds the "licenses" edges to the License entity.
func (cc *ComponentCreate) AddLicenses(l ...*License) *ComponentCreate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return cc.AddLicenseIDs(ids...)
}

// AddUseIDs adds the "uses" edge to the ReleaseComponent entity by IDs.
func (cc *ComponentCreate) AddUseIDs(ids ...int) *ComponentCreate {
	cc.mutation.AddUseIDs(ids...)
	return cc
}

// AddUses adds the "uses" edges to the ReleaseComponent entity.
func (cc *ComponentCreate) AddUses(r ...*ReleaseComponent) *ComponentCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return cc.AddUseIDs(ids...)
}

// Mutation returns the ComponentMutation object of the builder.
func (cc *ComponentCreate) Mutation() *ComponentMutation {
	return cc.mutation
}

// Save creates the Component in the database.
func (cc *ComponentCreate) Save(ctx context.Context) (*Component, error) {
	var (
		err  error
		node *Component
	)
	cc.defaults()
	if len(cc.hooks) == 0 {
		if err = cc.check(); err != nil {
			return nil, err
		}
		node, err = cc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ComponentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = cc.check(); err != nil {
				return nil, err
			}
			cc.mutation = mutation
			if node, err = cc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(cc.hooks) - 1; i >= 0; i-- {
			if cc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = cc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (cc *ComponentCreate) SaveX(ctx context.Context) *Component {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *ComponentCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *ComponentCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *ComponentCreate) defaults() {
	if _, ok := cc.mutation.Namespace(); !ok {
		v := component.DefaultNamespace
		cc.mutation.SetNamespace(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *ComponentCreate) check() error {
	if _, ok := cc.mutation.Scheme(); !ok {
		return &ValidationError{Name: "scheme", err: errors.New(`ent: missing required field "scheme"`)}
	}
	if v, ok := cc.mutation.Scheme(); ok {
		if err := component.SchemeValidator(v); err != nil {
			return &ValidationError{Name: "scheme", err: fmt.Errorf(`ent: validator failed for field "scheme": %w`, err)}
		}
	}
	if _, ok := cc.mutation.Namespace(); !ok {
		return &ValidationError{Name: "namespace", err: errors.New(`ent: missing required field "namespace"`)}
	}
	if _, ok := cc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "name"`)}
	}
	if v, ok := cc.mutation.Name(); ok {
		if err := component.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "name": %w`, err)}
		}
	}
	if _, ok := cc.mutation.Version(); !ok {
		return &ValidationError{Name: "version", err: errors.New(`ent: missing required field "version"`)}
	}
	if v, ok := cc.mutation.Version(); ok {
		if err := component.VersionValidator(v); err != nil {
			return &ValidationError{Name: "version", err: fmt.Errorf(`ent: validator failed for field "version": %w`, err)}
		}
	}
	if _, ok := cc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner", err: errors.New("ent: missing required edge \"owner\"")}
	}
	return nil
}

func (cc *ComponentCreate) sqlSave(ctx context.Context) (*Component, error) {
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (cc *ComponentCreate) createSpec() (*Component, *sqlgraph.CreateSpec) {
	var (
		_node = &Component{config: cc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: component.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: component.FieldID,
			},
		}
	)
	if value, ok := cc.mutation.Scheme(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldScheme,
		})
		_node.Scheme = value
	}
	if value, ok := cc.mutation.Namespace(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldNamespace,
		})
		_node.Namespace = value
	}
	if value, ok := cc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldName,
		})
		_node.Name = value
	}
	if value, ok := cc.mutation.Version(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldVersion,
		})
		_node.Version = value
	}
	if value, ok := cc.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldDescription,
		})
		_node.Description = value
	}
	if value, ok := cc.mutation.URL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldURL,
		})
		_node.URL = value
	}
	if value, ok := cc.mutation.Metadata(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: component.FieldMetadata,
		})
		_node.Metadata = value
	}
	if value, ok := cc.mutation.Labels(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: component.FieldLabels,
		})
		_node.Labels = value
	}
	if nodes := cc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   component.OwnerTable,
			Columns: []string{component.OwnerColumn},
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
		_node.component_owner = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   component.VulnerabilitiesTable,
			Columns: component.VulnerabilitiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: vulnerability.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.LicensesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   component.LicensesTable,
			Columns: component.LicensesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: license.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.UsesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   component.UsesTable,
			Columns: []string{component.UsesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releasecomponent.FieldID,
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

// ComponentCreateBulk is the builder for creating many Component entities in bulk.
type ComponentCreateBulk struct {
	config
	builders []*ComponentCreate
}

// Save creates the Component entities in the database.
func (ccb *ComponentCreateBulk) Save(ctx context.Context) ([]*Component, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Component, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ComponentMutation)
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
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *ComponentCreateBulk) SaveX(ctx context.Context) []*Component {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *ComponentCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *ComponentCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
