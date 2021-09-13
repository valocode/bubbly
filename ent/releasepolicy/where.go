// Code generated by entc, DO NOT EDIT.

package releasepolicy

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/valocode/bubbly/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
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
func IDGT(id int) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Module applies equality check predicate on the "module" field. It's identical to ModuleEQ.
func Module(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldModule), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.ReleasePolicy {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.ReleasePolicy {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// ModuleEQ applies the EQ predicate on the "module" field.
func ModuleEQ(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldModule), v))
	})
}

// ModuleNEQ applies the NEQ predicate on the "module" field.
func ModuleNEQ(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldModule), v))
	})
}

// ModuleIn applies the In predicate on the "module" field.
func ModuleIn(vs ...string) predicate.ReleasePolicy {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldModule), v...))
	})
}

// ModuleNotIn applies the NotIn predicate on the "module" field.
func ModuleNotIn(vs ...string) predicate.ReleasePolicy {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldModule), v...))
	})
}

// ModuleGT applies the GT predicate on the "module" field.
func ModuleGT(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldModule), v))
	})
}

// ModuleGTE applies the GTE predicate on the "module" field.
func ModuleGTE(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldModule), v))
	})
}

// ModuleLT applies the LT predicate on the "module" field.
func ModuleLT(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldModule), v))
	})
}

// ModuleLTE applies the LTE predicate on the "module" field.
func ModuleLTE(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldModule), v))
	})
}

// ModuleContains applies the Contains predicate on the "module" field.
func ModuleContains(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldModule), v))
	})
}

// ModuleHasPrefix applies the HasPrefix predicate on the "module" field.
func ModuleHasPrefix(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldModule), v))
	})
}

// ModuleHasSuffix applies the HasSuffix predicate on the "module" field.
func ModuleHasSuffix(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldModule), v))
	})
}

// ModuleEqualFold applies the EqualFold predicate on the "module" field.
func ModuleEqualFold(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldModule), v))
	})
}

// ModuleContainsFold applies the ContainsFold predicate on the "module" field.
func ModuleContainsFold(v string) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldModule), v))
	})
}

// HasOwner applies the HasEdge predicate on the "owner" edge.
func HasOwner() predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(OwnerTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, OwnerTable, OwnerColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOwnerWith applies the HasEdge predicate on the "owner" edge with a given conditions (other predicates).
func HasOwnerWith(preds ...predicate.Organization) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(OwnerInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, OwnerTable, OwnerColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProjects applies the HasEdge predicate on the "projects" edge.
func HasProjects() predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectsTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ProjectsTable, ProjectsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectsWith applies the HasEdge predicate on the "projects" edge with a given conditions (other predicates).
func HasProjectsWith(preds ...predicate.Project) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ProjectsTable, ProjectsPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRepos applies the HasEdge predicate on the "repos" edge.
func HasRepos() predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ReposTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ReposTable, ReposPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasReposWith applies the HasEdge predicate on the "repos" edge with a given conditions (other predicates).
func HasReposWith(preds ...predicate.Repo) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ReposInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ReposTable, ReposPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasViolations applies the HasEdge predicate on the "violations" edge.
func HasViolations() predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ViolationsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ViolationsTable, ViolationsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasViolationsWith applies the HasEdge predicate on the "violations" edge with a given conditions (other predicates).
func HasViolationsWith(preds ...predicate.ReleasePolicyViolation) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ViolationsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ViolationsTable, ViolationsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ReleasePolicy) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ReleasePolicy) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
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
func Not(p predicate.ReleasePolicy) predicate.ReleasePolicy {
	return predicate.ReleasePolicy(func(s *sql.Selector) {
		p(s.Not())
	})
}
