// Code generated by entc, DO NOT EDIT.

package releaselicense

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/valocode/bubbly/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
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
func IDGT(id int) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// HasLicense applies the HasEdge predicate on the "license" edge.
func HasLicense() predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LicenseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LicenseTable, LicenseColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLicenseWith applies the HasEdge predicate on the "license" edge with a given conditions (other predicates).
func HasLicenseWith(preds ...predicate.License) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LicenseInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LicenseTable, LicenseColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasComponent applies the HasEdge predicate on the "component" edge.
func HasComponent() predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ComponentTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ComponentTable, ComponentColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasComponentWith applies the HasEdge predicate on the "component" edge with a given conditions (other predicates).
func HasComponentWith(preds ...predicate.ReleaseComponent) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ComponentInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ComponentTable, ComponentColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRelease applies the HasEdge predicate on the "release" edge.
func HasRelease() predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ReleaseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ReleaseTable, ReleaseColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasReleaseWith applies the HasEdge predicate on the "release" edge with a given conditions (other predicates).
func HasReleaseWith(preds ...predicate.Release) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
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

// HasScans applies the HasEdge predicate on the "scans" edge.
func HasScans() predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ScansTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ScansTable, ScansColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasScansWith applies the HasEdge predicate on the "scans" edge with a given conditions (other predicates).
func HasScansWith(preds ...predicate.CodeScan) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ScansInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ScansTable, ScansColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ReleaseLicense) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ReleaseLicense) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
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
func Not(p predicate.ReleaseLicense) predicate.ReleaseLicense {
	return predicate.ReleaseLicense(func(s *sql.Selector) {
		p(s.Not())
	})
}
