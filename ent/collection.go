// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (a *AdapterQuery) CollectFields(ctx context.Context, satisfies ...string) *AdapterQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		a = a.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return a
}

func (a *AdapterQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *AdapterQuery {
	return a
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (a *ArtifactQuery) CollectFields(ctx context.Context, satisfies ...string) *ArtifactQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		a = a.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return a
}

func (a *ArtifactQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ArtifactQuery {
	return a
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ci *CodeIssueQuery) CollectFields(ctx context.Context, satisfies ...string) *CodeIssueQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		ci = ci.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return ci
}

func (ci *CodeIssueQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *CodeIssueQuery {
	return ci
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (cs *CodeScanQuery) CollectFields(ctx context.Context, satisfies ...string) *CodeScanQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		cs = cs.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return cs
}

func (cs *CodeScanQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *CodeScanQuery {
	return cs
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (c *ComponentQuery) CollectFields(ctx context.Context, satisfies ...string) *ComponentQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		c = c.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return c
}

func (c *ComponentQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ComponentQuery {
	return c
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (gc *GitCommitQuery) CollectFields(ctx context.Context, satisfies ...string) *GitCommitQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		gc = gc.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return gc
}

func (gc *GitCommitQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *GitCommitQuery {
	return gc
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (l *LicenseQuery) CollectFields(ctx context.Context, satisfies ...string) *LicenseQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		l = l.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return l
}

func (l *LicenseQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *LicenseQuery {
	return l
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (lu *LicenseUseQuery) CollectFields(ctx context.Context, satisfies ...string) *LicenseUseQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		lu = lu.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return lu
}

func (lu *LicenseUseQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *LicenseUseQuery {
	return lu
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pr *ProjectQuery) CollectFields(ctx context.Context, satisfies ...string) *ProjectQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		pr = pr.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return pr
}

func (pr *ProjectQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ProjectQuery {
	return pr
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (r *ReleaseQuery) CollectFields(ctx context.Context, satisfies ...string) *ReleaseQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		r = r.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return r
}

func (r *ReleaseQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ReleaseQuery {
	return r
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (rc *ReleaseComponentQuery) CollectFields(ctx context.Context, satisfies ...string) *ReleaseComponentQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		rc = rc.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return rc
}

func (rc *ReleaseComponentQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ReleaseComponentQuery {
	return rc
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (re *ReleaseEntryQuery) CollectFields(ctx context.Context, satisfies ...string) *ReleaseEntryQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		re = re.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return re
}

func (re *ReleaseEntryQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ReleaseEntryQuery {
	return re
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (rl *ReleaseLicenseQuery) CollectFields(ctx context.Context, satisfies ...string) *ReleaseLicenseQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		rl = rl.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return rl
}

func (rl *ReleaseLicenseQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ReleaseLicenseQuery {
	return rl
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (rp *ReleasePolicyQuery) CollectFields(ctx context.Context, satisfies ...string) *ReleasePolicyQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		rp = rp.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return rp
}

func (rp *ReleasePolicyQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ReleasePolicyQuery {
	return rp
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (rpv *ReleasePolicyViolationQuery) CollectFields(ctx context.Context, satisfies ...string) *ReleasePolicyViolationQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		rpv = rpv.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return rpv
}

func (rpv *ReleasePolicyViolationQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ReleasePolicyViolationQuery {
	return rpv
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (rv *ReleaseVulnerabilityQuery) CollectFields(ctx context.Context, satisfies ...string) *ReleaseVulnerabilityQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		rv = rv.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return rv
}

func (rv *ReleaseVulnerabilityQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ReleaseVulnerabilityQuery {
	return rv
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (r *RepoQuery) CollectFields(ctx context.Context, satisfies ...string) *RepoQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		r = r.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return r
}

func (r *RepoQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *RepoQuery {
	return r
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (tc *TestCaseQuery) CollectFields(ctx context.Context, satisfies ...string) *TestCaseQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		tc = tc.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return tc
}

func (tc *TestCaseQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *TestCaseQuery {
	return tc
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (tr *TestRunQuery) CollectFields(ctx context.Context, satisfies ...string) *TestRunQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		tr = tr.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return tr
}

func (tr *TestRunQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *TestRunQuery {
	return tr
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (v *VulnerabilityQuery) CollectFields(ctx context.Context, satisfies ...string) *VulnerabilityQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		v = v.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return v
}

func (v *VulnerabilityQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *VulnerabilityQuery {
	return v
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (vr *VulnerabilityReviewQuery) CollectFields(ctx context.Context, satisfies ...string) *VulnerabilityReviewQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		vr = vr.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return vr
}

func (vr *VulnerabilityReviewQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *VulnerabilityReviewQuery {
	return vr
}
