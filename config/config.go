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

// StoreConfig stores the configuration of a bubbly store, used
// to interact with a backend database
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

// ###########################################
// Agent
// ###########################################

type AgentDeploymentType string

const (
	SingleDeployment AgentDeploymentType = "single"
	// TODO: Implement
	// DistributedDeployment AgentDeploymentType = "distributed"
)

// AgentConfig stores the configuration of a bubbly agent
type AgentConfig struct {
	StoreConfig       *StoreConfig
	NATSServerConfig  *NATSServerConfig
	EnabledComponents *AgentComponentsToggle
	DeploymentType    AgentDeploymentType
}

type AgentComponentsToggle struct {
	UI        bool
	APIServer bool
	DataStore bool
	Worker    bool
}

// ##########################
// NATS
// ##########################

type NATSServerConfig struct {
	HTTPPort int
	Port     int
	Addr     string
}
