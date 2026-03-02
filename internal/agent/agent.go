package agent

import (
	"context"
	"log/slog"
	"time"

	"github.com/agentops/platform/internal/analyze"
	"github.com/agentops/platform/internal/execute"
	"github.com/agentops/platform/internal/learn"
	"github.com/agentops/platform/internal/monitor"
	"github.com/agentops/platform/internal/types"
	"github.com/agentops/platform/internal/verify"
	"github.com/google/uuid"
)

// Agent runs the core loop: Monitor → Analyze → Plan → Execute → Verify → Learn.
type Agent struct {
	stream   monitor.Stream
	analyzer analyze.Analyzer
	exec     execute.Executor
	verifier verify.Verifier
	learn    learn.Store
}

// New builds an agent with the given components.
func New(
	stream monitor.Stream,
	analyzer analyze.Analyzer,
	exec execute.Executor,
	verifier verify.Verifier,
	learnStore learn.Store,
) *Agent {
	return &Agent{
		stream:   stream,
		analyzer: analyzer,
		exec:     exec,
		verifier: verifier,
		learn:    learnStore,
	}
}

// Run starts the monitor and processes each event through the pipeline. Blocks until ctx is cancelled.
func (a *Agent) Run(ctx context.Context) error {
	if err := a.stream.Start(); err != nil {
		return err
	}
	defer a.stream.Stop()

	events := a.stream.Events()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case evt, ok := <-events:
			if !ok {
				return nil
			}
			a.processEvent(ctx, evt)
		}
	}
}

func (a *Agent) processEvent(ctx context.Context, evt *types.Event) {
	slog.Info("processing event", "event_id", evt.ID, "title", evt.Title)

	plan, err := a.analyzer.Analyze(ctx, evt)
	if err != nil {
		slog.Error("analyze failed", "event_id", evt.ID, "error", err)
		return
	}
	if plan == nil || len(plan.Actions) == 0 {
		slog.Info("no actions for event", "event_id", evt.ID)
		return
	}

	var audit []types.AuditEntry
	for _, action := range plan.Actions {
		entry := types.AuditEntry{
			ID:         uuid.New().String(),
			EventID:    evt.ID,
			PlanID:     plan.ID,
			ActionID:   action.ID,
			ActionType: action.Type,
			Reasoning:  plan.Reasoning,
			ResourceID: action.ResourceID,
			Timestamp:  time.Now(),
		}
		start := time.Now()
		result, err := a.exec.Execute(ctx, action, evt)
		entry.DurationMs = time.Since(start).Milliseconds()
		if err != nil {
			entry.Status = "failure"
			entry.Details = err.Error()
			audit = append(audit, entry)
			slog.Error("execute failed", "action_id", action.ID, "error", err)
			if rerr := a.exec.Rollback(ctx, action, evt); rerr != nil {
				slog.Error("rollback failed", "action_id", action.ID, "error", rerr)
			}
			continue
		}
		entry.Status = "success"
		entry.Details = result

		if verr := a.verifier.Verify(ctx, evt, action, result); verr != nil {
			entry.Status = "rolled_back"
			entry.RolledBack = true
			entry.Details = result + "; verify failed: " + verr.Error()
			audit = append(audit, entry)
			slog.Warn("verify failed, rolling back", "action_id", action.ID, "error", verr)
			_ = a.exec.Rollback(ctx, action, evt)
			continue
		}
		audit = append(audit, entry)
	}

	if err := a.learn.RecordOutcome(ctx, evt, plan, audit); err != nil {
		slog.Error("record outcome failed", "event_id", evt.ID, "error", err)
	}
}
