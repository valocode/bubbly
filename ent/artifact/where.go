// Code generated by entc, DO NOT EDIT.

package artifact

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/valocode/bubbly/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
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
func IDGT(id int) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Sha256 applies equality check predicate on the "sha256" field. It's identical to Sha256EQ.
func Sha256(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSha256), v))
	})
}

// Time applies equality check predicate on the "time" field. It's identical to TimeEQ.
func Time(v time.Time) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTime), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Artifact {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Artifact(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.Artifact {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Artifact(func(s *sql.Selector) {
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
func NameGT(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// Sha256EQ applies the EQ predicate on the "sha256" field.
func Sha256EQ(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSha256), v))
	})
}

// Sha256NEQ applies the NEQ predicate on the "sha256" field.
func Sha256NEQ(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSha256), v))
	})
}

// Sha256In applies the In predicate on the "sha256" field.
func Sha256In(vs ...string) predicate.Artifact {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Artifact(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSha256), v...))
	})
}

// Sha256NotIn applies the NotIn predicate on the "sha256" field.
func Sha256NotIn(vs ...string) predicate.Artifact {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Artifact(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSha256), v...))
	})
}

// Sha256GT applies the GT predicate on the "sha256" field.
func Sha256GT(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSha256), v))
	})
}

// Sha256GTE applies the GTE predicate on the "sha256" field.
func Sha256GTE(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSha256), v))
	})
}

// Sha256LT applies the LT predicate on the "sha256" field.
func Sha256LT(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSha256), v))
	})
}

// Sha256LTE applies the LTE predicate on the "sha256" field.
func Sha256LTE(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSha256), v))
	})
}

// Sha256Contains applies the Contains predicate on the "sha256" field.
func Sha256Contains(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldSha256), v))
	})
}

// Sha256HasPrefix applies the HasPrefix predicate on the "sha256" field.
func Sha256HasPrefix(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldSha256), v))
	})
}

// Sha256HasSuffix applies the HasSuffix predicate on the "sha256" field.
func Sha256HasSuffix(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldSha256), v))
	})
}

// Sha256EqualFold applies the EqualFold predicate on the "sha256" field.
func Sha256EqualFold(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldSha256), v))
	})
}

// Sha256ContainsFold applies the ContainsFold predicate on the "sha256" field.
func Sha256ContainsFold(v string) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldSha256), v))
	})
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v Type) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldType), v))
	})
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v Type) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldType), v))
	})
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...Type) predicate.Artifact {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Artifact(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldType), v...))
	})
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...Type) predicate.Artifact {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Artifact(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldType), v...))
	})
}

// TimeEQ applies the EQ predicate on the "time" field.
func TimeEQ(v time.Time) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTime), v))
	})
}

// TimeNEQ applies the NEQ predicate on the "time" field.
func TimeNEQ(v time.Time) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTime), v))
	})
}

// TimeIn applies the In predicate on the "time" field.
func TimeIn(vs ...time.Time) predicate.Artifact {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Artifact(func(s *sql.Selector) {
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
func TimeNotIn(vs ...time.Time) predicate.Artifact {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Artifact(func(s *sql.Selector) {
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
func TimeGT(v time.Time) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTime), v))
	})
}

// TimeGTE applies the GTE predicate on the "time" field.
func TimeGTE(v time.Time) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTime), v))
	})
}

// TimeLT applies the LT predicate on the "time" field.
func TimeLT(v time.Time) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTime), v))
	})
}

// TimeLTE applies the LTE predicate on the "time" field.
func TimeLTE(v time.Time) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTime), v))
	})
}

// MetadataIsNil applies the IsNil predicate on the "metadata" field.
func MetadataIsNil() predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldMetadata)))
	})
}

// MetadataNotNil applies the NotNil predicate on the "metadata" field.
func MetadataNotNil() predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldMetadata)))
	})
}

// HasRelease applies the HasEdge predicate on the "release" edge.
func HasRelease() predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ReleaseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ReleaseTable, ReleaseColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasReleaseWith applies the HasEdge predicate on the "release" edge with a given conditions (other predicates).
func HasReleaseWith(preds ...predicate.Release) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
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
func HasEntry() predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EntryTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, EntryTable, EntryColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEntryWith applies the HasEdge predicate on the "entry" edge with a given conditions (other predicates).
func HasEntryWith(preds ...predicate.ReleaseEntry) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
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

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Artifact) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Artifact) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
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
func Not(p predicate.Artifact) predicate.Artifact {
	return predicate.Artifact(func(s *sql.Selector) {
		p(s.Not())
	})
}
