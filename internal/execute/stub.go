package execute

import (
	"context"
	"log/slog"
	"time"

	"github.com/agentops/platform/internal/types"
)

// StubExecutor logs actions and simulates success for development.
type StubExecutor struct{}

// Execute implements Executor.
func (e *StubExecutor) Execute(ctx context.Context, action types.Action, evt *types.Event) (string, error) {
	slog.Info("execute (stub)", "action_id", action.ID, "type", action.Type, "resource_id", action.ResourceID)
	time.Sleep(100 * time.Millisecond) // simulate work
	return "stub success", nil
}

// Rollback implements Executor.
func (e *StubExecutor) Rollback(ctx context.Context, action types.Action, evt *types.Event) error {
	slog.Info("rollback (stub)", "action_id", action.ID)
	return nil
}
