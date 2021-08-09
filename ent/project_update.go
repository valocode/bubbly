// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/cverule"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/repo"
)

// ProjectUpdate is the builder for updating Project entities.
type ProjectUpdate struct {
	config
	hooks    []Hook
	mutation *ProjectMutation
}

// Where appends a list predicates to the ProjectUpdate builder.
func (pu *ProjectUpdate) Where(ps ...predicate.Project) *ProjectUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetName sets the "name" field.
func (pu *ProjectUpdate) SetName(s string) *ProjectUpdate {
	pu.mutation.SetName(s)
	return pu
}

// AddRepoIDs adds the "repos" edge to the Repo entity by IDs.
func (pu *ProjectUpdate) AddRepoIDs(ids ...int) *ProjectUpdate {
	pu.mutation.AddRepoIDs(ids...)
	return pu
}

// AddRepos adds the "repos" edges to the Repo entity.
func (pu *ProjectUpdate) AddRepos(r ...*Repo) *ProjectUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return pu.AddRepoIDs(ids...)
}

// AddReleaseIDs adds the "releases" edge to the Release entity by IDs.
func (pu *ProjectUpdate) AddReleaseIDs(ids ...int) *ProjectUpdate {
	pu.mutation.AddReleaseIDs(ids...)
	return pu
}

// AddReleases adds the "releases" edges to the Release entity.
func (pu *ProjectUpdate) AddReleases(r ...*Release) *ProjectUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return pu.AddReleaseIDs(ids...)
}

// AddCveRuleIDs adds the "cve_rules" edge to the CVERule entity by IDs.
func (pu *ProjectUpdate) AddCveRuleIDs(ids ...int) *ProjectUpdate {
	pu.mutation.AddCveRuleIDs(ids...)
	return pu
}

// AddCveRules adds the "cve_rules" edges to the CVERule entity.
func (pu *ProjectUpdate) AddCveRules(c ...*CVERule) *ProjectUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return pu.AddCveRuleIDs(ids...)
}

// Mutation returns the ProjectMutation object of the builder.
func (pu *ProjectUpdate) Mutation() *ProjectMutation {
	return pu.mutation
}

// ClearRepos clears all "repos" edges to the Repo entity.
func (pu *ProjectUpdate) ClearRepos() *ProjectUpdate {
	pu.mutation.ClearRepos()
	return pu
}

// RemoveRepoIDs removes the "repos" edge to Repo entities by IDs.
func (pu *ProjectUpdate) RemoveRepoIDs(ids ...int) *ProjectUpdate {
	pu.mutation.RemoveRepoIDs(ids...)
	return pu
}

// RemoveRepos removes "repos" edges to Repo entities.
func (pu *ProjectUpdate) RemoveRepos(r ...*Repo) *ProjectUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return pu.RemoveRepoIDs(ids...)
}

// ClearReleases clears all "releases" edges to the Release entity.
func (pu *ProjectUpdate) ClearReleases() *ProjectUpdate {
	pu.mutation.ClearReleases()
	return pu
}

// RemoveReleaseIDs removes the "releases" edge to Release entities by IDs.
func (pu *ProjectUpdate) RemoveReleaseIDs(ids ...int) *ProjectUpdate {
	pu.mutation.RemoveReleaseIDs(ids...)
	return pu
}

// RemoveReleases removes "releases" edges to Release entities.
func (pu *ProjectUpdate) RemoveReleases(r ...*Release) *ProjectUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return pu.RemoveReleaseIDs(ids...)
}

// ClearCveRules clears all "cve_rules" edges to the CVERule entity.
func (pu *ProjectUpdate) ClearCveRules() *ProjectUpdate {
	pu.mutation.ClearCveRules()
	return pu
}

// RemoveCveRuleIDs removes the "cve_rules" edge to CVERule entities by IDs.
func (pu *ProjectUpdate) RemoveCveRuleIDs(ids ...int) *ProjectUpdate {
	pu.mutation.RemoveCveRuleIDs(ids...)
	return pu
}

