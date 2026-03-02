package types

import "time"

// Source represents where an infrastructure event originated.
type Source string

const (
	SourceAWS   Source = "aws"
	SourceGCP  Source = "gcp"
	SourceK8s  Source = "k8s"
	SourceKafka Source = "kafka"
)

// Severity for prioritization.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// Event is a normalized infrastructure event from any source.
type Event struct {
	ID          string            `json:"id"`
	Source      Source            `json:"source"`
	Kind        string            `json:"kind"`        // incident, cost, security, deployment
	Severity    Severity          `json:"severity"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Raw         map[string]any    `json:"raw,omitempty"`
	ResourceID  string            `json:"resource_id"`
	Region      string            `json:"region,omitempty"`
	AccountID   string            `json:"account_id,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	Labels      map[string]string `json:"labels,omitempty"`
}
