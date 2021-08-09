// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/cve"
	"github.com/valocode/bubbly/ent/cverule"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/repo"
)

// CVERuleQuery is the builder for querying CVERule entities.
type CVERuleQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.CVERule
	// eager-loading edges.
	withCve     *CVEQuery
	withProject *ProjectQuery
	withRepo    *RepoQuery
	withFKs     bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CVERuleQuery builder.
func (crq *CVERuleQuery) Where(ps ...predicate.CVERule) *CVERuleQuery {
	crq.predicates = append(crq.predicates, ps...)
	return crq
}

// Limit adds a limit step to the query.
func (crq *CVERuleQuery) Limit(limit int) *CVERuleQuery {
	crq.limit = &limit
	return crq
}

// Offset adds an offset step to the query.
func (crq *CVERuleQuery) Offset(offset int) *CVERuleQuery {
	crq.offset = &offset
	return crq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (crq *CVERuleQuery) Unique(unique bool) *CVERuleQuery {
	crq.unique = &unique
	return crq
}

// Order adds an order step to the query.
func (crq *CVERuleQuery) Order(o ...OrderFunc) *CVERuleQuery {
	crq.order = append(crq.order, o...)
	return crq
}

// QueryCve chains the current query on the "cve" edge.
func (crq *CVERuleQuery) QueryCve() *CVEQuery {
	query := &CVEQuery{config: crq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := crq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := crq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(cverule.Table, cverule.FieldID, selector),
			sqlgraph.To(cve.Table, cve.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, cverule.CveTable, cverule.CveColumn),
		)
		fromU = sqlgraph.SetNeighbors(crq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryProject chains the current query on the "project" edge.
func (crq *CVERuleQuery) QueryProject() *ProjectQuery {
	query := &ProjectQuery{config: crq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := crq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := crq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(cverule.Table, cverule.FieldID, selector),
			sqlgraph.To(project.Table, project.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, cverule.ProjectTable, cverule.ProjectPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(crq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRepo chains the current query on the "repo" edge.
func (crq *CVERuleQuery) QueryRepo() *RepoQuery {
	query := &RepoQuery{config: crq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := crq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := crq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(cverule.Table, cverule.FieldID, selector),
			sqlgraph.To(repo.Table, repo.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, cverule.RepoTable, cverule.RepoPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(crq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CVERule entity from the query.
// Returns a *NotFoundError when no CVERule was found.
func (crq *CVERuleQuery) First(ctx context.Context) (*CVERule, error) {
	nodes, err := crq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{cverule.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (crq *CVERuleQuery) FirstX(ctx context.Context) *CVERule {
	node, err := crq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first CVERule ID from the query.
// Returns a *NotFoundError when no CVERule ID was found.
func (crq *CVERuleQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = crq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{cverule.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (crq *CVERuleQuery) FirstIDX(ctx context.Context) int {
	id, err := crq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single CVERule entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one CVERule entity is not found.
// Returns a *NotFoundError when no CVERule entities are found.
func (crq *CVERuleQuery) Only(ctx context.Context) (*CVERule, error) {
	nodes, err := crq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{cverule.Label}
	default:
		return nil, &NotSingularError{cverule.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (crq *CVERuleQuery) OnlyX(ctx context.Context) *CVERule {
	node, err := crq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only CVERule ID in the query.
// Returns a *NotSingularError when exactly one CVERule ID is not found.
// Returns a *NotFoundError when no entities are found.
func (crq *CVERuleQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = crq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{cverule.Label}
	default:
		err = &NotSingularError{cverule.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (crq *CVERuleQuery) OnlyIDX(ctx context.Context) int {
	id, err := crq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CVERules.
func (crq *CVERuleQuery) All(ctx context.Context) ([]*CVERule, error) {
	if err := crq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return crq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (crq *CVERuleQuery) AllX(ctx context.Context) []*CVERule {
	nodes, err := crq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of CVERule IDs.
func (crq *CVERuleQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := crq.Select(cverule.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (crq *CVERuleQuery) IDsX(ctx context.Context) []int {
	ids, err := crq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (crq *CVERuleQuery) Count(ctx context.Context) (int, error) {
	if err := crq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return crq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (crq *CVERuleQuery) CountX(ctx context.Context) int {
	count, err := crq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (crq *CVERuleQuery) Exist(ctx context.Context) (bool, error) {
	if err := crq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return crq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (crq *CVERuleQuery) ExistX(ctx context.Context) bool {
	exist, err := crq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CVERuleQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (crq *CVERuleQuery) Clone() *CVERuleQuery {
	if crq == nil {
		return nil
	}
	return &CVERuleQuery{
		config:      crq.config,
		limit:       crq.limit,
		offset:      crq.offset,
		order:       append([]OrderFunc{}, crq.order...),
		predicates:  append([]predicate.CVERule{}, crq.predicates...),
		withCve:     crq.withCve.Clone(),
		withProject: crq.withProject.Clone(),
		withRepo:    crq.withRepo.Clone(),
		// clone intermediate query.
		sql:  crq.sql.Clone(),
		path: crq.path,
	}
}

// WithCve tells the query-builder to eager-load the nodes that are connected to
// the "cve" edge. The optional arguments are used to configure the query builder of the edge.
func (crq *CVERuleQuery) WithCve(opts ...func(*CVEQuery)) *CVERuleQuery {
	query := &CVEQuery{config: crq.config}
	for _, opt := range opts {
		opt(query)
	}
	crq.withCve = query
	return crq
}

// WithProject tells the query-builder to eager-load the nodes that are connected to
// the "project" edge. The optional arguments are used to configure the query builder of the edge.
func (crq *CVERuleQuery) WithProject(opts ...func(*ProjectQuery)) *CVERuleQuery {
	query := &ProjectQuery{config: crq.config}
	for _, opt := range opts {
		opt(query)
	}
	crq.withProject = query
	return crq
}

// WithRepo tells the query-builder to eager-load the nodes that are connected to
// the "repo" edge. The optional arguments are used to configure the query builder of the edge.
func (crq *CVERuleQuery) WithRepo(opts ...func(*RepoQuery)) *CVERuleQuery {
	query := &RepoQuery{config: crq.config}
	for _, opt := range opts {
		opt(query)
	}
	crq.withRepo = query
	return crq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.CVERule.Query().
//		GroupBy(cverule.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (crq *CVERuleQuery) GroupBy(field string, fields ...string) *CVERuleGroupBy {
	group := &CVERuleGroupBy{config: crq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := crq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return crq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.CVERule.Query().
//		Select(cverule.FieldName).
//		Scan(ctx, &v)
//
func (crq *CVERuleQuery) Select(fields ...string) *CVERuleSelect {
	crq.fields = append(crq.fields, fields...)
	return &CVERuleSelect{CVERuleQuery: crq}
}

func (crq *CVERuleQuery) prepareQuery(ctx context.Context) error {
	for _, f := range crq.fields {
		if !cverule.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if crq.path != nil {
		prev, err := crq.path(ctx)
		if err != nil {
			return err
		}
		crq.sql = prev
	}
	return nil
}

func (crq *CVERuleQuery) sqlAll(ctx context.Context) ([]*CVERule, error) {
	var (
		nodes       = []*CVERule{}
		withFKs     = crq.withFKs
		_spec       = crq.querySpec()
		loadedTypes = [3]bool{
			crq.withCve != nil,
			crq.withProject != nil,
			crq.withRepo != nil,
		}
	)
	if crq.withCve != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, cverule.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &CVERule{config: crq.config}
		nodes = append(nodes, node)
		return node.scanValues(columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if err := sqlgraph.QueryNodes(ctx, crq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := crq.withCve; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*CVERule)
		for i := range nodes {
			if nodes[i].cve_rule_cve == nil {
				continue
			}
			fk := *nodes[i].cve_rule_cve
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(cve.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "cve_rule_cve" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Cve = n
			}
		}
	}

	if query := crq.withProject; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		ids := make(map[int]*CVERule, len(nodes))
		for _, node := range nodes {
			ids[node.ID] = node
			fks = append(fks, node.ID)
			node.Edges.Project = []*Project{}
		}
		var (
			edgeids []int
			edges   = make(map[int][]*CVERule)
		)
		_spec := &sqlgraph.EdgeQuerySpec{
			Edge: &sqlgraph.EdgeSpec{
				Inverse: false,
				Table:   cverule.ProjectTable,
				Columns: cverule.ProjectPrimaryKey,
			},
			Predicate: func(s *sql.Selector) {
				s.Where(sql.InValues(cverule.ProjectPrimaryKey[0], fks...))
			},
			ScanValues: func() [2]interface{} {
				return [2]interface{}{new(sql.NullInt64), new(sql.NullInt64)}
			},
			Assign: func(out, in interface{}) error {
				eout, ok := out.(*sql.NullInt64)
				if !ok || eout == nil {
					return fmt.Errorf("unexpected id value for edge-out")
				}
				ein, ok := in.(*sql.NullInt64)
				if !ok || ein == nil {
					return fmt.Errorf("unexpected id value for edge-in")
				}
				outValue := int(eout.Int64)
				inValue := int(ein.Int64)
				node, ok := ids[outValue]
				if !ok {
					return fmt.Errorf("unexpected node id in edges: %v", outValue)
				}
				if _, ok := edges[inValue]; !ok {
					edgeids = append(edgeids, inValue)
				}
				edges[inValue] = append(edges[inValue], node)
				return nil
			},
		}
		if err := sqlgraph.QueryEdges(ctx, crq.driver, _spec); err != nil {
			return nil, fmt.Errorf(`query edges "project": %w`, err)
		}
		query.Where(project.IDIn(edgeids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := edges[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected "project" node returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Project = append(nodes[i].Edges.Project, n)
			}
		}
	}

	if query := crq.withRepo; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		ids := make(map[int]*CVERule, len(nodes))
		for _, node := range nodes {
			ids[node.ID] = node
			fks = append(fks, node.ID)
			node.Edges.Repo = []*Repo{}
		}
		var (
			edgeids []int
			edges   = make(map[int][]*CVERule)
		)
		_spec := &sqlgraph.EdgeQuerySpec{
			Edge: &sqlgraph.EdgeSpec{
				Inverse: false,
				Table:   cverule.RepoTable,
				Columns: cverule.RepoPrimaryKey,
			},
			Predicate: func(s *sql.Selector) {
				s.Where(sql.InValues(cverule.RepoPrimaryKey[0], fks...))
			},
			ScanValues: func() [2]interface{} {
				return [2]interface{}{new(sql.NullInt64), new(sql.NullInt64)}
			},
			Assign: func(out, in interface{}) error {
				eout, ok := out.(*sql.NullInt64)
				if !ok || eout == nil {
					return fmt.Errorf("unexpected id value for edge-out")
				}
				ein, ok := in.(*sql.NullInt64)
				if !ok || ein == nil {
					return fmt.Errorf("unexpected id value for edge-in")
				}
				outValue := int(eout.Int64)
				inValue := int(ein.Int64)
				node, ok := ids[outValue]
				if !ok {
					return fmt.Errorf("unexpected node id in edges: %v", outValue)
				}
				if _, ok := edges[inValue]; !ok {
					edgeids = append(edgeids, inValue)
				}
				edges[inValue] = append(edges[inValue], node)
				return nil
			},
		}
		if err := sqlgraph.QueryEdges(ctx, crq.driver, _spec); err != nil {
			return nil, fmt.Errorf(`query edges "repo": %w`, err)
		}
		query.Where(repo.IDIn(edgeids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := edges[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected "repo" node returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Repo = append(nodes[i].Edges.Repo, n)
			}
		}
	}

	return nodes, nil
}

func (crq *CVERuleQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := crq.querySpec()
	return sqlgraph.CountNodes(ctx, crq.driver, _spec)
}

func (crq *CVERuleQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := crq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (crq *CVERuleQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cverule.Table,
			Columns: cverule.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: cverule.FieldID,
			},
		},
		From:   crq.sql,
		Unique: true,
	}
	if unique := crq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := crq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cverule.FieldID)
		for i := range fields {
			if fields[i] != cverule.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := crq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := crq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := crq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := crq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (crq *CVERuleQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(crq.driver.Dialect())
	t1 := builder.Table(cverule.Table)
	columns := crq.fields
	if len(columns) == 0 {
		columns = cverule.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if crq.sql != nil {
		selector = crq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	for _, p := range crq.predicates {
		p(selector)
	}
	for _, p := range crq.order {
		p(selector)
	}
	if offset := crq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := crq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CVERuleGroupBy is the group-by builder for CVERule entities.
type CVERuleGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (crgb *CVERuleGroupBy) Aggregate(fns ...AggregateFunc) *CVERuleGroupBy {
	crgb.fns = append(crgb.fns, fns...)
	return crgb
}

// Scan applies the group-by query and scans the result into the given value.
func (crgb *CVERuleGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := crgb.path(ctx)
	if err != nil {
		return err
	}
	crgb.sql = query
	return crgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (crgb *CVERuleGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := crgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (crgb *CVERuleGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(crgb.fields) > 1 {
		return nil, errors.New("ent: CVERuleGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := crgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (crgb *CVERuleGroupBy) StringsX(ctx context.Context) []string {
	v, err := crgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (crgb *CVERuleGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = crgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cverule.Label}
	default:
		err = fmt.Errorf("ent: CVERuleGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (crgb *CVERuleGroupBy) StringX(ctx context.Context) string {
	v, err := crgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (crgb *CVERuleGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(crgb.fields) > 1 {
		return nil, errors.New("ent: CVERuleGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := crgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (crgb *CVERuleGroupBy) IntsX(ctx context.Context) []int {
	v, err := crgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (crgb *CVERuleGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = crgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cverule.Label}
	default:
		err = fmt.Errorf("ent: CVERuleGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (crgb *CVERuleGroupBy) IntX(ctx context.Context) int {
	v, err := crgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (crgb *CVERuleGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(crgb.fields) > 1 {
		return nil, errors.New("ent: CVERuleGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := crgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (crgb *CVERuleGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := crgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (crgb *CVERuleGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = crgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cverule.Label}
	default:
		err = fmt.Errorf("ent: CVERuleGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (crgb *CVERuleGroupBy) Float64X(ctx context.Context) float64 {
	v, err := crgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (crgb *CVERuleGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(crgb.fields) > 1 {
		return nil, errors.New("ent: CVERuleGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := crgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (crgb *CVERuleGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := crgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (crgb *CVERuleGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = crgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cverule.Label}
	default:
		err = fmt.Errorf("ent: CVERuleGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (crgb *CVERuleGroupBy) BoolX(ctx context.Context) bool {
	v, err := crgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (crgb *CVERuleGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range crgb.fields {
		if !cverule.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := crgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := crgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (crgb *CVERuleGroupBy) sqlQuery() *sql.Selector {
	selector := crgb.sql.Select()
	aggregation := make([]string, 0, len(crgb.fns))
	for _, fn := range crgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(crgb.fields)+len(crgb.fns))
		for _, f := range crgb.fields {
			columns = append(columns, selector.C(f))
		}
		for _, c := range aggregation {
			columns = append(columns, c)
		}
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(crgb.fields...)...)
}

// CVERuleSelect is the builder for selecting fields of CVERule entities.
type CVERuleSelect struct {
	*CVERuleQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (crs *CVERuleSelect) Scan(ctx context.Context, v interface{}) error {
	if err := crs.prepareQuery(ctx); err != nil {
		return err
	}
	crs.sql = crs.CVERuleQuery.sqlQuery(ctx)
	return crs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (crs *CVERuleSelect) ScanX(ctx context.Context, v interface{}) {
	if err := crs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (crs *CVERuleSelect) Strings(ctx context.Context) ([]string, error) {
	if len(crs.fields) > 1 {
		return nil, errors.New("ent: CVERuleSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := crs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (crs *CVERuleSelect) StringsX(ctx context.Context) []string {
	v, err := crs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (crs *CVERuleSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = crs.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cverule.Label}
	default:
		err = fmt.Errorf("ent: CVERuleSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (crs *CVERuleSelect) StringX(ctx context.Context) string {
	v, err := crs.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (crs *CVERuleSelect) Ints(ctx context.Context) ([]int, error) {
	if len(crs.fields) > 1 {
		return nil, errors.New("ent: CVERuleSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := crs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (crs *CVERuleSelect) IntsX(ctx context.Context) []int {
	v, err := crs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (crs *CVERuleSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = crs.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cverule.Label}
	default:
		err = fmt.Errorf("ent: CVERuleSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (crs *CVERuleSelect) IntX(ctx context.Context) int {
	v, err := crs.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (crs *CVERuleSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(crs.fields) > 1 {
		return nil, errors.New("ent: CVERuleSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := crs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (crs *CVERuleSelect) Float64sX(ctx context.Context) []float64 {
	v, err := crs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (crs *CVERuleSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = crs.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cverule.Label}
	default:
		err = fmt.Errorf("ent: CVERuleSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (crs *CVERuleSelect) Float64X(ctx context.Context) float64 {
	v, err := crs.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (crs *CVERuleSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(crs.fields) > 1 {
		return nil, errors.New("ent: CVERuleSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := crs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (crs *CVERuleSelect) BoolsX(ctx context.Context) []bool {
	v, err := crs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (crs *CVERuleSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = crs.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cverule.Label}
	default:
		err = fmt.Errorf("ent: CVERuleSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (crs *CVERuleSelect) BoolX(ctx context.Context) bool {
	v, err := crs.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (crs *CVERuleSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := crs.sql.Query()
	if err := crs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
