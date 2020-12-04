package config

import "fmt"

// ServerConfig is a struct storing the server information.
type ServerConfig struct {
	Protocol string
	Port     string
	Host     string `validate:"required"`
	Auth     bool
	Token    string
}

func (s ServerConfig) HostURL() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
