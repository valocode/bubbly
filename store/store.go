package store

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-playground/validator/v10"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/migrate"
	"github.com/valocode/bubbly/env"

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

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// If the struct field is another struct embedded, then check if there
		// is an alias tag attached so that we can get a nicer name than the
		// horrible struct names, e.g. MyTypeModelCreate `alias:"type"`
		if name == "" && fld.Anonymous {
			return strings.SplitN(fld.Tag.Get("alias"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	s := &Store{
		client:    client,
		ctx:       context.Background(),
		validator: validate,
	}

	if err := s.initDB(); err != nil {
		return nil, err
	}

	return s, nil
}

type (
	Store struct {
		client    *ent.Client
		ctx       context.Context
		validator *validator.Validate
	}
)

func (s *Store) Client() *ent.Client {
	return s.client
}

func (s *Store) Close() error {
	return s.client.Close()
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

func (s *Store) initDB() error {
	//
	// Make sure default organisation exists: TODO
	//

	//
	// Make sure default project exists
	//
	_, err := s.client.Project.Create().
		SetName(config.DefaultReleaseProject).
		Save(s.ctx)
	// Constraint error is fine (in case it already exists). Everything else is not
	if !ent.IsConstraintError(err) {
		return err
	}

	return nil
}

func (s *Store) clientOrTx(tx *ent.Tx) *ent.Client {
	var client = s.client
	if tx != nil {
		client = tx.Client()
	}
	return client
}
