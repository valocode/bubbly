// Code generated by entc, DO NOT EDIT.

package gitcommit

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/valocode/bubbly/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
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
func IDGT(id int) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Hash applies equality check predicate on the "hash" field. It's identical to HashEQ.
func Hash(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldHash), v))
	})
}

// Branch applies equality check predicate on the "branch" field. It's identical to BranchEQ.
func Branch(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBranch), v))
	})
}

// Tag applies equality check predicate on the "tag" field. It's identical to TagEQ.
func Tag(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTag), v))
	})
}

// Time applies equality check predicate on the "time" field. It's identical to TimeEQ.
func Time(v time.Time) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTime), v))
	})
}

// HashEQ applies the EQ predicate on the "hash" field.
func HashEQ(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldHash), v))
	})
}

// HashNEQ applies the NEQ predicate on the "hash" field.
func HashNEQ(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldHash), v))
	})
}

// HashIn applies the In predicate on the "hash" field.
func HashIn(vs ...string) predicate.GitCommit {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GitCommit(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldHash), v...))
	})
}

// HashNotIn applies the NotIn predicate on the "hash" field.
func HashNotIn(vs ...string) predicate.GitCommit {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GitCommit(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldHash), v...))
	})
}

// HashGT applies the GT predicate on the "hash" field.
func HashGT(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldHash), v))
	})
}

// HashGTE applies the GTE predicate on the "hash" field.
func HashGTE(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldHash), v))
	})
}

// HashLT applies the LT predicate on the "hash" field.
func HashLT(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldHash), v))
	})
}

// HashLTE applies the LTE predicate on the "hash" field.
func HashLTE(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldHash), v))
	})
}

// HashContains applies the Contains predicate on the "hash" field.
func HashContains(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldHash), v))
	})
}

// HashHasPrefix applies the HasPrefix predicate on the "hash" field.
func HashHasPrefix(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldHash), v))
	})
}

// HashHasSuffix applies the HasSuffix predicate on the "hash" field.
func HashHasSuffix(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldHash), v))
	})
}

// HashEqualFold applies the EqualFold predicate on the "hash" field.
func HashEqualFold(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldHash), v))
	})
}

// HashContainsFold applies the ContainsFold predicate on the "hash" field.
func HashContainsFold(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldHash), v))
	})
}

// BranchEQ applies the EQ predicate on the "branch" field.
func BranchEQ(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBranch), v))
	})
}

// BranchNEQ applies the NEQ predicate on the "branch" field.
func BranchNEQ(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldBranch), v))
	})
}

// BranchIn applies the In predicate on the "branch" field.
func BranchIn(vs ...string) predicate.GitCommit {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GitCommit(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldBranch), v...))
	})
}

// BranchNotIn applies the NotIn predicate on the "branch" field.
func BranchNotIn(vs ...string) predicate.GitCommit {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GitCommit(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldBranch), v...))
	})
}

// BranchGT applies the GT predicate on the "branch" field.
func BranchGT(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldBranch), v))
	})
}

// BranchGTE applies the GTE predicate on the "branch" field.
func BranchGTE(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldBranch), v))
	})
}

// BranchLT applies the LT predicate on the "branch" field.
func BranchLT(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldBranch), v))
	})
}

// BranchLTE applies the LTE predicate on the "branch" field.
func BranchLTE(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldBranch), v))
	})
}

// BranchContains applies the Contains predicate on the "branch" field.
func BranchContains(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldBranch), v))
	})
}

// BranchHasPrefix applies the HasPrefix predicate on the "branch" field.
func BranchHasPrefix(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldBranch), v))
	})
}

// BranchHasSuffix applies the HasSuffix predicate on the "branch" field.
func BranchHasSuffix(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldBranch), v))
	})
}

// BranchEqualFold applies the EqualFold predicate on the "branch" field.
func BranchEqualFold(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldBranch), v))
	})
}

