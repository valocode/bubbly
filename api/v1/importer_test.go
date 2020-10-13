package v1

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func TestImporterJSON(t *testing.T) {
	formatString := `
	object({
		issues: list(object({
			engineId: string,
			ruleId: string,
			severity: string,
			type: string,
			primaryLocation: object({
				message: string,
				filePath: string,
				textRange: object({
					startLine: number,
					endLine: number,
					startColumn: number,
					endColumn: number
				})
			})
		}))
	})
	`
	expr, diags := hclsyntax.ParseExpression([]byte(formatString), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		t.Errorf("Failed to create expression: %s", diags.Error())
		t.FailNow()
	}
	importer := Importer{
		Spec: ImporterSpec{
			Type:   jsonImporterType,
			Format: expr,
		},
		Source: &JSONSource{
			File: "testdata/importer/json/sonarqube-example.json",
		},
	}

	val, err := importer.Resolve()
	if err != nil {
		t.Errorf("Failed to Resolve() JSON importer: %s", err.Error())
		t.Fail()
	}
	if val.IsNull() {
		t.Errorf("Received Null type value")
	}
	t.Logf("JSON Importer returned value: %s", val.GoString())
}

func TestImporterXML(t *testing.T) {
	formatString := `
	object({
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
	`
	expr, diags := hclsyntax.ParseExpression([]byte(formatString), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		t.Errorf("Failed to create expression: %s", diags.Error())
		t.FailNow()
	}
	importer := Importer{
		Spec: ImporterSpec{
			Type:   xmlImporterType,
			Format: expr,
		},
		Source: &XMLSource{
			File: "testdata/importer/xml/junit.xml",
		},
	}

	val, err := importer.Resolve()
	if err != nil {
		t.Errorf("Failed to Resolve() XML importer: %s", err.Error())
		t.Fail()
	}
	if val.IsNull() {
		t.Errorf("Received Null type value")
	}
	t.Logf("XML Importer returned value: %s", val.GoString())

}
