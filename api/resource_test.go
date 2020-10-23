package api

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/hashicorp/hcl/v2/hclparse"
// 	"github.com/verifa/bubbly/parser"
// 	"github.com/zclconf/go-cty/cty"
// )

// func TestResource(t *testing.T) {
// 	hclString := `
// resource "upload" "sonarqube" {
// 	api_version = "v1"
// 	spec {
// 		contents = translator.sonarqube.output
// 	}
// }

// resource "upload" "git" {
// 	api_version = "v1"
// 	spec {
// 		contents = translator.git.output
// 	}
// }

// // stage build
// bubbly apply -f . --filter upload.git

// // stage sonarqube
// bubbly apply -f . --filter upload.sonarqube

// resource "importer" "default" {
// 	api_version = "v1"
// 	spec {
// 		type = "xml"
// 		source {
// 			file = "testdata/importer/junit.xml"
// 		}
// 		format = object({
// 			testsuites: object({
// 				duration: number,
// 				testsuite: list(object({
// 					failures: number,
// 					name: string,
// 					package: string,
// 					tests: number,
// 					time: number,
// 					testcase: list(object({
// 						classname: string
// 						name: string
// 						time: number
// 					}))
// 				}))
// 			})
// 		})
// 	}
// }
// 	`
// 	main := &HCLMainType{}
// 	processHCL(t, hclString, main)

// 	for _, res := range main.ResourceBlocks {
// 		traversals := res.Resource.Variables()
// 		fmt.Printf("NEED TO RESOLVE: %v\n", traversals)
// 	}

// }

// func processHCL(t *testing.T, src string, val interface{}) {
// 	processHCLWithInputs(t, src, val, cty.NilVal)
// }

// func processHCLWithInputs(t *testing.T, src string, val interface{}, inputs cty.Value) {
// 	parser := hclparse.NewParser()
// 	file, diags := parser.ParseHCL([]byte(src), "test-file")
// 	if diags.HasErrors() {
// 		t.Errorf("Failed to parse HCL: %s", diags.Error())
// 	}
// 	s := parser.NewScope()
// 	if !inputs.IsNull() {
// 		s.SetInputs(inputs)
// 	}
// 	diags = s.Decode(file.Body, val)
// 	if diags.HasErrors() {
// 		t.Errorf("Failed to process HCL: %s", diags.Error())
// 		for i, d := range diags {
// 			t.Errorf("Diag (%d): %s", i, d.Error())
// 		}
// 		t.FailNow()
// 	}
// }
