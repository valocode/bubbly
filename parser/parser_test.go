package parser

import (
	"testing"

	"github.com/verifa/bubbly/api/core"
)

func TestParser(t *testing.T) {
	p, err := NewParser("testdata/example")
	if err != nil {
		t.Errorf("Could not initialise parser: %s", err.Error())
		t.FailNow()
	}

	err = p.Decode()
	if err != nil {
		t.Errorf("Failed to decode parser body: %s", err.Error())
		t.FailNow()
	}

	pipelineRunKinds := p.Resources[core.PipelineRunResourceKind]
	for _, resource := range pipelineRunKinds {
		pipelineRun := resource.(core.PipelineRun)
		t.Logf("Decoding resource: %s\n", pipelineRun.String())
		if err := pipelineRun.Decode(p.Scope.decodeResource); err != nil {
			t.Errorf(`Failed to decode pipelineRun "%s": %s`, pipelineRun.String(), err.Error())
		}

		pipelineID := pipelineRun.Pipeline()
		pipeline := p.Resources.Resource(core.PipelineResourceKind, pipelineID)
		if pipeline == nil {
			t.Errorf(`Referenced pipeline "%s" does not exist`, pipelineID)
		}

		// create a nested scope with inputs
		nestedScope := p.Scope.NestedScope(pipelineRun.Inputs())
		// TODO nestedScope inputs...
		if err := pipeline.Decode(nestedScope.decodeResource); err != nil {
			t.Errorf(`Failed to decode pipeline "%s": %s`, pipelineRun.String(), err.Error())
		}
	}
}
