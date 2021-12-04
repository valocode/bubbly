// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/organization"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasepolicy"
	"github.com/valocode/bubbly/ent/repo"
	schema "github.com/valocode/bubbly/ent/schema/types"
	"github.com/valocode/bubbly/ent/vulnerabilityreview"
)

// RepoUpdate is the builder for updating Repo entities.
type RepoUpdate struct {
	config
	hooks    []Hook
	mutation *RepoMutation
}

// Where appends a list predicates to the RepoUpdate builder.
func (ru *RepoUpdate) Where(ps ...predicate.Repo) *RepoUpdate {
	ru.mutation.Where(ps...)
	return ru
}

// SetName sets the "name" field.
func (ru *RepoUpdate) SetName(s string) *RepoUpdate {
	ru.mutation.SetName(s)
	return ru
}

// SetDefaultBranch sets the "default_branch" field.
func (ru *RepoUpdate) SetDefaultBranch(s string) *RepoUpdate {
	ru.mutation.SetDefaultBranch(s)
	return ru
}

// SetNillableDefaultBranch sets the "default_branch" field if the given value is not nil.
func (ru *RepoUpdate) SetNillableDefaultBranch(s *string) *RepoUpdate {
	if s != nil {
		ru.SetDefaultBranch(*s)
	}
	return ru
}

// SetLabels sets the "labels" field.
func (ru *RepoUpdate) SetLabels(s schema.Labels) *RepoUpdate {
	ru.mutation.SetLabels(s)
	return ru
}

// ClearLabels clears the value of the "labels" field.
func (ru *RepoUpdate) ClearLabels() *RepoUpdate {
	ru.mutation.ClearLabels()
	return ru
}

// SetOwnerID sets the "owner" edge to the Organization entity by ID.
func (ru *RepoUpdate) SetOwnerID(id int) *RepoUpdate {
	ru.mutation.SetOwnerID(id)
	return ru
}

// SetOwner sets the "owner" edge to the Organization entity.
func (ru *RepoUpdate) SetOwner(o *Organization) *RepoUpdate {
	return ru.SetOwnerID(o.ID)
}

// SetProjectID sets the "project" edge to the Project entity by ID.
func (ru *RepoUpdate) SetProjectID(id int) *RepoUpdate {
	ru.mutation.SetProjectID(id)
	return ru
}

// SetProject sets the "project" edge to the Project entity.
func (ru *RepoUpdate) SetProject(p *Project) *RepoUpdate {
	return ru.SetProjectID(p.ID)
}

// SetHeadID sets the "head" edge to the Release entity by ID.
func (ru *RepoUpdate) SetHeadID(id int) *RepoUpdate {
	ru.mutation.SetHeadID(id)
	return ru
}

// SetNillableHeadID sets the "head" edge to the Release entity by ID if the given value is not nil.
func (ru *RepoUpdate) SetNillableHeadID(id *int) *RepoUpdate {
	if id != nil {
		ru = ru.SetHeadID(*id)
	}
	return ru
}

// SetHead sets the "head" edge to the Release entity.
func (ru *RepoUpdate) SetHead(r *Release) *RepoUpdate {
	return ru.SetHeadID(r.ID)
}

// AddCommitIDs adds the "commits" edge to the GitCommit entity by IDs.
func (ru *RepoUpdate) AddCommitIDs(ids ...int) *RepoUpdate {
	ru.mutation.AddCommitIDs(ids...)
	return ru
}

// AddCommits adds the "commits" edges to the GitCommit entity.
func (ru *RepoUpdate) AddCommits(g ...*GitCommit) *RepoUpdate {
	ids := make([]int, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return ru.AddCommitIDs(ids...)
}

// AddVulnerabilityReviewIDs adds the "vulnerability_reviews" edge to the VulnerabilityReview entity by IDs.
func (ru *RepoUpdate) AddVulnerabilityReviewIDs(ids ...int) *RepoUpdate {
	ru.mutation.AddVulnerabilityReviewIDs(ids...)
	return ru
}

// AddVulnerabilityReviews adds the "vulnerability_reviews" edges to the VulnerabilityReview entity.
func (ru *RepoUpdate) AddVulnerabilityReviews(v ...*VulnerabilityReview) *RepoUpdate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return ru.AddVulnerabilityReviewIDs(ids...)
}

