// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/license"
	"github.com/valocode/bubbly/ent/organization"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/releasecomponent"
	schema "github.com/valocode/bubbly/ent/schema/types"
	"github.com/valocode/bubbly/ent/vulnerability"
)

// ComponentUpdate is the builder for updating Component entities.
type ComponentUpdate struct {
	config
	hooks    []Hook
	mutation *ComponentMutation
}

// Where appends a list predicates to the ComponentUpdate builder.
func (cu *ComponentUpdate) Where(ps ...predicate.Component) *ComponentUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetName sets the "name" field.
func (cu *ComponentUpdate) SetName(s string) *ComponentUpdate {
	cu.mutation.SetName(s)
	return cu
}

// SetVendor sets the "vendor" field.
func (cu *ComponentUpdate) SetVendor(s string) *ComponentUpdate {
	cu.mutation.SetVendor(s)
	return cu
}

// SetNillableVendor sets the "vendor" field if the given value is not nil.
func (cu *ComponentUpdate) SetNillableVendor(s *string) *ComponentUpdate {
	if s != nil {
		cu.SetVendor(*s)
	}
	return cu
}

// SetVersion sets the "version" field.
func (cu *ComponentUpdate) SetVersion(s string) *ComponentUpdate {
	cu.mutation.SetVersion(s)
	return cu
}

// SetDescription sets the "description" field.
func (cu *ComponentUpdate) SetDescription(s string) *ComponentUpdate {
	cu.mutation.SetDescription(s)
	return cu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cu *ComponentUpdate) SetNillableDescription(s *string) *ComponentUpdate {
	if s != nil {
		cu.SetDescription(*s)
	}
	return cu
}

// ClearDescription clears the value of the "description" field.
func (cu *ComponentUpdate) ClearDescription() *ComponentUpdate {
	cu.mutation.ClearDescription()
	return cu
}

// SetURL sets the "url" field.
func (cu *ComponentUpdate) SetURL(s string) *ComponentUpdate {
	cu.mutation.SetURL(s)
	return cu
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (cu *ComponentUpdate) SetNillableURL(s *string) *ComponentUpdate {
	if s != nil {
		cu.SetURL(*s)
	}
	return cu
}

// ClearURL clears the value of the "url" field.
func (cu *ComponentUpdate) ClearURL() *ComponentUpdate {
	cu.mutation.ClearURL()
	return cu
}

// SetMetadata sets the "metadata" field.
func (cu *ComponentUpdate) SetMetadata(s schema.Metadata) *ComponentUpdate {
	cu.mutation.SetMetadata(s)
	return cu
}

// ClearMetadata clears the value of the "metadata" field.
func (cu *ComponentUpdate) ClearMetadata() *ComponentUpdate {
	cu.mutation.ClearMetadata()
	return cu
}

// SetOwnerID sets the "owner" edge to the Organization entity by ID.
func (cu *ComponentUpdate) SetOwnerID(id int) *ComponentUpdate {
	cu.mutation.SetOwnerID(id)
	return cu
}

// SetOwner sets the "owner" edge to the Organization entity.
func (cu *ComponentUpdate) SetOwner(o *Organization) *ComponentUpdate {
	return cu.SetOwnerID(o.ID)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (cu *ComponentUpdate) AddVulnerabilityIDs(ids ...int) *ComponentUpdate {
	cu.mutation.AddVulnerabilityIDs(ids...)
	return cu
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (cu *ComponentUpdate) AddVulnerabilities(v ...*Vulnerability) *ComponentUpdate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return cu.AddVulnerabilityIDs(ids...)
}

// AddLicenseIDs adds the "licenses" edge to the License entity by IDs.
func (cu *ComponentUpdate) AddLicenseIDs(ids ...int) *ComponentUpdate {
	cu.mutation.AddLicenseIDs(ids...)
	return cu
}

// AddLicenses adds the "licenses" edges to the License entity.
func (cu *ComponentUpdate) AddLicenses(l ...*License) *ComponentUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return cu.AddLicenseIDs(ids...)
}

// AddUseIDs adds the "uses" edge to the ReleaseComponent entity by IDs.
func (cu *ComponentUpdate) AddUseIDs(ids ...int) *ComponentUpdate {
	cu.mutation.AddUseIDs(ids...)
	return cu
}

// AddUses adds the "uses" edges to the ReleaseComponent entity.
func (cu *ComponentUpdate) AddUses(r ...*ReleaseComponent) *ComponentUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return cu.AddUseIDs(ids...)
}

// Mutation returns the ComponentMutation object of the builder.
func (cu *ComponentUpdate) Mutation() *ComponentMutation {
	return cu.mutation
}

