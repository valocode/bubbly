package config

// NewDefaultConfig creates a Config struct from default configurations
func NewDefaultConfig() *Config {
	return &Config{
		ServerConfig: newDefaultServerConfig(),
	}
}

func newDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Protocol: "http",
		Port:     "8080",
		Auth:     false,
		Host:     "localhost",
	}
}
