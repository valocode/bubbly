package parser

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

func TestScope(t *testing.T) {
	basicHCLString := `

local "api_version" {
	value = "v1"
}

resource "extract" "junit" {
    api_version = local.api_version
    spec {
        input "file" {}
        type = "xml"
        source {
            file = self.spec.input.file
        }
        format = object({
            testsuites: object({
                duration: number,
                testsuite: list(object({
                    failures: number,
                    name: string,
                    package: string,
                    tests: number,
                    time: number,
                    testcase: list(object({
                        classname: string
                        name: string
                        time: number
                    }))
                }))
            })
        })
    }
}
`
	t.Run("Basic HCL example with optionals", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		b := core.HCLMainType{}
		processHCL(bCtx, t, basicHCLString, &b)
		// assert.Equal(t, b.BasicBlocks[0].FirstLabel, "first_label")
		// assert.Equal(t, b.BasicBlocks[0].SecondLabel, "second_label")
		// assert.Equal(t, b.BasicBlocks[0].Number, 42)
		// assert.Equal(t, b.BasicBlocks[0].String, "spiffing")
		// assert.Equal(t, b.BasicBlocks[0].OptionalString, "")
	})
}

// func TestDynamicComplex(t *testing.T) {
// 	dynamicHCLString := `
// listblock "sunny" "road" {
// 	strlist = ["beavus", "butthead"]
// }

// listblock "sunny" "garden" {
// 	strlist = ["bill", "ted"]
// }

// dynamic "basicblock" {
// 	for_each = listblock.sunny.garden.strlist
// 	iterator = it
// 	labels = [ "character", "${it.value}" ]
// 	content {
// 		number = length(listblock.sunny.garden.strlist)
// 		string = "Rock n roll ${it.value}!"
// 	}
// }

// optblock "testing" {
// 	strlist = [basicblock.character.bill.string, basicblock.character.ted.string]
// }
// 	`

// 	t.Run("Advanved HCL example with dynamic", func(t *testing.T) {
// 		b := BasicHCL{}
// 		processHCL(t, dynamicHCLString, &b)
// 		assert.Equal(t, b.OptionalBlock.Label, "testing")
// 		assert.Equal(t, b.OptionalBlock.StringList[0], "Rock n roll bill!")
// 		assert.Equal(t, b.OptionalBlock.StringList[1], "Rock n roll ted!")
// 	})
// }

// func TestInputs(t *testing.T) {
// 	dynamicHCLString := `

// // mixing up the order and use of input comes before the input
// listblock "sunny" "garden" {
// 	strlist = input.strlist_input
// }

// input "strlist_input" {
// 	description = "A list of strings"
// 	default = ["bill", "ted"]
// 	type = list(string)
// }

// dynamic "basicblock" {
// 	for_each = listblock.sunny.garden.strlist
// 	iterator = it
// 	labels = [ "character", "${it.value}" ]
// 	content {
// 		number = length(listblock.sunny.garden.strlist)
// 		string = "Rock n roll ${it.value}!"
// 	}
// }

// optblock "testing" {
// 	strlist = [basicblock.character.bill.string, basicblock.character.ted.string]
// }
// 	`

// 	t.Run("Test a scope with inputs", func(t *testing.T) {
// 		b := BasicHCL{}
// 		inputs := cty.ObjectVal(map[string]cty.Value{
// 			"strlist_input": cty.ListVal([]cty.Value{cty.StringVal("bill"), cty.StringVal("ted")}),
// 		})
// 		processHCLWithInputs(t, dynamicHCLString, &b, inputs)
// 		assert.Equal(t, b.OptionalBlock.Label, "testing")
// 		assert.Equal(t, b.OptionalBlock.StringList[0], "Rock n roll bill!")
// 		assert.Equal(t, b.OptionalBlock.StringList[1], "Rock n roll ted!")
// 	})
// 	t.Run("Test a scope with default input values", func(t *testing.T) {
// 		b := BasicHCL{}
// 		processHCL(t, dynamicHCLString, &b)
// 		assert.Equal(t, b.OptionalBlock.Label, "testing")
// 		assert.Equal(t, b.OptionalBlock.StringList[0], "Rock n roll bill!")
// 		assert.Equal(t, b.OptionalBlock.StringList[1], "Rock n roll ted!")
// 	})
// }

// func TestOutputs(t *testing.T) {
// 	dynamicHCLString := `

// // mixing up the order and use of input comes before the input
// listblock "sunny" "garden" {
// 	strlist = ["bill", "ted"]
// }

// dynamic "basicblock" {
// 	for_each = listblock.sunny.garden.strlist
// 	iterator = it
// 	labels = [ "character", "${it.value}" ]
// 	content {
// 		number = length(listblock.sunny.garden.strlist)
// 		string = "Rock n roll ${it.value}!"
// 	}
// }

// output "testing" {
// 	value = [basicblock.character.bill.string, basicblock.character.ted.string]
// }
// 	`

// 	t.Run("Test a sub module", func(t *testing.T) {
// 		b := BasicHCL{}
// 		processHCL(t, dynamicHCLString, &b)
// 		assert.Equal(t, b.Outputs[0].Value.Index(cty.NumberIntVal(0)), cty.StringVal("Rock n roll bill!"))
// 		assert.Equal(t, b.Outputs[0].Value.Index(cty.NumberIntVal(1)), cty.StringVal("Rock n roll ted!"))
// 		// assert.Equal(t, b.OptionalBlock.Label, "testing")
// 		// assert.Equal(t, b.OptionalBlock.StringList[0], "Rock n roll bill!")
// 		// assert.Equal(t, b.OptionalBlock.StringList[1], "Rock n roll ted!")
// 	})
// 	println(dynamicHCLString)
// }

// func init() {
// 	// setup our favourite logger
// 	bCtx.Logger.Logger = bCtx.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
// }

func processHCL(bCtx *env.BubblyContext, t *testing.T, src string, val interface{}) {
	processHCLWithInputs(bCtx, t, src, val, cty.NilVal)
}

func processHCLWithInputs(bCtx *env.BubblyContext, t *testing.T, src string, val interface{}, inputs cty.Value) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(src), "test-file")
	if diags.HasErrors() {
		t.Errorf("Failed to parse HCL: %s", diags.Error())
	}
	s := NewScope()
	if !inputs.IsNull() {
		s.SetInputs(inputs)
	}
	if err := s.decodeBody(bCtx, file.Body, val); err != nil {
		t.Errorf("Failed to process HCL: %s", err.Error())
		t.FailNow()

	}
}
