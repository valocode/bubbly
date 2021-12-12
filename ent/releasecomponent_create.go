// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/codescan"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasecomponent"
	"github.com/valocode/bubbly/ent/releaselicense"
	"github.com/valocode/bubbly/ent/releasevulnerability"
)

// ReleaseComponentCreate is the builder for creating a ReleaseComponent entity.
type ReleaseComponentCreate struct {
	config
	mutation *ReleaseComponentMutation
	hooks    []Hook
}

// SetType sets the "type" field.
func (rcc *ReleaseComponentCreate) SetType(r releasecomponent.Type) *ReleaseComponentCreate {
	rcc.mutation.SetType(r)
	return rcc
}

// SetNillableType sets the "type" field if the given value is not nil.
func (rcc *ReleaseComponentCreate) SetNillableType(r *releasecomponent.Type) *ReleaseComponentCreate {
	if r != nil {
		rcc.SetType(*r)
	}
	return rcc
}

// SetReleaseID sets the "release" edge to the Release entity by ID.
func (rcc *ReleaseComponentCreate) SetReleaseID(id int) *ReleaseComponentCreate {
	rcc.mutation.SetReleaseID(id)
	return rcc
}

// SetRelease sets the "release" edge to the Release entity.
func (rcc *ReleaseComponentCreate) SetRelease(r *Release) *ReleaseComponentCreate {
	return rcc.SetReleaseID(r.ID)
}

// AddScanIDs adds the "scans" edge to the CodeScan entity by IDs.
func (rcc *ReleaseComponentCreate) AddScanIDs(ids ...int) *ReleaseComponentCreate {
	rcc.mutation.AddScanIDs(ids...)
	return rcc
}

// AddScans adds the "scans" edges to the CodeScan entity.
func (rcc *ReleaseComponentCreate) AddScans(c ...*CodeScan) *ReleaseComponentCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return rcc.AddScanIDs(ids...)
}

// SetComponentID sets the "component" edge to the Component entity by ID.
func (rcc *ReleaseComponentCreate) SetComponentID(id int) *ReleaseComponentCreate {
	rcc.mutation.SetComponentID(id)
	return rcc
}

// SetComponent sets the "component" edge to the Component entity.
func (rcc *ReleaseComponentCreate) SetComponent(c *Component) *ReleaseComponentCreate {
	return rcc.SetComponentID(c.ID)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the ReleaseVulnerability entity by IDs.
func (rcc *ReleaseComponentCreate) AddVulnerabilityIDs(ids ...int) *ReleaseComponentCreate {
	rcc.mutation.AddVulnerabilityIDs(ids...)
	return rcc
}

// AddVulnerabilities adds the "vulnerabilities" edges to the ReleaseVulnerability entity.
func (rcc *ReleaseComponentCreate) AddVulnerabilities(r ...*ReleaseVulnerability) *ReleaseComponentCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rcc.AddVulnerabilityIDs(ids...)
}

// AddLicenseIDs adds the "licenses" edge to the ReleaseLicense entity by IDs.
func (rcc *ReleaseComponentCreate) AddLicenseIDs(ids ...int) *ReleaseComponentCreate {
	rcc.mutation.AddLicenseIDs(ids...)
	return rcc
}

// AddLicenses adds the "licenses" edges to the ReleaseLicense entity.
func (rcc *ReleaseComponentCreate) AddLicenses(r ...*ReleaseLicense) *ReleaseComponentCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rcc.AddLicenseIDs(ids...)
}

// Mutation returns the ReleaseComponentMutation object of the builder.
func (rcc *ReleaseComponentCreate) Mutation() *ReleaseComponentMutation {
	return rcc.mutation
}

