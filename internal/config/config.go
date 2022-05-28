package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppPort string     `yaml:"appPort" env-required:"false" env:"APP_PORT"`
	PG      DbConfig   `yaml:"postgres"`
	Auth    AuthConfig `yaml:"auth"`
}

type AuthConfig struct {
	JwtKey string `yaml:"jwtKey" env-required:"true" env:"JWT_KEY"`
}

type DbConfig struct {
	Host     string `yaml:"host" env-required:"false" env:"DB_HOST"`
	Port     int    `yaml:"port" env-required:"false" env:"DB_PORT"`
	DbName   string `yaml:"dbName" env-required:"false" env:"DB_NAME"`
	User     string `yaml:"user" env-required:"false" env:"DB_USER"`
	Password string `yaml:"password" env-required:"true" env:"DB_PASSWORD"`
	SslMode  string `yaml:"sslMode" env-required:"false" env:"DB_SSL_MODE" env-default:"disable"`
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
