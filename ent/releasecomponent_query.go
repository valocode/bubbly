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
	"github.com/valocode/bubbly/ent/codescan"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasecomponent"
	"github.com/valocode/bubbly/ent/releasevulnerability"
)

// ReleaseComponentQuery is the builder for querying ReleaseComponent entities.
type ReleaseComponentQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.ReleaseComponent
	// eager-loading edges.
	withRelease         *ReleaseQuery
	withScans           *CodeScanQuery
	withComponent       *ComponentQuery
	withVulnerabilities *ReleaseVulnerabilityQuery
	withFKs             bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ReleaseComponentQuery builder.
func (rcq *ReleaseComponentQuery) Where(ps ...predicate.ReleaseComponent) *ReleaseComponentQuery {
	rcq.predicates = append(rcq.predicates, ps...)
	return rcq
}

// Limit adds a limit step to the query.
func (rcq *ReleaseComponentQuery) Limit(limit int) *ReleaseComponentQuery {
	rcq.limit = &limit
	return rcq
}

// Offset adds an offset step to the query.
func (rcq *ReleaseComponentQuery) Offset(offset int) *ReleaseComponentQuery {
	rcq.offset = &offset
	return rcq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rcq *ReleaseComponentQuery) Unique(unique bool) *ReleaseComponentQuery {
	rcq.unique = &unique
	return rcq
}

// Order adds an order step to the query.
func (rcq *ReleaseComponentQuery) Order(o ...OrderFunc) *ReleaseComponentQuery {
	rcq.order = append(rcq.order, o...)
	return rcq
}

