package parser

import (
	"github.com/anthropics/anthropic-sdk-go/v2"	
)


// TODO: context management for management requests
func CreateMessage(client anthropic.Client, input string) (message anthropic.Message) {
	message = client.Messages.New(context.TODO(), createParams())
	return message
}

func createParams(model string, messages []anthropic.MessageParam) anthropic.MessageNewParams {
	return anthropic.MessageNewParams {
		Model: model,
		MaxTokens: 2480,
		Messages: messages,
	}
}
