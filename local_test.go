package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valocode/bubbly/api/core"
)

func TestResource(t *testing.T) {
	// str := `{"api_version":"","kind":"pipeline","metadata":{},"name":"gotest","spec":"\n        input \"file\" {}\n        // Each task in a pipeline has an output, similar to resources,\n        // so that task outputs can be referenced\n        task \"extract\" {\n            resource = \"extract/gotest\"\n            input \"file\" {\n                value = self.input.file\n            }\n        }\n        task \"transform\" {\n            resource = \"transform/gotest\"\n            input \"data\" {\n                value = self.task.extract.value\n            }\n        }\n        task \"load\" {\n            resource = \"load/gotest\"\n            input \"data\" {\n                value = self.task.transform.value\n            }\n        }\n    "}`
	str := `{"kind":"pipeline","spec":"\n        input \"file\" {}\n        // Each task in a pipeline has an output, similar to resources,\n        // so that task outputs can be referenced\n        task \"extract\" {\n            resource = \"extract/gotest\"\n            input \"file\" {\n                value = self.input.file\n            }\n        }\n        task \"transform\" {\n            resource = \"transform/gotest\"\n            input \"data\" {\n                value = self.task.extract.value\n            }\n        }\n        task \"load\" {\n            resource = \"load/gotest\"\n            input \"data\" {\n                value = self.task.transform.value\n            }\n        }\n    "}`
	var res core.ResourceBlock
	err := json.Unmarshal([]byte(str), &res)
	assert.NoError(t, err)
}
