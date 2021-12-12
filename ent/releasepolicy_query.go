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
	"github.com/valocode/bubbly/ent/organization"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/releasepolicy"
	"github.com/valocode/bubbly/ent/releasepolicyviolation"
)

// ReleasePolicyQuery is the builder for querying ReleasePolicy entities.
type ReleasePolicyQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.ReleasePolicy
	// eager-loading edges.
	withOwner      *OrganizationQuery
	withViolations *ReleasePolicyViolationQuery
	withFKs        bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ReleasePolicyQuery builder.
func (rpq *ReleasePolicyQuery) Where(ps ...predicate.ReleasePolicy) *ReleasePolicyQuery {
	rpq.predicates = append(rpq.predicates, ps...)
	return rpq
}

// Limit adds a limit step to the query.
func (rpq *ReleasePolicyQuery) Limit(limit int) *ReleasePolicyQuery {
	rpq.limit = &limit
	return rpq
}

// Offset adds an offset step to the query.
func (rpq *ReleasePolicyQuery) Offset(offset int) *ReleasePolicyQuery {
	rpq.offset = &offset
	return rpq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rpq *ReleasePolicyQuery) Unique(unique bool) *ReleasePolicyQuery {
	rpq.unique = &unique
	return rpq
}

// Order adds an order step to the query.
func (rpq *ReleasePolicyQuery) Order(o ...OrderFunc) *ReleasePolicyQuery {
	rpq.order = append(rpq.order, o...)
	return rpq
}

