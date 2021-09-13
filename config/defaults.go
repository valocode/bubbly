package config

import (
	"os"
	"strconv"

	"github.com/valocode/bubbly/auth"
)

// Default CLI configuration
const (
	DefaultCLINoColorToggle = false
	DefaultCLIDebugToggle   = false
)

// Default Bubbly Release configuration
const (
	DefaultBubblyDir      = ".bubbly"
	DefaultReleaseSpec    = "" // default spec is under bubbly dir
	DefaultReleaseProject = "default"
	DefaultOrganization   = "bubbly"
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
	DefaultBubblyAddr      = "http://localhost:8111"
	DefaultNATSAddr        = "localhost:4223"
)

// DefaultEnvStr reads a string value from the environment and falls back to the
// default value if not provided.
func DefaultEnvStr(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// DefaultEnvBool reads a boolean value from the environment and falls back to the
// default value if not provided.
func DefaultEnvBool(key string, defaultValue bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		val, _ := strconv.ParseBool(DefaultEnvStr(key, value))
		return val
	}
	return defaultValue
}

// DefaultReleaseConfig creates a ReleaseConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultReleaseConfig() *ReleaseConfig {
	return &ReleaseConfig{
		BubblyDir:   DefaultEnvStr("BUBBLY_DIR", DefaultBubblyDir),
		ReleaseSpec: DefaultEnvStr("BUBBLY_RELEASE", DefaultBubblyDir),
		Project:     DefaultEnvStr("BUBBLY_PROJECT", DefaultReleaseProject),
	}
}

// DefaultServerConfig creates a ServerConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: DefaultEnvStr("BUBBLY_HOST", DefaultServerHost),
		Port: DefaultEnvStr("BUBBLY_PORT", DefaultServerPort),
	}
}

// DefaultStoreConfig creates a StoreConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultStoreConfig() *StoreConfig {
	return &StoreConfig{
		// Default provider
		Provider: Provider(DefaultEnvStr("BUBBLY_DB", DefaultStoreProvider)),
		// Default configuration for Postgres
		PostgresAddr:     DefaultEnvStr("POSTGRES_ADDR", DefaultPostgresAddr),
		PostgresUser:     DefaultEnvStr("POSTGRES_USER", DefaultPostgresUser),
		PostgresPassword: DefaultEnvStr("POSTGRES_PASSWORD", DefaultPostgresPassword),
		PostgresDatabase: DefaultEnvStr("POSTGRES_DATABASE", DefaultPostgresDatabase),
		// Default configuration for CockroachDB
		CockroachAddr:     DefaultEnvStr("COCKROACH_ADDR", defaultCockroachAddr),
		CockroachUser:     DefaultEnvStr("COCKROACH_USER", defaultCockroachUser),
		CockroachPassword: DefaultEnvStr("COCKROACH_PASSWORD", defaultCockroachPassword),
		CockroachDatabase: DefaultEnvStr("COCKROACH_DATABASE", defaultCockroachDatabase),

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
		DefaultEnvStr("NATS_HTTP_PORT", DefaultNATSServerHTTPPort),
	)
	port, _ := strconv.Atoi(
		DefaultEnvStr("NATS_PORT", DefaultNATSServerPort),
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
		AuthToken:  DefaultEnvStr("BUBBLY_TOKEN", DefaultClientAuthToken),
		BubblyAddr: DefaultEnvStr("BUBBLY_ADDR", DefaultBubblyAddr),
	}
}

// ###########################################
// CLI
// ###########################################

func DefaultCLIConfig() *CLIConfig {
	return &CLIConfig{
		NoColor: DefaultEnvBool("NO_COLOR", DefaultCLINoColorToggle),
	}
}

// ###########################################
// Auth
// ###########################################

func DefaultAuthConfig() *auth.Config {
	return &auth.Config{
		ProviderURL:  DefaultEnvStr("AUTH_PROVIDER_URL", ""),
		ClientID:     DefaultEnvStr("AUTH_CLIENT_ID", ""),
		ClientSecret: DefaultEnvStr("AUTH_CLIENT_SECRET", ""),
		RedirectURL:  DefaultEnvStr("AUTH_REDIRECT_URL", ""),
		Scopes:       []string{"email"},
	}
}
