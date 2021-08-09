package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/gql"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/test"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"

	// required by schema hooks.
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/release"
	_ "github.com/valocode/bubbly/ent/runtime"
)

func main() {
	var (
		sqliteDB bool
		skipAll  bool
		skipCVE  bool
		skipSPDX bool
		client   *ent.Client
		provider = store.ProviderPostgres
	)
	flag.BoolVar(&sqliteDB, "local", false, "Whether to use a local sqlite DB with no external dependencies. Default false, which uses an external postgres DB")
	flag.BoolVar(&skipCVE, "skip-cve", false, "Skips cve data creation on startup")
	flag.BoolVar(&skipSPDX, "skip-spdx", false, "Skips spdx data creation on startup")
	flag.BoolVar(&skipAll, "skip-all", false, "Skips all data creation on startup")
	flag.Parse()

	if sqliteDB {
		provider = store.ProviderSqlite
	}

	store, err := store.New(provider)
	if err != nil {
		log.Fatal("error connecting to db: ", err)
	}

	client = store.Client()

	if !skipAll {
		if !skipCVE {
			fmt.Println("Fetching CVEs from NVD... This will take a few seconds...")
			if err := test.SaveCVEData(client); err != nil {
				log.Fatal("loading CVEs: ", err)
			}
			fmt.Println("Done!")
			fmt.Println("")
		}
		if !skipSPDX {
			fmt.Println("Fetching SPDX licenses from GitHub...")
			if err := test.SaveSPDXData(client); err != nil {
				log.Fatal("loading SPDX: ", err)
			}
			fmt.Println("Done!")
			fmt.Println("")

		}

		fmt.Println("Creating dummy releases...")
		if err := test.CreateDummyData(client); err != nil {
			log.Fatal("creating dummy data: ", err)
		}
		fmt.Println("Done!")
		fmt.Println("")

		fmt.Println("Evaluating releases...")
		if err := test.FailSomeRandomReleases(store); err != nil {
			log.Fatal("evaluating releases: ", err)
		}
		fmt.Println("Done!")
		fmt.Println("")
	}

	ctx := context.Background()
	releases := client.Release.Query().Where(release.HasCommitWith(
		gitcommit.TimeGTE(time.Now().AddDate(0, 0, -30)),
	)).AllX(ctx)

	fmt.Printf("Releases: %d\n", len(releases))

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	srv := handler.NewDefaultServer(gql.NewSchema(client))
	router.Handle("/", playground.Handler("Bubbly", "/query"))
	router.Handle("/query", srv)

	log.Println("listening on :8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatal("http server terminated", err)
	}
}