// AddPolicyIDs adds the "policies" edge to the ReleasePolicy entity by IDs.
func (ru *RepoUpdate) AddPolicyIDs(ids ...int) *RepoUpdate {
	ru.mutation.AddPolicyIDs(ids...)
	return ru
}

// AddPolicies adds the "policies" edges to the ReleasePolicy entity.
func (ru *RepoUpdate) AddPolicies(r ...*ReleasePolicy) *RepoUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ru.AddPolicyIDs(ids...)
}

// Mutation returns the RepoMutation object of the builder.
func (ru *RepoUpdate) Mutation() *RepoMutation {
	return ru.mutation
}

// ClearOwner clears the "owner" edge to the Organization entity.
func (ru *RepoUpdate) ClearOwner() *RepoUpdate {
	ru.mutation.ClearOwner()
	return ru
}

// ClearProject clears the "project" edge to the Project entity.
func (ru *RepoUpdate) ClearProject() *RepoUpdate {
	ru.mutation.ClearProject()
	return ru
}

// ClearHead clears the "head" edge to the Release entity.
func (ru *RepoUpdate) ClearHead() *RepoUpdate {
	ru.mutation.ClearHead()
	return ru
}

// ClearCommits clears all "commits" edges to the GitCommit entity.
func (ru *RepoUpdate) ClearCommits() *RepoUpdate {
	ru.mutation.ClearCommits()
	return ru
}

// RemoveCommitIDs removes the "commits" edge to GitCommit entities by IDs.
func (ru *RepoUpdate) RemoveCommitIDs(ids ...int) *RepoUpdate {
	ru.mutation.RemoveCommitIDs(ids...)
	return ru
}

// RemoveCommits removes "commits" edges to GitCommit entities.
func (ru *RepoUpdate) RemoveCommits(g ...*GitCommit) *RepoUpdate {
	ids := make([]int, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return ru.RemoveCommitIDs(ids...)
}

// ClearVulnerabilityReviews clears all "vulnerability_reviews" edges to the VulnerabilityReview entity.
func (ru *RepoUpdate) ClearVulnerabilityReviews() *RepoUpdate {
	ru.mutation.ClearVulnerabilityReviews()
	return ru
}

// RemoveVulnerabilityReviewIDs removes the "vulnerability_reviews" edge to VulnerabilityReview entities by IDs.
func (ru *RepoUpdate) RemoveVulnerabilityReviewIDs(ids ...int) *RepoUpdate {
	ru.mutation.RemoveVulnerabilityReviewIDs(ids...)
	return ru
}

// RemoveVulnerabilityReviews removes "vulnerability_reviews" edges to VulnerabilityReview entities.
func (ru *RepoUpdate) RemoveVulnerabilityReviews(v ...*VulnerabilityReview) *RepoUpdate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return ru.RemoveVulnerabilityReviewIDs(ids...)
}

// ClearPolicies clears all "policies" edges to the ReleasePolicy entity.
func (ru *RepoUpdate) ClearPolicies() *RepoUpdate {
	ru.mutation.ClearPolicies()
	return ru
}

// RemovePolicyIDs removes the "policies" edge to ReleasePolicy entities by IDs.
func (ru *RepoUpdate) RemovePolicyIDs(ids ...int) *RepoUpdate {
	ru.mutation.RemovePolicyIDs(ids...)
	return ru
}

