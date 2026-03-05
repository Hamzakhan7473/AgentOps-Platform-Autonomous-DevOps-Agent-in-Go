package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agentops/platform/internal/agent"
	"github.com/agentops/platform/internal/analyze"
	"github.com/agentops/platform/internal/execute"
	"github.com/agentops/platform/internal/learn"
	"github.com/agentops/platform/internal/monitor"
	"github.com/agentops/platform/internal/server"
	"github.com/agentops/platform/internal/verify"
	"github.com/agentops/platform/pkg/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		slog.Error("invalid config", "error", err)
		os.Exit(1)
	}

	// Stub stream for development (replace with Kafka/EventBridge in production)
	stream := monitor.NewStubStream(cfg.StubEventInterval)

	// LLM analyzer (skip if no API key — use stub that returns no-op plan)
	var analyzer analyze.Analyzer
	if cfg.LLMAPIKey != "" {
		analyzer = analyze.NewLLMAnalyzer(cfg.LLMBaseURL, cfg.LLMAPIKey, cfg.LLMModel)
	} else {
		analyzer = analyze.NewNoopAnalyzer()
	}

	exec := &execute.StubExecutor{}
	verifier := &verify.StubVerifier{}
	learnStore := learn.NewMemoryStore()

	ag := agent.New(stream, analyzer, exec, verifier, learnStore)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// HTTP server for health, readiness, and status (K8s/Docker-friendly)
	srv := server.New(fmt.Sprintf(":%d", cfg.HTTPPort))
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "error", err)
		}
	}()

	// Run agent loop in background
	go func() {
		slog.Info("agent starting")
		if err := ag.Run(ctx); err != nil && err != context.Canceled {
			slog.Error("agent stopped", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("http shutdown", "error", err)
	}
	slog.Info("agent stopped")
}
