package main

import (
	"dari/harness/agent/parser"
	"os"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

func main() {
	key := os.Getenv("LLM_API_KEY")
	model := os.Getenv("MODEL")
	client := anthropic.NewClient(
		option.WithAPIKey(key),
	)

	scanner := parser.InitScanner()

	var transcript []string

	for {
		message := agent.parser.CreateMessage()
		// add to transcript through conversation.go
		// send to model through agent client
		// determine next steps based on model response
	}
}
