package analyze

import (
	"context"
	"time"

	"github.com/agentops/platform/internal/types"
	"github.com/google/uuid"
)

// NoopAnalyzer returns an empty plan (no actions). Use when LLM is not configured.
type NoopAnalyzer struct{}

// NewNoopAnalyzer creates a no-op analyzer.
func NewNoopAnalyzer() *NoopAnalyzer {
	return &NoopAnalyzer{}
}

// Analyze implements Analyzer.
func (a *NoopAnalyzer) Analyze(ctx context.Context, evt *types.Event) (*types.Plan, error) {
	return &types.Plan{
		ID:        uuid.New().String(),
		EventID:   evt.ID,
		Reasoning: "no-op analyzer: no LLM configured",
		Actions:   nil,
		CreatedAt: time.Now(),
	}, nil
}
