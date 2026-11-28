package agent

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"dari/harness/api"
	herrors "dari/harness/errors"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

func Send(ctx context.Context, client Client, transcript []api.Message, input string) ([]api.Message, string, api.Usage, error) {
	request := append(slices.Clone(transcript), api.Message{Role: api.User, Text: input})
	message, err := client.api.Messages.New(ctx, newParams(client, toParams(request)))
	if err != nil {
		return transcript, "", api.Usage{}, herrors.New(herrors.APIRequest, err)
	}
	if err := checkStopReason(message.StopReason); err != nil {
		return transcript, "", api.Usage{}, err
	}
	reply := replyText(message)
	return append(request, api.Message{Role: api.Assistant, Text: reply}), reply, usageFrom(message.Usage), nil
}

func toParams(messages []api.Message) []anthropic.MessageParam {
	params := make([]anthropic.MessageParam, 0, len(messages))
	for _, message := range messages {
		params = append(params, toParam(message))
	}
	return params
}

func toParam(message api.Message) anthropic.MessageParam {
	block := anthropic.NewTextBlock(message.Text)
	if message.Role == api.Assistant {
		return anthropic.NewAssistantMessage(block)
	}
	return anthropic.NewUserMessage(block)
}

func newParams(client Client, messages []anthropic.MessageParam) anthropic.MessageNewParams {
	return anthropic.MessageNewParams{
		Model:     client.model,
		MaxTokens: client.maxTokens,
		Messages:  messages,
	}
}

func checkStopReason(reason anthropic.StopReason) error {
	switch reason {
	case anthropic.StopReasonEndTurn, anthropic.StopReasonStopSequence:
		return nil
	case anthropic.StopReasonMaxTokens:
		return herrors.New(herrors.Truncated, nil)
	case anthropic.StopReasonToolUse:
		return herrors.New(herrors.UnsupportedToolUse, nil)
	case anthropic.StopReasonRefusal:
		return herrors.New(herrors.Refused, nil)
	default:
		return herrors.New(herrors.UnexpectedStop, fmt.Errorf("stop reason %q", reason))
	}
}

func replyText(message *anthropic.Message) string {
	var text strings.Builder
	for _, block := range message.Content {
		if textBlock, ok := block.AsAny().(anthropic.TextBlock); ok {
			text.WriteString(textBlock.Text)
		}
	}
	return text.String()
}
