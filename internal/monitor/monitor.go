package monitor

import "github.com/agentops/platform/internal/types"

// Stream provides a channel of infrastructure events (AWS, GCP, K8s, Kafka).
type Stream interface {
	// Events returns a channel that receives normalized events. Caller should range over it.
	Events() <-chan *types.Event
	// Start begins streaming; implementers use goroutines internally.
	Start() error
	// Stop gracefully stops the stream.
	Stop() error
}
