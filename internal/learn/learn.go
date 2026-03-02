package learn

import (
	"context"

	"github.com/agentops/platform/internal/types"
)

// Store persists outcomes for future decisions (vector DB + relational).
type Store interface {
	// RecordOutcome saves event + plan + audit for learning and audit trail.
	RecordOutcome(ctx context.Context, evt *types.Event, plan *types.Plan, audit []types.AuditEntry) error
	// SearchSimilar finds past similar events (e.g. via vector search).
	SearchSimilar(ctx context.Context, evt *types.Event, limit int) ([]Outcome, error)
}

// Outcome is a past incident outcome for retrieval.
type Outcome struct {
	EventID   string
	PlanID    string
	Reasoning string
	Success   bool
}
