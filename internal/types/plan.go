package types

import "time"

// Action represents a single remediation step.
type Action struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`   // e.g. scale, restart, patch, rollback
	ResourceID  string            `json:"resource_id"`
	Params      map[string]string `json:"params"`
	Reason      string            `json:"reason"`
	RollbackDef string            `json:"rollback_def,omitempty"` // how to undo
}

// Plan is a multi-step remediation plan from the Analyze phase.
type Plan struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	Actions   []Action  `json:"actions"`
	Reasoning string    `json:"reasoning"`
	CreatedAt time.Time `json:"created_at"`
}
