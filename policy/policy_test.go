package policy

import (
	"testing"
)

type fakeResolver struct{}

// TODO: implement the resolver to work with tests...

func TestEvaluate(t *testing.T) {
	// ctx := context.Background()
	// bCtx := env.NewBubblyContext()
	// s, err := store.New(bCtx)
	// require.NoError(t, err)
	// client := s.Client()
	// {
	// 	err := s.PopulateStoreWithDummyData()
	// 	require.NoError(t, err)
	// }
	// dbRelease, err := client.Release.Query().Where(release.HasHeadOf()).WithHeadOf().Only(ctx)
	// require.NoError(t, err)

	// module := `
	// package bubbly

	// violation[{"msg": msg, "severity": severity}] {
	// 	# high_issues := get_severity_high_issues(input.issues)
	// 	high_issues := [issue | issue := input.issues[i]; issue.severity == "high"]
	// 	count(high_issues) > 0
	// 	msg := sprintf("%d high issue(s)", [count(high_issues)])
	// 	severity := "error"
	// 	# some i
	// 	# input.issues[i].severity == "high"
	// 	# msg := input.issues[i].id
	// 	# severity := sprintf("warning %d", [i])
	// }

	// get_severity_high_issues(issues) = result {
	// 	result = [issue | issue := issues[i]; issue.severity == "high"]
	// }

	// `

	// dbPolicy, err := client.ReleasePolicy.Create().
	// 	AddRepoIDs(dbRelease.Edges.HeadOf.ID).
	// 	SetName("code_issues_severity_high-asd asdasd").
	// 	SetInput(releasepolicy.InputCodeIssues).
	// 	SetModule(module).
	// 	Save(ctx)
	// require.NoError(t, err)

	// p := ent.NewReleasePolicyModelRead().FromEnt(dbPolicy)

	// evalErr := EvaluatePolicy(client, dbRelease.ID, p)
	// require.NoError(t, evalErr)
}
