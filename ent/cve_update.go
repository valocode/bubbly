// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/cve"
	"github.com/valocode/bubbly/ent/cverule"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/vulnerability"
)

// CVEUpdate is the builder for updating CVE entities.
type CVEUpdate struct {
	config
	hooks    []Hook
	mutation *CVEMutation
}

// Where appends a list predicates to the CVEUpdate builder.
func (cu *CVEUpdate) Where(ps ...predicate.CVE) *CVEUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetCveID sets the "cve_id" field.
func (cu *CVEUpdate) SetCveID(s string) *CVEUpdate {
	cu.mutation.SetCveID(s)
	return cu
}

// SetDescription sets the "description" field.
func (cu *CVEUpdate) SetDescription(s string) *CVEUpdate {
	cu.mutation.SetDescription(s)
	return cu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cu *CVEUpdate) SetNillableDescription(s *string) *CVEUpdate {
	if s != nil {
		cu.SetDescription(*s)
	}
	return cu
}

// ClearDescription clears the value of the "description" field.
func (cu *CVEUpdate) ClearDescription() *CVEUpdate {
	cu.mutation.ClearDescription()
	return cu
}

// SetSeverityScore sets the "severity_score" field.
func (cu *CVEUpdate) SetSeverityScore(f float64) *CVEUpdate {
	cu.mutation.ResetSeverityScore()
	cu.mutation.SetSeverityScore(f)
	return cu
}

// SetNillableSeverityScore sets the "severity_score" field if the given value is not nil.
func (cu *CVEUpdate) SetNillableSeverityScore(f *float64) *CVEUpdate {
	if f != nil {
		cu.SetSeverityScore(*f)
	}
	return cu
}

// AddSeverityScore adds f to the "severity_score" field.
func (cu *CVEUpdate) AddSeverityScore(f float64) *CVEUpdate {
	cu.mutation.AddSeverityScore(f)
	return cu
}

// SetSeverity sets the "severity" field.
func (cu *CVEUpdate) SetSeverity(c cve.Severity) *CVEUpdate {
	cu.mutation.SetSeverity(c)
	return cu
}

// SetNillableSeverity sets the "severity" field if the given value is not nil.
func (cu *CVEUpdate) SetNillableSeverity(c *cve.Severity) *CVEUpdate {
	if c != nil {
		cu.SetSeverity(*c)
	}
	return cu
}

// SetPublishedData sets the "published_data" field.
func (cu *CVEUpdate) SetPublishedData(t time.Time) *CVEUpdate {
	cu.mutation.SetPublishedData(t)
	return cu
}

// SetNillablePublishedData sets the "published_data" field if the given value is not nil.
func (cu *CVEUpdate) SetNillablePublishedData(t *time.Time) *CVEUpdate {
	if t != nil {
		cu.SetPublishedData(*t)
	}
	return cu
}

// ClearPublishedData clears the value of the "published_data" field.
func (cu *CVEUpdate) ClearPublishedData() *CVEUpdate {
	cu.mutation.ClearPublishedData()
	return cu
}

// SetModifiedData sets the "modified_data" field.
func (cu *CVEUpdate) SetModifiedData(t time.Time) *CVEUpdate {
	cu.mutation.SetModifiedData(t)
	return cu
}

// SetNillableModifiedData sets the "modified_data" field if the given value is not nil.
func (cu *CVEUpdate) SetNillableModifiedData(t *time.Time) *CVEUpdate {
	if t != nil {
		cu.SetModifiedData(*t)
	}
	return cu
}

// ClearModifiedData clears the value of the "modified_data" field.
func (cu *CVEUpdate) ClearModifiedData() *CVEUpdate {
	cu.mutation.ClearModifiedData()
	return cu
}

// AddComponentIDs adds the "components" edge to the Component entity by IDs.
func (cu *CVEUpdate) AddComponentIDs(ids ...int) *CVEUpdate {
	cu.mutation.AddComponentIDs(ids...)
	return cu
}

