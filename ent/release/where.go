// Code generated by entc, DO NOT EDIT.

package release

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/valocode/bubbly/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
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
func IDGT(id int) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Version applies equality check predicate on the "version" field. It's identical to VersionEQ.
func Version(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVersion), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Release {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Release(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.Release {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Release(func(s *sql.Selector) {
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
func NameGT(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// VersionEQ applies the EQ predicate on the "version" field.
func VersionEQ(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVersion), v))
	})
}

// VersionNEQ applies the NEQ predicate on the "version" field.
func VersionNEQ(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldVersion), v))
	})
}

// VersionIn applies the In predicate on the "version" field.
func VersionIn(vs ...string) predicate.Release {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Release(func(s *sql.Selector) {
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
func VersionNotIn(vs ...string) predicate.Release {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Release(func(s *sql.Selector) {
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
func VersionGT(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldVersion), v))
	})
}

// VersionGTE applies the GTE predicate on the "version" field.
func VersionGTE(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldVersion), v))
	})
}

// VersionLT applies the LT predicate on the "version" field.
func VersionLT(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldVersion), v))
	})
}

// VersionLTE applies the LTE predicate on the "version" field.
func VersionLTE(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldVersion), v))
	})
}

// VersionContains applies the Contains predicate on the "version" field.
func VersionContains(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldVersion), v))
	})
}

// VersionHasPrefix applies the HasPrefix predicate on the "version" field.
func VersionHasPrefix(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldVersion), v))
	})
}

// VersionHasSuffix applies the HasSuffix predicate on the "version" field.
func VersionHasSuffix(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldVersion), v))
	})
}

// VersionEqualFold applies the EqualFold predicate on the "version" field.
func VersionEqualFold(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldVersion), v))
	})
}

// VersionContainsFold applies the ContainsFold predicate on the "version" field.
func VersionContainsFold(v string) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldVersion), v))
	})
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v Status) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), v))
	})
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v Status) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStatus), v))
	})
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...Status) predicate.Release {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Release(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldStatus), v...))
	})
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...Status) predicate.Release {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Release(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldStatus), v...))
	})
}

// HasSubreleases applies the HasEdge predicate on the "subreleases" edge.
func HasSubreleases() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SubreleasesTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, SubreleasesTable, SubreleasesPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSubreleasesWith applies the HasEdge predicate on the "subreleases" edge with a given conditions (other predicates).
func HasSubreleasesWith(preds ...predicate.Release) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, SubreleasesTable, SubreleasesPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasDependencies applies the HasEdge predicate on the "dependencies" edge.
func HasDependencies() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(DependenciesTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, DependenciesTable, DependenciesPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDependenciesWith applies the HasEdge predicate on the "dependencies" edge with a given conditions (other predicates).
func HasDependenciesWith(preds ...predicate.Release) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, DependenciesTable, DependenciesPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCommit applies the HasEdge predicate on the "commit" edge.
func HasCommit() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CommitTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, CommitTable, CommitColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCommitWith applies the HasEdge predicate on the "commit" edge with a given conditions (other predicates).
func HasCommitWith(preds ...predicate.GitCommit) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CommitInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, CommitTable, CommitColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasHeadOf applies the HasEdge predicate on the "head_of" edge.
func HasHeadOf() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(HeadOfTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, HeadOfTable, HeadOfColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasHeadOfWith applies the HasEdge predicate on the "head_of" edge with a given conditions (other predicates).
func HasHeadOfWith(preds ...predicate.Repo) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(HeadOfInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, HeadOfTable, HeadOfColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLog applies the HasEdge predicate on the "log" edge.
func HasLog() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LogTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, LogTable, LogColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLogWith applies the HasEdge predicate on the "log" edge with a given conditions (other predicates).
func HasLogWith(preds ...predicate.ReleaseEntry) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LogInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, LogTable, LogColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasViolations applies the HasEdge predicate on the "violations" edge.
func HasViolations() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ViolationsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ViolationsTable, ViolationsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasViolationsWith applies the HasEdge predicate on the "violations" edge with a given conditions (other predicates).
func HasViolationsWith(preds ...predicate.ReleasePolicyViolation) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
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

// HasArtifacts applies the HasEdge predicate on the "artifacts" edge.
func HasArtifacts() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ArtifactsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ArtifactsTable, ArtifactsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasArtifactsWith applies the HasEdge predicate on the "artifacts" edge with a given conditions (other predicates).
func HasArtifactsWith(preds ...predicate.Artifact) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ArtifactsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ArtifactsTable, ArtifactsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasComponents applies the HasEdge predicate on the "components" edge.
func HasComponents() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ComponentsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ComponentsTable, ComponentsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasComponentsWith applies the HasEdge predicate on the "components" edge with a given conditions (other predicates).
func HasComponentsWith(preds ...predicate.ReleaseComponent) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ComponentsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ComponentsTable, ComponentsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasVulnerabilities applies the HasEdge predicate on the "vulnerabilities" edge.
func HasVulnerabilities() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(VulnerabilitiesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, VulnerabilitiesTable, VulnerabilitiesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasVulnerabilitiesWith applies the HasEdge predicate on the "vulnerabilities" edge with a given conditions (other predicates).
func HasVulnerabilitiesWith(preds ...predicate.ReleaseVulnerability) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
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

// HasCodeScans applies the HasEdge predicate on the "code_scans" edge.
func HasCodeScans() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CodeScansTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, CodeScansTable, CodeScansColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCodeScansWith applies the HasEdge predicate on the "code_scans" edge with a given conditions (other predicates).
func HasCodeScansWith(preds ...predicate.CodeScan) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CodeScansInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, CodeScansTable, CodeScansColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTestRuns applies the HasEdge predicate on the "test_runs" edge.
func HasTestRuns() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TestRunsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, TestRunsTable, TestRunsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTestRunsWith applies the HasEdge predicate on the "test_runs" edge with a given conditions (other predicates).
func HasTestRunsWith(preds ...predicate.TestRun) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TestRunsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, TestRunsTable, TestRunsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasVulnerabilityReviews applies the HasEdge predicate on the "vulnerability_reviews" edge.
func HasVulnerabilityReviews() predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(VulnerabilityReviewsTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, VulnerabilityReviewsTable, VulnerabilityReviewsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasVulnerabilityReviewsWith applies the HasEdge predicate on the "vulnerability_reviews" edge with a given conditions (other predicates).
func HasVulnerabilityReviewsWith(preds ...predicate.VulnerabilityReview) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(VulnerabilityReviewsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, VulnerabilityReviewsTable, VulnerabilityReviewsPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Release) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Release) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
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
func Not(p predicate.Release) predicate.Release {
	return predicate.Release(func(s *sql.Selector) {
		p(s.Not())
	})
}
