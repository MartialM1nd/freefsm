package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL   string
	SessionSecret string
	Port          string
	Env           string
	StaticPath    string
	SetupToken    string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		Port:          os.Getenv("PORT"),
		Env:           os.Getenv("ENV"),
		StaticPath:    os.Getenv("STATIC_PATH"),
		SetupToken:    os.Getenv("SETUP_TOKEN"),
	}

	// Defaults
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.Env == "" {
		cfg.Env = "development"
	}
	if cfg.StaticPath == "" {
		cfg.StaticPath = "ui/static"
	}

	// Required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is required")
	}

	return cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}
