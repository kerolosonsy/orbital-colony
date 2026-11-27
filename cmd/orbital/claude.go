package main

import (
	"context"
	"errors"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// On a hosted server the helper chat and the code score ask Claude through
// the host's API key (the local game uses the player's own Codex CLI).

const claudeModel = "claude-opus-5-5"

func claudeAsker(key string) func(context.Context, string) (string, error) {
	client := anthropic.NewClient(option.WithAPIKey(key))
	return func(ctx context.Context, prompt string) (string, error) {
		ctx, cancel := context.WithTimeout(ctx, chatTimeout)
		defer cancel()
		resp, err := client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
			Model:        claudeModel,
			MaxTokens:    16000,
			Betas:        []anthropic.AnthropicBeta{anthropic.AnthropicBetaServerSideFallback2026_07_01},
			Fallbacks:    anthropic.BetaFallbacksParamOfDefault(),
			OutputConfig: anthropic.BetaOutputConfigParam{Effort: anthropic.BetaOutputConfigEffortMedium},
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(prompt)),
			},
		})
		if err != nil {
			return "", err
		}
		if resp.StopReason == anthropic.BetaStopReasonRefusal {
			return "", errors.New("the model declined")
		}
		var b strings.Builder
		for _, block := range resp.Content {
			if text, ok := block.AsAny().(anthropic.BetaTextBlock); ok {
				b.WriteString(text.Text)
			}
		}
		if reply := strings.TrimSpace(b.String()); reply != "" {
			return reply, nil
		}
		return "", errors.New("empty reply")
	}
}
