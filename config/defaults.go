package config

import (
	"os"
	"strconv"
)

// Default CLI configuration
const (
	defaultCLIColorToggle = true
	DefaultDebugToggle    = false
)

// Default Bubbly API Server configuration
const (
	defaultAPIServerProtocol = "http"
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
	DefaultNATSServerAddr     = "localhost:4223"
)

// Default configuration for the bubbly agent components
const (
	DefaultAPIServerToggle  = false
	DefaultDataStoreToggle  = false
	DefaultWorkerToggle     = false
	DefaultNATSServerToggle = true
	DefaultDeploymentType   = SingleDeployment
)

func defaultEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// DefaultServerConfig creates a ServerConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Protocol: defaultEnv("BUBBLY_PROTOCOL", defaultAPIServerProtocol),
		Host:     defaultEnv("BUBBLY_HOST", DefaultAPIServerHost),
		Port:     defaultEnv("BUBBLY_PORT", DefaultAPIServerPort),
	}
}

// DefaultStoreConfig creates a StoreConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultStoreConfig() *StoreConfig {
	return &StoreConfig{
		// Default provider
		Provider: StoreProviderType(defaultEnv("BUBBLY_STORE_PROVIDER", DefaultStoreProvider)),
		// Default configuration for Postgres
		PostgresAddr:     defaultEnv("POSTGRES_ADDR", DefaultPostgresAddr),
		PostgresUser:     defaultEnv("POSTGRES_USER", DefaultPostgresUser),
		PostgresPassword: defaultEnv("POSTGRES_PASSWORD", DefaultPostgresPassword),
		PostgresDatabase: defaultEnv("POSTGRES_DATABASE", DefaultPostgresDatabase),
		// Default configuration for CockroachDB
		CockroachAddr:     defaultEnv("COCKROACH_ADDR", defaultCockroachAddr),
		CockroachUser:     defaultEnv("COCKROACH_USER", defaultCockroachUser),
		CockroachPassword: defaultEnv("COCKROACH_PASSWORD", defaultCockroachPassword),
		CockroachDatabase: defaultEnv("COCKROACH_DATABASE", defaultCockroachDatabase),

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
		StoreConfig:       DefaultStoreConfig(),
		NATSServerConfig:  DefaultNATSServerConfig(),
		EnabledComponents: DefaultAgentComponentsEnabled(),
		DeploymentType:    AgentDeploymentType(defaultEnv("AGENT_DEPLOYMENT_TYPE", DefaultDeploymentType.String())),
	}
}

// DefaultNATSServerConfig creates a NATSServerConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultNATSServerConfig() *NATSServerConfig {
	httpPort, _ := strconv.Atoi(
		defaultEnv("NATS_SERVER_HTTP_PORT", DefaultNATSServerHTTPPort),
	)
	port, _ := strconv.Atoi(
		defaultEnv("NATS_SERVER_PORT", DefaultNATSServerPort),
	)
	return &NATSServerConfig{
		HTTPPort: httpPort,
		Port:     port,
		Addr:     defaultEnv("NATS_SERVER_ADDR", DefaultNATSServerAddr),
	}
}

// DefaultAgentComponentsEnabled creates an AgentComponentsToggle struct
// instance with all components disabled
func DefaultAgentComponentsEnabled() *AgentComponentsToggle {
	apiServerToggle, _ := strconv.ParseBool(defaultEnv("AGENT_API_SERVER_TOGGLE", strconv.FormatBool(DefaultAPIServerToggle)))
	dataStoreToggle, _ := strconv.ParseBool(defaultEnv("AGENT_DATA_STORE_TOGGLE", strconv.FormatBool(DefaultDataStoreToggle)))
	workerToggle, _ := strconv.ParseBool(defaultEnv("AGENT_WORKER_TOGGLE", strconv.FormatBool(DefaultWorkerToggle)))
	natsServerToggle, _ := strconv.ParseBool(defaultEnv("AGENT_NATS_SERVER_TOGGLE", strconv.FormatBool(DefaultNATSServerToggle)))
	return &AgentComponentsToggle{
		APIServer:  apiServerToggle,
		DataStore:  dataStoreToggle,
		Worker:     workerToggle,
		NATSServer: natsServerToggle,
	}
}

func DefaultCLIConfig() *CLIConfig {
	color, _ := strconv.ParseBool(defaultEnv("COLOR", strconv.FormatBool(defaultCLIColorToggle)))
	return &CLIConfig{
		Color: color,
	}
}