// RemoveCveRules removes "cve_rules" edges to CVERule entities.
func (pu *ProjectUpdate) RemoveCveRules(c ...*CVERule) *ProjectUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return pu.RemoveCveRuleIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *ProjectUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(pu.hooks) == 0 {
		if err = pu.check(); err != nil {
			return 0, err
		}
		affected, err = pu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProjectMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = pu.check(); err != nil {
				return 0, err
			}
			pu.mutation = mutation
			affected, err = pu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(pu.hooks) - 1; i >= 0; i-- {
			if pu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = pu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (pu *ProjectUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *ProjectUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *ProjectUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pu *ProjectUpdate) check() error {
	if v, ok := pu.mutation.Name(); ok {
		if err := project.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	return nil
}

func (pu *ProjectUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   project.Table,
			Columns: project.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: project.FieldID,
			},
		},
	}
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: project.FieldName,
		})
	}
	if pu.mutation.ReposCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReposTable,
			Columns: []string{project.ReposColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repo.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedReposIDs(); len(nodes) > 0 && !pu.mutation.ReposCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReposTable,
			Columns: []string{project.ReposColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repo.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.ReposIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReposTable,
			Columns: []string{project.ReposColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repo.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.mutation.ReleasesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReleasesTable,
			Columns: []string{project.ReleasesColumn},
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
	if nodes := pu.mutation.RemovedReleasesIDs(); len(nodes) > 0 && !pu.mutation.ReleasesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReleasesTable,
			Columns: []string{project.ReleasesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.ReleasesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReleasesTable,
			Columns: []string{project.ReleasesColumn},
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
	if pu.mutation.CveRulesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   project.CveRulesTable,
			Columns: project.CveRulesPrimaryKey,
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
	if nodes := pu.mutation.RemovedCveRulesIDs(); len(nodes) > 0 && !pu.mutation.CveRulesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   project.CveRulesTable,
			Columns: project.CveRulesPrimaryKey,
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
	if nodes := pu.mutation.CveRulesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   project.CveRulesTable,
			Columns: project.CveRulesPrimaryKey,
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
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{project.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// ProjectUpdateOne is the builder for updating a single Project entity.
type ProjectUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ProjectMutation
}

// SetName sets the "name" field.
func (puo *ProjectUpdateOne) SetName(s string) *ProjectUpdateOne {
	puo.mutation.SetName(s)
	return puo
}

// AddRepoIDs adds the "repos" edge to the Repo entity by IDs.
func (puo *ProjectUpdateOne) AddRepoIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.AddRepoIDs(ids...)
	return puo
}

// AddRepos adds the "repos" edges to the Repo entity.
func (puo *ProjectUpdateOne) AddRepos(r ...*Repo) *ProjectUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return puo.AddRepoIDs(ids...)
}

// AddReleaseIDs adds the "releases" edge to the Release entity by IDs.
func (puo *ProjectUpdateOne) AddReleaseIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.AddReleaseIDs(ids...)
	return puo
}

// AddReleases adds the "releases" edges to the Release entity.
func (puo *ProjectUpdateOne) AddReleases(r ...*Release) *ProjectUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return puo.AddReleaseIDs(ids...)
}

// AddCveRuleIDs adds the "cve_rules" edge to the CVERule entity by IDs.
func (puo *ProjectUpdateOne) AddCveRuleIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.AddCveRuleIDs(ids...)
	return puo
}

// AddCveRules adds the "cve_rules" edges to the CVERule entity.
func (puo *ProjectUpdateOne) AddCveRules(c ...*CVERule) *ProjectUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return puo.AddCveRuleIDs(ids...)
}

// Mutation returns the ProjectMutation object of the builder.
func (puo *ProjectUpdateOne) Mutation() *ProjectMutation {
	return puo.mutation
}

// ClearRepos clears all "repos" edges to the Repo entity.
func (puo *ProjectUpdateOne) ClearRepos() *ProjectUpdateOne {
	puo.mutation.ClearRepos()
	return puo
}

