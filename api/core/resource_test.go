package core

/*
Most resource testing is done within the parser package. This is due to the fact that a ResourceContext, required by many of ResourceBlock's methods, is tightly coupled to the external parser package.
*/

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
)

// newResourceBlock is a helper function for testing the creation of
// ResourceBlocks from hcl files
func newResourceBlock(t *testing.T, filename string) ResourceBlock {
	input, err := os.Open(filename)
	if err != nil {
		t.Error()
	}
	defer input.Close()

	src, err := ioutil.ReadAll(input)
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(src), "test-file")
	if diags.HasErrors() {
		t.Error(diags.Error())
	}
	resWrap := &ResourceBlockHCLWrapper{}

	diags = gohcl.DecodeBody(file.Body, nil, resWrap)

	if diags.HasErrors() {
		t.Error(diags.Error())
	}
	return resWrap.ResourceBlock
}

// tests the conversion of HCL to a ResourceBlock independent of the
// bubbly/parse package
func TestReadHCLToResourceBlock(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		desc     string
		input    string
		expected ResourceBlock
	}{
		{
			desc:  "basic extract",
			input: "testdata/extract.hcl",
			expected: ResourceBlock{
				ResourceKind:       "extract",
				ResourceName:       "sonarqube",
				ResourceAPIVersion: "v1",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			rb := newResourceBlock(t, tc.input)

			assert.Equal(t, tc.expected.ResourceKind, string(rb.Kind()))
			assert.Equal(t, tc.expected.ResourceAPIVersion, rb.APIVersion())
			assert.Equal(t, tc.expected.ResourceName, rb.Name())
		})
	}
}

// tests the conversion of ResourceBlock to string (ResourceBlock.String())
func TestResourceBlockString(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		desc     string
		input    ResourceBlock
		expected string
	}{
		{
			desc: "basic extract",
			input: ResourceBlock{
				ResourceKind:       "extract",
				ResourceName:       "sonarqube",
				ResourceAPIVersion: "v1",
				Metadata: &Metadata{
					Labels: map[string]string{
						"environment": "prod",
					},
					Namespace: "qa",
				},
			},
			expected: "v1.qa.extract.sonarqube",
		},
		{
			desc: "basic extract without namespace",
			input: ResourceBlock{
				ResourceKind:       "extract",
				ResourceName:       "sonarqube",
				ResourceAPIVersion: "v1",
				Metadata: &Metadata{
					Labels: map[string]string{
						"environment": "prod",
					},
				},
			},
			expected: "v1.default.extract.sonarqube",
		},
		{
			desc: "basic extract without metadata",
			input: ResourceBlock{
				ResourceKind:       "extract",
				ResourceName:       "sonarqube",
				ResourceAPIVersion: "v1",
			},
			expected: "v1.default.extract.sonarqube",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			sID := tc.input.String()

			assert.Equal(t, tc.expected, sID)
		})
	}
}

// tests the conversion of ResourceBlock to string (ResourceBlock.String())
func TestResourceBlockLabels(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		desc     string
		input    ResourceBlock
		expected map[string]string
	}{
		{
			desc: "basic extract",
			input: ResourceBlock{
				ResourceKind:       "extract",
				ResourceName:       "sonarqube",
				ResourceAPIVersion: "v1",
				Metadata: &Metadata{
					Labels: map[string]string{
						"environment": "prod",
					},
					Namespace: "qa",
				},
			},
			expected: map[string]string{
				"environment": "prod",
			},
		},
		{
			desc: "basic extract without labels",
			input: ResourceBlock{
				ResourceKind:       "extract",
				ResourceName:       "sonarqube",
				ResourceAPIVersion: "v1",
				Metadata: &Metadata{
					Namespace: "qa",
				},
			},
			expected: nil,
		},
		{
			desc: "basic extract without metadata",
			input: ResourceBlock{
				ResourceKind:       "extract",
				ResourceName:       "sonarqube",
				ResourceAPIVersion: "v1",
			},
			expected: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			l := tc.input.Labels()

			assert.Equal(t, tc.expected, l)
		})
	}
}

// tests core.NewResourceIDFromString()
func TestNewResourceIDFromString(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		desc     string
		input    string
		expected ResourceID
	}{
		{
			desc:  "basic extract without namespace specified",
			input: "extract/sonarqube",
			expected: ResourceID{
				Kind: "extract",
				Name: "sonarqube",
			},
		},
		{
			desc:  "basic pipeline with namespace specified",
			input: "qa/pipeline/sonarqube",
			expected: ResourceID{
				Kind:      "pipeline",
				Name:      "sonarqube",
				Namespace: "qa",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			rID := NewResourceIDFromString(tc.input)

			t.Logf("Resource ID: %s", rID)

			assert.Equal(t, tc.expected.Kind, rID.Kind)
			assert.Equal(t, tc.expected.Name, rID.Name)
			assert.Equal(t, tc.expected.Namespace, rID.Namespace)
		})
	}
}

// tests ResourceID.String()
func TestString(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		desc     string
		input    ResourceID
		expected string
	}{
		{
			desc: "basic extract",
			input: ResourceID{
				Kind:      "extract",
				Name:      "sonarqube",
				Namespace: "qa",
			},
			expected: "qa/extract/sonarqube",
		},
		{
			desc: "basic pipeline",
			input: ResourceID{
				Kind:      "pipeline",
				Name:      "sonarqube",
				Namespace: "qa",
			},
			expected: "qa/pipeline/sonarqube",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			sID := tc.input.String()

			assert.Equal(t, tc.expected, sID)
		})
	}
}
