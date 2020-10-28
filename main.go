package main

import (
	"net/http"

	"github.com/verifa/bubbly/server"
	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var BUBBLY_API_PORT = ":8080"
var bubblyVersion = "0.0.1"

// The starting point for launching Bubbly's web API
func main() {
	server.SetVersion(bubblyVersion)
	// initialize the router's endpoints
	router = server.SetupRouter()

	// Execute Bubbly server (default port is 8080)
	// router.Run(BUBBLY_API_PORT)
	
	server := &http.Server{
		Addr: "localhost" + BUBBLY_API_PORT,
		Handler: router,
	}
	listenAndServe(server)
}

func listenAndServe(s *http.Server) error {
	var g errgroup.Group
	g.Go(func() error {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	return g.Wait()
}
