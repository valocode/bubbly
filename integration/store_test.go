// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/bubbly"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/store"
)

// TODO: this needs a proper test suite as well...

func TestStore(t *testing.T) {

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	err := bubbly.Apply(bCtx, "./testdata/testautomation/golang/pipeline.bubbly")
	assert.NoError(t, err)

	s, err := store.New(bCtx)
	assert.NoErrorf(t, err, "failed to create store")

	n, err := s.Query(`{
			test_case(status: "pass") {
				name
				status
			}
		}`)

	assert.NoError(t, err)

	b, err := json.MarshalIndent(n, "", "  ")
	assert.NoError(t, err)

	fmt.Println(string(b))
}
