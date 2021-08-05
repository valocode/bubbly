// Code generated by entc, DO NOT EDIT.

package cvescan

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/valocode/bubbly/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Tool applies equality check predicate on the "tool" field. It's identical to ToolEQ.
func Tool(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTool), v))
	})
}

// ToolEQ applies the EQ predicate on the "tool" field.
func ToolEQ(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTool), v))
	})
}

// ToolNEQ applies the NEQ predicate on the "tool" field.
func ToolNEQ(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTool), v))
	})
}

// ToolIn applies the In predicate on the "tool" field.
func ToolIn(vs ...string) predicate.CVEScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.CVEScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTool), v...))
	})
}

// ToolNotIn applies the NotIn predicate on the "tool" field.
func ToolNotIn(vs ...string) predicate.CVEScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.CVEScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTool), v...))
	})
}

// ToolGT applies the GT predicate on the "tool" field.
func ToolGT(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTool), v))
	})
}

// ToolGTE applies the GTE predicate on the "tool" field.
func ToolGTE(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTool), v))
	})
}

// ToolLT applies the LT predicate on the "tool" field.
func ToolLT(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTool), v))
	})
}

// ToolLTE applies the LTE predicate on the "tool" field.
func ToolLTE(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTool), v))
	})
}

// ToolContains applies the Contains predicate on the "tool" field.
func ToolContains(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldTool), v))
	})
}

// ToolHasPrefix applies the HasPrefix predicate on the "tool" field.
func ToolHasPrefix(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldTool), v))
	})
}

// ToolHasSuffix applies the HasSuffix predicate on the "tool" field.
func ToolHasSuffix(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldTool), v))
	})
}

// ToolEqualFold applies the EqualFold predicate on the "tool" field.
func ToolEqualFold(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldTool), v))
	})
}

// ToolContainsFold applies the ContainsFold predicate on the "tool" field.
func ToolContainsFold(v string) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldTool), v))
	})
}

// HasRelease applies the HasEdge predicate on the "release" edge.
func HasRelease() predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ReleaseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ReleaseTable, ReleaseColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasReleaseWith applies the HasEdge predicate on the "release" edge with a given conditions (other predicates).
func HasReleaseWith(preds ...predicate.Release) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ReleaseInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ReleaseTable, ReleaseColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEntry applies the HasEdge predicate on the "entry" edge.
func HasEntry() predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EntryTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, EntryTable, EntryColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEntryWith applies the HasEdge predicate on the "entry" edge with a given conditions (other predicates).
func HasEntryWith(preds ...predicate.ReleaseEntry) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EntryInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, EntryTable, EntryColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasVulnerabilities applies the HasEdge predicate on the "vulnerabilities" edge.
func HasVulnerabilities() predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(VulnerabilitiesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, VulnerabilitiesTable, VulnerabilitiesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasVulnerabilitiesWith applies the HasEdge predicate on the "vulnerabilities" edge with a given conditions (other predicates).
func HasVulnerabilitiesWith(preds ...predicate.Vulnerability) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(VulnerabilitiesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, VulnerabilitiesTable, VulnerabilitiesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.CVEScan) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.CVEScan) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.CVEScan) predicate.CVEScan {
	return predicate.CVEScan(func(s *sql.Selector) {
		p(s.Not())
	})
}