// ClearOwner clears the "owner" edge to the Organization entity.
func (cu *ComponentUpdate) ClearOwner() *ComponentUpdate {
	cu.mutation.ClearOwner()
	return cu
}

// ClearVulnerabilities clears all "vulnerabilities" edges to the Vulnerability entity.
func (cu *ComponentUpdate) ClearVulnerabilities() *ComponentUpdate {
	cu.mutation.ClearVulnerabilities()
	return cu
}

// RemoveVulnerabilityIDs removes the "vulnerabilities" edge to Vulnerability entities by IDs.
func (cu *ComponentUpdate) RemoveVulnerabilityIDs(ids ...int) *ComponentUpdate {
	cu.mutation.RemoveVulnerabilityIDs(ids...)
	return cu
}

// RemoveVulnerabilities removes "vulnerabilities" edges to Vulnerability entities.
func (cu *ComponentUpdate) RemoveVulnerabilities(v ...*Vulnerability) *ComponentUpdate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return cu.RemoveVulnerabilityIDs(ids...)
}

// ClearLicenses clears all "licenses" edges to the License entity.
func (cu *ComponentUpdate) ClearLicenses() *ComponentUpdate {
	cu.mutation.ClearLicenses()
	return cu
}

// RemoveLicenseIDs removes the "licenses" edge to License entities by IDs.
func (cu *ComponentUpdate) RemoveLicenseIDs(ids ...int) *ComponentUpdate {
	cu.mutation.RemoveLicenseIDs(ids...)
	return cu
}

// RemoveLicenses removes "licenses" edges to License entities.
func (cu *ComponentUpdate) RemoveLicenses(l ...*License) *ComponentUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return cu.RemoveLicenseIDs(ids...)
}

// ClearUses clears all "uses" edges to the ReleaseComponent entity.
func (cu *ComponentUpdate) ClearUses() *ComponentUpdate {
	cu.mutation.ClearUses()
	return cu
}

// RemoveUseIDs removes the "uses" edge to ReleaseComponent entities by IDs.
func (cu *ComponentUpdate) RemoveUseIDs(ids ...int) *ComponentUpdate {
	cu.mutation.RemoveUseIDs(ids...)
	return cu
}

// RemoveUses removes "uses" edges to ReleaseComponent entities.
func (cu *ComponentUpdate) RemoveUses(r ...*ReleaseComponent) *ComponentUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return cu.RemoveUseIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *ComponentUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(cu.hooks) == 0 {
		if err = cu.check(); err != nil {
			return 0, err
		}
		affected, err = cu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ComponentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = cu.check(); err != nil {
				return 0, err
			}
			cu.mutation = mutation
			affected, err = cu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(cu.hooks) - 1; i >= 0; i-- {
			if cu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = cu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (cu *ComponentUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *ComponentUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *ComponentUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cu *ComponentUpdate) check() error {
	if v, ok := cu.mutation.Name(); ok {
		if err := component.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	if v, ok := cu.mutation.Version(); ok {
		if err := component.VersionValidator(v); err != nil {
			return &ValidationError{Name: "version", err: fmt.Errorf("ent: validator failed for field \"version\": %w", err)}
		}
	}
	if _, ok := cu.mutation.OwnerID(); cu.mutation.OwnerCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"owner\"")
	}
	return nil
}

func (cu *ComponentUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   component.Table,
			Columns: component.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: component.FieldID,
			},
		},
	}
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldName,
		})
	}
	if value, ok := cu.mutation.Vendor(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldVendor,
		})
	}
	if value, ok := cu.mutation.Version(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldVersion,
		})
	}
	if value, ok := cu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldDescription,
		})
	}
	if cu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: component.FieldDescription,
		})
	}
	if value, ok := cu.mutation.URL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldURL,
		})
	}
	if cu.mutation.URLCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: component.FieldURL,
		})
	}
	if value, ok := cu.mutation.Metadata(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: component.FieldMetadata,
		})
	}
	if cu.mutation.MetadataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: component.FieldMetadata,
		})
	}
	if cu.mutation.OwnerCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.OwnerIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cu.mutation.VulnerabilitiesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedVulnerabilitiesIDs(); len(nodes) > 0 && !cu.mutation.VulnerabilitiesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cu.mutation.LicensesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedLicensesIDs(); len(nodes) > 0 && !cu.mutation.LicensesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.LicensesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cu.mutation.UsesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedUsesIDs(); len(nodes) > 0 && !cu.mutation.UsesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.UsesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{component.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// ComponentUpdateOne is the builder for updating a single Component entity.
type ComponentUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ComponentMutation
}

// SetName sets the "name" field.
func (cuo *ComponentUpdateOne) SetName(s string) *ComponentUpdateOne {
	cuo.mutation.SetName(s)
	return cuo
}

