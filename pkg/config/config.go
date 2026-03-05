package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds agent configuration from env.
type Config struct {
	// HTTP server
	HTTPPort int

	// LLM (OpenAI-compatible)
	LLMBaseURL string
	LLMAPIKey  string
	LLMModel   string

	// Monitor
	StubEventInterval time.Duration

	// Redis / PostgreSQL (for production state) — optional for v1
	RedisURL    string
	PostgresDSN string
}

// Validate returns an error if required settings are missing or invalid.
func (c *Config) Validate() error {
	if c.HTTPPort <= 0 || c.HTTPPort > 65535 {
		return errors.New("HTTP_PORT must be between 1 and 65535")
	}
	if c.StubEventInterval <= 0 {
		return fmt.Errorf("invalid STUB_EVENT_INTERVAL_SEC: %v", c.StubEventInterval)
	}
	return nil
}

// Load reads .env and builds Config.
func Load() (*Config, error) {
	_ = godotenv.Load()

	intervalSec := 30
	if s := os.Getenv("STUB_EVENT_INTERVAL_SEC"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			intervalSec = n
		}
	}

	port := 8080
	if s := os.Getenv("HTTP_PORT"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			port = n
		}
	}

	return &Config{
		HTTPPort:          port,
		LLMBaseURL:        os.Getenv("LLM_BASE_URL"),
		LLMAPIKey:         os.Getenv("LLM_API_KEY"),
		LLMModel:          getEnv("LLM_MODEL", "gpt-4o-mini"),
		StubEventInterval: time.Duration(intervalSec) * time.Second,
		RedisURL:          os.Getenv("REDIS_URL"),
		PostgresDSN:       os.Getenv("POSTGRES_DSN"),
	}, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
