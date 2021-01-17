package main

import (
	"log"
	"net/http"
)

func main() {
	var (
		addr = ":8112"
		fs   = http.FileServer(http.Dir("static"))
	)
	http.Handle("/", fs)
	log.Printf("storefront listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
