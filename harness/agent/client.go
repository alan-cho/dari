package agent

import (
	"github.com/anthropics/anthropic-sdk-go/v2"
)

func CreateClient() anthropic.Client {
	return anthropic.NewClient()
}
