package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr              string        `yaml:"addr" env:"SERVER_ADDR" env-default:":8080"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env:"SERVER_READ_HEADER_TIMEOUT" env-default:"2s"`
	ReadTimeout       time.Duration `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout      time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env:"SERVER_IDLE_TIMEOUT" env-default:"30s"`
	MaxHeaderBytes    int           `yaml:"max_header_bytes" env:"SERVER_MAX_HEADER_BYTES" env-default:"1048576"`
}

type Config struct {
	Server HTTPServer `yaml:"server" env:"server" env-required:"true"`
}

func MustLoad() (*Config, error) {
	const fn = "config.MustLoad"

	cfg := new(Config)

	err := cleanenv.ReadConfig("internal/config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: during read config.yaml: %w", fn, err)
	}

	err = cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		err = cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, fmt.Errorf("%s: during read .env: %w", fn, err)
		}
	}

	return cfg, nil
}
