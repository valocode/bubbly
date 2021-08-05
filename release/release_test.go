package release

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/graphql/handler"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/gql"
)

func graphqlRequest(query string) ([]byte, error) {
	queryData := map[string]string{
		"query": query,
	}
	jsonReq, err := json.Marshal(queryData)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("http://localhost:8081/query", "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func TestRelease(t *testing.T) {
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Configure the server and start listening on :8081.
	srv := handler.NewDefaultServer(gql.NewSchema(client))
	http.Handle("/query", srv)
	log.Println("listening on :8081")
	go func() {
		http.ListenAndServe(":8081", nil)
	}()

	{
		release, err := CreateReleaseFromSpec("../.bubbly.hcl")
		require.NoError(t, err)
		project := ent.NewProjectNode().SetName(release.Project)
		err = project.Graph().Save(client)
		require.NoError(t, err)

		t.Logf("release: %#v", release)

		relNode, err := release.Node(false)
		require.NoError(t, err)
		err = relNode.Graph().Save(client)
		require.NoError(t, err)
	}

	{
		rels := client.Release.Query().WithChecks().WithCommit().WithProject().AllX(ctx)
		for _, rel := range rels {
			t.Logf("release: %#v", rel)
		}
	}
}
