package config

// ServerConfig is a struct storing the server information.
type ServerConfig struct {
	Protocol string
	Port     string
	Host     string `validate:"required"`
	Auth     bool
	Token    string
}

// Config is a struct storing all config information.
type Config struct {
	ServerConfig *ServerConfig
}
