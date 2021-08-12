package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/migrate"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/gql"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"

	// required by schema hooks.
	_ "github.com/valocode/bubbly/ent/runtime"
)

func New(bCtx *env.BubblyContext) (*Store, error) {

	var (
		client *ent.Client
		err    error
	)
	// Connect to the provider's database RetryAttempts times, with a RetrySleep
	for connRetry := 1; connRetry <= 5; connRetry++ {
		switch bCtx.StoreConfig.Provider {
		case config.ProviderPostgres:
			var db *sql.DB
			// Create ent.Client and run the schema migration.
			db, err = sql.Open("pgx", "postgresql://postgres:postgres@127.0.0.1/bubbly")
			if err == nil {
				// Create an ent.Driver from `db`.
				drv := entsql.OpenDB(dialect.Postgres, db)
				client = ent.NewClient(ent.Driver(drv))
			}
		case config.ProviderSqlite:
			client, err = ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
		}
		// If there was no error the connection was successful
		if err == nil {
			break
		}
		fmt.Printf("Connection attempt %d to DB failed: %s\n", connRetry, err.Error())
		// Sleep for the specified amount of time
		time.Sleep(time.Second * time.Duration(time.Second*1))
	}
	if err != nil {
		return nil, fmt.Errorf("error opening db connection: %w", err)
	}

	ctx := context.Background()
	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx,
		// TODO: https://entgo.io/docs/migrate/#universal-ids
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return &Store{
		client: client,
		ctx:    context.Background(),
	}, nil
}

type (
	Store struct {
		client *ent.Client
		ctx    context.Context
	}
)

func (s *Store) Client() *ent.Client {
	return s.client
}

func (s *Store) Close() error {
	return s.client.Close()
}

func (s *Store) EvaluateRelease(id int) error {
	rel, err := s.client.Release.Get(s.ctx, id)
	if err != nil {
		return fmt.Errorf("error getting release with id %d: %w", id, err)
	}

	issues, err := rel.
		QueryCodeScans().
		QueryIssues().
		Where(codeissue.SeverityEQ(codeissue.SeverityHigh)).
		All(s.ctx)
	if err != nil {
		return fmt.Errorf("error getting issues for release: %w", err)
	}
	if len(issues) > 0 {
		_, err := s.client.Release.UpdateOneID(id).
			SetStatus(release.StatusBlocked).
			Save(s.ctx)
		if err != nil {
			return fmt.Errorf("error updating release: %w", err)
		}
	}

	// TODO: check test runs, CVEs, etc, etc
	return nil
}

func (s *Store) Query(query string) (json.RawMessage, error) {
	ctx := graphql.StartOperationTrace(context.Background())
	now := graphql.Now()
	exec := executor.New(gql.NewSchema(s.client))
	rc, errs := exec.CreateOperationContext(ctx, &graphql.RawParams{
		Query: query,
		ReadTime: graphql.TraceTiming{
			Start: now,
			End:   now,
		},
	})
	if errs != nil {
		return nil, errs
	}
	handler, rctx := exec.DispatchOperation(ctx, rc)
	response := handler(rctx)
	if response.Errors != nil {
		return nil, response.Errors
	}
	return response.Data, nil

}

func (s *Store) WithTx(fn func(tx *ent.Tx) error) error {
	tx, err := s.client.Tx(s.ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("error rolling back transaction: %w: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

// func (s *Store) SaveAdapterResult(result *adapter.Result) error {
// 	if result.ReleaseID == nil {
// 		return errors.New("must provide release id")
// 	}
// 	release, err := s.client.Release.Get(s.ctx, *result.ReleaseID)
// 	if err != nil {
// 		return fmt.Errorf("error getting release: %w", err)
// 	}

// 	if result.CodeScan != nil {
// 		rScan := result.CodeScan
// 		scan, err := s.client.CodeScan.Create().
// 			SetTool(*rScan.Tool).
// 			SetRelease(release).
// 			Save(s.ctx)
// 		if err != nil {
// 			return err
// 		}
// 		for _, rIssue := range rScan.Issues {
// 			s.client.CodeIssue.Create().
// 				SetScan(scan).
// 				SetRuleID(*rIssue.RuleID)
// 		}
// 	}

// }

func (s *Store) clientOrTx(tx *ent.Tx) *ent.Client {
	var client = s.client
	if tx != nil {
		client = tx.Client()
	}
	return client
}
