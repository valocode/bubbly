package parser

import (
	"fmt"
	"testing"

	"github.com/verifa/bubbly/env"

	"github.com/stretchr/testify/assert"

	"github.com/verifa/bubbly/api/core"
	v1 "github.com/verifa/bubbly/api/v1"
	"github.com/zclconf/go-cty/cty"
)

// TestExtractJSONConversion tests that the JSON representation of HCL
// extract resource is correct and that converted Resources match what is
// expected.
func TestExtractJSONConversion(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		desc     string
		input    string
		resource core.ResourceBlock
		expected map[string]interface{}
	}{
		{
			desc:  "basic JSON conversion for junit extract",
			input: "testdata/extracts/junit-extract.bubbly",
			expected: map[string]interface{}{
				"resourceJSON": string(`{"resource":{"extract":{"junit-extract":{"api_version":"v1","spec":{"input":{"file":{}},"source":{"file":"${self.input.file}","format":"object({testsuites=object({duration=number,testsuite=list(object({failures=number,name=string,package=string,testcase=list(object({classname=string,name=string,time=number})),tests=number,time=number}))})})"},"type":"xml"}}}}}`),
				"resource": &v1.Extract{
					ResourceBlock: &core.ResourceBlock{
						ResourceKind:       "extract",
						ResourceName:       "junit-extract",
						ResourceAPIVersion: "v1",
					},
				},
			},
		},
		{
			desc:  "basic JSON conversion for sonarqube extract",
			input: "testdata/extracts/sonarqube-extract.bubbly",
			expected: map[string]interface{}{
				"resourceJSON": string(`{"resource":{"extract":{"sonarqube-extract":{"api_version":"v1","spec":{"input":{"file":{}},"source":{"file":"${self.input.file}","format":"object({issues=list(object({engineId=string,primaryLocation=object({filePath=string,message=string,textRange=object({endColumn=number,endLine=number,startColumn=number,startLine=number})}),ruleId=string,severity=string,type=string}))})"},"type":"json"}}}}}`),
				"resource": &v1.Extract{
					ResourceBlock: &core.ResourceBlock{
						ResourceKind:       "extract",
						ResourceName:       "sonarqube-extract",
						ResourceAPIVersion: "v1",
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()
			p, err := NewParserFromFilename(tc.input)
			assert.NoError(t, err, fmt.Errorf("Failed to create parser: %w", err))

			err = p.Parse(bCtx)

			assert.NoError(t, err, fmt.Errorf("Failed to decode parser: %w", err))

			// create a new parser to load the JSON resources into
			p2 := newParser(nil, nil)

			for _, resMap := range p.Resources {
				for _, resource := range resMap {
					t.Logf("Converting resource %s to JSON", resource.String())
					bJSON, err := resource.JSON(p.Context(cty.NilVal))

					assert.NoError(t, err, fmt.Errorf("Failed to convert to resource to JSON %s: %w", resource.String(), err))

					t.Logf("Resource %s JSON representation: %s", resource.String(), bJSON)

					assert.Equal(t, tc.expected["resourceJSON"], string(bJSON))

					_, err = p2.JSONToResource(bCtx, bJSON)

					assert.NoError(t, err, fmt.Errorf("Failed to convert to json to resource %s: %w", resource.String(), err))

					// Now let's evaluate the resource
					expectedExtract := tc.expected["resource"].(*v1.Extract)
					assert.Equal(t, expectedExtract.ResourceKind, string(resource.Kind()))
					assert.Equal(t, expectedExtract.ResourceAPIVersion, resource.APIVersion())
					assert.Equal(t, expectedExtract.ResourceName, resource.Name())
					assert.Equal(t, expectedExtract.String(), resource.String())

					// Now let's evaluate the underlying ResourceBlock
					actualExtract := resource.(*v1.Extract)

					assert.Equal(t, expectedExtract.ResourceBlock.Kind(), actualExtract.ResourceBlock.Kind())
					assert.Equal(t, expectedExtract.ResourceBlock.APIVersion(), actualExtract.ResourceBlock.APIVersion())
					assert.Equal(t, expectedExtract.ResourceBlock.Name(), actualExtract.ResourceBlock.Name())
					assert.Equal(t, expectedExtract.ResourceBlock.String(), actualExtract.ResourceBlock.String())

					rbJSON, err := actualExtract.ResourceBlock.JSON(p.Context(cty.NilVal))

					assert.NoError(t, err, fmt.Errorf("Failed to convert %s ResourceBlock to JSON: %w", actualExtract.ResourceBlock.String(), err))

					assert.Equal(t, tc.expected["resourceJSON"], string(rbJSON))

				}
			}

			_, err = p2.GetResource(tc.expected["resource"].(*v1.Extract).Kind(), tc.expected["resource"].(*v1.Extract).Name())

			assert.NoError(t, err, fmt.Errorf("Couldn't get resource %s: %w", tc.resource.String(), err))

		})
	}

}

// TestApplyFromJSONParser tests that a valid HCL pipeline can:
// 1. Be parsed normally from its HCL representation
// 2. Be converted to a valid JSON representation
// 3. be decoded from JSON into Resource instances
// 4. be applied to the bubbly server
func TestApplyFromJSONParser(t *testing.T) {

	t.Parallel()
	tcs := []struct {
		desc      string
		testdata  string
		resources map[string]string
		inputs    map[string]cty.Value
	}{
		{
			desc:     "basic apply from json over junit pipeline",
			testdata: "../bubbly/testdata/junit",
			resources: map[string]string{
				"extract":   "junit-simple",
				"transform": "junit-simple",
			},
			inputs: map[string]cty.Value{
				"extract": cty.ObjectVal(map[string]cty.Value{
					"input": cty.ObjectVal(
						map[string]cty.Value{
							"data": cty.ListVal([]cty.Value{cty.StringVal("WALALALALA")}),
							"file": cty.StringVal("../bubbly/testdata/junit/junit.xml"),
						},
					),
				}),
				"transform": cty.ObjectVal(map[string]cty.Value{
					"input": cty.ObjectVal(
						map[string]cty.Value{
							"data": cty.ListVal([]cty.Value{cty.StringVal("WALALALALA")}),
						},
					),
				}),
				"load": cty.ObjectVal(map[string]cty.Value{
					"input": cty.ObjectVal(
						map[string]cty.Value{
							"data": cty.ListVal([]cty.Value{cty.StringVal("WALALALALA")}),
						},
					),
				}),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()

			// First, verify that the testdata can be parsed "normally"
			p, err := NewParserFromFilename(tc.testdata)
			assert.NoError(t, err, fmt.Errorf("Failed to create parser: %w", err))

			err = p.Parse(bCtx)

			assert.NoError(t, err, fmt.Errorf("Failed to decode parser: %w", err))

			// Next, test that each resource can be converted from HCL -> JSON -> Resource
			p2 := loadJSONResources(t, p, tc.testdata)

			// Finally, test that each resource can be applied given valid inputs
			inputs := tc.inputs["extract"]

			// extract apply

			res, err := p2.GetResource(core.ExtractResourceKind, tc.resources["extract"])

			assert.NoError(t, err, fmt.Errorf("Couldn't get %s resource %s: %w", core.ExtractResourceKind, tc.resources["extract"], err))

			out := res.Apply(bCtx, p2.Context(inputs))

			t.Logf("Resource %s ResourceOutput: %+v", res.String(), out.Output())

			assert.NoError(t, out.Error)

			// transform apply

			inputs = tc.inputs["transform"]

			res, err = p2.GetResource(core.TransformResourceKind, tc.resources["transform"])
			assert.NoError(t, err, fmt.Errorf("Couldn't get %s resource %s: %w", core.TransformResourceKind, tc.resources["transform"], err))

			out = res.Apply(bCtx, p2.Context(inputs))

			t.Logf("Resource %s ResourceOutput: %+v", res.String(), out.Output())

			assert.NoError(t, out.Error)

			// TODO: Figure out load step onwards.

			// load apply

			// inputs = tc.inputs["load"]

			// res, err = p2.GetResource(core.LoadResourceKind, "junit-simple")
			// assert.NoError(t, err, fmt.Errorf("Couldn't get resource %s: %w", "load/junit-simple", err))
			// out = res.Apply(p2.Context(inputs))

			// t.Logf("Resource %s ResourceOutput: %+v", res.String(), out.Output())

			// assert.NoError(t, out.Error)
		})
	}
}

// loadJSONResources is a convenience function for loading bubbly resources
// from HCL -> Resource -> JSON -> Resource
// Usage: when testing the conversion of Resource to JSON and back
// Returns a parser loaded with Resources from the provided path
func loadJSONResources(t *testing.T, p *Parser, path string) *Parser {
	// create a new parser to load the JSON resources located at `path` into
	p2 := newParser(nil, nil)
	for _, resMap := range p.Resources {
		for _, resource := range resMap {
			bCtx := env.NewBubblyContext()
			t.Logf("Converting resource %s to JSON", resource.String())
			bJSON, err := resource.JSON(p.Context(cty.NilVal))

			assert.NoError(t, err, fmt.Errorf("Failed to convert to json for resource %s: %w", resource.String(), err))

			_, err = p2.JSONToResource(bCtx, bJSON)

			assert.NoError(t, err, fmt.Errorf("Failed to convert json to resource %s: %w", resource.String(), err))
		}
	}

	return p2
}
