package verify

import (
	"context"

	"github.com/agentops/platform/internal/types"
)

// StubVerifier always reports success (no real checks).
type StubVerifier struct{}

// Verify implements Verifier.
func (v *StubVerifier) Verify(ctx context.Context, evt *types.Event, action types.Action, result string) error {
	return nil
}
