package store

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"

	bolt "go.etcd.io/bbolt"
)

var dbPath = "default.db"

const (
	schemaBucket      = "schema"
	schemaGraphBucket = "schemaGraph"
	triggerBucket     = "trigger"
)

func saveElement(tenantID string, bucketID string, element []byte) error {
	db, err := bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return fmt.Errorf("failed to open bolt DB: %w", err)
	}
	defer db.Close()

	updateErr := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketID))
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(tenantID), element)
		if err != nil {
			return fmt.Errorf("failed to store element in Bolt: %w", err)
		}
		return nil
	})
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func fetchElement(tenantID string, bucketID string) ([]byte, error) {
	db, err := bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt DB: %w", err)
	}
	var b []byte
	defer db.Close()
	updateErr := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		b = bucket.Get([]byte(tenantID))
		return nil
	})
	if updateErr != nil {
		return nil, fmt.Errorf("error fetching from bucket %s, %w", bucketID, updateErr)
	}
	return b, nil
}

func saveSchema(tenantID string, schema *graphql.Schema) error {
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	return saveElement(tenantID, schemaBucket, schemaBytes)
}

func saveTrigger(tenantID string, triggers []*trigger) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	err := enc.Encode(triggers)
	if err != nil {
		return fmt.Errorf("failed to encode trigger: %w", err)
	}

	return saveElement(tenantID, triggerBucket, buf.Bytes())
}

func saveSchemaGraph(tenantID string, graph *schemaGraph) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	err := enc.Encode(graph)
	if err != nil {
		return fmt.Errorf("failed to encode SchemaGraph: %w", err)
	}

	return saveElement(tenantID, schemaGraphBucket, buf.Bytes())
}

func fetchSchema(tenantID string) (*graphql.Schema, error) {
	results, err := fetchElement(tenantID, schemaBucket)
	if err != nil {
		return nil, err
	}
	s := graphql.Schema{}
	err = json.Unmarshal(results, &s)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall schema: %w", err)
	}
	return &s, nil
}

func fetchSchemaGraph(tenantID string) (*schemaGraph, error) {
	results, err := fetchElement(tenantID, schemaGraphBucket)
	if err != nil {
		return nil, err
	}
	sg := schemaGraph{}
	buf := bytes.NewBuffer(results)
	err = json.NewDecoder(buf).Decode(&sg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode SchemaGraph: %w", err)
	}
	return &sg, nil
}

func fetchTriggers(tenantID string) ([]*trigger, error) {
	results, err := fetchElement(tenantID, triggerBucket)
	if err != nil {
		return nil, err
	}
	var trs []*trigger

	err = json.Unmarshal(results, &trs)
	if err != nil {
		return nil, fmt.Errorf("failed to  unmarshal triggers: %w", err)
	}
	return trs, nil
}
