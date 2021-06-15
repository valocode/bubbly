package config

import (
	"os"
	"strconv"
)

// Default CLI configuration
const (
	DefaultCLIColorToggle = true
	DefaultDebugToggle    = false
)

// Default Bubbly API Server configuration
const (
	DefaultAPIServerProtocol = "http"
	DefaultAPIServerHost     = "127.0.0.1"
	DefaultAPIServerPort     = "8111"
)

// Default store configuration
const (
	DefaultStoreProvider = "postgres"
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

// Default configuration for the bubbly agent components
const (
	DefaultAPIServerToggle  = false
	DefaultDataStoreToggle  = false
	DefaultWorkerToggle     = false
	DefaultNATSServerToggle = true
	DefaultDeploymentType   = SingleDeployment
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
		Protocol: defaultEnvStr("BUBBLY_PROTOCOL", DefaultAPIServerProtocol),
		Host:     defaultEnvStr("BUBBLY_HOST", DefaultAPIServerHost),
		Port:     defaultEnvStr("BUBBLY_PORT", DefaultAPIServerPort),
	}
}

// DefaultStoreConfig creates a StoreConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultStoreConfig() *StoreConfig {
	return &StoreConfig{
		// Default provider
		Provider: StoreProviderType(defaultEnvStr("BUBBLY_STORE_PROVIDER", DefaultStoreProvider)),
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
// Agent
// ###########################################

// DefaultAgentConfig creates an AgentConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		NATSServerConfig:  DefaultNATSServerConfig(),
		EnabledComponents: DefaultAgentComponentsEnabled(),
		DeploymentType:    AgentDeploymentType(defaultEnvStr("AGENT_DEPLOYMENT_TYPE", DefaultDeploymentType.String())),
	}
}

// DefaultAgentComponentsEnabled creates an AgentComponentsToggle struct
// instance with all components disabled
func DefaultAgentComponentsEnabled() *AgentComponentsToggle {
	return &AgentComponentsToggle{
		APIServer:  defaultEnvBool("AGENT_API_SERVER_TOGGLE", DefaultAPIServerToggle),
		DataStore:  defaultEnvBool("AGENT_DATA_STORE_TOGGLE", DefaultDataStoreToggle),
		Worker:     defaultEnvBool("AGENT_WORKER_TOGGLE", DefaultWorkerToggle),
		NATSServer: defaultEnvBool("AGENT_NATS_SERVER_TOGGLE", DefaultNATSServerToggle),
	}
}

// ###########################################
// Auth
// ###########################################

func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		Authentication: defaultEnvBool("BUBBLY_AUTHENTICATION", false),
		MultiTenancy:   defaultEnvBool("BUBBLY_MULTITENANCY", false),
		AuthAddr:       defaultEnvStr("BUBBLY_AUTH_API", "http://bubbly-auth:1323/api/v1"),
	}
}

// ###########################################
// NATS
// ###########################################

// DefaultNATSServerConfig creates a NATSServerConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultNATSServerConfig() *NATSServerConfig {
	httpPort, _ := strconv.Atoi(
		defaultEnvStr("NATS_SERVER_HTTP_PORT", DefaultNATSServerHTTPPort),
	)
	port, _ := strconv.Atoi(
		defaultEnvStr("NATS_SERVER_PORT", DefaultNATSServerPort),
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
		ClientType: HTTPClientType,
		AuthToken:  defaultEnvStr("BUBBLY_TOKEN", DefaultClientAuthToken),
		BubblyAddr: defaultEnvStr("BUBBLY_ADDR", DefaultBubblyAddr),
		NATSAddr:   defaultEnvStr("BUBBLY_NATS_ADDR", DefaultNATSAddr),
	}
}

// ###########################################
// CLI
// ###########################################

func DefaultCLIConfig() *CLIConfig {
	return &CLIConfig{
		Color: defaultEnvBool("COLOR", DefaultCLIColorToggle),
	}
}
