package agent

import "github.com/anthropics/anthropic-sdk-go/v2"

type Usage struct {
	Prompt int64
	Output int64
}

func Total(usage Usage) int64 {
	return usage.Prompt + usage.Output
}

func MessageInput(turn Usage, previousTotal int64) int64 {
	return (Total(turn) - previousTotal) - turn.Output
}

func usageFrom(usage anthropic.Usage) Usage {
	prompt := usage.InputTokens + usage.CacheReadInputTokens + usage.CacheCreationInputTokens
	return Usage{Prompt: prompt, Output: usage.OutputTokens}
}
