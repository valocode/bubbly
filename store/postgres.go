package store

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/jackc/pgx/v4"
	"github.com/verifa/bubbly/env"
)

func newPostgres(bCtx *env.BubblyContext) (*postgres, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		bCtx.StoreConfig.PostgresUser,
		bCtx.StoreConfig.PostgresPassword,
		bCtx.StoreConfig.PostgresAddr,
		bCtx.StoreConfig.PostgresDatabase,
	)
	conn, err := psqlNewConn(bCtx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize connection to db: %w", err)
	}

	return &postgres{
		conn: conn,
	}, nil
}

type postgres struct {
	conn *pgx.Conn
}

func (p *postgres) Apply(schema *bubblySchema) error {
	tx, err := p.conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	err = psqlApplySchema(tx, schema)
	if err != nil {
		return fmt.Errorf("failed to apply tables: %w", err)
	}

	return tx.Commit(context.Background())
}

func (p *postgres) Save(schema *bubblySchema, tree dataTree) error {
	tx, err := p.conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	saveNode := func(node *dataNode) error {
		return psqlSaveNode(tx, node, schema)
	}
	if err := tree.traverse(saveNode); err != nil {
		return fmt.Errorf("failed to save data in postgres: %w", err)
	}

	return tx.Commit(context.Background())
}

func (p *postgres) ResolveScalar(params graphql.ResolveParams) (interface{}, error) {
	return psqlResolveScalar(p.conn, params)
}

func (p *postgres) ResolveList(params graphql.ResolveParams) (interface{}, error) {
	return psqlResolveParams(p.conn, params)
}