// AddComponents adds the "components" edges to the Component entity.
func (cu *CVEUpdate) AddComponents(c ...*Component) *CVEUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cu.AddComponentIDs(ids...)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (cu *CVEUpdate) AddVulnerabilityIDs(ids ...int) *CVEUpdate {
	cu.mutation.AddVulnerabilityIDs(ids...)
	return cu
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (cu *CVEUpdate) AddVulnerabilities(v ...*Vulnerability) *CVEUpdate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return cu.AddVulnerabilityIDs(ids...)
}

// AddRuleIDs adds the "rules" edge to the CVERule entity by IDs.
func (cu *CVEUpdate) AddRuleIDs(ids ...int) *CVEUpdate {
	cu.mutation.AddRuleIDs(ids...)
	return cu
}

// AddRules adds the "rules" edges to the CVERule entity.
func (cu *CVEUpdate) AddRules(c ...*CVERule) *CVEUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cu.AddRuleIDs(ids...)
}

// Mutation returns the CVEMutation object of the builder.
func (cu *CVEUpdate) Mutation() *CVEMutation {
	return cu.mutation
}

// ClearComponents clears all "components" edges to the Component entity.
func (cu *CVEUpdate) ClearComponents() *CVEUpdate {
	cu.mutation.ClearComponents()
	return cu
}

// RemoveComponentIDs removes the "components" edge to Component entities by IDs.
func (cu *CVEUpdate) RemoveComponentIDs(ids ...int) *CVEUpdate {
	cu.mutation.RemoveComponentIDs(ids...)
	return cu
}

// RemoveComponents removes "components" edges to Component entities.
func (cu *CVEUpdate) RemoveComponents(c ...*Component) *CVEUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cu.RemoveComponentIDs(ids...)
}

// ClearVulnerabilities clears all "vulnerabilities" edges to the Vulnerability entity.
func (cu *CVEUpdate) ClearVulnerabilities() *CVEUpdate {
	cu.mutation.ClearVulnerabilities()
	return cu
}

// RemoveVulnerabilityIDs removes the "vulnerabilities" edge to Vulnerability entities by IDs.
func (cu *CVEUpdate) RemoveVulnerabilityIDs(ids ...int) *CVEUpdate {
	cu.mutation.RemoveVulnerabilityIDs(ids...)
	return cu
}

// RemoveVulnerabilities removes "vulnerabilities" edges to Vulnerability entities.
func (cu *CVEUpdate) RemoveVulnerabilities(v ...*Vulnerability) *CVEUpdate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return cu.RemoveVulnerabilityIDs(ids...)
}

// ClearRules clears all "rules" edges to the CVERule entity.
func (cu *CVEUpdate) ClearRules() *CVEUpdate {
	cu.mutation.ClearRules()
	return cu
}

// RemoveRuleIDs removes the "rules" edge to CVERule entities by IDs.
func (cu *CVEUpdate) RemoveRuleIDs(ids ...int) *CVEUpdate {
	cu.mutation.RemoveRuleIDs(ids...)
	return cu
}

