package v1

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	testDataJSON "github.com/verifa/bubbly/api/v1/testdata/importer/json"
	testDataXML "github.com/verifa/bubbly/api/v1/testdata/importer/xml"
	"github.com/zclconf/go-cty/cty"
)

func TestImporterJSON(t *testing.T) {

	ctyType := testDataJSON.ExpectedType()
	expected := testDataJSON.ExpectedValue()

	source := jsonSource{
		File:   filepath.FromSlash(`testdata/importer/json/sonarqube-example.json`),
		Format: ctyType,
	}

	val, err := source.Resolve()

	assert.Nil(t, err, "failed to Resolve() the importer")
	assert.False(t, val.IsNull(), "the importer returned null type value")
	assert.False(t, val.Equals(expected).False(), "the importer returned unexpected value")
}

// runXMLSubtestHelper is a helper which runs tests for a variety of XML input files
func runXMLSubtestHelper(t *testing.T, xmlFile string, ctyType cty.Type, expected cty.Value) {

	t.Helper()

	source := xmlSource{
		File:   xmlFile,
		Format: ctyType,
	}

	val, err := source.Resolve()

	assert.Nil(t, err, "failed to Resolve() the importer")
	assert.False(t, val.IsNull(), "the importer returned null type value")
	assert.False(t, val.Equals(expected).False(), "the importer returned unexpected value")
}

// The XML format is different from JSON in a way that it
// does not have syntax for lists. So the XML parser does not
// know whether an element is by itself, or it's in a list of length one.
// This information is available only in cty.Type data structure we build from HCL
func TestImporterXML(t *testing.T) {

	ctyType := testDataXML.CtyType()

	// Easiest. Baseline. No "short" lists.
	t.Run("junit0", func(t *testing.T) {

		runXMLSubtestHelper(t,
			filepath.FromSlash(`testdata/importer/xml/junit0.xml`),
			ctyType,
			testDataXML.ExpectedValue0(),
		)
	})

	// Harder. A single "testsuite" element but multiple "testcase" elements therein.
	t.Run("junit1", func(t *testing.T) {

		runXMLSubtestHelper(t,
			filepath.FromSlash(`testdata/importer/xml/junit1.xml`),
			ctyType,
			testDataXML.ExpectedValue1(),
		)
	})

	// Hardest. A single "testsuite" element with a single "testcase" elements within.
	t.Run("junit2", func(t *testing.T) {

		runXMLSubtestHelper(t,
			filepath.FromSlash(`testdata/importer/xml/junit2.xml`),
			ctyType,
			testDataXML.ExpectedValue2(),
		)
	})

}

func TestImporterGit(t *testing.T) {

	source := gitSource{
		Directory: filepath.FromSlash(`testdata/importer/git/repo1.git`),
	}

	expected := cty.ObjectVal(map[string]cty.Value{
		"active_branch": cty.StringVal("master"),
		"branches": cty.ObjectVal(map[string]cty.Value{
			"local":  cty.ListVal([]cty.Value{cty.StringVal("dev"), cty.StringVal("master")}),
			"remote": cty.NullVal(cty.List(cty.String)),
		}),
		"commit_id": cty.StringVal("81411ea85f68f64f727f140400d7107786d93ba4"),
		"is_bare":   cty.True,
		"remotes": cty.ListValEmpty(cty.Object(map[string]cty.Type{
			"name": cty.String,
			"url":  cty.String,
		})),
		"tag": cty.StringVal("kawabunga"),
	})

	val, err := source.Resolve()

	assert.Nil(t, err, "failed to Resolve() the importer")
	assert.False(t, val.IsNull(), "the importer returned null type value")
	assert.False(t, val.Equals(expected).False(), "the importer returned unexpected value")
}
