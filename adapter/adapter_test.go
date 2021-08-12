package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type WInter struct {
	Type  OperationType `json:"-"`
	Inter `json:"operation,omitempty"`
}

func (w *WInter) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &jsonMap)
	if err != nil {
		return err
	}
	switch w.Type {
	case OpJSON:
		w.Inter = &A{}
	default:
		return errors.New("must set type of WInter")
	}
	fmt.Printf("about to unmarshal: %s\n", b)
	if js, ok := jsonMap["operation"]; ok && js != nil {
		err := json.Unmarshal(*js, w.Inter)
		if err != nil {
			return err
		}
	}
	return nil
}

type Inter interface {
	String() string
}

type A struct {
	Field string `json:"field"`
}

func (A) String() string { return "" }

type B struct {
	Channel string `json:"channel"`
}

func (B) String() string { return "" }

func TestJSON(t *testing.T) {
	// var i WInter
	// i.Inter = &A{Field: "value"}
	// b, err := json.MarshalIndent(i, "", "  ")
	// require.NoError(t, err)
	// t.Logf("json: %s", b)

	// {
	// 	var j WInter
	// 	j.Type = OpJSON
	// 	err := json.Unmarshal(b, &j)
	// 	require.NoError(t, err)
	// 	t.Logf("j: %#v", j.Inter)
	// }
	a, err := FromFile("testdata/adapters/gosec.adapt.hcl")
	require.NoError(t, err)
	t.Logf("adapter: %#v", a)

	b, err := json.MarshalIndent(a, "", "  ")
	require.NoError(t, err)
	t.Logf("json: %s", b)
}

// func TestRunAdapter(t *testing.T) {
// 	opts := RunOptions{
// 		Name:     "snyk",
// 		Filename: "testdata/adapters/gosec.adapt.hcl",
// 	}
// 	args := RunArgs{
// 		Filename: "testdata/adapters/gosec.json",
// 	}
// 	_, err := Run(opts, args)
// 	require.NoError(t, err)
// }

func TestAdapter(t *testing.T) {
	adapt, err := FromFile("testdata/adapters/gosec.adapt.hcl")
	require.NoError(t, err)

	args := RunArgs{
		Filename: "testdata/adapters/gosec.json",
	}
	result, err := adapt.Run(args)
	require.NoError(t, err)

	t.Logf("result: %#v", result)
	t.Logf("issues: %#v", result.CodeScan.Issues)
}

func TestTwo(t *testing.T) {
	var a string
	ty := reflect.TypeOf(a)
	// reflect.Zero(ty).String()
	t.Log("this: ", reflect.Zero(ty).Type().String())
}

// func TestAdapter(t *testing.T) {
// 	adapt, err := ParseAdapter("../bubbly/gosec.adapt.hcl")
// 	require.NoError(t, err)

// 	inputs := map[string]cty.Value{
// 		"json": cty.StringVal("../gosec.json"),
// 	}
// 	_, err = adapt.Run(inputs)
// 	require.NoError(t, err)
// }
