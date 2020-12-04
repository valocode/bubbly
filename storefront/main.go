package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/rs/zerolog"

	"github.com/verifa/bubbly/env"
	testData "github.com/verifa/bubbly/integration/testdata"
	"github.com/verifa/bubbly/server"
)

func main() {
	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	go func() {
		err := server.ListenAndServe(bCtx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	tables, err := testData.TestAutomationSchema("../integration")
	if err != nil {
		log.Fatal(err)
	}

	data, err := testData.TestAutomationData("../integration")
	if err != nil {
		log.Fatal(err)
	}

	s := server.GetStore()

	if err := s.Create(tables); err != nil {
		log.Fatal(err)
	}
	if err := s.Save(data); err != nil {
		log.Fatal(err)
	}

	schema := s.Schema()
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/graphql", h)
	http.Handle("/", fs)
	log.Println("storefront listening on :8111")
	log.Fatal(http.ListenAndServe(":8111", nil))
}
