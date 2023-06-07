// This code defines a configuration struct named "Config" and reads
// the configuration values from both a YAML file and environment
// variables using the "cleanenv" package.
package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// The "Config" struct has three nested structs: "App" containing
// the application name and version, "HTTP" containing the HTTP port
// to be used by the application, and "PG" containing the dialect and
// URL for the PostgreSQL database.
type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		PG   `yaml:"postgres"`
		Log  `yaml:"logger"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	PG struct {
		Dialect string `env-requered:"true" yaml:"dialect" env:"DIALECT"`
		URL     string `env-required:"true" yaml:"pg_url" env:"PG_URL"`
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}
)

// Creates a new config entity after reading the configuration values
// from the YAML file and environment variables.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