// SetVendor sets the "vendor" field.
func (cuo *ComponentUpdateOne) SetVendor(s string) *ComponentUpdateOne {
	cuo.mutation.SetVendor(s)
	return cuo
}

// SetNillableVendor sets the "vendor" field if the given value is not nil.
func (cuo *ComponentUpdateOne) SetNillableVendor(s *string) *ComponentUpdateOne {
	if s != nil {
		cuo.SetVendor(*s)
	}
	return cuo
}

// SetVersion sets the "version" field.
func (cuo *ComponentUpdateOne) SetVersion(s string) *ComponentUpdateOne {
	cuo.mutation.SetVersion(s)
	return cuo
}

// SetDescription sets the "description" field.
func (cuo *ComponentUpdateOne) SetDescription(s string) *ComponentUpdateOne {
	cuo.mutation.SetDescription(s)
	return cuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cuo *ComponentUpdateOne) SetNillableDescription(s *string) *ComponentUpdateOne {
	if s != nil {
		cuo.SetDescription(*s)
	}
	return cuo
}

// ClearDescription clears the value of the "description" field.
func (cuo *ComponentUpdateOne) ClearDescription() *ComponentUpdateOne {
	cuo.mutation.ClearDescription()
	return cuo
}

// SetURL sets the "url" field.
func (cuo *ComponentUpdateOne) SetURL(s string) *ComponentUpdateOne {
	cuo.mutation.SetURL(s)
	return cuo
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (cuo *ComponentUpdateOne) SetNillableURL(s *string) *ComponentUpdateOne {
	if s != nil {
		cuo.SetURL(*s)
	}
	return cuo
}

// ClearURL clears the value of the "url" field.
func (cuo *ComponentUpdateOne) ClearURL() *ComponentUpdateOne {
	cuo.mutation.ClearURL()
	return cuo
}

// SetMetadata sets the "metadata" field.
func (cuo *ComponentUpdateOne) SetMetadata(s schema.Metadata) *ComponentUpdateOne {
	cuo.mutation.SetMetadata(s)
	return cuo
}

// ClearMetadata clears the value of the "metadata" field.
func (cuo *ComponentUpdateOne) ClearMetadata() *ComponentUpdateOne {
	cuo.mutation.ClearMetadata()
	return cuo
}

// SetOwnerID sets the "owner" edge to the Organization entity by ID.
func (cuo *ComponentUpdateOne) SetOwnerID(id int) *ComponentUpdateOne {
	cuo.mutation.SetOwnerID(id)
	return cuo
}

// SetOwner sets the "owner" edge to the Organization entity.
func (cuo *ComponentUpdateOne) SetOwner(o *Organization) *ComponentUpdateOne {
	return cuo.SetOwnerID(o.ID)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (cuo *ComponentUpdateOne) AddVulnerabilityIDs(ids ...int) *ComponentUpdateOne {
	cuo.mutation.AddVulnerabilityIDs(ids...)
	return cuo
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (cuo *ComponentUpdateOne) AddVulnerabilities(v ...*Vulnerability) *ComponentUpdateOne {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return cuo.AddVulnerabilityIDs(ids...)
}

// AddLicenseIDs adds the "licenses" edge to the License entity by IDs.
func (cuo *ComponentUpdateOne) AddLicenseIDs(ids ...int) *ComponentUpdateOne {
	cuo.mutation.AddLicenseIDs(ids...)
	return cuo
}

// AddLicenses adds the "licenses" edges to the License entity.
func (cuo *ComponentUpdateOne) AddLicenses(l ...*License) *ComponentUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return cuo.AddLicenseIDs(ids...)
}

// AddUseIDs adds the "uses" edge to the ReleaseComponent entity by IDs.
func (cuo *ComponentUpdateOne) AddUseIDs(ids ...int) *ComponentUpdateOne {
	cuo.mutation.AddUseIDs(ids...)
	return cuo
}

// AddUses adds the "uses" edges to the ReleaseComponent entity.
func (cuo *ComponentUpdateOne) AddUses(r ...*ReleaseComponent) *ComponentUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return cuo.AddUseIDs(ids...)
}

// Mutation returns the ComponentMutation object of the builder.
func (cuo *ComponentUpdateOne) Mutation() *ComponentMutation {
	return cuo.mutation
}

// ClearOwner clears the "owner" edge to the Organization entity.
func (cuo *ComponentUpdateOne) ClearOwner() *ComponentUpdateOne {
	cuo.mutation.ClearOwner()
	return cuo
}

// ClearVulnerabilities clears all "vulnerabilities" edges to the Vulnerability entity.
func (cuo *ComponentUpdateOne) ClearVulnerabilities() *ComponentUpdateOne {
	cuo.mutation.ClearVulnerabilities()
	return cuo
}

// RemoveVulnerabilityIDs removes the "vulnerabilities" edge to Vulnerability entities by IDs.
func (cuo *ComponentUpdateOne) RemoveVulnerabilityIDs(ids ...int) *ComponentUpdateOne {
	cuo.mutation.RemoveVulnerabilityIDs(ids...)
	return cuo
}

// RemoveVulnerabilities removes "vulnerabilities" edges to Vulnerability entities.
func (cuo *ComponentUpdateOne) RemoveVulnerabilities(v ...*Vulnerability) *ComponentUpdateOne {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return cuo.RemoveVulnerabilityIDs(ids...)
}

// ClearLicenses clears all "licenses" edges to the License entity.
func (cuo *ComponentUpdateOne) ClearLicenses() *ComponentUpdateOne {
	cuo.mutation.ClearLicenses()
	return cuo
}

// RemoveLicenseIDs removes the "licenses" edge to License entities by IDs.
func (cuo *ComponentUpdateOne) RemoveLicenseIDs(ids ...int) *ComponentUpdateOne {
	cuo.mutation.RemoveLicenseIDs(ids...)
	return cuo
}

// RemoveLicenses removes "licenses" edges to License entities.
func (cuo *ComponentUpdateOne) RemoveLicenses(l ...*License) *ComponentUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return cuo.RemoveLicenseIDs(ids...)
}

