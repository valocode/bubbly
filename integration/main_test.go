// +build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"

	schemaApplyCmd "github.com/valocode/bubbly/cmd/schema/apply"
	"github.com/valocode/bubbly/env"
)

func TestMain(m *testing.M) {

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	schemaApplyCmd, _ := schemaApplyCmd.NewCmdApply(bCtx)

	schemaApplyCmd.SetArgs([]string{"-f", filepath.FromSlash("./testdata/schema/schema.bubbly")})
	schemaApplyCmd.SilenceUsage = true

	err := schemaApplyCmd.Execute()
	if err != nil {
		bCtx.Logger.Fatal().Err(err).Msg("failed to apply schema")
	}

	// Run the tests in this module
	os.Exit(m.Run())
}
