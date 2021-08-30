package adapter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRego(t *testing.T) {
	result, err := RunFromFile(
		"./testdata/adapters/gosec.rego",
		WithInputFiles("./testdata/adapters/gosec.json"),
		WithTracing(true),
	)
	require.NoError(t, err)
	t.Logf("result: %#v", result.CodeScan)
	for _, trace := range result.Traces {
		fmt.Println(trace)
	}
}

// func TestLoader(t *testing.T) {
// 	module := `
// 	package main
// 	a = "abc"

// 	b = "abc"

// 	`
// 	ctx := context.Background()
// 	r := rego.New(
// 		// rego.Query("data.main.a"),
// 		rego.Query("{\"a\": data.main.a, \"b\": data.main.b}"),
// 		rego.Module(adapterModule, module),
// 	)
// 	query, err := r.PrepareForEval(ctx)
// 	require.NoError(t, err)

// 	// Run evaluation
// 	rs, err := query.Eval(ctx)
// 	require.NoError(t, err)

// 	fmt.Printf("resultset: %#v\n\n", rs)
// 	for _, r := range rs {
// 		for _, expr := range r.Expressions {
// 			resultMap, ok := expr.Value.(map[string]interface{})
// 			if !ok {
// 				t.Fatal("not ok")
// 			}
// 			fmt.Println("a: ", resultMap["a"])
// 			fmt.Println("b: ", resultMap["b"])
// 		}
// 	}
// }
