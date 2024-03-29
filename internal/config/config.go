package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppPort  string     `yaml:"appPort" env-required:"false" env:"PORT"`
	PG       DbConfig   `yaml:"postgres"`
	Auth     AuthConfig `yaml:"auth"`
	LogDebug bool       `env:"LOG_DEBUG" env-default:"false"`
}

type AuthConfig struct {
	JwtKey        string `yaml:"jwtKey" env-required:"true" env:"JWT_KEY"`
	JwtTtlMinutes int    `yaml:"jwtTtlMinutes" env-default:"60"`
}

type DbConfig struct {
	Host     string `yaml:"host" env-required:"false" env:"DB_HOST"`
	Port     int    `yaml:"port" env-required:"false" env:"DB_PORT"`
	DbName   string `yaml:"dbName" env-required:"false" env:"DB_NAME"`
	User     string `yaml:"user" env-required:"false" env:"DB_USER"`
	Password string `yaml:"password" env-required:"false" env:"DB_PASSWORD"`
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
