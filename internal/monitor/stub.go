package monitor

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/agentops/platform/internal/types"
	"github.com/google/uuid"
)

// StubStream is a placeholder that emits sample events for local development.
type StubStream struct {
	interval time.Duration
	events   chan *types.Event
	done     chan struct{}
	once     sync.Once
}

// NewStubStream creates a stream that emits a sample event every interval.
func NewStubStream(interval time.Duration) *StubStream {
	return &StubStream{
		interval: interval,
		events:   make(chan *types.Event, 64),
		done:     make(chan struct{}),
	}
}

// Events returns the event channel.
func (s *StubStream) Events() <-chan *types.Event {
	return s.events
}

// Start begins emitting events in a goroutine.
func (s *StubStream) Start() error {
	s.once.Do(func() {
		go s.run()
	})
	return nil
}

func (s *StubStream) run() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			close(s.events)
			return
		case <-ticker.C:
			evt := &types.Event{
				ID:          uuid.New().String(),
				Source:      types.SourceAWS,
				Kind:        "incident",
				Severity:    types.SeverityHigh,
				Title:       "High CPU on web-tier",
				Description: "CPU utilization > 90% for 5m",
				ResourceID:  "i-0123456789abcdef0",
				Region:      "us-east-1",
				Timestamp:   time.Now(),
			}
			slog.Info("stub event emitted", "event_id", evt.ID, "title", evt.Title)
			select {
			case s.events <- evt:
			case <-s.done:
				close(s.events)
				return
			}
		}
	}
}

// Stop gracefully stops the stream.
func (s *StubStream) Stop() error {
	select {
	case <-s.done:
		return nil
	default:
		close(s.done)
	}
	return nil
}

// Ensure StubStream implements Stream (compile-time check).
var _ Stream = (*StubStream)(nil)
