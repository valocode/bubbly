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
	tcs := []struct {
		desc     string
		input    string
		expected ResourceBlock
	}{
		{
			desc:  "basic importer",
			input: "testdata/importer.hcl",
			expected: ResourceBlock{
				ResourceKind:       "importer",
				ResourceName:       "sonarqube",
				ResourceAPIVersion: "v1",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			rb := newResourceBlock(t, tc.input)

			assert.Equal(t, tc.expected.ResourceKind, string(rb.Kind()))
			assert.Equal(t, tc.expected.ResourceAPIVersion, rb.APIVersion())
			assert.Equal(t, tc.expected.ResourceName, rb.Name())
		})
	}
}

// tests the conversion of ResourceBlock to string (ResourceBlock.String())
func TestResourceBlockString(t *testing.T) {
	tcs := []struct {
		desc     string
		input    ResourceBlock
		expected string
	}{
		{
			desc: "basic importer",
			input: ResourceBlock{
				ResourceKind:       "importer",
				ResourceName:       "sonarqube",
				ResourceAPIVersion: "v1",
			},
			expected: "v1.importer.sonarqube",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			sID := tc.input.String()

			assert.Equal(t, tc.expected, sID)
		})
	}
}

// tests core.NewResourceIDFromString()
func TestNewResourceIDFromString(t *testing.T) {
	tcs := []struct {
		desc     string
		input    string
		expected ResourceID
	}{
		{
			desc:  "basic importer",
			input: "importer/sonarqube",
			expected: ResourceID{
				Kind: "importer",
				Name: "sonarqube",
			},
		},
		{
			desc:  "basic importer",
			input: "pipeline/sonarqube",
			expected: ResourceID{
				Kind: "pipeline",
				Name: "sonarqube",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			rID := NewResourceIDFromString(tc.input)

			assert.Equal(t, tc.expected.Kind, rID.Kind)
			assert.Equal(t, tc.expected.Name, rID.Name)
		})
	}
}

// tests ResourceID.String()
func TestString(t *testing.T) {
	tcs := []struct {
		desc     string
		input    ResourceID
		expected string
	}{
		{
			desc: "basic importer",
			input: ResourceID{
				Kind: "importer",
				Name: "sonarqube",
			},
			expected: "importer/sonarqube",
		},
		{
			desc: "basic importer",
			input: ResourceID{
				Kind: "pipeline",
				Name: "sonarqube",
			},
			expected: "pipeline/sonarqube",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			sID := tc.input.String()

			assert.Equal(t, tc.expected, sID)
		})
	}
}
