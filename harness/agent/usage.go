package agent

import (
	"dari/harness/api"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

func usageFrom(usage anthropic.Usage) api.Usage {
	prompt := usage.InputTokens + usage.CacheReadInputTokens + usage.CacheCreationInputTokens
	return api.Usage{Prompt: prompt, Output: usage.OutputTokens}
}
