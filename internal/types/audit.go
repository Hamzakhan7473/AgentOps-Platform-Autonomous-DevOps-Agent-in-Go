package types

import "time"

// AuditEntry records every action for compliance.
type AuditEntry struct {
	ID          string    `json:"id"`
	EventID     string    `json:"event_id"`
	PlanID      string    `json:"plan_id"`
	ActionID    string    `json:"action_id"`
	ActionType  string    `json:"action_type"`
	Reasoning   string    `json:"reasoning"`
	Status      string    `json:"status"` // pending, success, failure, rolled_back
	Details     string    `json:"details,omitempty"`
	ResourceID  string    `json:"resource_id"`
	Timestamp   time.Time `json:"timestamp"`
	DurationMs  int64     `json:"duration_ms,omitempty"`
	RolledBack  bool      `json:"rolled_back"`
}
