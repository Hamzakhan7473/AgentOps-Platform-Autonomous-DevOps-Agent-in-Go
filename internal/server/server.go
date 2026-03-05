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
		"status":    "running",
		"ready":     ready,
		"uptime_s":  time.Since(s.start).Seconds(),
		"service":   "agentops-agent",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}
