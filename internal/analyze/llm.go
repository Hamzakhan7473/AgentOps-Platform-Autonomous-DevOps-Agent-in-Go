package analyze

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/agentops/platform/internal/types"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

// LLMAnalyzer uses an OpenAI-compatible API to classify events and generate plans.
type LLMAnalyzer struct {
	client *openai.Client
	model  string
}

// NewLLMAnalyzer creates an analyzer that calls the given OpenAI-compatible endpoint.
func NewLLMAnalyzer(baseURL, apiKey, model string) *LLMAnalyzer {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	return &LLMAnalyzer{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}
}

// Analyze implements Analyzer.
func (a *LLMAnalyzer) Analyze(ctx context.Context, evt *types.Event) (*types.Plan, error) {
	prompt := a.buildPrompt(evt)
	resp, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: a.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Temperature: 0.2,
	})
	if err != nil {
		return nil, fmt.Errorf("llm completion: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("empty llm response")
	}
	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	return a.parsePlan(content, evt.ID)
}

const systemPrompt = `You are an autonomous DevOps agent. Given an infrastructure event (incident, cost, security, or deployment), you must:
1. Classify severity and root cause.
2. Output a JSON remediation plan with this exact structure (no markdown, no code fence):
{"reasoning":"...","actions":[{"type":"...","resource_id":"...","params":{},"reason":"...","rollback_def":"..."}]}
If no action is needed, return {"reasoning":"...","actions":[]}.
Types: scale, restart, patch, rollback, terminate, resize, rotate_secret, etc.`

func (a *LLMAnalyzer) buildPrompt(evt *types.Event) string {
	return fmt.Sprintf("Event: source=%s kind=%s severity=%s title=%s description=%s resource_id=%s region=%s",
		evt.Source, evt.Kind, evt.Severity, evt.Title, evt.Description, evt.ResourceID, evt.Region)
}

func (a *LLMAnalyzer) parsePlan(content string, eventID string) (*types.Plan, error) {
	// Strip markdown code block if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSpace(content)
	var raw struct {
		Reasoning string `json:"reasoning"`
		Actions   []struct {
			Type       string            `json:"type"`
			ResourceID string            `json:"resource_id"`
			Params     map[string]string `json:"params"`
			Reason     string            `json:"reason"`
			RollbackDef string           `json:"rollback_def"`
		} `json:"actions"`
	}
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return nil, fmt.Errorf("parse plan json: %w", err)
	}
	plan := &types.Plan{
		ID:        uuid.New().String(),
		EventID:   eventID,
		Reasoning: raw.Reasoning,
		Actions:   make([]types.Action, 0, len(raw.Actions)),
		CreatedAt: time.Now(),
	}
	for _, ac := range raw.Actions {
		if ac.ResourceID == "" {
			ac.ResourceID = "unknown"
		}
		plan.Actions = append(plan.Actions, types.Action{
			ID:          uuid.New().String(),
			Type:        ac.Type,
			ResourceID:  ac.ResourceID,
			Params:      ac.Params,
			Reason:      ac.Reason,
			RollbackDef: ac.RollbackDef,
		})
	}
	return plan, nil
}
