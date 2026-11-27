package agent

import (
	"os"

	herrors "dari/harness/errors"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

const modelEnv = "MODEL"

type Client struct {
	api   anthropic.Client
	model string
}

func CreateClient() (Client, error) {
	model := os.Getenv(modelEnv)
	if model == "" {
		return Client{}, herrors.New(herrors.MissingModel, nil)
	}
	return Client{api: anthropic.NewClient(), model: model}, nil
}
