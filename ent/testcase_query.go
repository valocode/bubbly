// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/testcase"
	"github.com/valocode/bubbly/ent/testrun"
)

// TestCaseQuery is the builder for querying TestCase entities.
type TestCaseQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.TestCase
	// eager-loading edges.
	withRun *TestRunQuery
	withFKs bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the TestCaseQuery builder.
func (tcq *TestCaseQuery) Where(ps ...predicate.TestCase) *TestCaseQuery {
	tcq.predicates = append(tcq.predicates, ps...)
	return tcq
}

// Limit adds a limit step to the query.
func (tcq *TestCaseQuery) Limit(limit int) *TestCaseQuery {
	tcq.limit = &limit
	return tcq
}

// Offset adds an offset step to the query.
func (tcq *TestCaseQuery) Offset(offset int) *TestCaseQuery {
	tcq.offset = &offset
	return tcq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (tcq *TestCaseQuery) Unique(unique bool) *TestCaseQuery {
	tcq.unique = &unique
	return tcq
}

// Order adds an order step to the query.
func (tcq *TestCaseQuery) Order(o ...OrderFunc) *TestCaseQuery {
	tcq.order = append(tcq.order, o...)
	return tcq
}

// QueryRun chains the current query on the "run" edge.
func (tcq *TestCaseQuery) QueryRun() *TestRunQuery {
	query := &TestRunQuery{config: tcq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(testcase.Table, testcase.FieldID, selector),
			sqlgraph.To(testrun.Table, testrun.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, testcase.RunTable, testcase.RunColumn),
		)
		fromU = sqlgraph.SetNeighbors(tcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first TestCase entity from the query.
// Returns a *NotFoundError when no TestCase was found.
func (tcq *TestCaseQuery) First(ctx context.Context) (*TestCase, error) {
	nodes, err := tcq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{testcase.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (tcq *TestCaseQuery) FirstX(ctx context.Context) *TestCase {
	node, err := tcq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first TestCase ID from the query.
// Returns a *NotFoundError when no TestCase ID was found.
func (tcq *TestCaseQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tcq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{testcase.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (tcq *TestCaseQuery) FirstIDX(ctx context.Context) int {
	id, err := tcq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single TestCase entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one TestCase entity is not found.
// Returns a *NotFoundError when no TestCase entities are found.
func (tcq *TestCaseQuery) Only(ctx context.Context) (*TestCase, error) {
	nodes, err := tcq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{testcase.Label}
	default:
		return nil, &NotSingularError{testcase.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (tcq *TestCaseQuery) OnlyX(ctx context.Context) *TestCase {
	node, err := tcq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only TestCase ID in the query.
// Returns a *NotSingularError when exactly one TestCase ID is not found.
// Returns a *NotFoundError when no entities are found.
func (tcq *TestCaseQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tcq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{testcase.Label}
	default:
		err = &NotSingularError{testcase.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (tcq *TestCaseQuery) OnlyIDX(ctx context.Context) int {
	id, err := tcq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of TestCases.
func (tcq *TestCaseQuery) All(ctx context.Context) ([]*TestCase, error) {
	if err := tcq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return tcq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (tcq *TestCaseQuery) AllX(ctx context.Context) []*TestCase {
	nodes, err := tcq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of TestCase IDs.
func (tcq *TestCaseQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := tcq.Select(testcase.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (tcq *TestCaseQuery) IDsX(ctx context.Context) []int {
	ids, err := tcq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (tcq *TestCaseQuery) Count(ctx context.Context) (int, error) {
	if err := tcq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return tcq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (tcq *TestCaseQuery) CountX(ctx context.Context) int {
	count, err := tcq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (tcq *TestCaseQuery) Exist(ctx context.Context) (bool, error) {
	if err := tcq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return tcq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (tcq *TestCaseQuery) ExistX(ctx context.Context) bool {
	exist, err := tcq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the TestCaseQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (tcq *TestCaseQuery) Clone() *TestCaseQuery {
	if tcq == nil {
		return nil
	}
	return &TestCaseQuery{
		config:     tcq.config,
		limit:      tcq.limit,
		offset:     tcq.offset,
		order:      append([]OrderFunc{}, tcq.order...),
		predicates: append([]predicate.TestCase{}, tcq.predicates...),
		withRun:    tcq.withRun.Clone(),
		// clone intermediate query.
		sql:  tcq.sql.Clone(),
		path: tcq.path,
	}
}

// WithRun tells the query-builder to eager-load the nodes that are connected to
// the "run" edge. The optional arguments are used to configure the query builder of the edge.
func (tcq *TestCaseQuery) WithRun(opts ...func(*TestRunQuery)) *TestCaseQuery {
	query := &TestRunQuery{config: tcq.config}
	for _, opt := range opts {
		opt(query)
	}
	tcq.withRun = query
	return tcq
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
//	client.TestCase.Query().
//		GroupBy(testcase.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (tcq *TestCaseQuery) GroupBy(field string, fields ...string) *TestCaseGroupBy {
	group := &TestCaseGroupBy{config: tcq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := tcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return tcq.sqlQuery(ctx), nil
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
//	client.TestCase.Query().
//		Select(testcase.FieldName).
//		Scan(ctx, &v)
//
func (tcq *TestCaseQuery) Select(field string, fields ...string) *TestCaseSelect {
	tcq.fields = append([]string{field}, fields...)
	return &TestCaseSelect{TestCaseQuery: tcq}
}

func (tcq *TestCaseQuery) prepareQuery(ctx context.Context) error {
	for _, f := range tcq.fields {
		if !testcase.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if tcq.path != nil {
		prev, err := tcq.path(ctx)
		if err != nil {
			return err
		}
		tcq.sql = prev
	}
	return nil
}

func (tcq *TestCaseQuery) sqlAll(ctx context.Context) ([]*TestCase, error) {
	var (
		nodes       = []*TestCase{}
		withFKs     = tcq.withFKs
		_spec       = tcq.querySpec()
		loadedTypes = [1]bool{
			tcq.withRun != nil,
		}
	)
	if tcq.withRun != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, testcase.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &TestCase{config: tcq.config}
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
	if err := sqlgraph.QueryNodes(ctx, tcq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := tcq.withRun; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*TestCase)
		for i := range nodes {
			if nodes[i].test_case_run == nil {
				continue
			}
			fk := *nodes[i].test_case_run
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(testrun.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "test_case_run" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Run = n
			}
		}
	}

	return nodes, nil
}

func (tcq *TestCaseQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := tcq.querySpec()
	return sqlgraph.CountNodes(ctx, tcq.driver, _spec)
}

func (tcq *TestCaseQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := tcq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (tcq *TestCaseQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   testcase.Table,
			Columns: testcase.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: testcase.FieldID,
			},
		},
		From:   tcq.sql,
		Unique: true,
	}
	if unique := tcq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := tcq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, testcase.FieldID)
		for i := range fields {
			if fields[i] != testcase.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := tcq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := tcq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := tcq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := tcq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (tcq *TestCaseQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(tcq.driver.Dialect())
	t1 := builder.Table(testcase.Table)
	columns := tcq.fields
	if len(columns) == 0 {
		columns = testcase.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if tcq.sql != nil {
		selector = tcq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	for _, p := range tcq.predicates {
		p(selector)
	}
	for _, p := range tcq.order {
		p(selector)
	}
	if offset := tcq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := tcq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// TestCaseGroupBy is the group-by builder for TestCase entities.
type TestCaseGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (tcgb *TestCaseGroupBy) Aggregate(fns ...AggregateFunc) *TestCaseGroupBy {
	tcgb.fns = append(tcgb.fns, fns...)
	return tcgb
}

// Scan applies the group-by query and scans the result into the given value.
func (tcgb *TestCaseGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := tcgb.path(ctx)
	if err != nil {
		return err
	}
	tcgb.sql = query
	return tcgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (tcgb *TestCaseGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := tcgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (tcgb *TestCaseGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(tcgb.fields) > 1 {
		return nil, errors.New("ent: TestCaseGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := tcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (tcgb *TestCaseGroupBy) StringsX(ctx context.Context) []string {
	v, err := tcgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (tcgb *TestCaseGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = tcgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{testcase.Label}
	default:
		err = fmt.Errorf("ent: TestCaseGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (tcgb *TestCaseGroupBy) StringX(ctx context.Context) string {
	v, err := tcgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (tcgb *TestCaseGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(tcgb.fields) > 1 {
		return nil, errors.New("ent: TestCaseGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := tcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (tcgb *TestCaseGroupBy) IntsX(ctx context.Context) []int {
	v, err := tcgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (tcgb *TestCaseGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = tcgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{testcase.Label}
	default:
		err = fmt.Errorf("ent: TestCaseGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (tcgb *TestCaseGroupBy) IntX(ctx context.Context) int {
	v, err := tcgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (tcgb *TestCaseGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(tcgb.fields) > 1 {
		return nil, errors.New("ent: TestCaseGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := tcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (tcgb *TestCaseGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := tcgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (tcgb *TestCaseGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = tcgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{testcase.Label}
	default:
		err = fmt.Errorf("ent: TestCaseGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (tcgb *TestCaseGroupBy) Float64X(ctx context.Context) float64 {
	v, err := tcgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (tcgb *TestCaseGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(tcgb.fields) > 1 {
		return nil, errors.New("ent: TestCaseGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := tcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (tcgb *TestCaseGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := tcgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (tcgb *TestCaseGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = tcgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{testcase.Label}
	default:
		err = fmt.Errorf("ent: TestCaseGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (tcgb *TestCaseGroupBy) BoolX(ctx context.Context) bool {
	v, err := tcgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (tcgb *TestCaseGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range tcgb.fields {
		if !testcase.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := tcgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tcgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (tcgb *TestCaseGroupBy) sqlQuery() *sql.Selector {
	selector := tcgb.sql.Select()
	aggregation := make([]string, 0, len(tcgb.fns))
	for _, fn := range tcgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(tcgb.fields)+len(tcgb.fns))
		for _, f := range tcgb.fields {
			columns = append(columns, selector.C(f))
		}
		for _, c := range aggregation {
			columns = append(columns, c)
		}
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(tcgb.fields...)...)
}

// TestCaseSelect is the builder for selecting fields of TestCase entities.
type TestCaseSelect struct {
	*TestCaseQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (tcs *TestCaseSelect) Scan(ctx context.Context, v interface{}) error {
	if err := tcs.prepareQuery(ctx); err != nil {
		return err
	}
	tcs.sql = tcs.TestCaseQuery.sqlQuery(ctx)
	return tcs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (tcs *TestCaseSelect) ScanX(ctx context.Context, v interface{}) {
	if err := tcs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (tcs *TestCaseSelect) Strings(ctx context.Context) ([]string, error) {
	if len(tcs.fields) > 1 {
		return nil, errors.New("ent: TestCaseSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := tcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (tcs *TestCaseSelect) StringsX(ctx context.Context) []string {
	v, err := tcs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (tcs *TestCaseSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = tcs.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{testcase.Label}
	default:
		err = fmt.Errorf("ent: TestCaseSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (tcs *TestCaseSelect) StringX(ctx context.Context) string {
	v, err := tcs.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (tcs *TestCaseSelect) Ints(ctx context.Context) ([]int, error) {
	if len(tcs.fields) > 1 {
		return nil, errors.New("ent: TestCaseSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := tcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (tcs *TestCaseSelect) IntsX(ctx context.Context) []int {
	v, err := tcs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (tcs *TestCaseSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = tcs.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{testcase.Label}
	default:
		err = fmt.Errorf("ent: TestCaseSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (tcs *TestCaseSelect) IntX(ctx context.Context) int {
	v, err := tcs.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (tcs *TestCaseSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(tcs.fields) > 1 {
		return nil, errors.New("ent: TestCaseSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := tcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (tcs *TestCaseSelect) Float64sX(ctx context.Context) []float64 {
	v, err := tcs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (tcs *TestCaseSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = tcs.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{testcase.Label}
	default:
		err = fmt.Errorf("ent: TestCaseSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (tcs *TestCaseSelect) Float64X(ctx context.Context) float64 {
	v, err := tcs.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (tcs *TestCaseSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(tcs.fields) > 1 {
		return nil, errors.New("ent: TestCaseSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := tcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (tcs *TestCaseSelect) BoolsX(ctx context.Context) []bool {
	v, err := tcs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (tcs *TestCaseSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = tcs.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{testcase.Label}
	default:
		err = fmt.Errorf("ent: TestCaseSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (tcs *TestCaseSelect) BoolX(ctx context.Context) bool {
	v, err := tcs.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (tcs *TestCaseSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := tcs.sql.Query()
	if err := tcs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
