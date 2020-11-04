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
func runXMLSubtestHelper(t *testing.T, xmlFile string, ctyType cty.Type, expected cty.Value) {

	t.Helper()

	source := xmlSource{
		File:   xmlFile,
		Format: ctyType,
	}

	val, err := source.Resolve()

	if err != nil {
		t.Error(errors.Wrap(err, "failed to Resolve() XML importer"))
	}
	if val.IsNull() {
		t.Error(errors.New("source.Resolve() returned Null type value"))
	}
	if val.Equals(expected).False() {
		t.Error(errors.New(fmt.Sprintf("XML Importer returned unexpected value,\n\nExpected:\n\n\t%s\n\nActual:\n\n\t%s", expected.GoString(), val.GoString())))
	}

	t.Logf("XML Importer returned value: %s", val.GoString())
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
			"testdata/importer/xml/junit0.xml",
			ctyType,
			testDataXML.ExpectedValue0(),
		)
	})

	// Harder. A single "testsuite" element but multiple "testcase" elements therein.
	t.Run("junit1", func(t *testing.T) {

		runXMLSubtestHelper(t,
			"testdata/importer/xml/junit1.xml",
			ctyType,
			testDataXML.ExpectedValue1(),
		)
	})

	// Hardest. A single "testsuite" element with a single "testcase" elements within.
	t.Run("junit2", func(t *testing.T) {

		runXMLSubtestHelper(t,
			"testdata/importer/xml/junit2.xml",
			ctyType,
			testDataXML.ExpectedValue2(),
		)
	})

}
