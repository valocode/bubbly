package store

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/graphql-go/graphql"
	"github.com/jackc/pgx/v4"
	"github.com/verifa/bubbly/env"
)

func newCockroachdb(bCtx *env.BubblyContext) (*cockroachdb, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		bCtx.StoreConfig.CockroachUser,
		bCtx.StoreConfig.CockroachPassword,
		bCtx.StoreConfig.CockroachAddr,
		bCtx.StoreConfig.CockroachDatabase,
	)
	conn, err := psqlNewConn(bCtx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize connection to db: %w", err)
	}

	return &cockroachdb{
		conn: conn,
	}, nil
}

type cockroachdb struct {
	conn *pgx.Conn
}

func (c *cockroachdb) Apply(schema *bubblySchema) error {

	err := crdbpgx.ExecuteTx(context.Background(), c.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return psqlApplySchema(tx, schema)
	})
	if err != nil {
		return fmt.Errorf("failed to apply tables: %w", err)
	}

	// If the apply was successful then store the current state of the tables.
	// TODO: how does this work in distributed mode? If another store were to
	// update the schema, this would need to be sync'd
	return nil
}

func (c *cockroachdb) Save(sc *saveContext) error {

	err := crdbpgx.ExecuteTx(context.Background(), c.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return psqlSaveData(tx, sc)
	})
	if err != nil {
		return fmt.Errorf("failed to save data in cockroachdb: %w", err)
	}
	return nil
}

func (c *cockroachdb) ResolveScalar(params graphql.ResolveParams) (interface{}, error) {
	return psqlResolveScalar(c.conn, params)
}

func (c *cockroachdb) ResolveList(params graphql.ResolveParams) (interface{}, error) {
	return psqlResolveParams(c.conn, params)
}
