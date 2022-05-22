package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppPort string     `yaml:"appPort"`
	PG      DbConfig   `yaml:"postgres"`
	Auth    AuthConfig `yaml:"auth"`
}

type AuthConfig struct {
	JwtKey string `env-required:"true" yaml:"jwtKey" env:"JWT_KEY"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	DbName   string `yaml:"dbName"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `env-required:"true" yaml:"password" env:"DB_PASSWORD"`
}

func Read() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig("./internal/config/config.yml", cfg); err != nil {
		return nil, fmt.Errorf("read config from yaml failed: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config from env - error: %w", err)
	}

	return cfg, nil
}
