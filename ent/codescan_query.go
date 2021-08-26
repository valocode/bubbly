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
	"github.com/valocode/bubbly/ent/codescan"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasecomponent"
	"github.com/valocode/bubbly/ent/releaseentry"
	"github.com/valocode/bubbly/ent/releasevulnerability"
)

// CodeScanQuery is the builder for querying CodeScan entities.
type CodeScanQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.CodeScan
	// eager-loading edges.
	withRelease         *ReleaseQuery
	withEntry           *ReleaseEntryQuery
	withIssues          *CodeIssueQuery
	withVulnerabilities *ReleaseVulnerabilityQuery
	withComponents      *ReleaseComponentQuery
	withFKs             bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CodeScanQuery builder.
func (csq *CodeScanQuery) Where(ps ...predicate.CodeScan) *CodeScanQuery {
	csq.predicates = append(csq.predicates, ps...)
	return csq
}

// Limit adds a limit step to the query.
func (csq *CodeScanQuery) Limit(limit int) *CodeScanQuery {
	csq.limit = &limit
	return csq
}

// Offset adds an offset step to the query.
func (csq *CodeScanQuery) Offset(offset int) *CodeScanQuery {
	csq.offset = &offset
	return csq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (csq *CodeScanQuery) Unique(unique bool) *CodeScanQuery {
	csq.unique = &unique
	return csq
}

// Order adds an order step to the query.
func (csq *CodeScanQuery) Order(o ...OrderFunc) *CodeScanQuery {
	csq.order = append(csq.order, o...)
	return csq
}

// QueryRelease chains the current query on the "release" edge.
func (csq *CodeScanQuery) QueryRelease() *ReleaseQuery {
	query := &ReleaseQuery{config: csq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := csq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := csq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(codescan.Table, codescan.FieldID, selector),
			sqlgraph.To(release.Table, release.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, codescan.ReleaseTable, codescan.ReleaseColumn),
		)
		fromU = sqlgraph.SetNeighbors(csq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEntry chains the current query on the "entry" edge.
func (csq *CodeScanQuery) QueryEntry() *ReleaseEntryQuery {
	query := &ReleaseEntryQuery{config: csq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := csq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := csq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(codescan.Table, codescan.FieldID, selector),
			sqlgraph.To(releaseentry.Table, releaseentry.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, codescan.EntryTable, codescan.EntryColumn),
		)
		fromU = sqlgraph.SetNeighbors(csq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryIssues chains the current query on the "issues" edge.
func (csq *CodeScanQuery) QueryIssues() *CodeIssueQuery {
	query := &CodeIssueQuery{config: csq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := csq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := csq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(codescan.Table, codescan.FieldID, selector),
			sqlgraph.To(codeissue.Table, codeissue.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, codescan.IssuesTable, codescan.IssuesColumn),
		)
		fromU = sqlgraph.SetNeighbors(csq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryVulnerabilities chains the current query on the "vulnerabilities" edge.
func (csq *CodeScanQuery) QueryVulnerabilities() *ReleaseVulnerabilityQuery {
	query := &ReleaseVulnerabilityQuery{config: csq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := csq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := csq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(codescan.Table, codescan.FieldID, selector),
			sqlgraph.To(releasevulnerability.Table, releasevulnerability.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, codescan.VulnerabilitiesTable, codescan.VulnerabilitiesColumn),
		)
		fromU = sqlgraph.SetNeighbors(csq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryComponents chains the current query on the "components" edge.
func (csq *CodeScanQuery) QueryComponents() *ReleaseComponentQuery {
	query := &ReleaseComponentQuery{config: csq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := csq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := csq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(codescan.Table, codescan.FieldID, selector),
			sqlgraph.To(releasecomponent.Table, releasecomponent.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, codescan.ComponentsTable, codescan.ComponentsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(csq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CodeScan entity from the query.
// Returns a *NotFoundError when no CodeScan was found.
func (csq *CodeScanQuery) First(ctx context.Context) (*CodeScan, error) {
	nodes, err := csq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{codescan.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (csq *CodeScanQuery) FirstX(ctx context.Context) *CodeScan {
	node, err := csq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first CodeScan ID from the query.
// Returns a *NotFoundError when no CodeScan ID was found.
func (csq *CodeScanQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = csq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{codescan.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (csq *CodeScanQuery) FirstIDX(ctx context.Context) int {
	id, err := csq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single CodeScan entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one CodeScan entity is not found.
// Returns a *NotFoundError when no CodeScan entities are found.
func (csq *CodeScanQuery) Only(ctx context.Context) (*CodeScan, error) {
	nodes, err := csq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{codescan.Label}
	default:
		return nil, &NotSingularError{codescan.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (csq *CodeScanQuery) OnlyX(ctx context.Context) *CodeScan {
	node, err := csq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only CodeScan ID in the query.
// Returns a *NotSingularError when exactly one CodeScan ID is not found.
// Returns a *NotFoundError when no entities are found.
func (csq *CodeScanQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = csq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{codescan.Label}
	default:
		err = &NotSingularError{codescan.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (csq *CodeScanQuery) OnlyIDX(ctx context.Context) int {
	id, err := csq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CodeScans.
func (csq *CodeScanQuery) All(ctx context.Context) ([]*CodeScan, error) {
	if err := csq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return csq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (csq *CodeScanQuery) AllX(ctx context.Context) []*CodeScan {
	nodes, err := csq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of CodeScan IDs.
func (csq *CodeScanQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := csq.Select(codescan.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (csq *CodeScanQuery) IDsX(ctx context.Context) []int {
	ids, err := csq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (csq *CodeScanQuery) Count(ctx context.Context) (int, error) {
	if err := csq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return csq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (csq *CodeScanQuery) CountX(ctx context.Context) int {
	count, err := csq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (csq *CodeScanQuery) Exist(ctx context.Context) (bool, error) {
	if err := csq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return csq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (csq *CodeScanQuery) ExistX(ctx context.Context) bool {
	exist, err := csq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CodeScanQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (csq *CodeScanQuery) Clone() *CodeScanQuery {
	if csq == nil {
		return nil
	}
	return &CodeScanQuery{
		config:              csq.config,
		limit:               csq.limit,
		offset:              csq.offset,
		order:               append([]OrderFunc{}, csq.order...),
		predicates:          append([]predicate.CodeScan{}, csq.predicates...),
		withRelease:         csq.withRelease.Clone(),
		withEntry:           csq.withEntry.Clone(),
		withIssues:          csq.withIssues.Clone(),
		withVulnerabilities: csq.withVulnerabilities.Clone(),
		withComponents:      csq.withComponents.Clone(),
		// clone intermediate query.
		sql:  csq.sql.Clone(),
		path: csq.path,
	}
}

// WithRelease tells the query-builder to eager-load the nodes that are connected to
// the "release" edge. The optional arguments are used to configure the query builder of the edge.
func (csq *CodeScanQuery) WithRelease(opts ...func(*ReleaseQuery)) *CodeScanQuery {
	query := &ReleaseQuery{config: csq.config}
	for _, opt := range opts {
		opt(query)
	}
	csq.withRelease = query
	return csq
}

// WithEntry tells the query-builder to eager-load the nodes that are connected to
// the "entry" edge. The optional arguments are used to configure the query builder of the edge.
func (csq *CodeScanQuery) WithEntry(opts ...func(*ReleaseEntryQuery)) *CodeScanQuery {
	query := &ReleaseEntryQuery{config: csq.config}
	for _, opt := range opts {
		opt(query)
	}
	csq.withEntry = query
	return csq
}

// WithIssues tells the query-builder to eager-load the nodes that are connected to
// the "issues" edge. The optional arguments are used to configure the query builder of the edge.
func (csq *CodeScanQuery) WithIssues(opts ...func(*CodeIssueQuery)) *CodeScanQuery {
	query := &CodeIssueQuery{config: csq.config}
	for _, opt := range opts {
		opt(query)
	}
	csq.withIssues = query
	return csq
}

// WithVulnerabilities tells the query-builder to eager-load the nodes that are connected to
// the "vulnerabilities" edge. The optional arguments are used to configure the query builder of the edge.
func (csq *CodeScanQuery) WithVulnerabilities(opts ...func(*ReleaseVulnerabilityQuery)) *CodeScanQuery {
	query := &ReleaseVulnerabilityQuery{config: csq.config}
	for _, opt := range opts {
		opt(query)
	}
	csq.withVulnerabilities = query
	return csq
}

// WithComponents tells the query-builder to eager-load the nodes that are connected to
// the "components" edge. The optional arguments are used to configure the query builder of the edge.
func (csq *CodeScanQuery) WithComponents(opts ...func(*ReleaseComponentQuery)) *CodeScanQuery {
	query := &ReleaseComponentQuery{config: csq.config}
	for _, opt := range opts {
		opt(query)
	}
	csq.withComponents = query
	return csq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Tool string `json:"tool,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.CodeScan.Query().
//		GroupBy(codescan.FieldTool).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (csq *CodeScanQuery) GroupBy(field string, fields ...string) *CodeScanGroupBy {
	group := &CodeScanGroupBy{config: csq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := csq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return csq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Tool string `json:"tool,omitempty"`
//	}
//
//	client.CodeScan.Query().
//		Select(codescan.FieldTool).
//		Scan(ctx, &v)
//
func (csq *CodeScanQuery) Select(fields ...string) *CodeScanSelect {
	csq.fields = append(csq.fields, fields...)
	return &CodeScanSelect{CodeScanQuery: csq}
}

func (csq *CodeScanQuery) prepareQuery(ctx context.Context) error {
	for _, f := range csq.fields {
		if !codescan.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if csq.path != nil {
		prev, err := csq.path(ctx)
		if err != nil {
			return err
		}
		csq.sql = prev
	}
	return nil
}

func (csq *CodeScanQuery) sqlAll(ctx context.Context) ([]*CodeScan, error) {
	var (
		nodes       = []*CodeScan{}
		withFKs     = csq.withFKs
		_spec       = csq.querySpec()
		loadedTypes = [5]bool{
			csq.withRelease != nil,
			csq.withEntry != nil,
			csq.withIssues != nil,
			csq.withVulnerabilities != nil,
			csq.withComponents != nil,
		}
	)
	if csq.withRelease != nil || csq.withEntry != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, codescan.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &CodeScan{config: csq.config}
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
	if err := sqlgraph.QueryNodes(ctx, csq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := csq.withRelease; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*CodeScan)
		for i := range nodes {
			if nodes[i].code_scan_release == nil {
				continue
			}
			fk := *nodes[i].code_scan_release
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
				return nil, fmt.Errorf(`unexpected foreign-key "code_scan_release" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Release = n
			}
		}
	}

	if query := csq.withEntry; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*CodeScan)
		for i := range nodes {
			if nodes[i].release_entry_code_scan == nil {
				continue
			}
			fk := *nodes[i].release_entry_code_scan
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(releaseentry.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "release_entry_code_scan" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Entry = n
			}
		}
	}

	if query := csq.withIssues; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*CodeScan)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
			nodes[i].Edges.Issues = []*CodeIssue{}
		}
		query.withFKs = true
		query.Where(predicate.CodeIssue(func(s *sql.Selector) {
			s.Where(sql.InValues(codescan.IssuesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.code_issue_scan
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "code_issue_scan" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "code_issue_scan" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Issues = append(node.Edges.Issues, n)
		}
	}

	if query := csq.withVulnerabilities; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*CodeScan)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
			nodes[i].Edges.Vulnerabilities = []*ReleaseVulnerability{}
		}
		query.withFKs = true
		query.Where(predicate.ReleaseVulnerability(func(s *sql.Selector) {
			s.Where(sql.InValues(codescan.VulnerabilitiesColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.release_vulnerability_scan
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "release_vulnerability_scan" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "release_vulnerability_scan" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Vulnerabilities = append(node.Edges.Vulnerabilities, n)
		}
	}

	if query := csq.withComponents; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		ids := make(map[int]*CodeScan, len(nodes))
		for _, node := range nodes {
			ids[node.ID] = node
			fks = append(fks, node.ID)
			node.Edges.Components = []*ReleaseComponent{}
		}
		var (
			edgeids []int
			edges   = make(map[int][]*CodeScan)
		)
		_spec := &sqlgraph.EdgeQuerySpec{
			Edge: &sqlgraph.EdgeSpec{
				Inverse: true,
				Table:   codescan.ComponentsTable,
				Columns: codescan.ComponentsPrimaryKey,
			},
			Predicate: func(s *sql.Selector) {
				s.Where(sql.InValues(codescan.ComponentsPrimaryKey[1], fks...))
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
		if err := sqlgraph.QueryEdges(ctx, csq.driver, _spec); err != nil {
			return nil, fmt.Errorf(`query edges "components": %w`, err)
		}
		query.Where(releasecomponent.IDIn(edgeids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := edges[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected "components" node returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Components = append(nodes[i].Edges.Components, n)
			}
		}
	}

	return nodes, nil
}

func (csq *CodeScanQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := csq.querySpec()
	return sqlgraph.CountNodes(ctx, csq.driver, _spec)
}

func (csq *CodeScanQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := csq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (csq *CodeScanQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   codescan.Table,
			Columns: codescan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: codescan.FieldID,
			},
		},
		From:   csq.sql,
		Unique: true,
	}
	if unique := csq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := csq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, codescan.FieldID)
		for i := range fields {
			if fields[i] != codescan.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := csq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := csq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := csq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := csq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (csq *CodeScanQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(csq.driver.Dialect())
	t1 := builder.Table(codescan.Table)
	columns := csq.fields
	if len(columns) == 0 {
		columns = codescan.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if csq.sql != nil {
		selector = csq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	for _, p := range csq.predicates {
		p(selector)
	}
	for _, p := range csq.order {
		p(selector)
	}
	if offset := csq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := csq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CodeScanGroupBy is the group-by builder for CodeScan entities.
type CodeScanGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (csgb *CodeScanGroupBy) Aggregate(fns ...AggregateFunc) *CodeScanGroupBy {
	csgb.fns = append(csgb.fns, fns...)
	return csgb
}

// Scan applies the group-by query and scans the result into the given value.
func (csgb *CodeScanGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := csgb.path(ctx)
	if err != nil {
		return err
	}
	csgb.sql = query
	return csgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (csgb *CodeScanGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := csgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (csgb *CodeScanGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(csgb.fields) > 1 {
		return nil, errors.New("ent: CodeScanGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := csgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (csgb *CodeScanGroupBy) StringsX(ctx context.Context) []string {
	v, err := csgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (csgb *CodeScanGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = csgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{codescan.Label}
	default:
		err = fmt.Errorf("ent: CodeScanGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (csgb *CodeScanGroupBy) StringX(ctx context.Context) string {
	v, err := csgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (csgb *CodeScanGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(csgb.fields) > 1 {
		return nil, errors.New("ent: CodeScanGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := csgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (csgb *CodeScanGroupBy) IntsX(ctx context.Context) []int {
	v, err := csgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (csgb *CodeScanGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = csgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{codescan.Label}
	default:
		err = fmt.Errorf("ent: CodeScanGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (csgb *CodeScanGroupBy) IntX(ctx context.Context) int {
	v, err := csgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (csgb *CodeScanGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(csgb.fields) > 1 {
		return nil, errors.New("ent: CodeScanGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := csgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (csgb *CodeScanGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := csgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (csgb *CodeScanGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = csgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{codescan.Label}
	default:
		err = fmt.Errorf("ent: CodeScanGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (csgb *CodeScanGroupBy) Float64X(ctx context.Context) float64 {
	v, err := csgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (csgb *CodeScanGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(csgb.fields) > 1 {
		return nil, errors.New("ent: CodeScanGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := csgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (csgb *CodeScanGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := csgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (csgb *CodeScanGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = csgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{codescan.Label}
	default:
		err = fmt.Errorf("ent: CodeScanGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (csgb *CodeScanGroupBy) BoolX(ctx context.Context) bool {
	v, err := csgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (csgb *CodeScanGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range csgb.fields {
		if !codescan.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := csgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := csgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (csgb *CodeScanGroupBy) sqlQuery() *sql.Selector {
	selector := csgb.sql.Select()
	aggregation := make([]string, 0, len(csgb.fns))
	for _, fn := range csgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(csgb.fields)+len(csgb.fns))
		for _, f := range csgb.fields {
			columns = append(columns, selector.C(f))
		}
		for _, c := range aggregation {
			columns = append(columns, c)
		}
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(csgb.fields...)...)
}

// CodeScanSelect is the builder for selecting fields of CodeScan entities.
type CodeScanSelect struct {
	*CodeScanQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (css *CodeScanSelect) Scan(ctx context.Context, v interface{}) error {
	if err := css.prepareQuery(ctx); err != nil {
		return err
	}
	css.sql = css.CodeScanQuery.sqlQuery(ctx)
	return css.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (css *CodeScanSelect) ScanX(ctx context.Context, v interface{}) {
	if err := css.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (css *CodeScanSelect) Strings(ctx context.Context) ([]string, error) {
	if len(css.fields) > 1 {
		return nil, errors.New("ent: CodeScanSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := css.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (css *CodeScanSelect) StringsX(ctx context.Context) []string {
	v, err := css.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (css *CodeScanSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = css.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{codescan.Label}
	default:
		err = fmt.Errorf("ent: CodeScanSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (css *CodeScanSelect) StringX(ctx context.Context) string {
	v, err := css.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (css *CodeScanSelect) Ints(ctx context.Context) ([]int, error) {
	if len(css.fields) > 1 {
		return nil, errors.New("ent: CodeScanSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := css.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (css *CodeScanSelect) IntsX(ctx context.Context) []int {
	v, err := css.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (css *CodeScanSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = css.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{codescan.Label}
	default:
		err = fmt.Errorf("ent: CodeScanSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (css *CodeScanSelect) IntX(ctx context.Context) int {
	v, err := css.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (css *CodeScanSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(css.fields) > 1 {
		return nil, errors.New("ent: CodeScanSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := css.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (css *CodeScanSelect) Float64sX(ctx context.Context) []float64 {
	v, err := css.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (css *CodeScanSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = css.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{codescan.Label}
	default:
		err = fmt.Errorf("ent: CodeScanSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (css *CodeScanSelect) Float64X(ctx context.Context) float64 {
	v, err := css.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (css *CodeScanSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(css.fields) > 1 {
		return nil, errors.New("ent: CodeScanSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := css.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (css *CodeScanSelect) BoolsX(ctx context.Context) []bool {
	v, err := css.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (css *CodeScanSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = css.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{codescan.Label}
	default:
		err = fmt.Errorf("ent: CodeScanSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (css *CodeScanSelect) BoolX(ctx context.Context) bool {
	v, err := css.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (css *CodeScanSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := css.sql.Query()
	if err := css.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
