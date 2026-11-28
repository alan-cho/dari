package api

import herrors "dari/harness/errors"

type Role string

const (
	User      Role = "user"
	Assistant Role = "assistant"
)

type Message struct {
	Role Role   `json:"role"`
	Text string `json:"text"`
}

type SendRequest struct {
	Messages []Message `json:"messages"`
	Input    string    `json:"input"`
}

type Usage struct {
	Prompt int64 `json:"prompt"`
	Output int64 `json:"output"`
}

type SendResponse struct {
	Messages []Message `json:"messages"`
	Reply    string    `json:"reply"`
	Usage    Usage     `json:"usage"`
}

type ErrorResponse struct {
	Code        herrors.Code `json:"code"`
	Description string       `json:"description"`
}
