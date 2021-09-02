package config

import "fmt"

// ###########################################
// Release
// ###########################################

// ReleaseConfig is a struct storing local bubbly runtime configs, such as the
// localtion of the .bubbly directory containing the release specification and
// local adapters
type ReleaseConfig struct {
	// BubblyDir points to the .bubbly directory
	BubblyDir string
	// ReleaseSpec points to the release specification file explicitly
	// (default $BUBBLY_DIR/release.json)
	ReleaseSpec string
	// Project defines the bubbly project
	Project string
}

// ###########################################
// Server
// ###########################################

// ServerConfig is a struct storing the server information.
type ServerConfig struct {
	Port string
	Host string
	UI   bool
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

func (c *ClientConfig) V1() string {
	return c.BubblyAddr + "/api/v1"
}

// ##########################
// CLI
// ##########################

type CLIConfig struct {
	NoColor bool
	Debug   bool
}
