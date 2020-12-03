package main

import (
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/handler"

	testData "github.com/verifa/bubbly/integration/testdata"
	"github.com/verifa/bubbly/server"
)

func main() {
	// initialize the router's endpoints
	router := server.SetupRouter()

	hostURL := os.Getenv("BUBBLY_ADDR")

	log.Printf("Started server on: %s", hostURL)

	serv := &http.Server{
		Addr:    hostURL,
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe(serv)

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
