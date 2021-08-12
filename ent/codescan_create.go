// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/codescan"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasecomponent"
	"github.com/valocode/bubbly/ent/releaseentry"
	"github.com/valocode/bubbly/ent/releasevulnerability"
)

// CodeScanCreate is the builder for creating a CodeScan entity.
type CodeScanCreate struct {
	config
	mutation *CodeScanMutation
	hooks    []Hook
}

// SetTool sets the "tool" field.
func (csc *CodeScanCreate) SetTool(s string) *CodeScanCreate {
	csc.mutation.SetTool(s)
	return csc
}

// SetReleaseID sets the "release" edge to the Release entity by ID.
func (csc *CodeScanCreate) SetReleaseID(id int) *CodeScanCreate {
	csc.mutation.SetReleaseID(id)
	return csc
}

// SetRelease sets the "release" edge to the Release entity.
func (csc *CodeScanCreate) SetRelease(r *Release) *CodeScanCreate {
	return csc.SetReleaseID(r.ID)
}

// SetEntryID sets the "entry" edge to the ReleaseEntry entity by ID.
func (csc *CodeScanCreate) SetEntryID(id int) *CodeScanCreate {
	csc.mutation.SetEntryID(id)
	return csc
}

// SetNillableEntryID sets the "entry" edge to the ReleaseEntry entity by ID if the given value is not nil.
func (csc *CodeScanCreate) SetNillableEntryID(id *int) *CodeScanCreate {
	if id != nil {
		csc = csc.SetEntryID(*id)
	}
	return csc
}

// SetEntry sets the "entry" edge to the ReleaseEntry entity.
func (csc *CodeScanCreate) SetEntry(r *ReleaseEntry) *CodeScanCreate {
	return csc.SetEntryID(r.ID)
}

// AddIssueIDs adds the "issues" edge to the CodeIssue entity by IDs.
func (csc *CodeScanCreate) AddIssueIDs(ids ...int) *CodeScanCreate {
	csc.mutation.AddIssueIDs(ids...)
	return csc
}

// AddIssues adds the "issues" edges to the CodeIssue entity.
func (csc *CodeScanCreate) AddIssues(c ...*CodeIssue) *CodeScanCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return csc.AddIssueIDs(ids...)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the ReleaseVulnerability entity by IDs.
func (csc *CodeScanCreate) AddVulnerabilityIDs(ids ...int) *CodeScanCreate {
	csc.mutation.AddVulnerabilityIDs(ids...)
	return csc
}

// AddVulnerabilities adds the "vulnerabilities" edges to the ReleaseVulnerability entity.
func (csc *CodeScanCreate) AddVulnerabilities(r ...*ReleaseVulnerability) *CodeScanCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return csc.AddVulnerabilityIDs(ids...)
}

// AddComponentIDs adds the "components" edge to the ReleaseComponent entity by IDs.
func (csc *CodeScanCreate) AddComponentIDs(ids ...int) *CodeScanCreate {
	csc.mutation.AddComponentIDs(ids...)
	return csc
}

// AddComponents adds the "components" edges to the ReleaseComponent entity.
func (csc *CodeScanCreate) AddComponents(r ...*ReleaseComponent) *CodeScanCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return csc.AddComponentIDs(ids...)
}

// Mutation returns the CodeScanMutation object of the builder.
func (csc *CodeScanCreate) Mutation() *CodeScanMutation {
	return csc.mutation
}

// Save creates the CodeScan in the database.
func (csc *CodeScanCreate) Save(ctx context.Context) (*CodeScan, error) {
	var (
		err  error
		node *CodeScan
	)
	if len(csc.hooks) == 0 {
		if err = csc.check(); err != nil {
			return nil, err
		}
		node, err = csc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CodeScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = csc.check(); err != nil {
				return nil, err
			}
			csc.mutation = mutation
			if node, err = csc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(csc.hooks) - 1; i >= 0; i-- {
			if csc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = csc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, csc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (csc *CodeScanCreate) SaveX(ctx context.Context) *CodeScan {
	v, err := csc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (csc *CodeScanCreate) Exec(ctx context.Context) error {
	_, err := csc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (csc *CodeScanCreate) ExecX(ctx context.Context) {
	if err := csc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (csc *CodeScanCreate) check() error {
	if _, ok := csc.mutation.Tool(); !ok {
		return &ValidationError{Name: "tool", err: errors.New(`ent: missing required field "tool"`)}
	}
	if v, ok := csc.mutation.Tool(); ok {
		if err := codescan.ToolValidator(v); err != nil {
			return &ValidationError{Name: "tool", err: fmt.Errorf(`ent: validator failed for field "tool": %w`, err)}
		}
	}
	if _, ok := csc.mutation.ReleaseID(); !ok {
		return &ValidationError{Name: "release", err: errors.New("ent: missing required edge \"release\"")}
	}
	return nil
}

func (csc *CodeScanCreate) sqlSave(ctx context.Context) (*CodeScan, error) {
	_node, _spec := csc.createSpec()
	if err := sqlgraph.CreateNode(ctx, csc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (csc *CodeScanCreate) createSpec() (*CodeScan, *sqlgraph.CreateSpec) {
	var (
		_node = &CodeScan{config: csc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: codescan.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: codescan.FieldID,
			},
		}
	)
	if value, ok := csc.mutation.Tool(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: codescan.FieldTool,
		})
		_node.Tool = value
	}
	if nodes := csc.mutation.ReleaseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   codescan.ReleaseTable,
			Columns: []string{codescan.ReleaseColumn},
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
		_node.code_scan_release = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := csc.mutation.EntryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   codescan.EntryTable,
			Columns: []string{codescan.EntryColumn},
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
		_node.release_entry_code_scan = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := csc.mutation.IssuesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   codescan.IssuesTable,
			Columns: []string{codescan.IssuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: codeissue.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := csc.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   codescan.VulnerabilitiesTable,
			Columns: codescan.VulnerabilitiesPrimaryKey,
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
	if nodes := csc.mutation.ComponentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   codescan.ComponentsTable,
			Columns: codescan.ComponentsPrimaryKey,
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

// CodeScanCreateBulk is the builder for creating many CodeScan entities in bulk.
type CodeScanCreateBulk struct {
	config
	builders []*CodeScanCreate
}

// Save creates the CodeScan entities in the database.
func (cscb *CodeScanCreateBulk) Save(ctx context.Context) ([]*CodeScan, error) {
	specs := make([]*sqlgraph.CreateSpec, len(cscb.builders))
	nodes := make([]*CodeScan, len(cscb.builders))
	mutators := make([]Mutator, len(cscb.builders))
	for i := range cscb.builders {
		func(i int, root context.Context) {
			builder := cscb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CodeScanMutation)
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
					_, err = mutators[i+1].Mutate(root, cscb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, cscb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, cscb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (cscb *CodeScanCreateBulk) SaveX(ctx context.Context) []*CodeScan {
	v, err := cscb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cscb *CodeScanCreateBulk) Exec(ctx context.Context) error {
	_, err := cscb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cscb *CodeScanCreateBulk) ExecX(ctx context.Context) {
	if err := cscb.Exec(ctx); err != nil {
		panic(err)
	}
}
