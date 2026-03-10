package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL   string
	SessionSecret string
	Port          string
	Env           string
	StaticPath    string
	SetupToken    string
}

func Load(configFile string) (*Config, error) {
	cfg := &Config{}

	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if idx := strings.Index(line, "="); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				value := strings.TrimSpace(line[idx+1:])
				value = strings.Trim(value, "\"'")

				switch key {
				case "DATABASE_URL":
					cfg.DatabaseURL = value
				case "SESSION_SECRET":
					cfg.SessionSecret = value
				case "PORT":
					cfg.Port = value
				case "ENV":
					cfg.Env = value
				case "STATIC_PATH":
					cfg.StaticPath = value
				case "SETUP_TOKEN":
					cfg.SetupToken = value
				}
			}
		}
	}

	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.DatabaseURL = v
	}
	if v := os.Getenv("SESSION_SECRET"); v != "" {
		cfg.SessionSecret = v
	}
	if v := os.Getenv("PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("ENV"); v != "" {
		cfg.Env = v
	}
	if v := os.Getenv("STATIC_PATH"); v != "" {
		cfg.StaticPath = v
	}
	if v := os.Getenv("SETUP_TOKEN"); v != "" {
		cfg.SetupToken = v
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
