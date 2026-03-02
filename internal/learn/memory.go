package learn

import (
	"context"
	"sync"

	"github.com/agentops/platform/internal/types"
)

// MemoryStore is an in-memory implementation for development (no vector search).
type MemoryStore struct {
	mu      sync.RWMutex
	entries []outcomeEntry
}

type outcomeEntry struct {
	evt   *types.Event
	plan  *types.Plan
	audit []types.AuditEntry
}

// NewMemoryStore creates an in-memory learn store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{entries: make([]outcomeEntry, 0)}
}

// RecordOutcome implements Store.
func (m *MemoryStore) RecordOutcome(ctx context.Context, evt *types.Event, plan *types.Plan, audit []types.AuditEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, outcomeEntry{evt: evt, plan: plan, audit: audit})
	return nil
}

// SearchSimilar returns recent outcomes (no vector similarity in stub).
func (m *MemoryStore) SearchSimilar(ctx context.Context, evt *types.Event, limit int) ([]Outcome, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n := len(m.entries)
	if limit <= 0 || limit > n {
		limit = n
	}
	start := n - limit
	if start < 0 {
		start = 0
	}
	out := make([]Outcome, 0, limit)
	for i := n - 1; i >= start && len(out) < limit; i-- {
		e := m.entries[i]
		success := true
		for _, a := range e.audit {
			if a.Status == "failure" || a.RolledBack {
				success = false
				break
			}
		}
		out = append(out, Outcome{
			EventID:   e.evt.ID,
			PlanID:    e.plan.ID,
			Reasoning: e.plan.Reasoning,
			Success:   success,
		})
	}
	return out, nil
}
