package resource

import (
	"os"
	"path"

	"github.com/tidwall/buntdb"
)

// TODO change these to cfg
var dbPath = MustBuntDBPath()
var dbIndex = "index"
var dbIndexName = "name"

// MustBuntDBPath returns the path to the DB
func MustBuntDBPath() string {
	url := "bunt.db"
	e, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dbpath := path.Join(path.Dir(e), url)
	return dbpath
}

func newBuntdb(cfg Config) (provider, error) {
	db, err := buntdb.Open(dbPath)
	if err != nil {
		return nil, err
	}
	db.CreateIndex(dbIndex, "*", buntdb.IndexJSON(dbIndexName))
	return &buntDBProvider{
		db: db,
	}, nil
}

type buntDBProvider struct {
	db *buntdb.DB
}

func (p *buntDBProvider) Query(key string) (string, error) {
	var value string
	err := p.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		value = val
		return nil
	})
	if err != nil {
		return "", err
	}
	return value, nil
}

func (p *buntDBProvider) Save(key string, value string) error {
	err := p.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
