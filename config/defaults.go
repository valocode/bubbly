package config

// NewDefaultConfig creates a Config struct from default configurations
func NewDefaultConfig() *ServerConfig {
	return &ServerConfig{
		Protocol: "http",
		Port:     "8111",
		Auth:     false,
		Host:     "localhost",
	}
}