// QueryOwner chains the current query on the "owner" edge.
func (rpq *ReleasePolicyQuery) QueryOwner() *OrganizationQuery {
	query := &OrganizationQuery{config: rpq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rpq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(releasepolicy.Table, releasepolicy.FieldID, selector),
			sqlgraph.To(organization.Table, organization.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, releasepolicy.OwnerTable, releasepolicy.OwnerColumn),
		)
		fromU = sqlgraph.SetNeighbors(rpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryViolations chains the current query on the "violations" edge.
func (rpq *ReleasePolicyQuery) QueryViolations() *ReleasePolicyViolationQuery {
	query := &ReleasePolicyViolationQuery{config: rpq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rpq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(releasepolicy.Table, releasepolicy.FieldID, selector),
			sqlgraph.To(releasepolicyviolation.Table, releasepolicyviolation.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, releasepolicy.ViolationsTable, releasepolicy.ViolationsColumn),
		)
		fromU = sqlgraph.SetNeighbors(rpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ReleasePolicy entity from the query.
// Returns a *NotFoundError when no ReleasePolicy was found.
func (rpq *ReleasePolicyQuery) First(ctx context.Context) (*ReleasePolicy, error) {
	nodes, err := rpq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{releasepolicy.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rpq *ReleasePolicyQuery) FirstX(ctx context.Context) *ReleasePolicy {
	node, err := rpq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ReleasePolicy ID from the query.
// Returns a *NotFoundError when no ReleasePolicy ID was found.
func (rpq *ReleasePolicyQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rpq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{releasepolicy.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rpq *ReleasePolicyQuery) FirstIDX(ctx context.Context) int {
	id, err := rpq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ReleasePolicy entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one ReleasePolicy entity is not found.
// Returns a *NotFoundError when no ReleasePolicy entities are found.
func (rpq *ReleasePolicyQuery) Only(ctx context.Context) (*ReleasePolicy, error) {
	nodes, err := rpq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{releasepolicy.Label}
	default:
		return nil, &NotSingularError{releasepolicy.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rpq *ReleasePolicyQuery) OnlyX(ctx context.Context) *ReleasePolicy {
	node, err := rpq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ReleasePolicy ID in the query.
// Returns a *NotSingularError when exactly one ReleasePolicy ID is not found.
// Returns a *NotFoundError when no entities are found.
func (rpq *ReleasePolicyQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rpq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{releasepolicy.Label}
	default:
		err = &NotSingularError{releasepolicy.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rpq *ReleasePolicyQuery) OnlyIDX(ctx context.Context) int {
	id, err := rpq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ReleasePolicies.
func (rpq *ReleasePolicyQuery) All(ctx context.Context) ([]*ReleasePolicy, error) {
	if err := rpq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return rpq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (rpq *ReleasePolicyQuery) AllX(ctx context.Context) []*ReleasePolicy {
	nodes, err := rpq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ReleasePolicy IDs.
func (rpq *ReleasePolicyQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := rpq.Select(releasepolicy.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rpq *ReleasePolicyQuery) IDsX(ctx context.Context) []int {
	ids, err := rpq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rpq *ReleasePolicyQuery) Count(ctx context.Context) (int, error) {
	if err := rpq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return rpq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (rpq *ReleasePolicyQuery) CountX(ctx context.Context) int {
	count, err := rpq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rpq *ReleasePolicyQuery) Exist(ctx context.Context) (bool, error) {
	if err := rpq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return rpq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (rpq *ReleasePolicyQuery) ExistX(ctx context.Context) bool {
	exist, err := rpq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ReleasePolicyQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rpq *ReleasePolicyQuery) Clone() *ReleasePolicyQuery {
	if rpq == nil {
		return nil
	}
	return &ReleasePolicyQuery{
		config:         rpq.config,
		limit:          rpq.limit,
		offset:         rpq.offset,
		order:          append([]OrderFunc{}, rpq.order...),
		predicates:     append([]predicate.ReleasePolicy{}, rpq.predicates...),
		withOwner:      rpq.withOwner.Clone(),
		withViolations: rpq.withViolations.Clone(),
		// clone intermediate query.
		sql:  rpq.sql.Clone(),
		path: rpq.path,
	}
}

// WithOwner tells the query-builder to eager-load the nodes that are connected to
// the "owner" edge. The optional arguments are used to configure the query builder of the edge.
func (rpq *ReleasePolicyQuery) WithOwner(opts ...func(*OrganizationQuery)) *ReleasePolicyQuery {
	query := &OrganizationQuery{config: rpq.config}
	for _, opt := range opts {
		opt(query)
	}
	rpq.withOwner = query
	return rpq
}

// WithViolations tells the query-builder to eager-load the nodes that are connected to
// the "violations" edge. The optional arguments are used to configure the query builder of the edge.
func (rpq *ReleasePolicyQuery) WithViolations(opts ...func(*ReleasePolicyViolationQuery)) *ReleasePolicyQuery {
	query := &ReleasePolicyViolationQuery{config: rpq.config}
	for _, opt := range opts {
		opt(query)
	}
	rpq.withViolations = query
	return rpq
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
//	client.ReleasePolicy.Query().
//		GroupBy(releasepolicy.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (rpq *ReleasePolicyQuery) GroupBy(field string, fields ...string) *ReleasePolicyGroupBy {
	group := &ReleasePolicyGroupBy{config: rpq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := rpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return rpq.sqlQuery(ctx), nil
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
//	client.ReleasePolicy.Query().
//		Select(releasepolicy.FieldName).
//		Scan(ctx, &v)
//
func (rpq *ReleasePolicyQuery) Select(fields ...string) *ReleasePolicySelect {
	rpq.fields = append(rpq.fields, fields...)
	return &ReleasePolicySelect{ReleasePolicyQuery: rpq}
}

func (rpq *ReleasePolicyQuery) prepareQuery(ctx context.Context) error {
	for _, f := range rpq.fields {
		if !releasepolicy.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if rpq.path != nil {
		prev, err := rpq.path(ctx)
		if err != nil {
			return err
		}
		rpq.sql = prev
	}
	return nil
}

func (rpq *ReleasePolicyQuery) sqlAll(ctx context.Context) ([]*ReleasePolicy, error) {
	var (
		nodes       = []*ReleasePolicy{}
		withFKs     = rpq.withFKs
		_spec       = rpq.querySpec()
		loadedTypes = [2]bool{
			rpq.withOwner != nil,
			rpq.withViolations != nil,
		}
	)
	if rpq.withOwner != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, releasepolicy.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &ReleasePolicy{config: rpq.config}
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
	if err := sqlgraph.QueryNodes(ctx, rpq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := rpq.withOwner; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*ReleasePolicy)
		for i := range nodes {
			if nodes[i].release_policy_owner == nil {
				continue
			}
			fk := *nodes[i].release_policy_owner
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(organization.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "release_policy_owner" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Owner = n
			}
		}
	}

	if query := rpq.withViolations; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*ReleasePolicy)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
			nodes[i].Edges.Violations = []*ReleasePolicyViolation{}
		}
		query.withFKs = true
		query.Where(predicate.ReleasePolicyViolation(func(s *sql.Selector) {
			s.Where(sql.InValues(releasepolicy.ViolationsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.release_policy_violation_policy
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "release_policy_violation_policy" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "release_policy_violation_policy" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Violations = append(node.Edges.Violations, n)
		}
	}

	return nodes, nil
}

func (rpq *ReleasePolicyQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rpq.querySpec()
	return sqlgraph.CountNodes(ctx, rpq.driver, _spec)
}

func (rpq *ReleasePolicyQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := rpq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (rpq *ReleasePolicyQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   releasepolicy.Table,
			Columns: releasepolicy.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: releasepolicy.FieldID,
			},
		},
		From:   rpq.sql,
		Unique: true,
	}
	if unique := rpq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := rpq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, releasepolicy.FieldID)
		for i := range fields {
			if fields[i] != releasepolicy.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := rpq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rpq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rpq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rpq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rpq *ReleasePolicyQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rpq.driver.Dialect())
	t1 := builder.Table(releasepolicy.Table)
	columns := rpq.fields
	if len(columns) == 0 {
		columns = releasepolicy.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if rpq.sql != nil {
		selector = rpq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	for _, p := range rpq.predicates {
		p(selector)
	}
	for _, p := range rpq.order {
		p(selector)
	}
	if offset := rpq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rpq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ReleasePolicyGroupBy is the group-by builder for ReleasePolicy entities.
type ReleasePolicyGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rpgb *ReleasePolicyGroupBy) Aggregate(fns ...AggregateFunc) *ReleasePolicyGroupBy {
	rpgb.fns = append(rpgb.fns, fns...)
	return rpgb
}

// Scan applies the group-by query and scans the result into the given value.
func (rpgb *ReleasePolicyGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := rpgb.path(ctx)
	if err != nil {
		return err
	}
	rpgb.sql = query
	return rpgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (rpgb *ReleasePolicyGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := rpgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (rpgb *ReleasePolicyGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(rpgb.fields) > 1 {
		return nil, errors.New("ent: ReleasePolicyGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := rpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (rpgb *ReleasePolicyGroupBy) StringsX(ctx context.Context) []string {
	v, err := rpgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rpgb *ReleasePolicyGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = rpgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasepolicy.Label}
	default:
		err = fmt.Errorf("ent: ReleasePolicyGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (rpgb *ReleasePolicyGroupBy) StringX(ctx context.Context) string {
	v, err := rpgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (rpgb *ReleasePolicyGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(rpgb.fields) > 1 {
		return nil, errors.New("ent: ReleasePolicyGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := rpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (rpgb *ReleasePolicyGroupBy) IntsX(ctx context.Context) []int {
	v, err := rpgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rpgb *ReleasePolicyGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = rpgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasepolicy.Label}
	default:
		err = fmt.Errorf("ent: ReleasePolicyGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (rpgb *ReleasePolicyGroupBy) IntX(ctx context.Context) int {
	v, err := rpgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (rpgb *ReleasePolicyGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(rpgb.fields) > 1 {
		return nil, errors.New("ent: ReleasePolicyGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := rpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (rpgb *ReleasePolicyGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := rpgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rpgb *ReleasePolicyGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = rpgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasepolicy.Label}
	default:
		err = fmt.Errorf("ent: ReleasePolicyGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (rpgb *ReleasePolicyGroupBy) Float64X(ctx context.Context) float64 {
	v, err := rpgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (rpgb *ReleasePolicyGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(rpgb.fields) > 1 {
		return nil, errors.New("ent: ReleasePolicyGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := rpgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (rpgb *ReleasePolicyGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := rpgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rpgb *ReleasePolicyGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = rpgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasepolicy.Label}
	default:
		err = fmt.Errorf("ent: ReleasePolicyGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (rpgb *ReleasePolicyGroupBy) BoolX(ctx context.Context) bool {
	v, err := rpgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (rpgb *ReleasePolicyGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range rpgb.fields {
		if !releasepolicy.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := rpgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rpgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (rpgb *ReleasePolicyGroupBy) sqlQuery() *sql.Selector {
	selector := rpgb.sql.Select()
	aggregation := make([]string, 0, len(rpgb.fns))
	for _, fn := range rpgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(rpgb.fields)+len(rpgb.fns))
		for _, f := range rpgb.fields {
			columns = append(columns, selector.C(f))
		}
		for _, c := range aggregation {
			columns = append(columns, c)
		}
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(rpgb.fields...)...)
}

// ReleasePolicySelect is the builder for selecting fields of ReleasePolicy entities.
type ReleasePolicySelect struct {
	*ReleasePolicyQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (rps *ReleasePolicySelect) Scan(ctx context.Context, v interface{}) error {
	if err := rps.prepareQuery(ctx); err != nil {
		return err
	}
	rps.sql = rps.ReleasePolicyQuery.sqlQuery(ctx)
	return rps.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (rps *ReleasePolicySelect) ScanX(ctx context.Context, v interface{}) {
	if err := rps.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (rps *ReleasePolicySelect) Strings(ctx context.Context) ([]string, error) {
	if len(rps.fields) > 1 {
		return nil, errors.New("ent: ReleasePolicySelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := rps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (rps *ReleasePolicySelect) StringsX(ctx context.Context) []string {
	v, err := rps.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (rps *ReleasePolicySelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = rps.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasepolicy.Label}
	default:
		err = fmt.Errorf("ent: ReleasePolicySelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (rps *ReleasePolicySelect) StringX(ctx context.Context) string {
	v, err := rps.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (rps *ReleasePolicySelect) Ints(ctx context.Context) ([]int, error) {
	if len(rps.fields) > 1 {
		return nil, errors.New("ent: ReleasePolicySelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := rps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (rps *ReleasePolicySelect) IntsX(ctx context.Context) []int {
	v, err := rps.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (rps *ReleasePolicySelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = rps.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasepolicy.Label}
	default:
		err = fmt.Errorf("ent: ReleasePolicySelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (rps *ReleasePolicySelect) IntX(ctx context.Context) int {
	v, err := rps.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (rps *ReleasePolicySelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(rps.fields) > 1 {
		return nil, errors.New("ent: ReleasePolicySelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := rps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (rps *ReleasePolicySelect) Float64sX(ctx context.Context) []float64 {
	v, err := rps.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (rps *ReleasePolicySelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = rps.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasepolicy.Label}
	default:
		err = fmt.Errorf("ent: ReleasePolicySelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (rps *ReleasePolicySelect) Float64X(ctx context.Context) float64 {
	v, err := rps.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (rps *ReleasePolicySelect) Bools(ctx context.Context) ([]bool, error) {
	if len(rps.fields) > 1 {
		return nil, errors.New("ent: ReleasePolicySelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := rps.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (rps *ReleasePolicySelect) BoolsX(ctx context.Context) []bool {
	v, err := rps.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (rps *ReleasePolicySelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = rps.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasepolicy.Label}
	default:
		err = fmt.Errorf("ent: ReleasePolicySelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (rps *ReleasePolicySelect) BoolX(ctx context.Context) bool {
	v, err := rps.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (rps *ReleasePolicySelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := rps.sql.Query()
	if err := rps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
