package analyze

import (
	"context"

	"github.com/agentops/platform/internal/types"
)

// Analyzer classifies and prioritizes events and produces a remediation plan.
type Analyzer interface {
	// Analyze takes an event and returns a plan (may be empty if no action needed).
	Analyze(ctx context.Context, evt *types.Event) (*types.Plan, error)
}
