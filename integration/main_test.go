// +build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/env"

	testData "github.com/verifa/bubbly/integration/testdata"
)

func TestMain(m *testing.M) {

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	client, err := client.New(bCtx)
	if err != nil {
		bCtx.Logger.Fatal().Err(err).Msg("failed to create client")
	}

	tables, err := testData.TestSchema(filepath.FromSlash("./testdata/schema/schema.bubbly"))
	if err != nil {
		bCtx.Logger.Fatal().Err(err).Msg("failed to parse schema")
	}

	tableBytes, err := json.Marshal(tables)
	if err != nil {
		bCtx.Logger.Fatal().Err(err).Msg("failed to json marshal schema")
	}

	if err := client.PostSchema(bCtx, tableBytes); err != nil {
		bCtx.Logger.Fatal().Err(err).Msg("failed to post schema to bubbly server")
	}

	// Run the tests in this module
	os.Exit(m.Run())
}
