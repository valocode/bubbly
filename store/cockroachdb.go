package store

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/graphql-go/graphql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

var _ provider = (*cockroachdb)(nil)

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

func (c *cockroachdb) Apply(tenant string, schema *bubblySchema) error {

	err := crdbpgx.ExecuteTx(context.Background(), c.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return psqlApplySchema(tx, tenant, schema)
	})
	if err != nil {
		return fmt.Errorf("failed to apply tables: %w", err)
	}

	return nil
}

func (c *cockroachdb) Migrate(tenant string, schema *bubblySchema, cl schemaUpdates) error {
	pgSchema := psqlBubblySchemaPrefix + tenant
	migration, err := psqlGenerateMigration(config.CockroachDBStore, pgSchema, cl)
	if err != nil {
		return fmt.Errorf("failed to generate migration list: %w", err)
	}
	return psqlMigrate(c.pool, tenant, schema, migration)
}

func (c *cockroachdb) Save(bCtx *env.BubblyContext, tenant string, graph *schemaGraph, tree dataTree) error {

	err := crdbpgx.ExecuteTx(context.Background(), c.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		saveNode := func(bCtx *env.BubblyContext, node *dataNode, blocks *core.DataBlocks) error {
			// Check that the data node we are saving exists in the schema graph.
			// Otherwise it does not exist in our schema
			tNode, ok := graph.NodeIndex[node.Data.TableName]
			if !ok {
				return fmt.Errorf("data block refers to non-existing table: %s", node.Data.TableName)
			}
			return psqlSaveNode(tx, tenant, node, *tNode.table)
		}

		_, err := tree.traverse(bCtx, saveNode)

		return err
	})
	if err != nil {
		return fmt.Errorf("failed to save data in cockroachdb: %w", err)
	}

	return nil
}

func (c *cockroachdb) ResolveQuery(tenant string, graph *schemaGraph, params graphql.ResolveParams) (interface{}, error) {
	return psqlResolveRootQueries(c.pool, tenant, graph, params)
}

func (c *cockroachdb) Tenants() ([]string, error) {
	return psqlTenantSchemas(c.pool)
}

func (c *cockroachdb) CreateTenant(name string) error {
	return psqlCreateSchema(c.pool, name)
}

func (c *cockroachdb) HasTable(tenant string, table core.Table) (bool, error) {
	return psqlHasTable(c.pool, tenant, table)
}
