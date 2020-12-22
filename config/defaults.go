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

// DefaultServerConfig creates a Config struct from default configurations
func DefaultServerConfig() *ServerConfig {
	auth, _ := strconv.ParseBool(defaultEnv("BUBBLY_ENABLE_AUTH", "false"))
	return &ServerConfig{
		Protocol: defaultEnv("BUBBLY_PROTOCOL", "http"),
		Host:     defaultEnv("BUBBLY_HOST", "localhost"),
		Port:     defaultEnv("BUBBLY_PORT", "8111"),
		Auth:     auth,
	}
}

func DefaultStoreConfig() *StoreConfig {
	return &StoreConfig{
		Provider:         StoreProviderType(defaultEnv("BUBBLY_STORE_PROVIDER", string(PostgresStore))),
		PostgresAddr:     defaultEnv("POSTGRES_ADDR", "postgres:5432"),
		PostgresUser:     defaultEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: defaultEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDatabase: defaultEnv("POSTGRES_DATABASE", "bubbly"),
	}
}

func DefaultResourceConfig() *ResourceConfig {
	return &ResourceConfig{
		Provider: ResourceProviderType(defaultEnv("BUBBLY_RESOURCE_PROVIDER", string(BuntdbResourceProvider))),
	}
}
