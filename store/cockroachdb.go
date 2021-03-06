package store

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/graphql-go/graphql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

// FIXME: because Roach provider was heavy dependent on Postgres,
//        it now also uses pgpool connection pool. Is that ok?

func newCockroachdb(bCtx *env.BubblyContext) (*cockroachdb, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		bCtx.StoreConfig.CockroachUser,
		bCtx.StoreConfig.CockroachPassword,
		bCtx.StoreConfig.CockroachAddr,
		bCtx.StoreConfig.CockroachDatabase,
	)
	pool, err := psqlNewPool(bCtx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize connection to db: %w", err)
	}

	return &cockroachdb{
		pool: pool,
	}, nil
}

type cockroachdb struct {
	pool *pgxpool.Pool
}

func (c *cockroachdb) Apply(schema *bubblySchema) error {

	err := crdbpgx.ExecuteTx(context.Background(), c.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return psqlApplySchema(tx, schema)
	})
	if err != nil {
		return fmt.Errorf("failed to apply tables: %w", err)
	}

	return nil
}

func (c *cockroachdb) Save(bCtx *env.BubblyContext, schema *bubblySchema, tree dataTree) error {

	err := crdbpgx.ExecuteTx(context.Background(), c.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		saveNode := func(bCtx *env.BubblyContext, node *dataNode, blocks *core.DataBlocks) error {
			return psqlSaveNode(tx, node, schema)
		}

		_, err := tree.traverse(bCtx, saveNode)

		return err
	})
	if err != nil {
		return fmt.Errorf("failed to save data in cockroachdb: %w", err)
	}

	return nil
}

func (c *cockroachdb) ResolveQuery(graph *schemaGraph, params graphql.ResolveParams) (interface{}, error) {
	return psqlResolveRootQueries(c.pool, graph, params)
}

func (c *cockroachdb) HasTable(table core.Table) (bool, error) {
	return psqlHasTable(c.pool, table)
}

func (c *cockroachdb) GenerateMigration(bCtx *env.BubblyContext, cl Changelog) (migration, error) {
	return psqlGenerateMigration(bCtx, cl)
}

func (c *cockroachdb) Migrate(migrationList migration) error {
	return psqlMigrate(c.pool, migrationList)
}
