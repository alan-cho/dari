package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

// maxReadBytes caps how much of a file is returned, so one large file
// can't flood the context window.
const maxReadBytes = 100_000

type ReadFile struct{}

func (ReadFile) Name() string { return "read_file" }

func (ReadFile) Definition() anthropic.ToolParam {
	return anthropic.ToolParam{
		Name:        "read_file",
		Description: anthropic.String("Read a text file from the local filesystem and return its contents. Use a path relative to the current working directory or an absolute path."),
		InputSchema: anthropic.ToolInputSchemaParam{
			Properties: map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path of the file to read",
				},
			},
			Required: []string{"path"},
		},
	}
}

func (ReadFile) Run(_ context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	if in.Path == "" {
		return "", fmt.Errorf("path is required")
	}

	data, err := os.ReadFile(in.Path)
	if err != nil {
		return "", err
	}
	if len(data) > maxReadBytes {
		return string(data[:maxReadBytes]) + fmt.Sprintf("\n\n[truncated: file is %d bytes, showing first %d]", len(data), maxReadBytes), nil
	}
	return string(data), nil
}
