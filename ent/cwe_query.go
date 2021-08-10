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
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/cwe"
	"github.com/valocode/bubbly/ent/predicate"
)

// CWEQuery is the builder for querying CWE entities.
type CWEQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.CWE
	// eager-loading edges.
	withIssues *CodeIssueQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CWEQuery builder.
func (cq *CWEQuery) Where(ps ...predicate.CWE) *CWEQuery {
	cq.predicates = append(cq.predicates, ps...)
	return cq
}

// Limit adds a limit step to the query.
func (cq *CWEQuery) Limit(limit int) *CWEQuery {
	cq.limit = &limit
	return cq
}

// Offset adds an offset step to the query.
func (cq *CWEQuery) Offset(offset int) *CWEQuery {
	cq.offset = &offset
	return cq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (cq *CWEQuery) Unique(unique bool) *CWEQuery {
	cq.unique = &unique
	return cq
}

// Order adds an order step to the query.
func (cq *CWEQuery) Order(o ...OrderFunc) *CWEQuery {
	cq.order = append(cq.order, o...)
	return cq
}

// QueryIssues chains the current query on the "issues" edge.
func (cq *CWEQuery) QueryIssues() *CodeIssueQuery {
	query := &CodeIssueQuery{config: cq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := cq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := cq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(cwe.Table, cwe.FieldID, selector),
			sqlgraph.To(codeissue.Table, codeissue.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, cwe.IssuesTable, cwe.IssuesPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(cq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CWE entity from the query.
// Returns a *NotFoundError when no CWE was found.
func (cq *CWEQuery) First(ctx context.Context) (*CWE, error) {
	nodes, err := cq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{cwe.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (cq *CWEQuery) FirstX(ctx context.Context) *CWE {
	node, err := cq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first CWE ID from the query.
// Returns a *NotFoundError when no CWE ID was found.
func (cq *CWEQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = cq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{cwe.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (cq *CWEQuery) FirstIDX(ctx context.Context) int {
	id, err := cq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single CWE entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one CWE entity is not found.
// Returns a *NotFoundError when no CWE entities are found.
func (cq *CWEQuery) Only(ctx context.Context) (*CWE, error) {
	nodes, err := cq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{cwe.Label}
	default:
		return nil, &NotSingularError{cwe.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (cq *CWEQuery) OnlyX(ctx context.Context) *CWE {
	node, err := cq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only CWE ID in the query.
// Returns a *NotSingularError when exactly one CWE ID is not found.
// Returns a *NotFoundError when no entities are found.
func (cq *CWEQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = cq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{cwe.Label}
	default:
		err = &NotSingularError{cwe.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (cq *CWEQuery) OnlyIDX(ctx context.Context) int {
	id, err := cq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CWEs.
func (cq *CWEQuery) All(ctx context.Context) ([]*CWE, error) {
	if err := cq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return cq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (cq *CWEQuery) AllX(ctx context.Context) []*CWE {
	nodes, err := cq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of CWE IDs.
func (cq *CWEQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := cq.Select(cwe.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (cq *CWEQuery) IDsX(ctx context.Context) []int {
	ids, err := cq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (cq *CWEQuery) Count(ctx context.Context) (int, error) {
	if err := cq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return cq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (cq *CWEQuery) CountX(ctx context.Context) int {
	count, err := cq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (cq *CWEQuery) Exist(ctx context.Context) (bool, error) {
	if err := cq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return cq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (cq *CWEQuery) ExistX(ctx context.Context) bool {
	exist, err := cq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CWEQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (cq *CWEQuery) Clone() *CWEQuery {
	if cq == nil {
		return nil
	}
	return &CWEQuery{
		config:     cq.config,
		limit:      cq.limit,
		offset:     cq.offset,
		order:      append([]OrderFunc{}, cq.order...),
		predicates: append([]predicate.CWE{}, cq.predicates...),
		withIssues: cq.withIssues.Clone(),
		// clone intermediate query.
		sql:  cq.sql.Clone(),
		path: cq.path,
	}
}

// WithIssues tells the query-builder to eager-load the nodes that are connected to
// the "issues" edge. The optional arguments are used to configure the query builder of the edge.
func (cq *CWEQuery) WithIssues(opts ...func(*CodeIssueQuery)) *CWEQuery {
	query := &CodeIssueQuery{config: cq.config}
	for _, opt := range opts {
		opt(query)
	}
	cq.withIssues = query
	return cq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CweID string `json:"cwe_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.CWE.Query().
//		GroupBy(cwe.FieldCweID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (cq *CWEQuery) GroupBy(field string, fields ...string) *CWEGroupBy {
	group := &CWEGroupBy{config: cq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := cq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return cq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CweID string `json:"cwe_id,omitempty"`
//	}
//
//	client.CWE.Query().
//		Select(cwe.FieldCweID).
//		Scan(ctx, &v)
//
func (cq *CWEQuery) Select(fields ...string) *CWESelect {
	cq.fields = append(cq.fields, fields...)
	return &CWESelect{CWEQuery: cq}
}

func (cq *CWEQuery) prepareQuery(ctx context.Context) error {
	for _, f := range cq.fields {
		if !cwe.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if cq.path != nil {
		prev, err := cq.path(ctx)
		if err != nil {
			return err
		}
		cq.sql = prev
	}
	return nil
}

func (cq *CWEQuery) sqlAll(ctx context.Context) ([]*CWE, error) {
	var (
		nodes       = []*CWE{}
		_spec       = cq.querySpec()
		loadedTypes = [1]bool{
			cq.withIssues != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &CWE{config: cq.config}
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
	if err := sqlgraph.QueryNodes(ctx, cq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := cq.withIssues; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		ids := make(map[int]*CWE, len(nodes))
		for _, node := range nodes {
			ids[node.ID] = node
			fks = append(fks, node.ID)
			node.Edges.Issues = []*CodeIssue{}
		}
		var (
			edgeids []int
			edges   = make(map[int][]*CWE)
		)
		_spec := &sqlgraph.EdgeQuerySpec{
			Edge: &sqlgraph.EdgeSpec{
				Inverse: true,
				Table:   cwe.IssuesTable,
				Columns: cwe.IssuesPrimaryKey,
			},
			Predicate: func(s *sql.Selector) {
				s.Where(sql.InValues(cwe.IssuesPrimaryKey[1], fks...))
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
		if err := sqlgraph.QueryEdges(ctx, cq.driver, _spec); err != nil {
			return nil, fmt.Errorf(`query edges "issues": %w`, err)
		}
		query.Where(codeissue.IDIn(edgeids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := edges[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected "issues" node returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Issues = append(nodes[i].Edges.Issues, n)
			}
		}
	}

	return nodes, nil
}

func (cq *CWEQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := cq.querySpec()
	return sqlgraph.CountNodes(ctx, cq.driver, _spec)
}

func (cq *CWEQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := cq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (cq *CWEQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cwe.Table,
			Columns: cwe.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: cwe.FieldID,
			},
		},
		From:   cq.sql,
		Unique: true,
	}
	if unique := cq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := cq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cwe.FieldID)
		for i := range fields {
			if fields[i] != cwe.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := cq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := cq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := cq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := cq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (cq *CWEQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(cq.driver.Dialect())
	t1 := builder.Table(cwe.Table)
	columns := cq.fields
	if len(columns) == 0 {
		columns = cwe.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if cq.sql != nil {
		selector = cq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	for _, p := range cq.predicates {
		p(selector)
	}
	for _, p := range cq.order {
		p(selector)
	}
	if offset := cq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := cq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CWEGroupBy is the group-by builder for CWE entities.
type CWEGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (cgb *CWEGroupBy) Aggregate(fns ...AggregateFunc) *CWEGroupBy {
	cgb.fns = append(cgb.fns, fns...)
	return cgb
}

// Scan applies the group-by query and scans the result into the given value.
func (cgb *CWEGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := cgb.path(ctx)
	if err != nil {
		return err
	}
	cgb.sql = query
	return cgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (cgb *CWEGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := cgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (cgb *CWEGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(cgb.fields) > 1 {
		return nil, errors.New("ent: CWEGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := cgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (cgb *CWEGroupBy) StringsX(ctx context.Context) []string {
	v, err := cgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (cgb *CWEGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = cgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cwe.Label}
	default:
		err = fmt.Errorf("ent: CWEGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (cgb *CWEGroupBy) StringX(ctx context.Context) string {
	v, err := cgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (cgb *CWEGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(cgb.fields) > 1 {
		return nil, errors.New("ent: CWEGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := cgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (cgb *CWEGroupBy) IntsX(ctx context.Context) []int {
	v, err := cgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (cgb *CWEGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = cgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cwe.Label}
	default:
		err = fmt.Errorf("ent: CWEGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (cgb *CWEGroupBy) IntX(ctx context.Context) int {
	v, err := cgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (cgb *CWEGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(cgb.fields) > 1 {
		return nil, errors.New("ent: CWEGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := cgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (cgb *CWEGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := cgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (cgb *CWEGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = cgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cwe.Label}
	default:
		err = fmt.Errorf("ent: CWEGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (cgb *CWEGroupBy) Float64X(ctx context.Context) float64 {
	v, err := cgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (cgb *CWEGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(cgb.fields) > 1 {
		return nil, errors.New("ent: CWEGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := cgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (cgb *CWEGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := cgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (cgb *CWEGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = cgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cwe.Label}
	default:
		err = fmt.Errorf("ent: CWEGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (cgb *CWEGroupBy) BoolX(ctx context.Context) bool {
	v, err := cgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (cgb *CWEGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range cgb.fields {
		if !cwe.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := cgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (cgb *CWEGroupBy) sqlQuery() *sql.Selector {
	selector := cgb.sql.Select()
	aggregation := make([]string, 0, len(cgb.fns))
	for _, fn := range cgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(cgb.fields)+len(cgb.fns))
		for _, f := range cgb.fields {
			columns = append(columns, selector.C(f))
		}
		for _, c := range aggregation {
			columns = append(columns, c)
		}
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(cgb.fields...)...)
}

// CWESelect is the builder for selecting fields of CWE entities.
type CWESelect struct {
	*CWEQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (cs *CWESelect) Scan(ctx context.Context, v interface{}) error {
	if err := cs.prepareQuery(ctx); err != nil {
		return err
	}
	cs.sql = cs.CWEQuery.sqlQuery(ctx)
	return cs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (cs *CWESelect) ScanX(ctx context.Context, v interface{}) {
	if err := cs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (cs *CWESelect) Strings(ctx context.Context) ([]string, error) {
	if len(cs.fields) > 1 {
		return nil, errors.New("ent: CWESelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := cs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (cs *CWESelect) StringsX(ctx context.Context) []string {
	v, err := cs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (cs *CWESelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = cs.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cwe.Label}
	default:
		err = fmt.Errorf("ent: CWESelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (cs *CWESelect) StringX(ctx context.Context) string {
	v, err := cs.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (cs *CWESelect) Ints(ctx context.Context) ([]int, error) {
	if len(cs.fields) > 1 {
		return nil, errors.New("ent: CWESelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := cs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (cs *CWESelect) IntsX(ctx context.Context) []int {
	v, err := cs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (cs *CWESelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = cs.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cwe.Label}
	default:
		err = fmt.Errorf("ent: CWESelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (cs *CWESelect) IntX(ctx context.Context) int {
	v, err := cs.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (cs *CWESelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(cs.fields) > 1 {
		return nil, errors.New("ent: CWESelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := cs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (cs *CWESelect) Float64sX(ctx context.Context) []float64 {
	v, err := cs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (cs *CWESelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = cs.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cwe.Label}
	default:
		err = fmt.Errorf("ent: CWESelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (cs *CWESelect) Float64X(ctx context.Context) float64 {
	v, err := cs.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (cs *CWESelect) Bools(ctx context.Context) ([]bool, error) {
	if len(cs.fields) > 1 {
		return nil, errors.New("ent: CWESelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := cs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (cs *CWESelect) BoolsX(ctx context.Context) []bool {
	v, err := cs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (cs *CWESelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = cs.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{cwe.Label}
	default:
		err = fmt.Errorf("ent: CWESelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (cs *CWESelect) BoolX(ctx context.Context) bool {
	v, err := cs.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (cs *CWESelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := cs.sql.Query()
	if err := cs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}