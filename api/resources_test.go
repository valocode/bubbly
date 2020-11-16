package api

import (
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/verifa/bubbly/api/core"
	v1 "github.com/verifa/bubbly/api/v1"

	"github.com/stretchr/testify/assert"
)

// Tests api.NewResources
func TestNewResources(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		desc     string
		expected *Resources
	}{
		{
			desc: "base set up of NewResources",
			expected: &Resources{
				"importer":     map[string]core.Resource{},
				"pipeline":     map[string]core.Resource{},
				"pipeline_run": map[string]core.Resource{},
				"publish":      map[string]core.Resource{},
				"translator":   map[string]core.Resource{},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			resources := NewResources()

			assert.NotNil(t, resources)

			assert.Equal(t, tc.expected, resources)
		})
	}
}

// Tests api.NewResourcesFromBlocks
func TestNewResourcesFromBlocks(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		desc              string
		input             core.ResourceBlocks
		expectedResource  core.Resource
		expectedResources *Resources
		expectedSuccess   bool
	}{
		{
			desc: "basic importer",
			input: core.ResourceBlocks{
				&core.ResourceBlock{
					ResourceKind:       "importer",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
			},
			expectedResource: &v1.Importer{
				ResourceBlock: &core.ResourceBlock{
					ResourceKind:       "importer",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
			},
			expectedResources: &Resources{
				"importer": map[string]core.Resource{
					"sonarqube": &v1.Importer{
						ResourceBlock: &core.ResourceBlock{
							ResourceKind:       "importer",
							ResourceName:       "sonarqube",
							ResourceAPIVersion: "v1",
						},
					},
				},
				"pipeline":     map[string]core.Resource{},
				"pipeline_run": map[string]core.Resource{},
				"publish":      map[string]core.Resource{},
				"translator":   map[string]core.Resource{},
			},
		},
		{
			desc: "basic all resource types",
			input: core.ResourceBlocks{
				&core.ResourceBlock{
					ResourceKind:       "importer",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
				&core.ResourceBlock{
					ResourceKind:       "translator",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
				&core.ResourceBlock{
					ResourceKind:       "publish",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
				&core.ResourceBlock{
					ResourceKind:       "pipeline",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
				&core.ResourceBlock{
					ResourceKind:       "pipeline_run",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
			},
			expectedResource: &v1.Importer{
				ResourceBlock: &core.ResourceBlock{
					ResourceKind:       "pipeline",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
			},
			expectedResources: &Resources{
				"importer": map[string]core.Resource{
					"sonarqube": &v1.Importer{
						ResourceBlock: &core.ResourceBlock{
							ResourceKind:       "importer",
							ResourceName:       "sonarqube",
							ResourceAPIVersion: "v1",
						},
					},
				},
				"translator": map[string]core.Resource{
					"sonarqube": &v1.Translator{
						ResourceBlock: &core.ResourceBlock{
							ResourceKind:       "translator",
							ResourceName:       "sonarqube",
							ResourceAPIVersion: "v1",
						},
					},
				},
				"pipeline": map[string]core.Resource{
					"sonarqube": &v1.Pipeline{
						ResourceBlock: &core.ResourceBlock{
							ResourceKind:       "pipeline",
							ResourceName:       "sonarqube",
							ResourceAPIVersion: "v1",
						},
						Tasks: core.Tasks{},
					},
				},
				"pipeline_run": map[string]core.Resource{
					"sonarqube": &v1.PipelineRun{
						ResourceBlock: &core.ResourceBlock{
							ResourceKind:       "pipeline_run",
							ResourceName:       "sonarqube",
							ResourceAPIVersion: "v1",
						},
					},
				},
				"publish": map[string]core.Resource{
					"sonarqube": &v1.Publish{
						ResourceBlock: &core.ResourceBlock{
							ResourceKind:       "publish",
							ResourceName:       "sonarqube",
							ResourceAPIVersion: "v1",
						},
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			resources := NewResourcesFromBlocks(tc.input)

			assert.NotNil(t, resources)

			if resources != nil {
				t.Logf("resources: %+v", resources)
			}

			assert.NotNil(t, resources.Get(tc.expectedResource.Kind(), tc.expectedResource.Name()))

			assert.Equal(t, tc.expectedResources, resources)
		})
	}
}

// Tests Resources.Get
func TestGetResource(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		desc      string
		resources *Resources
		input     core.Resource

		expectedSuccess bool
	}{
		{
			desc: "basic Get",
			resources: &Resources{
				"importer": map[string]core.Resource{
					"sonarqube": &v1.Importer{
						ResourceBlock: &core.ResourceBlock{
							ResourceKind:       "importer",
							ResourceName:       "sonarqube",
							ResourceAPIVersion: "v1",
						},
					},
				},
			},
			input: &v1.Importer{
				ResourceBlock: &core.ResourceBlock{
					ResourceKind:       "importer",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
			},
			expectedSuccess: true,
		},
		{
			desc: "basic unsuccessful Get",
			resources: &Resources{
				"importer": map[string]core.Resource{
					"sonarqube": &v1.Importer{
						ResourceBlock: &core.ResourceBlock{
							ResourceKind:       "importer",
							ResourceName:       "sonarqube",
							ResourceAPIVersion: "v1",
						},
					},
				},
			},
			input: &v1.Importer{
				ResourceBlock: &core.ResourceBlock{
					ResourceKind:       "publisher",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
			},
			expectedSuccess: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			r := tc.resources.Get(tc.input.Kind(), tc.input.Name())

			if tc.expectedSuccess {
				assert.NotNil(t, r)
			} else {
				assert.Nil(t, r)
			}
		})
	}
}

// Tests failures cases of Resources.NewResource
func TestNewResourceFails(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		desc          string
		input         core.ResourceBlocks
		expectedError string
	}{
		{
			desc: "basic duplicate resource creation",
			input: core.ResourceBlocks{
				&core.ResourceBlock{
					ResourceKind:       "importer",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
				&core.ResourceBlock{
					ResourceKind:       "importer",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
			},
			expectedError: "resource v1.default.importer.sonarqube already exists",
		},
		{
			desc: "basic unsupported resource creation",
			input: core.ResourceBlocks{
				&core.ResourceBlock{
					ResourceKind:       "destroyer",
					ResourceName:       "sonarqube",
					ResourceAPIVersion: "v1",
				},
			},
			expectedError: "resource not supported: destroyer",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			resources := NewResources()

			assert.NotNil(t, resources)

			var result error
			for _, r := range tc.input {
				_, err := resources.NewResource(r)

				result = multierror.Append(result, err)
			}

			assert.Error(t, result)

			assert.Contains(t, result.Error(), tc.expectedError)
		})
	}
}
