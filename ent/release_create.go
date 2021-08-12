// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/artifact"
	"github.com/valocode/bubbly/ent/codescan"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasecomponent"
	"github.com/valocode/bubbly/ent/releaseentry"
	"github.com/valocode/bubbly/ent/releasevulnerability"
	"github.com/valocode/bubbly/ent/testrun"
	"github.com/valocode/bubbly/ent/vulnerabilityreview"
)

// ReleaseCreate is the builder for creating a Release entity.
type ReleaseCreate struct {
	config
	mutation *ReleaseMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (rc *ReleaseCreate) SetName(s string) *ReleaseCreate {
	rc.mutation.SetName(s)
	return rc
}

// SetVersion sets the "version" field.
func (rc *ReleaseCreate) SetVersion(s string) *ReleaseCreate {
	rc.mutation.SetVersion(s)
	return rc
}

// SetStatus sets the "status" field.
func (rc *ReleaseCreate) SetStatus(r release.Status) *ReleaseCreate {
	rc.mutation.SetStatus(r)
	return rc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (rc *ReleaseCreate) SetNillableStatus(r *release.Status) *ReleaseCreate {
	if r != nil {
		rc.SetStatus(*r)
	}
	return rc
}

// AddSubreleaseIDs adds the "subreleases" edge to the Release entity by IDs.
func (rc *ReleaseCreate) AddSubreleaseIDs(ids ...int) *ReleaseCreate {
	rc.mutation.AddSubreleaseIDs(ids...)
	return rc
}

// AddSubreleases adds the "subreleases" edges to the Release entity.
func (rc *ReleaseCreate) AddSubreleases(r ...*Release) *ReleaseCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rc.AddSubreleaseIDs(ids...)
}

// AddDependencyIDs adds the "dependencies" edge to the Release entity by IDs.
func (rc *ReleaseCreate) AddDependencyIDs(ids ...int) *ReleaseCreate {
	rc.mutation.AddDependencyIDs(ids...)
	return rc
}

// AddDependencies adds the "dependencies" edges to the Release entity.
func (rc *ReleaseCreate) AddDependencies(r ...*Release) *ReleaseCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rc.AddDependencyIDs(ids...)
}

// SetCommitID sets the "commit" edge to the GitCommit entity by ID.
func (rc *ReleaseCreate) SetCommitID(id int) *ReleaseCreate {
	rc.mutation.SetCommitID(id)
	return rc
}

// SetCommit sets the "commit" edge to the GitCommit entity.
func (rc *ReleaseCreate) SetCommit(g *GitCommit) *ReleaseCreate {
	return rc.SetCommitID(g.ID)
}

// AddLogIDs adds the "log" edge to the ReleaseEntry entity by IDs.
func (rc *ReleaseCreate) AddLogIDs(ids ...int) *ReleaseCreate {
	rc.mutation.AddLogIDs(ids...)
	return rc
}

// AddLog adds the "log" edges to the ReleaseEntry entity.
func (rc *ReleaseCreate) AddLog(r ...*ReleaseEntry) *ReleaseCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rc.AddLogIDs(ids...)
}

// AddArtifactIDs adds the "artifacts" edge to the Artifact entity by IDs.
func (rc *ReleaseCreate) AddArtifactIDs(ids ...int) *ReleaseCreate {
	rc.mutation.AddArtifactIDs(ids...)
	return rc
}

// AddArtifacts adds the "artifacts" edges to the Artifact entity.
func (rc *ReleaseCreate) AddArtifacts(a ...*Artifact) *ReleaseCreate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return rc.AddArtifactIDs(ids...)
}

// AddComponentIDs adds the "components" edge to the ReleaseComponent entity by IDs.
func (rc *ReleaseCreate) AddComponentIDs(ids ...int) *ReleaseCreate {
	rc.mutation.AddComponentIDs(ids...)
	return rc
}

// AddComponents adds the "components" edges to the ReleaseComponent entity.
func (rc *ReleaseCreate) AddComponents(r ...*ReleaseComponent) *ReleaseCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rc.AddComponentIDs(ids...)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the ReleaseVulnerability entity by IDs.
func (rc *ReleaseCreate) AddVulnerabilityIDs(ids ...int) *ReleaseCreate {
	rc.mutation.AddVulnerabilityIDs(ids...)
	return rc
}

// AddVulnerabilities adds the "vulnerabilities" edges to the ReleaseVulnerability entity.
func (rc *ReleaseCreate) AddVulnerabilities(r ...*ReleaseVulnerability) *ReleaseCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rc.AddVulnerabilityIDs(ids...)
}

// AddCodeScanIDs adds the "code_scans" edge to the CodeScan entity by IDs.
func (rc *ReleaseCreate) AddCodeScanIDs(ids ...int) *ReleaseCreate {
	rc.mutation.AddCodeScanIDs(ids...)
	return rc
}

// AddCodeScans adds the "code_scans" edges to the CodeScan entity.
func (rc *ReleaseCreate) AddCodeScans(c ...*CodeScan) *ReleaseCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return rc.AddCodeScanIDs(ids...)
}

// AddTestRunIDs adds the "test_runs" edge to the TestRun entity by IDs.
func (rc *ReleaseCreate) AddTestRunIDs(ids ...int) *ReleaseCreate {
	rc.mutation.AddTestRunIDs(ids...)
	return rc
}