// Save creates the ReleaseComponent in the database.
func (rcc *ReleaseComponentCreate) Save(ctx context.Context) (*ReleaseComponent, error) {
	var (
		err  error
		node *ReleaseComponent
	)
	rcc.defaults()
	if len(rcc.hooks) == 0 {
		if err = rcc.check(); err != nil {
			return nil, err
		}
		node, err = rcc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ReleaseComponentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rcc.check(); err != nil {
				return nil, err
			}
			rcc.mutation = mutation
			if node, err = rcc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(rcc.hooks) - 1; i >= 0; i-- {
			if rcc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = rcc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rcc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rcc *ReleaseComponentCreate) SaveX(ctx context.Context) *ReleaseComponent {
	v, err := rcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rcc *ReleaseComponentCreate) Exec(ctx context.Context) error {
	_, err := rcc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcc *ReleaseComponentCreate) ExecX(ctx context.Context) {
	if err := rcc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rcc *ReleaseComponentCreate) defaults() {
	if _, ok := rcc.mutation.GetType(); !ok {
		v := releasecomponent.DefaultType
		rcc.mutation.SetType(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rcc *ReleaseComponentCreate) check() error {
	if _, ok := rcc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "type"`)}
	}
	if v, ok := rcc.mutation.GetType(); ok {
		if err := releasecomponent.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "type": %w`, err)}
		}
	}
	if _, ok := rcc.mutation.ReleaseID(); !ok {
		return &ValidationError{Name: "release", err: errors.New("ent: missing required edge \"release\"")}
	}
	if len(rcc.mutation.ScansIDs()) == 0 {
		return &ValidationError{Name: "scans", err: errors.New("ent: missing required edge \"scans\"")}
	}
	if _, ok := rcc.mutation.ComponentID(); !ok {
		return &ValidationError{Name: "component", err: errors.New("ent: missing required edge \"component\"")}
	}
	return nil
}

func (rcc *ReleaseComponentCreate) sqlSave(ctx context.Context) (*ReleaseComponent, error) {
	_node, _spec := rcc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rcc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (rcc *ReleaseComponentCreate) createSpec() (*ReleaseComponent, *sqlgraph.CreateSpec) {
	var (
		_node = &ReleaseComponent{config: rcc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: releasecomponent.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: releasecomponent.FieldID,
			},
		}
	)
	if value, ok := rcc.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: releasecomponent.FieldType,
		})
		_node.Type = value
	}
	if nodes := rcc.mutation.ReleaseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   releasecomponent.ReleaseTable,
			Columns: []string{releasecomponent.ReleaseColumn},
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
		_node.release_component_release = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rcc.mutation.ScansIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   releasecomponent.ScansTable,
			Columns: releasecomponent.ScansPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: codescan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rcc.mutation.ComponentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   releasecomponent.ComponentTable,
			Columns: []string{releasecomponent.ComponentColumn},
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
		_node.release_component_component = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rcc.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   releasecomponent.VulnerabilitiesTable,
			Columns: []string{releasecomponent.VulnerabilitiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releasevulnerability.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rcc.mutation.LicensesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   releasecomponent.LicensesTable,
			Columns: []string{releasecomponent.LicensesColumn},
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

// ReleaseComponentCreateBulk is the builder for creating many ReleaseComponent entities in bulk.
type ReleaseComponentCreateBulk struct {
	config
	builders []*ReleaseComponentCreate
}

// Save creates the ReleaseComponent entities in the database.
func (rccb *ReleaseComponentCreateBulk) Save(ctx context.Context) ([]*ReleaseComponent, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rccb.builders))
	nodes := make([]*ReleaseComponent, len(rccb.builders))
	mutators := make([]Mutator, len(rccb.builders))
	for i := range rccb.builders {
		func(i int, root context.Context) {
			builder := rccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ReleaseComponentMutation)
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
					_, err = mutators[i+1].Mutate(root, rccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, rccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rccb *ReleaseComponentCreateBulk) SaveX(ctx context.Context) []*ReleaseComponent {
	v, err := rccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rccb *ReleaseComponentCreateBulk) Exec(ctx context.Context) error {
	_, err := rccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rccb *ReleaseComponentCreateBulk) ExecX(ctx context.Context) {
	if err := rccb.Exec(ctx); err != nil {
		panic(err)
	}
}
