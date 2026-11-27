package parser

import (
	"context"

	herrors "dari/harness/errors"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

// TODO: context management for management requests
func CreateMessage(client anthropic.Client, model string, input string) (*anthropic.Message, error) {
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(input)),
	}
	message, err := client.Messages.New(context.TODO(), createParams(model, messages))
	if err != nil {
		return nil, herrors.New(herrors.APIRequest, err)
	}
	return message, nil
}

func createParams(model string, messages []anthropic.MessageParam) anthropic.MessageNewParams {
	return anthropic.MessageNewParams{
		Model:     model,
		MaxTokens: 2480,
		Messages:  messages,
	}
}
