package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds agent configuration from env.
type Config struct {
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

// Load reads .env and builds Config.
func Load() (*Config, error) {
	_ = godotenv.Load()

	intervalSec := 30
	if s := os.Getenv("STUB_EVENT_INTERVAL_SEC"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			intervalSec = n
		}
	}

	return &Config{
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