// RemoveRepoIDs removes the "repos" edge to Repo entities by IDs.
func (puo *ProjectUpdateOne) RemoveRepoIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.RemoveRepoIDs(ids...)
	return puo
}

// RemoveRepos removes "repos" edges to Repo entities.
func (puo *ProjectUpdateOne) RemoveRepos(r ...*Repo) *ProjectUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return puo.RemoveRepoIDs(ids...)
}

// ClearReleases clears all "releases" edges to the Release entity.
func (puo *ProjectUpdateOne) ClearReleases() *ProjectUpdateOne {
	puo.mutation.ClearReleases()
	return puo
}

// RemoveReleaseIDs removes the "releases" edge to Release entities by IDs.
func (puo *ProjectUpdateOne) RemoveReleaseIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.RemoveReleaseIDs(ids...)
	return puo
}

// RemoveReleases removes "releases" edges to Release entities.
func (puo *ProjectUpdateOne) RemoveReleases(r ...*Release) *ProjectUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return puo.RemoveReleaseIDs(ids...)
}

// ClearCveRules clears all "cve_rules" edges to the CVERule entity.
func (puo *ProjectUpdateOne) ClearCveRules() *ProjectUpdateOne {
	puo.mutation.ClearCveRules()
	return puo
}

// RemoveCveRuleIDs removes the "cve_rules" edge to CVERule entities by IDs.
func (puo *ProjectUpdateOne) RemoveCveRuleIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.RemoveCveRuleIDs(ids...)
	return puo
}

// RemoveCveRules removes "cve_rules" edges to CVERule entities.
func (puo *ProjectUpdateOne) RemoveCveRules(c ...*CVERule) *ProjectUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return puo.RemoveCveRuleIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *ProjectUpdateOne) Select(field string, fields ...string) *ProjectUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Project entity.
func (puo *ProjectUpdateOne) Save(ctx context.Context) (*Project, error) {
	var (
		err  error
		node *Project
	)
	if len(puo.hooks) == 0 {
		if err = puo.check(); err != nil {
			return nil, err
		}
		node, err = puo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProjectMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = puo.check(); err != nil {
				return nil, err
			}
			puo.mutation = mutation
			node, err = puo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(puo.hooks) - 1; i >= 0; i-- {
			if puo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = puo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, puo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (puo *ProjectUpdateOne) SaveX(ctx context.Context) *Project {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *ProjectUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *ProjectUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (puo *ProjectUpdateOne) check() error {
	if v, ok := puo.mutation.Name(); ok {
		if err := project.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	return nil
}

func (puo *ProjectUpdateOne) sqlSave(ctx context.Context) (_node *Project, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   project.Table,
			Columns: project.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: project.FieldID,
			},
		},
	}
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Project.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, project.FieldID)
		for _, f := range fields {
			if !project.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != project.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := puo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: project.FieldName,
		})
	}
	if puo.mutation.ReposCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReposTable,
			Columns: []string{project.ReposColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repo.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedReposIDs(); len(nodes) > 0 && !puo.mutation.ReposCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReposTable,
			Columns: []string{project.ReposColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repo.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.ReposIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReposTable,
			Columns: []string{project.ReposColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repo.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.mutation.ReleasesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReleasesTable,
			Columns: []string{project.ReleasesColumn},
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
	if nodes := puo.mutation.RemovedReleasesIDs(); len(nodes) > 0 && !puo.mutation.ReleasesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReleasesTable,
			Columns: []string{project.ReleasesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.ReleasesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ReleasesTable,
			Columns: []string{project.ReleasesColumn},
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
	if puo.mutation.CveRulesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   project.CveRulesTable,
			Columns: project.CveRulesPrimaryKey,
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
	if nodes := puo.mutation.RemovedCveRulesIDs(); len(nodes) > 0 && !puo.mutation.CveRulesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   project.CveRulesTable,
			Columns: project.CveRulesPrimaryKey,
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
	if nodes := puo.mutation.CveRulesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   project.CveRulesTable,
			Columns: project.CveRulesPrimaryKey,
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
	_node = &Project{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{project.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
