package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/verifa/bubbly/api/core"
	v1 "github.com/verifa/bubbly/api/v1"
	"github.com/zclconf/go-cty/cty"
)

// TestJunitJSON tests that the JSON representation of HCL resources is correct
// and that converted Resources match what is expected.
func TestJSONConversion(t *testing.T) {
	tcs := []struct {
		desc     string
		input    string
		resource core.ResourceBlock
		want     string
		expected map[string]interface{}
	}{
		{
			desc:  "basic importer JSON",
			input: "testdata/importers/junit-importer.bubbly",
			expected: map[string]interface{}{
				"resourceJSON": string(`{"resource":{"importer":{"junit-importer":{"api_version":"v1","spec":{"input":{"file":{}},"source":{"file":"${self.input.file}","format":"object({testsuites=object({duration=number,testsuite=list(object({failures=number,name=string,package=string,testcase=list(object({classname=string,name=string,time=number})),tests=number,time=number}))})})"},"type":"xml"}}}}}`),
				"resource": &v1.Importer{
					ResourceBlock: &core.ResourceBlock{
						ResourceKind:       "importer",
						ResourceName:       "junit-importer",
						ResourceAPIVersion: "v1",
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			p, err := NewParserFromFilename(tc.input)
			assert.NoError(t, err, fmt.Errorf("Failed to create parser: %w", err))

			err = p.Parse()

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

					_, err = p2.JSONToResource(bJSON)

					assert.NoError(t, err, fmt.Errorf("Failed to convert to json to resource %s: %w", resource.String(), err))

					// Now let's evaluate the resource
					expectedImporter := tc.expected["resource"].(*v1.Importer)
					assert.Equal(t, expectedImporter.ResourceKind, string(resource.Kind()))
					assert.Equal(t, expectedImporter.ResourceAPIVersion, resource.APIVersion())
					assert.Equal(t, expectedImporter.ResourceName, resource.Name())
					assert.Equal(t, expectedImporter.String(), resource.String())

					// Now let's evaluate the underlying ResourceBlock
					actualImporter := resource.(*v1.Importer)

					assert.Equal(t, expectedImporter.ResourceBlock.Kind(), actualImporter.ResourceBlock.Kind())
					assert.Equal(t, expectedImporter.ResourceBlock.APIVersion(), actualImporter.ResourceBlock.APIVersion())
					assert.Equal(t, expectedImporter.ResourceBlock.Name(), actualImporter.ResourceBlock.Name())
					assert.Equal(t, expectedImporter.ResourceBlock.String(), actualImporter.ResourceBlock.String())

					rbJSON, err := actualImporter.ResourceBlock.JSON(p.Context(cty.NilVal))

					assert.NoError(t, err, fmt.Errorf("Failed to convert %s ResourceBlock to JSON: %w", actualImporter.ResourceBlock.String(), err))

					assert.Equal(t, tc.expected["resourceJSON"], string(rbJSON))

				}
			}

			_, err = p2.GetResource(tc.expected["resource"].(*v1.Importer).Kind(), tc.expected["resource"].(*v1.Importer).Name())

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
	t.Run("test parser.JSON for all parser/testdata files", func(t *testing.T) {
		t.Parallel()

		// First, verify that the testdata can be parsed "normally"
		p, err := NewParserFromFilename("testdata/example-full-pipeline")
		assert.NoError(t, err, fmt.Errorf("Failed to create parser: %w", err))

		err = p.Parse()

		assert.NoError(t, err, fmt.Errorf("Failed to decode parser: %w", err))

		// Next, test that each resource can be converted from HCL -> JSON -> Resource
		p2 := loadJSONResources(t, p, "testdata/example-full-pipeline")

		// Finally, test that each resource can be applied given valid inputs
		inputs := cty.ObjectVal(map[string]cty.Value{
			"input": cty.ObjectVal(
				map[string]cty.Value{
					"data": cty.ListVal([]cty.Value{cty.StringVal("WALALALALA")}),
					"file": cty.StringVal("./testdata/example-full-pipeline/junit.xml"),
				},
			),
		})

		// importer apply

		res, err := p2.GetResource(core.ImporterResourceKind, "junit-simple")

		assert.NoError(t, err, fmt.Errorf("Couldn't get resource %s: %w", "importer/junit-simple", err))

		out := res.Apply(p2.Context(inputs))

		t.Logf("Resource %s ResourceOutput: %+v", res.String(), out.Output())

		assert.NoError(t, out.Error)

		// translator apply

		res, err = p2.GetResource(core.TranslatorResourceKind, "junit-simple")
		assert.NoError(t, err, fmt.Errorf("Couldn't get resource %s: %w", "translator/junit-simple", err))
		out = res.Apply(p2.Context(inputs))

		t.Logf("Resource %s ResourceOutput: %+v", res.String(), out.Output())

		assert.NoError(t, out.Error)

		// TODO: Figure out publish step onwards.

		// publish apply

		// res, err = p2.GetResource(core.PublishResourceKind, "junit-simple")
		// assert.NoError(t, err, fmt.Errorf("Couldn't get resource %s: %w", "publisher/junit-simple", err))
		// inputs = cty.ObjectVal(map[string]cty.Value{
		// 	"input": cty.ObjectVal(
		// 		map[string]cty.Value{
		// 			"data": cty.StringVal(""),
		// 		},
		// 	),
		// })
		// out = res.Apply(p2.Context(inputs))

		// t.Logf("Resource %s ResourceOutput: %+v", res.String(), out.Output())

		// assert.NoError(t, out.Error)

		// for _, resMap := range p.Resources {
		// 	for _, res := range resMap {
		// 		out := res.Apply(p.Context(inputs))

		// 		t.Logf("Resource %s ResourceOutput: %+v", res.String(), out.Output())

		// 		assert.NoError(t, out.Error)
		// 	}
		// }

		// res, err := p.GetResource(core.TranslatorResourceKind, "junit")
		// if err != nil {
		// 	t.Errorf("Couldnt get resource: %s", "junit")
		// 	t.FailNow()
		// }

		// if out := res.Apply(p.Context(inputs)); out.Error != nil {
		// 	t.Errorf("Failed to decode translator at the end: %s", out.Error.Error())
		// }

	})

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
			t.Logf("Converting resource %s to JSON", resource.String())
			bJSON, err := resource.JSON(p.Context(cty.NilVal))

			assert.NoError(t, err, fmt.Errorf("Failed to convert to json for resource %s: %w", resource.String(), err))

			_, err = p2.JSONToResource(bJSON)

			assert.NoError(t, err, fmt.Errorf("Failed to convert json to resource %s: %w", resource.String(), err))
		}
	}

	return p2
}
