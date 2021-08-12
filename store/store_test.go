package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"
)

func TestAdapterResults(t *testing.T) {
	store, err := New(env.NewBubblyContext())
	require.NoError(t, err)

	apt, err := adapter.FromFile("../adapter/testdata/adapters/gosec.adapt.hcl")
	require.NoError(t, err)

	adapterModel, err := apt.Model()
	require.NoError(t, err)

	dbAdapter, err := store.SaveAdapter(&api.AdapterSaveRequest{
		AdapterModel: adapterModel,
	})
	if err != nil {
		t.Logf("err: %#v", err)
	}
	require.NoError(t, err)
	t.Logf("Saved adapter: %s", dbAdapter.String())
	//
	// TO RUN
	//
	// args := adapter.RunArgs{
	// 	Filename: "../adapter/testdata/adapters/gosec.json",
	// }
	// result, err := apt.Run(args)
	// require.NoError(t, err)

	// {
	// 	commit := "12345"
	// 	scanReq := api.CodeScanRequest{
	// 		CodeScan: result.CodeScan,
	// 		Commit:   &commit,
	// 	}
	// 	codeScan, err := store.SaveCodeScan(&scanReq)
	// 	require.NoError(t, err)
	// 	t.Logf("Created code scan: %s", codeScan.String())
	// }

}