// RemoveRules removes "rules" edges to CVERule entities.
func (cu *CVEUpdate) RemoveRules(c ...*CVERule) *CVEUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cu.RemoveRuleIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *CVEUpdate) Save(ctx context.Context) (int, error) {
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
			mutation, ok := m.(*CVEMutation)
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
func (cu *CVEUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CVEUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CVEUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cu *CVEUpdate) check() error {
	if v, ok := cu.mutation.CveID(); ok {
		if err := cve.CveIDValidator(v); err != nil {
			return &ValidationError{Name: "cve_id", err: fmt.Errorf("ent: validator failed for field \"cve_id\": %w", err)}
		}
	}
	if v, ok := cu.mutation.Severity(); ok {
		if err := cve.SeverityValidator(v); err != nil {
			return &ValidationError{Name: "severity", err: fmt.Errorf("ent: validator failed for field \"severity\": %w", err)}
		}
	}
	return nil
}

func (cu *CVEUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cve.Table,
			Columns: cve.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: cve.FieldID,
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
	if value, ok := cu.mutation.CveID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cve.FieldCveID,
		})
	}
	if value, ok := cu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cve.FieldDescription,
		})
	}
	if cu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: cve.FieldDescription,
		})
	}
	if value, ok := cu.mutation.SeverityScore(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: cve.FieldSeverityScore,
		})
	}
	if value, ok := cu.mutation.AddedSeverityScore(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: cve.FieldSeverityScore,
		})
	}
	if value, ok := cu.mutation.Severity(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: cve.FieldSeverity,
		})
	}
	if value, ok := cu.mutation.PublishedData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cve.FieldPublishedData,
		})
	}
	if cu.mutation.PublishedDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: cve.FieldPublishedData,
		})
	}
	if value, ok := cu.mutation.ModifiedData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cve.FieldModifiedData,
		})
	}
	if cu.mutation.ModifiedDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: cve.FieldModifiedData,
		})
	}
	if cu.mutation.ComponentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cve.ComponentsTable,
			Columns: cve.ComponentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: component.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedComponentsIDs(); len(nodes) > 0 && !cu.mutation.ComponentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cve.ComponentsTable,
			Columns: cve.ComponentsPrimaryKey,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.ComponentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cve.ComponentsTable,
			Columns: cve.ComponentsPrimaryKey,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cu.mutation.VulnerabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.VulnerabilitiesTable,
			Columns: []string{cve.VulnerabilitiesColumn},
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
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.VulnerabilitiesTable,
			Columns: []string{cve.VulnerabilitiesColumn},
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
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.VulnerabilitiesTable,
			Columns: []string{cve.VulnerabilitiesColumn},
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
	if cu.mutation.RulesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.RulesTable,
			Columns: []string{cve.RulesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: cverule.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedRulesIDs(); len(nodes) > 0 && !cu.mutation.RulesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.RulesTable,
			Columns: []string{cve.RulesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RulesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.RulesTable,
			Columns: []string{cve.RulesColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cve.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// CVEUpdateOne is the builder for updating a single CVE entity.
type CVEUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *CVEMutation
}

// SetCveID sets the "cve_id" field.
func (cuo *CVEUpdateOne) SetCveID(s string) *CVEUpdateOne {
	cuo.mutation.SetCveID(s)
	return cuo
}

// SetDescription sets the "description" field.
func (cuo *CVEUpdateOne) SetDescription(s string) *CVEUpdateOne {
	cuo.mutation.SetDescription(s)
	return cuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cuo *CVEUpdateOne) SetNillableDescription(s *string) *CVEUpdateOne {
	if s != nil {
		cuo.SetDescription(*s)
	}
	return cuo
}

// ClearDescription clears the value of the "description" field.
func (cuo *CVEUpdateOne) ClearDescription() *CVEUpdateOne {
	cuo.mutation.ClearDescription()
	return cuo
}

// SetSeverityScore sets the "severity_score" field.
func (cuo *CVEUpdateOne) SetSeverityScore(f float64) *CVEUpdateOne {
	cuo.mutation.ResetSeverityScore()
	cuo.mutation.SetSeverityScore(f)
	return cuo
}

// SetNillableSeverityScore sets the "severity_score" field if the given value is not nil.
func (cuo *CVEUpdateOne) SetNillableSeverityScore(f *float64) *CVEUpdateOne {
	if f != nil {
		cuo.SetSeverityScore(*f)
	}
	return cuo
}

// AddSeverityScore adds f to the "severity_score" field.
func (cuo *CVEUpdateOne) AddSeverityScore(f float64) *CVEUpdateOne {
	cuo.mutation.AddSeverityScore(f)
	return cuo
}

// SetSeverity sets the "severity" field.
func (cuo *CVEUpdateOne) SetSeverity(c cve.Severity) *CVEUpdateOne {
	cuo.mutation.SetSeverity(c)
	return cuo
}

// SetNillableSeverity sets the "severity" field if the given value is not nil.
func (cuo *CVEUpdateOne) SetNillableSeverity(c *cve.Severity) *CVEUpdateOne {
	if c != nil {
		cuo.SetSeverity(*c)
	}
	return cuo
}

// SetPublishedData sets the "published_data" field.
func (cuo *CVEUpdateOne) SetPublishedData(t time.Time) *CVEUpdateOne {
	cuo.mutation.SetPublishedData(t)
	return cuo
}

// SetNillablePublishedData sets the "published_data" field if the given value is not nil.
func (cuo *CVEUpdateOne) SetNillablePublishedData(t *time.Time) *CVEUpdateOne {
	if t != nil {
		cuo.SetPublishedData(*t)
	}
	return cuo
}

// ClearPublishedData clears the value of the "published_data" field.
func (cuo *CVEUpdateOne) ClearPublishedData() *CVEUpdateOne {
	cuo.mutation.ClearPublishedData()
	return cuo
}

// SetModifiedData sets the "modified_data" field.
func (cuo *CVEUpdateOne) SetModifiedData(t time.Time) *CVEUpdateOne {
	cuo.mutation.SetModifiedData(t)
	return cuo
}

// SetNillableModifiedData sets the "modified_data" field if the given value is not nil.
func (cuo *CVEUpdateOne) SetNillableModifiedData(t *time.Time) *CVEUpdateOne {
	if t != nil {
		cuo.SetModifiedData(*t)
	}
	return cuo
}

// ClearModifiedData clears the value of the "modified_data" field.
func (cuo *CVEUpdateOne) ClearModifiedData() *CVEUpdateOne {
	cuo.mutation.ClearModifiedData()
	return cuo
}

// AddComponentIDs adds the "components" edge to the Component entity by IDs.
func (cuo *CVEUpdateOne) AddComponentIDs(ids ...int) *CVEUpdateOne {
	cuo.mutation.AddComponentIDs(ids...)
	return cuo
}

// AddComponents adds the "components" edges to the Component entity.
func (cuo *CVEUpdateOne) AddComponents(c ...*Component) *CVEUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cuo.AddComponentIDs(ids...)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (cuo *CVEUpdateOne) AddVulnerabilityIDs(ids ...int) *CVEUpdateOne {
	cuo.mutation.AddVulnerabilityIDs(ids...)
	return cuo
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (cuo *CVEUpdateOne) AddVulnerabilities(v ...*Vulnerability) *CVEUpdateOne {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return cuo.AddVulnerabilityIDs(ids...)
}

// AddRuleIDs adds the "rules" edge to the CVERule entity by IDs.
func (cuo *CVEUpdateOne) AddRuleIDs(ids ...int) *CVEUpdateOne {
	cuo.mutation.AddRuleIDs(ids...)
	return cuo
}

// AddRules adds the "rules" edges to the CVERule entity.
func (cuo *CVEUpdateOne) AddRules(c ...*CVERule) *CVEUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cuo.AddRuleIDs(ids...)
}

// Mutation returns the CVEMutation object of the builder.
func (cuo *CVEUpdateOne) Mutation() *CVEMutation {
	return cuo.mutation
}

// ClearComponents clears all "components" edges to the Component entity.
func (cuo *CVEUpdateOne) ClearComponents() *CVEUpdateOne {
	cuo.mutation.ClearComponents()
	return cuo
}

// RemoveComponentIDs removes the "components" edge to Component entities by IDs.
func (cuo *CVEUpdateOne) RemoveComponentIDs(ids ...int) *CVEUpdateOne {
	cuo.mutation.RemoveComponentIDs(ids...)
	return cuo
}

// RemoveComponents removes "components" edges to Component entities.
func (cuo *CVEUpdateOne) RemoveComponents(c ...*Component) *CVEUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cuo.RemoveComponentIDs(ids...)
}

// ClearVulnerabilities clears all "vulnerabilities" edges to the Vulnerability entity.
func (cuo *CVEUpdateOne) ClearVulnerabilities() *CVEUpdateOne {
	cuo.mutation.ClearVulnerabilities()
	return cuo
}

// RemoveVulnerabilityIDs removes the "vulnerabilities" edge to Vulnerability entities by IDs.
func (cuo *CVEUpdateOne) RemoveVulnerabilityIDs(ids ...int) *CVEUpdateOne {
	cuo.mutation.RemoveVulnerabilityIDs(ids...)
	return cuo
}

// RemoveVulnerabilities removes "vulnerabilities" edges to Vulnerability entities.
func (cuo *CVEUpdateOne) RemoveVulnerabilities(v ...*Vulnerability) *CVEUpdateOne {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return cuo.RemoveVulnerabilityIDs(ids...)
}

// ClearRules clears all "rules" edges to the CVERule entity.
func (cuo *CVEUpdateOne) ClearRules() *CVEUpdateOne {
	cuo.mutation.ClearRules()
	return cuo
}

// RemoveRuleIDs removes the "rules" edge to CVERule entities by IDs.
func (cuo *CVEUpdateOne) RemoveRuleIDs(ids ...int) *CVEUpdateOne {
	cuo.mutation.RemoveRuleIDs(ids...)
	return cuo
}

// RemoveRules removes "rules" edges to CVERule entities.
func (cuo *CVEUpdateOne) RemoveRules(c ...*CVERule) *CVEUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cuo.RemoveRuleIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *CVEUpdateOne) Select(field string, fields ...string) *CVEUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated CVE entity.
func (cuo *CVEUpdateOne) Save(ctx context.Context) (*CVE, error) {
	var (
		err  error
		node *CVE
	)
	if len(cuo.hooks) == 0 {
		if err = cuo.check(); err != nil {
			return nil, err
		}
		node, err = cuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CVEMutation)
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
func (cuo *CVEUpdateOne) SaveX(ctx context.Context) *CVE {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *CVEUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CVEUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cuo *CVEUpdateOne) check() error {
	if v, ok := cuo.mutation.CveID(); ok {
		if err := cve.CveIDValidator(v); err != nil {
			return &ValidationError{Name: "cve_id", err: fmt.Errorf("ent: validator failed for field \"cve_id\": %w", err)}
		}
	}
	if v, ok := cuo.mutation.Severity(); ok {
		if err := cve.SeverityValidator(v); err != nil {
			return &ValidationError{Name: "severity", err: fmt.Errorf("ent: validator failed for field \"severity\": %w", err)}
		}
	}
	return nil
}

func (cuo *CVEUpdateOne) sqlSave(ctx context.Context) (_node *CVE, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cve.Table,
			Columns: cve.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: cve.FieldID,
			},
		},
	}
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing CVE.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cve.FieldID)
		for _, f := range fields {
			if !cve.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != cve.FieldID {
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
	if value, ok := cuo.mutation.CveID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cve.FieldCveID,
		})
	}
	if value, ok := cuo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cve.FieldDescription,
		})
	}
	if cuo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: cve.FieldDescription,
		})
	}
	if value, ok := cuo.mutation.SeverityScore(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: cve.FieldSeverityScore,
		})
	}
	if value, ok := cuo.mutation.AddedSeverityScore(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: cve.FieldSeverityScore,
		})
	}
	if value, ok := cuo.mutation.Severity(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: cve.FieldSeverity,
		})
	}
	if value, ok := cuo.mutation.PublishedData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cve.FieldPublishedData,
		})
	}
	if cuo.mutation.PublishedDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: cve.FieldPublishedData,
		})
	}
	if value, ok := cuo.mutation.ModifiedData(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cve.FieldModifiedData,
		})
	}
	if cuo.mutation.ModifiedDataCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: cve.FieldModifiedData,
		})
	}
	if cuo.mutation.ComponentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cve.ComponentsTable,
			Columns: cve.ComponentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: component.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedComponentsIDs(); len(nodes) > 0 && !cuo.mutation.ComponentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cve.ComponentsTable,
			Columns: cve.ComponentsPrimaryKey,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.ComponentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cve.ComponentsTable,
			Columns: cve.ComponentsPrimaryKey,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cuo.mutation.VulnerabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.VulnerabilitiesTable,
			Columns: []string{cve.VulnerabilitiesColumn},
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
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.VulnerabilitiesTable,
			Columns: []string{cve.VulnerabilitiesColumn},
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
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.VulnerabilitiesTable,
			Columns: []string{cve.VulnerabilitiesColumn},
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
	if cuo.mutation.RulesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.RulesTable,
			Columns: []string{cve.RulesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: cverule.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedRulesIDs(); len(nodes) > 0 && !cuo.mutation.RulesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.RulesTable,
			Columns: []string{cve.RulesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RulesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   cve.RulesTable,
			Columns: []string{cve.RulesColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &CVE{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cve.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
