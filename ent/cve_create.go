// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/cve"
	"github.com/valocode/bubbly/ent/cverule"
	"github.com/valocode/bubbly/ent/vulnerability"
)

// CVECreate is the builder for creating a CVE entity.
type CVECreate struct {
	config
	mutation *CVEMutation
	hooks    []Hook
}

// SetCveID sets the "cve_id" field.
func (cc *CVECreate) SetCveID(s string) *CVECreate {
	cc.mutation.SetCveID(s)
	return cc
}

// SetDescription sets the "description" field.
func (cc *CVECreate) SetDescription(s string) *CVECreate {
	cc.mutation.SetDescription(s)
	return cc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cc *CVECreate) SetNillableDescription(s *string) *CVECreate {
	if s != nil {
		cc.SetDescription(*s)
	}
	return cc
}

// SetSeverityScore sets the "severity_score" field.
func (cc *CVECreate) SetSeverityScore(f float64) *CVECreate {
	cc.mutation.SetSeverityScore(f)
	return cc
}

// SetNillableSeverityScore sets the "severity_score" field if the given value is not nil.
func (cc *CVECreate) SetNillableSeverityScore(f *float64) *CVECreate {
	if f != nil {
		cc.SetSeverityScore(*f)
	}
	return cc
}

// SetSeverity sets the "severity" field.
func (cc *CVECreate) SetSeverity(c cve.Severity) *CVECreate {
	cc.mutation.SetSeverity(c)
	return cc
}

// SetNillableSeverity sets the "severity" field if the given value is not nil.
func (cc *CVECreate) SetNillableSeverity(c *cve.Severity) *CVECreate {
	if c != nil {
		cc.SetSeverity(*c)
	}
	return cc
}

// SetPublishedData sets the "published_data" field.
func (cc *CVECreate) SetPublishedData(t time.Time) *CVECreate {
	cc.mutation.SetPublishedData(t)
	return cc
}

// SetNillablePublishedData sets the "published_data" field if the given value is not nil.
func (cc *CVECreate) SetNillablePublishedData(t *time.Time) *CVECreate {
	if t != nil {
		cc.SetPublishedData(*t)
	}
	return cc
}

// SetModifiedData sets the "modified_data" field.
func (cc *CVECreate) SetModifiedData(t time.Time) *CVECreate {
	cc.mutation.SetModifiedData(t)
	return cc
}

// SetNillableModifiedData sets the "modified_data" field if the given value is not nil.
func (cc *CVECreate) SetNillableModifiedData(t *time.Time) *CVECreate {
	if t != nil {
		cc.SetModifiedData(*t)
	}
	return cc
}

// AddComponentIDs adds the "components" edge to the Component entity by IDs.
func (cc *CVECreate) AddComponentIDs(ids ...int) *CVECreate {
	cc.mutation.AddComponentIDs(ids...)
	return cc
}

// AddComponents adds the "components" edges to the Component entity.
func (cc *CVECreate) AddComponents(c ...*Component) *CVECreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cc.AddComponentIDs(ids...)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (cc *CVECreate) AddVulnerabilityIDs(ids ...int) *CVECreate {
	cc.mutation.AddVulnerabilityIDs(ids...)
	return cc
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (cc *CVECreate) AddVulnerabilities(v ...*Vulnerability) *CVECreate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return cc.AddVulnerabilityIDs(ids...)
}

// AddRuleIDs adds the "rules" edge to the CVERule entity by IDs.
func (cc *CVECreate) AddRuleIDs(ids ...int) *CVECreate {
	cc.mutation.AddRuleIDs(ids...)
	return cc
}

// AddRules adds the "rules" edges to the CVERule entity.
func (cc *CVECreate) AddRules(c ...*CVERule) *CVECreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cc.AddRuleIDs(ids...)
}

// Mutation returns the CVEMutation object of the builder.
func (cc *CVECreate) Mutation() *CVEMutation {
	return cc.mutation
}

// Save creates the CVE in the database.
func (cc *CVECreate) Save(ctx context.Context) (*CVE, error) {
	var (
		err  error
		node *CVE
	)
	if err := cc.defaults(); err != nil {
		return nil, err
	}
	if len(cc.hooks) == 0 {
		if err = cc.check(); err != nil {
			return nil, err
		}
		node, err = cc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CVEMutation)
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
func (cc *CVECreate) SaveX(ctx context.Context) *CVE {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *CVECreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *CVECreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *CVECreate) defaults() error {
	if _, ok := cc.mutation.SeverityScore(); !ok {
		v := cve.DefaultSeverityScore
		cc.mutation.SetSeverityScore(v)
	}
	if _, ok := cc.mutation.Severity(); !ok {
		v := cve.DefaultSeverity
		cc.mutation.SetSeverity(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (cc *CVECreate) check() error {
	if _, ok := cc.mutation.CveID(); !ok {
		return &ValidationError{Name: "cve_id", err: errors.New(`ent: missing required field "cve_id"`)}
	}
	if v, ok := cc.mutation.CveID(); ok {
		if err := cve.CveIDValidator(v); err != nil {
			return &ValidationError{Name: "cve_id", err: fmt.Errorf(`ent: validator failed for field "cve_id": %w`, err)}
		}
	}
	if _, ok := cc.mutation.SeverityScore(); !ok {
		return &ValidationError{Name: "severity_score", err: errors.New(`ent: missing required field "severity_score"`)}
	}
	if _, ok := cc.mutation.Severity(); !ok {
		return &ValidationError{Name: "severity", err: errors.New(`ent: missing required field "severity"`)}
	}
	if v, ok := cc.mutation.Severity(); ok {
		if err := cve.SeverityValidator(v); err != nil {
			return &ValidationError{Name: "severity", err: fmt.Errorf(`ent: validator failed for field "severity": %w`, err)}
		}
	}
	return nil
}

func (cc *CVECreate) sqlSave(ctx context.Context) (*CVE, error) {
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

func (cc *CVECreate) createSpec() (*CVE, *sqlgraph.CreateSpec) {
	var (
		_node = &CVE{config: cc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: cve.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: cve.FieldID,
			},
		}
	)
	if value, ok := cc.mutation.CveID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cve.FieldCveID,
		})
		_node.CveID = value
	}
	if value, ok := cc.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cve.FieldDescription,
		})
		_node.Description = value
	}
	if value, ok := cc.mutation.SeverityScore(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: cve.FieldSeverityScore,
		})
		_node.SeverityScore = value
	}
	if value, ok := cc.mutation.Severity(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: cve.FieldSeverity,
		})
		_node.Severity = value
	}
	if value, ok := cc.mutation.PublishedData(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cve.FieldPublishedData,
		})
		_node.PublishedData = value
	}
	if value, ok := cc.mutation.ModifiedData(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cve.FieldModifiedData,
		})
		_node.ModifiedData = value
	}
	if nodes := cc.mutation.ComponentsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.RulesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// CVECreateBulk is the builder for creating many CVE entities in bulk.
type CVECreateBulk struct {
	config
	builders []*CVECreate
}

// Save creates the CVE entities in the database.
func (ccb *CVECreateBulk) Save(ctx context.Context) ([]*CVE, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*CVE, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CVEMutation)
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
func (ccb *CVECreateBulk) SaveX(ctx context.Context) []*CVE {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *CVECreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *CVECreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}