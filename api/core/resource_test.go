package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceJSON(t *testing.T) {
	// bCtx := env.NewBubblyContext()
	res := ResourceBlock{
		ResourceKind:       string(ExtractResourceKind),
		ResourceName:       "test-extract",
		ResourceAPIVersion: "v1",
		// SpecRaw:            []byte(""),
		SpecRaw: "",
	}

	b, err := json.Marshal(res)
	assert.NoError(t, err)
	t.Logf("bytes: %s", string(b))
}
