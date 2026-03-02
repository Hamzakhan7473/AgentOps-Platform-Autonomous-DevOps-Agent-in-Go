package execute

import (
	"context"

	"github.com/agentops/platform/internal/types"
)

// Executor runs remediation actions (e.g. via AWS SDK, kubectl, GCP client).
type Executor interface {
	// Execute runs a single action and returns outcome (error nil = success).
	Execute(ctx context.Context, action types.Action, evt *types.Event) (result string, err error)
	// Rollback undoes an action when verification fails.
	Rollback(ctx context.Context, action types.Action, evt *types.Event) error
}
