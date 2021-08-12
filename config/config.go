package config

import "fmt"

// ###########################################
// Store
// ###########################################

// ServerConfig is a struct storing the server information.
type ServerConfig struct {
	Port string
	Host string
}

func (s ServerConfig) HostURL() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

// ###########################################
// Store
// ###########################################

// Provider is a store provider.
type Provider string

const (
	ProviderSqlite    Provider = "sqlite"
	ProviderPostgres  Provider = "postgres"
	ProviderCockroach Provider = "cockroachdb"
)

func (_type Provider) String() string {
	return string(_type)
}

// StoreConfig stores the configuration of a bubbly store, used
// to interact with a backend database
type StoreConfig struct {
	Provider Provider

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

// ##########################
// NATS
// ##########################

type NATSServerConfig struct {
	HTTPPort int
	Port     int
}

// ##########################
// Client
// ##########################

type ClientConfig struct {
	AuthToken  string
	BubblyAddr string
	NATSAddr   string
}

// ##########################
// CLI
// ##########################

type CLIConfig struct {
	NoColor bool
	Debug   bool
}
