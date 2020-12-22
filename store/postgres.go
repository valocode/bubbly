package store

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/graphql-go/graphql"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

func newPostgres(bCtx *env.BubblyContext) (provider, error) {
	db := pg.Connect(&pg.Options{
		Addr:     bCtx.StoreConfig.PostgresAddr,
		User:     bCtx.StoreConfig.PostgresUser,
		Password: bCtx.StoreConfig.PostgresPassword,
		Database: bCtx.StoreConfig.PostgresDatabase,
	})

	// Attempt to create a table to hold our typeInfo. This table
	// probably already exists unless this is the first time
	// that the server has been booted against this data store.
	err := db.Model((*typeInfo)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create type info table in postgres: %w", err)
	}

	types, err := currentSchemaTypesPostgres(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get current schema types: %w", err)
	}

	return &postgres{
		db:    db,
		types: types,
	}, nil
}

type postgres struct {
	db    *pg.DB
	types map[string]schemaType
}

func (p *postgres) Create(tables core.Tables) error {
	var types map[string]schemaType
	err := p.db.RunInTransaction(func(tx *pg.Tx) error {
		info := &typeInfo{
			Tables: tables,
		}
		if _, err := tx.Model(info).Insert(); err != nil {
			return fmt.Errorf("failed to insert type info: %w", err)
		}

		types = newSchemaTypes(tables)
		for n, t := range types {
			// Create a query set to the type of our dynamic
			// struct and give it the name we were given. Normally,
			// pg would derive the name from the struct.
			q := p.db.Model(t.Empty()).Table(n)
			if err := q.CreateTable(nil); err != nil {
				return fmt.Errorf("failed to create postgres table: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create tables in postgres: %w", err)
	}

	// If we successfully committed the transaction we
	// are safe to move forward with this info.
	p.types = types

	return nil
}

func (p *postgres) Save(data core.DataBlocks) (core.Tables, error) {
	if p.types == nil {
		return nil, errors.New("postgres has no type information")
	}

	err := p.db.RunInTransaction(func(tx *pg.Tx) error {
		return p.save(tx, data, "", 0)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save data in postgres: %w", err)
	}

	return currentCoreTablesPostgres(p.db)
}

func (p *postgres) save(tx *pg.Tx, data core.DataBlocks, parentName string, parentID int64) error {
	for _, d := range data {
		// Retrieve the schema type for the data
		// we are trying to insert.
		st, ok := p.types[d.TableName]
		if !ok {
			return fmt.Errorf("unkown type: %s", d.TableName)
		}

		// Create a new instance of a predefined struct
		// that corresponds to both the data and the schema.
		n, err := st.New(d, parentName, parentID)
		if err != nil {
			return fmt.Errorf("falied to create instance of %s: %w", d.TableName, err)
		}

		// Insert the data
		if _, err := tx.Model(n).Table(d.TableName).Insert(); err != nil {
			return fmt.Errorf("falied to insert %s: %w", d.TableName, err)
		}

		// Recursively insert all sub-data.
		if err := p.save(tx, d.Data, d.TableName, schemaTypeID(n)); err != nil {
			return err
		}
	}

	return nil
}

func (p *postgres) ResolveScalar(params graphql.ResolveParams) (interface{}, error) {
	var (
		tableName = params.Info.FieldName
		n         = p.types[tableName].Empty()
		q         = p.db.Model(n).Table(tableName)
	)

	q = applyArgsPostgres(q, params.Args)

	if err := q.Last(); err != nil {
		return nil, fmt.Errorf("failed to resolve scalar %s: %w", tableName, err)
	}

	return n, nil
}

func (p *postgres) ResolveList(params graphql.ResolveParams) (interface{}, error) {
	var (
		tableName = params.Info.FieldName
		parent    = params.Source
		n         = p.types[tableName].EmptySlice()
		q         = p.db.Model(n).Table(tableName)
	)

	q = applyArgsPostgres(q, params.Args)

	if isValidParent(parent) {
		var (
			parentIDField = schemaTypeName(parent) + "_id"
			parentID      = strconv.FormatInt(schemaTypeID(parent), 10)
		)
		q = q.Where(parentIDField+" = ?", parentID)
	}

	if err := q.Select(); err != nil {
		return nil, fmt.Errorf("failed to resolve list %s: %w", tableName, err)
	}

	return n, nil
}

func (p *postgres) LastValue(tableName, field string) (cty.Value, error) {
	n, err := p.ResolveScalar(graphql.ResolveParams{
		Info: graphql.ResolveInfo{
			FieldName: tableName,
		},
	})
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to resolve ref as scalar: %w", err)
	}

	val, err := schemaTypeVal(n, field)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to coerce value for %s: %w", field, err)
	}

	return val, nil
}

func applyArgsPostgres(q *orm.Query, args map[string]interface{}) *orm.Query {
	for k, v := range args {
		if k != filterName {
			q = q.Where(k+" = ?", v)
			continue
		}

		hasSuffix := func(n string, filter string) (bool, string) {
			if !strings.HasSuffix(n, filter) {
				return false, ""
			}
			return true, strings.ReplaceAll(n, filter, "")
		}

		filter := v.(map[string]interface{})
		for n, val := range filter {
			if ok, f := hasSuffix(n, filterGreaterThan); ok {
				q = q.Where(f+" > ?", val)
				continue
			}
			if ok, f := hasSuffix(n, filterLessThan); ok {
				q = q.Where(f+" < ?", val)
				continue
			}
			if ok, f := hasSuffix(n, filterGreaterThanOrEqualTo); ok {
				q = q.Where(f+" >= ?", val)
				continue
			}
			if ok, f := hasSuffix(n, filterLessThanOrEqualTo); ok {
				q = q.Where(f+" <= ?", val)
				continue
			}
			// Important: test "_not_in" first so we don't
			// accidentally match "_in".
			if ok, f := hasSuffix(n, filterNotIn); ok {
				q = q.Where(f+" NOT IN (?)", pg.In(val))
				continue
			}
			if ok, f := hasSuffix(n, filterIn); ok {
				q = q.Where(f+" IN (?)", pg.In(val))
			}
		}
	}

	return q
}

func currentCoreTablesPostgres(db *pg.DB) (core.Tables, error) {
	// Try to load the most recent typeInfo.
	// TODO(andrewhare): We need a plan for how to handle
	// migrations from one typeInfo to another.
	var info typeInfo
	err := db.Model(&info).Last()
	if err != nil {
		return nil, fmt.Errorf("failed to get type info from postgres: %w", err)
	}

	return info.Tables, nil
}

func currentSchemaTypesPostgres(db *pg.DB) (map[string]schemaType, error) {
	tables, err := currentCoreTablesPostgres(db)
	if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return newSchemaTypes(tables), nil
}
