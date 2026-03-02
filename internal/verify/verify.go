package verify

import (
	"context"

	"github.com/agentops/platform/internal/types"
)

// Verifier checks that remediation succeeded (e.g. re-check metrics, health).
type Verifier interface {
	// Verify returns nil if the issue is resolved; otherwise error indicates failure (trigger rollback).
	Verify(ctx context.Context, evt *types.Event, action types.Action, result string) error
}