// AddTestRuns adds the "test_runs" edges to the TestRun entity.
func (rc *ReleaseCreate) AddTestRuns(t ...*TestRun) *ReleaseCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return rc.AddTestRunIDs(ids...)
}

// AddVulnerabilityReviewIDs adds the "vulnerability_reviews" edge to the VulnerabilityReview entity by IDs.
func (rc *ReleaseCreate) AddVulnerabilityReviewIDs(ids ...int) *ReleaseCreate {
	rc.mutation.AddVulnerabilityReviewIDs(ids...)
	return rc
}

// AddVulnerabilityReviews adds the "vulnerability_reviews" edges to the VulnerabilityReview entity.
func (rc *ReleaseCreate) AddVulnerabilityReviews(v ...*VulnerabilityReview) *ReleaseCreate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return rc.AddVulnerabilityReviewIDs(ids...)
}

// Mutation returns the ReleaseMutation object of the builder.
func (rc *ReleaseCreate) Mutation() *ReleaseMutation {
	return rc.mutation
}

// Save creates the Release in the database.
func (rc *ReleaseCreate) Save(ctx context.Context) (*Release, error) {
	var (
		err  error
		node *Release
	)
	rc.defaults()
	if len(rc.hooks) == 0 {
		if err = rc.check(); err != nil {
			return nil, err
		}
		node, err = rc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ReleaseMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rc.check(); err != nil {
				return nil, err
			}
			rc.mutation = mutation
			if node, err = rc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(rc.hooks) - 1; i >= 0; i-- {
			if rc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = rc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rc *ReleaseCreate) SaveX(ctx context.Context) *Release {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rc *ReleaseCreate) Exec(ctx context.Context) error {
	_, err := rc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rc *ReleaseCreate) ExecX(ctx context.Context) {
	if err := rc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rc *ReleaseCreate) defaults() {
	if _, ok := rc.mutation.Status(); !ok {
		v := release.DefaultStatus
		rc.mutation.SetStatus(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *ReleaseCreate) check() error {
	if _, ok := rc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "name"`)}
	}
	if v, ok := rc.mutation.Name(); ok {
		if err := release.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "name": %w`, err)}
		}
	}
	if _, ok := rc.mutation.Version(); !ok {
		return &ValidationError{Name: "version", err: errors.New(`ent: missing required field "version"`)}
	}
	if v, ok := rc.mutation.Version(); ok {
		if err := release.VersionValidator(v); err != nil {
			return &ValidationError{Name: "version", err: fmt.Errorf(`ent: validator failed for field "version": %w`, err)}
		}
	}
	if _, ok := rc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "status"`)}
	}
	if v, ok := rc.mutation.Status(); ok {
		if err := release.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "status": %w`, err)}
		}
	}
	if _, ok := rc.mutation.CommitID(); !ok {
		return &ValidationError{Name: "commit", err: errors.New("ent: missing required edge \"commit\"")}
	}
	return nil
}

func (rc *ReleaseCreate) sqlSave(ctx context.Context) (*Release, error) {
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (rc *ReleaseCreate) createSpec() (*Release, *sqlgraph.CreateSpec) {
	var (
		_node = &Release{config: rc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: release.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: release.FieldID,
			},
		}
	)
	if value, ok := rc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: release.FieldName,
		})
		_node.Name = value
	}
	if value, ok := rc.mutation.Version(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: release.FieldVersion,
		})
		_node.Version = value
	}
	if value, ok := rc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: release.FieldStatus,
		})
		_node.Status = value
	}
	if nodes := rc.mutation.SubreleasesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   release.SubreleasesTable,
			Columns: release.SubreleasesPrimaryKey,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.DependenciesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   release.DependenciesTable,
			Columns: release.DependenciesPrimaryKey,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.CommitIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   release.CommitTable,
			Columns: []string{release.CommitColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: gitcommit.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.git_commit_release = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.LogIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   release.LogTable,
			Columns: []string{release.LogColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.ArtifactsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   release.ArtifactsTable,
			Columns: []string{release.ArtifactsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: artifact.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.ComponentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   release.ComponentsTable,
			Columns: []string{release.ComponentsColumn},
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
	if nodes := rc.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   release.VulnerabilitiesTable,
			Columns: []string{release.VulnerabilitiesColumn},
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
	if nodes := rc.mutation.CodeScansIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   release.CodeScansTable,
			Columns: []string{release.CodeScansColumn},
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
	if nodes := rc.mutation.TestRunsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   release.TestRunsTable,
			Columns: []string{release.TestRunsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: testrun.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.VulnerabilityReviewsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   release.VulnerabilityReviewsTable,
			Columns: release.VulnerabilityReviewsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: vulnerabilityreview.FieldID,
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

// ReleaseCreateBulk is the builder for creating many Release entities in bulk.
type ReleaseCreateBulk struct {
	config
	builders []*ReleaseCreate
}

// Save creates the Release entities in the database.
func (rcb *ReleaseCreateBulk) Save(ctx context.Context) ([]*Release, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Release, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ReleaseMutation)
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
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *ReleaseCreateBulk) SaveX(ctx context.Context) []*Release {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rcb *ReleaseCreateBulk) Exec(ctx context.Context) error {
	_, err := rcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcb *ReleaseCreateBulk) ExecX(ctx context.Context) {
	if err := rcb.Exec(ctx); err != nil {
		panic(err)
	}
}
