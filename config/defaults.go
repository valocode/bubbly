package config

import (
	"os"
	"strconv"
)

// Default CLI configuration
const (
	DefaultCLINoColorToggle = false
	DefaultCLIDebugToggle   = false
)

// Default Bubbly Server configuration
const (
	DefaultServerHost = ""
	DefaultServerPort = "8111"
)

// Default store configuration
const (
	DefaultStoreProvider = "sqlite"
	DefaultRetryAttempts = 5
	DefaultRetrySleep    = 1
)

// Default store configuration for Postgres
const (
	DefaultPostgresAddr     = "postgres:5432"
	DefaultPostgresUser     = "postgres"
	DefaultPostgresPassword = "postgres"
	DefaultPostgresDatabase = "bubbly"
)

// Default store configuration for CockroachDB
const (
	defaultCockroachAddr     = "cockroachdb:26257"
	defaultCockroachUser     = "root"
	defaultCockroachPassword = "admin"
	defaultCockroachDatabase = "defaultdb"
)

// Default configuration for NATS Server
const (
	DefaultNATSServerHTTPPort = "8222"
	DefaultNATSServerPort     = "4223"
)

// Default configuration for the bubbly client config
const (
	DefaultClientAuthToken = ""
	DefaultBubblyAddr      = "http://localhost:8111/api/v1"
	DefaultNATSAddr        = "localhost:4223"
)

// defaultEnvStr reads a string value from the environment and falls back to the
// default value if not provided.
func defaultEnvStr(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// defaultEnvBool reads a boolean value from the environment and falls back to the
// default value if not provided.
func defaultEnvBool(key string, defaultValue bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		val, _ := strconv.ParseBool(defaultEnvStr(key, value))
		return val
	}
	return defaultValue
}

// DefaultServerConfig creates a ServerConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: defaultEnvStr("BUBBLY_HOST", DefaultServerHost),
		Port: defaultEnvStr("BUBBLY_PORT", DefaultServerPort),
	}
}

// DefaultStoreConfig creates a StoreConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultStoreConfig() *StoreConfig {
	return &StoreConfig{
		// Default provider
		Provider: Provider(defaultEnvStr("BUBBLY_STORE", DefaultStoreProvider)),
		// Default configuration for Postgres
		PostgresAddr:     defaultEnvStr("POSTGRES_ADDR", DefaultPostgresAddr),
		PostgresUser:     defaultEnvStr("POSTGRES_USER", DefaultPostgresUser),
		PostgresPassword: defaultEnvStr("POSTGRES_PASSWORD", DefaultPostgresPassword),
		PostgresDatabase: defaultEnvStr("POSTGRES_DATABASE", DefaultPostgresDatabase),
		// Default configuration for CockroachDB
		CockroachAddr:     defaultEnvStr("COCKROACH_ADDR", defaultCockroachAddr),
		CockroachUser:     defaultEnvStr("COCKROACH_USER", defaultCockroachUser),
		CockroachPassword: defaultEnvStr("COCKROACH_PASSWORD", defaultCockroachPassword),
		CockroachDatabase: defaultEnvStr("COCKROACH_DATABASE", defaultCockroachDatabase),

		// Default retry configs, so retry every 1 second up to 5 times
		RetrySleep:    DefaultRetrySleep,
		RetryAttempts: DefaultRetryAttempts,
	}
}

// ###########################################
// NATS
// ###########################################

// DefaultNATSServerConfig creates a NATSServerConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultNATSServerConfig() *NATSServerConfig {
	httpPort, _ := strconv.Atoi(
		defaultEnvStr("NATS_HTTP_PORT", DefaultNATSServerHTTPPort),
	)
	port, _ := strconv.Atoi(
		defaultEnvStr("NATS_PORT", DefaultNATSServerPort),
	)
	return &NATSServerConfig{
		HTTPPort: httpPort,
		Port:     port,
	}
}

// ###########################################
// ClientConfig
// ###########################################

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		AuthToken:  defaultEnvStr("BUBBLY_TOKEN", DefaultClientAuthToken),
		BubblyAddr: defaultEnvStr("BUBBLY_ADDR", DefaultBubblyAddr),
	}
}

// ###########################################
// CLI
// ###########################################

func DefaultCLIConfig() *CLIConfig {
	return &CLIConfig{
		NoColor: defaultEnvBool("NO_COLOR", DefaultCLINoColorToggle),
	}
}
