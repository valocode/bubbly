package test

import (
	"context"
	"log"
	"testing"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/ent"

	// required by schema hooks.
	_ "github.com/valocode/bubbly/ent/runtime"
)

func TestDummyData(t *testing.T) {
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

	{
		err := CreateDummyData(client)
		require.NoError(t, err)
	}

	// artifact := client.Artifact.Query().WithEntry().FirstX(ctx)
	// client.ReleaseEntry.UpdateOneID(artifact.Edges.Entry.ID).SetTime(time.Now().AddDate(-5, 0, 0)).S
	// cases := client.TestCase.Query().WithRun().AllX(ctx)
	// for _, c := range cases {
	// 	t.Logf("case: %s -- %s", c.String(), c.Edges.Run.String())
	// }
}
