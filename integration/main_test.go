// +build integration

package integration

import (
	"log"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/verifa/bubbly/env"
	testData "github.com/verifa/bubbly/integration/testdata"
	"github.com/verifa/bubbly/server"
)

func TestMain(m *testing.M) {

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	bCtx.Logger.Debug().Msgf("Initializing the store")
	if err := server.InitStore(); err != nil {
		bCtx.Logger.Fatal().Err(err).Msgf("failed to create store")
	}

	bCtx.Logger.Debug().Msgf("Starting server on: %s", bCtx.ServerConfig.HostURL())

	go func() {
		err := server.ListenAndServe(bCtx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	tables, err := testData.TestAutomationSchema(".")
	if err != nil {
		log.Fatal(err)
	}

	// this should not be needed in the future... but currently we need to
	// create the schema by accessing the store directly from the bubbly server
	s := server.GetStore()

	err = s.Create(tables)
	if err != nil {
		log.Fatal(err)
	}

	// Stores that don't have type information can't save anything.
	// assert.Contains(t, server.GetStore().Save(nil).Error(), "no type information")
	os.Exit(m.Run())
}