// BranchContainsFold applies the ContainsFold predicate on the "branch" field.
func BranchContainsFold(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldBranch), v))
	})
}

// TagEQ applies the EQ predicate on the "tag" field.
func TagEQ(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTag), v))
	})
}

// TagNEQ applies the NEQ predicate on the "tag" field.
func TagNEQ(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTag), v))
	})
}

// TagIn applies the In predicate on the "tag" field.
func TagIn(vs ...string) predicate.GitCommit {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GitCommit(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTag), v...))
	})
}

// TagNotIn applies the NotIn predicate on the "tag" field.
func TagNotIn(vs ...string) predicate.GitCommit {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GitCommit(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTag), v...))
	})
}

// TagGT applies the GT predicate on the "tag" field.
func TagGT(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTag), v))
	})
}

// TagGTE applies the GTE predicate on the "tag" field.
func TagGTE(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTag), v))
	})
}

// TagLT applies the LT predicate on the "tag" field.
func TagLT(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTag), v))
	})
}

// TagLTE applies the LTE predicate on the "tag" field.
func TagLTE(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTag), v))
	})
}

// TagContains applies the Contains predicate on the "tag" field.
func TagContains(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldTag), v))
	})
}

// TagHasPrefix applies the HasPrefix predicate on the "tag" field.
func TagHasPrefix(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldTag), v))
	})
}

// TagHasSuffix applies the HasSuffix predicate on the "tag" field.
func TagHasSuffix(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldTag), v))
	})
}

// TagIsNil applies the IsNil predicate on the "tag" field.
func TagIsNil() predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldTag)))
	})
}

// TagNotNil applies the NotNil predicate on the "tag" field.
func TagNotNil() predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldTag)))
	})
}

// TagEqualFold applies the EqualFold predicate on the "tag" field.
func TagEqualFold(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldTag), v))
	})
}

// TagContainsFold applies the ContainsFold predicate on the "tag" field.
func TagContainsFold(v string) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldTag), v))
	})
}

// TimeEQ applies the EQ predicate on the "time" field.
func TimeEQ(v time.Time) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTime), v))
	})
}

// TimeNEQ applies the NEQ predicate on the "time" field.
func TimeNEQ(v time.Time) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTime), v))
	})
}

// TimeIn applies the In predicate on the "time" field.
func TimeIn(vs ...time.Time) predicate.GitCommit {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GitCommit(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTime), v...))
	})
}

// TimeNotIn applies the NotIn predicate on the "time" field.
func TimeNotIn(vs ...time.Time) predicate.GitCommit {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GitCommit(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTime), v...))
	})
}

// TimeGT applies the GT predicate on the "time" field.
func TimeGT(v time.Time) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTime), v))
	})
}

// TimeGTE applies the GTE predicate on the "time" field.
func TimeGTE(v time.Time) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTime), v))
	})
}

// TimeLT applies the LT predicate on the "time" field.
func TimeLT(v time.Time) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTime), v))
	})
}

// TimeLTE applies the LTE predicate on the "time" field.
func TimeLTE(v time.Time) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTime), v))
	})
}

// HasRepository applies the HasEdge predicate on the "repository" edge.
func HasRepository() predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(RepositoryTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, RepositoryTable, RepositoryColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRepositoryWith applies the HasEdge predicate on the "repository" edge with a given conditions (other predicates).
func HasRepositoryWith(preds ...predicate.Repository) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(RepositoryInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, RepositoryTable, RepositoryColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRelease applies the HasEdge predicate on the "release" edge.
func HasRelease() predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ReleaseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, ReleaseTable, ReleaseColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasReleaseWith applies the HasEdge predicate on the "release" edge with a given conditions (other predicates).
func HasReleaseWith(preds ...predicate.Release) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ReleaseInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, ReleaseTable, ReleaseColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.GitCommit) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.GitCommit) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
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
func Not(p predicate.GitCommit) predicate.GitCommit {
	return predicate.GitCommit(func(s *sql.Selector) {
		p(s.Not())
	})
}
