package agent

import (
	"context"
	"fmt"
	"slices"
	"strings"

	herrors "dari/harness/errors"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

const maxTokens = 4096

func Send(ctx context.Context, client Client, transcript []anthropic.MessageParam, input string) ([]anthropic.MessageParam, string, Usage, error) {
	request := append(slices.Clone(transcript), anthropic.NewUserMessage(anthropic.NewTextBlock(input)))
	message, err := client.api.Messages.New(ctx, newParams(client.model, request))
	if err != nil {
		return transcript, "", Usage{}, herrors.New(herrors.APIRequest, err)
	}
	if err := checkStopReason(message.StopReason); err != nil {
		return transcript, "", Usage{}, err
	}
	return append(request, message.ToParam()), replyText(message), usageFrom(message.Usage), nil
}

func newParams(model string, messages []anthropic.MessageParam) anthropic.MessageNewParams {
	return anthropic.MessageNewParams{
		Model:     model,
		MaxTokens: maxTokens,
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