// ClearUses clears all "uses" edges to the ReleaseComponent entity.
func (cuo *ComponentUpdateOne) ClearUses() *ComponentUpdateOne {
	cuo.mutation.ClearUses()
	return cuo
}

// RemoveUseIDs removes the "uses" edge to ReleaseComponent entities by IDs.
func (cuo *ComponentUpdateOne) RemoveUseIDs(ids ...int) *ComponentUpdateOne {
	cuo.mutation.RemoveUseIDs(ids...)
	return cuo
}

// RemoveUses removes "uses" edges to ReleaseComponent entities.
func (cuo *ComponentUpdateOne) RemoveUses(r ...*ReleaseComponent) *ComponentUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return cuo.RemoveUseIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *ComponentUpdateOne) Select(field string, fields ...string) *ComponentUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Component entity.
func (cuo *ComponentUpdateOne) Save(ctx context.Context) (*Component, error) {
	var (
		err  error
		node *Component
	)
	if len(cuo.hooks) == 0 {
		if err = cuo.check(); err != nil {
			return nil, err
		}
		node, err = cuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ComponentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = cuo.check(); err != nil {
				return nil, err
			}
			cuo.mutation = mutation
			node, err = cuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(cuo.hooks) - 1; i >= 0; i-- {
			if cuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = cuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *ComponentUpdateOne) SaveX(ctx context.Context) *Component {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *ComponentUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *ComponentUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cuo *ComponentUpdateOne) check() error {
	if v, ok := cuo.mutation.Name(); ok {
		if err := component.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	if v, ok := cuo.mutation.Version(); ok {
		if err := component.VersionValidator(v); err != nil {
			return &ValidationError{Name: "version", err: fmt.Errorf("ent: validator failed for field \"version\": %w", err)}
		}
	}
	if _, ok := cuo.mutation.OwnerID(); cuo.mutation.OwnerCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"owner\"")
	}
	return nil
}

func (cuo *ComponentUpdateOne) sqlSave(ctx context.Context) (_node *Component, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   component.Table,
			Columns: component.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: component.FieldID,
			},
		},
	}
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Component.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, component.FieldID)
		for _, f := range fields {
			if !component.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != component.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldName,
		})
	}
	if value, ok := cuo.mutation.Vendor(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldVendor,
		})
	}
	if value, ok := cuo.mutation.Version(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldVersion,
		})
	}
	if value, ok := cuo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldDescription,
		})
	}
	if cuo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: component.FieldDescription,
		})
	}
	if value, ok := cuo.mutation.URL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: component.FieldURL,
		})
	}
	if cuo.mutation.URLCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: component.FieldURL,
		})
	}
	if value, ok := cuo.mutation.Metadata(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: component.FieldMetadata,
		})
	}
	if cuo.mutation.MetadataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: component.FieldMetadata,
		})
	}
	if cuo.mutation.OwnerCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.OwnerIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cuo.mutation.VulnerabilitiesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedVulnerabilitiesIDs(); len(nodes) > 0 && !cuo.mutation.VulnerabilitiesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cuo.mutation.LicensesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedLicensesIDs(); len(nodes) > 0 && !cuo.mutation.LicensesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.LicensesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cuo.mutation.UsesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedUsesIDs(); len(nodes) > 0 && !cuo.mutation.UsesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.UsesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Component{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{component.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
