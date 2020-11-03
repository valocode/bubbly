package server

import (
	"fmt"
	"os"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/interim"
	"github.com/zclconf/go-cty/cty"
)

var router *gin.Engine
var db *interim.DB

// SetupRouter returns a pointer to a gin engine after setting up middleware
// and initializing routes
func SetupRouter() *gin.Engine {
	// Enable Production Mode
	// gin.SetMode(gin.ReleaseMode)

	// Initialize Logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	// SETUP DB
	dbErr := initDb()
	if dbErr != nil {
		log.Error().Msg("Error setting up DB: " + dbErr.Error())
	}
	// Initialize Router
	// router := gin.Default()  // Sets the Gin defaults
	router := gin.New() // Use a blank Gin server with no middleware loaded
	router.Use(logger.SetLogger())
	router.Use(gin.Recovery())
	router.Use(VersionMiddleware())

	// Initialize HTTP Routes
	InitializeRoutes(router)

	return router
}

func initDb() error {
	var err error
	db, err = interim.NewDB(getSchema())
	if err != nil {
		return fmt.Errorf("failed to create memDB: %w", err)
	}
	return nil
}

// GetDb returns a pointer to the DB
func GetDb() *interim.DB {
	return db
}

func getSchema() []core.Table {
	tables := []core.Table{
		{
			Name: "product",
			Fields: []core.TableField{
				{
					Name:   "Name",
					Type:   cty.String,
					Unique: true,
				},
			},
		},
		{
			Name: "project",
			Fields: []core.TableField{
				{
					Name:   "Name",
					Type:   cty.String,
					Unique: true,
				},
			},
		},
		{
			Name: "repository",
			Fields: []core.TableField{
				{
					Name:   "Id",
					Type:   cty.String,
					Unique: true,
				},
				{
					Name:   "Url",
					Type:   cty.String,
					Unique: false,
				},
			},
		},
		{
			Name: "repository_version",
			Fields: []core.TableField{
				{
					Name:   "Id",
					Type:   cty.String,
					Unique: true,
				},
				{
					Name:   "Tag",
					Type:   cty.String,
					Unique: false,
				},
				{
					Name:   "Branch",
					Type:   cty.String,
					Unique: false,
				},
				{
					Name:   "RepositoryId",
					Type:   cty.String,
					Unique: false,
				},
			},
		},
	}
	return tables
}