// QueryRelease chains the current query on the "release" edge.
func (rcq *ReleaseComponentQuery) QueryRelease() *ReleaseQuery {
	query := &ReleaseQuery{config: rcq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(releasecomponent.Table, releasecomponent.FieldID, selector),
			sqlgraph.To(release.Table, release.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, releasecomponent.ReleaseTable, releasecomponent.ReleaseColumn),
		)
		fromU = sqlgraph.SetNeighbors(rcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryScans chains the current query on the "scans" edge.
func (rcq *ReleaseComponentQuery) QueryScans() *CodeScanQuery {
	query := &CodeScanQuery{config: rcq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(releasecomponent.Table, releasecomponent.FieldID, selector),
			sqlgraph.To(codescan.Table, codescan.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, releasecomponent.ScansTable, releasecomponent.ScansPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(rcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryComponent chains the current query on the "component" edge.
func (rcq *ReleaseComponentQuery) QueryComponent() *ComponentQuery {
	query := &ComponentQuery{config: rcq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(releasecomponent.Table, releasecomponent.FieldID, selector),
			sqlgraph.To(component.Table, component.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, releasecomponent.ComponentTable, releasecomponent.ComponentColumn),
		)
		fromU = sqlgraph.SetNeighbors(rcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryVulnerabilities chains the current query on the "vulnerabilities" edge.
func (rcq *ReleaseComponentQuery) QueryVulnerabilities() *ReleaseVulnerabilityQuery {
	query := &ReleaseVulnerabilityQuery{config: rcq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(releasecomponent.Table, releasecomponent.FieldID, selector),
			sqlgraph.To(releasevulnerability.Table, releasevulnerability.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, releasecomponent.VulnerabilitiesTable, releasecomponent.VulnerabilitiesPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(rcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ReleaseComponent entity from the query.
// Returns a *NotFoundError when no ReleaseComponent was found.
func (rcq *ReleaseComponentQuery) First(ctx context.Context) (*ReleaseComponent, error) {
	nodes, err := rcq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{releasecomponent.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rcq *ReleaseComponentQuery) FirstX(ctx context.Context) *ReleaseComponent {
	node, err := rcq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ReleaseComponent ID from the query.
// Returns a *NotFoundError when no ReleaseComponent ID was found.
func (rcq *ReleaseComponentQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rcq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{releasecomponent.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rcq *ReleaseComponentQuery) FirstIDX(ctx context.Context) int {
	id, err := rcq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ReleaseComponent entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one ReleaseComponent entity is not found.
// Returns a *NotFoundError when no ReleaseComponent entities are found.
func (rcq *ReleaseComponentQuery) Only(ctx context.Context) (*ReleaseComponent, error) {
	nodes, err := rcq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{releasecomponent.Label}
	default:
		return nil, &NotSingularError{releasecomponent.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rcq *ReleaseComponentQuery) OnlyX(ctx context.Context) *ReleaseComponent {
	node, err := rcq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ReleaseComponent ID in the query.
// Returns a *NotSingularError when exactly one ReleaseComponent ID is not found.
// Returns a *NotFoundError when no entities are found.
func (rcq *ReleaseComponentQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rcq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{releasecomponent.Label}
	default:
		err = &NotSingularError{releasecomponent.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rcq *ReleaseComponentQuery) OnlyIDX(ctx context.Context) int {
	id, err := rcq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ReleaseComponents.
func (rcq *ReleaseComponentQuery) All(ctx context.Context) ([]*ReleaseComponent, error) {
	if err := rcq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return rcq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (rcq *ReleaseComponentQuery) AllX(ctx context.Context) []*ReleaseComponent {
	nodes, err := rcq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ReleaseComponent IDs.
func (rcq *ReleaseComponentQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := rcq.Select(releasecomponent.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rcq *ReleaseComponentQuery) IDsX(ctx context.Context) []int {
	ids, err := rcq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rcq *ReleaseComponentQuery) Count(ctx context.Context) (int, error) {
	if err := rcq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return rcq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (rcq *ReleaseComponentQuery) CountX(ctx context.Context) int {
	count, err := rcq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rcq *ReleaseComponentQuery) Exist(ctx context.Context) (bool, error) {
	if err := rcq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return rcq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (rcq *ReleaseComponentQuery) ExistX(ctx context.Context) bool {
	exist, err := rcq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ReleaseComponentQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rcq *ReleaseComponentQuery) Clone() *ReleaseComponentQuery {
	if rcq == nil {
		return nil
	}
	return &ReleaseComponentQuery{
		config:              rcq.config,
		limit:               rcq.limit,
		offset:              rcq.offset,
		order:               append([]OrderFunc{}, rcq.order...),
		predicates:          append([]predicate.ReleaseComponent{}, rcq.predicates...),
		withRelease:         rcq.withRelease.Clone(),
		withScans:           rcq.withScans.Clone(),
		withComponent:       rcq.withComponent.Clone(),
		withVulnerabilities: rcq.withVulnerabilities.Clone(),
		// clone intermediate query.
		sql:  rcq.sql.Clone(),
		path: rcq.path,
	}
}

// WithRelease tells the query-builder to eager-load the nodes that are connected to
// the "release" edge. The optional arguments are used to configure the query builder of the edge.
func (rcq *ReleaseComponentQuery) WithRelease(opts ...func(*ReleaseQuery)) *ReleaseComponentQuery {
	query := &ReleaseQuery{config: rcq.config}
	for _, opt := range opts {
		opt(query)
	}
	rcq.withRelease = query
	return rcq
}

// WithScans tells the query-builder to eager-load the nodes that are connected to
// the "scans" edge. The optional arguments are used to configure the query builder of the edge.
func (rcq *ReleaseComponentQuery) WithScans(opts ...func(*CodeScanQuery)) *ReleaseComponentQuery {
	query := &CodeScanQuery{config: rcq.config}
	for _, opt := range opts {
		opt(query)
	}
	rcq.withScans = query
	return rcq
}

// WithComponent tells the query-builder to eager-load the nodes that are connected to
// the "component" edge. The optional arguments are used to configure the query builder of the edge.
func (rcq *ReleaseComponentQuery) WithComponent(opts ...func(*ComponentQuery)) *ReleaseComponentQuery {
	query := &ComponentQuery{config: rcq.config}
	for _, opt := range opts {
		opt(query)
	}
	rcq.withComponent = query
	return rcq
}

// WithVulnerabilities tells the query-builder to eager-load the nodes that are connected to
// the "vulnerabilities" edge. The optional arguments are used to configure the query builder of the edge.
func (rcq *ReleaseComponentQuery) WithVulnerabilities(opts ...func(*ReleaseVulnerabilityQuery)) *ReleaseComponentQuery {
	query := &ReleaseVulnerabilityQuery{config: rcq.config}
	for _, opt := range opts {
		opt(query)
	}
	rcq.withVulnerabilities = query
	return rcq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (rcq *ReleaseComponentQuery) GroupBy(field string, fields ...string) *ReleaseComponentGroupBy {
	group := &ReleaseComponentGroupBy{config: rcq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := rcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return rcq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (rcq *ReleaseComponentQuery) Select(fields ...string) *ReleaseComponentSelect {
	rcq.fields = append(rcq.fields, fields...)
	return &ReleaseComponentSelect{ReleaseComponentQuery: rcq}
}

func (rcq *ReleaseComponentQuery) prepareQuery(ctx context.Context) error {
	for _, f := range rcq.fields {
		if !releasecomponent.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if rcq.path != nil {
		prev, err := rcq.path(ctx)
		if err != nil {
			return err
		}
		rcq.sql = prev
	}
	return nil
}

func (rcq *ReleaseComponentQuery) sqlAll(ctx context.Context) ([]*ReleaseComponent, error) {
	var (
		nodes       = []*ReleaseComponent{}
		withFKs     = rcq.withFKs
		_spec       = rcq.querySpec()
		loadedTypes = [4]bool{
			rcq.withRelease != nil,
			rcq.withScans != nil,
			rcq.withComponent != nil,
			rcq.withVulnerabilities != nil,
		}
	)
	if rcq.withRelease != nil || rcq.withComponent != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, releasecomponent.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &ReleaseComponent{config: rcq.config}
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
	if err := sqlgraph.QueryNodes(ctx, rcq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := rcq.withRelease; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*ReleaseComponent)
		for i := range nodes {
			if nodes[i].release_component_release == nil {
				continue
			}
			fk := *nodes[i].release_component_release
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(release.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "release_component_release" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Release = n
			}
		}
	}

	if query := rcq.withScans; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		ids := make(map[int]*ReleaseComponent, len(nodes))
		for _, node := range nodes {
			ids[node.ID] = node
			fks = append(fks, node.ID)
			node.Edges.Scans = []*CodeScan{}
		}
		var (
			edgeids []int
			edges   = make(map[int][]*ReleaseComponent)
		)
		_spec := &sqlgraph.EdgeQuerySpec{
			Edge: &sqlgraph.EdgeSpec{
				Inverse: false,
				Table:   releasecomponent.ScansTable,
				Columns: releasecomponent.ScansPrimaryKey,
			},
			Predicate: func(s *sql.Selector) {
				s.Where(sql.InValues(releasecomponent.ScansPrimaryKey[0], fks...))
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
		if err := sqlgraph.QueryEdges(ctx, rcq.driver, _spec); err != nil {
			return nil, fmt.Errorf(`query edges "scans": %w`, err)
		}
		query.Where(codescan.IDIn(edgeids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := edges[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected "scans" node returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Scans = append(nodes[i].Edges.Scans, n)
			}
		}
	}

	if query := rcq.withComponent; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*ReleaseComponent)
		for i := range nodes {
			if nodes[i].release_component_component == nil {
				continue
			}
			fk := *nodes[i].release_component_component
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(component.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "release_component_component" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Component = n
			}
		}
	}

	if query := rcq.withVulnerabilities; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		ids := make(map[int]*ReleaseComponent, len(nodes))
		for _, node := range nodes {
			ids[node.ID] = node
			fks = append(fks, node.ID)
			node.Edges.Vulnerabilities = []*ReleaseVulnerability{}
		}
		var (
			edgeids []int
			edges   = make(map[int][]*ReleaseComponent)
		)
		_spec := &sqlgraph.EdgeQuerySpec{
			Edge: &sqlgraph.EdgeSpec{
				Inverse: true,
				Table:   releasecomponent.VulnerabilitiesTable,
				Columns: releasecomponent.VulnerabilitiesPrimaryKey,
			},
			Predicate: func(s *sql.Selector) {
				s.Where(sql.InValues(releasecomponent.VulnerabilitiesPrimaryKey[1], fks...))
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
		if err := sqlgraph.QueryEdges(ctx, rcq.driver, _spec); err != nil {
			return nil, fmt.Errorf(`query edges "vulnerabilities": %w`, err)
		}
		query.Where(releasevulnerability.IDIn(edgeids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := edges[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected "vulnerabilities" node returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Vulnerabilities = append(nodes[i].Edges.Vulnerabilities, n)
			}
		}
	}

	return nodes, nil
}

func (rcq *ReleaseComponentQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rcq.querySpec()
	return sqlgraph.CountNodes(ctx, rcq.driver, _spec)
}

func (rcq *ReleaseComponentQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := rcq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (rcq *ReleaseComponentQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   releasecomponent.Table,
			Columns: releasecomponent.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: releasecomponent.FieldID,
			},
		},
		From:   rcq.sql,
		Unique: true,
	}
	if unique := rcq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := rcq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, releasecomponent.FieldID)
		for i := range fields {
			if fields[i] != releasecomponent.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := rcq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rcq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rcq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rcq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rcq *ReleaseComponentQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rcq.driver.Dialect())
	t1 := builder.Table(releasecomponent.Table)
	columns := rcq.fields
	if len(columns) == 0 {
		columns = releasecomponent.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if rcq.sql != nil {
		selector = rcq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	for _, p := range rcq.predicates {
		p(selector)
	}
	for _, p := range rcq.order {
		p(selector)
	}
	if offset := rcq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rcq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ReleaseComponentGroupBy is the group-by builder for ReleaseComponent entities.
type ReleaseComponentGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rcgb *ReleaseComponentGroupBy) Aggregate(fns ...AggregateFunc) *ReleaseComponentGroupBy {
	rcgb.fns = append(rcgb.fns, fns...)
	return rcgb
}

// Scan applies the group-by query and scans the result into the given value.
func (rcgb *ReleaseComponentGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := rcgb.path(ctx)
	if err != nil {
		return err
	}
	rcgb.sql = query
	return rcgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (rcgb *ReleaseComponentGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := rcgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (rcgb *ReleaseComponentGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(rcgb.fields) > 1 {
		return nil, errors.New("ent: ReleaseComponentGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := rcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (rcgb *ReleaseComponentGroupBy) StringsX(ctx context.Context) []string {
	v, err := rcgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rcgb *ReleaseComponentGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = rcgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasecomponent.Label}
	default:
		err = fmt.Errorf("ent: ReleaseComponentGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (rcgb *ReleaseComponentGroupBy) StringX(ctx context.Context) string {
	v, err := rcgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (rcgb *ReleaseComponentGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(rcgb.fields) > 1 {
		return nil, errors.New("ent: ReleaseComponentGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := rcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (rcgb *ReleaseComponentGroupBy) IntsX(ctx context.Context) []int {
	v, err := rcgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rcgb *ReleaseComponentGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = rcgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasecomponent.Label}
	default:
		err = fmt.Errorf("ent: ReleaseComponentGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (rcgb *ReleaseComponentGroupBy) IntX(ctx context.Context) int {
	v, err := rcgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (rcgb *ReleaseComponentGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(rcgb.fields) > 1 {
		return nil, errors.New("ent: ReleaseComponentGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := rcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (rcgb *ReleaseComponentGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := rcgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rcgb *ReleaseComponentGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = rcgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasecomponent.Label}
	default:
		err = fmt.Errorf("ent: ReleaseComponentGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (rcgb *ReleaseComponentGroupBy) Float64X(ctx context.Context) float64 {
	v, err := rcgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (rcgb *ReleaseComponentGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(rcgb.fields) > 1 {
		return nil, errors.New("ent: ReleaseComponentGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := rcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (rcgb *ReleaseComponentGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := rcgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rcgb *ReleaseComponentGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = rcgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasecomponent.Label}
	default:
		err = fmt.Errorf("ent: ReleaseComponentGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (rcgb *ReleaseComponentGroupBy) BoolX(ctx context.Context) bool {
	v, err := rcgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (rcgb *ReleaseComponentGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range rcgb.fields {
		if !releasecomponent.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := rcgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rcgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (rcgb *ReleaseComponentGroupBy) sqlQuery() *sql.Selector {
	selector := rcgb.sql.Select()
	aggregation := make([]string, 0, len(rcgb.fns))
	for _, fn := range rcgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(rcgb.fields)+len(rcgb.fns))
		for _, f := range rcgb.fields {
			columns = append(columns, selector.C(f))
		}
		for _, c := range aggregation {
			columns = append(columns, c)
		}
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(rcgb.fields...)...)
}

// ReleaseComponentSelect is the builder for selecting fields of ReleaseComponent entities.
type ReleaseComponentSelect struct {
	*ReleaseComponentQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (rcs *ReleaseComponentSelect) Scan(ctx context.Context, v interface{}) error {
	if err := rcs.prepareQuery(ctx); err != nil {
		return err
	}
	rcs.sql = rcs.ReleaseComponentQuery.sqlQuery(ctx)
	return rcs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (rcs *ReleaseComponentSelect) ScanX(ctx context.Context, v interface{}) {
	if err := rcs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (rcs *ReleaseComponentSelect) Strings(ctx context.Context) ([]string, error) {
	if len(rcs.fields) > 1 {
		return nil, errors.New("ent: ReleaseComponentSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := rcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (rcs *ReleaseComponentSelect) StringsX(ctx context.Context) []string {
	v, err := rcs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (rcs *ReleaseComponentSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = rcs.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasecomponent.Label}
	default:
		err = fmt.Errorf("ent: ReleaseComponentSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (rcs *ReleaseComponentSelect) StringX(ctx context.Context) string {
	v, err := rcs.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (rcs *ReleaseComponentSelect) Ints(ctx context.Context) ([]int, error) {
	if len(rcs.fields) > 1 {
		return nil, errors.New("ent: ReleaseComponentSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := rcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (rcs *ReleaseComponentSelect) IntsX(ctx context.Context) []int {
	v, err := rcs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (rcs *ReleaseComponentSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = rcs.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasecomponent.Label}
	default:
		err = fmt.Errorf("ent: ReleaseComponentSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (rcs *ReleaseComponentSelect) IntX(ctx context.Context) int {
	v, err := rcs.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (rcs *ReleaseComponentSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(rcs.fields) > 1 {
		return nil, errors.New("ent: ReleaseComponentSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := rcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (rcs *ReleaseComponentSelect) Float64sX(ctx context.Context) []float64 {
	v, err := rcs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (rcs *ReleaseComponentSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = rcs.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasecomponent.Label}
	default:
		err = fmt.Errorf("ent: ReleaseComponentSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (rcs *ReleaseComponentSelect) Float64X(ctx context.Context) float64 {
	v, err := rcs.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (rcs *ReleaseComponentSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(rcs.fields) > 1 {
		return nil, errors.New("ent: ReleaseComponentSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := rcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (rcs *ReleaseComponentSelect) BoolsX(ctx context.Context) []bool {
	v, err := rcs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (rcs *ReleaseComponentSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = rcs.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{releasecomponent.Label}
	default:
		err = fmt.Errorf("ent: ReleaseComponentSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (rcs *ReleaseComponentSelect) BoolX(ctx context.Context) bool {
	v, err := rcs.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (rcs *ReleaseComponentSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := rcs.sql.Query()
	if err := rcs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
