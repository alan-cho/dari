package agent

import (
	"dari/harness/config"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

const (
	modelEnv     = "MODEL"
	maxTokensEnv = "MAX_TOKENS"
)

type Client struct {
	api       anthropic.Client
	model     string
	maxTokens int64
}

func CreateClient() (Client, error) {
	model, err := config.Require(modelEnv)
	if err != nil {
		return Client{}, err
	}
	maxTokens, err := config.RequirePositiveInt(maxTokensEnv)
	if err != nil {
		return Client{}, err
	}
	return Client{api: anthropic.NewClient(), model: model, maxTokens: maxTokens}, nil
}
