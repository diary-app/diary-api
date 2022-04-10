package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port string `env:"APP_PORT,required"`
}

func Read() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		panic(fmt.Sprintf("Failed to read config: %v", err))
	}

	return cfg

}
