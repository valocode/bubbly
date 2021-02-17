package config

import (
	"fmt"
)

// ServerConfig is a struct storing the server information.
type ServerConfig struct {
	Protocol string
	Port     string
	Host     string
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
	// CockroachDBStore is a CockroachDB provider.
	CockroachDBStore = "cockroachdb"
)

// StoreConfig stores the configuration of a bubbly store, used
// to interact with a backend database
type StoreConfig struct {
	Provider StoreProviderType

	PostgresAddr     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string

	CockroachAddr     string
	CockroachUser     string
	CockroachPassword string
	CockroachDatabase string

	RetrySleep    int
	RetryAttempts int
}

// ###########################################
// Agent
// ###########################################

type AgentDeploymentType string

const (
	SingleDeployment AgentDeploymentType = "single"
	// TODO: Implement
	// DistributedDeployment AgentDeploymentType = "distributed"
)

func (a AgentDeploymentType) String() string {
	switch a {
	case SingleDeployment:
		return "single"
	default:
		return "unsupported"
	}
}

// AgentConfig stores the configuration of a bubbly agent
type AgentConfig struct {
	StoreConfig       *StoreConfig
	NATSServerConfig  *NATSServerConfig
	EnabledComponents *AgentComponentsToggle
	DeploymentType    AgentDeploymentType
}

type AgentComponentsToggle struct {
	APIServer  bool
	DataStore  bool
	Worker     bool
	NATSServer bool
}

// ##########################
// NATS
// ##########################

type NATSServerConfig struct {
	HTTPPort int
	Port     int
	Addr     string
}
