// +build integration

package integration

import (
	"log"
	"net/http"
	"os"
	"testing"

	testData "github.com/verifa/bubbly/integration/testdata"
	"github.com/verifa/bubbly/server"
)

var hostURL string

func TestMain(m *testing.M) {

	hostURL = "localhost:8112"
	// initialize the router's endpoints
	router := server.SetupRouter()

	serv := &http.Server{
		Addr:    hostURL,
		Handler: router,
	}

	log.Printf("Starting server on: %s", hostURL)

	go func() {
		err := server.ListenAndServe(serv)
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
