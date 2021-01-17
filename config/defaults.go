package config

import (
	"os"
	"strconv"
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
	auth, _ := strconv.ParseBool(defaultEnv("BUBBLY_ENABLE_AUTH", "false"))
	return &ServerConfig{
		Protocol: defaultEnv("BUBBLY_PROTOCOL", "http"),
		Host:     defaultEnv("BUBBLY_HOST", "localhost"),
		Port:     defaultEnv("BUBBLY_PORT", "8111"),
		Auth:     auth,
	}
}

// DefaultStoreConfig creates a StoreConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultStoreConfig() *StoreConfig {
	return &StoreConfig{
		// Default provider
		Provider: StoreProviderType(defaultEnv("BUBBLY_STORE_PROVIDER", string(PostgresStore))),
		// Default configuration for Postgres
		PostgresAddr:     defaultEnv("POSTGRES_ADDR", "postgres:5432"),
		PostgresUser:     defaultEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: defaultEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDatabase: defaultEnv("POSTGRES_DATABASE", "bubbly"),
		// Default configuration for CockroachDB
		CockroachAddr:     defaultEnv("COCKROACH_ADDR", "cockroachdb:26257"),
		CockroachUser:     defaultEnv("COCKROACH_USER", "root"),
		CockroachPassword: defaultEnv("COCKROACH_PASSWORD", "admin"),
		CockroachDatabase: defaultEnv("COCKROACH_DATABASE", "defaultdb"),
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
	}
}

// DefaultNATSServerConfig creates a NATSServerConfig struct from defaults
// or, preferentially, from provided environment variables.
func DefaultNATSServerConfig() *NATSServerConfig {
	httpPort, _ := strconv.Atoi(
		defaultEnv(
			"NATS_SERVER_HTTP_PORT",
			"8222"),
	)
	port, _ := strconv.Atoi(
		defaultEnv(
			"NATS_SERVER_PORT",
			"4223"),
	)
	return &NATSServerConfig{
		HTTPPort: httpPort,
		Port:     port,
		Addr:     defaultEnv("NATS_SERVER_ADDR", "localhost:4223"),
	}
}

// DefaultAgentComponentsEnabled creates an AgentComponentsToggle struct
// instance with all components disabled
func DefaultAgentComponentsEnabled() *AgentComponentsToggle {
	return &AgentComponentsToggle{
		UI:        false,
		APIServer: false,
		DataStore: false,
		Worker:    false,
	}
}
