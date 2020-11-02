package v1

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	testDataJSON "github.com/verifa/bubbly/api/v1/testdata/importer/json"
	testDataXML "github.com/verifa/bubbly/api/v1/testdata/importer/xml"
	"github.com/zclconf/go-cty/cty"
)

func TestImporterJSON(t *testing.T) {

	ctyType := testDataJSON.ExpectedType()
	expVal := testDataJSON.ExpectedValue()

	source := jsonSource{
		File:   "testdata/importer/json/sonarqube-example.json",
		Format: ctyType,
	}

	val, err := source.Resolve()
	if err != nil {
		t.Errorf("Failed to Resolve() JSON importer: %s", err.Error())
	}
	if val.IsNull() {
		t.Errorf("Received Null type value")
	}
	if val.Equals(expVal).False() {
		t.Errorf("JSON Importer returned unexpected value.\n\nExpected:\n\n\t%s\n\nActual:\n\n\t%s",
			expVal.GoString(), val.GoString())
	}

	t.Logf("JSON Importer returned value: %s", val.GoString())
}

// runXMLSubtestHelper is a helper which runs tests for a variety of XML input files
func runXMLSubtestHelper(t *testing.T, xmlFile string, ctyType cty.Type, expected cty.Value) string {
	t.Helper()

	source := xmlSource{
		File:   xmlFile,
		Format: ctyType,
	}

	val, err := source.Resolve()

	if err != nil {
		t.Error(errors.Wrap(err, "failed to Resolve() XML importer"))
		return ""
	}
	if val.IsNull() {
		t.Error(errors.New("source.Resolve() returned Null type value"))
		return ""
	}
	if val.Equals(expected).False() {
		t.Error(errors.New(fmt.Sprintf("XML Importer returned unexpected value,\n\nExpected:\n\n\t%s\n\nActual:\n\n\t%s", expected.GoString(), val.GoString())))
		return ""
	}

	t.Logf("XML Importer returned value: %s", val.GoString())
	return ""
}
func TestImporterXML(t *testing.T) {

	ctyType := testDataXML.ExpectedType()

	// Baseline XML file with no surprises
	t.Run("baseline", func(t *testing.T) {

		xmlFileName := "testdata/importer/xml/junit.xml"
		expected := testDataXML.ExpectedValue()

		runXMLSubtestHelper(t, xmlFileName, ctyType, expected)
	})

	// This XML file has a list of length one, as it has only one instance of "testsuite" element
	// FIXME GitHub issue #28
	t.Run("oneElementList", func(t *testing.T) {

		xmlFileName := "testdata/importer/xml/junit-one-element.xml"
		expected := testDataXML.ExpectedValueOneElement()

		runXMLSubtestHelper(t, xmlFileName, ctyType, expected)
	})
}
