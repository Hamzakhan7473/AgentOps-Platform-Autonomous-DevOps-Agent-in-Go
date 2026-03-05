package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_defaults(t *testing.T) {
	// Clear env to get defaults
	os.Unsetenv("HTTP_PORT")
	os.Unsetenv("STUB_EVENT_INTERVAL_SEC")
	os.Unsetenv("LLM_MODEL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTPPort != 8080 {
		t.Errorf("HTTPPort: got %d, want 8080", cfg.HTTPPort)
	}
	if cfg.LLMModel != "gpt-4o-mini" {
		t.Errorf("LLMModel: got %q, want gpt-4o-mini", cfg.LLMModel)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{"valid", &Config{HTTPPort: 8080, StubEventInterval: time.Second}, false},
		{"port zero", &Config{HTTPPort: 0, StubEventInterval: time.Second}, true},
		{"port too high", &Config{HTTPPort: 70000, StubEventInterval: time.Second}, true},
		{"interval zero", &Config{HTTPPort: 8080, StubEventInterval: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