// RemovePolicies removes "policies" edges to ReleasePolicy entities.
func (ru *RepoUpdate) RemovePolicies(r ...*ReleasePolicy) *RepoUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ru.RemovePolicyIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ru *RepoUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ru.hooks) == 0 {
		if err = ru.check(); err != nil {
			return 0, err
		}
		affected, err = ru.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RepoMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ru.check(); err != nil {
				return 0, err
			}
			ru.mutation = mutation
			affected, err = ru.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ru.hooks) - 1; i >= 0; i-- {
			if ru.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ru.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ru.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (ru *RepoUpdate) SaveX(ctx context.Context) int {
	affected, err := ru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ru *RepoUpdate) Exec(ctx context.Context) error {
	_, err := ru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ru *RepoUpdate) ExecX(ctx context.Context) {
	if err := ru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ru *RepoUpdate) check() error {
	if v, ok := ru.mutation.Name(); ok {
		if err := repo.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	if v, ok := ru.mutation.DefaultBranch(); ok {
		if err := repo.DefaultBranchValidator(v); err != nil {
			return &ValidationError{Name: "default_branch", err: fmt.Errorf("ent: validator failed for field \"default_branch\": %w", err)}
		}
	}
	if _, ok := ru.mutation.OwnerID(); ru.mutation.OwnerCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"owner\"")
	}
	if _, ok := ru.mutation.ProjectID(); ru.mutation.ProjectCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"project\"")
	}
	return nil
}

