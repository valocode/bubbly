package config

import (
	"fmt"
)

// ServerConfig is a struct storing the server information.
type ServerConfig struct {
	Protocol string
	Port     string
	Host     string `validate:"required"`
	Auth     bool
	Token    string
}

func (s ServerConfig) HostURL() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

// ###########################################
// Store
// ###########################################

// StoreProviderType is a store provider.
type StoreProviderType string

const (
	// PostgresStore is a Postgres provider.
	PostgresStore StoreProviderType = "postgres"
)

// Config is a top-level struct for holding
// all provider-specific config.
type StoreConfig struct {
	Provider StoreProviderType

	PostgresAddr     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
}

// ###########################################
// Resources
// ###########################################
type ResourceConfig struct {
	Provider ResourceProviderType
}

// ResourceProviderType is a resource manager provider
type ResourceProviderType string

const (
	// Etcd is an Etcd provider
	EtcdResourceProvider ResourceProviderType = "etcd"
	// Buntdb is a Buntdb provider
	BuntdbResourceProvider ResourceProviderType = "buntdb"
)
