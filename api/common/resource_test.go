package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valocode/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func TestCompareInputs(t *testing.T) {
	tests := []struct {
		name          string
		decls         core.InputDeclarations
		inputs        cty.Value
		expectError   bool
		expectedValue cty.Value
	}{
		{
			name: "basic test",
			decls: core.InputDeclarations{
				&core.InputDeclaration{Name: "input1"},
			},
			inputs: cty.ObjectVal(map[string]cty.Value{
				"input": cty.ObjectVal(map[string]cty.Value{
					"input1": cty.StringVal("empty"),
				}),
			}),
			expectError: false,
			expectedValue: cty.ObjectVal(map[string]cty.Value{
				"input": cty.ObjectVal(map[string]cty.Value{
					"input1": cty.StringVal("empty"),
				}),
			}),
		},
		{
			name: "use defaults test",
			decls: core.InputDeclarations{
				&core.InputDeclaration{Name: "input1", Default: cty.StringVal("empty")},
			},
			inputs:      cty.EmptyObjectVal,
			expectError: false,
			expectedValue: cty.ObjectVal(map[string]cty.Value{
				"input": cty.ObjectVal(map[string]cty.Value{
					"input1": cty.StringVal("empty"),
				}),
			}),
		},
		{
			name: "expect error",
			decls: core.InputDeclarations{
				&core.InputDeclaration{Name: "input1"},
			},
			inputs:        cty.EmptyObjectVal,
			expectError:   true,
			expectedValue: cty.NilVal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retInputs, err := compareInputsWithDecls(tt.decls, tt.inputs)
			assert.Equalf(t, tt.expectedValue, retInputs, "returned inputs did not match expected inputs")
			if tt.expectError {
				assert.Errorf(t, err, "expected an error but did not receive one")
			} else {
				assert.NoErrorf(t, err, "unexpected error")
			}
		})
	}
}