func (ru *RepoUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   repo.Table,
			Columns: repo.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: repo.FieldID,
			},
		},
	}
	if ps := ru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ru.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: repo.FieldName,
		})
	}
	if value, ok := ru.mutation.DefaultBranch(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: repo.FieldDefaultBranch,
		})
	}
	if value, ok := ru.mutation.Labels(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: repo.FieldLabels,
		})
	}
	if ru.mutation.LabelsCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: repo.FieldLabels,
		})
	}
	if ru.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   repo.OwnerTable,
			Columns: []string{repo.OwnerColumn},
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
	if nodes := ru.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   repo.OwnerTable,
			Columns: []string{repo.OwnerColumn},
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
	if ru.mutation.ProjectCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   repo.ProjectTable,
			Columns: []string{repo.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   repo.ProjectTable,
			Columns: []string{repo.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ru.mutation.HeadCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   repo.HeadTable,
			Columns: []string{repo.HeadColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: release.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.HeadIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   repo.HeadTable,
			Columns: []string{repo.HeadColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ru.mutation.CommitsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   repo.CommitsTable,
			Columns: []string{repo.CommitsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: gitcommit.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.RemovedCommitsIDs(); len(nodes) > 0 && !ru.mutation.CommitsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   repo.CommitsTable,
			Columns: []string{repo.CommitsColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.CommitsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   repo.CommitsTable,
			Columns: []string{repo.CommitsColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ru.mutation.VulnerabilityReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.VulnerabilityReviewsTable,
			Columns: repo.VulnerabilityReviewsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: vulnerabilityreview.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.RemovedVulnerabilityReviewsIDs(); len(nodes) > 0 && !ru.mutation.VulnerabilityReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.VulnerabilityReviewsTable,
			Columns: repo.VulnerabilityReviewsPrimaryKey,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.VulnerabilityReviewsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.VulnerabilityReviewsTable,
			Columns: repo.VulnerabilityReviewsPrimaryKey,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ru.mutation.PoliciesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.PoliciesTable,
			Columns: repo.PoliciesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releasepolicy.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.RemovedPoliciesIDs(); len(nodes) > 0 && !ru.mutation.PoliciesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.PoliciesTable,
			Columns: repo.PoliciesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releasepolicy.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.PoliciesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.PoliciesTable,
			Columns: repo.PoliciesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releasepolicy.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{repo.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// RepoUpdateOne is the builder for updating a single Repo entity.
type RepoUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *RepoMutation
}

// SetName sets the "name" field.
func (ruo *RepoUpdateOne) SetName(s string) *RepoUpdateOne {
	ruo.mutation.SetName(s)
	return ruo
}

// SetDefaultBranch sets the "default_branch" field.
func (ruo *RepoUpdateOne) SetDefaultBranch(s string) *RepoUpdateOne {
	ruo.mutation.SetDefaultBranch(s)
	return ruo
}

// SetNillableDefaultBranch sets the "default_branch" field if the given value is not nil.
func (ruo *RepoUpdateOne) SetNillableDefaultBranch(s *string) *RepoUpdateOne {
	if s != nil {
		ruo.SetDefaultBranch(*s)
	}
	return ruo
}

// SetLabels sets the "labels" field.
func (ruo *RepoUpdateOne) SetLabels(s schema.Labels) *RepoUpdateOne {
	ruo.mutation.SetLabels(s)
	return ruo
}

// ClearLabels clears the value of the "labels" field.
func (ruo *RepoUpdateOne) ClearLabels() *RepoUpdateOne {
	ruo.mutation.ClearLabels()
	return ruo
}

// SetOwnerID sets the "owner" edge to the Organization entity by ID.
func (ruo *RepoUpdateOne) SetOwnerID(id int) *RepoUpdateOne {
	ruo.mutation.SetOwnerID(id)
	return ruo
}

// SetOwner sets the "owner" edge to the Organization entity.
func (ruo *RepoUpdateOne) SetOwner(o *Organization) *RepoUpdateOne {
	return ruo.SetOwnerID(o.ID)
}

// SetProjectID sets the "project" edge to the Project entity by ID.
func (ruo *RepoUpdateOne) SetProjectID(id int) *RepoUpdateOne {
	ruo.mutation.SetProjectID(id)
	return ruo
}

// SetProject sets the "project" edge to the Project entity.
func (ruo *RepoUpdateOne) SetProject(p *Project) *RepoUpdateOne {
	return ruo.SetProjectID(p.ID)
}

// SetHeadID sets the "head" edge to the Release entity by ID.
func (ruo *RepoUpdateOne) SetHeadID(id int) *RepoUpdateOne {
	ruo.mutation.SetHeadID(id)
	return ruo
}

// SetNillableHeadID sets the "head" edge to the Release entity by ID if the given value is not nil.
func (ruo *RepoUpdateOne) SetNillableHeadID(id *int) *RepoUpdateOne {
	if id != nil {
		ruo = ruo.SetHeadID(*id)
	}
	return ruo
}

// SetHead sets the "head" edge to the Release entity.
func (ruo *RepoUpdateOne) SetHead(r *Release) *RepoUpdateOne {
	return ruo.SetHeadID(r.ID)
}

// AddCommitIDs adds the "commits" edge to the GitCommit entity by IDs.
func (ruo *RepoUpdateOne) AddCommitIDs(ids ...int) *RepoUpdateOne {
	ruo.mutation.AddCommitIDs(ids...)
	return ruo
}

// AddCommits adds the "commits" edges to the GitCommit entity.
func (ruo *RepoUpdateOne) AddCommits(g ...*GitCommit) *RepoUpdateOne {
	ids := make([]int, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return ruo.AddCommitIDs(ids...)
}

// AddVulnerabilityReviewIDs adds the "vulnerability_reviews" edge to the VulnerabilityReview entity by IDs.
func (ruo *RepoUpdateOne) AddVulnerabilityReviewIDs(ids ...int) *RepoUpdateOne {
	ruo.mutation.AddVulnerabilityReviewIDs(ids...)
	return ruo
}

// AddVulnerabilityReviews adds the "vulnerability_reviews" edges to the VulnerabilityReview entity.
func (ruo *RepoUpdateOne) AddVulnerabilityReviews(v ...*VulnerabilityReview) *RepoUpdateOne {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return ruo.AddVulnerabilityReviewIDs(ids...)
}

// AddPolicyIDs adds the "policies" edge to the ReleasePolicy entity by IDs.
func (ruo *RepoUpdateOne) AddPolicyIDs(ids ...int) *RepoUpdateOne {
	ruo.mutation.AddPolicyIDs(ids...)
	return ruo
}

// AddPolicies adds the "policies" edges to the ReleasePolicy entity.
func (ruo *RepoUpdateOne) AddPolicies(r ...*ReleasePolicy) *RepoUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ruo.AddPolicyIDs(ids...)
}

// Mutation returns the RepoMutation object of the builder.
func (ruo *RepoUpdateOne) Mutation() *RepoMutation {
	return ruo.mutation
}

// ClearOwner clears the "owner" edge to the Organization entity.
func (ruo *RepoUpdateOne) ClearOwner() *RepoUpdateOne {
	ruo.mutation.ClearOwner()
	return ruo
}

// ClearProject clears the "project" edge to the Project entity.
func (ruo *RepoUpdateOne) ClearProject() *RepoUpdateOne {
	ruo.mutation.ClearProject()
	return ruo
}

// ClearHead clears the "head" edge to the Release entity.
func (ruo *RepoUpdateOne) ClearHead() *RepoUpdateOne {
	ruo.mutation.ClearHead()
	return ruo
}

// ClearCommits clears all "commits" edges to the GitCommit entity.
func (ruo *RepoUpdateOne) ClearCommits() *RepoUpdateOne {
	ruo.mutation.ClearCommits()
	return ruo
}

// RemoveCommitIDs removes the "commits" edge to GitCommit entities by IDs.
func (ruo *RepoUpdateOne) RemoveCommitIDs(ids ...int) *RepoUpdateOne {
	ruo.mutation.RemoveCommitIDs(ids...)
	return ruo
}

// RemoveCommits removes "commits" edges to GitCommit entities.
func (ruo *RepoUpdateOne) RemoveCommits(g ...*GitCommit) *RepoUpdateOne {
	ids := make([]int, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return ruo.RemoveCommitIDs(ids...)
}

// ClearVulnerabilityReviews clears all "vulnerability_reviews" edges to the VulnerabilityReview entity.
func (ruo *RepoUpdateOne) ClearVulnerabilityReviews() *RepoUpdateOne {
	ruo.mutation.ClearVulnerabilityReviews()
	return ruo
}

// RemoveVulnerabilityReviewIDs removes the "vulnerability_reviews" edge to VulnerabilityReview entities by IDs.
func (ruo *RepoUpdateOne) RemoveVulnerabilityReviewIDs(ids ...int) *RepoUpdateOne {
	ruo.mutation.RemoveVulnerabilityReviewIDs(ids...)
	return ruo
}

// RemoveVulnerabilityReviews removes "vulnerability_reviews" edges to VulnerabilityReview entities.
func (ruo *RepoUpdateOne) RemoveVulnerabilityReviews(v ...*VulnerabilityReview) *RepoUpdateOne {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return ruo.RemoveVulnerabilityReviewIDs(ids...)
}

// ClearPolicies clears all "policies" edges to the ReleasePolicy entity.
func (ruo *RepoUpdateOne) ClearPolicies() *RepoUpdateOne {
	ruo.mutation.ClearPolicies()
	return ruo
}

// RemovePolicyIDs removes the "policies" edge to ReleasePolicy entities by IDs.
func (ruo *RepoUpdateOne) RemovePolicyIDs(ids ...int) *RepoUpdateOne {
	ruo.mutation.RemovePolicyIDs(ids...)
	return ruo
}

// RemovePolicies removes "policies" edges to ReleasePolicy entities.
func (ruo *RepoUpdateOne) RemovePolicies(r ...*ReleasePolicy) *RepoUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ruo.RemovePolicyIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ruo *RepoUpdateOne) Select(field string, fields ...string) *RepoUpdateOne {
	ruo.fields = append([]string{field}, fields...)
	return ruo
}

// Save executes the query and returns the updated Repo entity.
func (ruo *RepoUpdateOne) Save(ctx context.Context) (*Repo, error) {
	var (
		err  error
		node *Repo
	)
	if len(ruo.hooks) == 0 {
		if err = ruo.check(); err != nil {
			return nil, err
		}
		node, err = ruo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RepoMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ruo.check(); err != nil {
				return nil, err
			}
			ruo.mutation = mutation
			node, err = ruo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ruo.hooks) - 1; i >= 0; i-- {
			if ruo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ruo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ruo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (ruo *RepoUpdateOne) SaveX(ctx context.Context) *Repo {
	node, err := ruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ruo *RepoUpdateOne) Exec(ctx context.Context) error {
	_, err := ruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruo *RepoUpdateOne) ExecX(ctx context.Context) {
	if err := ruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ruo *RepoUpdateOne) check() error {
	if v, ok := ruo.mutation.Name(); ok {
		if err := repo.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	if v, ok := ruo.mutation.DefaultBranch(); ok {
		if err := repo.DefaultBranchValidator(v); err != nil {
			return &ValidationError{Name: "default_branch", err: fmt.Errorf("ent: validator failed for field \"default_branch\": %w", err)}
		}
	}
	if _, ok := ruo.mutation.OwnerID(); ruo.mutation.OwnerCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"owner\"")
	}
	if _, ok := ruo.mutation.ProjectID(); ruo.mutation.ProjectCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"project\"")
	}
	return nil
}

func (ruo *RepoUpdateOne) sqlSave(ctx context.Context) (_node *Repo, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   repo.Table,
			Columns: repo.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: repo.FieldID,
			},
		},
	}
	id, ok := ruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Repo.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := ruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, repo.FieldID)
		for _, f := range fields {
			if !repo.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != repo.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ruo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: repo.FieldName,
		})
	}
	if value, ok := ruo.mutation.DefaultBranch(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: repo.FieldDefaultBranch,
		})
	}
	if value, ok := ruo.mutation.Labels(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: repo.FieldLabels,
		})
	}
	if ruo.mutation.LabelsCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: repo.FieldLabels,
		})
	}
	if ruo.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   repo.OwnerTable,
			Columns: []string{repo.OwnerColumn},
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
	if nodes := ruo.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   repo.OwnerTable,
			Columns: []string{repo.OwnerColumn},
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
	if ruo.mutation.ProjectCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   repo.ProjectTable,
			Columns: []string{repo.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   repo.ProjectTable,
			Columns: []string{repo.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ruo.mutation.HeadCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   repo.HeadTable,
			Columns: []string{repo.HeadColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: release.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.HeadIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   repo.HeadTable,
			Columns: []string{repo.HeadColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ruo.mutation.CommitsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   repo.CommitsTable,
			Columns: []string{repo.CommitsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: gitcommit.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.RemovedCommitsIDs(); len(nodes) > 0 && !ruo.mutation.CommitsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   repo.CommitsTable,
			Columns: []string{repo.CommitsColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.CommitsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   repo.CommitsTable,
			Columns: []string{repo.CommitsColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ruo.mutation.VulnerabilityReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.VulnerabilityReviewsTable,
			Columns: repo.VulnerabilityReviewsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: vulnerabilityreview.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.RemovedVulnerabilityReviewsIDs(); len(nodes) > 0 && !ruo.mutation.VulnerabilityReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.VulnerabilityReviewsTable,
			Columns: repo.VulnerabilityReviewsPrimaryKey,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.VulnerabilityReviewsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.VulnerabilityReviewsTable,
			Columns: repo.VulnerabilityReviewsPrimaryKey,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ruo.mutation.PoliciesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.PoliciesTable,
			Columns: repo.PoliciesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releasepolicy.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.RemovedPoliciesIDs(); len(nodes) > 0 && !ruo.mutation.PoliciesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.PoliciesTable,
			Columns: repo.PoliciesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releasepolicy.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.PoliciesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   repo.PoliciesTable,
			Columns: repo.PoliciesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: releasepolicy.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Repo{config: ruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{repo.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
