package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// Server runs the HTTP API (health, readiness, status).
type Server struct {
	addr   string
	mu     sync.RWMutex
	ready  bool
	start  time.Time
	mux    *http.ServeMux
	server *http.Server
}

// New creates an HTTP server bound to addr (e.g. ":8080").
func New(addr string) *Server {
	s := &Server{addr: addr, start: time.Now()}
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/ready", s.handleReady)
	s.mux.HandleFunc("/api/status", s.handleStatus)
	s.mux.HandleFunc("/api/incidents", s.handleIncidents)
	s.server = &http.Server{Addr: addr, Handler: s.mux}
	return s
}

// SetReady marks the server as ready (e.g. after agent has started).
func (s *Server) SetReady(ready bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ready = ready
}

// Start listens and serves. Blocks until Shutdown is called.
func (s *Server) Start() error {
	s.SetReady(true)
	slog.Info("http server listening", "addr", s.addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleReady(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	ready := s.ready
	s.mu.RUnlock()
	if ready {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	} else {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
	}
}

func (s *Server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	ready := s.ready
	s.mu.RUnlock()
	out := map[string]interface{}{
		"status":   "running",
		"ready":    ready,
		"uptime_s": time.Since(s.start).Seconds(),
		"service":  "agentops-agent",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

// handleIncidents returns a small static list of example incidents.
// This is a starting point for a richer incidents API backed by the agent's state.
func (s *Server) handleIncidents(w http.ResponseWriter, _ *http.Request) {
	type incident struct {
		ID          string  `json:"id"`
		Kind        string  `json:"kind"`
		Severity    string  `json:"severity"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		ResourceID  string  `json:"resource_id"`
		Region      string  `json:"region,omitempty"`
		AgeSeconds  float64 `json:"age_seconds"`
	}

	now := time.Now()
	examples := []incident{
		{
			ID:          "example-1",
			Kind:        "incident",
			Severity:    "high",
			Title:       "High CPU on web-tier",
			Description: "CPU utilization > 90% for 5m",
			ResourceID:  "i-0123456789abcdef0",
			Region:      "us-east-1",
			AgeSeconds:  time.Since(now.Add(-5 * time.Minute)).Seconds(),
		},
		{
			ID:          "example-2",
			Kind:        "cost",
			Severity:    "medium",
			Title:       "Idle dev database",
			Description: "RDS instance idle for 24h; candidate for stop",
			ResourceID:  "rds-dev-001",
			Region:      "us-west-2",
			AgeSeconds:  time.Since(now.Add(-24 * time.Hour)).Seconds(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items": examples,
	})
}
