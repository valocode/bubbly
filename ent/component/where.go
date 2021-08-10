// Code generated by entc, DO NOT EDIT.

package component

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/valocode/bubbly/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
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
func IDGT(id int) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Vendor applies equality check predicate on the "vendor" field. It's identical to VendorEQ.
func Vendor(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVendor), v))
	})
}

// Version applies equality check predicate on the "version" field. It's identical to VersionEQ.
func Version(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVersion), v))
	})
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDescription), v))
	})
}

// URL applies equality check predicate on the "url" field. It's identical to URLEQ.
func URL(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldURL), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Component {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Component(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.Component {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Component(func(s *sql.Selector) {
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
func NameGT(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// VendorEQ applies the EQ predicate on the "vendor" field.
func VendorEQ(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVendor), v))
	})
}

// VendorNEQ applies the NEQ predicate on the "vendor" field.
func VendorNEQ(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldVendor), v))
	})
}

// VendorIn applies the In predicate on the "vendor" field.
func VendorIn(vs ...string) predicate.Component {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Component(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldVendor), v...))
	})
}

// VendorNotIn applies the NotIn predicate on the "vendor" field.
func VendorNotIn(vs ...string) predicate.Component {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Component(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldVendor), v...))
	})
}

// VendorGT applies the GT predicate on the "vendor" field.
func VendorGT(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldVendor), v))
	})
}

// VendorGTE applies the GTE predicate on the "vendor" field.
func VendorGTE(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldVendor), v))
	})
}

// VendorLT applies the LT predicate on the "vendor" field.
func VendorLT(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldVendor), v))
	})
}

// VendorLTE applies the LTE predicate on the "vendor" field.
func VendorLTE(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldVendor), v))
	})
}

// VendorContains applies the Contains predicate on the "vendor" field.
func VendorContains(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldVendor), v))
	})
}

// VendorHasPrefix applies the HasPrefix predicate on the "vendor" field.
func VendorHasPrefix(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldVendor), v))
	})
}

// VendorHasSuffix applies the HasSuffix predicate on the "vendor" field.
func VendorHasSuffix(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldVendor), v))
	})
}

// VendorEqualFold applies the EqualFold predicate on the "vendor" field.
func VendorEqualFold(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldVendor), v))
	})
}

// VendorContainsFold applies the ContainsFold predicate on the "vendor" field.
func VendorContainsFold(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldVendor), v))
	})
}

// VersionEQ applies the EQ predicate on the "version" field.
func VersionEQ(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVersion), v))
	})
}

// VersionNEQ applies the NEQ predicate on the "version" field.
func VersionNEQ(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldVersion), v))
	})
}

// VersionIn applies the In predicate on the "version" field.
func VersionIn(vs ...string) predicate.Component {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Component(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldVersion), v...))
	})
}

// VersionNotIn applies the NotIn predicate on the "version" field.
func VersionNotIn(vs ...string) predicate.Component {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Component(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldVersion), v...))
	})
}

// VersionGT applies the GT predicate on the "version" field.
func VersionGT(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldVersion), v))
	})
}

// VersionGTE applies the GTE predicate on the "version" field.
func VersionGTE(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldVersion), v))
	})
}

// VersionLT applies the LT predicate on the "version" field.
func VersionLT(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldVersion), v))
	})
}

// VersionLTE applies the LTE predicate on the "version" field.
func VersionLTE(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldVersion), v))
	})
}

// VersionContains applies the Contains predicate on the "version" field.
func VersionContains(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldVersion), v))
	})
}

// VersionHasPrefix applies the HasPrefix predicate on the "version" field.
func VersionHasPrefix(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldVersion), v))
	})
}

// VersionHasSuffix applies the HasSuffix predicate on the "version" field.
func VersionHasSuffix(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldVersion), v))
	})
}

// VersionEqualFold applies the EqualFold predicate on the "version" field.
func VersionEqualFold(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldVersion), v))
	})
}

// VersionContainsFold applies the ContainsFold predicate on the "version" field.
func VersionContainsFold(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldVersion), v))
	})
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDescription), v))
	})
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDescription), v))
	})
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Component {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Component(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldDescription), v...))
	})
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Component {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Component(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldDescription), v...))
	})
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldDescription), v))
	})
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldDescription), v))
	})
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldDescription), v))
	})
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldDescription), v))
	})
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldDescription), v))
	})
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldDescription), v))
	})
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldDescription), v))
	})
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldDescription)))
	})
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldDescription)))
	})
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldDescription), v))
	})
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldDescription), v))
	})
}

// URLEQ applies the EQ predicate on the "url" field.
func URLEQ(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldURL), v))
	})
}

// URLNEQ applies the NEQ predicate on the "url" field.
func URLNEQ(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldURL), v))
	})
}

// URLIn applies the In predicate on the "url" field.
func URLIn(vs ...string) predicate.Component {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Component(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldURL), v...))
	})
}

// URLNotIn applies the NotIn predicate on the "url" field.
func URLNotIn(vs ...string) predicate.Component {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Component(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldURL), v...))
	})
}

// URLGT applies the GT predicate on the "url" field.
func URLGT(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldURL), v))
	})
}

// URLGTE applies the GTE predicate on the "url" field.
func URLGTE(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldURL), v))
	})
}

// URLLT applies the LT predicate on the "url" field.
func URLLT(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldURL), v))
	})
}

// URLLTE applies the LTE predicate on the "url" field.
func URLLTE(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldURL), v))
	})
}

// URLContains applies the Contains predicate on the "url" field.
func URLContains(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldURL), v))
	})
}

// URLHasPrefix applies the HasPrefix predicate on the "url" field.
func URLHasPrefix(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldURL), v))
	})
}

// URLHasSuffix applies the HasSuffix predicate on the "url" field.
func URLHasSuffix(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldURL), v))
	})
}

// URLIsNil applies the IsNil predicate on the "url" field.
func URLIsNil() predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldURL)))
	})
}

// URLNotNil applies the NotNil predicate on the "url" field.
func URLNotNil() predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldURL)))
	})
}

// URLEqualFold applies the EqualFold predicate on the "url" field.
func URLEqualFold(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldURL), v))
	})
}

// URLContainsFold applies the ContainsFold predicate on the "url" field.
func URLContainsFold(v string) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldURL), v))
	})
}

// HasCves applies the HasEdge predicate on the "cves" edge.
func HasCves() predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CvesTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, CvesTable, CvesPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCvesWith applies the HasEdge predicate on the "cves" edge with a given conditions (other predicates).
func HasCvesWith(preds ...predicate.CVE) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CvesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, CvesTable, CvesPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLicenses applies the HasEdge predicate on the "licenses" edge.
func HasLicenses() predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LicensesTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, LicensesTable, LicensesPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLicensesWith applies the HasEdge predicate on the "licenses" edge with a given conditions (other predicates).
func HasLicensesWith(preds ...predicate.License) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LicensesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, LicensesTable, LicensesPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasUses applies the HasEdge predicate on the "uses" edge.
func HasUses() predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UsesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, UsesTable, UsesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUsesWith applies the HasEdge predicate on the "uses" edge with a given conditions (other predicates).
func HasUsesWith(preds ...predicate.ComponentUse) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UsesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, UsesTable, UsesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Component) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Component) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
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
func Not(p predicate.Component) predicate.Component {
	return predicate.Component(func(s *sql.Selector) {
		p(s.Not())
	})
}